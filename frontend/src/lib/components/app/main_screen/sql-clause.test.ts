// Assertions for the clause-bar completion engine.
//
// The repo has no test runner, so this file is deliberately dependency-free and
// runs on Node's built-in type stripping (Node 22.6+):
//
//   node --experimental-strip-types \
//     src/lib/components/app/main_screen/sql-clause.test.ts
//
// Same rationale as app/table_schema_test.go: this code's whole output is SQL
// text and caret positions, which compiling cannot check.

import {
  completionAt, handlePairKey, handleBackspace, applySuggestion, scanQuotes, quoteIdent
} from './sql-clause.ts';

let pass = 0, fail = 0;
const eq = (name: string, got: unknown, want: unknown) => {
  const g = JSON.stringify(got), w = JSON.stringify(want);
  if (g === w) { pass++; } else { fail++; console.log(`FAIL ${name}\n   got  ${g}\n   want ${w}`); }
};

const cols = [
  { name: 'id', dataType: 'integer', isPrimaryKey: true, isNullable: false },
  { name: 'email', dataType: 'text', isNullable: false },
  { name: 'status', dataType: 'user_status', isNullable: false, enumValues: ['active', 'archived'] },
  { name: 'created_at', dataType: 'timestamptz', isNullable: true },
  { name: 'is_admin', dataType: 'boolean', isNullable: true },
  { name: 'Full Name', dataType: 'text', isNullable: true }
];
const labels = (t: string, c: number, k: any = 'where') =>
  completionAt(k, cols, t, c).items.slice(0, 4).map(i => i.label);

// --- caret-aware pair insertion (the old code appended to end of string) ---
const cursor = (t: string, s: number, e = s) => (k: string) => {
  const r = handlePairKey(k, t, s, e);
  return r ? r.text.slice(0, r.caret) + '|' + r.text.slice(r.caret) : 'PASSTHROUGH';
};
eq('( auto-closes at end', cursor('a = b', 5)('('), 'a = b(|)');
eq('( before word does not close', cursor('a = b', 4)('('), 'PASSTHROUGH');
eq('( mid-word does not close', cursor('abc', 1)('('), 'PASSTHROUGH');
eq("' auto-closes", cursor('status = ', 9)("'"), "status = '|'");
eq("' after word does not close", cursor("don", 3)("'"), 'PASSTHROUGH');
eq('wrap selection in parens', cursor('a or b', 0, 6)('('), '(|a or b)');
eq('wrap selection in quotes', cursor('abc', 0, 3)("'"), "'|abc'");
eq('type over close paren', cursor('f()', 2)(')'), 'f()|');
eq('bracket inside literal is text', cursor("x = 'a", 6)('('), 'PASSTHROUGH');

// --- backspace deletes both halves of an empty pair ---
const bs = (t: string, c: number) => { const r = handleBackspace(t, c, c); return r ? r.text : 'PASSTHROUGH'; };
eq('backspace empties ()', bs('f()', 2), 'f');
eq("backspace empties ''", bs("a=''", 3), 'a=');
eq('backspace with content kept', bs('f(x)', 3), 'PASSTHROUGH');

// --- quote scanning ---
eq('inside literal', scanQuotes("a = 'act", 8).state, 'single');
eq('escaped quote stays inside', scanQuotes("a = 'it''s", 10).state, 'single');
eq('closed literal', scanQuotes("a = 'x' and", 11).state, 'none');

// --- filtered, ranked completion ---
eq('prefix filters columns', labels('sta', 3), ['status']);
eq('subsequence match', labels('ca', 2), ['created_at']);
eq('is matches column then keywords', labels('is', 2), ['is_admin', 'is null', 'is not null']);
eq('order clause words', labels('de', 2, 'orderBy'), ['desc']);
eq('select aggregates', labels('cou', 3, 'select'), ['count(*)', 'count()']);

// --- value suggestions after an operator ---
eq('enum values after =', labels("status = ", 9), ["'active'", "'archived'", 'id', 'email']);
eq('booleans after =', labels('is_admin = ', 11).slice(0, 3), ['true', 'false', 'null']);
eq('enum filtered inside literal', labels("status = 'arc", 13), ["'archived'"]);
eq('IN () offers enums', labels('status in (', 11).slice(0, 2), ["'active'", "'archived'"]);
eq('no values without operator', labels('sta', 3).includes("'active'"), false);

// --- insertion replaces the word at the caret, never appends ---
const accept = (text: string, caret: number, pick = 0, clause: any = 'where') => {
  const c = completionAt(clause, cols, text, caret);
  const e = applySuggestion(text, c.range, c.items[pick]);
  return e.text.slice(0, e.caret) + '|' + e.text.slice(e.caret);
};
eq('replaces prefix mid-string', accept('sta = 1', 3), 'status| = 1');
eq('replaces prefix at caret', accept('id = 1 and ema', 14), 'id = 1 and email|');
eq('enum value inserted', accept('status = ', 9), "status = 'active'|");
eq('caret lands inside count()', accept('cou', 3, 1, 'select'), 'count(|)');

// --- identifiers that need quoting ---
eq('quotes name with space', quoteIdent('Full Name'), '"Full Name"');
eq('quotes reserved word', quoteIdent('order'), '"order"');
eq('leaves plain name alone', quoteIdent('created_at'), 'created_at');
eq('completion quotes on accept', accept('Full', 4, 0, 'select'), '"Full Name"|');

console.log(`\n${pass} passed, ${fail} failed`);
process.exit(fail ? 1 : 0);
