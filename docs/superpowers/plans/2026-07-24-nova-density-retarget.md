# Nova Density Retarget Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Retarget gsxui's density metrics (heights, paddings, gaps, radii, font sizes, icon sizes) from shadcn's new-york-v4 registry to the **nova** style — what ui.shadcn.com renders today — per the user's "Full density retarget" decision (2026-07-24).

**Architecture:** Metric-token substitution over the existing components. Markup structure, slots, behavior JS, and the theme color system stay as-is; only class-string metric tokens change. The single source of truth for every token is the committed map **`docs/superpowers/plans/2026-07-24-nova-density-map.md`** — each task below names its map sections, and the implementer applies each section's `delta:` list exactly, consulting `nova:` for the full target strings. Already shipped out-of-band (do NOT redo): Geist font, dialog/sheet overlay (`bg-black/10` + `backdrop-blur-xs`), nova radio paint, popover-family discrete transitions.

**Tech Stack:** gsx components (`ui/*.gsx`), `go tool gsx generate`, Go pin tests, Tailwind v4 + tw-animate-css, tailwind-merge-go.

## Global Constraints

- Metric tokens ONLY: h-/w-/size-/min-w-/max-h-/p*/gap-/space-/text-size/rounded-/`[&_svg…]:size-*` (+ the two explicitly-listed non-metric base additions for button). Never change colors, focus-ring colors, dark: color variants, semantics, or behavior JS in this plan.
- Shadow deltas: apply exactly the removals the map's per-component `notes:` and the `## Shadow-presence footnote (roll-up)` section list (nova drops `shadow-xs` in several places); add no new shadows.
- ADAPT (plan-wide decision): nova's `has-data-[icon=inline-start/end]:p*` directional icon paddings require `data-icon` stamps we don't have. Keep gsxui's existing `has-[>svg]:px-*` selector MECHANISM and substitute nova's numeric values (e.g. button default `has-[>svg]:px-3` → `has-[>svg]:px-2`, using the inline-start value). Record this in each affected component's doc comment.
- Out of scope (ledger as roadmap items, do not implement): the map's "Markup prerequisites" item 3 (new parts: AlertAction, AlertDialogMedia, PopoverHeader/Title/Description, ProgressLabel/Value, AvatarBadge/AvatarGroupCount, Item size-xs, ItemGroup gaps) and item 4 (new size axes: alert-dialog/avatar/card/native-select/switch sizes).
- After every component batch: `go tool gsx generate`, `go test ./...`, re-pin failing class pins with the splice procedure (Task 0), `make check` green, one commit per task.
- The dev loop may already have regenerated `.x.go` (gsx generate then prints "up to date") — that is normal; commit whatever `.x.go` state `make check` verifies.

---

### Task 0: Re-pin helper script (shared by all tasks)

**Files:**
- Create: `scripts/repin.py`

**Interfaces:**
- Produces: `python3 scripts/repin.py <TestA|TestB|…> <test_file1> [test_file2…]` — runs the named Go tests, extracts every `got:`/`want:` class-attribute pair from the failure output, and splices the actual class strings into the given `_test.go` files. Used verbatim by Tasks 1–7.

- [ ] **Step 1: Write the script**

```python
#!/usr/bin/env python3
"""Splice actual rendered class strings into pinned `want` literals.

Usage: python3 scripts/repin.py 'TestFoo|TestBar' ui/foo_test.go ui/bar_test.go
Run from the repo root AFTER `go tool gsx generate`. Idempotent: with no
failing pins it replaces nothing.
"""
import pathlib, re, subprocess, sys

run_expr, files = sys.argv[1], sys.argv[2:]
r = subprocess.run(["go", "test", "./ui", "-run", run_expr],
                   capture_output=True, text=True)
out = r.stdout + r.stderr
gots = re.findall(r"got: (<[^\n]*)", out)
wants = re.findall(r"want: (<[^\n]*)", out)
assert len(gots) == len(wants), (len(gots), len(wants))
pairs = []
for g, w in zip(gots, wants):
    gm, wm = re.search(r'class="([^"]*)"', g), re.search(r'class="([^"]*)"', w)
    if gm and wm and gm.group(1) != wm.group(1):
        pairs.append((wm.group(1), gm.group(1)))
for f in files:
    p = pathlib.Path(f)
    s = p.read_text()
    n = 0
    for wc, gc in pairs:
        if wc in s:
            s = s.replace(wc, gc)
            n += 1
    p.write_text(s)
    print(f, "replaced", n)
```

- [ ] **Step 2: Verify it is a no-op on a green tree**

Run: `go test ./ui && python3 scripts/repin.py 'TestNothingMatches' ui/button_test.go`
Expected: tests pass; script prints `ui/button_test.go replaced 0`.

- [ ] **Step 3: Commit**

```bash
git add scripts/repin.py
git commit -m "chore: shared pin-splice helper for the nova density retarget"
```

### Task 1: Button core

