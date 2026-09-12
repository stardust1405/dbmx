<script lang="ts">
	import { tick } from 'svelte';
	import { cn } from '$lib/utils.js';
	import {
		applySuggestion,
		completionAt,
		handleBackspace,
		handlePairKey,
		highlight,
		type ClauseKind,
		type ColumnInfo,
		type Completion,
		type Suggestion
	} from './sql-clause';

	let {
		value = $bindable(''),
		clause,
		columns = [],
		label,
		placeholder = '',
		id,
		onrun
	}: {
		value?: string;
		clause: ClauseKind;
		columns?: ColumnInfo[];
		label: string;
		placeholder?: string;
		id: string;
		onrun?: () => void;
	} = $props();

	let inputEl = $state<HTMLInputElement | null>(null);
	let listEl = $state<HTMLDivElement | null>(null);
	let open = $state(false);
	let active = $state(0);
	let completion = $state<Completion | null>(null);
	let box = $state({ left: 0, top: 0, width: 0, flip: false });

	const listId = $derived(`${id}-suggestions`);
	const items = $derived(completion?.items ?? []);
	const prefix = $derived(completion?.prefix ?? '');

	/** Appends the popup to <body> so resizable panes can't clip it. */
	function portal(node: HTMLElement) {
		document.body.appendChild(node);
		return {
			destroy() {
				node.remove();
			}
		};
	}

	function place() {
		if (!inputEl) return;
		const rect = inputEl.getBoundingClientRect();
		const below = window.innerHeight - rect.bottom;
		const flip = below < 220 && rect.top > below;
		box = {
			left: rect.left,
			top: flip ? rect.top - 6 : rect.bottom + 6,
			width: Math.max(rect.width, 280),
			flip
		};
	}

	/** Recomputes candidates from the live caret position. */
	function refresh(force = false) {
		if (!inputEl) return;
		const caret = inputEl.selectionStart ?? value.length;
		const next = completionAt(clause, columns, value, caret);
		completion = next;
		active = 0;

		if (next.items.length === 0) {
			open = false;
			return;
		}
		if (force || next.prefix.length > 0 || next.eager) {
			open = true;
			place();
		}
	}

	function close() {
		open = false;
	}

	async function commit(edit: { text: string; caret: number; caretEnd?: number }) {
		value = edit.text;
		await tick();
		inputEl?.setSelectionRange(edit.caret, edit.caretEnd ?? edit.caret);
	}

	async function accept(suggestion: Suggestion) {
		if (!completion) return;
		const edit = applySuggestion(value, completion.range, suggestion);
		open = false;
		completion = null;
		await commit(edit);
		inputEl?.focus();
		// Deliberately no re-query here: the accepted word would read as the next
		// prefix and reopen the popup on the item just chosen. Typing reopens it,
		// and typing `'` after `status =` lands straight on the enum values.
	}

	function move(delta: number) {
		if (items.length === 0) return;
		active = (active + delta + items.length) % items.length;
	}

	async function onkeydown(event: KeyboardEvent) {
		if (!inputEl) return;

		if (event.key === 'Escape' && open) {
			event.preventDefault();
			event.stopPropagation();
			close();
			return;
		}

		if (event.key === ' ' && (event.ctrlKey || event.metaKey)) {
			event.preventDefault();
			refresh(true);
			return;
		}

		if (open && items.length > 0) {
			if (event.key === 'ArrowDown') {
				event.preventDefault();
				move(1);
				return;
			}
			if (event.key === 'ArrowUp') {
				event.preventDefault();
				move(-1);
				return;
			}
			if (event.key === 'Home' && !event.shiftKey) {
				event.preventDefault();
				active = 0;
				return;
			}
			if (event.key === 'End' && !event.shiftKey) {
				event.preventDefault();
				active = items.length - 1;
				return;
			}
			if (event.key === 'Enter' || event.key === 'Tab') {
				event.preventDefault();
				accept(items[active]);
				return;
			}
		}

		if (event.key === 'ArrowDown' && !open) {
			event.preventDefault();
			refresh(true);
			return;
		}

		if (event.key === 'Enter') {
			event.preventDefault();
			close();
			onrun?.();
			return;
		}

		const start = inputEl.selectionStart ?? 0;
		const end = inputEl.selectionEnd ?? start;

		if (event.key === 'Backspace') {
			const edit = handleBackspace(value, start, end);
			if (edit) {
				event.preventDefault();
				await commit(edit);
				refresh(false);
			}
			return;
		}

		if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
			const edit = handlePairKey(event.key, value, start, end);
			if (edit) {
				event.preventDefault();
				await commit(edit);
				refresh(false);
			}
		}
	}

	function oninput() {
		refresh(false);
	}

	/** Arrow/Home/End and clicks move the caret without changing the text. */
	function syncCaret(event: KeyboardEvent) {
		if (!open) return;
		if (['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key)) refresh(false);
	}

	$effect(() => {
		if (!open || !listEl) return;
		const el = listEl.querySelector<HTMLElement>(`[data-index="${active}"]`);
		el?.scrollIntoView({ block: 'nearest' });
	});

	$effect(() => {
		if (!open) return;
		const reposition = () => place();
		window.addEventListener('resize', reposition);
		window.addEventListener('scroll', reposition, true);
		return () => {
			window.removeEventListener('resize', reposition);
			window.removeEventListener('scroll', reposition, true);
		};
	});

	const KIND_GLYPH: Record<Suggestion['kind'], string> = {
		column: '■',
		value: '◆',
		function: 'ƒ',
		keyword: '·'
	};
</script>

<div class="clause-cell">
	<label class="clause-label" for={id}>{label}</label>
	<input
		{id}
		bind:this={inputEl}
		bind:value
		type="text"
		{placeholder}
		spellcheck="false"
		autocomplete="off"
		autocapitalize="off"
		autocorrect="off"
		class="clause-input"
		role="combobox"
		aria-expanded={open}
		aria-controls={listId}
		aria-autocomplete="list"
		aria-activedescendant={open && items.length > 0 ? `${listId}-${active}` : undefined}
		{oninput}
		{onkeydown}
		onkeyup={syncCaret}
		onclick={() => open && refresh(false)}
		onblur={close}
	/>
</div>

{#if open && items.length > 0}
	<div
		use:portal
		id={listId}
		role="listbox"
		aria-label="{label} suggestions"
		tabindex="-1"
		class="clause-popup"
		style:left="{box.left}px"
		style:top="{box.top}px"
		style:width="{box.width}px"
		style:transform={box.flip ? 'translateY(-100%)' : 'none'}
		bind:this={listEl}
		onmousedown={(event) => event.preventDefault()}
	>
		{#each items as item, index (item.kind + item.label)}
			<!-- Keyboard access lives on the input (aria-activedescendant); options are pointer-only. -->
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<div
				id="{listId}-{index}"
				data-index={index}
				role="option"
				aria-selected={index === active}
				tabindex="-1"
				class={cn('clause-option', index === active && 'is-active')}
				onclick={() => accept(item)}
				onmouseenter={() => (active = index)}
			>
				<span class="clause-glyph" aria-hidden="true">{KIND_GLYPH[item.kind]}</span>
				<span class="clause-name">
					{#each highlight(item.label, prefix) as part}
						<span class={part.hit ? 'clause-hit' : ''}>{part.text}</span>
					{/each}
				</span>
				{#if item.badge}
					<span class="clause-badge">{item.badge}</span>
				{/if}
				{#if item.detail}
					<span class="clause-detail">{item.detail}</span>
				{/if}
			</div>
		{/each}
	</div>
{/if}

<style>
	.clause-cell {
		display: flex;
		align-items: stretch;
		min-width: 0;
	}

	.clause-label {
		flex: 0 0 4.25rem;
		display: flex;
		align-items: center;
		justify-content: flex-end;
		padding-inline: 0.625rem;
		font-size: 0.75rem;
		line-height: 1;
		color: hsl(var(--muted-foreground));
		border-right: 1px solid hsl(var(--border));
		user-select: none;
		cursor: text;
	}

	.clause-input {
		flex: 1 1 auto;
		min-width: 0;
		background: transparent;
		border: 0;
		outline: 0;
		padding: 0.4rem 0.625rem;
		font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
		font-size: 0.8125rem;
		color: hsl(var(--foreground));
	}

	.clause-input::placeholder {
		color: hsl(var(--muted-foreground) / 0.55);
	}

	/* The focused cell is marked by its gutter and ground, not a ring around the box. */
	.clause-cell:focus-within {
		background: hsl(var(--accent) / 0.55);
	}

	.clause-cell:focus-within .clause-label {
		color: hsl(var(--foreground));
		border-right-color: hsl(var(--foreground) / 0.45);
	}

	/* Keyboard users still get an unmistakable marker. */
	.clause-input:focus-visible {
		box-shadow: inset 2px 0 0 hsl(var(--foreground) / 0.55);
	}

	.clause-popup {
		position: fixed;
		z-index: 60;
		max-height: 17rem;
		overflow-y: auto;
		padding: 0.25rem;
		border: 1px solid hsl(var(--border));
		border-radius: calc(var(--radius) - 2px);
		background: hsl(var(--popover));
		color: hsl(var(--popover-foreground));
		box-shadow:
			0 1px 2px hsl(var(--foreground) / 0.06),
			0 8px 24px hsl(var(--foreground) / 0.12);
	}

	.clause-option {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.25rem 0.4rem;
		border-radius: calc(var(--radius) - 4px);
		border-left: 2px solid transparent;
		cursor: pointer;
		font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
		font-size: 0.78125rem;
		line-height: 1.5;
	}

	.clause-option.is-active {
		background: hsl(var(--accent));
		border-left-color: hsl(var(--foreground) / 0.6);
	}

	.clause-glyph {
		flex: 0 0 auto;
		width: 0.75rem;
		text-align: center;
		font-size: 0.5625rem;
		color: hsl(var(--muted-foreground) / 0.7);
	}

	.clause-name {
		flex: 1 1 auto;
		min-width: 0;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		color: hsl(var(--muted-foreground));
	}

	.clause-hit {
		color: hsl(var(--foreground));
		font-weight: 600;
	}

	.clause-badge {
		flex: 0 0 auto;
		padding: 0 0.3rem;
		border: 1px solid hsl(var(--border));
		border-radius: 2px;
		font-size: 0.625rem;
		color: hsl(var(--muted-foreground));
	}

	.clause-detail {
		flex: 0 0 auto;
		max-width: 45%;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
		font-size: 0.6875rem;
		color: hsl(var(--muted-foreground) / 0.75);
	}
</style>
