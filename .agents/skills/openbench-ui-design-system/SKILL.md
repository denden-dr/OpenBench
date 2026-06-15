---
name: openbench-ui-design-system
description: Use when styling components, building dashboard layouts, managing viewport constraints, or applying the Neubrutalism design language with Tailwind CSS v4.
version: 1.0.0
---

# OpenBench UI Design System

## Overview
OpenBench uses a "Neubrutalism" aesthetic powered by Tailwind CSS v4. This guide covers the visual language, typography, and responsive layout constraints necessary to build consistent dashboards.

## Core Principles (Neubrutalism)
- High contrast, thick borders (e.g., `border-4 border-neubrutalism-charcoal`).
- Solid, unblurred drop shadows (`shadow-neubrutalism-md`).
- Vibrant background colors clashing with stark whites and charcoals.
- Tactile, "pushable" button interactions.

## Color Palette (Tailwind V4 Theme Overrides)
Use specific palette classes:
- `neubrutalism-charcoal`: `#1C1C1C`
- `neubrutalism-green`: `#4ADE80`
- `neubrutalism-yellow`: `#FBBF24`
- `neubrutalism-pink`: `#F472B6`
- `neubrutalism-blue`: `#60A5FA`

## Component Styling Guidelines
1. **Cards and Containers**: Thick border and solid drop shadow.
2. **Buttons**: When pressed (`active:`), they should move down (`active:translate-y-1 active:shadow-none`).
3. **Inputs & Textareas**: Distinct focus states with thick rings (`focus:ring-4 focus:ring-neubrutalism-charcoal`), avoiding default Tailwind focus rings.
4. **Active States & Hover Contrasts**: When overriding active links with brand colors, bind corresponding dark/light hover shades. Do not apply generic `hover:bg-zinc-100` to all states, as it causes contrast failures on active elements.
5. **Loading Skeletons (Asynchronous States)**: Instead of generic thin spinners, use stark Neubrutalist pulse blocks (`animate-pulse bg-zinc-200 border-4 border-neubrutalism-charcoal`) for content loading asynchronously. This maintains the structural aesthetic even during network delays.

## Responsive Dashboard Layouts
1. **The Wrapper**: For true dashboards, lock the wrapper strictly to viewport height (`h-screen w-full flex overflow-hidden`).
2. **Sidebar Positioning**: Use `sticky top-0 h-screen` or `fixed` for sidebars. Avoid `static` sidebars inside `min-h-screen` wrappers, as they will stretch infinitely and push footers out of view.
3. **Independent Scrolling**: Once the app is locked to `h-screen`, handle scrolling locally (e.g., `<main class="flex-grow h-full overflow-y-auto">`).
4. **Mobile Drawers**: Use an overlay and absolute/fixed positioning to transition sidebars on smaller screens.

## Common Mistakes to Avoid
- Mixing `min-h-screen` with a `flex-col` sidebar causing the sidebar to stretch off-screen.
- Forgetting `overflow-y-auto` on the `<main>` container.
- Using `static` positioning for sidebars inside flexible wrappers.
- Creating a `tailwind.config.js` (Tailwind v4 uses CSS `@theme` properties).
