# Sidebar example polish

## Problem

The Sidebar documentation currently exposes three visual defects:

1. `SidebarSeparator` combines the base horizontal separator's `w-full` with
   sidebar `mx-2`, so its right edge extends beyond the sidebar.
2. The Basic example places two actions and one badge in the same menu item.
   All three controls intentionally occupy the same absolute position, and the
   badge is visually displaced over the active submenu row.
3. The collapsible examples render `Acme Inc` as an ordinary wrapping `div`.
   In icon mode the available width is 3rem, so the name wraps and leaks.

## Design

- Make the sidebar separator's horizontal selector as specific as the base
  separator orientation selector. The later sidebar rule can then replace
  `w-full` with `w-auto` while retaining the intended two-sided inset.
- Keep each Basic menu item to one trailing decoration. The Inbox keeps its
  hover action and submenu, Calendar demonstrates a badge, and Search
  demonstrates an always-visible action. This preserves the API examples
  without overlapping controls.
- Render the example brand through a large `SidebarMenuButton`: a fixed square
  `A` mark followed by `Acme Inc`. Existing icon-collapse behavior shrinks the
  button to the mark and clips the name, matching the rest of the menu.

The separator correction belongs to the component CSS. The menu arrangement
and brand treatment remain example-owned because the library cannot infer
consumer branding or decide which trailing menu decoration a consumer wants.

## Verification

Extend the Sidebar documentation browser test to assert that:

- the separator is inset on both sides of the desktop sidebar;
- Basic menu decorations stay within their own button row and do not overlap;
- the icon-collapsed brand exposes the `A` mark without overflowing its button.

Regenerate GSX and highlighted example output, then run the focused browser
spec and the repository check.
