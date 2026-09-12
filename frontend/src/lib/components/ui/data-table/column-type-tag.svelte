<script lang="ts">
	import type { model } from '$lib/wailsjs/go/models';
	import { columnTypeLabel } from './data-table.svelte.js';
	import ColumnTypeIcon from './column-type-icon.svelte';

	let {
		columnType,
		class: className = ''
	}: { columnType?: model.ColumnType; class?: string } = $props();

	// A result column with no type behind it -- an error row, the "Rows Affected"
	// of a write, a tab restored from a version that did not record types -- gets
	// no tag at all rather than a placeholder.
	const hasType = $derived(Boolean(columnType?.dataType));

	// The full name postgres prints sits in the tooltip, so a header that shows
	// "varchar" can still be read as "character varying" on hover.
	const typeTitle = $derived(columnTypeLabel(columnType));

	const keyBadge = $derived(
		columnType?.isPrimaryKey
			? columnType.isCompositeKey
				? { label: 'CK', title: 'part of the composite primary key' }
				: { label: 'PK', title: 'primary key' }
			: undefined
	);

	const foreignKeyTitle = $derived(
		columnType?.foreignKeyTable
			? `foreign key to ${columnType.foreignKeyTable}`
			: 'foreign key'
	);
</script>

{#snippet badge(label: string, title: string, tone: string)}
	<span
		{title}
		class="rounded-[3px] border gap-1 font-medium leading-none {tone}"
	>
		{label}
	</span>
{/snippet}

{#if hasType}
	<span class="inline-flex shrink-0 items-center justify-center gap-1 {className}">
		<!-- Deliberately inline rather than flex: the icon aligns itself to this
		     text's baseline, which only works while the two share a line box. The
		     leading matches the column name's own 16px line box, so centring the two
		     as flex items lines their baselines up as well. -->
		<span
			title={typeTitle}
			class="text-muted-foreground/70 shrink-0 whitespace-nowrap font-normal mr-2"
		>
			<ColumnTypeIcon {columnType} class="mr-1" />
			{columnType?.displayType}
		</span>

		<!-- Key badges carry a tinted border so they read ahead of the type they sit
		     beside; nullability is the common case, so it stays borderless and quiet. -->
		{#if keyBadge}
			{@render badge(
				keyBadge.label,
				keyBadge.title,
				'border-amber-500/40 text-amber-600 dark:text-amber-400 align-[-0.2em]'
			)}
		{/if}
		{#if columnType?.isForeignKey}
			{@render badge('FK', foreignKeyTitle, 'border-sky-500/40 text-sky-600 dark:text-sky-400 align-[-0.2em]')}
		{/if}
		{#if columnType?.isNullable}
			{@render badge('NULL', 'nullable', 'text-muted-foreground/60 border-transparent align-[-0.2em]')}
		{/if}
	</span>
{/if}
