<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import {
		DropColumn,
		DropConstraint,
		DropIndex,
		GetTableInfo
	} from '$lib/wailsjs/go/app/Connections.js';
	import { toast } from 'svelte-sonner';
	import type { ColumnDef, RowData } from '@tanstack/table-core';
	import DataTableManage from './data-table-manage.svelte';
	import TableDDL from './table_ddl.svelte';
	import SchemaColumnSheet from './schema-column-sheet.svelte';
	import SchemaIndexSheet from './schema-index-sheet.svelte';
	import SchemaConstraintSheet from './schema-constraint-sheet.svelte';

	// Accept pool id and table name as props
	let { tabID, tabName = '' } = $props();
	let selectedView = $state('structure');

	type Grid = { columns: ColumnDef<RowData, unknown>[]; rows: RowData[] };
	const emptyGrid: Grid = { columns: [], rows: [] };

	let structure = $state<Grid>(emptyGrid);
	let indexes = $state<Grid>(emptyGrid);
	let rules = $state<Grid>(emptyGrid);

	let columnSheetOpen = $state(false);
	let editingColumn = $state('');
	let indexSheetOpen = $state(false);
	let editingIndex = $state('');
	let constraintSheetOpen = $state(false);
	let editingConstraint = $state('');

	// One confirmation dialog serves all three views: the wording and the call differ,
	// the shape of the decision -- drop this named thing, optionally cascading -- does not.
	type Pending = { kind: 'column' | 'index' | 'constraint'; name: string };
	let pendingDrop = $state<Pending | null>(null);
	let cascade = $state(false);
	let dropping = $state(false);

	const DROP_COPY = {
		column: {
			title: 'Drop column',
			body: 'The column and every value in it are removed.',
			cascade: 'Also drop whatever depends on it — views, indexes, foreign keys.'
		},
		index: {
			title: 'Drop index',
			body: 'Queries relying on it fall back to scanning the table.',
			cascade: 'Also drop whatever depends on it.'
		},
		constraint: {
			title: 'Drop constraint',
			body: 'The rule stops being enforced. Rows already in the table are untouched.',
			cascade: 'Also drop whatever depends on it — foreign keys pointing at this key.'
		}
	};

	onMount(() => {
		getTableInfo();
	});

	// Section -> tanstack grid. The backend hands every section back in the same
	// { columns, rows } shape, so one mapper covers structure, indexes and rules.
	function toGrid(section: { columns?: string[]; rows?: { column: string; value: string }[][] }): Grid {
		const columns = (section.columns ?? []).map((column) => ({
			accessorKey: column,
			id: column, // Unique ID guaranteed!
			header: column
		}));

		const rows = (section.rows ?? []).map((row) => {
			const cell: Record<string, any> = {};
			for (const resultCell of row) {
				if (resultCell.column && resultCell.value) {
					cell[resultCell.column] = resultCell.value;
				}
			}
			return cell;
		});

		return { columns, rows } as Grid;
	}

	function getTableInfo() {
		if (tabName == '') {
			toast.error('Please select a table to execute the query', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}

		GetTableInfo(tabID, tabName)
			.then((response) => {
				// Assigned rather than appended, so a refresh after an edit replaces the
				// previous read instead of stacking a second copy of every row onto it.
				structure = response.structure ? toGrid(response.structure) : emptyGrid;
				indexes = response.indexes ? toGrid(response.indexes) : emptyGrid;
				rules = response.rules ? toGrid(response.rules) : emptyGrid;
			})
			.catch((error) => {
				toast.error('Failed to get table info', {
					description: String(error),
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});
	}

	function openColumnSheet(name: string) {
		editingColumn = name;
		columnSheetOpen = true;
	}

	function openIndexSheet(name: string) {
		editingIndex = name;
		indexSheetOpen = true;
	}

	function openConstraintSheet(name: string) {
		editingConstraint = name;
		constraintSheetOpen = true;
	}

	function confirmDrop(kind: Pending['kind'], name: string) {
		pendingDrop = { kind, name };
		cascade = false;
	}

	function drop() {
		if (!pendingDrop) return;
		const { kind, name } = pendingDrop;

		dropping = true;
		const request =
			kind === 'column'
				? DropColumn(tabID, tabName, name, cascade)
				: kind === 'index'
					? DropIndex(tabID, tabName, name, cascade)
					: DropConstraint(tabID, tabName, name, cascade);

		request
			.then(() => {
				toast.success(`${DROP_COPY[kind].title.replace('Drop', 'Dropped')} ${name}`, {
					description: `Removed from ${tabName}.`
				});
				pendingDrop = null;
				getTableInfo();
			})
			.catch((error) => {
				toast.error(`Failed to drop the ${kind}`, { description: String(error) });
			})
			.finally(() => {
				dropping = false;
			});
	}

	// An index postgres created for a constraint is owned by that constraint: it is
	// edited and dropped under Rules, not here. A standalone unique index is not --
	// that one this view owns.
	function indexIsLocked(row: Record<string, any>): boolean {
		return row['is_constraint'] === 'true';
	}
</script>

<SchemaColumnSheet
	{tabID}
	tableName={tabName}
	bind:open={columnSheetOpen}
	columnName={editingColumn}
	onSaved={getTableInfo}
/>
<SchemaIndexSheet
	{tabID}
	tableName={tabName}
	bind:open={indexSheetOpen}
	indexName={editingIndex}
	onSaved={getTableInfo}
/>
<SchemaConstraintSheet
	{tabID}
	tableName={tabName}
	bind:open={constraintSheetOpen}
	constraintName={editingConstraint}
	onSaved={getTableInfo}
/>

<Dialog.Root
	open={pendingDrop !== null}
	onOpenChange={(next) => {
		if (!next) pendingDrop = null;
	}}
>
	<Dialog.Content class="sm:max-w-md">
		{#if pendingDrop}
			{@const copy = DROP_COPY[pendingDrop.kind]}
			<Dialog.Header>
				<Dialog.Title>{copy.title} {pendingDrop.name}?</Dialog.Title>
				<Dialog.Description>
					{copy.body} This cannot be undone.
				</Dialog.Description>
			</Dialog.Header>
			<label class="flex items-start gap-2 text-sm">
				<Checkbox
					class="mt-0.5"
					checked={cascade}
					onCheckedChange={(checked: boolean) => (cascade = checked)}
				/>
				<span class="text-muted-foreground">{copy.cascade}</span>
			</label>
			<Dialog.Footer>
				<Button variant="outline" onclick={() => (pendingDrop = null)} disabled={dropping}>
					Cancel
				</Button>
				<Button variant="destructive" onclick={drop} disabled={dropping}>
					{dropping ? 'Dropping...' : 'Drop'}
				</Button>
			</Dialog.Footer>
		{/if}
	</Dialog.Content>
</Dialog.Root>

<div class="flex h-full flex-row overflow-hidden">
	<!-- Tab Management -->
	<div class="m-2 flex w-20 flex-col overflow-hidden">
		<div class="mb-0.5 flex flex-1 items-center justify-center">
			<Button
				variant={selectedView === 'structure' ? 'default' : 'outline'}
				class="flex h-full w-full flex-col items-center justify-center"
				onclick={() => (selectedView = 'structure')}
			>
				{#each 'structure'.split('') as letter}
					<span class="flex-shrink-0 text-[16px] leading-none">{letter}</span>
				{/each}
			</Button>
		</div>
		<div class="mt-0.5 flex flex-1 items-center justify-center">
			<Button
				variant={selectedView === 'indexes' ? 'default' : 'outline'}
				class="flex h-full w-full flex-col items-center justify-center"
				onclick={() => (selectedView = 'indexes')}
			>
				{#each 'indexes'.split('') as letter}
					<span class="flex-shrink-0 text-[16px] leading-none">{letter}</span>
				{/each}
			</Button>
		</div>
		<div class="mt-0.5 flex flex-1 items-center justify-center">
			<Button
				variant={selectedView === 'rules' ? 'default' : 'outline'}
				class="flex h-full w-full flex-col items-center justify-center"
				onclick={() => (selectedView = 'rules')}
			>
				{#each 'rules'.split('') as letter}
					<span class="flex-shrink-0 text-[16px] leading-none">{letter}</span>
				{/each}
			</Button>
		</div>
		<div class="mt-0.5 flex flex-1 items-center justify-center">
			<Button
				variant={selectedView === 'sql' ? 'default' : 'outline'}
				class="flex h-full w-full flex-col items-center justify-center"
				onclick={() => (selectedView = 'sql')}
			>
				{#each 'sql'.split('') as letter}
					<span class="flex-shrink-0 text-[16px] leading-none">{letter}</span>
				{/each}
			</Button>
		</div>
	</div>

	<!-- Tab Content -->
	<div class="flex h-full min-h-0 flex-[15] overflow-hidden">
		{#if selectedView === 'structure'}
			<div class="flex h-full min-h-0 flex-1 flex-col overflow-hidden pr-2 pb-2">
				<div class="flex h-1 min-h-0 flex-1 flex-col overflow-hidden">
					{#if structure.columns.length > 0}
						<DataTableManage
							rows={structure.rows}
							columns={structure.columns}
							nameKey="column_name"
							addLabel="Add column"
							onAdd={() => openColumnSheet('')}
							onEdit={openColumnSheet}
							onDelete={(name) => confirmDrop('column', name)}
						/>
					{:else}
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
					{/if}
				</div>
			</div>
		{:else if selectedView === 'indexes'}
			<div class="flex h-full min-h-0 flex-1 flex-col overflow-hidden pr-2 pb-2">
				<div class="flex h-1 min-h-0 flex-1 flex-col overflow-hidden">
					{#if indexes.columns.length > 0}
						<DataTableManage
							rows={indexes.rows}
							columns={indexes.columns}
							nameKey="index_name"
							addLabel="Add index"
							onAdd={() => openIndexSheet('')}
							onEdit={openIndexSheet}
							onDelete={(name) => confirmDrop('index', name)}
							isLocked={indexIsLocked}
						/>
					{:else}
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
					{/if}
				</div>
			</div>
		{:else if selectedView === 'rules'}
			<div class="flex h-full min-h-0 flex-1 flex-col overflow-hidden pr-2 pb-2">
				<div class="flex h-1 min-h-0 flex-1 flex-col overflow-hidden">
					{#if rules.columns.length > 0}
						<DataTableManage
							rows={rules.rows}
							columns={rules.columns}
							nameKey="constraint_name"
							addLabel="Add constraint"
							onAdd={() => openConstraintSheet('')}
							onEdit={openConstraintSheet}
							onDelete={(name) => confirmDrop('constraint', name)}
						/>
					{:else}
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
						<Skeleton class="my-3 h-[40px] w-full" />
					{/if}
				</div>
			</div>
		{:else if selectedView === 'sql'}
			<div class="flex h-full min-h-0 flex-1 flex-col overflow-hidden pr-2 pb-2">
				<!-- Keyed on the table so switching tables tears the editor down and
				     refetches, rather than leaving the previous table's definition on screen. -->
				{#key tabName}
					<TableDDL {tabID} tableName={tabName} />
				{/key}
			</div>
		{/if}
	</div>
</div>
