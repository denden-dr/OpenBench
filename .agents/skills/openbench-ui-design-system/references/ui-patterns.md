# UI Patterns

## Files To Inspect First

- Theme tokens: `apps/frontend/src/routes/layout.css`
- Shared components: `apps/frontend/src/lib/components/Button.svelte`, `Card.svelte`, `Input.svelte`
- Admin shell: `apps/frontend/src/routes/admin/+layout.svelte`
- Public tracker: `apps/frontend/src/routes/tracker`

## Theme Tokens

OpenBench uses Tailwind CSS v4. Theme values live in CSS:

- `--color-neubrutalism-yellow`
- `--color-neubrutalism-green`
- `--color-neubrutalism-pink`
- `--color-neubrutalism-charcoal`
- `--color-neubrutalism-bg`
- `--shadow-neubrutalism-sm|md|lg`

Do not add `tailwind.config.js` for new tokens unless the project has intentionally moved away from CSS `@theme`.

## Visual Language

- Use `border-4 border-neubrutalism-charcoal` for primary surfaces.
- Use hard shadows such as `shadow-neubrutalism-md`, not blurred shadows.
- Use `rounded-none` unless a route already establishes another radius.
- Use tactile button states with hard-shadow movement.
- Use lucide icons inside buttons and compact action controls when an icon exists.

## Cards And Surfaces

- Cards are for individual repeated items, modals, and framed controls.
- Page sections should be direct layouts or full-width bands, not cards inside cards.
- Reuse shared `Card` for normal content panels; use raw markup only when a table, modal, or shell needs precise structure.

## Dashboards

- Prefer dense, scannable dashboards over marketing-style sections.
- For app shells, constrain the root to viewport height when the sidebar and main region need independent scrolling.
- Use sticky or fixed sidebars; avoid static sidebars inside long `min-h-screen` flex layouts.
- Put `overflow-y-auto` on the scrolling region, usually `main`.

## Active Navigation

- Active colored links need matching hover colors:
  - Yellow active: amber hover
  - Green active: emerald hover
  - Pink active: pink hover with readable text
- Do not apply generic `hover:bg-zinc-100` to active colored links.

## Loading States

- Use skeleton blocks with borders and `animate-pulse`.
- Avoid small generic spinners for content regions.
- Preserve layout dimensions while loading to prevent visible shifts.

## Responsive Checks

- Verify text does not overflow buttons, status chips, cards, or table cells.
- Keep admin tables horizontally scrollable when columns are dense.
- Mobile drawers need overlay/fixed behavior and an obvious close path.
