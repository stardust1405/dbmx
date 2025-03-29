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

	import * as Select from '$lib/components/ui/select/index.js';

	// Import our custom components
	import SqlEditor from '$lib/components/app/main_screen/sql_editor.svelte';
	import { model } from '$lib/wailsjs/go/models';
	import {
		AddTab,
		DeleteTab,
		GetAllTabs,
		SetActiveTab,
		UpdateTabEditorContent
	} from '$lib/wailsjs/go/app/Tabs';
	import { SvelteMap } from 'svelte/reactivity';
	import { activeDBs, suggestions } from './tabs.svelte.ts';

	let editorHeight = $state(50); // Percentage of the container height
	let outputHeight = $state(50); // Percentage of the container height

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
	let editor = $state('');
	let output = $state<string[]>([]);

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
	let tabsMap = $state<SvelteMap<number, model.Tab>>(new SvelteMap<number, model.Tab>());

	// Tab related operations
	function getAllTabs() {
		GetAllTabs().then((tabs) => {
			if (!tabs) {
				return;
			}
			for (const tab of tabs) {
				tabsMap.set(tab.ID, tab);

				// Set active tab properties
				if (tab.IsActive) {
					tabID = tab.ID;
					tabName = tab.Name;
					editor = tab.Editor;
				}
			}
		});
	}

	// Write a function to call addTab when pressed cmd + t
	$effect(() => {
		document.addEventListener('keydown', (event: KeyboardEvent) => {
			if (event.key === 't' && event.metaKey) {
				addTab();
			}
		});
	});

	function addTab() {
		// Send default values for now in activeDBID and activeDB
		AddTab(0, '').then((tab) => {
			tabsMap.set(tab.ID, tab);

			tabID = tab.ID;
			tabName = tab.Name;
			editor = tab.Editor;
		});
	}

	function deleteTab(id: number) {
		tabsMap.delete(id);
		DeleteTab(id).then((tab) => {
			if (tab) {
				tabsMap.set(tab.ID, tab);

				tabID = tab.ID;
				tabName = tab.Name;
				editor = tab.Editor;
			}
		});
		console.log(tabsMap);
	}

	function setActiveTab(id: number) {
		SetActiveTab(id).then((tab) => {
			tabsMap.set(tab.ID, tab);

			tabID = tab.ID;
			tabName = tab.Name;
			editor = tab.Editor;
		});
	}

	let selectedDBDisplay = $state('Connect to a database');
	let currentColor = $state('');
	let activePoolID = $state('');

	$effect(() => {
		if (activeDBs.length == 0) {
			selectedDBDisplay = 'Connect to a database';
			currentColor = '';
		}
	});

	function getColorClass(color: string): string {
		const colorMap: Record<string, string> = {
			'bg-purple-500': 'bg-purple-500',
			'bg-indigo-500': 'bg-indigo-500',
			'bg-emerald-500': 'bg-emerald-500',
			'bg-red-500': 'bg-red-500',
			'bg-blue-500': 'bg-blue-500',
			'bg-green-500': 'bg-green-500',
			'bg-yellow-500': 'bg-yellow-500',
			'bg-orange-500': 'bg-orange-500',
			'bg-pink-500': 'bg-pink-500'
		};
		return colorMap[color] || '';
	}

	import { toast } from 'svelte-sonner';
	import { ExecuteQuery } from '$lib/wailsjs/go/app/Connections.js';

	let selectedText = $state('');

	let columns = $state<string[]>([]);
	let rows = $state<string[][]>([]);

	function executeQuery() {
		if (selectedText.trim() == '') {
			toast.error('Please select a query to execute', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}
		if (activePoolID == '') {
			toast.error('Please select a database to execute the query', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}
		// Execute query
		ExecuteQuery(activePoolID, selectedText)
			.then((result) => {
				// Update tab output
				columns = result.columns;
				rows = result.rows;
			})
			.catch((error) => {
				// Handle errors from the ExecuteQuery call
				toast.error('Query Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
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
						<Select.Root type="single" name="activeDatabase">
							<Select.Trigger
								class="{getColorClass(currentColor)} w-[180px] bg-opacity-20 hover:bg-opacity-25"
							>
								{selectedDBDisplay}
							</Select.Trigger>
							<Select.Content>
								<Select.Group>
									{#each activeDBs as activeDB}
										<Select.Item
											onclick={() => {
												activePoolID = activeDB.PoolID;
												currentColor = activeDB.Colour;
												selectedDBDisplay = activeDB.PostgresConnectionName + ' - ' + activeDB.Name;
											}}
											class="{getColorClass(activeDB.Colour)} bg-opacity-20 hover:bg-opacity-25"
											value={activeDB.ID}
											label={activeDB.Name}
											>{activeDB.PostgresConnectionName} - {activeDB.Name}</Select.Item
										>
									{/each}
								</Select.Group>
							</Select.Content>
						</Select.Root>
						<div class="flex gap-2">
							<Button variant="outline" size="sm" onclick={resetSplitView}>
								<ArrowsMaximize size={16} class="mr-2" /> Reset Split
							</Button>
							<Button variant="default" size="sm" onclick={executeQuery}>Execute Query</Button>
						</div>
					</div>

					<Resizable.ResizablePaneGroup direction="vertical" class="h-full">
						<!-- SQL Editor Pane -->
						<Resizable.ResizablePane
							defaultSize={editorHeight}
							minSize={10}
							class="rsz-pane overflow-hidden rounded-md border"
						>
							<SqlEditor bind:value={editor} bind:selectedText {suggestions} />
						</Resizable.ResizablePane>

						<Resizable.ResizableHandle />

						<!-- Output Pane -->
						<Resizable.ResizablePane
							defaultSize={outputHeight}
							minSize={10}
							class="rsz-pane overflow-auto"
						>
							{#each columns as column}
								{column} |
							{/each}
							{#each rows as row}
								{#each row as cell}
									{cell} |
								{/each}
							{/each}
						</Resizable.ResizablePane>
					</Resizable.ResizablePaneGroup>
				</div>
			</Tabs.Content>
		</div>
	{/if}
</Tabs.Root>
