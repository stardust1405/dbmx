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
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Label } from '$lib/components/ui/label/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Popover from '$lib/components/ui/popover/index.js';
	import { Badge } from '$lib/components/ui/badge/index.js';

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

	import { toast } from 'svelte-sonner';
	import { ExecuteQuery, GetTableData } from '$lib/wailsjs/go/app/Connections.js';

	import DataTable from './data-table.svelte';
	import { columns, rows } from '$lib/state.svelte';

	let editorHeight = $state(50); // Percentage of the container height
	let outputHeight = $state(50); // Percentage of the container height

	let queryLoading = $state(false);

	// Handle Tabs

	// Active tab properties
	let {
		tabID = $bindable(0),
		tabName = $bindable(''),
		tabType = $bindable(''),
		tabDBName = $bindable(''),
		tabTableDBPoolID = $bindable(''),
		tabPostgresConnName = $bindable(''),
		tabPostgresConnID = $bindable(0),
		select = $bindable(''),
		limit = $bindable(''),
		offset = $bindable(''),
		where = $bindable(''),
		orderBy = $bindable(''),
		groupBy = $bindable('')
	} = $props();
	let editor = $state('');

	// Table view tab state (for Data/Structure/Indexes)
	let tableViewTab = $state('data');

	// Multi-select input state for SELECT
	let selectedColumns = $state<string[]>([]);
	let inputValue = $state('');
	let showSuggestions = $state(false);
	let filteredSuggestions = $state<string[]>([]);
	let selectedIndex = $state(-1); // For arrow key navigation
	let originalInputValue = $state(''); // Store original input when navigating
	let focusedChipIndex = $state(-1); // For chip navigation
	let editingChipIndex = $state(-1); // For chip editing
	let editingChipValue = $state(''); // Value being edited

	// Multi-select input state for WHERE
	let selectedWhereColumns = $state<string[]>([]);
	let whereInputValue = $state('');
	let showWhereSuggestions = $state(false);
	let filteredWhereSuggestions = $state<string[]>([]);
	let selectedWhereIndex = $state(-1);
	let originalWhereInputValue = $state('');
	let focusedWhereChipIndex = $state(-1);
	let editingWhereChipIndex = $state(-1);
	let editingWhereChipValue = $state('');

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
					tabType = tab.Type;

					// Properties for table view tab
					tabDBName = tab.DBName || '';
					tabTableDBPoolID = tab.ActiveDBID || '';
					tabPostgresConnName = tab.PostgresConnName || '';
					tabPostgresConnID = tab.PostgresConnID || 0;

					select = tab.Select;
					limit = tab.Limit;
					offset = tab.Offset;
					where = tab.Where;
					orderBy = tab.OrderBy;
					groupBy = tab.GroupBy;

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

					if (tab.rows && tab.columns) {
						// Update columns
						for (const column of tab.columns) {
							columns.set([
								...$columns,
								{
									accessorKey: column,
									header: column
								}
							]);
						}

						for (const row of tab.rows) {
							let cell: Record<string, any> = {};
							for (const resultCell of row) {
								if (resultCell.column && resultCell.value) {
									cell[resultCell.column] = resultCell.value;
								}
							}
							rows.set([...$rows, cell]);
						}
					}
				}
			}
		});

		columns.set([]);
		rows.set([]);
	}

	export function addTab(
		tableName: string = '',
		postgresConnID: number = 0,
		dbName: string = '',
		tableDBPoolID: string = '',
		postgresConnName: string = ''
	) {
		let tableType = 'editor';
		let tabDisplayName = 'Editor';

		// If table's pool id is provided, use that instead of the active pool id
		// This is used when we want to open a table in a new tab
		// In case of table, we want the pool id of the table's database

		// First assign the active global pool id
		let poolID: string = $activePoolID;

		// If table's pool id is provided, use that instead of the active pool id
		if (tableName.trim() != '') {
			tableType = 'table';
			tabDisplayName = tableName;
			poolID = tableDBPoolID;
		}
		// Send default values for now in activeDBID and activeDB
		AddTab(
			poolID,
			$selectedDBDisplay,
			$currentColor,
			tabDisplayName,
			tableType,
			postgresConnID,
			dbName,
			postgresConnName
		)
			.then((tab) => {
				queryLoading = false;
				tabsMap.set(tab.ID, tab);

				tabID = tab.ID;
				tabName = tab.Name;
				tabType = tab.Type;

				// Properties for table view tab
				tabDBName = tab.DBName || '';
				tabTableDBPoolID = tab.ActiveDBID || '';
				tabPostgresConnName = tab.PostgresConnName || '';
				tabPostgresConnID = tab.PostgresConnID || 0;

				select = tab.Select;
				limit = tab.Limit;
				offset = tab.Offset;
				where = tab.Where;
				orderBy = tab.OrderBy;
				groupBy = tab.GroupBy;

				editor = tab.Editor;
			})
			.catch((error) => {
				toast.error('Failed to add tab', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});

		$selectedDBDisplay = $selectedDBDisplay;
		$activePoolID = $activePoolID;
		$currentColor = $currentColor;

		columns.set([]);
		rows.set([]);
	}

	function deleteTab(id: number) {
		// Delete the old tab from the map
		tabsMap.delete(id);
		DeleteTab(id)
			.then((tab) => {
				if (tab) {
					queryLoading = false;

					// Set the new tab as active
					tabsMap.set(tab.ID, tab);

					tabID = tab.ID;
					tabName = tab.Name;
					tabType = tab.Type;

					// Properties for table view tab
					tabDBName = tab.DBName || '';
					tabTableDBPoolID = tab.ActiveDBID || '';
					tabPostgresConnName = tab.PostgresConnName || '';
					tabPostgresConnID = tab.PostgresConnID || 0;

					select = tab.Select;
					limit = tab.Limit;
					offset = tab.Offset;
					where = tab.Where;
					orderBy = tab.OrderBy;
					groupBy = tab.GroupBy;

					editor = tab.Editor;

					$selectedDBDisplay = tab.ActiveDB || 'Connect to a database';
					$activePoolID = tab.ActiveDBID || '';
					$currentColor = tab.ActiveDBColor || '';

					if (tab.rows && tab.columns) {
						// Update columns
						for (const column of tab.columns) {
							columns.set([
								...$columns,
								{
									accessorKey: column,
									header: column
								}
							]);
						}

						for (const row of tab.rows) {
							let cell: Record<string, any> = {};
							for (const resultCell of row) {
								if (resultCell.column && resultCell.value) {
									cell[resultCell.column] = resultCell.value;
								}
							}
							rows.set([...$rows, cell]);
						}
					}
				}
			})
			.catch((error) => {
				toast.error('Failed to delete tab', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});

		columns.set([]);
		rows.set([]);
	}

	function setActiveTab(id: number) {
		SetActiveTab(id)
			.then((tab) => {
				queryLoading = false;
				tabsMap.set(tab.ID, tab);

				tabID = tab.ID;
				tabName = tab.Name;
				tabType = tab.Type;

				// Properties for table view tab
				tabDBName = tab.DBName || '';
				tabTableDBPoolID = tab.ActiveDBID || '';
				tabPostgresConnName = tab.PostgresConnName || '';
				tabPostgresConnID = tab.PostgresConnID || 0;

				select = tab.Select;
				limit = tab.Limit;
				offset = tab.Offset;
				where = tab.Where;
				orderBy = tab.OrderBy;
				groupBy = tab.GroupBy;

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

				if (tab.rows && tab.columns) {
					// Update columns
					for (const column of tab.columns) {
						columns.set([
							...$columns,
							{
								accessorKey: column,
								header: column
							}
						]);
					}

					for (const row of tab.rows) {
						let cell: Record<string, any> = {};
						for (const resultCell of row) {
							if (resultCell.column && resultCell.value) {
								cell[resultCell.column] = resultCell.value;
							}
						}
						rows.set([...$rows, cell]);
					}
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

		columns.set([]);
		rows.set([]);
		inputValue = '';
		showSuggestions = true;
		originalInputValue = '';
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

	function executeQuery() {
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

		queryLoading = true;
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
					columns.set([
						...$columns,
						{
							accessorKey: column,
							header: column
						}
					]);
				}

				for (const row of result.rows) {
					let cell: Record<string, any> = {};
					for (const resultCell of row) {
						if (resultCell.column && resultCell.value) {
							cell[resultCell.column] = resultCell.value;
						}
					}
					rows.set([...$rows, cell]);
				}
				queryLoading = false;
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

		columns.set([]);
		rows.set([]);
	}

	function getTableData() {
		if (tabTableDBPoolID == '') {
			toast.error('Please select a database to execute the query', {
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}

		queryLoading = true;
		// Execute query
		GetTableData(tabTableDBPoolID, tabID, tabName, select, limit, offset, where, orderBy, groupBy)
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
					columns.set([
						...$columns,
						{
							accessorKey: column,
							header: column
						}
					]);
				}

				for (const row of result.rows) {
					let cell: Record<string, any> = {};
					for (const resultCell of row) {
						if (resultCell.column && resultCell.value) {
							cell[resultCell.column] = resultCell.value;
						}
					}
					rows.set([...$rows, cell]);
				}
				queryLoading = false;
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

		columns.set([]);
		rows.set([]);
	}

	document.addEventListener('keydown', handleKeyDown);
	function handleKeyDown(event: KeyboardEvent) {
		if (event.altKey && event.key === 'Enter') {
			event.preventDefault();
			if (tabType == 'table') {
				getTableData();
			} else {
				executeQuery();
			}
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

	function selectActiveDB(activeDBDisplay: string, poolID: string, activeDBColor: string) {
		$selectedDBDisplay = activeDBDisplay;
		$activePoolID = poolID;
		$currentColor = activeDBColor;

		SaveActiveDBProps(tabID, $activePoolID, $selectedDBDisplay, $currentColor);
	}

	// Multi-select input functions
	function handleInputChange(value: string) {
		inputValue = value;
		originalInputValue = value; // Store the original typed value
		selectedIndex = -1; // Reset selection when typing
		showSuggestions = true;
		showWhereSuggestions = false;
		filteredSuggestions = Array.from($suggestions).filter(
			(s) => s.toLowerCase().includes(value.toLowerCase()) && !selectedColumns.includes(s)
		);
	}

	function selectColumn(column: string) {
		if (!selectedColumns.includes(column)) {
			selectedColumns = [...selectedColumns, column];
			updateSelectValue();
		}
		inputValue = '';
		showSuggestions = true;
		selectedIndex = -1;
	}

	function removeColumn(column: string) {
		selectedColumns = selectedColumns.filter((c) => c !== column);
		updateSelectValue();
	}

	function updateSelectValue() {
		select = selectedColumns.join(', ');
	}

	function handleInputKeyDown(event: KeyboardEvent) {
		// Handle chip navigation when input is empty
		if (!inputValue && selectedColumns.length > 0) {
			if (event.key === 'ArrowLeft') {
				event.preventDefault();
				focusedChipIndex = Math.max(0, selectedColumns.length - 1);
				return;
			}
		}

		if (event.key === 'ArrowDown') {
			event.preventDefault();
			if (filteredSuggestions.length > 0) {
				// Store original input value when starting navigation
				if (selectedIndex === -1) {
					originalInputValue = inputValue;
				}
				selectedIndex = Math.min(selectedIndex + 1, filteredSuggestions.length - 1);
				// Update input with selected suggestion
				inputValue = filteredSuggestions[selectedIndex];
				showSuggestions = true;
				scrollToSelectedItem();
			}
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			if (filteredSuggestions.length > 0) {
				// Store original input value when starting navigation
				if (selectedIndex === -1) {
					originalInputValue = inputValue;
				}
				selectedIndex = Math.max(selectedIndex - 1, -1);
				// Update input with selected suggestion or restore original
				if (selectedIndex === -1) {
					inputValue = originalInputValue;
				} else {
					inputValue = filteredSuggestions[selectedIndex];
				}
				scrollToSelectedItem();
			}
		} else if (event.key === 'Enter') {
			event.preventDefault();
			if (selectedIndex >= 0 && filteredSuggestions[selectedIndex]) {
				selectColumn(filteredSuggestions[selectedIndex]);
			} else if (inputValue.trim()) {
				selectColumn(inputValue.trim());
			}
		} else if (event.key === 'Escape') {
			showSuggestions = false;
			selectedIndex = -1;
			// Restore original input value
			inputValue = originalInputValue;
		} else if (event.key === 'Backspace' && !inputValue && selectedColumns.length > 0) {
			// Remove last chip when backspace is pressed and input is empty
			event.preventDefault();
			const lastColumn = selectedColumns[selectedColumns.length - 1];
			removeColumn(lastColumn);
		}
	}

	function handleChipKeyDown(event: KeyboardEvent, chipIndex: number) {
		// Don't handle navigation if we're currently editing this chip
		if (editingChipIndex === chipIndex) {
			return;
		}

		if (event.key === 'ArrowLeft') {
			event.preventDefault();
			focusedChipIndex = Math.max(0, chipIndex - 1);
		} else if (event.key === 'ArrowRight') {
			event.preventDefault();
			if (chipIndex < selectedColumns.length - 1) {
				focusedChipIndex = chipIndex + 1;
			} else {
				// Move focus back to input
				focusedChipIndex = -1;
				const input = document.querySelector('.select-input') as HTMLInputElement;
				input?.focus();
			}
		} else if (event.key === ' ' || event.key === 'Enter') {
			event.preventDefault();
			startEditingChip(chipIndex);
		} else if (event.key === 'Delete' || event.key === 'Backspace') {
			event.preventDefault();
			removeColumn(selectedColumns[chipIndex]);
			focusedChipIndex = Math.min(chipIndex, selectedColumns.length - 2);
		}
	}

	function startEditingChip(chipIndex: number) {
		editingChipIndex = chipIndex;
		editingChipValue = selectedColumns[chipIndex];
		setTimeout(() => {
			const input = document.querySelector('.chip-edit-input') as HTMLInputElement;
			input?.focus();
			// Position cursor at the end instead of selecting all text
			input?.setSelectionRange(input.value.length, input.value.length);
		}, 0);
	}

	function finishEditingChip() {
		const wasEditingIndex = editingChipIndex;
		if (editingChipIndex >= 0 && editingChipValue.trim()) {
			selectedColumns[editingChipIndex] = editingChipValue.trim();
			updateSelectValue();
		}
		editingChipIndex = -1;
		editingChipValue = '';
		// Restore focus to the chip that was being edited
		setTimeout(() => {
			focusedChipIndex = wasEditingIndex;
			const chip = document.querySelector(`[data-chip-index="${wasEditingIndex}"]`) as HTMLElement;
			chip?.focus();
		}, 10);
	}

	function cancelEditingChip() {
		const wasEditingIndex = editingChipIndex;
		editingChipIndex = -1;
		editingChipValue = '';
		// Restore focus to the chip that was being edited
		setTimeout(() => {
			focusedChipIndex = wasEditingIndex;
			// Ensure the chip element gets focus
			const chip = document.querySelector(`[data-chip-index="${wasEditingIndex}"]`) as HTMLElement;
			chip?.focus();
		}, 10); // Slightly longer delay to ensure DOM is updated
	}

	function scrollToSelectedItem() {
		if (selectedIndex >= 0) {
			// Use setTimeout to ensure DOM is updated
			setTimeout(() => {
				const suggestionsList = document.getElementById('suggestions-list');
				const selectedItem = suggestionsList?.querySelector(`[data-index="${selectedIndex}"]`);
				if (selectedItem && suggestionsList) {
					selectedItem.scrollIntoView({
						behavior: 'smooth',
						block: 'nearest'
					});
				}
			}, 0);
		}
	}

	// WHERE input functions (similar to SELECT functions)
	function handleWhereInputChange(value: string) {
		whereInputValue = value;
		originalWhereInputValue = value;
		selectedWhereIndex = -1;
		showWhereSuggestions = true;
		showSuggestions = false;
		filteredWhereSuggestions = Array.from($suggestions).filter(
			(s) => s.toLowerCase().includes(value.toLowerCase()) && !selectedWhereColumns.includes(s)
		);
	}

	function selectWhereColumn(column: string) {
		if (!selectedWhereColumns.includes(column)) {
			selectedWhereColumns = [...selectedWhereColumns, column];
			updateWhereValue();
		}
		whereInputValue = '';
		showWhereSuggestions = false;
		selectedWhereIndex = -1;
	}

	function removeWhereColumn(column: string) {
		selectedWhereColumns = selectedWhereColumns.filter((c) => c !== column);
		updateWhereValue();
	}

	function updateWhereValue() {
		where = selectedWhereColumns.join(' AND ');
	}

	function handleWhereInputKeyDown(event: KeyboardEvent) {
		// Handle chip navigation when input is empty
		if (!whereInputValue && selectedWhereColumns.length > 0) {
			if (event.key === 'ArrowLeft') {
				event.preventDefault();
				focusedWhereChipIndex = Math.max(0, selectedWhereColumns.length - 1);
				return;
			}
		}

		if (event.key === 'ArrowDown') {
			event.preventDefault();
			if (filteredWhereSuggestions.length > 0) {
				if (selectedWhereIndex === -1) {
					originalWhereInputValue = whereInputValue;
				}
				selectedWhereIndex = Math.min(selectedWhereIndex + 1, filteredWhereSuggestions.length - 1);
				whereInputValue = filteredWhereSuggestions[selectedWhereIndex];
				showWhereSuggestions = true;
				scrollToSelectedWhereItem();
			}
		} else if (event.key === 'ArrowUp') {
			event.preventDefault();
			if (filteredWhereSuggestions.length > 0) {
				if (selectedWhereIndex === -1) {
					originalWhereInputValue = whereInputValue;
				}
				selectedWhereIndex = Math.max(selectedWhereIndex - 1, -1);
				if (selectedWhereIndex === -1) {
					whereInputValue = originalWhereInputValue;
				} else {
					whereInputValue = filteredWhereSuggestions[selectedWhereIndex];
				}
				scrollToSelectedWhereItem();
			}
		} else if (event.key === 'Enter') {
			event.preventDefault();
			if (selectedWhereIndex >= 0 && filteredWhereSuggestions[selectedWhereIndex]) {
				selectWhereColumn(filteredWhereSuggestions[selectedWhereIndex]);
			} else if (whereInputValue.trim()) {
				selectWhereColumn(whereInputValue.trim());
			}
		} else if (event.key === 'Escape') {
			showWhereSuggestions = false;
			selectedWhereIndex = -1;
			whereInputValue = originalWhereInputValue;
		} else if (event.key === 'Backspace' && !whereInputValue && selectedWhereColumns.length > 0) {
			event.preventDefault();
			const lastColumn = selectedWhereColumns[selectedWhereColumns.length - 1];
			removeWhereColumn(lastColumn);
		}
	}

	function handleWhereChipKeyDown(event: KeyboardEvent, chipIndex: number) {
		// Don't handle navigation if we're currently editing this chip
		if (editingWhereChipIndex === chipIndex) {
			return;
		}

		if (event.key === 'ArrowLeft') {
			event.preventDefault();
			focusedWhereChipIndex = Math.max(0, chipIndex - 1);
		} else if (event.key === 'ArrowRight') {
			event.preventDefault();
			if (chipIndex < selectedWhereColumns.length - 1) {
				focusedWhereChipIndex = chipIndex + 1;
			} else {
				// Move focus back to input
				focusedWhereChipIndex = -1;
				const input = document.querySelector('.where-input') as HTMLInputElement;
				input?.focus();
			}
		} else if (event.key === ' ' || event.key === 'Enter') {
			event.preventDefault();
			startEditingWhereChip(chipIndex);
		} else if (event.key === 'Delete' || event.key === 'Backspace') {
			event.preventDefault();
			removeWhereColumn(selectedWhereColumns[chipIndex]);
			focusedWhereChipIndex = Math.min(chipIndex, selectedWhereColumns.length - 2);
		}
	}

	function startEditingWhereChip(chipIndex: number) {
		editingWhereChipIndex = chipIndex;
		editingWhereChipValue = selectedWhereColumns[chipIndex];
		setTimeout(() => {
			const input = document.querySelector('.where-chip-edit-input') as HTMLInputElement;
			input?.focus();
			// Position cursor at the end instead of selecting all text
			input?.setSelectionRange(input.value.length, input.value.length);
		}, 0);
	}

	function finishEditingWhereChip() {
		const wasEditingIndex = editingWhereChipIndex;
		if (editingWhereChipIndex >= 0 && editingWhereChipValue.trim()) {
			selectedWhereColumns[editingWhereChipIndex] = editingWhereChipValue.trim();
			updateWhereValue();
		}
		editingWhereChipIndex = -1;
		editingWhereChipValue = '';
		// Restore focus to the chip that was being edited
		setTimeout(() => {
			focusedWhereChipIndex = wasEditingIndex;
			const chip = document.querySelector(
				`[data-where-chip-index="${wasEditingIndex}"]`
			) as HTMLElement;
			chip?.focus();
		}, 10);
	}

	function cancelEditingWhereChip() {
		const wasEditingIndex = editingWhereChipIndex;
		editingWhereChipIndex = -1;
		editingWhereChipValue = '';
		// Restore focus to the chip that was being edited
		setTimeout(() => {
			focusedWhereChipIndex = wasEditingIndex;
			// Ensure the chip element gets focus
			const chip = document.querySelector(
				`[data-where-chip-index="${wasEditingIndex}"]`
			) as HTMLElement;
			chip?.focus();
		}, 10); // Slightly longer delay to ensure DOM is updated
	}

	function scrollToSelectedWhereItem() {
		if (selectedWhereIndex >= 0) {
			setTimeout(() => {
				const suggestionsList = document.getElementById('where-suggestions-list');
				const selectedItem = suggestionsList?.querySelector(`[data-index="${selectedWhereIndex}"]`);
				if (selectedItem && suggestionsList) {
					selectedItem.scrollIntoView({
						behavior: 'smooth',
						block: 'nearest'
					});
				}
			}, 0);
		}
	}

	$effect(() => {
		if ($activeDBs.length == 0) {
			$selectedDBDisplay = 'Connect to a database';
			$currentColor = '';
		}
	});

	// Initialize selectedColumns from select value
	$effect(() => {
		if (select && select.trim()) {
			selectedColumns = select
				.split(',')
				.map((s) => s.trim())
				.filter((s) => s);
		} else {
			selectedColumns = [];
		}
	});

	// Initialize selectedWhereColumns from where value
	$effect(() => {
		if (where && where.trim()) {
			selectedWhereColumns = where
				.split(' AND ')
				.map((s) => s.trim())
				.filter((s) => s);
		} else {
			selectedWhereColumns = [];
		}
	});

	// Handle clicking outside to close suggestions
	$effect(() => {
		function handleClickOutside(event: MouseEvent) {
			const target = event.target as Element;
			if (!target.closest('.multi-select-container')) {
				showSuggestions = false;
				showWhereSuggestions = false;
			}
		}

		if (showSuggestions || showWhereSuggestions) {
			document.addEventListener('click', handleClickOutside);
			return () => document.removeEventListener('click', handleClickOutside);
		}
	});

	// Focus management for chips
	$effect(() => {
		if (focusedChipIndex >= 0) {
			const chip = document.querySelector(`[data-chip-index="${focusedChipIndex}"]`) as HTMLElement;
			chip?.focus();
		}
	});

	$effect(() => {
		if (focusedWhereChipIndex >= 0) {
			const chip = document.querySelector(
				`[data-where-chip-index="${focusedWhereChipIndex}"]`
			) as HTMLElement;
			chip?.focus();
		}
	});

	// Call UpdateTabEditorContent on editor change
	let editorUpdateTimer: number | undefined;
	$effect(() => {
		// Explicitly reference editor to ensure reactivity
		const _ = editor;
		const selectQuery = select;
		const limitQuery = limit;
		const offsetQuery = offset;
		const whereQuery = where;
		const orderByQuery = orderBy;
		const groupByQuery = groupBy;

		// Clear any existing timeout to debounce rapid changes
		if (editorUpdateTimer) clearTimeout(editorUpdateTimer);

		// Set a new timeout to update the content after typing stops
		editorUpdateTimer = setTimeout(() => {
			UpdateTabEditorContent(tabID, editor, select, limit, offset, where, orderBy, groupBy);
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
</script>

<div class="my-2 flex h-full flex-1 flex-col overflow-hidden rounded-md border">
	<Tabs.Root value={tabID.toString()} class="flex h-full flex-1 flex-col overflow-hidden">
		<!-- Tabs visible in the header -->
		<header class="flex h-12 items-center gap-2 border-b px-4">
			<Sidebar.Trigger class="-ml-1" />
			<Separator orientation="vertical" />
			<Tabs.List class="thin-scrollbar scrollbar-thin overflow-x-auto overflow-y-hidden">
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
					onclick={() => addTab()}
				>
					<Plus size={16} /> Add Tab
				</button>
			</Tabs.List>
		</header>

		{#if tabsMap.size > 0}
			<!-- Main Content on screen -->

			{#if tabType == 'table'}
				<div class="flex h-screen flex-1 flex-col">
					<div class="flex h-full flex-1 flex-col justify-center">
						<!-- Breadcrumb -->
						<div class="mt-1 flex items-center justify-between">
							<div class="flex items-center px-2">
								<Breadcrumb.Root>
									<Breadcrumb.List>
										<Breadcrumb.Item>
											<Breadcrumb.Link>{tabPostgresConnName}</Breadcrumb.Link>
										</Breadcrumb.Item>
										<Breadcrumb.Separator />
										<Breadcrumb.Item>
											<Breadcrumb.Link>{tabDBName}</Breadcrumb.Link>
										</Breadcrumb.Item>
										<Breadcrumb.Separator />
										<Breadcrumb.Item>
											<Breadcrumb.Page>{tabName}</Breadcrumb.Page>
										</Breadcrumb.Item>
									</Breadcrumb.List>
								</Breadcrumb.Root>
								{#if tabTableDBPoolID === ''}
									<Label class="ml-2 text-red-500">Disconnected</Label>
								{:else}
									<Label class="ml-2 text-green-500">Connected</Label>
								{/if}
							</div>
							<div class="flex px-2">
								<Tabs.Root value={tableViewTab}>
									<Tabs.List class="flex items-center justify-center gap-2">
										<Tabs.Trigger value="data" onclick={() => (tableViewTab = 'data')}
											>Data</Tabs.Trigger
										>
										<Tabs.Trigger value="structure" onclick={() => (tableViewTab = 'structure')}
											>Structure</Tabs.Trigger
										>
										<Tabs.Trigger value="indexes" onclick={() => (tableViewTab = 'indexes')}
											>Indexes</Tabs.Trigger
										>
									</Tabs.List>
								</Tabs.Root>
							</div>
						</div>

						<!-- Content based on selected tab -->
						{#if tableViewTab === 'data'}
							<div class="flex h-screen flex-1 flex-col">
								<div class="h-18 flex flex-col">
									<div class="flex flex-1 items-center gap-2 p-1">
										<Label for="select">Select</Label>
										<div class="multi-select-container relative w-full">
											<div
												class="border-input bg-background flex w-full cursor-text flex-wrap items-center gap-1 rounded-md border px-3 py-2 text-sm"
												role="combobox"
												tabindex="-1"
												aria-expanded={showSuggestions}
												aria-controls="suggestions-list"
												onclick={(e) => {
													// Focus the input when clicking anywhere in the container
													const input = e.currentTarget.querySelector(
														'.select-input'
													) as HTMLInputElement;
													input?.focus();
													showSuggestions = true;
												}}
											>
												<!-- Selected chips -->
												{#each selectedColumns as column, index}
													<div
														class="bg-secondary focus:ring-ring flex items-center gap-1 rounded-md px-2.5 py-0.5 text-xs transition-colors focus:outline-none {focusedChipIndex ===
														index
															? 'ring-ring ring-1'
															: ''}"
														tabindex="0"
														role="button"
														data-chip-index={index}
														onkeydown={(e) => handleChipKeyDown(e, index)}
														ondblclick={() => startEditingChip(index)}
														onclick={() => {
															focusedChipIndex = index;
														}}
													>
														{#if editingChipIndex === index}
															<input
																type="text"
																class="chip-edit-input w-20 border-none bg-transparent text-xs outline-none"
																bind:value={editingChipValue}
																onblur={finishEditingChip}
																onkeydown={(e) => {
																	if (e.key === 'Enter') {
																		e.preventDefault();
																		finishEditingChip();
																	} else if (e.key === 'Escape') {
																		e.preventDefault();
																		e.stopPropagation();
																		cancelEditingChip();
																	}
																}}
															/>
														{:else}
															{column}
														{/if}
														<button
															type="button"
															class="hover:bg-muted ml-1 rounded-full"
															onclick={(e) => {
																e.stopPropagation();
																removeColumn(column);
															}}
														>
															<X class="h-3 w-3" />
														</button>
													</div>
												{/each}
												<!-- Input field -->
												<input
													type="text"
													placeholder={selectedColumns.length === 0 ? 'Select columns...' : ''}
													class="select-input placeholder:text-muted-foreground flex-1 bg-transparent outline-none"
													bind:value={inputValue}
													oninput={(e) => handleInputChange(e.currentTarget.value)}
													onkeydown={handleInputKeyDown}
													onfocus={() => {
														showSuggestions = true;
														focusedChipIndex = -1;
													}}
												/>
											</div>

											<!-- Suggestions dropdown -->
											{#if showSuggestions && (filteredSuggestions.length > 0 || inputValue.trim())}
												<div
													id="suggestions-list"
													class="bg-popover absolute left-0 right-0 top-full z-50 mt-1 rounded-md border p-0 shadow-md"
												>
													{#if filteredSuggestions.length > 0}
														<div class="max-h-60 overflow-auto">
															{#each filteredSuggestions as suggestion, index}
																<button
																	type="button"
																	class="hover:bg-muted w-full px-3 py-2 text-left text-sm {selectedIndex ===
																	index
																		? 'bg-muted'
																		: ''}"
																	data-index={index}
																	onclick={() => selectColumn(suggestion)}
																	onmouseenter={() => {
																		selectedIndex = index;
																		inputValue = suggestion;
																	}}
																	onmouseleave={() => {
																		if (selectedIndex === index) {
																			inputValue = originalInputValue;
																			selectedIndex = -1;
																		}
																	}}
																>
																	{suggestion}
																</button>
															{/each}
														</div>
													{:else if inputValue.trim()}
														<div class="text-muted-foreground px-3 py-2 text-sm">
															No suggestions found
														</div>
													{/if}
												</div>
											{/if}
										</div>
									</div>
									<div class="flex flex-1 items-center gap-2 p-1">
										<Label for="where">Where</Label>
										<div class="multi-select-container relative w-full">
											<div
												class="border-input bg-background flex min-h-10 w-full cursor-text flex-wrap items-center gap-1 rounded-md border px-3 py-2 text-sm"
												role="combobox"
												tabindex="-1"
												aria-expanded={showWhereSuggestions}
												aria-controls="where-suggestions-list"
												onclick={(e) => {
													// Focus the input when clicking anywhere in the container
													const input = e.currentTarget.querySelector(
														'.where-input'
													) as HTMLInputElement;
													input?.focus();
													showWhereSuggestions = true;
												}}
											>
												<!-- Selected chips -->
												{#each selectedWhereColumns as column, index}
													<div
														class="bg-secondary flex items-center gap-1 rounded-md px-2.5 py-0.5 text-xs transition-colors focus:outline-none {focusedWhereChipIndex ===
														index
															? 'ring-ring ring-1'
															: ''}"
														tabindex="0"
														role="button"
														data-where-chip-index={index}
														onkeydown={(e) => handleWhereChipKeyDown(e, index)}
														ondblclick={() => startEditingWhereChip(index)}
														onclick={() => {
															focusedWhereChipIndex = index;
														}}
													>
														{#if editingWhereChipIndex === index}
															<input
																type="text"
																class="where-chip-edit-input w-20 border-none bg-transparent text-xs outline-none"
																bind:value={editingWhereChipValue}
																onblur={finishEditingWhereChip}
																onkeydown={(e) => {
																	if (e.key === 'Enter') {
																		e.preventDefault();
																		finishEditingWhereChip();
																	} else if (e.key === 'Escape') {
																		e.preventDefault();
																		e.stopPropagation();
																		cancelEditingWhereChip();
																	}
																}}
															/>
														{:else}
															{column}
														{/if}
														<button
															type="button"
															class="hover:bg-muted ml-1 rounded-full"
															onclick={(e) => {
																e.stopPropagation();
																removeWhereColumn(column);
															}}
														>
															<X class="h-3 w-3" />
														</button>
													</div>
												{/each}
												<!-- Input field -->
												<input
													type="text"
													placeholder={selectedWhereColumns.length === 0
														? 'Where conditions...'
														: ''}
													class="where-input placeholder:text-muted-foreground flex-1 bg-transparent outline-none"
													bind:value={whereInputValue}
													oninput={(e) => handleWhereInputChange(e.currentTarget.value)}
													onkeydown={handleWhereInputKeyDown}
													onfocus={() => {
														showWhereSuggestions = true;
														focusedWhereChipIndex = -1;
													}}
												/>
											</div>

											<!-- Suggestions dropdown -->
											{#if showWhereSuggestions && (filteredWhereSuggestions.length > 0 || whereInputValue.trim())}
												<div
													id="where-suggestions-list"
													class="bg-popover absolute left-0 right-0 top-full z-50 mt-1 rounded-md border p-0 shadow-md"
												>
													{#if filteredWhereSuggestions.length > 0}
														<div class="max-h-60 overflow-auto">
															{#each filteredWhereSuggestions as suggestion, index}
																<button
																	type="button"
																	class="hover:bg-muted w-full px-3 py-2 text-left text-sm {selectedWhereIndex ===
																	index
																		? 'bg-muted'
																		: ''}"
																	data-index={index}
																	onclick={() => selectWhereColumn(suggestion)}
																	onmouseenter={() => {
																		selectedWhereIndex = index;
																		whereInputValue = suggestion;
																	}}
																	onmouseleave={() => {
																		if (selectedWhereIndex === index) {
																			whereInputValue = originalWhereInputValue;
																			selectedWhereIndex = -1;
																		}
																	}}
																>
																	{suggestion}
																</button>
															{/each}
														</div>
													{:else if whereInputValue.trim()}
														<div class="text-muted-foreground px-3 py-2 text-sm">
															No suggestions found
														</div>
													{/if}
												</div>
											{/if}
										</div>
									</div>
									<div class="flex flex-1 items-center p-1 pt-0">
										<div class="flex flex-1 items-center gap-2 p-1">
											<Label for="orderBy">Order</Label>
											<Input
												type="text"
												id="orderBy"
												placeholder="Order By"
												class="w-full focus-visible:ring-0"
												bind:value={orderBy}
											/>
										</div>
										<div class="flex flex-1 items-center gap-2 p-1">
											<Label for="groupBy">Group</Label>
											<Input
												type="text"
												id="groupBy"
												placeholder="Group By"
												class="w-full focus-visible:ring-0"
												bind:value={groupBy}
											/>
										</div>
										<div class="flex items-center gap-2 p-1">
											<Label for="limit">Limit</Label>
											<Input
												type="text"
												id="limit"
												placeholder="Limit"
												class="w-24 focus-visible:ring-0"
												bind:value={limit}
											/>
										</div>
										<div class="flex items-center gap-2 p-1">
											<Label for="offset">Offset</Label>
											<Input
												type="text"
												id="offset"
												placeholder="Offset"
												class="w-24 focus-visible:ring-0"
												bind:value={offset}
											/>
										</div>
									</div>
								</div>
								<div class="flex flex-1 overflow-auto rounded-md border">
									<div class="h-full w-full overflow-auto">
										{#if $columns.length > 0}
											<DataTable
												data={$rows}
												columns={$columns}
												{queryLoading}
												query={$selectedQuery}
											/>
										{:else if queryLoading}
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
											<Skeleton class="my-3 h-[40px] w-full" />
										{/if}
									</div>
								</div>
							</div>
						{:else if tableViewTab === 'structure'}
							<div class="mt-2 flex flex-1 flex-col px-2">
								<p>Structure here</p>
							</div>
						{:else if tableViewTab === 'indexes'}
							<div class="mt-1 flex flex-1 flex-col px-2">
								<p>Indexes here</p>
							</div>
						{/if}
					</div>
				</div>
			{:else}
				<div class="flex h-screen flex-1 flex-col rounded-md">
					<Tabs.Content value={tabID.toString()} class="flex-1 overflow-hidden">
						<div class="flex h-full flex-col">
							<!-- Active DB Selector and Execute Query Button -->
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

							<!-- Resizable Panes for Editor and Output -->
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
										suggestions={Array.from($suggestions)}
									/>
								</Resizable.Pane>

								<Resizable.ResizableHandle withHandle />

								<!-- Output Pane -->
								<Resizable.Pane
									defaultSize={outputHeight}
									minSize={10}
									class="rsz-pane rounded-md border"
								>
									<div class="h-full">
										{#if $columns.length > 0}
											<DataTable
												data={$rows}
												columns={$columns}
												{queryLoading}
												query={$selectedQuery}
											/>
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
		{/if}
	</Tabs.Root>
</div>

<style>
	:global(.thin-scrollbar) {
		/* Firefox */
		scrollbar-width: thin;
		scrollbar-color: #4a5568 #2d3748;
	}

	:global(.thin-scrollbar::-webkit-scrollbar) {
		height: 1px; /* Horizontal scrollbar height */
		width: 1px; /* Vertical scrollbar width */
	}

	:global(.thin-scrollbar::-webkit-scrollbar-track) {
		background: #2d3748;
		border-radius: 1px;
	}

	:global(.thin-scrollbar::-webkit-scrollbar-thumb) {
		background: #4a5568;
		border-radius: 1px;
	}

	:global(.thin-scrollbar::-webkit-scrollbar-thumb:hover) {
		background: #718096;
	}
</style>
