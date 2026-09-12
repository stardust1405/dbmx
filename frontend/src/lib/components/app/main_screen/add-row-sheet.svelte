<script lang="ts">
	import * as Sheet from '$lib/components/ui/sheet/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { SvelteMap } from 'svelte/reactivity';
	import { cn } from '$lib/utils.js';
	import type { model } from '$lib/wailsjs/go/models';
	import { GetTableColumnsMeta, InsertRow } from '$lib/wailsjs/go/app/Connections';

	let {
		tabID,
		tableName,
		open = $bindable(false),
		onInserted
	}: {
		tabID: number;
		tableName: string;
		open?: boolean;
		onInserted?: () => void;
	} = $props();

	// A column is either left to its database default, set to SQL NULL, set to the
	// empty string, or given a literal parsed by the column's own postgres type.
	// 'empty' is deliberately distinct from a blank 'value' box: postgres NOT NULL
	// permits '', so a blank box would otherwise silently insert one on a column
	// the form has just labelled "required".
	type FieldMode = 'default' | 'null' | 'empty' | 'value';
	type FieldState = { mode: FieldMode; value: string };

	let columnsMeta = $state<model.ColumnMeta[]>([]);
	let fields = $state(new SvelteMap<string, FieldState>());
	let loading = $state(false);
	let saving = $state(false);
	let showErrors = $state(false);

	const editableColumns = $derived(columnsMeta.filter((column) => !column.isReadOnly));

	$effect(() => {
		if (open) {
			loadColumns();
		}
	});

	function initialMode(column: model.ColumnMeta): FieldMode {
		if (column.hasDefault) return 'default';
		if (column.isNullable) return 'null';
		return 'value';
	}

	function loadColumns() {
		loading = true;
		saving = false;
		showErrors = false;
		GetTableColumnsMeta(tabID, tableName)
			.then((meta) => {
				columnsMeta = meta ?? [];
				const next = new SvelteMap<string, FieldState>();
				for (const column of columnsMeta) {
					next.set(column.name, { mode: initialMode(column), value: '' });
				}
				fields = next;
			})
			.catch((error) => {
				toast.error('Failed to load table columns', { description: String(error) });
				open = false;
			})
			.finally(() => {
				loading = false;
			});
	}

	function setMode(name: string, mode: FieldMode) {
		const field = fields.get(name);
		if (!field) return;
		fields.set(name, { mode, value: mode === 'value' ? field.value : '' });
	}

	function setValue(name: string, value: string) {
		fields.set(name, { mode: 'value', value });
	}

	// pg_type.typcategory drives the widget, so every postgres type is handled without
	// enumerating type names: B boolean, E enum, plus a text box for the long ones.
	function widgetFor(column: model.ColumnMeta): 'boolean' | 'enum' | 'textarea' | 'input' {
		if (column.category === 'B') return 'boolean';
		if (column.enumValues && column.enumValues.length > 0) return 'enum';
		if (['text', 'json', 'jsonb', 'xml'].includes(column.castType.toLowerCase())) return 'textarea';
		return 'input';
	}

	function placeholderFor(column: model.ColumnMeta): string {
		if (column.category === 'A') return '{value1,value2}';
		switch (column.castType.toLowerCase()) {
			case 'json':
			case 'jsonb':
				return '{ "key": "value" }';
			case 'uuid':
				return '00000000-0000-0000-0000-000000000000';
			case 'bytea':
				return '\\x48656c6c6f';
			case 'date':
				return 'YYYY-MM-DD';
			case 'time without time zone':
			case 'time with time zone':
				return 'HH:MM:SS';
			case 'timestamp without time zone':
			case 'timestamp with time zone':
				return 'YYYY-MM-DD HH:MM:SS';
			case 'interval':
				return '1 day 02:03:04';
			case 'macaddr':
			case 'macaddr8':
				return '08:00:2b:01:02:03';
			case 'money':
				return '1000.00';
		}
		switch (column.category) {
			case 'N':
				return '0';
			case 'I':
				return '192.168.0.1';
			case 'R':
				return '[1,10)';
			case 'G':
				return '(1,2)';
			case 'V':
				return '1010';
			default:
				return '';
		}
	}

	// Postgres accepts '' in a NOT NULL string column, so a blank box is ambiguous
	// rather than invalid. Treat it as unfinished for every type and point at the
	// explicit alternatives; everything else is left to the real postgres error.
	function errorFor(column: model.ColumnMeta): string {
		const field = fields.get(column.name);
		if (!field || field.mode !== 'value' || field.value !== '') return '';

		const alternatives: string[] = [];
		if (column.hasDefault) alternatives.push('Default');
		if (column.isNullable) alternatives.push('NULL');
		if (allowsEmptyString(column)) alternatives.push('Empty');

		return alternatives.length > 0
			? `Enter a value, or tick ${alternatives.join(' / ')}.`
			: 'A value is required.';
	}

	// Only the string family round-trips '' as a value; '' is invalid input for
	// numbers, dates, json, uuid and the rest, so postgres would reject it anyway.
	function allowsEmptyString(column: model.ColumnMeta): boolean {
		return column.category === 'S';
	}

	// What the empty input shows: the standing non-value choice, or a format hint.
	function hintFor(column: model.ColumnMeta): string {
		switch (fields.get(column.name)?.mode) {
			case 'default':
				return `DEFAULT ${column.defaultValue}`;
			case 'null':
				return 'NULL';
			case 'empty':
				return "''";
			default:
				return placeholderFor(column);
		}
	}

	const invalidColumns = $derived(editableColumns.filter((column) => errorFor(column) !== ''));

	function insert() {
		showErrors = true;
		if (invalidColumns.length > 0) {
			toast.error('Cannot insert row', {
				description: `${invalidColumns.length} column(s) still need a value.`
			});
			return;
		}

		const payload: model.InsertValue[] = [];
		for (const column of editableColumns) {
			const field = fields.get(column.name);
			if (!field || field.mode === 'default') continue;
			payload.push({
				columnName: column.name,
				value: field.mode === 'null' ? undefined : field.mode === 'empty' ? '' : field.value
			});
		}

		saving = true;
		InsertRow(tabID, tableName, payload)
			.then(() => {
				toast.success('Row inserted', { description: `A new row was added to ${tableName}.` });
				open = false;
				onInserted?.();
			})
			.catch((error) => {
				toast.error('Failed to insert row', { description: String(error) });
				// Only cleared on failure. On success the sheet is already sliding out, and
				// restoring the label mid-animation just flickers; opening resets it.
				saving = false;
			});
	}
