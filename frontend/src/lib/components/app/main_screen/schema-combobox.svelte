<script module lang="ts">
	export type ComboboxOption = {
		value: string;
		/** Shown to the right of the value: an enum's labels, a type's kind. */
		detail?: string;
		/** Options carrying the same group are listed together, groups in first-seen order. */
		group?: string;
	};
</script>

<script lang="ts">
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Command from '$lib/components/ui/command/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import { ChevronsUpDown, Check } from 'lucide-svelte';
	import { cn } from '$lib/utils.js';

	// A schema form picks from a catalog list that is never complete: a type can carry a
	// modifier the list doesn't hold ("character varying(255)"), an array suffix, or a
	// name only this database has. So the search box doubles as the input -- what is
	// typed is offered as the value alongside the matches.
	let {
		value = $bindable(''),
		options = [],
		placeholder = 'Select',
		searchPlaceholder = 'Search or type a value',
		emptyText = 'No match. Press Enter to use what you typed.',
		id,
		disabled = false,
		invalid = false,
		allowCustom = true
	}: {
		value?: string;
		options?: ComboboxOption[];
		placeholder?: string;
		searchPlaceholder?: string;
		emptyText?: string;
		id?: string;
		disabled?: boolean;
		invalid?: boolean;
		allowCustom?: boolean;
	} = $props();

	let open = $state(false);
	let search = $state('');

	const groups = $derived.by(() => {
		const byGroup = new Map<string, ComboboxOption[]>();
		for (const option of options) {
			const key = option.group ?? '';
			const existing = byGroup.get(key);
			if (existing) existing.push(option);
			else byGroup.set(key, [option]);
		}
		return [...byGroup];
	});

	const typed = $derived(search.trim());
	const isNewValue = $derived(
		allowCustom && typed !== '' && !options.some((option) => option.value === typed)
	);

	function choose(next: string) {
		value = next;
		open = false;
		search = '';
	}
</script>

<Popover.Root bind:open>
	<Popover.Trigger {id} {disabled}>
		{#snippet child({ props })}
			<Button
				{...props}
				variant="outline"
				role="combobox"
				aria-expanded={open}
				aria-invalid={invalid}
				class="w-full justify-between font-normal"
			>
				<span class={cn('truncate', value === '' && 'text-muted-foreground')}>
					{value === '' ? placeholder : value}
				</span>
				<ChevronsUpDown class="opacity-50" />
			</Button>
		{/snippet}
	</Popover.Trigger>
	<Popover.Content class="w-(--bits-popover-anchor-width) p-0" align="start">
		<Command.Root shouldFilter={true}>
			<Command.Input placeholder={searchPlaceholder} bind:value={search} />
			<Command.List class="max-h-64">
				<Command.Empty>{emptyText}</Command.Empty>
				{#if isNewValue}
					<Command.Group heading="Use as typed">
						<!-- forceMount keeps this item reachable while Command filters the rest
						     away, which is exactly when it matters. -->
						<Command.Item value={typed} onSelect={() => choose(typed)} forceMount>
							<Check class="opacity-0" />
							<span class="truncate font-mono text-xs">{typed}</span>
						</Command.Item>
					</Command.Group>
				{/if}
				{#each groups as [heading, groupOptions] (heading)}
					<Command.Group heading={heading || undefined}>
						{#each groupOptions as option (option.value)}
							<Command.Item value={option.value} onSelect={() => choose(option.value)}>
								<Check class={cn(value !== option.value && 'opacity-0')} />
								<span class="truncate">{option.value}</span>
								{#if option.detail}
									<Badge
										variant="secondary"
										class="ml-auto max-w-40 truncate px-1.5 py-0 text-[10px] font-normal"
									>
										{option.detail}
									</Badge>
								{/if}
							</Command.Item>
						{/each}
					</Command.Group>
				{/each}
			</Command.List>
		</Command.Root>
	</Popover.Content>
</Popover.Root>
