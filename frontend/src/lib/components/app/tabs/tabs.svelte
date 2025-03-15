<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import X from 'lucide-svelte/icons/x';
	import Plus from 'lucide-svelte/icons/plus';
	import ArrowsMaximize from 'lucide-svelte/icons/maximize-2';
	import * as Resizable from '$lib/components/ui/resizable/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { onMount } from 'svelte';

	// Import our custom components
	import SqlEditor from '$lib/components/app/main_screen/sql_editor.svelte';
	import SqlEditor2 from '$lib/components/app/main_screen/sql_editor2.svelte';
	import QueryOutput from '$lib/components/app/main_screen/query_output.svelte';
	import { model } from '$lib/wailsjs/go/models';
	import {
		AddTab,
		DeleteTab,
		GetAllTabs,
		SetActiveTab,
		UpdateTabEditorContent
	} from '$lib/wailsjs/go/app/Tabs';

	let editorHeight = $state(50); // Percentage of the container height
	let outputHeight = $state(50); // Percentage of the container height
	let keywords = $state([
		'SELECT',
		'FROM',
		'WHERE',
		'INNER',
		'LEFT',
		'RIGHT',
		'JOIN',
		'AND',
		'OR',
		'NOT',
		'IN',
		'LIKE',
		'BETWEEN',
		'ORDER BY',
		'LIMIT',
		'location_wise_inventory'
	]);

	function resetSplitView() {
		// Force a reset of the pane sizes
		editorHeight = 50;
		outputHeight = 50;

		// Force a re-render of the panes
		setTimeout(() => {
			const panes = document.querySelectorAll('.rsz-pane');
			if (panes.length >= 2) {
				// Update the style directly if needed
				panes[0].setAttribute('style', 'height: 50%');
				panes[1].setAttribute('style', 'height: 50%');
			}
		}, 0);
	}

	// Handle Tabs

	// Active tab properties
	let tabID = $state(0);
	let tabName = $state('');
	let editor = $state('select * from');
	let output = $state('');
	let activeDBID = $state(0);
	let activeDB = $state('');

	onMount(() => {
		getAllTabs();
	});

	// Call UpdateTabEditorContent on editor change
	let editorUpdateTimer: number | undefined;
	$effect(() => {
		// Explicitly reference editor to ensure reactivity
		const _ = editor;

		// Clear any existing timeout to debounce rapid changes
		if (editorUpdateTimer) clearTimeout(editorUpdateTimer);

		// Set a new timeout to update the content after typing stops
		editorUpdateTimer = setTimeout(() => {
			UpdateTabEditorContent(tabID, editor);
		}, 500);
	});

	// Declare tabsMap as a reactive state variable
	let tabsMap = $state<Map<number, model.Tab>>(new Map<number, model.Tab>());

	// Tab related operations
	function getAllTabs() {
		GetAllTabs().then((tabs) => {
			for (const tab of tabs) {
				tabsMap.set(tab.ID, tab);
				tabsMap = new Map(tabsMap); // Reassign to trigger reactivity

				// Set active tab properties
				if (tab.IsActive) {
					tabID = tab.ID;
					tabName = tab.Name;
					editor = tab.Editor;
					output = tab.Output;
					activeDBID = tab.ActiveDBID || 0;
					activeDB = tab.ActiveDB || '';
				}
			}
		});
	}

	function addTab() {
		// Send default values for now in activeDBID and activeDB
		AddTab(0, '').then((tab) => {
			tabsMap.set(tab.ID, tab);
			tabsMap = new Map(tabsMap); // Reassign to trigger reactivity

			tabID = tab.ID;
			tabName = tab.Name;
			editor = tab.Editor;
			output = tab.Output;
			activeDBID = tab.ActiveDBID || 0;
			activeDB = tab.ActiveDB || '';
		});
	}

	function deleteTab(id: number) {
		tabsMap.delete(id);
		tabsMap = new Map(tabsMap); // Reassign to trigger reactivity
		DeleteTab(id).then((tab) => {
			if (tab) {
				tabsMap.set(tab.ID, tab);
				tabsMap = new Map(tabsMap); // Reassign to trigger reactivity

				tabID = tab.ID;
				tabName = tab.Name;
				editor = tab.Editor;
				output = tab.Output;
				activeDBID = tab.ActiveDBID || 0;
				activeDB = tab.ActiveDB || '';
			}
		});
		console.log(tabsMap);
	}

	function setActiveTab(id: number) {
		SetActiveTab(id).then((tab) => {
			tabsMap.set(tab.ID, tab);
			tabsMap = new Map(tabsMap); // Reassign to trigger reactivity

			tabID = tab.ID;
			tabName = tab.Name;
			editor = tab.Editor;
			output = tab.Output;
			activeDBID = tab.ActiveDBID || 0;
			activeDB = tab.ActiveDB || '';
		});
	}
</script>

<Tabs.Root value={tabID.toString()}>
	<!-- Tabs visible in the header -->
	<header class="flex h-14 shrink-0 items-center gap-2 overflow-auto border-b px-4">
		<Sidebar.Trigger class="-ml-1" />
		<Separator orientation="vertical" />
		<Tabs.List>
			{#each Array.from(tabsMap.entries()) as [key, tab]}
				<div class="mr-2 flex rounded-sm bg-slate-800">
					<Tabs.Trigger
						value={tab.ID.toString()}
						class="flex items-center"
						onclick={() => setActiveTab(tab.ID)}
					>
						{tab.Name}
					</Tabs.Trigger>
					<button
						class="rounded-r-sm bg-slate-900 px-2 py-1 text-slate-300 hover:text-red-700"
						onclick={() => deleteTab(tab.ID)}
					>
						<X size={16} />
					</button>
				</div>
			{/each}
			<button
				class="ml-2 flex items-center gap-1 text-blue-500 hover:text-blue-700"
				onclick={addTab}
			>
				<Plus size={16} /> Add Tab
			</button>
		</Tabs.List>
	</header>

	{#if tabsMap.size > 0}
		<!-- Main Content on screen -->
		<div class="flex h-screen flex-1 flex-col gap-4 p-4">
			<Tabs.Content value={tabID.toString()} class="flex-1 overflow-hidden">
				<div class="flex h-full flex-col">
					<div class="mb-2 flex items-center justify-between">
						<h2 class="text-lg font-semibold">{tabName}</h2>
						<div class="flex gap-2">
							<Button variant="outline" size="sm" onclick={resetSplitView}>
								<ArrowsMaximize size={16} class="mr-2" /> Reset Split
							</Button>
							<Button variant="default" size="sm">Execute Query</Button>
						</div>
					</div>

					<Resizable.ResizablePaneGroup direction="vertical" class="h-full">
						<!-- SQL Editor Pane -->
						<Resizable.ResizablePane
							defaultSize={editorHeight}
							minSize={10}
							class="rsz-pane overflow-hidden rounded-md border"
						>
							<!-- <SqlEditor id={tabID.toString()} value={editor} /> -->
							<SqlEditor2 bind:value={editor} {keywords} />
						</Resizable.ResizablePane>

						<Resizable.ResizableHandle />

						<!-- Output Pane -->
						<Resizable.ResizablePane
							defaultSize={outputHeight}
							minSize={10}
							class="rsz-pane overflow-auto"
						>
							<h1>{output}</h1>
							<!-- <QueryOutput outputs={outputContent} /> -->
						</Resizable.ResizablePane>
					</Resizable.ResizablePaneGroup>
				</div>
			</Tabs.Content>
		</div>
	{/if}
</Tabs.Root>
