// The Monaco themes shared by the SQL editor and the read-only table-definition viewer.
// They live here rather than inside a component so both register the identical palette:
// Monaco's theme registry is global to the loader's single Monaco instance, so whichever
// component mounts first defines them for the other.
import type { editor } from 'monaco-editor';

export const AuroraSQLTheme: editor.IStandaloneThemeData = {
	base: 'vs-dark',
	inherit: true,
	rules: [
		// base / comments
		{ token: '', foreground: 'E6E9EF', background: '000000' },
		{ token: 'comment', foreground: '7D8696', fontStyle: 'italic' },
		{ token: 'comment.sql', foreground: '7D8696', fontStyle: 'italic' },

		// keywords (SELECT, WHERE, JOIN, LIMIT…)
		{ token: 'keyword', foreground: '8CAAEE' }, // soft azure
		{ token: 'keyword.sql', foreground: '8CAAEE' },

		// types (INT, VARCHAR…), NULL/TRUE/FALSE
		{ token: 'type', foreground: '8BD5CA' }, // teal
		{ token: 'type.sql', foreground: '8BD5CA' },
		{ token: 'predefined', foreground: 'F2CDCD' }, // NULL/TRUE/FALSE (rose)
		{ token: 'predefined.sql', foreground: 'F2CDCD' },

		// strings & numbers
		{ token: 'string', foreground: 'A6E3A1' }, // mint
		{ token: 'string.sql', foreground: 'A6E3A1' },
		{ token: 'string.escape', foreground: 'B5E8E0' },
		{ token: 'number', foreground: 'F5A97F' }, // coral/peach
		{ token: 'number.sql', foreground: 'F5A97F' },

		// functions (COUNT, SUM, MAX…)
		{ token: 'entity.name.function', foreground: 'EED49F' }, // amber
		{ token: 'support.function', foreground: 'EED49F' },

		// identifiers (tables/columns), quoted names
		{ token: 'identifier', foreground: 'DCE1EA' },
		{ token: 'identifier.sql', foreground: 'DCE1EA' },
		{ token: 'identifier.quote', foreground: 'AAB4CF' },

		// operators / punctuation
		{ token: 'operator', foreground: 'BAC2DE' },
		{ token: 'operator.sql', foreground: 'BAC2DE' },
		{ token: 'delimiter', foreground: '6C7086' },
		{ token: 'delimiter.sql', foreground: '6C7086' },
		{ token: 'delimiter.bracket', foreground: 'B4BEFE' },

		// errors (muted)
		{ token: 'invalid', foreground: 'F38BA8' },
		{ token: 'invalid.deprecated', foreground: 'F9E2AF' }
	],
	colors: {
		// canvas & text
		'editor.background': '#000000',
		'editor.foreground': '#E6E9EF',

		// gutter / guides
		'editorLineNumber.foreground': '#2A3340',
		'editorLineNumber.activeForeground': '#9DB2CE',
		'editorIndentGuide.background': '#121212',
		'editorIndentGuide.activeBackground': '#1E1E1E',
		'editorWhitespace.foreground': '#202020',
		editorLineHighlightBackground: '#0A0C10',

		// caret / selection / word highlights
		'editorCursor.foreground': '#E6E9EF',
		'editor.selectionBackground': '#1E2A3A',
		'editor.inactiveSelectionBackground': '#10151C',
		'editor.wordHighlightBackground': '#0B122066',
		'editor.wordHighlightStrongBackground': '#0B122099',

		// matches / search
		'editor.findMatchBackground': '#22D3EE66',
		'editor.findMatchHighlightBackground': '#6EE7B780',
		'editor.findRangeHighlightBackground': '#37415166',

		// brackets
		'editorBracketMatch.background': '#101318',
		'editorBracketMatch.border': '#2A3646',

		// autocomplete
		'editorSuggestWidget.background': '#0B0D11',
		'editorSuggestWidget.border': '#1F2633',
		'editorSuggestWidget.foreground': '#D7DBE2',
		'editorSuggestWidget.selectedBackground': '#151B24',
		'editorSuggestWidget.highlightForeground': '#8CAAEE',

		// hover / peek
		'editorHoverWidget.background': '#0B0D11',
		'editorHoverWidget.border': '#1F2633',

		// scrollbar
		'scrollbarSlider.background': '#2A2A2AAA',
		'scrollbarSlider.hoverBackground': '#3A3A3AAA',
		'scrollbarSlider.activeBackground': '#4A4A4AAA',

		// markers
		'editorError.foreground': '#F38BA8',
		'editorWarning.foreground': '#EED49F',
		'editorInfo.foreground': '#8CAAEE'
	}
};

