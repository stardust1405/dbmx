<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import { Plus, X } from 'lucide-svelte';
	import SchemaCombobox, { type ComboboxOption } from './schema-combobox.svelte';

	// An ordered list of column references. Index keys are not always bare columns --
	// "lower(email)", "id DESC NULLS LAST", "tags jsonb_path_ops" are all valid -- so
	// each entry is a free-text combobox seeded with the table's column names rather
	// than a select that would lose anything postgres rendered back.
	let {
		entries = $bindable([]),
		columnNames = [],
		addLabel = 'Add column',
		placeholder = 'Column or expression',
		idPrefix = 'expression'
	}: {
		entries?: string[];
		columnNames?: string[];
		addLabel?: string;
		placeholder?: string;
		idPrefix?: string;
	} = $props();

	const options = $derived<ComboboxOption[]>(
		columnNames.map((name) => ({ value: name, group: 'Columns' }))
	);

	function setEntry(position: number, next: string) {
		entries = entries.map((entry, index) => (index === position ? next : entry));
	}

	function removeEntry(position: number) {
		entries = entries.filter((_, index) => index !== position);
	}
</script>

<div class="flex flex-col gap-2">
	{#each entries as entry, position (position)}
		<div class="flex items-center gap-2">
			<div class="min-w-0 flex-1">
				<SchemaCombobox
					id={`${idPrefix}-${position}`}
					bind:value={() => entry, (next) => setEntry(position, next)}
					{options}
					{placeholder}
					searchPlaceholder="Column name or expression"
					emptyText="No column matches. Press Enter to use what you typed."
				/>
			</div>
			<Button
				variant="ghost"
				size="icon"
				class="size-8 shrink-0"
				aria-label={`Remove entry ${position + 1}`}
				onclick={() => removeEntry(position)}
			>
				<X />
			</Button>
		</div>
	{/each}
	<Button variant="outline" size="sm" class="h-8 self-start" onclick={() => (entries = [...entries, ''])}>
		<Plus data-icon="inline-start" />
		{addLabel}
	</Button>
</div>
