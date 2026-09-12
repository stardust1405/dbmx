<script lang="ts">
	import SqlClauseInput from './sql-clause-input.svelte';
	import type { ColumnInfo } from './sql-clause';
	import { GetTableColumnsMeta } from '$lib/wailsjs/go/app/Connections';

	let {
		tabID,
		tableName,
		columnNames = [],
		select = $bindable(''),
		where = $bindable(''),
		orderBy = $bindable(''),
		groupBy = $bindable(''),
		onrun
	}: {
		tabID: number;
		tableName: string;
		columnNames?: string[];
		select?: string;
		where?: string;
		orderBy?: string;
		groupBy?: string;
		onrun?: () => void;
	} = $props();

	let meta = $state<ColumnInfo[]>([]);

	// Names alone still drive completion if the metadata call fails or hasn't landed.
	// The list is normalised rather than trusted: it arrives from Go, where a nil
	// slice crosses the boundary as null, which no prop default fills in.
	const columns = $derived<ColumnInfo[]>(
		meta.length > 0 ? meta : (columnNames ?? []).map((name) => ({ name }))
	);

	$effect(() => {
		const id = tabID;
		const table = tableName;
		if (!id || !table) return;

		// One <ClauseBar> serves every tab, so drop the previous table's metadata
		// immediately rather than suggesting its columns until the fetch lands.
		meta = [];

		let cancelled = false;
		GetTableColumnsMeta(id, table)
			.then((columnsMeta) => {
				if (cancelled) return;
				meta = (columnsMeta ?? []).map((column) => ({
					name: column.name,
					dataType: column.dataType,
					isPrimaryKey: column.isPrimaryKey,
					isNullable: column.isNullable,
					enumValues: column.enumValues ?? []
				}));
			})
			.catch(() => {
				// Completion degrades to plain column names; the grid still loads.
				if (!cancelled) meta = [];
			});

		return () => {
			cancelled = true;
		};
	});
</script>

<div class="clause-bar">
	<div class="clause-slot slot-tl">
		<SqlClauseInput
			id="clause-select-{tabID}"
			clause="select"
			label="select"
			placeholder="*"
			columns={columns}
			bind:value={select}
			{onrun}
		/>
	</div>
	<div class="clause-slot slot-tr">
		<SqlClauseInput
			id="clause-where-{tabID}"
			clause="where"
			label="where"
			placeholder="all rows"
			columns={columns}
			bind:value={where}
			{onrun}
		/>
	</div>
	<div class="clause-slot slot-bl">
		<SqlClauseInput
			id="clause-order-{tabID}"
			clause="orderBy"
			label="order"
			placeholder="table order"
			columns={columns}
			bind:value={orderBy}
			{onrun}
		/>
	</div>
	<div class="clause-slot slot-br">
		<SqlClauseInput
			id="clause-group-{tabID}"
			clause="groupBy"
			label="group"
			placeholder="no grouping"
			columns={columns}
			bind:value={groupBy}
			{onrun}
		/>
	</div>
</div>

<style>
	.clause-bar {
		display: grid;
		grid-template-columns: 1fr 1fr;
		border: 1px solid hsl(var(--border));
		border-radius: calc(var(--radius) - 2px);
		background: hsl(var(--background));
		overflow: hidden;
	}

	.clause-slot {
		min-width: 0;
	}

	/* Hairlines between cells only — no boxes around each input. */
	.slot-tl,
	.slot-tr {
		border-bottom: 1px solid hsl(var(--border));
	}

	.slot-tl,
	.slot-bl {
		border-right: 1px solid hsl(var(--border));
	}

	@media (max-width: 720px) {
		.clause-bar {
			grid-template-columns: 1fr;
		}

		.slot-tl,
		.slot-bl {
			border-right: 0;
		}

		.slot-bl {
			border-bottom: 1px solid hsl(var(--border));
		}
	}
</style>
