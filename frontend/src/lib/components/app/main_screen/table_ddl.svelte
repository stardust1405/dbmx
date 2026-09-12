<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';
	import { mode } from 'mode-watcher';
	import { get } from 'svelte/store';
	import { toast } from 'svelte-sonner';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { Check, Copy } from 'lucide-svelte';
	import { editorFontSize } from '$lib/state.svelte';
	import { GetTableDDL } from '$lib/wailsjs/go/app/Connections.js';
	import { ClipboardSetText } from '$lib/wailsjs/runtime/runtime.js';
	import { defineSQLThemes, sqlThemeFor } from './monaco-sql-theme';
	import type * as MonacoNS from 'monaco-editor';

	let { tabID, tableName = '' }: { tabID: number; tableName?: string } = $props();

	let editorContainer: HTMLElement;
	let editor: MonacoNS.editor.IStandaloneCodeEditor | null = null;
	let monacoInstance: typeof import('monaco-editor') | null = null;

	let ddl = $state('');
	let loading = $state(true);
	let copied = $state(false);
	let copiedTimer: ReturnType<typeof setTimeout> | null = null;

	onMount(async () => {
		const [monaco] = await Promise.all([loader.init(), loadDDL()]);
		monacoInstance = monaco;
		defineSQLThemes(monaco);

		editor = monaco.editor.create(editorContainer, {
			value: ddl,
			language: 'sql',
			theme: sqlThemeFor(mode.current),
			// The definition is generated, so the editor is a viewer: editing it would
			// imply the change lands on the table, which it does not. Selection and copy
			// still work, and readOnly keeps Monaco's own context menu honest.
			readOnly: true,
			domReadOnly: true,
			automaticLayout: true,
			minimap: { enabled: false },
			fontSize: get(editorFontSize),
			wordWrap: 'on',
			scrollBeyondLastLine: false,
			renderLineHighlight: 'none',
			smoothScrolling: true,
			fixedOverflowWidgets: true
		});
	});

	// Follow the app's light/dark mode, matching the SQL editor.
	$effect(() => {
		const currentMode = mode.current;
		if (!monacoInstance || !editor) return;
		monacoInstance.editor.setTheme(sqlThemeFor(currentMode));
	});

	async function loadDDL() {
		if (!tableName) {
			loading = false;
			return;
		}
		try {
			ddl = await GetTableDDL(tabID, tableName);
			editor?.setValue(ddl);
		} catch (error) {
			toast.error('Failed to load table definition', { description: String(error) });
		} finally {
			loading = false;
		}
	}

	function copyDDL() {
		if (!ddl) return;
		// The Wails clipboard runtime is used rather than navigator.clipboard, which is
		// gated on a secure context and permissions the webview does not reliably grant.
		ClipboardSetText(ddl)
			.then((ok) => {
				if (!ok) {
					toast.error('Failed to copy definition');
					return;
				}
				copied = true;
				if (copiedTimer) clearTimeout(copiedTimer);
				copiedTimer = setTimeout(() => (copied = false), 2000);
			})
			.catch((error) => {
				toast.error('Failed to copy definition', { description: String(error) });
			});
	}

	onDestroy(() => {
		if (copiedTimer) clearTimeout(copiedTimer);
		editor?.dispose();
		editor = null;
	});
</script>

<div class="flex h-full min-h-0 flex-1 flex-col overflow-hidden">
	<div class="flex items-center justify-between px-1 pb-2">
		<span class="text-muted-foreground truncate text-sm">
			{tableName ? `Definition of ${tableName}` : 'No table selected'}
		</span>
		<Button
			variant="outline"
			size="sm"
			class="h-8"
			onclick={copyDDL}
			disabled={loading || !ddl}
		>
			{#if copied}
				<Check data-icon="inline-start" />
				Copied
			{:else}
				<Copy data-icon="inline-start" />
				Copy
			{/if}
		</Button>
	</div>

	<div class="relative flex h-1 min-h-0 flex-1 overflow-hidden rounded-xl border">
		<!-- Monaco is mounted unconditionally so its container keeps a stable size; the
		     skeleton covers it while the definition is still being fetched. -->
		<div bind:this={editorContainer} class="h-full w-full"></div>
		{#if loading}
			<div class="bg-background absolute inset-0 flex flex-col gap-3 p-3">
				<Skeleton class="h-[24px] w-1/3" />
				<Skeleton class="h-[24px] w-2/3" />
				<Skeleton class="h-[24px] w-1/2" />
				<Skeleton class="h-[24px] w-3/5" />
			</div>
		{/if}
	</div>
</div>
