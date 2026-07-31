# Sidebar example polish implementation plan

1. Add browser regressions for separator geometry, non-overlapping Basic menu
   decorations, and the icon-collapsed brand mark.
2. Confirm the new assertions fail against the current examples.
3. Increase the Sidebar separator rule's orientation specificity so `w-auto`
   wins over the base horizontal `w-full`.
4. Split the Basic example's action and badge demonstrations across distinct
   menu items.
5. Add one shared example brand component and reuse it from Variants and
   Persisted.
6. Regenerate GSX and highlighted source output.
7. Run the focused Sidebar browser test, relevant Go/style checks, and the
   complete repository check.
8. Reinspect the running documentation preview at the reported viewport.
