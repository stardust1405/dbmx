import type { model } from "$lib/wailsjs/go/models";

export let activeDBs = $state<Array<model.Database>>([]);

export let suggestions = $state<Array<string>>([
    'SELECT',
    'FROM',
    'WHERE',
    'INNER',
    'LEFT',
    'RIGHT',
    'JOIN',
    'AND',
    'OR',
    'NOT',
    'IN',
    'LIKE',
    'BETWEEN',
    'ORDER BY',
    'LIMIT'
]);

// Handle DB selection
export function addActiveDB(db: model.Database) {
    activeDBs.push(db); // Directly mutate the array
    suggestions.push(...db.Tables);
}

export function removeActiveDB(dbID: string, tables: Array<string>) {
    activeDBs.splice(
        activeDBs.findIndex((db) => db.ID === dbID),
        1
    ); // Use splice to remove the item
    suggestions.splice(
        suggestions.findIndex((table) => tables.includes(table)),
        tables.length
    );
}
