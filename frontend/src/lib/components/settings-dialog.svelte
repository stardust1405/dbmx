<script lang="ts">
	import * as Breadcrumb from '$lib/components/ui/breadcrumb/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Dialog from '$lib/components/ui/dialog/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import Lock from 'lucide-svelte/icons/lock';
	import Label from './ui/label/label.svelte';
	import Input from './ui/input/input.svelte';
	import { AddPostgresConnection, TestConnectPostgres } from '$lib/wailsjs/go/app/Connections';
	import { model } from '$lib/wailsjs/go/models';
	import { toast } from 'svelte-sonner';

	const data = {
		nav: [
			{ name: 'PostgreSQL', icon: Lock },
			{ name: 'SQLite3', icon: Lock },
			{ name: 'MongoDB', icon: Lock },
			{ name: 'DynamoDB', icon: Lock }
		]
	};

	let open = $state(false);
	let selected_db = $state(data.nav[0].name);

	// Postgres connection data
	let name = $state('');
	let env = $state('');
	let host = $state('');
	let port = $state('');
	let username = $state('');
	let password = $state('');
	let database = $state('');

	let connectPostgres = async (
		name: String,
		env: String,
		host: String,
		port: String,
		username: String,
		password: String,
		database: String
	) => {
		let missingFields = [];
		if (name.trim() === '') {
			missingFields.push('name');
		}
		if (env.trim() === '') {
			missingFields.push('env');
		}
		if (host.trim() === '') {
			missingFields.push('host');
		}
		if (port.trim() === '') {
			missingFields.push('port');
		}
		if (username.trim() === '') {
			missingFields.push('username');
		}
		if (password.trim() === '') {
			missingFields.push('password');
		}

		if (missingFields.length > 0) {
			toast.error('Missing required field(s)', {
				description: `${missingFields.join(', ')}`,
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}

		const connectPostgresData = new model.Postgres({
			Name: name,
			Env: env,
			Host: host,
			Port: port,
			Username: username,
			Password: password,
			Database: database
		});
		try {
			// Await the result of the TestConnectPostgres call
			await AddPostgresConnection(connectPostgresData);

			// If successful, show success message
			toast.success('Success', {
				description: 'The connection was saved successfully',
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
		} catch (error) {
			let errorMessage = 'An error occurred while trying to connect.';

			// Check if error is an instance of Error
			if (error instanceof Error) {
				errorMessage = error.message; // Access message property safely
			} else if (typeof error === 'string') {
				errorMessage = error; // If it's a string, use it directly
			}
			// Handle errors from the TestConnectPostgres call
			toast.error('Connection Failed', {
				description: errorMessage,
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
		}
	};

	let testConnectPostgres = async (
		name: String,
		env: String,
		host: String,
		port: String,
		username: String,
		password: String,
		database: String
	) => {
		let missingFields = [];
		if (name.trim() === '') {
			missingFields.push('name');
		}
		if (env.trim() === '') {
			missingFields.push('env');
		}
		if (host.trim() === '') {
			missingFields.push('host');
		}
		if (port.trim() === '') {
			missingFields.push('port');
		}
		if (username.trim() === '') {
			missingFields.push('username');
		}
		if (password.trim() === '') {
			missingFields.push('password');
		}

		if (missingFields.length > 0) {
			toast.error('Missing required field(s)', {
				description: `${missingFields.join(', ')}`,
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
			return;
		}

		const testConnectPostgresData = new model.Postgres({
			Name: name,
			Env: env,
			Host: host,
			Port: port,
			Username: username,
			Password: password,
			Database: database
		});
		try {
			// Await the result of the TestConnectPostgres call
			await TestConnectPostgres(testConnectPostgresData);

			// If successful, show success message
			toast.success('Test Successful', {
				description: 'You can successfully connect using the given credentials',
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
		} catch (error) {
			let errorMessage = 'An error occurred while trying to connect.';

			// Check if error is an instance of Error
			if (error instanceof Error) {
				errorMessage = error.message; // Access message property safely
			} else if (typeof error === 'string') {
				errorMessage = error; // If it's a string, use it directly
			}
			// Handle errors from the TestConnectPostgres call
			toast.error('Connection Failed', {
				description: errorMessage,
				action: {
					label: 'OK',
					onClick: () => console.info('OK')
				}
			});
		}

		const result = await TestConnectPostgres(testConnectPostgresData);
	};
</script>

<Dialog.Root bind:open>
	<Dialog.Trigger>
		{#snippet child({ props })}
			<Button size="sm" {...props}>New Connection</Button>
		{/snippet}
	</Dialog.Trigger>
	<Dialog.Content class="overflow-hidden p-0 md:max-h-[500px] md:max-w-[700px] lg:max-w-[800px]">
		<Dialog.Title class="sr-only">Settings</Dialog.Title>
		<Dialog.Description class="sr-only">Customize your settings here.</Dialog.Description>
		<Sidebar.Provider class="items-start">
			<Sidebar.Root collapsible="none" class="hidden md:flex">
				<Sidebar.Content>
					<Sidebar.Group>
						<Sidebar.GroupContent>
							<Sidebar.Menu>
								{#each data.nav as item (item.name)}
									<Sidebar.MenuItem>
										<Sidebar.MenuButton
											isActive={item.name === 'Messages & media'}
											onclick={() => {
												selected_db = item.name;
											}}
										>
											{#snippet child({ props })}
												<a href="##" {...props}>
													<item.icon />
													<span>{item.name}</span>
												</a>
											{/snippet}
										</Sidebar.MenuButton>
									</Sidebar.MenuItem>
								{/each}
							</Sidebar.Menu>
						</Sidebar.GroupContent>
					</Sidebar.Group>
				</Sidebar.Content>
			</Sidebar.Root>
			<main class="flex h-[480px] flex-1 flex-col overflow-hidden">
				<header
					class="flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12"
				>
					<div class="flex items-center gap-2 px-4">
						<Breadcrumb.Root>
							<Breadcrumb.List>
								<Breadcrumb.Item class="hidden md:block">
									<Breadcrumb.Link href="#">New Connection</Breadcrumb.Link>
								</Breadcrumb.Item>
								<Breadcrumb.Separator class="hidden md:block" />
								<Breadcrumb.Item>
									<Breadcrumb.Page>{selected_db}</Breadcrumb.Page>
								</Breadcrumb.Item>
							</Breadcrumb.List>
						</Breadcrumb.Root>
					</div>
				</header>
				<div class="flex flex-1 flex-col gap-4 overflow-y-auto p-2 pt-0">
					<div class="">
						<form>
							<div class="grid w-full items-center gap-4">
								<div class="flex flex-col space-y-1.5">
									<Label for="name">Name</Label>
									<Input bind:value={name} id="name" placeholder="Name of your connection" />
									<Label for="environment">Environment</Label>
									<Input
										bind:value={env}
										id="environment"
										placeholder="Local / Staging / Production"
									/>
									<Label for="host">Host</Label>
									<Input bind:value={host} id="host" placeholder="127.0.0.1" />
									<Label for="port">Port</Label>
									<Input bind:value={port} id="port" placeholder="5432" />
									<Label for="username">Username</Label>
									<Input bind:value={username} id="username" placeholder="username" />
									<Label for="password">Password</Label>
									<Input bind:value={password} id="password" placeholder="password" />
									<Label for="database">Database</Label>
									<Input bind:value={database} id="database" placeholder="database name" />
								</div>
								<div class="flex justify-between space-x-4">
									<Button
										class="flex-1"
										onclick={() =>
											testConnectPostgres(
												name.trim(),
												env.trim(),
												host.trim(),
												port.trim(),
												username.trim(),
												password.trim(),
												database.trim()
											)}
										variant="secondary">Test</Button
									>
									<Button
										class="flex-1"
										onclick={() =>
											connectPostgres(
												name.trim(),
												env.trim(),
												host.trim(),
												port.trim(),
												username.trim(),
												password.trim(),
												database.trim()
											)}
										variant="destructive">Save Connection</Button
									>
								</div>
							</div>
						</form>
					</div>
				</div>
			</main>
		</Sidebar.Provider>
	</Dialog.Content>
</Dialog.Root>
