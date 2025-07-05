<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import X from 'lucide-svelte/icons/x';
	import Plus from 'lucide-svelte/icons/plus';
	import * as Resizable from '$lib/components/ui/resizable/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { onMount, type ComponentProps } from 'svelte';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	import * as Select from '$lib/components/ui/select/index.js';

	// Import our custom components
	import SqlEditor from '$lib/components/app/main_screen/sql_editor.svelte';
	import {
		AddTab,
		DeleteTab,
		SetActiveTab,
		UpdateTabEditorContent,
		SaveActiveDBProps,
		GetAllTabs
	} from '$lib/wailsjs/go/app/Tabs';
	import { activeDBs, suggestions } from '$lib/state.svelte';
	import {
		tabsMap,
		selectedDBDisplay,
		currentColor,
		activePoolID,
		selectedQuery
	} from '$lib/state.svelte';

	let editorHeight = $state(50); // Percentage of the container height
	let outputHeight = $state(50); // Percentage of the container height

	// Handle Tabs

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

	// Write a function to call addTab when pressed cmd + t
	$effect(() => {
		document.addEventListener('keydown', (event: KeyboardEvent) => {
			if (event.key === 't' && event.metaKey) {
				addTab();
			}
		});
	});

	import type { ColumnDef, RowData } from '@tanstack/table-core';

	// Active tab properties
	let { tabID = $bindable(0), tabName = $bindable('') } = $props();
	let editor = $state('');

	let columns = $state<ColumnDef<RowData, unknown>[]>([]);
	let rows = $state<RowData[]>([]);

	onMount(() => {
		getAllTabs();
	});

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

					console.log($activeDBs.length);

					if ($activeDBs.length == 0) {
						$selectedDBDisplay = 'Connect to a database';
						$currentColor = '';
						$activePoolID = '';
					} else {
						console.log('in else', tab.ActiveDB);
						$selectedDBDisplay = tab.ActiveDB || 'Connect to a database';
						$activePoolID = tab.ActiveDBID || '';
						$currentColor = tab.ActiveDBColor || '';
					}

					// Update columns
					for (const column of tab.columns) {
						columns.push({
							accessorKey: column,
							header: column
						});
					}

					for (const row of tab.rows) {
						let cell: Record<string, any> = {};
						for (const resultCell of row) {
							if (resultCell.column && resultCell.value) {
								cell[resultCell.column] = resultCell.value;
							}
						}
						rows.push(cell);
					}
				}
			}
		});

		columns = [];
		rows = [];
	}

	function addTab() {
		queryLoading = false;
		// Send default values for now in activeDBID and activeDB
		AddTab($activePoolID, $selectedDBDisplay, $currentColor).then((tab) => {
			tabsMap.set(tab.ID, tab);

			tabID = tab.ID;
			tabName = tab.Name;
			editor = tab.Editor;
		});

		$selectedDBDisplay = $selectedDBDisplay;
		$activePoolID = $activePoolID;
		$currentColor = $currentColor;

		columns = [];
		rows = [];
	}

	function deleteTab(id: number) {
		tabsMap.delete(id);
		DeleteTab(id).then((tab) => {
			if (tab) {
				tabsMap.set(tab.ID, tab);

				tabID = tab.ID;
				tabName = tab.Name;
				editor = tab.Editor;

				$selectedDBDisplay = tab.ActiveDB || 'Connect to a database';
				$activePoolID = tab.ActiveDBID || '';
				$currentColor = tab.ActiveDBColor || '';

				// Update columns
				for (const column of tab.columns) {
					columns.push({
						accessorKey: column,
						header: column
					});
				}

				for (const row of tab.rows) {
					let cell: Record<string, any> = {};
					for (const resultCell of row) {
						if (resultCell.column && resultCell.value) {
							cell[resultCell.column] = resultCell.value;
						}
					}
					rows.push(cell);
				}
			}
		});

		columns = [];
		rows = [];
	}

	function setActiveTab(id: number) {
		SetActiveTab(id)
			.then((tab) => {
				tabsMap.set(tab.ID, tab);

				tabID = tab.ID;
				tabName = tab.Name;
				editor = tab.Editor;

				if ($activeDBs.length == 0) {
					$selectedDBDisplay = 'Connect to a database';
					$currentColor = '';
					$activePoolID = '';
				} else {
					$selectedDBDisplay = tab.ActiveDB || 'Connect to a database';
					$activePoolID = tab.ActiveDBID || '';
					$currentColor = tab.ActiveDBColor || '';
				}

				// Update columns
				for (const column of tab.columns) {
					columns.push({
						accessorKey: column,
						header: column
					});
				}

				for (const row of tab.rows) {
					let cell: Record<string, any> = {};
					for (const resultCell of row) {
						if (resultCell.column && resultCell.value) {
							cell[resultCell.column] = resultCell.value;
						}
					}
					rows.push(cell);
				}
			})
			.catch((error) => {
				toast.error('Failed to set active tab', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});

		columns = [];
		rows = [];
	}

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

	import DataTable from './data-table.svelte';

	let queryLoading = $state(false);

	function executeQuery() {
		queryLoading = true;
		if ($selectedQuery.trim() == '') {
			toast.error('Please select a query to execute', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}
		if ($activePoolID == '') {
			toast.error('Please select a database to execute the query', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}
		// Execute query
		ExecuteQuery($activePoolID, $selectedQuery, tabID)
			.then((result) => {
				if (!result.ok) {
					queryLoading = false;
					toast.error('Query Failed', {
						description: result.message,
						action: {
							label: 'OK',
							onClick: () => console.info('OK')
						}
					});
					return;
				}

				// Update columns
				for (const column of result.columns) {
					columns.push({
						accessorKey: column,
						header: column
					});
				}

				for (const row of result.rows) {
					let cell: Record<string, any> = {};
					for (const resultCell of row) {
						if (resultCell.column && resultCell.value) {
							cell[resultCell.column] = resultCell.value;
						}
					}
					rows.push(cell);
				}
			})
			.catch((error) => {
				queryLoading = false;
				// Handle errors from the ExecuteQuery call
				toast.error('Query Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});

		columns = [];
		rows = [];
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.altKey && event.key === 'Enter') {
			event.preventDefault();
			executeQuery();
		}
		// Command+S (Mac) or Ctrl+S (Windows/Linux)
		if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 's') {
			event.preventDefault();
			// Your custom logic here
			toast.success('Not Needed! 😂', {
				description: 'Your queries are saved automatically 😂',
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
		}
	}

	document.addEventListener('keydown', handleKeyDown);

	$effect(() => {
		if ($activeDBs.length == 0) {
			$selectedDBDisplay = 'Connect to a database';
			$currentColor = '';
		}
	});

	export function selectActiveDB(activeDBDisplay: string, poolID: string, activeDBColor: string) {
		$selectedDBDisplay = activeDBDisplay;
		$activePoolID = poolID;
		$currentColor = activeDBColor;

		SaveActiveDBProps(tabID, $activePoolID, $selectedDBDisplay, $currentColor);
	}
</script>

<div class="flex h-full flex-1 flex-col overflow-hidden">
	<Tabs.Root value={tabID.toString()} class="flex h-full flex-1 flex-col overflow-hidden">
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
			<div class="flex h-screen flex-1 flex-col px-4">
				<Tabs.Content value={tabID.toString()} class="flex-1 overflow-hidden">
					<div class="flex h-full flex-col">
						<div class="flex items-center justify-between">
							<Select.Root type="single" name="activeDatabase">
								<Select.Trigger
									class="{getColorClass(
										$currentColor
									)} prevent:default w-auto bg-opacity-20 hover:bg-opacity-25"
								>
									{$selectedDBDisplay}
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										{#each $activeDBs as activeDB}
											<Select.Item
												onclick={() =>
													selectActiveDB(
														activeDB.PostgresConnectionName + ' - ' + activeDB.Name,
														activeDB.PoolID,
														activeDB.Colour
													)}
												class="{getColorClass(activeDB.Colour)} bg-opacity-20 hover:bg-opacity-25"
												value={activeDB.ID}
												label={activeDB.Name}
												>{activeDB.PostgresConnectionName} - {activeDB.Name}</Select.Item
											>
										{/each}
									</Select.Group>
								</Select.Content>
							</Select.Root>
							<div class="flex">
								<Button variant="default" size="sm" onclick={executeQuery}>Execute Query</Button>
							</div>
						</div>

						<Resizable.ResizablePaneGroup direction="vertical" class="h-full">
							<!-- SQL Editor Pane -->
							<Resizable.Pane
								defaultSize={editorHeight}
								minSize={10}
								class="rsz-pane my-3 overflow-hidden rounded-md border"
							>
								<SqlEditor
									bind:value={editor}
									bind:selectedQuery={$selectedQuery}
									bind:suggestions={$suggestions}
								/>
							</Resizable.Pane>

							<Resizable.ResizableHandle withHandle />

							<!-- Output Pane -->
							<Resizable.Pane
								defaultSize={outputHeight}
								minSize={10}
								class="rsz-pane overflow-auto"
							>
								<div class="h-full overflow-auto">
									{#if columns.length > 0}
										<DataTable data={rows} {columns} {queryLoading} />
									{:else if queryLoading}
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
										<Skeleton class="my-3 h-[40px] w-full" />
									{/if}
								</div>
							</Resizable.Pane>
						</Resizable.ResizablePaneGroup>
					</div>
				</Tabs.Content>
			</div>
		{/if}
	</Tabs.Root>
</div>
