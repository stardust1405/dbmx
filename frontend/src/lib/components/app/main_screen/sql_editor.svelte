<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';

	let editorContainer: HTMLElement;
	let editor: import('monaco-editor').editor.IStandaloneCodeEditor;
	let completionProviderDisposable: import('monaco-editor').IDisposable;
	let isInitialized = false;

	let { value = $bindable(), height = '100%', width = '100%', keywords = [] } = $props();

	// This effect will run whenever the value changes from outside
	$effect(() => {
		const _ = value;

		// Only update the editor if it exists and is initialized
		if (editor && isInitialized) {
			// Prevent infinite loops by checking if the editor value is different from the prop value
			const currentValue = editor.getValue();
			if (currentValue !== value) {
				// Update the editor model with the new value
				editor.setValue(value);
			}
		}
	});

	onMount(async () => {
		const monaco = await loader.init();

		monaco.languages.register({ id: 'sql' });

		completionProviderDisposable = monaco.languages.registerCompletionItemProvider('sql', {
			provideCompletionItems: (model, position) => {
				const word = model.getWordAtPosition(position);
				if (!word) {
					return { suggestions: [] };
				}

				const sqlKeywords = keywords;

				const prefix = word.word.toLowerCase();
				const suggestions = sqlKeywords.filter((keyword) =>
					keyword.toLowerCase().startsWith(prefix)
				);

				return {
					suggestions: suggestions.map((keyword) => ({
						label: keyword,
						kind: monaco.languages.CompletionItemKind.Keyword,
						insertText: keyword,
						range: {
							startLineNumber: position.lineNumber,
							startColumn: position.column - prefix.length,
							endLineNumber: position.lineNumber,
							endColumn: position.column
						}
					}))
				};
			}
		});

		editor = monaco.editor.create(editorContainer, {
			value: value,
			language: 'sql',
			theme: 'vs-dark',
			minimap: { enabled: false },
			fontSize: 14,
			automaticLayout: true
		});

		// Update value when editor content changes
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});

		// Manually trigger suggestions on key press
		editor.onKeyDown((e) => {
			if (e.keyCode === 32 || e.keyCode === 9) {
				// Space or Tab
				// Use the correct command for triggering suggestions
				editor.trigger('keyboard', 'editor.action.triggerSuggest', null);
			}
		});

		// Mark as initialized after editor is created
		isInitialized = true;
	});

	onDestroy(() => {
		editor?.dispose();
		completionProviderDisposable?.dispose();
	});
</script>

<div bind:this={editorContainer} class="sql-editor" style="height: {height}; width: {width};"></div>