</script>

<Sheet.Root bind:open>
	<Sheet.Content side="right" class="flex w-full flex-col gap-0 p-0 sm:max-w-lg">
		<Sheet.Header class="border-b px-6 py-4 pr-12 text-left">
			<Sheet.Title>Add row to {tableName}</Sheet.Title>
			<Sheet.Description>
				Leave a column on its default, set it to NULL or an empty string, or type a value. Each
				value is parsed by that column's own postgres type.
			</Sheet.Description>
		</Sheet.Header>

		<div class="flex-1 overflow-y-auto px-6 py-5">
			{#if loading}
				<div class="flex flex-col gap-6">
					{#each [1, 2, 3, 4, 5] as placeholder (placeholder)}
						<div class="flex flex-col gap-2">
							<Skeleton class="h-4 w-40" />
							<Skeleton class="h-9 w-full" />
						</div>
					{/each}
				</div>
			{:else}
				<Field.FieldGroup>
					{#each columnsMeta as column (column.name)}
						{@const field = fields.get(column.name)}
						{@const error = showErrors ? errorFor(column) : ''}
						<Field.Field data-invalid={error !== ''}>
							<div class="flex flex-wrap items-center justify-between gap-2">
								<Field.FieldLabel for={`add-row-${column.name}`} class="flex items-center gap-2">
									{column.name}
									<span class="text-muted-foreground text-xs font-normal">{column.dataType}</span>
									{#if column.isPrimaryKey}
										<Badge variant="secondary" class="px-1.5 py-0 text-[10px]">PK</Badge>
									{/if}
									{#if !column.isNullable && !column.hasDefault && !column.isReadOnly}
										<Badge variant="outline" class="px-1.5 py-0 text-[10px]">required</Badge>
									{/if}
									{#if column.isReadOnly}
										<Badge variant="outline" class="px-1.5 py-0 text-[10px]">generated</Badge>
									{/if}
								</Field.FieldLabel>

								{#if !column.isReadOnly}
									<div class="text-muted-foreground flex items-center gap-3 text-xs">
										{#if column.hasDefault}
											<label class="flex items-center gap-1.5">
												<Checkbox
													checked={field?.mode === 'default'}
													onCheckedChange={(checked: boolean) =>
														setMode(column.name, checked ? 'default' : 'value')}
												/>
												Default
											</label>
										{/if}
										{#if column.isNullable}
											<label class="flex items-center gap-1.5">
												<Checkbox
													checked={field?.mode === 'null'}
													onCheckedChange={(checked: boolean) =>
														setMode(column.name, checked ? 'null' : 'value')}
												/>
												NULL
											</label>
										{/if}
										{#if allowsEmptyString(column)}
											<label class="flex items-center gap-1.5">
												<Checkbox
													checked={field?.mode === 'empty'}
													onCheckedChange={(checked: boolean) =>
														setMode(column.name, checked ? 'empty' : 'value')}
												/>
												Empty
											</label>
										{/if}
									</div>
								{/if}
							</div>

							{#if column.isReadOnly}
								<Input
									id={`add-row-${column.name}`}
									disabled
									value=""
									placeholder="Generated by the database"
								/>
							{:else if widgetFor(column) === 'boolean' || widgetFor(column) === 'enum'}
								{@const options =
									widgetFor(column) === 'boolean' ? ['true', 'false'] : column.enumValues}
								<Select.Root
									type="single"
									value={field?.mode === 'value' ? field.value : ''}
									onValueChange={(value) => setValue(column.name, value)}
								>
									<Select.Trigger
										id={`add-row-${column.name}`}
										class="w-full"
										aria-invalid={error !== ''}
									>
										{#if field?.mode === 'value' && field.value !== ''}
											{field.value}
										{:else if field?.mode === 'value'}
											<span class="text-muted-foreground">Select a value</span>
										{:else}
											<span class="text-muted-foreground">{hintFor(column)}</span>
										{/if}
									</Select.Trigger>
									<Select.Content>
										<Select.Group>
											{#each options as option (option)}
												<Select.Item value={option}>{option}</Select.Item>
											{/each}
										</Select.Group>
									</Select.Content>
								</Select.Root>
							{:else if widgetFor(column) === 'textarea'}
								<Textarea
									id={`add-row-${column.name}`}
									rows={3}
									class="font-mono text-xs"
									aria-invalid={error !== ''}
									value={field?.mode === 'value' ? field.value : ''}
									placeholder={hintFor(column)}
									oninput={(e: Event & { currentTarget: HTMLTextAreaElement }) =>
										setValue(column.name, e.currentTarget.value)}
								/>
							{:else}
								<Input
									id={`add-row-${column.name}`}
									inputmode={column.category === 'N' ? 'decimal' : undefined}
									aria-invalid={error !== ''}
									value={field?.mode === 'value' ? field.value : ''}
									placeholder={hintFor(column)}
									oninput={(e: Event & { currentTarget: HTMLInputElement }) =>
										setValue(column.name, e.currentTarget.value)}
								/>
							{/if}

							{#if error}
								<Field.FieldError>{error}</Field.FieldError>
							{:else if column.comment}
								<Field.FieldDescription>{column.comment}</Field.FieldDescription>
							{/if}
						</Field.Field>
					{/each}
				</Field.FieldGroup>
			{/if}
		</div>

		<Sheet.Footer class="flex-row justify-end gap-2 border-t px-6 py-4 sm:space-x-0">
			<Button variant="outline" onclick={() => (open = false)} disabled={saving}>Cancel</Button>
			<!-- The spinner is wrapped and absolutely positioned so the button's box never
			     changes: a direct-child <svg> would trip the button's has-[>svg]:px-3 rule
			     and, with transition-all, animate the padding on every toggle. -->
			<Button onclick={insert} disabled={loading || saving} class="relative">
				<span class={cn(saving && 'invisible')}>Insert row</span>
				{#if saving}
					<span class="absolute inset-0 flex items-center justify-center">
						<Spinner />
					</span>
				{/if}
			</Button>
		</Sheet.Footer>
	</Sheet.Content>
</Sheet.Root>
