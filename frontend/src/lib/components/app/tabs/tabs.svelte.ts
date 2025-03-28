import type { model } from "$lib/wailsjs/go/models";

export let activeDBs = $state<Array<model.Database>>([]);

// Handle DB selection
export function addActiveDB(db: model.Database) {
    activeDBs.push(db); // Directly mutate the array
}

export function removeActiveDB(dbID: string) {
    activeDBs.splice(
        activeDBs.findIndex((db) => db.ID === dbID),
        1
    ); // Use splice to remove the item
}
