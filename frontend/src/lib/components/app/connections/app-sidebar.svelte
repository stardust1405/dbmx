<script lang="ts" module>
	import { onMount } from 'svelte';
	import type { model } from '$lib/wailsjs/go/models';

	// This is sample data.
	const data = {
		changes: [
			{
				file: 'README.md',
				state: 'M'
			},
			{
				file: 'routes/+page.svelte',
				state: 'U'
			},
			{
				file: 'routes/+layout.svelte',
				state: 'M'
			}
		],
		tree: [
			// ['Postgres - local', 'inventorydb', 'reportingdb'],
			[
				{
					name: 'Stage Inventory',
					env: 'local',
					colour: 'purple'
				}
				// {
				// 	name: 'inventorydb',
				// 	env: 'local',
				// 	colour: 'purple'
				// },
				// {
				// 	name: 'reportingdb',
				// 	env: 'local',
				// 	colour: 'purple'
				// }
			],
			[
				{
					name: 'Prod Inventory',
					env: 'prod',
					colour: 'purple'
				}
				// {
				// 	name: 'inventorydb',
				// 	env: 'local',
				// 	colour: 'purple'
				// },
				// {
				// 	name: 'reportingdb',
				// 	env: 'local',
				// 	colour: 'purple'
				// }
			]
			// ['lib', ['components', 'button.svelte', 'card.svelte'], 'utils.ts'],
			// [
			// 	'routes',
			// 	['hello', '+page.svelte', '+page.ts'],
			// 	'+page.svelte',
			// 	'+page.server.ts',
			// 	'+layout.svelte'
			// ],
			// ['static', 'favicon.ico', 'svelte.svg']
		]
	};
</script>

<script lang="ts">
	import * as Collapsible from '$lib/components/ui/collapsible/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import ChevronRight from 'lucide-svelte/icons/chevron-right';
	import Plus from 'lucide-svelte/icons/plus';
	import Activity from 'lucide-svelte/icons/activity';
	import Table2 from 'lucide-svelte/icons/table-2';
	import RefreshCw from 'lucide-svelte/icons/refresh-cw';
	import SettingsDialog from '$lib/components/settings-dialog.svelte';
	import { type ComponentProps } from 'svelte';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	let postgresConnections: Array<model.PostgresConnection> = $state([]);

	import {
		EstablishPostgresConnection,
		GetPostgresConnections
	} from '$lib/wailsjs/go/app/Connections';
	import Button from '$lib/components/ui/button/button.svelte';

	onMount(() => {
		GetPostgresConnections().then((connections) => (postgresConnections = connections));
	});

	const refresh = () => {
		GetPostgresConnections().then((connections) => (postgresConnections = connections));
	};

	$effect(() => {
		console.log(databases);
	});

	let databases: Array<model.Database> = $state([]);

	let loading = $state(true);

	function establishConnection(id: number) {
		if (databases.length > 0) {
			return;
		}
		EstablishPostgresConnection(id).then((dbs) => {
			databases = dbs;
			loading = false;
		});
	}

	function getColorClass(color: string): string {
		const colorMap: Record<string, string> = {
			'bg-red-500': 'bg-red-500',
			'bg-blue-500': 'bg-blue-500',
			'bg-green-500': 'bg-green-500'
		};
		return colorMap[color] || '';
	}
</script>

<Sidebar.Root bind:ref {...restProps}>
	<Sidebar.Content>
		<!-- <Sidebar.Group>
			<Sidebar.GroupLabel>Favourites</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each data.changes as item, index (index)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton>
								<File />
								{item.file}
							</Sidebar.MenuButton>
							<Sidebar.MenuBadge>{item.state}</Sidebar.MenuBadge>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group> -->
		<Sidebar.Group>
			<Sidebar.GroupLabel>
				<div class="flex w-full items-center justify-between">
					Connections
					<Button size="sm" variant="ghost" onclick={() => refresh()}>
						<RefreshCw />
					</Button>
				</div>
			</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each postgresConnections as connection, index (index)}
						<Sidebar.MenuItem
							class="{getColorClass(connection.Colour)} bg-opacity-20 hover:bg-opacity-25"
						>
							<Collapsible.Root>
								<Collapsible.Trigger onclick={() => establishConnection(connection.ID)}>
									<Sidebar.MenuButton>
										<ChevronRight className="transition-transform" />

										{connection.Name}
									</Sidebar.MenuButton>
								</Collapsible.Trigger>
								<Collapsible.Content>
									<Sidebar.MenuSub>
										{#if loading}
											<p>Loading...</p>
										{:else if databases.length > 0}
											{#each databases as database, index}
												<Collapsible.Root>
													<Collapsible.Trigger>
														<Sidebar.MenuButton>
															<ChevronRight className="transition-transform" />
															{database.Name}
															{#if database.IsActive}
																<Activity color="#4fff4d" />
															{/if}
														</Sidebar.MenuButton>
													</Collapsible.Trigger>
													<Collapsible.Content>
														{#if database.Tables}
															{#each database.Tables as table, index}
																<Sidebar.MenuSub>
																	<Sidebar.MenuButton>
																		<Table2 color="#fd6868" strokeWidth={2} size={25} />
																		<p>{table}</p>
																	</Sidebar.MenuButton>
																</Sidebar.MenuSub>
															{/each}
														{:else}
															<Sidebar.MenuSub>No tables found.</Sidebar.MenuSub>
														{/if}
													</Collapsible.Content>
												</Collapsible.Root>
											{/each}
										{:else}
											<p>No databases found.</p>
										{/if}
									</Sidebar.MenuSub>
								</Collapsible.Content>
							</Collapsible.Root>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					<Plus />
					<SettingsDialog />
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Footer>
	<Sidebar.Rail />
</Sidebar.Root>

<!-- eslint-disable-next-line @typescript-eslint/no-explicit-any
{#snippet Tree({ item }: { item: string | any[] })}
	{@const [name, ...items] = Array.isArray(item) ? item : [item]}
	{#if !items.length}
		<Sidebar.MenuButton
			isActive={name === 'button.svelte'}
			class="data-[active=true]:bg-transparent"
		>
			<File />
			{name.name} - {name.env}
		</Sidebar.MenuButton>
	{:else}
		<Sidebar.MenuItem>
			<Collapsible.Root
				class="group/collapsible [&[data-state=open]>button>svg:first-child]:rotate-90"
				open={name === 'lib' || name === 'components'}
			>
				<Collapsible.Trigger>
					{#snippet child({ props })}
						<Sidebar.MenuButton {...props}>
							<ChevronRight className="transition-transform" />
							<Folder />
							{name.name} - {name.env}
						</Sidebar.MenuButton>
					{/snippet}
				</Collapsible.Trigger>
				<Collapsible.Content>
					<Sidebar.MenuSub>
						{#each items as subItem, index (index)}
							{@render Tree({ item: subItem })}
						{/each}
					</Sidebar.MenuSub>
				</Collapsible.Content>
			</Collapsible.Root>
		</Sidebar.MenuItem>
	{/if}
{/snippet} -->
