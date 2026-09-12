<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Switch } from '$lib/components/ui/switch/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { cn } from '$lib/utils.js';
	import type { model } from '$lib/wailsjs/go/models';
	import {
		CreateIndex,
		GetIndexDefinition,
		GetSchemaEditorOptions,
		GetTableColumnsMeta,
		UpdateIndex
	} from '$lib/wailsjs/go/app/Connections';
	import SchemaExpressionList from './schema-expression-list.svelte';

	let {
		tabID,
		tableName,
		open = $bindable(false),
		/** The index being edited; empty opens the sheet as "add index". */
		indexName = '',
		onSaved
	}: {
		tabID: number;
		tableName: string;
		open?: boolean;
		indexName?: string;
		onSaved?: () => void;
	} = $props();

	const isEdit = $derived(indexName !== '');

	function blankIndex(): model.IndexDefinition {
		return {
			name: '',
			method: 'btree',
			isUnique: false,
			columns: [''],
			include: [],
			where: '',
			comment: '',
			isConstraint: false,
			isPrimary: false
		};
	}

	let form = $state<model.IndexDefinition>(blankIndex());
	let methods = $state<string[]>(['btree']);
	let columnNames = $state<string[]>([]);
	let loading = $state(false);
	let saving = $state(false);
	let showErrors = $state(false);

	$effect(() => {
		if (open) {
			load();
		}
	});

	async function load() {
		loading = true;
		saving = false;
		showErrors = false;
		form = blankIndex();

		try {
			const [options, meta, existing] = await Promise.all([
				GetSchemaEditorOptions(tabID),
				GetTableColumnsMeta(tabID, tableName),
				indexName ? GetIndexDefinition(tabID, tableName, indexName) : Promise.resolve(null)
			]);

			methods = options.indexMethods?.length ? options.indexMethods : ['btree'];
			columnNames = (meta ?? []).map((column) => column.name);

			if (existing) {
				form = {
					...blankIndex(),
					...existing,
					columns: existing.columns?.length ? [...existing.columns] : [''],
					include: existing.include ? [...existing.include] : []
				};
			} else if (!methods.includes(form.method)) {
				form.method = methods[0];
			}
		} catch (error) {
			toast.error('Failed to open the index editor', { description: String(error) });
			open = false;
		} finally {
			loading = false;
		}
	}

	const nameError = $derived(form.name.trim() === '' ? 'An index name is required.' : '');
	const columnsError = $derived(
		form.columns.filter((column) => column.trim() !== '').length === 0
			? 'An index needs at least one column.'
			: ''
	);
	const hasErrors = $derived(nameError !== '' || columnsError !== '');

	// Postgres can rename an index but not redefine one, so any other change is a drop
	// and a create. Saying so up front is fairer than a surprise on a large table.
	const willRebuild = $derived(isEdit && !form.isConstraint);

	function save() {
		showErrors = true;
		if (hasErrors) return;

		const payload: model.IndexDefinition = {
			...form,
			name: form.name.trim(),
			columns: form.columns.map((column) => column.trim()).filter(Boolean),
			include: form.include.map((column) => column.trim()).filter(Boolean),
			where: form.where.trim(),
			comment: form.comment.trim()
		};

		saving = true;
		const request = isEdit
			? UpdateIndex(tabID, tableName, indexName, payload)
			: CreateIndex(tabID, tableName, payload);

		request
			.then(() => {
				toast.success(isEdit ? 'Index updated' : 'Index created', {
					description: `${payload.name} on ${tableName}.`
				});
				open = false;
				onSaved?.();
			})
			.catch((error) => {
				toast.error(isEdit ? 'Failed to update the index' : 'Failed to create the index', {
					description: String(error)
				});
				saving = false;
			});
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="flex w-full flex-col gap-0 p-0 sm:max-w-lg">
		<Sheet.Header class="border-b px-6 py-4 pr-12 text-left">
			<Sheet.Title>
				{isEdit ? `Edit ${indexName}` : `Add an index to ${tableName}`}
			</Sheet.Title>
			<Sheet.Description>
				{form.isConstraint
					? 'This index backs a constraint, so it is edited under Rules instead.'
					: willRebuild
						? 'Anything beyond the name is applied by rebuilding the index, in one transaction.'
						: 'Key columns may be expressions, and an index can cover a subset of the rows.'}
			</Sheet.Description>
		</Sheet.Header>

		<div class="flex-1 overflow-y-auto px-6 py-5">
			{#if loading}
				<div class="flex flex-col gap-6">
					{#each [1, 2, 3, 4] as placeholder (placeholder)}
						<div class="flex flex-col gap-2">
							<Skeleton class="h-4 w-32" />
							<Skeleton class="h-9 w-full" />
						</div>
					{/each}
				</div>
			{:else}
				<Field.FieldGroup>
					<Field.Field data-invalid={showErrors && nameError !== ''}>
						<Field.FieldLabel for="index-name">Name</Field.FieldLabel>
						<Input
							id="index-name"
							bind:value={form.name}
							placeholder={`${tableName}_column_idx`}
							disabled={form.isConstraint}
							aria-invalid={showErrors && nameError !== ''}
						/>
						{#if showErrors && nameError}
							<Field.FieldError>{nameError}</Field.FieldError>
						{/if}
					</Field.Field>

					<Field.Field>
						<Field.FieldLabel for="index-method">Method</Field.FieldLabel>
						<Select.Root
							type="single"
							value={form.method}
							onValueChange={(next) => (form.method = next)}
						>
							<Select.Trigger id="index-method" class="w-full" disabled={form.isConstraint}>
								{form.method}
							</Select.Trigger>
							<Select.Content>
								{#each methods as method (method)}
									<Select.Item value={method}>{method}</Select.Item>
								{/each}
							</Select.Content>
						</Select.Root>
					</Field.Field>

					<Field.Field orientation="horizontal">
						<Field.FieldContent>
							<Field.FieldLabel for="index-unique">Unique</Field.FieldLabel>
							<Field.FieldDescription>
								Rejects duplicate values across the key columns.
							</Field.FieldDescription>
						</Field.FieldContent>
						<Switch
							id="index-unique"
							bind:checked={form.isUnique}
							disabled={form.isConstraint}
						/>
					</Field.Field>

					<Field.Field data-invalid={showErrors && columnsError !== ''}>
						<Field.FieldLabel>Key columns</Field.FieldLabel>
						<SchemaExpressionList
							bind:entries={form.columns}
							{columnNames}
							idPrefix="index-column"
							addLabel="Add key column"
						/>
						{#if showErrors && columnsError}
							<Field.FieldError>{columnsError}</Field.FieldError>
						{:else}
							<Field.FieldDescription>
								An entry can carry an expression and an ordering — <code>lower(email)</code>,
								<code>created_at DESC NULLS LAST</code>.
							</Field.FieldDescription>
						{/if}
					</Field.Field>

					<Field.Field>
						<Field.FieldLabel>Included columns</Field.FieldLabel>
						<SchemaExpressionList
							bind:entries={form.include}
							{columnNames}
							idPrefix="index-include"
							addLabel="Add included column"
						/>
						<Field.FieldDescription>
							Stored in the index but not part of the key, so a query can be answered without
							reading the table.
						</Field.FieldDescription>
					</Field.Field>

					<Field.Field>
						<Field.FieldLabel for="index-where">Only index rows where</Field.FieldLabel>
						<Textarea
							id="index-where"
							rows={2}
							class="font-mono text-xs"
							bind:value={form.where}
							placeholder="deleted_at IS NULL"
						/>
						<Field.FieldDescription>
							Leave empty to index every row.
						</Field.FieldDescription>
					</Field.Field>

					<Field.Field>
						<Field.FieldLabel for="index-comment">Comment</Field.FieldLabel>
						<Textarea id="index-comment" rows={2} bind:value={form.comment} />
					</Field.Field>
				</Field.FieldGroup>
			{/if}
		</div>

		<Sheet.Footer class="flex-row justify-end gap-2 border-t px-6 py-4 sm:space-x-0">
			<Button variant="outline" onclick={() => (open = false)} disabled={saving}>Cancel</Button>
			<Button
				onclick={save}
				disabled={loading || saving || form.isConstraint}
				class="relative"
			>
				<span class={cn(saving && 'invisible')}>{isEdit ? 'Save changes' : 'Create index'}</span>
				{#if saving}
					<span class="absolute inset-0 flex items-center justify-center">
						<Spinner />
					</span>
				{/if}
			</Button>
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
