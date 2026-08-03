# Home Page Component Showcase

**Date:** 2026-08-03
**Status:** Approved

## Problem

The landing page's component section (`site/pages/home.gsx:61-106`) shows only
Button, Badge, and Dialog as flat variant rows. Button and Badge render almost
identically at a glance, and the section undersells a library of 50+
components. It reads as swatches, not as a UI library in use.

## Goal

Replace the flat section with a bento grid of composed demo cards that show
gsxui components working together in realistic app UI, shadcn.com-style.

## Design

### Section structure

- Replace the current `#components` section in `home.gsx` entirely.
- Single heading (e.g. "Built with gsxui") plus a "Browse all components →"
  link to the components index.
- Grid container: `grid gap-4 md:grid-cols-2`. Four cards, 2×2 on desktop,
  stacked on mobile. Natural height differences give the bento feel; no
  row-span tricks.

### The four cards

Each card is a `ui.Card` composing real library components with static inline
demo data:

1. **SignInCard** — CardHeader (title + description), Field/Label + Input for
   email and password, Checkbox ("Remember me"), CardFooter with a full-width
   default Button and a ghost Button.
2. **SettingsCard** — 2–3 rows of Label + Switch, a NativeSelect (e.g. theme
   choice), a Slider (e.g. density), with Separators between rows.
3. **StatsCard** — stat lines with Progress bars, status Badges, and an
   Avatar + name row.
4. **OverlaysCard** — Tabs wrapping a row of interaction triggers: the
   existing Dialog demo (moved here from home.gsx), a DropdownMenu (following
   the trigger-slot pattern in `site/examples/dropdown/basic.gsx`), a Tooltip
   on a button, and a Toast trigger. Demonstrates server-rendered
   interactivity with no client framework.

### Implementation location

New package `site/examples/showcase/` with one component per card
(`SignInCard`, `SettingsCard`, `StatsCard`, `OverlaysCard`), imported by
`home.gsx`. This matches the existing `site/examples/` convention, keeps
`home.gsx` small, and makes each card independently testable. No new registry
entries; no changes to `ui/`.

### Out of scope

- No changes to the hero, install snippet, nav, or footer.
- No new ui components or registry additions.
- No per-card links to component docs (the single "Browse all components"
  link covers discovery).

## Testing

- Render test asserting the home page contains the four showcase cards
  (extend existing pages/examples test patterns).
- Playwright check on the home page: dialog opens/closes and dropdown opens
  from the OverlaysCard (note the Playwright config trap documented in the
  verification-gates memory).
- Verify dark mode and RTL rendering don't break (RTL follows the project's
  physical-side ruling).
