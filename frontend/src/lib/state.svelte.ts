import { SvelteMap } from 'svelte/reactivity';
import { model } from '$lib/wailsjs/go/models';

export let postgresConnectionsMap = $state(new SvelteMap<number, model.PostgresConnection>());
export let connectionDatabasesMap = $state(new SvelteMap<number, string[]>());
export let databasesMap = $state(new SvelteMap<string, model.Database>());
export let loadingMap = $state<SvelteMap<number, boolean>>(new SvelteMap<number, boolean>());
export let dbLoadingMap = $state<SvelteMap<string, boolean>>(new SvelteMap<string, boolean>());