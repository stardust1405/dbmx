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
	import QueryOutput from '$lib/components/app/main_screen/query_output.svelte';

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

	// Store the previous tab id
	let previousTabId = $state('');
	let editorContent = $state('select * from users\nselect * from users');
</script>

<Tabs.Root value={'tab 1'}>
	<!-- Tabs visible in the header -->
	<header class="flex h-14 shrink-0 items-center gap-2 overflow-auto border-b px-4">
		<Sidebar.Trigger class="-ml-1" />
		<Separator orientation="vertical" />
		<Tabs.List>
			<!-- {#each Object.entries(tabs) as [id, tab]} -->
			<div class="mr-2 flex rounded-sm bg-slate-800">
				<Tabs.Trigger value={'tab 1'} class="flex items-center">
					{'tab 1'}
				</Tabs.Trigger>
				<button class="rounded-r-sm bg-slate-900 px-2 py-1 text-slate-300 hover:text-red-700">
					<X size={16} />
				</button>
			</div>
			<!-- {/each} -->
			<button class="ml-2 flex items-center gap-1 text-blue-500 hover:text-blue-700">
				<Plus size={16} /> Add Tab
			</button>
		</Tabs.List>
	</header>

	<!-- Main Content on screen -->
	<div class="flex h-screen flex-1 flex-col gap-4 p-4">
		<Tabs.Content value={'tab 1'} class="flex-1 overflow-hidden">
			<div class="flex h-full flex-col">
				<div class="mb-2 flex items-center justify-between">
					<h2 class="text-lg font-semibold">{'tab 1'}</h2>
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
						<SqlEditor id={'tab 1'} value={editorContent} />
					</Resizable.ResizablePane>

					<Resizable.ResizableHandle />

					<!-- Output Pane -->
					<Resizable.ResizablePane
						defaultSize={outputHeight}
						minSize={10}
						class="rsz-pane overflow-auto"
					>
						<!-- <QueryOutput outputs={tab.outputData} /> -->
					</Resizable.ResizablePane>
				</Resizable.ResizablePaneGroup>
			</div>
		</Tabs.Content>
	</div>
</Tabs.Root>
