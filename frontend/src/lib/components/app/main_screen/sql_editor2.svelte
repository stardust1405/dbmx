<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';

	let editorContainer: HTMLElement;
	let editor: import('monaco-editor').editor.IStandaloneCodeEditor;
	let completionProviderDisposable: import('monaco-editor').IDisposable;

	const { value = $bindable(), height = '100%', width = '100%', keywords = [] } = $props();

	onMount(async () => {
		console.log('Initializing editor');
		const monaco = await loader.init();

		monaco.languages.register({ id: 'sql' });

		completionProviderDisposable = monaco.languages.registerCompletionItemProvider('sql', {
			provideCompletionItems: (model, position) => {
				console.log('Completion provider called');
				const word = model.getWordAtPosition(position);
				if (!word) {
					console.log('No word found');
					return { suggestions: [] };
				}

				const sqlKeywords = keywords;

				const prefix = word.word.toLowerCase();
				const suggestions = sqlKeywords.filter((keyword) =>
					keyword.toLowerCase().startsWith(prefix)
				);

				console.log('Suggestions:', suggestions);

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
			automaticLayout: true
		});

		// Manually trigger suggestions on key press
		editor.onKeyDown((e) => {
			if (e.keyCode === 32 || e.keyCode === 9) {
				// Space or Tab
				editor.trigger('editor.action.triggerSuggest', 'sql', null);
			}
		});
	});

	onDestroy(() => {
		console.log('Destroying editor');
		editor?.dispose();
		completionProviderDisposable?.dispose();
	});
</script>

<div bind:this={editorContainer} class="sql-editor" style="height: {height}; width: {width};"></div>
