<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Download } from 'lucide-svelte';
	import { Spinner } from '$lib/components/ui/spinner/index.js';
	import { toast } from 'svelte-sonner';
	import { tabsMap } from '$lib/state.svelte';
	import { ExportRows, ExportTable } from '$lib/wailsjs/go/app/Export';
	import type { model } from '$lib/wailsjs/go/models';

	type ExportFormat = 'csv' | 'json';

	let {
		tabID,
		// Names the suggested file. Go appends a timestamp and the extension.
		fileName,
		// Only the table view can export beyond what is on screen: it pages
		// server-side, so the rows it holds are one page of a result set the
		// database still has in full. Leaving tableName empty -- which the
		// editor does -- offers the in-memory rows alone.
		tableName = '',
		select = '',
		where = '',
		orderBy = '',
		groupBy = ''
	}: {
		tabID: number;
		fileName: string;
		tableName?: string;
		select?: string;
		where?: string;
		orderBy?: string;
		groupBy?: string;
	} = $props();

	const isPaged = $derived(tableName !== '');

	let running = $state(false);

	function report(result: model.ExportResult) {
		// A dismissed save dialog is the user changing their mind, not a failure.
		if (result.canceled) {
			return;
		}
		if (!result.ok) {
			toast.error('Export failed', { description: result.message });
			return;
		}
		toast.success(`Exported ${result.rows} ${result.rows === 1 ? 'row' : 'rows'}`, {
			description: result.path
		});
	}

	function fail(error: unknown) {
		toast.error('Export failed', { description: String(error) });
	}

	function exportVisible(format: ExportFormat) {
		const tab = tabsMap.get(tabID);
		const columns = tab?.columns ?? [];
		if (columns.length === 0) {
			toast.error('Nothing to export', {
				description: 'Run a query first.'
			});
			return;
		}

		running = true;
		ExportRows(format, fileName, columns, tab?.rows ?? [])
			.then(report)
			.catch(fail)
			.finally(() => (running = false));
	}

	function exportEverything(format: ExportFormat) {
		running = true;
		ExportTable(tabID, format, fileName, tableName, select, where, orderBy, groupBy)
			.then(report)
			.catch(fail)
			.finally(() => (running = false));
	}
</script>

<DropdownMenu.Root>
	<DropdownMenu.Trigger disabled={running}>
		{#snippet child({ props })}
			<Button {...props} variant="outline" size="sm" class="h-8" disabled={running}>
				{#if running}
					<Spinner data-icon="inline-start" />
				{:else}
					<Download data-icon="inline-start" />
				{/if}
				Export
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content align="start" class="w-52">
		<DropdownMenu.Group>
			{#if isPaged}
				<DropdownMenu.GroupHeading>This page</DropdownMenu.GroupHeading>
			{/if}
			<DropdownMenu.Item onSelect={() => exportVisible('csv')}>CSV</DropdownMenu.Item>
			<DropdownMenu.Item onSelect={() => exportVisible('json')}>JSON</DropdownMenu.Item>
		</DropdownMenu.Group>
		{#if isPaged}
			<DropdownMenu.Separator />
			<DropdownMenu.Group>
				<DropdownMenu.GroupHeading>Every matching row</DropdownMenu.GroupHeading>
				<DropdownMenu.Item onSelect={() => exportEverything('csv')}>CSV</DropdownMenu.Item>
				<DropdownMenu.Item onSelect={() => exportEverything('json')}>JSON</DropdownMenu.Item>
			</DropdownMenu.Group>
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>