**Files:**
- Modify: `ui/button.gsx` (base string + `sizeClass` switch)
- Test: `ui/button_test.go`, plus consumers' pins: `ui/pagination_test.go`, `ui/button-group_test.go`, `ui/input-group_test.go`, `ui/dialog_test.go`, `ui/alert-dialog_test.go`, `ui/empty_test.go`, `ui/item_test.go`

**Map sections:** `## button` (apply the full `delta:` list — including the three base additions `rounded-lg`, `border border-transparent bg-clip-padding`, `active:not-aria-[haspopup]:translate-y-px`, base `gap-2` removal, and all seven size rewrites; use the `has-[>svg]` ADAPT from Global Constraints for the directional paddings; drop outline's `shadow-xs` per `notes:`).

- [ ] **Step 1: Apply the button deltas to `ui/button.gsx`** — base string and every case of the `sizeClass` switch, exactly per the map. Document the `has-[>svg]` ADAPT and the transparent-border box note in the component comment.
- [ ] **Step 2: Regenerate and test** — Run: `go tool gsx generate && go test ./...` Expected: FAILs limited to class pins.
- [ ] **Step 3: Re-pin** — Run: `python3 scripts/repin.py 'TestButton|TestPagination|TestButtonGroup|TestInputGroup|TestDialog|TestAlertDialog|TestEmpty|TestItem' ui/button_test.go ui/pagination_test.go ui/button-group_test.go ui/input-group_test.go ui/dialog_test.go ui/alert-dialog_test.go ui/empty_test.go ui/item_test.go` then `go test ./...` Expected: PASS.
- [ ] **Step 4: Check and commit**

```bash
make check
git add -A && git commit -m "feat(nova): button density — h-8/px-2.5 default, nova radii, press effect"
```

### Task 2: Form controls

**Files:**
- Modify: `ui/input.gsx`, `ui/textarea.gsx`, `ui/select.gsx`, `ui/checkbox.gsx`, `ui/toggle.gsx`, `ui/field.gsx`
- Test: matching `ui/*_test.go` for each

**Map sections:** `## input`, `## textarea`, `## select (native)`, `## checkbox`, `## toggle`, `## field`.

- [ ] **Step 1: Apply each section's `delta:` list** to its component (`select.gsx` keeps its wrapper/chevron mechanism — only the metric tokens on the `<select>` and chevron change).
- [ ] **Step 2: Regenerate, test, re-pin** — `go tool gsx generate && go test ./...`; then `python3 scripts/repin.py 'TestInput|TestTextarea|TestSelect|TestCheckbox|TestToggle|TestField' ui/input_test.go ui/textarea_test.go ui/select_test.go ui/checkbox_test.go ui/toggle_test.go ui/field_test.go`; `go test ./...` Expected: PASS.
- [ ] **Step 3: Check and commit** — `make check`; `git add -A && git commit -m "feat(nova): form-control density — input/textarea/select/checkbox/toggle/field"`

### Task 3: Menus and floating content

**Files:**
- Modify: `ui/dropdown.gsx`, `ui/context-menu.gsx`, `ui/popover.gsx`, `ui/hover-card.gsx`, `ui/tooltip.gsx`
- Test: matching `ui/*_test.go`

**Map sections:** `## dropdown`, `## context-menu`, `## popover`, `## hover-card`, `## tooltip`. Do NOT touch the discrete-transition blocks added 2026-07-24 — only the size/padding/radius tokens around them.

- [ ] **Step 1: Apply each section's `delta:` list** (content min-widths/padding, item `px/py/gap/rounded`, label/shortcut/separator metrics).
- [ ] **Step 2: Regenerate, test, re-pin** — `go tool gsx generate && go test ./...`; `python3 scripts/repin.py 'TestDropdown|TestContextMenu|TestPopover|TestHoverCard|TestTooltip' ui/dropdown_test.go ui/context-menu_test.go ui/popover_test.go ui/hover-card_test.go ui/tooltip_test.go`; `go test ./...` Expected: PASS.
- [ ] **Step 3: Check and commit** — `make check`; `git add -A && git commit -m "feat(nova): menu/floating density — dropdown, context-menu, popover, hover-card, tooltip"`

### Task 4: Dialog family content

**Files:**
- Modify: `ui/dialog.gsx`, `ui/alert-dialog.gsx`, `ui/sheet.gsx`
- Test: `ui/dialog_test.go`, `ui/alert-dialog_test.go`, `ui/sheet_test.go`

**Map sections:** `## dialog`, `## alert-dialog`, `## sheet` — content paddings/gaps/max-widths/title sizes only. The `backdrop:*` token run is ALREADY nova (bg-black/10 + blur-xs, shipped b85a49c): leave it byte-identical.

- [ ] **Step 1: Apply the three `delta:` lists** (panel `p-*`/`gap-*`/`rounded-*`/`sm:max-w-*`, header/footer gaps, title text size, close-button offsets).
- [ ] **Step 2: Regenerate, test, re-pin** — `go tool gsx generate && go test ./...`; `python3 scripts/repin.py 'TestDialog|TestAlertDialog|TestSheet' ui/dialog_test.go ui/alert-dialog_test.go ui/sheet_test.go`; `go test ./...` Expected: PASS.
- [ ] **Step 3: Check and commit** — `make check`; `git add -A && git commit -m "feat(nova): dialog-family density — dialog, alert-dialog, sheet"`

### Task 5: Containers and grouping

**Files:**
- Modify: `ui/card.gsx`, `ui/alert.gsx`, `ui/accordion.gsx`, `ui/tabs.gsx`, `ui/empty.gsx`, `ui/item.gsx`, `ui/input-group.gsx`, `ui/button-group.gsx`, `ui/badge.gsx`
- Test: matching `ui/*_test.go`

**Map sections:** `## card`, `## alert`, `## accordion`, `## tabs`, `## empty`, `## item`, `## input-group`, `## button-group`, `## badge`.

- [ ] **Step 1: Apply each `delta:` list.** For button-group, ALSO port the structural corner mechanism per its `notes:`: replace the zero-inner-corners selectors with nova's restore-outer-corner pattern (`[&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-lg!` horizontal / `rounded-b-lg!` vertical, per the `nova:` block) — this is the one structural change the map flags as not-a-token-swap.
- [ ] **Step 2: Regenerate, test, re-pin** — `go tool gsx generate && go test ./...`; `python3 scripts/repin.py 'TestCard|TestAlert|TestAccordion|TestTabs|TestEmpty|TestItem|TestInputGroup|TestButtonGroup|TestBadge' ui/card_test.go ui/alert_test.go ui/accordion_test.go ui/tabs_test.go ui/empty_test.go ui/item_test.go ui/input-group_test.go ui/button-group_test.go ui/badge_test.go`; `go test ./...` Expected: PASS.
- [ ] **Step 3: Check and commit** — `make check`; `git add -A && git commit -m "feat(nova): container density — card, alert, accordion, tabs, empty, item, groups, badge"`

### Task 6: Navigation and indicators

**Files:**
- Modify: `ui/breadcrumb.gsx`, `ui/pagination.gsx`, `ui/progress.gsx`
- Test: `ui/breadcrumb_test.go`, `ui/pagination_test.go`, `ui/progress_test.go`

**Map sections:** `## breadcrumb`, `## pagination`, `## progress`. (Pagination's box size largely lands via Task 1's button `icon: size-9 → size-8`; this task applies its remaining own-token deltas, e.g. the ellipsis cell.)

- [ ] **Step 1: Apply each `delta:` list.**
- [ ] **Step 2: Regenerate, test, re-pin** — `go tool gsx generate && go test ./...`; `python3 scripts/repin.py 'TestBreadcrumb|TestPagination|TestProgress' ui/breadcrumb_test.go ui/pagination_test.go ui/progress_test.go`; `go test ./...` Expected: PASS.
- [ ] **Step 3: Check and commit** — `make check`; `git add -A && git commit -m "feat(nova): nav/indicator density — breadcrumb, pagination, progress"`

### Task 7: Examples sweep, ledger, deploy

**Files:**
- Modify: `site/examples/**/*.gsx` (only where an example hard-codes a now-wrong metric, e.g. `class="size-8"` overrides on buttons that are now size-8 by default, spacing built around h-9 controls)
- Modify: `docs/jsx-parity.md` (new top-level `## nova density` entry), `docs/component-roadmap.md` (add the out-of-scope nova parts/size-axes from the map's "Markup prerequisites" items 3–4 as roadmap items)
- Test: `go test ./...` (site example pins, if any, re-pinned via `scripts/repin.py` with the failing test names)

- [ ] **Step 1: Sweep examples** — `grep -rn "size-9\|h-9\|size-8\|px-4" site/examples/` and fix only genuinely stale overrides (an override that duplicates the new default gets removed; a deliberate size deviation stays).
- [ ] **Step 2: Ledger** — `docs/jsx-parity.md` gains `## nova density`: one entry stating the retarget decision, the map/plan file paths, the `has-[>svg]` ADAPT, the button-group structural port, and that colors/behavior stayed new-york-v4. Roadmap gains the deferred nova parts/size axes.
- [ ] **Step 3: Full check and commit** — `make check`; `git add -A && git commit -m "feat(nova): examples sweep + ledger for the density retarget"`
- [ ] **Step 4: Push and deploy** — `git push origin main`; confirm `gh run list --repo gsxhq/gsxui --limit 1` reaches success. User performs the visual pass against ui.shadcn.com.

## Self-Review

- Spec coverage: all 27 delta components from the map are assigned (T1 button; T2 input/textarea/select/checkbox/toggle/field; T3 dropdown/context-menu/popover/hover-card/tooltip; T4 dialog/alert-dialog/sheet; T5 card/alert/accordion/tabs/empty/item/input-group/button-group/badge; T6 breadcrumb/pagination/progress) = 27 ✓; no-delta and no-counterpart components untouched ✓; radio/overlay/font/transitions already shipped and excluded ✓.
- Placeholders: none — every step names its exact command, file list, and map section; token values live in the committed map the briefs reference.
- Type consistency: no new Go surface; `scripts/repin.py` interface defined once in Task 0 and used identically in Tasks 1–7.
