<script lang="ts">
	import { Input } from '$lib/components/ui/command';
	import Label from '$lib/components/ui/label/label.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
</script>

<form>
	<div class="grid w-full items-center gap-4">
		<div class="flex flex-col space-y-1.5">
			<Label for="name">Name</Label>
			<Input id="name" placeholder="Name of your connection" />
			<Label for="environment">Environment</Label>
			<Input id="environment" placeholder="Local / Staging / Production" />
			<Label for="host">Host</Label>
			<Input id="host" placeholder="127.0.0.1" />
			<Label for="username">Username</Label>
			<Input id="username" placeholder="username" />
			<Label for="password">Password</Label>
			<Input id="password" placeholder="password" />
			<Label for="database">Database</Label>
			<Input id="database" placeholder="database name" />
		</div>
		<Button variant="destructive">Connect</Button>
	</div>
</form>
