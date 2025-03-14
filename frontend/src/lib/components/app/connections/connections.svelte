<script lang="ts">
	import { onMount } from 'svelte';
	import type { model } from '$lib/wailsjs/go/models';
	import { Button } from '$lib/components/ui/button/index.js';

	import { GetPostgresConnections } from '$lib/wailsjs/go/app/Connections';

	let postgresConnections = $state<Array<model.Postgres>>();

	onMount(() => {
		GetPostgresConnections().then((connections) => (postgresConnections = connections));
	});

	var refresh = () => {
		GetPostgresConnections().then((connections) => (postgresConnections = connections));
	};
</script>

<div class="flex w-1/3 flex-col">
	<div class="flex flex-row justify-center">
		<div class="flex self-center">
			<h1 class="m-2">Connections</h1>
		</div>
		<div class="flex self-center">
			<Button size="sm" onclick={() => refresh()}>Refresh</Button>
		</div>
	</div>
	<div class="h-full flex-col overflow-auto">
		<!-- each connection row -->
		{#if postgresConnections}
			{#each postgresConnections as connection, i (i)}
				<div class="mx-8 mt-4 flex h-16 flex-row justify-between">
					<div class="ml-10 flex flex-col justify-center">
						<p class="text-base">{connection.Name}</p>
						<p class="text-sm">Postgres - {connection.Env}</p>
					</div>
					<div class="flex w-24 items-center justify-center">
						<Button size="sm" class="self-center">Connect</Button>
					</div>
				</div>
			{/each}
		{/if}
	</div>
</div>
