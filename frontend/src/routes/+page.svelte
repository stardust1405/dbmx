<script lang="ts">
	import AppSidebar from '$lib/components/app/sidebar/sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Tabs from '$lib/components/app/main_screen/tabs.svelte';

	// Active tab properties
	let tabID = $state(0);
	let tabName = $state('');
	let tabType = $state('');
</script>

<Sidebar.Provider>
	<AppSidebar bind:tabID bind:tabName />
	<Sidebar.Inset>
		<Tabs bind:tabID bind:tabName bind:tabType />
	</Sidebar.Inset>
</Sidebar.Provider>
