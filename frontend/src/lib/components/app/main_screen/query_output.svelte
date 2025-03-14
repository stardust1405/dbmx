<script lang="ts">
	import {
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow
	} from '$lib/components/ui/table/index.js';

	export let outputs: Array<{
		id: string;
		columns: string[];
		rows: Record<string, string>[];
	}> = [];
</script>

<div class="h-full w-full overflow-auto p-2">
	{#if outputs.length > 0}
		{#each outputs as output}
			<div class="mb-4">
				<h3 class="mb-2 text-sm font-medium">Result #{output.id}</h3>
				<div class="rounded-md border">
					<Table>
						<TableHeader>
							<TableRow>
								{#each output.columns as column}
									<TableHead>{column}</TableHead>
								{/each}
							</TableRow>
						</TableHeader>
						<TableBody>
							{#each output.rows as row}
								<TableRow>
									{#each output.columns as column}
										<TableCell>{row[column]}</TableCell>
									{/each}
								</TableRow>
							{/each}
						</TableBody>
					</Table>
				</div>
			</div>
		{/each}
	{:else}
		<div class="flex h-full items-center justify-center">
			<p class="text-muted-foreground">Execute a query to see results</p>
		</div>
	{/if}
</div>
