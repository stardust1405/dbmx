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
	import { SvelteMap } from 'svelte/reactivity';

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();

	// let postgresConnections: Array<model.PostgresConnection> = $state([]);
	let postgresConnectionsMap = $state<SvelteMap<number, model.PostgresConnection>>(
		new SvelteMap<number, model.PostgresConnection>()
	);
	let connectionDatabasesMap = $state<SvelteMap<number, string[]>>(
		new SvelteMap<number, string[]>()
	);
	let databasesMap = $state<SvelteMap<string, model.Database>>(
		new SvelteMap<string, model.Database>()
	);

	import {
		EstablishPostgresConnection,
		EstablishPostgresDatabaseConnection,
		GetPostgresConnections
	} from '$lib/wailsjs/go/app/Connections';
	import Button from '$lib/components/ui/button/button.svelte';
	import { toast } from 'svelte-sonner';
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';

	onMount(() => {
		GetPostgresConnections()
			.then((connections) => {
				for (let connection of connections) {
					postgresConnectionsMap.set(connection.ID, connection);
				}
			})
			.catch((error) => {
				// Handle errors from the EstablishPostgresDatabaseConnection call
				toast.error('Connection Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});
	});

	const refresh = () => {
		GetPostgresConnections()
			.then((connections) => {
				for (let connection of connections) {
					postgresConnectionsMap.set(connection.ID, connection);
				}
			})
			.catch((error) => {
				// Handle errors from the EstablishPostgresDatabaseConnection call
				toast.error('Connection Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});
	};

	let loadingMap = $state<SvelteMap<number, boolean>>(new SvelteMap<number, boolean>());
	let dbLoadingMap = $state<SvelteMap<string, boolean>>(new SvelteMap<string, boolean>());

	function establishDatabaseConnection(id: number, dbID: string) {
		dbLoadingMap.set(dbID, true);

		let db = databasesMap.get(dbID);
		// Check if connection is already established
		if (db?.IsActive) {
			dbLoadingMap.set(dbID, false);
			return;
		}

		// Establish connection
		EstablishPostgresDatabaseConnection(id, dbID, db?.Name || '')
			.then((db) => {
				dbLoadingMap.set(dbID, false);
				databasesMap.set(db.ID, db);
			})
			.catch((error) => {
				dbLoadingMap.set(dbID, false);

				toast.error('Connection Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});
	}

	function establishConnection(id: number) {
		loadingMap.set(id, true);
		// Check if connection is already established
		if ((connectionDatabasesMap.get(id)?.length || 0) > 0) {
			loadingMap.set(id, false);
			return;
		}

		// Establish connection
		EstablishPostgresConnection(id)
			.then((dbs) => {
				for (let db of dbs) {
					databasesMap.set(db.ID, db);
				}
				connectionDatabasesMap.set(
					id,
					dbs.map((db) => db.ID)
				);
				loadingMap.set(id, false);
			})
			.catch((error) => {
				loadingMap.set(id, false);
				// Handle errors from the EstablishPostgresDatabaseConnection call
				toast.error('Connection Failed', {
					description: error,
					action: {
						label: 'OK',
						onClick: () => console.info('OK')
					}
				});
			});
	}

	function getColorClass(color: string): string {
		const colorMap: Record<string, string> = {
			'bg-purple-500': 'bg-purple-500',
			'bg-indigo-500': 'bg-indigo-500',
			'bg-emerald-500': 'bg-emerald-500'
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
					{#each Array.from(postgresConnectionsMap.entries()) as [key, connection]}
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
										{#if loadingMap.get(connection.ID)}
											<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
											<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
											<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
										{:else if (connectionDatabasesMap.get(connection.ID) || []).length > 0}
											{#each connectionDatabasesMap.get(connection.ID) || [] as databaseID}
												<Collapsible.Root>
													<Collapsible.Trigger
														onclick={() => establishDatabaseConnection(connection.ID, databaseID)}
													>
														<Sidebar.MenuButton>
															<ChevronRight className="transition-transform" />
															{databasesMap.get(databaseID)?.Name}
															{#if databasesMap.get(databaseID)?.IsActive}
																<Activity color="#4fff4d" />
															{/if}
														</Sidebar.MenuButton>
													</Collapsible.Trigger>
													<Collapsible.Content>
														<Sidebar.MenuSub>
															{#if dbLoadingMap.get(databaseID)}
																<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
																<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
																<Skeleton class="my-1 h-[20px] w-[100px] rounded-full" />
															{:else if (databasesMap.get(databaseID)?.Tables || []).length > 0}
																{#each databasesMap.get(databaseID)?.Tables || [] as table}
																	<Sidebar.MenuButton>
																		<Table2 color="#fd6868" strokeWidth={2} size={25} />
																		<p>{table}</p>
																	</Sidebar.MenuButton>
																{/each}
															{:else}
																No tables found
															{/if}
														</Sidebar.MenuSub>
													</Collapsible.Content>
												</Collapsible.Root>
											{/each}
										{:else}
											No databases found
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
