<script lang="ts">
	import AppSidebar from '$lib/components/app/sidebar/sidebar.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Tabs from '$lib/components/app/main_screen/tabs.svelte';
</script>

<Sidebar.Provider>
	<AppSidebar />
	<Sidebar.Inset>
		<Tabs />
	</Sidebar.Inset>
</Sidebar.Provider>
