<script lang="ts">
	// import Monaco from 'svelte-monaco';
	import { createEventDispatcher } from 'svelte';
	import { onMount } from 'svelte';
	import type monaco from 'monaco-editor';

	const { id = '', value = '', height = '100%', width = '100%' } = $props();

	// Track the editor instance
	let editor: monaco.editor.IStandaloneCodeEditor | null = null;
	// Track if the editor has been initialized
	let isInitialized = $state(false);
	// Local copy of the value to handle updates
	let localValue = $state(value);

	// Update local value when prop value changes
	$effect(() => {
		console.log(`[${id}] Value prop changed:`, value);

		// Only update if the value has changed and we're not in the middle of editing
		if (value !== localValue) {
			localValue = value;

			// If the editor is initialized, update its content
			if (isInitialized && editor) {
				// const model = editor.getModel();
				// if (model) {
				// 	const currentValue = model.getValue();
				// 	if (currentValue !== value) {
				// 		console.log(`[${id}] Updating editor model value`);
				// 		model.setValue(value);
				// 	}
				// }
			}
		}
	});

	const dispatch = createEventDispatcher<{
		change: string;
		execute: string;
	}>();

	function handleChange(event: CustomEvent<string>) {
		// Update our local value
		localValue = event.detail;

		// Only dispatch if the value has actually changed
		if (value !== event.detail) {
			console.log(`[${id}] Dispatching change event:`, event.detail);
			dispatch('change', event.detail);
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		// Execute query on Ctrl+Enter or Cmd+Enter
		if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
			event.preventDefault();
			dispatch('execute', value);
		}
	}

	// PostgreSQL keywords for autocomplete suggestions
	const sqlKeywords = [
		// SQL commands
		'SELECT',
		'INSERT',
		'UPDATE',
		'DELETE',
		'CREATE',
		'ALTER',
		'DROP',
		'TRUNCATE',
		'TABLE',
		'VIEW',
		'INDEX',
		'SEQUENCE',
		'FUNCTION',
		'PROCEDURE',
		'TRIGGER',
		'DATABASE',
		'SCHEMA',
		'GRANT',
		'REVOKE',
		'COMMIT',
		'ROLLBACK',
		'SAVEPOINT',

		// Clauses
		'FROM',
		'WHERE',
		'GROUP BY',
		'HAVING',
		'ORDER BY',
		'LIMIT',
		'OFFSET',
		'DISTINCT',
		'UNION',
		'UNION ALL',
		'INTERSECT',
		'EXCEPT',

		// Joins
		'JOIN',
		'INNER JOIN',
		'LEFT JOIN',
		'RIGHT JOIN',
		'FULL JOIN',
		'CROSS JOIN',
		'NATURAL JOIN',
		'LEFT OUTER JOIN',
		'RIGHT OUTER JOIN',
		'FULL OUTER JOIN',

		// Operators
		'AND',
		'OR',
		'NOT',
		'IN',
		'BETWEEN',
		'LIKE',
		'ILIKE',
		'IS NULL',
		'IS NOT NULL',
		'EXISTS',
		'ANY',
		'ALL',
		'SOME',

		// Functions
		'COUNT',
		'SUM',
		'AVG',
		'MIN',
		'MAX',
		'COALESCE',
		'NULLIF',
		'CURRENT_DATE',
		'CURRENT_TIME',
		'CURRENT_TIMESTAMP',
		'NOW()',
		'EXTRACT',
		'DATE_PART',
		'TO_CHAR',
		'TO_DATE',
		'TO_TIMESTAMP',

		// Data types
		'INTEGER',
		'BIGINT',
		'SMALLINT',
		'DECIMAL',
		'NUMERIC',
		'REAL',
		'DOUBLE PRECISION',
		'CHAR',
		'VARCHAR',
		'TEXT',
		'BYTEA',
		'TIMESTAMP',
		'DATE',
		'TIME',
		'INTERVAL',
		'BOOLEAN',
		'ENUM',
		'UUID',
		'JSON',
		'JSONB',
		'ARRAY',
		'HSTORE'
	];

	// Register SQL autocomplete provider when the component mounts
	onMount(() => {
		if (typeof window !== 'undefined' && 'monaco' in window) {
			try {
				// Use type assertion to avoid TypeScript errors
				const monaco = (window as any).monaco;

				// Register custom SQL completion provider
				monaco.languages.registerCompletionItemProvider('sql', {
					provideCompletionItems: () => {
						return {
							suggestions: sqlKeywords.map((keyword) => ({
								label: keyword,
								kind: monaco.languages.CompletionItemKind.Keyword,
								insertText: keyword,
								range: null // Will be filled in by Monaco
							}))
						};
					}
				});
			} catch (error) {
				console.error('Error registering SQL completion provider:', error);
			}
		}
	});
</script>

<svelte:window on:keydown={handleKeyDown} />

<!-- The Monaco editor itself manages focus and keyboard interactions -->
<div class="h-full w-full" style="height: {height}; width: {width};" data-editor-id={id}>
	<!-- <Monaco
		theme="vs-dark"
		value={localValue}
		on:change={handleChange}
		on:mount={(e) => {
			editor = e.detail;
			isInitialized = true;
			console.log(`[${id}] Editor mounted with value:`, localValue);
		}}
		options={{
			language: 'sql',
			// Editor appearance
			minimap: { enabled: false },
			fontFamily: 'monospace',
			fontSize: 14,
			lineNumbers: 'on',
			folding: true,
			wordWrap: 'on',

			// Editor behavior
			autoIndent: 'full',
			formatOnPaste: true,
			formatOnType: true,
			scrollBeyondLastLine: false,
			smoothScrolling: true,
			cursorBlinking: 'smooth',
			cursorStyle: 'line',

			// SQL-specific options
			quickSuggestions: true,
			suggestOnTriggerCharacters: true,
			parameterHints: { enabled: true },
			tabSize: 2
		}}
	/> -->
</div>

<style>
	:global(.monaco-editor) {
		height: 100% !important;
		width: 100% !important;
	}
</style>
