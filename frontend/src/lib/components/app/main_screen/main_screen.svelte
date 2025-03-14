<script lang="ts">
	import * as Tabs from '$lib/components/ui/tabs/index.js';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { Separator } from '$lib/components/ui/separator/index.js';
	import X from 'lucide-svelte/icons/x';
	import Plus from 'lucide-svelte/icons/plus';
	import * as Resizable from '$lib/components/ui/resizable/index.js';
	import { Textarea } from '$lib/components/ui/textarea/index.js';
	import { onMount } from 'svelte';
	import {
		Table,
		TableBody,
		TableCell,
		TableHead,
		TableHeader,
		TableRow
	} from '$lib/components/ui/table/index.js';
	interface TabData {
		label: string;
	}

	let tabs = $state<Record<string, TabData>>({
		account: { label: 'Account' },
		password: { label: 'Password' },
		abc1: { label: 'abc1' }
	});

	let currentTab = $state('account');

	function removeTab(id: string) {
		const newTabs = { ...tabs };
		delete newTabs[id];
		tabs = newTabs;
	}

	function addTab() {
		const newId = `tab-${Object.keys(tabs).length + 1}`;
		tabs = {
			...tabs,
			[newId]: { label: `Tab ${Object.keys(tabs).length + 1}` }
		};
		currentTab = newId;
	}

	let text = $state('');
	let lineCount = $state(1);
	let lineNumbersDiv: HTMLDivElement | null = null;
	let textarea: HTMLTextAreaElement | null = null;

	function updateLineNumbers(event: Event) {
		const target = event.target as HTMLTextAreaElement;
		const start = target.selectionStart;
		const end = target.selectionEnd;
		const value = target.value.substring(0, start);
		const after = target.value.substring(end);
		const lines = value.split('\n').length;
		lineCount = lines;
	}

	function syncScroll() {
		if (lineNumbersDiv && textarea) {
			lineNumbersDiv.scrollTop = textarea.scrollTop;
		}
	}

	const keywords = $state([
		'SELECT',
		'FROM',
		'WHERE',
		'JOIN',
		'GROUP BY',
		'ORDER BY',
		'LIMIT',
		'DELETE',
		'INNER JOIN',
		'OUTER JOIN',
		'INSERT',
		'UPDATE',
		'CREATE',
		'ALTER',
		'DROP',
		'TRUNCATE',
		'WITH',
		'AS',
		'ON',
		'USING',
		'HAVING',
		'UNION',
		'UNION ALL',
		'EXCEPT',
		'INTERSECT',
		'DISTINCT',
		'ALL',
		'CASE',
		'WHEN',
		'THEN',
		'ELSE',
		'END',
		'EXISTS',
		'BETWEEN',
		'IN',
		'LIKE',
		'ILIKE',
		'SIMILAR TO',
		'IS NULL',
		'IS NOT NULL',
		'AND',
		'OR',
		'NOT',
		'ANY',
		'SOME',
		'CAST',
		'COALESCE',
		'NULLIF',
		'GREATEST',
		'LEAST',
		'OVER',
		'PARTITION BY',
		'RANK',
		'ROW_NUMBER',
		'DENSE_RANK',
		'FIRST_VALUE',
		'LAST_VALUE',
		'LAG',
		'LEAD',
		'WINDOW',
		'FILTER',
		'ARRAY',
		'JSON',
		'JSONB',
		'ARRAY_AGG',
		'STRING_AGG',
		'UNNEST',
		'GENERATE_SERIES',
		'EXPLAIN',
		'ANALYZE',
		'VACUUM',
		'REINDEX',
		'CLUSTER',
		'CHECKPOINT',
		'BEGIN',
		'COMMIT',
		'ROLLBACK',
		'SAVEPOINT',
		'RELEASE',
		'SET',
		'SHOW',
		'RESET',
		'COPY',
		'PREPARE',
		'EXECUTE',
		'DEALLOCATE',
		'GRANT',
		'REVOKE',
		'CASCADE',
		'RESTRICT',
		'FOREIGN KEY',
		'PRIMARY KEY',
		'UNIQUE',
		'CHECK',
		'DEFAULT',
		'REFERENCES',
		'INDEX',
		'SEQUENCE',
		'VIEW',
		'MATERIALIZED VIEW',
		'FUNCTION',
		'PROCEDURE',
		'TRIGGER',
		'EVENT TRIGGER',
		'RULE',
		'TYPE',
		'DOMAIN',
		'SCHEMA',
		'EXTENSION',
		'TABLESPACE',
		'DATABASE',
		'ROLE',
		'USER',
		'GROUP',
		'TABLESAMPLE',
		'LATERAL',
		'RETURNING',
		'VALUES',
		'OFFSET',
		'FETCH',
		'FOR',
		'NOWAIT',
		'SKIP LOCKED',
		'CONSTRAINT',
		'TEMPORARY',
		'TEMP',
		'IF EXISTS',
		'IF NOT EXISTS',
		'ASC',
		'DESC',
		'NULLS FIRST',
		'NULLS LAST',
		'ONLY',
		'RECURSIVE',
		'CURRENT_DATE',
		'CURRENT_TIME',
		'CURRENT_TIMESTAMP',
		'LOCALTIME',
		'LOCALTIMESTAMP',
		'INTERVAL',
		'TRUE',
		'FALSE',
		'NULL',
		'ZONE',
		'AT TIME ZONE',
		'OVERLAPS',
		'EXTRACT',
		'DATE_PART',
		'TO_CHAR',
		'TO_DATE',
		'TO_TIMESTAMP',
		'AGE',
		'NOW',
		'CURRENT_CATALOG',
		'CURRENT_ROLE',
		'CURRENT_SCHEMA',
		'CURRENT_USER',
		'SESSION_USER',
		'VERSION',
		'DEFERRABLE',
		'NOT DEFERRABLE',
		'IMMEDIATE',
		'DEFERRED',
		'INITIALLY',
		'ENABLE',
		'DISABLE',
		'REPLICA',
		'IDENTITY',
		'GENERATED',
		'ALWAYS',
		'BY DEFAULT',
		'STORED',
		'WITHOUT',
		'OPTIONS',
		'INHERITS',
		'OWNER',
		'FILLFACTOR',
		'TOAST',
		'STORAGE',
		'WITH OIDS',
		'WITHOUT OIDS',
		'LOGGED',
		'UNLOGGED',
		'PERSISTENT',
		'UNPERSISTENT',
		'PARALLEL',
		'SERIALIZABLE',
		'REPEATABLE READ',
		'READ COMMITTED',
		'READ UNCOMMITTED',
		'ISOLATION',
		'LEVEL',
		'TRANSACTION',
		'READ WRITE',
		'READ ONLY',
		'WORK',
		'COMMITTED',
		'ABORT',
		'DISCARD',
		'PLPGSQL',
		'SECURITY',
		'INVOKER',
		'DEFINER',
		'VOLATILE',
		'STABLE',
		'IMMUTABLE',
		'LEAKPROOF',
		'COST',
		'ROWS',
		'SUPPORT',
		'SAFE',
		'FORCE',
		'NO FORCE',
		'VALID',
		'INVALID',
		'REPLACE',
		'CONCURRENTLY',
		'METHOD',
		'HASH',
		'BTREE',
		'GIST',
		'GIN',
		'SPGIST',
		'BRIN',
		'OPERATOR',
		'FAMILY',
		'CLASS',
		'ACCESS',
		'SHARE',
		'EXCLUSIVE',
		'NO KEY UPDATE',
		'KEY SHARE',
		'KEY',
		'NO KEY',
		'OF',
		'DO',
		'INSTEAD',
		'NOTHING',
		'RETURNS',
		'LANGUAGE',
		'CALLED',
		'ON NULL',
		'RETURN',
		'QUERY',
		'PERFORM',
		'RAISE',
		'ASSERT',
		'LOOP',
		'EXIT',
		'CONTINUE',
		'RETURN NEXT',
		'RETURN QUERY',
		'FOREACH',
		'SLICE',
		'EXCEPTION',
		'DECLARE',
		'IF',
		'ELSIF',
		'END IF',
		'WHILE'
	]);
