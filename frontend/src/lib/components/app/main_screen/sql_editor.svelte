<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';

	let editorContainer: HTMLElement;
	let editor: import('monaco-editor').editor.IStandaloneCodeEditor;
	let completionProviderDisposable: import('monaco-editor').IDisposable;
	let isInitialized = false;

	let {
		value = $bindable(),
		selectedQuery = $bindable(),
		height = '100%',
		width = '100%',
		suggestions = $bindable([])
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

		// On initialization
		let decorationCollection = editor.createDecorationsCollection([]);

		// Select the query on click and highlight it
		editor.onMouseDown((e) => {
			const position = editor.getPosition();
			if (!position) return;

			let currentLineNumber = position.lineNumber;
			let queryStartLineNumber = currentLineNumber;
			let queryEndLineNumber = currentLineNumber;

			// Move to the start of the query
			while (true) {
				if (currentLineNumber === 1) {
					queryStartLineNumber = 1;
					break;
				}

				const lineContent = model.getLineContent(currentLineNumber - 1).trim();
				if (lineContent === '') {
					queryStartLineNumber = currentLineNumber;
					break;
				}

				currentLineNumber--;
			}

			// Move to the end of the query
			while (true) {
				if (currentLineNumber === model.getLineCount()) {
					queryEndLineNumber = model.getLineCount();
					break;
				}

				const lineContent = model.getLineContent(currentLineNumber + 1).trim();
				if (lineContent === '') {
					queryEndLineNumber = currentLineNumber;
					break;
				}

				currentLineNumber++;
			}

			const selection = new monaco.Range(
				queryStartLineNumber,
				1,
				queryEndLineNumber,
				model.getLineLength(queryEndLineNumber) + 1
			);

			// On each mouse click
			decorationCollection.set([
				{
					range: selection,
					options: {
						isWholeLine: true,
						className: 'bg-green-100 bg-opacity-5',
						glyphMarginClassName: 'bg-green-500 bg-opacity-20'
					}
				}
			]);

			// Get content from the selection
			selectedQuery = model.getValueInRange(selection).trim();
		});

		// Select the query when cursor is placed on it using arrow keys
		editor.onKeyUp((e) => {
			const position = editor.getPosition();
			if (!position) return;

			let currentLineNumber = position.lineNumber;
			let queryStartLineNumber = currentLineNumber;
			let queryEndLineNumber = currentLineNumber;

			// Move to the start of the query
			while (true) {
				if (currentLineNumber === 1) {
					queryStartLineNumber = 1;
					break;
				}

				const lineContent = model.getLineContent(currentLineNumber - 1).trim();
				if (lineContent === '') {
					queryStartLineNumber = currentLineNumber;
					break;
				}

				currentLineNumber--;
			}

			// Move to the end of the query
			while (true) {
				if (currentLineNumber === model.getLineCount()) {
					queryEndLineNumber = model.getLineCount();
					break;
				}

				const lineContent = model.getLineContent(currentLineNumber + 1).trim();
				if (lineContent === '') {
					queryEndLineNumber = currentLineNumber;
					break;
				}

				currentLineNumber++;
			}

			const selection = new monaco.Range(
				queryStartLineNumber,
				1,
				queryEndLineNumber,
				model.getLineLength(queryEndLineNumber) + 1
			);

			// On each mouse click
			decorationCollection.set([
				{
					range: selection,
					options: {
						isWholeLine: true,
						className: 'bg-green-100 bg-opacity-5',
						glyphMarginClassName: 'bg-green-500 bg-opacity-20'
					}
				}
			]);

			selectedQuery = model.getValueInRange(selection).trim();
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
			selectedQuery = selection.toString().trim();
		}
	}
</script>

<div
	onselect={handleSelection}
	bind:this={editorContainer}
	class="sql-editor"
	style="height: {height}; width: {width};"
></div>
