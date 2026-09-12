// Completion + bracket-matching engine for the table data clause bar.
//
// Everything here is a pure function of (text, caret) so the Svelte component
// stays a thin shell around it: the editing rules are the part worth getting
// right, and they are far easier to reason about without a DOM attached.

export type ClauseKind = 'select' | 'where' | 'orderBy' | 'groupBy';

export type SuggestionKind = 'column' | 'keyword' | 'function' | 'value';

export interface ColumnInfo {
	name: string;
	dataType?: string;
	isPrimaryKey?: boolean;
	isNullable?: boolean;
	enumValues?: string[];
}

export interface Suggestion {
	label: string;
	insert: string;
	kind: SuggestionKind;
	detail?: string;
	badge?: string;
	/** Moves the caret back from the end of `insert` — used to land inside `count(|)`. */
	caretOffset?: number;
}

export interface Edit {
	text: string;
	caret: number;
	caretEnd?: number;
}

export interface Completion {
	/** Slice of the text the accepted suggestion replaces. */
	range: { start: number; end: number };
	prefix: string;
	items: Suggestion[];
	/** True when the popup should open with no prefix typed (e.g. just after `status =`). */
	eager: boolean;
}

const WORD = /[A-Za-z0-9_$]/;

// Postgres words that force an identifier to be double-quoted.
const RESERVED = new Set([
	'all', 'analyse', 'analyze', 'and', 'any', 'array', 'as', 'asc', 'asymmetric', 'both', 'case',
	'cast', 'check', 'collate', 'column', 'constraint', 'create', 'current_catalog', 'current_date',
	'current_role', 'current_time', 'current_timestamp', 'current_user', 'default', 'deferrable',
	'desc', 'distinct', 'do', 'else', 'end', 'except', 'false', 'fetch', 'for', 'foreign', 'from',
	'grant', 'group', 'having', 'in', 'initially', 'intersect', 'into', 'lateral', 'leading', 'limit',
	'localtime', 'localtimestamp', 'not', 'null', 'offset', 'on', 'only', 'or', 'order', 'placing',
	'primary', 'references', 'returning', 'select', 'session_user', 'some', 'symmetric', 'table',
	'then', 'to', 'trailing', 'true', 'union', 'unique', 'user', 'using', 'variadic', 'when', 'where',
	'window', 'with'
]);

const SAFE_IDENT = /^[a-z_][a-z0-9_$]*$/;

