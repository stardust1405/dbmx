<script lang="ts">
	import AppSidebar from '$lib/components/app/sidebar/sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Tabs from '$lib/components/app/main_screen/tabs.svelte';

	// Active tab properties
	let tabID = $state(0);
	let tabName = $state('');
	let tabType = $state('');

	// Reference to the Tabs component
	let tabsComponent: Tabs;

	// Function to handle adding a new tab from sidebar
	function handleAddTab(tableName?: string) {
		if (tabsComponent && tabsComponent.addTab) {
			tabsComponent.addTab(tableName);
		}
	}
</script>

<Sidebar.Provider>
	<AppSidebar bind:tabID bind:tabName onAddTab={handleAddTab} />
	<Sidebar.Inset>
		<Tabs bind:this={tabsComponent} bind:tabID bind:tabName bind:tabType />
	</Sidebar.Inset>
</Sidebar.Provider>
