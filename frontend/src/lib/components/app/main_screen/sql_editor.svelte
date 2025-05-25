<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';

	let editorContainer: HTMLElement;
	let editor: import('monaco-editor').editor.IStandaloneCodeEditor;
	let completionProviderDisposable: import('monaco-editor').IDisposable;
	let isInitialized = false;

	let {
		value = $bindable(),
		selectedText = $bindable(),
		height = '100%',
		width = '100%',
		suggestions = []
	} = $props();

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

				const prefix = word.word.toLowerCase();
				const filteredSuggestions = suggestions.filter((keyword) =>
					keyword.toLowerCase().startsWith(prefix)
				);

				return {
					suggestions: filteredSuggestions.map((keyword) => ({
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
			automaticLayout: true,
			glyphMargin: true
		});

		const model = editor.getModel();
		if (!model) return;

		editor.onMouseDown((e) => {
			const position = editor.getPosition();
			if (!position) return;

			let currentLineContent = model.getLineContent(position.lineNumber).trim();
			selectedText = currentLineContent;
		});

		editor.onKeyUp((e) => {
			const position = editor.getPosition();
			if (!position) return;

			let currentLineContent = model.getLineContent(position.lineNumber).trim();
			selectedText = currentLineContent;
		});

		// Update value when editor content changes
		editor.onDidChangeModelContent(() => {
			value = editor.getValue();
		});

		// editor.createDecorationsCollection([
		// 	{
		// 		range: new monaco.Range(1, 1, 1, 1),
		// 		options: {
		// 			isWholeLine: true,
		// 			className: 'bg-green-100 bg-opacity-5',
		// 			glyphMarginClassName: 'bg-green-500 bg-opacity-20'
		// 		}
		// 	},
		// 	{
		// 		range: new monaco.Range(3, 1, 6, 1),
		// 		options: {
		// 			isWholeLine: true,
		// 			className: 'bg-green-100 bg-opacity-5',
		// 			glyphMarginClassName: 'bg-green-500 bg-opacity-20'
		// 		}
		// 	},
		// 	{
		// 		range: new monaco.Range(8, 1, 11, 1),
		// 		options: {
		// 			isWholeLine: true,
		// 			className: 'bg-green-100 bg-opacity-5',
		// 			glyphMarginClassName: 'bg-green-500 bg-opacity-20'
		// 		}
		// 	}
		// ]);

		// Mark as initialized after editor is created
		isInitialized = true;
	});

	onDestroy(() => {
		editor?.dispose();
		completionProviderDisposable?.dispose();
	});

	// Function to update the state with the selected text
	function handleSelection() {
		const selection = window.getSelection();
		if (selection && selection.toString().trim() !== '') {
			selectedText = selection.toString().trim();
		}
	}
</script>

<div
	onselect={handleSelection}
	bind:this={editorContainer}
	class="sql-editor"
	style="height: {height}; width: {width};"
></div>
