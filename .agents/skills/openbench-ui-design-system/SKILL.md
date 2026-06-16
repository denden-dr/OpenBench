---
name: openbench-ui-design-system
description: Apply OpenBench visual design and layout rules. Use when styling Svelte components, building admin dashboards, public tracker pages, responsive sidebars, loading states, or Tailwind CSS v4 Neubrutalism UI in apps/frontend.
---

# OpenBench UI Design System

## Operating Rule

Reuse the established OpenBench component language. Inspect `apps/frontend/src/lib/components` and the closest route page before styling new UI.

## Workflow

1. Check tokens in `apps/frontend/src/routes/layout.css`; Tailwind v4 theme values live in CSS `@theme`, not `tailwind.config.js`.
2. Use high-contrast borders, hard shadows, tactile button states, and explicit focus rings consistently.
3. Keep dashboard layouts work-focused: dense, scannable, predictable navigation and local scrolling.
4. Use stark skeleton blocks for async loading instead of thin generic spinners.
5. Verify mobile and desktop layout constraints for text overflow, sidebar behavior, and independent scroll regions.

## Load References

- Read `references/ui-patterns.md` before changing shared components, dashboard shells, public tracker pages, navigation, tables, modals, or loading states.

## Hard Checks

- Do not create `tailwind.config.js` for theme tokens.
- Do not use generic hover colors on active colored navigation states.
- Do not nest cards inside cards.
- Do not let page-level dashboard wrappers create double-page scroll unless the route intentionally needs it.