</script>

<Resizable.PaneGroup direction="vertical">
	<Resizable.Pane minSize={10} defaultSize={70} class="flex flex-col">
		<Tabs.Content value={currentTab} class="flex flex-1 flex-col">
			<h2 class="mb-2 text-lg font-semibold">Editor</h2>
			<div class="relative flex-1 overflow-auto">
				<div class="relative h-full w-full">
					<div class="relative flex h-full w-full">
						<!-- Line Numbers -->
						<div class="relative flex h-full w-full">
							<!-- Line Numbers -->
							<div class="relative flex h-full w-full">
								<!-- Line Numbers -->
								<div
									bind:this={lineNumbersDiv}
									class="absolute left-0 top-0 min-h-full w-8 select-none overflow-auto bg-slate-800 p-2 pt-[1.2rem] text-right font-mono text-sm text-slate-400"
								>
									<div class="relative">
										{#each Array(lineCount) as _, i}
											<div class="h-[1.5rem] leading-[1.5rem]">{i + 1}</div>
										{/each}
									</div>
								</div>

								<!-- Textarea -->
								<textarea
									bind:this={textarea}
									class="min-h-full w-full resize-none overflow-auto rounded-md border border-slate-700 bg-slate-900 pl-12 pr-4 pt-[1.2rem] font-mono text-sm leading-[1.5rem] text-slate-100"
									rows={lineCount}
									placeholder="Enter your SQL query here..."
									bind:value={text}
									oninput={updateLineNumbers}
									onscroll={syncScroll}
									onkeydown={(e) => {
										if (e.key === 'Tab') {
											e.preventDefault();
											const target = e.target as HTMLTextAreaElement | null;
											if (!target) return;
											let start = target.selectionStart;
											let end = target.selectionEnd;
											const after = target.value.substring(end);
											const before = target.value.substring(0, start);

											// Get the current line content up to cursor position
											const lineStart = before.lastIndexOf('\n') + 1;
											const currentLine = before.substring(lineStart);

											// Get the word at cursor position
											const cursorPosInLine = start - lineStart;
											const words = currentLine.split(/\s+/);
											let wordStart = 0;
											let currentWord = '';
											for (const word of words) {
												if (
													cursorPosInLine >= wordStart &&
													cursorPosInLine <= wordStart + word.length
												) {
													currentWord = word;
													break;
												}
												wordStart += word.length + 1; // +1 for the space
											}

											if (currentWord.length > 0) {
												const matches = keywords.filter((k) =>
													k.toLowerCase().startsWith(currentWord.toLowerCase())
												);
												if (matches.length > 0) {
													if (matches.length > 1) {
														// Show dropdown if multiple matches
														const dropdown = document.createElement('div');
														dropdown.className =
															'absolute bg-slate-800 border border-slate-700 rounded-md p-1 z-50';
														dropdown.style.left = `${target.getBoundingClientRect().left + (start - lineStart) * 8}px`;
														dropdown.style.top = `${target.getBoundingClientRect().top + lineCount * 24}px`;
														dropdown.tabIndex = -1;

														let selectedIndex = 0;
														const options: HTMLDivElement[] = [];

														matches.forEach((match, index) => {
															const option = document.createElement('div');
															option.className = `px-2 py-1 hover:bg-slate-700 cursor-pointer ${
																index === selectedIndex ? 'bg-slate-700' : ''
															}`;
															option.textContent = match;
															option.onclick = () => {
																selectOption(match);
															};
															options.push(option);
															dropdown.appendChild(option);
														});

														const selectOption = (match: string) => {
															const newValue =
																before.substring(0, start - currentWord.length) + match + after;
															target.value = newValue;
															target.selectionStart = target.selectionEnd =
																start - currentWord.length + match.length;
															text = target.value;
															dropdown.remove();
														};

														const removeDropdown = () => {
															dropdown.remove();
															document.removeEventListener('click', handleOutsideClick);
														};

														const handleOutsideClick = (event: MouseEvent) => {
															if (!dropdown.contains(event.target as Node)) {
																removeDropdown();
															}
														};

														dropdown.onkeydown = (e) => {
															if (e.key === 'ArrowDown') {
																e.preventDefault();
																selectedIndex = (selectedIndex + 1) % matches.length;
																options.forEach((opt, i) => {
																	opt.classList.toggle('bg-slate-700', i === selectedIndex);
																});
															} else if (e.key === 'ArrowUp') {
																e.preventDefault();
																selectedIndex =
																	(selectedIndex - 1 + matches.length) % matches.length;
																options.forEach((opt, i) => {
																	opt.classList.toggle('bg-slate-700', i === selectedIndex);
																});
															} else if (e.key === 'Enter') {
																e.preventDefault();
																selectOption(matches[selectedIndex]);
															} else if (e.key === 'Escape') {
																e.preventDefault();
																removeDropdown();
															}
														};

														document.body.appendChild(dropdown);
														document.addEventListener('click', handleOutsideClick);
														dropdown.focus();
													} else {
														// Auto-complete if only one match
														const completion = matches[0];
														const newValue =
															before.substring(0, start - currentWord.length) + completion + after;
														target.value = newValue;
														target.selectionStart = target.selectionEnd =
															start - currentWord.length + completion.length;
														text = target.value;
													}
												}
											}
										}
									}}
								></textarea>
							</div>
						</div>
					</div>
				</div>
			</div></Tabs.Content
		>
	</Resizable.Pane>
	<Resizable.Handle />
	<Resizable.Pane minSize={10} defaultSize={30} class="flex flex-col">
		<Tabs.Content value={currentTab} class="flex flex-1 flex-col">
			<h2 class="mb-2 text-lg font-semibold">Table</h2>
			<div class="overflow-auto rounded-md border">
				<Table>
					<TableHeader>
						<TableRow>
							<TableHead class="w-[100px]">id</TableHead>
							<TableHead>name</TableHead>
							<TableHead>email</TableHead>
							<TableHead class="text-right">created_at</TableHead>
						</TableRow>
					</TableHeader>
					<TableBody>
						{#each Array(5) as _, i}
							<TableRow>
								<TableCell class="font-medium">{i + 1}</TableCell>
								<TableCell>User {i + 1}</TableCell>
								<TableCell>user{i + 1}@example.com</TableCell>
								<TableCell class="text-right"
									>2023-01-{i + 1 < 10 ? '0' + (i + 1) : i + 1}</TableCell
								>
							</TableRow>
						{/each}
					</TableBody>
				</Table>
			</div>
		</Tabs.Content>
	</Resizable.Pane>
</Resizable.PaneGroup>
