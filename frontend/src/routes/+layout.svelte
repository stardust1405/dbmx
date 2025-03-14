<script lang="ts">
	import { ModeWatcher } from 'mode-watcher';
	import '../app.css';
	import { Toaster } from '$lib/components/ui/sonner/index.js';
	let { children } = $props();
</script>

<Toaster />

<ModeWatcher />
{@render children()}
