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
    'LIMIT',
    'IS',
    'true',
    'false',
    'NULL',
    'COUNT()',
    'SUM()',
    'AVG()',
    'MAX()',
    'MIN()',
    'UPDATE',
    'INSERT',
    'DELETE',
    'ALTER',
    'CREATE',
    'DROP',
    'TRUNCATE',
    'RENAME',
    'EXPLAIN',
    'EXPLAIN ANALYZE',
    'EXPLAIN QUERY PLAN',
    'EXPLAIN CACHE',
    'EXPLAIN CONFIG',
    'EXPLAIN DEBUG',
    'EXPLAIN TRIGGER',
    'SET'
]);

// Handle DB selection
export function addActiveDB(db: model.Database) {
    activeDBs.push(db); // Directly mutate the array
    suggestions.push(...db.Tables);
    suggestions.push(...db.Columns);
}

export function removeActiveDB(dbID: string, tables: Array<string>, columns: Array<string>) {
    activeDBs.splice(
        activeDBs.findIndex((db) => db.ID === dbID),
        1
    ); // Use splice to remove the item
    suggestions.splice(
        suggestions.findIndex((table) => tables.includes(table)),
        tables.length
    );
    suggestions.splice(
        suggestions.findIndex((column) => columns.includes(column)),
        columns.length
    );
}