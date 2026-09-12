<script lang="ts">
	import type { model } from '$lib/wailsjs/go/models';
	import { columnTypeIcon, columnTypeLabel } from './data-table.svelte.js';

	let {
		columnType,
		class: className = ''
	}: { columnType?: model.ColumnType; class?: string } = $props();

	const Icon = $derived(columnTypeIcon(columnType));
	const label = $derived(columnTypeLabel(columnType));
</script>

<!-- The glyph is the only place a column's type appears, so it is labelled rather
     than hidden: the tooltip spells the type out for a pointer, the aria-label for
     a screen reader. -->
<!-- The icon is set like a glyph, not like a box: inline-block on the text
     baseline, dropped 0.2em so its ink centres on the cap height of the type name
     beside it. Centring it in the line box instead -- which is what any flex
     alignment does -- lands it about 2px low, because a line box reserves
     descender space below the baseline that the text's own ink never uses.
     The svg is display:block so the span's height is the glyph's 12px exactly,
     with no inline descender gap of its own underneath. -->
<span
	role="img"
	aria-label={label}
	title={label}
	class="text-muted-foreground/70 inline-block align-[-0.2em] {className}"
>
	<Icon size={12} class="block" />
</span>
