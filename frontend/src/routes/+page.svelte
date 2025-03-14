<script lang="ts">
	import AppSidebar from '$lib/components/app/connections/app-sidebar.svelte';
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import Button from '$lib/components/ui/button/button.svelte';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Tabs from '$lib/components/app/tabs/tabs.svelte';
</script>

<Sidebar.Provider>
	<AppSidebar />
	<Sidebar.Inset>
		<Tabs />
	</Sidebar.Inset>
</Sidebar.Provider>