export const AuroraSQLLightTheme: editor.IStandaloneThemeData = {
	base: 'vs',
	inherit: true,
	rules: [
		{ token: '', foreground: '1E293B', background: 'FFFFFF' },
		{ token: 'comment', foreground: '94A3B8', fontStyle: 'italic' },
		{ token: 'comment.sql', foreground: '94A3B8', fontStyle: 'italic' },
		{ token: 'keyword', foreground: '2563EB' },
		{ token: 'keyword.sql', foreground: '2563EB' },
		{ token: 'type', foreground: '0D9488' },
		{ token: 'type.sql', foreground: '0D9488' },
		{ token: 'predefined', foreground: 'BE185D' },
		{ token: 'predefined.sql', foreground: 'BE185D' },
		{ token: 'string', foreground: '16A34A' },
		{ token: 'string.sql', foreground: '16A34A' },
		{ token: 'string.escape', foreground: '0D9488' },
		{ token: 'number', foreground: 'EA580C' },
		{ token: 'number.sql', foreground: 'EA580C' },
		{ token: 'entity.name.function', foreground: 'D97706' },
		{ token: 'support.function', foreground: 'D97706' },
		{ token: 'identifier', foreground: '334155' },
		{ token: 'identifier.sql', foreground: '334155' },
		{ token: 'identifier.quote', foreground: '64748B' },
		{ token: 'operator', foreground: '475569' },
		{ token: 'operator.sql', foreground: '475569' },
		{ token: 'delimiter', foreground: '94A3B8' },
		{ token: 'delimiter.sql', foreground: '94A3B8' },
		{ token: 'delimiter.bracket', foreground: '7C3AED' },
		{ token: 'invalid', foreground: 'DC2626' },
		{ token: 'invalid.deprecated', foreground: 'D97706' }
	],
	colors: {
		'editor.background': '#FFFFFF',
		'editor.foreground': '#1E293B',
		'editorLineNumber.foreground': '#CBD5E1',
		'editorLineNumber.activeForeground': '#475569',
		'editorIndentGuide.background': '#F1F5F9',
		'editorIndentGuide.activeBackground': '#E2E8F0',
		'editorWhitespace.foreground': '#F1F5F9',
		editorLineHighlightBackground: '#F8FAFC',
		'editorCursor.foreground': '#1E293B',
		'editor.selectionBackground': '#DBEAFE',
		'editor.inactiveSelectionBackground': '#EFF6FF',
		'editor.wordHighlightBackground': '#DBEAFE66',
		'editor.wordHighlightStrongBackground': '#DBEAFE99',
		'editor.findMatchBackground': '#FDE68A',
		'editor.findMatchHighlightBackground': '#BBF7D080',
		'editor.findRangeHighlightBackground': '#F1F5F966',
		'editorBracketMatch.background': '#EDE9FE',
		'editorBracketMatch.border': '#C4B5FD',
		'editorSuggestWidget.background': '#FFFFFF',
		'editorSuggestWidget.border': '#E2E8F0',
		'editorSuggestWidget.foreground': '#334155',
		'editorSuggestWidget.selectedBackground': '#F1F5F9',
		'editorSuggestWidget.highlightForeground': '#2563EB',
		'editorHoverWidget.background': '#FFFFFF',
		'editorHoverWidget.border': '#E2E8F0',
		'scrollbarSlider.background': '#CBD5E1AA',
		'scrollbarSlider.hoverBackground': '#94A3B8AA',
		'scrollbarSlider.activeBackground': '#64748BAA',
		'editorError.foreground': '#DC2626',
		'editorWarning.foreground': '#D97706',
		'editorInfo.foreground': '#2563EB'
	}
};

/** Registers both themes on a Monaco instance. Safe to call more than once. */
export function defineSQLThemes(monaco: typeof import('monaco-editor')): void {
	monaco.editor.defineTheme('aurora-sql', AuroraSQLTheme);
	monaco.editor.defineTheme('aurora-sql-light', AuroraSQLLightTheme);
}

/** The theme name matching the app's current light/dark mode. */
export function sqlThemeFor(currentMode: string | undefined): string {
	return currentMode === 'light' ? 'aurora-sql-light' : 'aurora-sql';
}