export function quoteIdent(name: string): string {
	if (SAFE_IDENT.test(name) && !RESERVED.has(name)) return name;
	return '"' + name.replace(/"/g, '""') + '"';
}

export function unquoteIdent(raw: string): string {
	const trimmed = raw.trim();
	if (trimmed.startsWith('"') && trimmed.endsWith('"') && trimmed.length > 1) {
		return trimmed.slice(1, -1).replace(/""/g, '"');
	}
	return trimmed;
}

function quoteLiteral(value: string): string {
	return "'" + value.replace(/'/g, "''") + "'";
}

/**
 * Where the caret sits relative to string literals and quoted identifiers.
 * `openIndex` is the position of the quote that opened the run we're inside.
 */
export function scanQuotes(
	text: string,
	caret: number
): { state: 'none' | 'single' | 'double'; openIndex: number } {
	let state: 'none' | 'single' | 'double' = 'none';
	let openIndex = -1;
	let i = 0;
	while (i < caret) {
		const ch = text[i];
		if (state === 'none') {
			if (ch === "'") {
				state = 'single';
				openIndex = i;
			} else if (ch === '"') {
				state = 'double';
				openIndex = i;
			}
			i++;
			continue;
		}
		const quote = state === 'single' ? "'" : '"';
		if (ch === quote) {
			// A doubled quote is an escape, not a terminator.
			if (text[i + 1] === quote) {
				i += 2;
				continue;
			}
			state = 'none';
			openIndex = -1;
		}
		i++;
	}
	return { state, openIndex };
}

export function wordRangeAt(text: string, caret: number): { start: number; word: string } {
	let start = caret;
	while (start > 0 && WORD.test(text[start - 1])) start--;
	return { start, word: text.slice(start, caret) };
}

// ---------------------------------------------------------------------------
// Bracket and quote pairing
// ---------------------------------------------------------------------------

const PAIRS: Record<string, string> = { '(': ')', '[': ']', "'": "'", '"': '"' };
const CLOSERS = new Set([')', ']', "'", '"']);

/**
 * Monaco-style pairing. Returns the edit to apply, or null to let the browser
 * insert the character normally.
 */
export function handlePairKey(key: string, text: string, start: number, end: number): Edit | null {
	const close = PAIRS[key];

	// Typing an opener with a selection wraps the selection instead of replacing it.
	if (close && start !== end) {
		const selected = text.slice(start, end);
		return {
			text: text.slice(0, start) + key + selected + close + text.slice(end),
			caret: start + 1,
			caretEnd: end + 1
		};
	}

	if (start !== end) return null;

	const prev = start > 0 ? text[start - 1] : '';
	const next = start < text.length ? text[start] : '';
	const isQuote = key === "'" || key === '"';
	const { state } = scanQuotes(text, start);
	const insideThisQuote =
		(key === "'" && state === 'single') || (key === '"' && state === 'double');

	// Type over a closer that is already sitting under the caret.
	if (CLOSERS.has(key) && next === key) {
		if (isQuote && !insideThisQuote) {
			// Opening a fresh literal that happens to abut one — fall through to open it.
		} else {
			return { text, caret: start + 1 };
		}
	}

	if (!close) return null;

	if (isQuote) {
		// Closing the literal we're inside: a single character, no pairing.
		if (insideThisQuote) return null;
		// Never auto-close against a word — avoids turning don|t into don''t.
		if (WORD.test(prev) || prev === key) return null;
		if (WORD.test(next)) return null;
	} else {
		// Inside a string literal a bracket is just text.
		if (state !== 'none') return null;
		if (WORD.test(next)) return null;
	}

	return { text: text.slice(0, start) + key + close + text.slice(start), caret: start + 1 };
}

/** Backspace between the two halves of an empty pair removes both. */
export function handleBackspace(text: string, start: number, end: number): Edit | null {
	if (start !== end || start === 0) return null;
	const prev = text[start - 1];
	const close = PAIRS[prev];
	if (close && text[start] === close) {
		return { text: text.slice(0, start - 1) + text.slice(start + 1), caret: start - 1 };
	}
	return null;
}

// ---------------------------------------------------------------------------
// Vocabulary
// ---------------------------------------------------------------------------

function keyword(label: string, insert = label, caretOffset = 0): Suggestion {
	return { label, insert, kind: 'keyword', caretOffset };
}

function fn(label: string, insert: string, caretOffset: number, detail: string): Suggestion {
	return { label, insert, kind: 'function', detail, caretOffset };
}

const CLAUSE_WORDS: Record<ClauseKind, Suggestion[]> = {
	select: [
		keyword('*'),
		keyword('distinct'),
		keyword('as'),
		fn('count(*)', 'count(*)', 0, 'aggregate'),
		fn('count()', 'count()', -1, 'aggregate'),
		fn('sum()', 'sum()', -1, 'aggregate'),
		fn('avg()', 'avg()', -1, 'aggregate'),
		fn('min()', 'min()', -1, 'aggregate'),
		fn('max()', 'max()', -1, 'aggregate'),
		fn('coalesce()', 'coalesce()', -1, 'function'),
		fn('nullif()', 'nullif()', -1, 'function'),
		fn('lower()', 'lower()', -1, 'function'),
		fn('upper()', 'upper()', -1, 'function'),
		fn('length()', 'length()', -1, 'function'),
		fn('now()', 'now()', 0, 'function')
	],
	where: [
		keyword('and'),
		keyword('or'),
		keyword('not'),
		keyword('is null'),
		keyword('is not null'),
		keyword('in ()', 'in ()', -1),
		keyword('not in ()', 'not in ()', -1),
		keyword('between'),
		keyword('like'),
		keyword('ilike'),
		keyword('similar to'),
		keyword('true'),
		keyword('false'),
		keyword('null'),
		fn('lower()', 'lower()', -1, 'function'),
		fn('upper()', 'upper()', -1, 'function')
	],
	orderBy: [
		keyword('asc'),
		keyword('desc'),
		keyword('nulls first'),
		keyword('nulls last')
	],
	groupBy: [
		fn('rollup()', 'rollup()', -1, 'grouping'),
		fn('cube()', 'cube()', -1, 'grouping'),
		fn('grouping sets ()', 'grouping sets ()', -1, 'grouping')
	]
};

// Identifier, then a comparison operator, then optional whitespace/open paren.
const VALUE_OP =
	/(?:^|[\s(,])((?:"(?:[^"]|"")+")|[A-Za-z_][\w$]*)\s*(=|<>|!=|>=|<=|>|<|(?:not\s+)?i?like|(?:not\s+)?in|is(?:\s+not)?)\s*\(?\s*$/i;

function valueContext(head: string): { column: string; op: string } | null {
	const match = VALUE_OP.exec(head);
	if (!match) return null;
	return { column: unquoteIdent(match[1]), op: match[2].toLowerCase().replace(/\s+/g, ' ') };
}

function columnSuggestion(column: ColumnInfo): Suggestion {
	const badges: string[] = [];
	if (column.isPrimaryKey) badges.push('key');
	if (column.enumValues && column.enumValues.length > 0) badges.push('enum');
	return {
		label: column.name,
		insert: quoteIdent(column.name),
		kind: 'column',
		detail: column.dataType,
		badge: badges.join(' ') || undefined
	};
}

function valueSuggestions(column: ColumnInfo | undefined, op: string): Suggestion[] {
	const items: Suggestion[] = [];
	if (!column) return items;

	for (const value of column.enumValues ?? []) {
		items.push({
			label: quoteLiteral(value),
			insert: quoteLiteral(value),
			kind: 'value',
			detail: 'enum value'
		});
	}

	const type = (column.dataType ?? '').toLowerCase();
	if (type.includes('bool')) {
		items.push({ label: 'true', insert: 'true', kind: 'value', detail: column.dataType });
		items.push({ label: 'false', insert: 'false', kind: 'value', detail: column.dataType });
	}
	if (op.startsWith('is')) {
		items.push({ label: 'null', insert: 'null', kind: 'value' });
		items.push({ label: 'not null', insert: 'not null', kind: 'value' });
	} else if (column.isNullable) {
		items.push({ label: 'null', insert: 'null', kind: 'value' });
	}
	if (type.includes('date') || type.includes('time')) {
		items.push({ label: 'now()', insert: 'now()', kind: 'value', detail: 'current timestamp' });
	}
	return items;
}

// ---------------------------------------------------------------------------
// Ranking
// ---------------------------------------------------------------------------

const KIND_RANK: Record<SuggestionKind, number> = { value: 0, column: 1, function: 2, keyword: 3 };

/** Lower is better; null drops the candidate. */
function score(label: string, prefix: string): number | null {
	if (!prefix) return 0;
	const haystack = label.toLowerCase();
	const needle = prefix.toLowerCase();
	const at = haystack.indexOf(needle);
	if (at === 0) return 0;
	if (at > 0) return /[\s_('"]/.test(haystack[at - 1]) ? 1 : 2;

	// Fall back to a subsequence match so `ca` still finds `created_at`.
	let i = 0;
	for (const ch of haystack) {
		if (ch === needle[i]) i++;
		if (i === needle.length) break;
	}
	return i === needle.length ? 3 : null;
}

function rank(items: Suggestion[], prefix: string, valueFirst: boolean): Suggestion[] {
	// Ties fall back to the order candidates were built in: table order for
	// columns, authored order for keywords. That keeps count(*) above count()
	// and true/false above null, which sorting by label length does not.
	const scored: Array<{ item: Suggestion; score: number; index: number }> = [];
	items.forEach((item, index) => {
		const s = score(item.label, prefix);
		if (s === null) return;
		scored.push({ item, score: s, index });
	});

	const weight = (kind: SuggestionKind) =>
		valueFirst ? KIND_RANK[kind] : KIND_RANK[kind] + (kind === 'value' ? 10 : 0);

	scored.sort((a, b) => {
		if (a.score !== b.score) return a.score - b.score;
		const ra = weight(a.item.kind);
		const rb = weight(b.item.kind);
		if (ra !== rb) return ra - rb;
		return a.index - b.index;
	});
	return scored.map((entry) => entry.item);
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

const EMPTY: Completion = { range: { start: 0, end: 0 }, prefix: '', items: [], eager: false };

export function completionAt(
	clause: ClauseKind,
	columns: ColumnInfo[],
	text: string,
	caret: number
): Completion {
	const byName = new Map(columns.map((c) => [c.name.toLowerCase(), c]));
	const { state, openIndex } = scanQuotes(text, caret);

	// Inside a string literal: only offer values for the column being compared.
	if (state === 'single') {
		const head = text.slice(0, openIndex);
		const ctx = valueContext(head);
		if (!ctx) return EMPTY;
		const prefix = text.slice(openIndex + 1, caret);
		const items = valueSuggestions(byName.get(ctx.column.toLowerCase()), ctx.op).filter(
			(item) => item.kind === 'value' && item.insert.startsWith("'")
		);
		return {
			range: { start: openIndex, end: caret },
			prefix,
			items: rank(items, prefix, true),
			eager: true
		};
	}

	// Inside a quoted identifier: complete column names, re-quoting on accept.
	if (state === 'double') {
		const prefix = text.slice(openIndex + 1, caret);
		const items = columns.map(columnSuggestion);
		return {
			range: { start: openIndex, end: caret },
			prefix,
			items: rank(items, prefix, false),
			eager: true
		};
	}

	const { start, word } = wordRangeAt(text, caret);
	const head = text.slice(0, start);
	const ctx = clause === 'where' ? valueContext(head) : null;

	if (ctx) {
		const column = byName.get(ctx.column.toLowerCase());
		const items = [...valueSuggestions(column, ctx.op), ...columns.map(columnSuggestion)];
		return {
			range: { start, end: caret },
			prefix: word,
			items: rank(items, word, true),
			eager: true
		};
	}

	const items = [...columns.map(columnSuggestion), ...CLAUSE_WORDS[clause]];
	return {
		range: { start, end: caret },
		prefix: word,
		items: rank(items, word, false),
		eager: false
	};
}

/** Applies an accepted suggestion, returning the new text and caret. */
export function applySuggestion(
	text: string,
	range: { start: number; end: number },
	suggestion: Suggestion
): Edit {
	const next = text.slice(0, range.start) + suggestion.insert + text.slice(range.end);
	return {
		text: next,
		caret: range.start + suggestion.insert.length + (suggestion.caretOffset ?? 0)
	};
}

/** Splits a label around the matched prefix so the popup can highlight it. */
export function highlight(label: string, prefix: string): Array<{ text: string; hit: boolean }> {
	if (!prefix) return [{ text: label, hit: false }];
	const at = label.toLowerCase().indexOf(prefix.toLowerCase());
	if (at < 0) return [{ text: label, hit: false }];
	const parts = [
		{ text: label.slice(0, at), hit: false },
		{ text: label.slice(at, at + prefix.length), hit: true },
		{ text: label.slice(at + prefix.length), hit: false }
	];
	return parts.filter((part) => part.text.length > 0);
}
