<script lang="ts">
	import loader from '@monaco-editor/loader';
	import { onMount, onDestroy } from 'svelte';

	let editorContainer: HTMLElement;
	let editor: import('monaco-editor').editor.IStandaloneCodeEditor;
	let completionProviderDisposable: import('monaco-editor').IDisposable;
	let isInitialized = false;

	// Paste the Upstream Sunburst theme JSON here or import it
	const UpstreamSunburstTheme: import('monaco-editor').editor.IStandaloneThemeData = {
		base: 'vs-dark',
		inherit: true,
		rules: [
			{
				background: '000000',
				token: ''
			},
			{
				foreground: 'ffffff',
				background: '0f0f0f',
				token: 'text'
			},
			{
				background: '000000',
				token: 'source.ruby.rails.embedded.html'
			},
			{
				foreground: 'ffffff',
				background: '101010',
				token: 'text.html.ruby'
			},
			{
				foreground: 'ccff33',
				token: 'constant.numeric.ruby'
			},
			{
				foreground: 'ffffff',
				background: '000000',
				token: 'source'
			},
			{
				foreground: '9933cc',
				token: 'comment'
			},
			{
				foreground: '339999',
				token: 'constant'
			},
			{
				foreground: 'ff6600',
				token: 'keyword'
			},
			{
				foreground: 'edf8f9',
				token: 'keyword.preprocessor'
			},
			{
				foreground: 'ffffff',
				token: 'keyword.preprocessor directive'
			},
			{
				foreground: 'ffcc00',
				token: 'entity.name.function'
			},
			{
				foreground: 'ffcc00',
				token: 'storage.type.function.js'
			},
			{
				fontStyle: 'italic',
				token: 'variable.parameter'
			},
			{
				foreground: '772cb7',
				background: '070707',
				token: 'source comment.block'
			},
			{
				foreground: 'ffffff',
				token: 'variable.other'
			},
			{
				foreground: '999966',
				token: 'support.function.activerecord.rails'
			},
			{
				foreground: '66ff00',
				token: 'string'
			},
			{
				foreground: 'aaaaaa',
				token: 'string constant.character.escape'
			},
			{
				foreground: '000000',
				background: 'cccc33',
				token: 'string.interpolated'
			},
			{
				foreground: '44b4cc',
				token: 'string.regexp'
			},
			{
				foreground: 'cccc33',
				token: 'string.literal'
			},
			{
				foreground: '555555',
				token: 'string.interpolated constant.character.escape'
			},
			{
				fontStyle: 'underline',
				token: 'entity.name.class'
			},
			{
				fontStyle: 'underline',
				token: 'support.class.js'
			},
			{
				fontStyle: 'italic underline',
				token: 'entity.other.inherited-class'
			},
			{
				foreground: 'ff6600',
				token: 'meta.tag.inline.any.html'
			},
			{
				foreground: 'ff6600',
				token: 'meta.tag.block.any.html'
			},
			{
				foreground: '99cc99',
				fontStyle: 'italic',
				token: 'entity.other.attribute-name'
			},
			{
				foreground: 'dde93d',
				token: 'keyword.other'
			},
			{
				foreground: 'ff6600',
				token: 'meta.selector.css'
			},
			{
				foreground: 'ff6600',
				token: 'entity.other.attribute-name.pseudo-class.css'
			},
			{
				foreground: 'ff6600',
				token: 'entity.name.tag.wildcard.css'
			},
			{
				foreground: 'ff6600',
				token: 'entity.other.attribute-name.id.css'
			},
			{
				foreground: 'ff6600',
				token: 'entity.other.attribute-name.class.css'
			},
			{
				foreground: '999966',
				token: 'support.type.property-name.css'
			},
			{
				foreground: 'ffffff',
				token: 'keyword.other.unit.css'
			},
			{
				foreground: 'ffffff',
				token: 'constant.other.rgb-value.css'
			},
			{
				foreground: 'ffffff',
				token: 'constant.numeric.css'
			},
			{
				foreground: 'ffffff',
				token: 'support.function.event-handler.js'
			},
			{
				foreground: 'ffffff',
				token: 'keyword.operator.js'
			},
			{
				foreground: 'cccc66',
				token: 'keyword.control.js'
			},
			{
				foreground: 'ffffff',
				token: 'support.class.prototype.js'
			},
			{
				foreground: 'ff6600',
				token: 'object.property.function.prototype.js'
			}
		],
		colors: {
			'editor.foreground': '#FFFFFF',
			'editor.background': '#000000',
			'editor.selectionBackground': '#35493CE0',
			'editorCursor.foreground': '#FFFFFF',
			'editorWhitespace.foreground': '#404040'
		}
	};

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

		monaco.editor.defineTheme('sunburst', UpstreamSunburstTheme);

		editor = monaco.editor.create(editorContainer, {
			value: value,
			language: 'graphql',
			theme: 'sunburst',
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
