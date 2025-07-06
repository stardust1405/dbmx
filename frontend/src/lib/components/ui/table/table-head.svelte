<script lang="ts">
	import type { HTMLThAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLThAttributes> = $props();
</script>

<th bind:this={ref} class={cn('px-5 align-middle', className)} {...restProps}>
	{@render children?.()}
</th>
