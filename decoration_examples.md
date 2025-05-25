# Monaco Editor Line Decorations Guide

## Overview
This guide shows how to add decorations to all non-empty lines in Monaco Editor, excluding empty lines.

## Implementation

### 1. Basic Line Decorations
The `createLineDecorations()` function adds decorations to all lines containing text:

```typescript
function createLineDecorations() {
    if (!editor || !isInitialized || !monacoInstance) return;

    const model = editor.getModel();
    if (!model) return;

    const totalLines = model.getLineCount();
    const decorations = [];

    for (let lineNumber = 1; lineNumber <= totalLines; lineNumber++) {
        const lineContent = model.getLineContent(lineNumber);
        
        // Only add decoration if line has content (not empty or just whitespace)
        if (lineContent.trim() !== '') {
            decorations.push({
                range: new monacoInstance.Range(lineNumber, 1, lineNumber, 1),
                options: {
                    isWholeLine: true,
                    className: 'line-highlight',
                    glyphMarginClassName: 'glyph-margin-highlight',
                    glyphMarginHoverMessage: { value: 'Click to execute this line' },
                    hoverMessage: { value: 'Non-empty line with SQL content' }
                }
            });
        }
    }

    // Clear existing decorations and create new ones
    if (decorationsCollection) {
        decorationsCollection.clear();
    }
    decorationsCollection = editor.createDecorationsCollection(decorations);
}
```

### 2. Advanced Decorations with Conditional Styling
The `createAlternativeLineDecorations()` function provides different styles based on content:

```typescript
function createAlternativeLineDecorations() {
    // Different decoration styles based on SQL keywords
    if (lineContent.trim().toUpperCase().startsWith('SELECT')) {
        className = 'bg-green-100';  // Green for SELECT queries
        glyphClassName = 'bg-green-500';
    } else if (lineContent.trim().toUpperCase().startsWith('UPDATE') || 
               lineContent.trim().toUpperCase().startsWith('INSERT') ||
               lineContent.trim().toUpperCase().startsWith('DELETE')) {
        className = 'bg-red-100';    // Red for write operations
        glyphClassName = 'bg-red-500';
    }
}
```

## Decoration Options

### Range Options
- `new monacoInstance.Range(lineNumber, 1, lineNumber, 1)` - Single line
- `new monacoInstance.Range(startLine, 1, endLine, 1)` - Multiple lines

### Decoration Options
- `isWholeLine: true` - Decorates the entire line
- `className` - CSS class for line background
- `glyphMarginClassName` - CSS class for glyph margin (left side)
- `hoverMessage` - Tooltip on line hover
- `glyphMarginHoverMessage` - Tooltip on glyph margin hover

## CSS Styles

### Blue Highlight (Default)
```css
:global(.line-highlight) {
    background-color: rgba(59, 130, 246, 0.08) !important;
    border-left: 3px solid rgba(59, 130, 246, 0.4);
    transition: background-color 0.2s ease;
}
```

### Green Highlight (SELECT queries)
```css
:global(.bg-green-100) {
    background-color: rgba(34, 197, 94, 0.1) !important;
    border-left: 2px solid rgba(34, 197, 94, 0.3);
}
```

### Red Highlight (Write operations)
```css
:global(.bg-red-100) {
    background-color: rgba(239, 68, 68, 0.1) !important;
    border-left: 2px solid rgba(239, 68, 68, 0.3);
}
```

## Usage

### Automatic Updates
The decorations update automatically when:
- Editor content changes (`onDidChangeModelContent`)
- Editor is initialized

### Manual Updates
Call `createLineDecorations()` whenever you want to refresh decorations.

### Switching Decoration Styles
Replace `createLineDecorations()` with `createAlternativeLineDecorations()` in the event handlers to use conditional styling.

## Key Features

✅ **Excludes empty lines** - Only decorates lines with content
✅ **Real-time updates** - Decorations update as you type
✅ **Customizable styles** - Easy to modify colors and effects
✅ **Hover tooltips** - Interactive feedback
✅ **Glyph margin support** - Left margin decorations
✅ **Performance optimized** - Uses Monaco's decoration collections
