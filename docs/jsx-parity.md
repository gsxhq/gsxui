# JSX parity ledger

Divergences between gsxui components and their shadcn/ui reference, both
directions. Full audit: gsxhq docs repo, specs/2026-07-22-gsx-over-jsx-audit.md.

## packaging
- NOTE (single package by design): `ui/` is one flat `package ui` — every
  component is `ui.Button`, `ui.Card`, `ui.DialogContent`, and so on, not
  `button.Button`/`card.Card`/`dialog.DialogContent` behind a per-component
  import. shadcn's own component names already carry their namespace
  (`Button`, `CardHeader`, `DropdownMenuContent`), so a per-package prefix
  just stutters (`button.Button`, `card.CardHeader`); a consumer package
  literally named after its lone export (`package button` exporting
  `button.Button`) also forced every real caller into an import alias the
  moment it needed more than one component, since Go packages can't share a
  name with their own directory-mate. Two components — `select` and
  `switch` — additionally couldn't be per-component packages at all: both
  names are reserved Go keywords, illegal as a package name or an import
  alias. One flat package sidesteps the stutter, the forced aliasing, and
  the keyword problem in a single move.

## rendering
- FINDING (gsx upstream, RESOLVED 2026-07-24): gsx used to emit authored self-closing syntax verbatim — `<div/>` reached HTML as `<div/>`, which browsers parse as an OPEN tag, silently nesting all following siblings (bit dropdown's separator and skeleton in production; filed as gsxhq/gsx#144). gsx now serializes canonically (#156): self-closed non-voids are EXPANDED to explicit open/close pairs, void elements are emitted bare (`<input …>`, no `/`), and svg-subtree elements authored self-closed expand too (`<path/>` → `<path></path>`, equally valid in foreign content). The explicit-close discipline in .gsx sources is no longer load-bearing, and render pins now assert the canonical forms.

## animations
- FINDING (2026-07-24, fixed): every `animate-in`/`animate-out`/`fade-*`/`zoom-*`/`slide-in-*` token carried by dialog, dropdown, and tooltip was silently dead from v1 — those utilities are NOT Tailwind v4 core; they come from the `tw-animate-css` package, which shadcn's v4 globals.css imports and ours didn't. Ported classes pinned byte-for-byte can still be inert if the keyframes layer behind them never ships: class-string tests, `make check`, and even open/close behavior checks all pass without it (the components work, they just snap). Fixed by `@import "tw-animate-css"` in assets/gsxui.css + web/site.css and an npm dep; getting-started now tells consumers to install it. Verified live: computed `animation-name` was `none` before, keyframed after.
- FINDING (2026-07-24, Tier 2 browser pass, fixed): dialog.js's animated close awaited every `getAnimations({subtree:true})` `finished` promise with no cap — an unbounded child animation (ui.Spinner's `animate-spin` is `iterations:Infinity`, so its promise never resolves) or a backgrounded tab's frozen animation clock left the dialog open forever in `data-state="closed"`. `requestClose` now waits only on the dialog element's OWN animations, raced against a 600ms cap (comfortably above every shipped exit: dialog/alert-dialog 200ms, sheet 300ms) — the dialog always closes. Applies to all three `<dialog>` consumers (dialog, alert-dialog, sheet).
- NOTE (2026-07-24, resolved — was the "enter animations only" backlog item): `hidePopover()` display-nones the element immediately, so `data-[state=closed]:animate-out` never got a painted frame — popover-based components (dropdown, context-menu, popover, tooltip, hover-card) had working enters and silently dead exits. Resolved by ADAPT: those five now animate with discrete TRANSITIONS instead of tw-animate keyframes — base carries the closed geometry (`opacity-0 scale-95`), `open:` (Tailwind's variant matches `:popover-open`) the open one, `transition-[opacity,scale,translate,display,overlay] transition-discrete` keeps the element rendered and in the top layer through the exit (`allow-discrete` defers the discrete display/overlay flips to transition end), and `starting:open:*` (`@starting-style`) supplies the enter-from state, with the ported `slide-in-from-*-2` side offsets mapped to `data-[side=…]:starting:open:*translate-*-2`. Same fade/zoom/slide geometry as the shadcn tokens; graceful degradation — browsers without `@starting-style`/`allow-discrete` snap open/closed but stay fully functional. Also relevant: `data-state="open"` is stamped synchronously BEFORE `showPopover()` in the toggle-driven trio (dropdown/popover/context-menu) — the popover `toggle` event is a queued task, and a paint in the gap flashed one closed-state frame (routine in inactive windows, occasional in active ones). The `<dialog>` trio (dialog/alert-dialog/sheet) keeps tw-animate keyframes: its close is JS-intercepted, so exits get their frames.

## nova density
- DECISION (2026-07-24, "Full density retarget"): gsxui's visual-parity target flips from shadcn's **new-york-v4** registry style to **nova** — the style ui.shadcn.com actually renders as its default today. new-york-v4 remains the STRUCTURAL/markup reference (slots, parts, behavior JS, the color/theme system, dark-mode variants) — this retarget touches metric tokens ONLY: heights, widths, paddings, gaps, radii, font sizes, icon sizes, and the shadow-presence deltas nova's own CSS drops. Colors, focus-ring colors, dark: variants, and behavior are explicitly out of scope and untouched. Full rationale, per-component before/after token tables, and the task-by-task execution log live in the committed map and plan: `docs/superpowers/plans/2026-07-24-nova-density-map.md` (source of truth for every token; nova CSS reference `style-nova.css`, structural reference `radix/ui/*.tsx`) and `docs/superpowers/plans/2026-07-24-nova-density-retarget.md` (the 8-task execution plan, `.superpowers/sdd/progress.md` "NOVA DENSITY RETARGET" section carries the per-task adjudication notes). 27 of gsxui's components carried a metric delta and were retargeted (button, input, textarea, native-select, checkbox, toggle, field, dropdown, context-menu, popover, hover-card, tooltip, dialog, alert-dialog, sheet, card, alert, accordion, tabs, empty, item, input-group, button-group, badge, breadcrumb, pagination, progress); 8 had no delta and 4 have no nova counterpart (aspect-ratio, collapsible, spinner, icon) and were left untouched.
- ADAPT (icon-adjacent padding, button + toggle): nova replaces the old `has-[>svg]:px-*` selector with directional `has-data-[icon=inline-start]:pl-*` / `has-data-[icon=inline-end]:pr-*`, which requires stamping `data-icon="inline-start|inline-end"` on icon children — a markup change out of scope for a metric-only retarget. Both components keep the existing `has-[>svg]:px-*` MECHANISM and substitute nova's numeric value in its place (e.g. button default `has-[>svg]:px-3` → `has-[>svg]:px-2`, nova's inline-start figure) — same selector shape, nova's number. See each component's own doc comment in `ui/button.gsx`/`ui/toggle.gsx` for the inline note.
- PORT (button-group corner mechanism): nova inverts how button-group zeroes/restores corner radii at the ends of a row. gsxui's pre-nova mechanism (zero every inner corner via `[&>*:not(:first-child)]:rounded-l-none` / `[&>*:not(:last-child)]:rounded-r-none`, vertical analog `rounded-t/b-none`) is KEPT as the base — nova's own CSS starts from the same zero-inner-corners posture. What's ADDED on top is nova's supplementary restore-outer-corner rule, which force-restores the true end child's outer radius (`[&>[data-slot]:not(:has(~[data-slot]))]:rounded-r-lg!` horizontal / `rounded-b-lg!` vertical) rather than relying solely on the never-touched first/last-child selectors — this is a structural port (new selector shape), not a token swap, because gsxui's `button-group` also composes `native-select`'s trigger which needs the same restore. `## button-group` carries the per-selector detail.
- NOT ADOPTED (notes-only, out of scope — color/structure, not metric): nova's popover-family surfaces (dropdown/context-menu/popover/hover-card/tooltip content, dialog, alert-dialog, card) swap `border` for `ring-1 ring-foreground/10` — a color/box-model change (ring sits outside the border box, so the inner content box shifts by 1px/side), not a metric-token substitution, so it stayed `border` throughout this retarget. Likewise nova's several display-model changes (e.g. alert's `grid-cols-[0_1fr]` + conditional grid-cols collapsing to a plain `grid` + `has-[>svg]:grid-cols-[auto_1fr]`) are prose/notes in the map, not applied `delta:` tokens — ledgered per-component as NOTE entries there, not implemented here. Both classes of change are candidates for a future, separately-scoped pass; this retarget's Global Constraint was "metric tokens ONLY."
- DELTA (border→ring-1 swaps NOT adopted — restated for searchability): see NOT ADOPTED above; `border` is retained everywhere nova moved to `ring-1`.
- BEHAVIOR NOTE (card footer, now unconditional band): nova's `CardFooter` is always a full-bleed bordered band (`border-t p-(--card-spacing)`, root gains `has-data-[slot=card-footer]:pb-0` so the band sits flush against the rounded bottom corners) rather than the old conditional `[.border-t]:pt-6` (border/padding only applied when a caller also passed a `border-t` class). Every `CardFooter` now renders with its top border and full padding regardless of caller classes — a visible behavior change for any existing caller that relied on the old opt-in-via-`border-t` convention, not just a spacing tweak. Ledgered in detail under `## card`.
- SHADOW-PRESENCE (removals where nova drops a shadow, no shadow additions): button outline variant's `shadow-xs`, input's `shadow-xs`, input-group's `shadow-xs`, checkbox's `shadow-xs`, button-group text's `shadow-xs`, dialog/alert-dialog content's `shadow-lg` (alert-dialog additionally has no shadow at all, ring-1 only), textarea's `shadow-xs`, native-select's `shadow-xs`, toggle-outline's `shadow-xs`, card's `shadow-sm`, and radio/switch's `shadow-xs` (`## radio`, `## switch`) are all removed to match nova's CSS, which drops them outright rather than replacing them with something else. No component gained a new shadow in this retarget.
- Verification: each of the 7 implementation tasks (`ND Task 1`–`ND Task 6` plus this Task 7) ran `go tool gsx generate && go test ./...`, re-pinning any failing class-string pins via `scripts/repin.py` against the map's stated deltas, then `make check`, one commit per task — the full adjudication trail is `.superpowers/sdd/progress.md`'s "NOVA DENSITY RETARGET" section.

## alert
- WIN: `cva()` variant map replaced by `switch` inside `class={}` (default |
  destructive), the same idiom as badge/button. No `data-variant` stamp —
  shadcn's own Alert doesn't stamp one either (unlike Badge/Button), so
  there's nothing to port.

## alert-dialog
- NOTE (source revision, scoping only — corrected 2026-07-24): the `size` prop (`"default"|"sm"`) and `AlertDialogMedia` that `registry/new-york-v4/ui/alert-dialog.tsx` has grown at its current HEAD (`f31ed8198`) are out of scope for this task (its own wording, e.g. "shadcn wraps buttonVariants the same way" for Action/Cancel, matches the pre-refactor shape — the earlier commit `f1dd9c690` "update dark mode colors" calls `buttonVariants()` directly, current HEAD wraps `<Button asChild>` with `variant`/`size` passthrough instead), the same scoping call as avatar's `AvatarBadge`/`AvatarGroup`/`AvatarGroupCount` (`## avatar`). **This is scoping only, not a source-revision substitution**: every class STRING below is the current HEAD's, with only the genuinely `size`/Media-conditional selectors stripped — see the per-part breakdown in the GAP entry below for exactly which tokens that means, part by part (an earlier version of this port ported Header/Description from the older, pre-refactor revision's class strings wholesale instead of stripping conditionals off the current one; Header's pre-refactor `flex flex-col gap-2 text-center sm:text-left` is a materially different — and wrong, per this binding-to-disk rule — recipe from current HEAD's unconditional `grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center` base, caught in review and fixed; Footer/Title/Content's current-HEAD unconditional bases happened to already coincide with what was ported, so those needed no code change, only this note; Description's token order was also corrected to match, no visual difference).
- WIN: `AlertDialog`/`AlertDialogTrigger`/`AlertDialogContent` compose `ui.Dialog`/`ui.DialogContent` directly instead of re-deriving a second `<dialog>` machinery — an alert dialog IS a dialog that (a) cannot be light-dismissed by an outside click and (b) never renders the injected close X; every other behavior (top-layer stacking, Esc-to-close, trigger/content proximity wiring, `data-state`, aria wiring) is identical, so `ui/dialog.js` is reused unmodified except for the one opt-out below. This is also what makes `alert-dialog`'s derived dependency chain `[button dialog]` (`registry.Deps("alert-dialog")`) pull `ui/dialog.js` transitively for the CLI, even though `alert-dialog` has no behavior module of its own (`HasJS("alert-dialog")` is false).
- MECHANISM (the one `ui/dialog.js` change this task makes): a `<dialog>` carrying `data-gsxui-dialog-static` is skipped by the backdrop-click light-dismiss handler — one early `if (dialog.hasAttribute("data-gsxui-dialog-static")) return;` ahead of the existing rect check. Esc (the `cancel` listener), the close button/`data-gsxui-dialog-close` handler, and the `toggle`-driven state sync are all untouched — only the backdrop-click path is gated. `AlertDialogContent` stamps the attribute unconditionally; any `ui.DialogContent` can opt into the same behavior directly.
- NOTE (parity, the reason for the mechanism above): Radix's own `AlertDialog` ignores outside/backdrop clicks entirely while Esc still closes it (`AlertDialogPrimitive.Content` sets `onPointerDownOutside`/`onInteractOutside` to `preventDefault` by default, `onEscapeKeyDown` does not). `data-gsxui-dialog-static` plus native `<dialog>`'s own Esc-fires-`cancel` behavior reproduces exactly that split — outside clicks are inert, Esc still works — without a second dismissal mechanism to keep in sync with Dialog's own.
- GAP: `AlertDialogPortal`/`AlertDialogOverlay` are not ported — same reasoning as `## dialog`'s own WIN: the native `<dialog>` top layer replaces Portal and `::backdrop` (already on `DialogContent`, inherited via composition here) replaces Overlay. Nothing under `AlertDialogContent` in this port ever needs a portal or a separate overlay element to target.
- GAP (out of scope, see the source-revision NOTE above): `AlertDialogContent`'s `size` prop (`"default"|"sm"`) and `AlertDialogMedia` are not ported — `data-size` and `data-slot="alert-dialog-media"` are never stamped anywhere in this port, so every selector conditioned on either can never match. Per part, against current HEAD (`f31ed8198`):
  - `AlertDialogHeader` drops only the three conditional tokens — `has-data-[slot=alert-dialog-media]:grid-rows-[auto_auto_1fr]`, `has-data-[slot=alert-dialog-media]:gap-x-6`, and the `sm:group-data-[size=default]/alert-dialog-content:*` triple (`place-items-start`, `text-left`, `has-data-[slot=alert-dialog-media]:grid-rows-[auto_1fr]`) — and keeps its **unconditional base verbatim**: `grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center`. This base is NOT dead weight (unlike the conditional selectors) — it is Header's actual layout, always centered (see the ADAPT entry below).
  - `AlertDialogFooter` drops `group-data-[size=sm]/alert-dialog-content:grid` and `group-data-[size=sm]/alert-dialog-content:grid-cols-2`; its unconditional base, `flex flex-col-reverse gap-2 sm:flex-row sm:justify-end`, is ported verbatim and happens to be byte-identical to `DialogFooter`'s own (`## dialog`) — coincidence, not composition (still a separate element, not `DialogFooter` composed).
  - `AlertDialogTitle` drops `sm:group-data-[size=default]/alert-dialog-content:group-has-data-[slot=alert-dialog-media]/alert-dialog-content:col-start-2`; its unconditional base, `text-lg font-semibold` in the source, is retargeted by the nova density retarget to `text-base font-medium` (see the SHADOW-PRESENCE/`## nova density` entries and this component's own class string) rather than ported verbatim.
  - `AlertDialogContent` drops `group/alert-dialog-content` (a marker class solely for the `group-data-[size=*]/alert-dialog-content:*` selectors above, all dropped — nothing left to target it), `data-[size=sm]:max-w-xs`, and `data-[size=default]:sm:max-w-lg` (conditional; `DialogContent`'s own unconditional `sm:max-w-sm` already supplies a comparable max-width regardless — see the ADAPT entry below, which now DOES pass its own `max-w-xs sm:max-w-sm` tokens per the nova density retarget UPDATE).
  - `AlertDialogDescription` is unaffected by the `size`/Media refactor (`text-sm text-muted-foreground` unconditional in both revisions).
- ADAPT (Header layout, always centered): `AlertDialogHeader`'s `grid grid-rows-[auto_1fr] place-items-center gap-1.5 text-center` base — two grid rows (title then description, in DOM/source order; `grid-auto-flow`'s default `row` placement needs no explicit `grid-row` on either child), centered on both axes — renders identically at every viewport width in this port, since the one thing that would left-align it at `sm`+ (the `size=default`-conditional pair dropped above) is out of scope. Visually verified via `site/examples/alertdialog/basic.gsx`'s rendered markup: title stacks directly above description, both horizontally centered, no left-alignment breakpoint — matches shadcn's own `size="sm"` appearance at every width (the one variant this port effectively always renders, since `size` itself isn't ported) rather than its `size="default"` (the actual shadcn default, left-aligned at `sm`+) — a real, deliberate visual GAP downstream of the `size`-prop scoping decision above, not a copy error.
- ADAPT: `AlertDialogContent` passes no `class` attr into its `ui.DialogContent` call at all. Diffed token-for-token, every utility the (pre-refactor) source's content class carries — `bg-background`, the fixed-centering/sizing set, `sm:max-w-lg`, and all six `data-[state=…]:animate-in/out`/`fade-*`/`zoom-*` tokens — is already present in `DialogContent`'s own base class (both dialogs share one centered-card recipe upstream). The one token that is not shared, a bare `grid`, is dropped rather than re-supplied: `DialogContent`'s own `open:grid` exists specifically so content stays `display:none` while the native `<dialog>` is closed (`## dialog`'s own ADAPT); an unscoped `grid` alongside a `open:`-scoped one is not a tailwind-merge conflict (variant scope is part of the conflict key — the same non-collision as accordion's rotate override, `## accordion`) and would silently defeat that mechanism. A first pass tried literally re-passing the full (minus-`grid`) string anyway to exercise the merge path per the task brief; it round-tripped every positional/spacing utility correctly (later occurrence wins, matching `## accordion`'s documented conflict semantics) but left the six `data-[state=…]` animation tokens **duplicated** in the output — the configured merger has no conflict-group entry for `tw-animate-css`'s classes at all (`## animations`' FINDING), so two byte-identical copies of the same token are not recognized as a merge conflict and both survive. Passing nothing sidesteps the duplication and produces the same net class as the (correctly deduped) alternative — at the time this entry was written, `AlertDialogContent`'s rendered class was exactly `DialogContent`'s own default. UPDATE (2026-07-24, nova density retarget): no longer the whole story — `AlertDialogContent` now DOES pass one class string, `class="max-w-xs sm:max-w-sm"`, into its `ui.DialogContent` call, merged through the same tailwind-merge pass reasoned about above (the `grid`/animation-token reasoning for passing nothing else is unaffected). This overrides `DialogContent`'s own `max-w-[calc(100%-2rem)] ... sm:max-w-sm` width pair with the nova density map's alert-dialog delta (`sm:max-w-lg → max-w-xs` + `sm:max-w-sm`, `docs/superpowers/plans/2026-07-24-nova-density-map.md` `## alert-dialog`) — alert dialogs render narrower than plain dialogs below the `sm` breakpoint, matching nova's `size="default"` width without porting the `size` prop/axis itself (still out of scope, GAP above). Every other non-variant token (`p-4`, `rounded-xl`, `border`, `text-sm`, the animate-in/out set) still arrives for free via composition, unchanged. Pinned as such in `TestAlertDialogContentPinned` — `role="alertdialog"`, `data-gsxui-dialog-static`, `hideCloseButton`'s permanently-true injected-X suppression, and now the max-width override distinguish an alert dialog from a plain one.
- WIN (data-slot overrides): `AlertDialog`, `AlertDialogContent`, `AlertDialogAction`, and `AlertDialogCancel` each override their composed component's own `data-slot` (`"dialog"`→`"alert-dialog"`, `"dialog-content"`→`"alert-dialog-content"`, `"button"`→`"alert-dialog-action"`/`"alert-dialog-cancel"`) via the same explicit-non-parameter-attribute mechanism `ItemSeparator`/`FieldLabel` use to override `Separator`/`Label`'s own slot (`## item`, `## field`) — passed as an ordinary literal attr ahead of `{ attrs... }`, which the composed component's own `if !attrs.Has("data-slot")` guard then skips.
- MECHANISM: `AlertDialogTrigger` renders its own `<button>` rather than composing `DialogTrigger` — same button-in-button HTML trap, same fix (`## dialog` FINDING): children must be phrasing content, and a styled trigger should be `<ui.Button data-gsxui-dialog-trigger>` directly, not a `Button` nested inside this wrapper. The site example (`site/examples/alertdialog/basic.gsx`) demonstrates the recommended idiom rather than the wrapper, mirroring `ui/dialog`'s own example.
- GAP (narrow): `asChild` is not ported anywhere in this component — `AlertDialogAction`/`AlertDialogCancel` wrap `ui.Button` directly (composition, not Slot-cloning) instead of shadcn's `<Button asChild><AlertDialogPrimitive.Action/></Button>`, and `data-gsxui-dialog-close` is the attribute-idiom replacement for what `asChild` wires up in the original. Same narrow gap already covered by `## dialog`'s MECHANISM entry.
- Registry: `alert-dialog.gsx` imports nothing from `ui/icon`; `AlertDialog` calls `ui.Dialog` and `AlertDialogAction`/`AlertDialogCancel` call `ui.Button` (flat-package intra-package edges, same shape as `## dialog`'s own `button` dep) — `registry.Deps("alert-dialog") == ["button", "dialog"]`, pinned in `internal/registry/registry_test.go`.

## aspect-ratio
- ADAPT: shadcn's `AspectRatio` is a bare passthrough onto Radix's `AspectRatioPrimitive.Root`, which renders the classic padding-hack pair — an outer `position:relative; width:100%; padding-bottom:calc(100% / ratio)` div (`ratio` set as an inline style computed from a numeric `ratio` prop) wrapping an inner `position:absolute; inset:0` div that actually holds the children. This port drops both divs and both computed styles for the single-property CSS `aspect-ratio` declaration on one div — no wrapper-within-wrapper, no percentage arithmetic to reproduce, and (a real capability gap Radix's hack technique predates) native browser support for intrinsic-size media inside the ratio box that the padding-hack couldn't give for free.
- ADAPT: `ratio` is a Go `string`, not a `float`/`number` — callers write the literal CSS `aspect-ratio` value, e.g. `ratio="16 / 9"` (`aspect-ratio` also accepts a bare number, `ratio="1.5"`), instead of a numeric ratio the component would have to stringify itself.
- MECHANISM: the composed value is wrapped in `gsx.RawCSS(ratio)` to opt it out of gw's CSS value filter — a conservative character-blocklist port of `html/template`'s CSS sanitizer that rejects `/` outright (along with `(`, `)`, `;`, and several other characters that never appear in a valid `aspect-ratio` value either). `/` is not incidental punctuation here; it is `aspect-ratio`'s own required `<width> "/" <height>` separator, so no value most callers would ever write could pass the filter unmodified. `ratio` is treated as trusted, developer-authored layout intent — the same trust boundary every `ui/*.gsx` component already extends to its own class strings (Tailwind's `class` attribute receives no injection filtering either), not sanitized end-user request data.

## avatar
- ADAPT: AvatarImage adds `absolute inset-0` (not in shadcn) — the image overlays the in-flow fallback, so the no-JS/pre-JS state renders correctly (fallback behind, image covers when loaded); `ui/avatar/avatar.js` syncs image/fallback visibility via capture-delegated `load`/`error`. Radix's mount-gated rendering can't exist server-side. RACE (verified in-browser, fixed): images that settle before this module imports (data URIs, cached, fast-local) already fired their `load`/`error` event before the delegated listener existed, so neither image nor fallback visibility ever synced — resolved with a settle-sweep (`img.complete` check over `[data-gsxui-avatar-image]`) run once at import and again on `window` `load`, covering the gap between; the capture-delegated handlers still cover every settle after import, including HTMX swaps.
- GAP: `AvatarBadge`, `AvatarGroup`, `AvatarGroupCount` (added to shadcn's
  registry after the base three parts) are not ported — out of scope for
  this task; only `Avatar`/`AvatarImage`/`AvatarFallback` per the task
  brief.
- ADAPT: `Avatar`'s `size` prop (default/sm/lg) is dropped along with it, so
  `data-size` is never stamped and the size-keyed selectors that depend on
  it are dead weight — `group/avatar` and `data-[size=lg]:size-10
  data-[size=sm]:size-6` are dropped from `Avatar`'s class, and
  `group-data-[size=sm]/avatar:text-xs` from `AvatarFallback`'s (the same
  "drop the selector, don't ship dead CSS" call as dialog's close-button
  `data-[state=open]:...` ADAPT). `size-*` stays fully overridable via the
  ordinary caller-class-merge mechanism on `Avatar`/`AvatarFallback`
  directly.
- FINDING: gsx's image-sink URL sanitizer requires the literal `;base64,` marker in data: URLs — percent-encoded data:image/svg+xml,... forms are blocked to about:invalid#gsx. The ergonomic authoring path is gsx's own `dataURL` std filter: author the image as plain `[]byte` and write `src={svgBytes |> dataURL("image/svg+xml")}` — the filter assembles the base64 form and the sink re-validates it (the f-string form ``src=f`data:image/png;base64,@{bytes}` `` also works; `gsx.RawURL` remains the per-value vouch for a pre-built URL). The avatar examples use the filter. The block itself is still silent — no generate-time diagnostic, the broken image is the first signal; upstream proposal filed as gsxhq/gsx#154 (generate-time diagnostic for always-blocked literals + strictly-validated percent-encoded acceptance).

## badge
- WIN: `cva()` variant map replaced by `switch` inside `class={}`.
- GAP (narrow): shadcn's `asChild` tag-swapping (render the badge as an `<a>`)
  has no gsx equivalent — no dynamic tag. Behavior-attachment uses of
  `asChild` are covered by the data-attribute mechanism (see dialog).

## breadcrumb
- Straight port; no dropped tokens. shadcn's own `breadcrumb.tsx` has no Radix primitive underneath either — every part (`Breadcrumb`/`BreadcrumbList`/`BreadcrumbItem`/`BreadcrumbLink`/`BreadcrumbPage`/`BreadcrumbSeparator`/`BreadcrumbEllipsis`) is already a plain styled `<nav>`/`<ol>`/`<li>`/`<a>`/`<span>`, so this is the same "package-namespaced compound parts" shape as card, not an ADAPT.
- GAP (narrow): `BreadcrumbLink`'s `asChild` tag-swapping is dropped — it always renders a real `<a>`, which is shadcn's own default (`const Comp = asChild ? Slot.Root : "a"`) for the dominant/only realistic use anyway. Same narrow gap as button's `asChild`; behavior-attachment uses are covered by the data-attribute mechanism (see dialog).
- MECHANISM: `BreadcrumbSeparator`'s default child (shadcn's `{children ?? <ChevronRight />}`) and `BreadcrumbEllipsis`'s icon (`MoreHorizontal`) both come from `ui/icon` — `icon.ChevronRight` and `icon.Ellipsis`. Lucide renamed `MoreHorizontal` to `"ellipsis"` (verified against `ui/icon/icon_data.go`'s `"ellipsis"` entry: three `<circle>`s, the horizontal three-dot glyph, not `EllipsisVertical`'s rotated form) — the same rename precedent as Spinner's `Loader2Icon`/`icon.LoaderCircle` (see `## spinner`). This import is the `breadcrumb` → `icon` dependency `internal/registry` derives and `registry_test.go` pins.

## button
- GAP (narrow): `asChild` tag-swapping (no dynamic tag: `const Comp = asChild ? Slot : "button"`). Ported as `href` param rendering `<a>` — covers the dominant use. Behavior-attachment uses of `asChild` are covered by the data-attribute mechanism (see dialog).
- WIN: `type="button"` before `{ attrs... }` makes it an overridable default — positional spread precedence replaces prop-ordering conventions.
- WIN: `cva()` replaced by plain Go variant/size funcs shared by both branches.

## button-group
- WIN: shadcn's `buttonGroupVariants` cva map (`orientation`: horizontal/vertical) picks between two entirely static class blocks by the JS-resolved prop value — there are no `data-[orientation=...]:` selectors in `registry/new-york-v4/ui/button-group.tsx` to preserve. Ported via a `switch` inside `class={}`, the same idiom as badge/alert, not a CSS-side data-attribute selector.
- ADAPT: `data-orientation` is still stamped on `ButtonGroup`'s root via the house `|> default("horizontal")` pattern (see button.gsx/dropdown.gsx), for consistency with every other data-variant stamp in this codebase. shadcn's own `data-orientation={orientation}` leaves the attribute entirely unset when `orientation` is undefined (no defaulting on the raw prop, only inside the cva call) — a narrow, cosmetic divergence.
- GAP (narrow): `ButtonGroupText`'s `asChild` tag-swapping is dropped, same shape as button's own `asChild`; always renders a `<div>`.
- NOTE: `ButtonGroupText` carries no `data-slot` in shadcn's own source either (every other button-group part does) — ported as-is, not "fixed," per the token-for-token rule.
- WIN: `ButtonGroupSeparator` calls `ui.Separator` directly (flat package, no re-implementation) — the `button-group` → `separator` dependency `internal/registry` derives and `registry_test.go` pins. Its `orientation = "vertical"` default (the opposite of `Separator`'s own `"horizontal"`) is resolved before the call via the same `|> default` mechanism, since `Separator`'s own default is baked into its own component.

## card
- Straight port; package-namespaced compound parts (`card.CardHeader`) replace module exports. No divergences.

## collapsible
- WIN: Radix `CollapsiblePrimitive.Root`/`Trigger`/`Content` (client `open`/`onOpenChange` state, no visual output of its own) is replaced by a single un-grouped native `<details>`/`<summary>` pair — exactly Accordion's mechanism (see `## accordion` above), minus Accordion's `name`-grouping concern: shadcn's `collapsible.tsx` is a lone Root+Trigger+Content triple, not a Root+Item pair, so there is no exclusive-open-within-a-group behavior to model and no `name` param on `Collapsible`.
- WIN: `Collapsible`'s `open bool` server-renders the initial expanded/collapsed state as the native `open` boolean attribute (Go zero value `false` = collapsed, matching shadcn's `defaultOpen` unset) — it is `defaultOpen`, not a controlled `open`/`onOpenChange` pair; thereafter, toggling is entirely native `<details>` behavior, no hydration step.
- NOTE: `registry/new-york-v4/ui/collapsible.tsx` is a bare `data-slot` passthrough — none of its three parts carry a single class string. `Collapsible` and `CollapsibleContent` port with zero classes of their own (nothing to carry token-for-token, and nothing added).
- ADAPT: `CollapsibleTrigger` gains `list-none [&::-webkit-details-marker]:hidden`, classes absent from shadcn's source, for the same reason Accordion's trigger carries them (`## accordion` ADAPT above): a real `<summary>` draws a native disclosure triangle marker (`::marker`, suppressed by `list-none`; WebKit's separate `::-webkit-details-marker`) that Radix's `<button>`-based trigger never had to suppress. Without both, shadcn's chevron renders alongside a browser-drawn triangle.
- MECHANISM (asChild replaced by direct composition, not the data-attribute idiom): shadcn's own demo (`registry/new-york-v4/examples/collapsible-demo.tsx`) wraps a real `<Button>` inside `<CollapsibleTrigger asChild>` so an actual button element becomes the trigger while Radix's Trigger contributes only behavior/ARIA. Here neither `asChild` nor the `data-gsxui-*`-attribute idiom (dialog's MECHANISM, `## dialog` above) applies: a bare `<summary>` already IS the clickable disclosure control (native semantics, nothing to attach a handler to), and `<button>` is valid phrasing content inside `<summary>` (unlike `DialogTrigger`'s button-in-button trap, `## dialog` FINDING) — so callers compose `<ui.Button>` (or anything) directly as `CollapsibleTrigger`'s child with no wrapper, cloning, or special attribute at all: `<ui.CollapsibleTrigger><ui.Button variant="ghost" size="icon"><icon.ChevronsUpDown/></ui.Button></ui.CollapsibleTrigger>`.
- GAP: no `data-state` is stamped anywhere in this port — there is no `collapsible.js`, same as Accordion has no `accordion.js`. shadcn ships no `data-[state=open]`/`data-[state=closed]` selectors on Collapsible in the first place (unlike Accordion's Radix-sourced `animate-accordion-down`/`up`, which motivated Accordion's CSS-only `::details-content` replacement — see `## accordion` MECHANISM), so there is no animation block to port or replace: CSS consumers wanting open/closed-driven styling target the ancestor `<details>`'s native `[open]` attribute instead, the same substitution Accordion's chevron rotation uses (`[[open]_&]:...`-shaped selectors, not `data-[state=open]:...`).
- Registry: `collapsible.gsx` imports nothing from `ui/icon` and calls no other component — `registry.Deps("collapsible")` is empty, same shape as `kbd`/`aspect-ratio`/`progress`. (The site example composes `ui.Button` and `ui/icon`, but `internal/registry` only scans `ui/*.gsx`, not `site/examples/`, so that composition has no effect on derived deps.)

## dialog
- WIN: Radix Portal/Overlay replaced by native <dialog> top layer + ::backdrop; Esc handling is browser-native.
- ADAPT: DialogContent uses `open:grid` instead of shadcn's `grid` — content stays in the DOM when closed (no Radix unmount), so an ungated display utility would override the UA's closed-dialog `display:none`.
- ADAPT: `text-foreground` added to DialogContent's classes — native <dialog> gets UA `color: CanvasText` and does not inherit the themed body color (Radix's <div> content does); without it dark mode renders wrong text color.
- ADAPT: close button drops shadcn's `data-[state=open]:bg-accent data-[state=open]:text-muted-foreground` — we stamp `data-state` on the <dialog> element, not the close button, so those selectors are dead in this DOM.
- CONVENTION (decided 2026-07-22): gsx keeps Go zero-value semantics — no default parameter values in the language (designs using exported consts and an init-time registry were considered and rejected: exported-symbol pollution vs. runtime lookup, neither worth it). Name parameters so the zero value is the shadcn default: bools invert (`hideCloseButton`), numerics document "0 means the default" (e.g. upcoming `sideOffset`), strings use `""` = default.
- MECHANISM (decided 2026-07-22): `asChild`/Slot is not ported and not needed for behavior attachment — the `data-gsxui-*` attributes are each interactive component's PUBLIC contract, and fallthrough attrs deliver them to any element or component: `<ui.Button data-gsxui-dialog-trigger>Open</ui.Button>` makes your styled Button the trigger, no cloning, no wrapper. Document this idiom prominently per interactive component; only tag-swapping remains a (narrow, accepted) gap.
- FINDING (2026-07-24, fixed): nesting `<ui.Button>` inside `<ui.DialogTrigger>` — the shape shadcn's `asChild` invites — is invalid HTML (button in button): the parser hoists the inner button out as a SIBLING, so the wired trigger renders empty/unclickable and the visible button is orphaned. The landing page shipped this and its dialog silently never opened. The data-attribute idiom above is not just nicer, it is the only correct form for a styled trigger. `TestNoNestedButtons` renders every page and fails on any button-in-button markup; DialogTrigger's doc comment now warns its children must be phrasing content. (Same emit-what-the-parser-rejects family as self-closed non-voids, gsxhq/gsx#144 — a nested-interactive-content lint could catch both upstream.)
- GAP: Radix client context (trigger↔content wiring) replaced by closest("[data-gsxui-dialog]") proximity in JS.
- NOTE: controlled open/onOpenChange not ported; JS CustomEvents (gsxui:open/close) + dialog.showModal()/close() are the programmatic API. State + events ride ToggleEvent (Baseline 2024) — the native close event proved undelivered in current Chrome, and toggle also covers programmatic open/close, which close-based wiring never could.
- A11Y: aria-labelledby/-describedby/-controls stamped by JS with lazily generated ids (authored ids/aria always win); aria-haspopup + initial aria-expanded server-rendered on DialogTrigger; aria-expanded synced on toggle.
- ADAPT: exit animations run by stamping data-state="closed" before close() (Esc included, via intercepted cancel); a direct programmatic dialog.close() skips the exit animation — native-immediacy divergence, accepted.
- WIN: DialogFooter's showCloseButton ports with zero friction — shadcn defaults it false, so the Go zero value IS the default; its Close button uses the data-attribute idiom (<ui.Button data-gsxui-dialog-close>) where shadcn needs <DialogClose asChild><Button>>.

## checkbox
- ADAPT (native-first): Radix's `CheckboxPrimitive.Root` (button role, `aria-checked`, hidden real input, `data-[state=checked]`) is replaced by a real `<input type="checkbox">` — form-native, zero JS, browser `:checked`/`:disabled` truth. `data-slot="checkbox"` and every color/focus/disabled/aria-invalid token from `registry/new-york-v4/ui/checkbox.tsx` are carried over verbatim; `appearance-none` is added to suppress the UA checkbox box so the custom border/background show instead (mechanical necessity of the native-input swap, not present on the Radix version since it renders no native control at all).
- MECHANISM: the `Indicator`/`CheckIcon` child is not a child at all on a void `<input>` — its checkmark becomes a `checked:bg-[url('data:image/svg+xml...')]` data-URI background (lucide's `check` path, `M20 6 9 17l-5-5`, verified byte-for-byte against `ui/icon/icon_data.go`'s `"check"` entry) plus `checked:bg-center checked:bg-no-repeat checked:bg-[length:12px_12px]`, swapped in only under the `checked:` variant. `data-[state=checked]:text-primary-foreground` is dropped (it colored the now-nonexistent Indicator child's icon; our glyph color lives inside the data-URI strokes instead).
- FINDING (2026-07-24, fixed): `dark:data-[state=checked]:bg-primary` was ALSO dropped as "redundant" — wrongly. The dark custom variant compiles to `:is(.dark *)`, whose class specificity (0,2,0) beats bare `:checked` (0,1,1), so without the explicit dark override `dark:bg-input/30` won the cascade and a checked box rendered 4.5%-alpha in dark mode instead of primary. Restored as `dark:checked:bg-primary dark:checked:border-primary`, plus a second dark check URI stroking the dark `--primary-foreground` value (`oklch(0.205 0 0)`) — primary flips near-white in dark, where the light URI's white stroke would vanish. Lesson: shadcn's seemingly-duplicate dark overrides are cascade arbitration against `dark:bg-input/30`, not noise. Both strokes remain static text (data-URIs can't read CSS vars) — the currentColor-mask backlog item still covers arbitrary custom themes.
- FINDING (verified, corrected twice — see checkbox_test.go's `TestCheckboxDataURIDecodesToValidSVG`): the data-URI's embedded SVG markup needs spaces (attribute and path-data separators), but a literal space inside a `class="..."` value is a token boundary by the HTML spec itself — every whitespace-splitting class tool, including the configured `merge.Merge` (`tailwind-merge-go`), tears the token apart at each space. Reproduced directly against `merge.Merge`: with real spaces, `checked:bg-primary` was dropped entirely (its `bg` category lost to the accidentally-split `bg-[url(...` fragment) and interior tokens `stroke-width`/`stroke-linecap` were dead-code-eliminated as fake "conflicting" utilities against `stroke-linejoin`, corrupting the SVG. The first fix wrote the spaces as `_`, Tailwind's escape for whitespace in arbitrary values — WRONG for this token: Tailwind v4 deliberately does NOT convert underscores inside `url()` values (real URLs contain them), so the underscores reached the browser verbatim, the SVG parsed as `<svg_xmlns=...` (invalid), and the checkmark silently never painted — a solid primary square when checked (shipped broken; caught by eye, invisible to render-string tests that pinned the underscore form). The second fix percent-encoded the spaces as `%20` — which held until the dark-variant URI needed a stroke of `oklch(0.205 0 0)`: percent-encoded parens in the emitted url() broke vite's postcss parse of Tailwind's output (and the broken candidate then stuck in the dev server's accumulate-only Tailwind candidate cache until a restart). Final form: **base64 payloads** (`data:image/svg+xml;base64,...`) — `[A-Za-z0-9+/=]` only, nothing for tailwind-merge, the class scanner, postcss, or the browser to split, convert, or mis-parse. The pin test proves the browser's view instead of a byte form: it extracts each rendered URI, base64-decodes it, and XML-parses the result. (`bg-[length:12px_12px]` in the same string correctly keeps `_` — it is not a `url()` value, so Tailwind converts it to a space; same for radio's `radial-gradient` underscores below.)
- ADAPT (theming constraint, accepted): the check glyph's data-URI hard-codes `stroke="white"` in its embedded SVG markup. A data-URI is static text baked into the class string at build time — it cannot reference CSS custom properties, so the mark can never follow `--primary-foreground` the way radio's `currentColor`-based gradient follows `--primary` (see radio's MECHANISM entry above, which explains exactly why radio could take the live-`currentColor` route and checkbox's data-URI approach can't). For themes where `--primary-foreground` isn't near-white, the checkmark is wrong/low-contrast against the checked background. Backlog: swap to a `currentColor` CSS-mask (`mask-image`/`-webkit-mask-image` referencing the check SVG, painted via `background-color: currentColor`) in the Plan 4 theming work — masks, unlike data-URI `fill`/`stroke`, paint with the element's live computed color.

## radio
- ADAPT (native-first): Radix's `RadioGroupPrimitive.Item` (button role, `aria-checked`, Indicator) is replaced by a real `<input type="radio">` — same rationale as checkbox: form-native, zero JS, browser `:checked`/`:disabled` truth. Token set carried verbatim from `registry/new-york-v4/ui/radio-group.tsx`'s `RadioGroupItem` (`aspect-square size-4 shrink-0 rounded-full border border-input text-primary shadow-xs transition-[color,box-shadow] outline-none` plus focus-visible/disabled/aria-invalid/dark tokens); `appearance-none` is added for the same mechanical reason as checkbox.
- ADAPT (native-first): shadcn's `RadioGroup` root (a `grid gap-3` wrapper coordinating Radix's roving-tabindex/keyboard-nav group) is not ported. Native `<input type="radio" name="...">` siblings already form a group via the browser's own `name` attribute — arrow-key roving and single-selection are UA behavior, no JS required. The `grid gap-3` layout is not component behavior; it is the caller's own wrapper `<div>`.
- MECHANISM: the `Indicator`/`CircleIcon` child (`fill-primary`) is not a child at all on a void `<input>` — its filled dot becomes `checked:bg-[radial-gradient(circle_closest-side,currentColor_45%,transparent_50%)]`, painted in `currentColor` rather than checkbox's data-URI approach: a data-URI is static text and can't reference the caller's CSS custom properties (`--primary` can be themed/overridden per caller, per dark mode, etc.), but a `currentColor`-based gradient computed live on the element can. `text-primary` is kept (not dropped) specifically to make `currentColor` resolve to the primary color on this element — it is load-bearing for the dot's color, unlike checkbox's icon color which came from an explicit `fill`/`stroke` in the data-URI itself. (An earlier version of this port wrongly copied checkbox's fill-the-whole-circle-plus-white-dot recipe — `checked:bg-primary checked:border-primary` plus a white-circle data-URI — which does not match shadcn's actual outlined-circle-with-colored-dot visual; corrected here.) The gradient's spaces use Tailwind's `_` escape — valid HERE precisely where it was wrong for checkbox's data-URI: this is not a `url()` value, so Tailwind converts `_` to spaces in the emitted CSS (see checkbox's FINDING entry above) — verified the radial-gradient string round-trips `merge.Merge` unchanged (`TestRadioNoStraySpaceInGradient`).

## switch
- FINDING (2026-07-24, fixed — found by the dark-token audit, checkbox's twin): shadcn gates its dark track color to `dark:data-[state=unchecked]:bg-input/80`; our port carried it ungated as `dark:bg-input/80`, which outweighs `checked:bg-primary` in dark mode (the dark custom variant's `:is(.dark *)` carries class specificity (0,2,0) vs bare `:checked` (0,1,1)) — a checked track lost its primary fill in dark. Fixed with the explicit-override idiom (`dark:checked:bg-primary`), consistent with checkbox. The audit diffed every component's `dark:*` tokens against its shadcn source: switch was the only remaining hit (tabs' two missing dark tokens belong to the unported `line` variant, ledgered below).
- ADAPT (native-first): Radix's `SwitchPrimitive.Root` + separate `SwitchPrimitive.Thumb` span replaced by one real `<input type="checkbox" role="switch">` — form-native, zero JS, browser `:checked`/`:disabled` truth. `role="switch"` is kept explicitly (a bare checkbox input has no switch semantics of its own) since it's the one piece of Radix's ARIA contract a native checkbox doesn't supply for free. The `data-size="sm"|"default"` variant and its `group/switch` + `group-data-[size]/switch:` machinery are dropped — Task 5 ships default size only (no size variant asked for), and `group/switch` existed solely to let a *sibling* Thumb element read the Root's size; with the thumb now the same element's own pseudo-element, there is no sibling to target and the group plumbing is dead weight.
- MECHANISM: the Thumb span becomes this input's own `::before` pseudo-element (`before:` variant, thumb-span→before:) — track and thumb are the same DOM node instead of parent/child. A pseudo-element renders nothing without an explicit `content` (default `content: normal` produces no box at all, unlike a real child element which always has one) — Tailwind's `before:content-['']` is therefore load-bearing here and has no analog on the Radix Thumb span. Native `checked:` replaces `data-[state=checked]`/`data-[state=unchecked]:` throughout; an unchecked-specific class is unnecessary (the bare, unprefixed utility already covers the unchecked default, exactly as for checkbox/radio). `ring-0` itself is dropped as a dead reset (nothing `before:` could have a ring here).

## empty
- Straight port; no dropped tokens. shadcn's own `empty.tsx` has no Radix primitive underneath either — every part (`Empty`/`EmptyHeader`/`EmptyMedia`/`EmptyTitle`/`EmptyDescription`/`EmptyContent`) is already a plain styled `<div>`, the same "package-namespaced compound parts" shape as card/breadcrumb/pagination.
- WIN: `EmptyMedia`'s `emptyMediaVariants` cva map (`variant`: default/icon) picks between two entirely static class blocks by the JS-resolved value — no `data-[variant=...]` selectors in `registry/new-york-v4/ui/empty.tsx` to preserve — so it ports as a `switch` inside `class={}`, the same idiom as badge/button-group/item.
- NOTE: `EmptyMedia`'s `data-slot` is `"empty-icon"` in shadcn's own source, not `"empty-media"` — ported as-is (token-for-token), the same unmatched-`data-slot` call as `ButtonGroupText` (see `## button-group`).
- NOTE: `EmptyDescription` renders a `<div>`, matching shadcn's own actual returned element — its TypeScript prop type reads `React.ComponentProps<"p">` but the JSX it returns is a `<div>`, the same shipped-type/element mismatch already noted for `Kbd`/`KbdGroup` (see `## kbd`). Ported verbatim, tag included, per the token-for-token rule — contrast with `ItemDescription` (`## item`), whose shadcn source really does render a `<p>`.

## field
- Straight port; no dropped tokens. shadcn's own `field.tsx` has no Radix primitive underneath its own parts either (`Label`/`Separator` are composed, not Radix wrappers) — the same "package-namespaced compound parts" shape as card/breadcrumb/empty/item.
- WIN: `FieldLegend`'s `variant` prop only ever *types* the value — both `data-[variant=legend]:text-base` and `data-[variant=label]:text-sm` are present unconditionally in shadcn's own class string; `data-variant` plus Tailwind's `data-[variant=...]` selectors are what actually pick one. No switch needed — the same "single verbatim class string dispatches on the stamped attribute" shape as `Separator`'s `data-orientation` (see `## separator`), NOT a static-block cva switch like item/button-group/empty's variant maps.
- WIN: `Field`'s `fieldVariants` cva map (`orientation`: vertical/horizontal/responsive) picks between three static class blocks by the JS-resolved value — ported as a `switch` inside `class={}`, the usual idiom. UNLIKE button-group's `orientation` (nothing downstream ever reads its `data-orientation`), `Field`'s own `data-orientation` IS read by a sibling: `FieldDescription`'s `group-has-[[data-orientation=horizontal]]/field:text-balance` selector keys directly off it — so both the switch AND the `|> default("vertical")` attribute stamp are load-bearing here, not merely stamp-everything consistency.
- WIN: `FieldLabel` composes `ui.Label` and `FieldSeparator` composes `ui.Separator` directly (flat package, no re-implementation) — the `field` → `[label separator]` dependency `internal/registry` derives and `registry_test.go` pins. Both override the composed component's own `data-slot` (`"label"`→`"field-label"`, `"separator"` stays on the inner `Separator` — `FieldSeparator`'s own root carries `data-slot="field-separator"` instead) via the same explicit-non-parameter-attribute mechanism as `ItemSeparator`/`ButtonGroupSeparator` overriding `Separator`'s own `data-slot` (see `## item`, `## button-group`).
- NOTE: `FieldTitle` renders a `<div>` sharing `FieldLabel`'s exact `data-slot` value (`"field-label"`) in shadcn's own source — not a distinct `"field-title"` slot. Ported as-is (token-for-token), the same unmatched/shared-`data-slot` call as `EmptyMedia`'s `"empty-icon"` (`## empty`) and `ButtonGroupText`'s missing `data-slot` (`## button-group`).
- MECHANISM: `FieldSeparator`'s `data-content={!!children}` boolean stamp ports as `data-content={children != nil}` — gsx renders a bool expression as literal `"true"`/`"false"` attribute text, the same mechanism as `pagination.gsx`'s `data-active={isActive}` (see `## pagination`). The optional `field-separator-content` label span is gated by the same `children != nil` check via `{ if … }`.
- ADAPT: `FieldError`'s react-hook-form `errors` prop (`Array<{message?: string} | undefined>`, deduplicated and rendered as either a single message or a `<ul>` via `useMemo`) is dropped — there is no react-hook-form in a server-rendered gsx tree to produce that shape. `FieldError` keeps shadcn's other content path, plain `children`, unconditionally; a caller with more than one message renders its own `<ul>` child (the exact markup shadcn's own multi-error branch would have produced). `FieldError` renders nothing (not even a wrapper) when `children` is nil — the gsx `{ if cond { … } }` control-flow form (no `else`) standing in for shadcn's `if (!content) return null`.

## icon
- WIN: shadcn/templUI ports wrap each Lucide React component (or a `<template>`-per-icon component) individually; gsx's tag-callable values (`func(attrs ...gsx.Attr) gsx.Node`) let a single generated `New(name)` factory back every icon var (`var ChevronDown = New("chevron-down")`), so `<icon.ChevronDown class="size-4"/>` is both markup-callable and a plain Go value, generated from one shared `svgIcon` component instead of 1,748 near-duplicate wrapper components.
- WIN: `aria-hidden="true"` is authored before `{ attrs... }` in `svgIcon` — positional spread precedence (the same idiom as badge's `data-variant` and dialog's `data-state`) makes it an overridable default: a caller's own `aria-hidden` (e.g. `aria-hidden="false"` alongside `aria-label`) wins with no conditional logic.
- MECHANISM: unknown icon names are a render-time error (`New("nope")` → `unknown icon "nope"`), never a silently empty `<svg>` — mirrors the hard-error idiom used elsewhere in gsxui for unrecognized identifiers, so a typo'd icon name fails loudly instead of shipping a blank glyph.

## input
- Straight port. `type="text"` is authored before `{ attrs... }` — the same overridable-default idiom as button's `type="button"`, so `type="email"` etc. at the call site replaces rather than duplicates it.

## input-group
- Straight port; no dropped tokens. shadcn's own `input-group.tsx` has no Radix primitive underneath either — every part is already a plain styled element wrapping `Input`/`Textarea`/`Button`.
- WIN: `InputGroupAddon`'s `inputGroupAddonVariants` cva map (`align`: inline-start/inline-end/block-start/block-end) picks between four static class blocks by the JS-resolved value — no `data-[align=...]:` selectors in `registry/new-york-v4/ui/input-group.tsx` to preserve for *this* class string, so it ports as a `switch` inside `class={}`, the usual idiom. `data-align` is still stamped (`|> default("inline-start")`) — not merely for stamp-everything consistency, but because `InputGroup`'s OWN class string (a separate, sibling class computation) keys off it directly: `has-[>[data-align=inline-start]]:[&>input]:pl-2` and its three siblings.
- WIN: `InputGroupInput`/`InputGroupTextarea`/`InputGroupButton` compose `ui.Input`/`ui.Textarea`/`ui.Button` directly (flat package, no re-implementation) — the `input-group` → `[button input textarea]` dependency `internal/registry` derives and `registry_test.go` pins. `InputGroupInput`/`InputGroupTextarea` override the composed component's own `data-slot` (`"input"`/`"textarea"` → `"input-group-control"`) via the same explicit-non-parameter-attribute mechanism as `ItemSeparator` overriding `Separator`'s `data-slot` (see `## item`) — the attribute `InputGroup`'s own `has-[[data-slot=input-group-control]:focus-visible]:...` selectors key off.
- MECHANISM: `InputGroupButton` mirrors a subtlety in shadcn's own source rather than "fixing" it: its `size` prop is typed `Omit<ComponentProps<typeof Button>, "size">` — deliberately NEVER forwarded to the underlying `Button`'s own `size` prop. `Button` renders with its own default size classes, and `inputGroupButtonVariants({size})`'s overlay classes are merged on top by `cn()` — tailwind-merge, not a size prop, is what actually resolves the visible height/padding/rounding (e.g. the xs overlay's `px-*` beats `Button`'s own default `px-*`, but any `Button` default token the overlay never mentions survives untouched). This port reproduces that exactly: `size` is used only for the switch and the `data-size` stamp, never passed to `Button`'s own `size` param; `variant` defaults to `"ghost"` (`Button`'s own zero-value default is `"default"`/primary) and IS forwarded to `Button`'s own `variant` param, matching shadcn's `variant = "ghost"` passthrough. `data-size` is set as an explicit non-parameter attribute on the `<Button>` call so it overrides `Button`'s own internal `data-size={size}` stamp (which would otherwise read `Button`'s own, deliberately-unset, `size` param) — same competing-defaults override mechanism as `ItemSeparator`/`ButtonGroupSeparator` overriding `Separator`'s `data-slot`. UPDATE (2026-07-24, nova density retarget): the literal token strings this entry used to cite are stale — `Button`'s own default size is now `h-8 gap-1.5 px-2.5 has-[>svg]:px-2` (was `h-9 px-4 py-2 has-[>svg]:px-3`) and `InputGroupButton`'s own default overlay is now `h-6 gap-1 rounded-[calc(var(--radius)-3px)] px-1.5 has-[>svg]:px-1.5` (was `rounded-[calc(var(--radius)-5px)]` `px-2`/`has-[>svg]:px-2`) — the merge MECHANISM described above is unchanged, only the numbers are; see `docs/superpowers/plans/2026-07-24-nova-density-map.md` `## input-group` for the full before/after.
- NOTE: `InputGroupText` carries no `data-slot` in shadcn's own source (unlike every other input-group part) — ported as-is, the same unmatched-`data-slot` call as `ButtonGroupText` (see `## button-group`).
- GAP: `InputGroupAddon`'s `onClick` handler (focuses the group's own `<input>` when a click lands anywhere in the addon except on a nested `<button>`) is client JS with no equivalent here — zero client JS for this component, per the Tier 1 plan's Tech Stack constraint. The addon renders and styles identically; it just isn't click-to-focus.
- GAP (narrow): shadcn's `InputGroupButton` accepts every `Button` prop except `size` (`Omit<ComponentProps<typeof Button>, "size">`), so it can render as a link via `Button`'s own mechanisms; this port's signature (`variant`, `size`, `children`, `attrs`) hard-codes `href=""` in the composed `Button` call, and `Button`'s `<a>`-vs-`<button>` tag choice is gated on that Go param — not on anything `attrs` can supply — so an `InputGroupButton` can never render as an `<a>`. `disabled` still works (flows through `attrs` into `Button`'s spread). Add an `href` param if a real use appears; none of shadcn's own input-group demos exercise it.

## item
- Straight port; no dropped tokens on `ItemGroup`/`ItemContent`/`ItemTitle`/`ItemDescription`/`ItemActions`/`ItemHeader`/`ItemFooter`. shadcn's own `item.tsx` has no Radix primitive underneath its own parts either (`Separator` is composed, not a Radix wrapper) — the same "package-namespaced compound parts" shape as card/breadcrumb/empty.
- WIN: `Item`'s `itemVariants` cva map (`variant`: default/outline/muted, `size`: default/sm) and `ItemMedia`'s `itemMediaVariants` (`variant`: default/icon/image) all pick between static class blocks by the JS-resolved prop values — no `data-[variant=...]`/`data-[size=...]` selectors in `registry/new-york-v4/ui/item.tsx` to preserve — so all three port as `switch`es inside `class={}`, the same idiom as badge/button-group/empty. `Item`'s two switches (variant, size) are inlined directly rather than extracted into shared helper functions the way `button.gsx`'s `variantClass`/`sizeClass` are — no sibling component reuses this pair the way `pagination.gsx` reuses button's.
- GAP (narrow): `Item`'s `asChild` tag-swapping (no dynamic tag) is dropped — always renders a `<div>`. Same narrow gap as button's own `asChild` (see `## button`); behavior-attachment uses of `asChild` are covered by the data-attribute mechanism (see dialog).
- WIN: `ItemSeparator` composes `ui.Separator` directly (flat package, no re-implementation) — the `item` → `separator` dependency `internal/registry` derives and `registry_test.go` pins.
- MECHANISM: shadcn's `ItemSeparator` types its props as `React.ComponentProps<typeof Separator>` and hardcodes `orientation="horizontal"` before spreading `{...props}` after it, so a caller-supplied `orientation` prop wins there. The port exposes `orientation` as a real Go param (`orientation |> default("horizontal")`), not left to `attrs` — `attrs` is untyped fallthrough onto `Separator`'s own rendered `<div>`, not a hook into the `orientation` argument compiled into the call to `Separator`, so only an explicit param can actually reproduce that override capability. Same competing-defaults mechanism as `ButtonGroupSeparator`'s own `orientation |> default("vertical")` (see `## button-group`), mirrored here with the opposite default.
- NOTE: `ItemDescription` renders a real `<p>`, matching shadcn's own source exactly — contrast with `EmptyDescription` (`## empty`), whose TypeScript prop type also says `"p"` but whose actual shadcn element is a `<div>`.

## kbd
- Straight port; no dropped tokens. Both `Kbd` and `KbdGroup` render onto real `<kbd>` elements, verbatim against `registry/new-york-v4/ui/kbd.tsx` — browsers freely nest `<kbd>` inside `<kbd>`, which is exactly how `KbdGroup` models a compound shortcut ("Ctrl Shift K") as a group of individual `Kbd`s.
- NOTE: shadcn's own `KbdGroup` types its props as `React.ComponentProps<"div">` but the component itself renders a `<kbd>` element, not a `<div>` — carried over verbatim, tag included, since the task brief calls for token-for-token parity and this is shadcn's actual shipped behavior, not a typo we're free to "fix."
- NOTE: `Kbd`'s `[[data-slot=tooltip-content]_&]:...` tokens are a real, exercisable selector in this port — nesting a `Kbd` inside a `ui.TooltipContent` (which does stamp `data-slot="tooltip-content"`) activates them exactly as shadcn intends.

## label
- Straight port of the rendered markup: shadcn wraps Radix's `LabelPrimitive.Root`, which is itself a plain `<label>`.
- GAP (narrow, accepted): Radix's `onMouseDown` handler, which calls `preventDefault()` on multi-click to stop text selection inside the label, is not ported (no client JS for this component). Low impact — the base class already carries `select-none`, which suppresses text selection via CSS regardless.

## pagination
- Straight port; no dropped tokens. shadcn's own `pagination.tsx` has no Radix primitive underneath either — every part (`Pagination`/`PaginationContent`/`PaginationItem`/`PaginationLink`/`PaginationPrevious`/`PaginationNext`/`PaginationEllipsis`) is already a plain styled `<nav>`/`<ul>`/`<li>`/`<a>`/`<span>`, the same shape as breadcrumb.
- WIN: `PaginationLink` composes button.gsx's package-private `base`/`variantClass`/`sizeClass` helpers directly (flat package, no re-implementation of `buttonVariants`) — the same call `Button` itself makes. This is the `pagination` → `button` dependency `internal/registry` derives (via `declIndex`, the same intra-package-edge shape as dialog → button) and `registry_test.go` pins.
- MECHANISM: `isActive`'s conditional `aria-current="page"` — present only when true, entirely absent (not merely empty) when false — uses gsx's `{ if cond { attr=value } }` conditional-attribute syntax, standing in for shadcn's `aria-current={isActive ? "page" : undefined}`.
- ADAPT: `size` defaults to `"icon"` (`PaginationLinkProps`' own `size = "icon"` default) rather than Button's `"default"` zero-value size — resolved in a small `{{ }}` prelude block (same mechanism as tabs.gsx's `state`/`tabindex` locals) before the class composition, since gsx's `|> default` pipe only applies at an expression's own top level and can't nest inside a `sizeClass(...)` call argument.
- WIN: `PaginationPrevious`/`PaginationNext` hardcode their own icon+label content exactly like shadcn's versions (no `children` param) — matching React's own behavior, where a component's literal JSX children always win over anything spread from `...props`, so a caller-supplied `children` prop would have been silently ignored in the original too.
- MECHANISM: `ChevronLeft`/`ChevronRight` (`PaginationPrevious`/`PaginationNext`) and `Ellipsis` (`PaginationEllipsis`, Lucide's `MoreHorizontal` — see breadcrumb's own `Ellipsis` entry) come from `ui/icon` — the `pagination` → `icon` dependency `internal/registry` derives and `registry_test.go` pins.

## progress
- ADAPT: shadcn's `Progress` wraps Radix's `ProgressPrimitive.Root`/`ProgressPrimitive.Indicator` pair (`registry/new-york-v4/ui/progress.tsx`), driven by no client state of its own (a plain controlled `value` prop) but still routed through Radix's component machinery. This port replaces both with two plain `<div>`s and zero client JS — `role="progressbar"` plus `aria-valuemin="0"`/`aria-valuemax="100"`/`aria-valuenow={value}` (stamped directly, no `data-state`/`data-value`/`data-max` — nothing in shadcn's own class strings selects on those, so there is nothing for them to drive here) replace what Radix's `Root` stamps internally. Both divs' class strings are carried token-for-token.
- WIN: `value` is a Go zero-value `float64` (0–100) — the zero value (0) already matches shadcn's own `value || 0` JS fallback, so an unset `Progress` renders identically to shadcn's unset case (indicator fully translated off-screen), no extra handling needed.
- GAP: `children` is not a param — shadcn's own `Progress` never renders a `children` prop either (its JSX literal-child list is exactly the one `ProgressPrimitive.Indicator`, spread `{...props}` notwithstanding: a literal child between tags always wins over a same-named prop spread before it), so there is nothing to port. Same "children omitted where the element can't have any" call as `Skeleton`/`Separator`.
- MECHANISM: the indicator's fill is Radix's `style={{ transform: translateX(-${100 - (value || 0)}%) }}` — ported as the identical `translateX` mechanism (not a `width` swap, which `transition-all` would animate differently), built from `strconv.FormatFloat(100-value, 'f', -1, 64)` since a `float64` can't itself concatenate into a string. The composed `style` value is wrapped in `gsx.RawCSS` (same precedent as `AspectRatio`'s `ratio` property, see `## aspect-ratio`) to opt the whole expression out of gw's CSS value filter, which blocklists `(`/`)` — punctuation `translateX(...)`'s function-call syntax requires, not injected data; the percentage is trusted, developer-computed layout intent, the same trust boundary aspect-ratio's `ratio` already extends.

## separator
- ADAPT: Radix's `decorative` prop (default `true`, flips `role="separator"` +
  `aria-orientation` when `false`) is not ported — `Separator` always renders
  `role="none"`, matching shadcn's default usage. No orientation param needed
  for a semantic separator variant; callers wanting a semantic (non-decorative)
  separator fall through `attrs` to set `role`/`aria-orientation` themselves.
- WIN: no variant switch — the single verbatim class string dispatches on
  `data-orientation` via Tailwind's `data-[orientation=...]` selectors, so
  `orientation` only needs to stamp the attribute.

## sheet
- WIN: `Sheet`/`SheetTrigger` compose `ui.Dialog` directly / mirror
  `AlertDialogTrigger`'s own doc-comment idiom (see `## alert-dialog`) —
  `Sheet` overrides the composed `data-slot` (`"dialog"`→`"sheet"`) via the
  same explicit-non-parameter-attribute mechanism as `AlertDialog` (`##
  alert-dialog` WIN), so `sheet → dialog` derives and the CLI vendors
  `ui/dialog.js` — trigger/content proximity wiring, Esc-to-close,
  `data-state` stamping, and the `getAnimations`-gated exit-animation wait
  are all reused unmodified. This is also the mechanism that gives sheet
  working exit animations at all (unlike, e.g., the popover family's
  ledgered-inert exit animations).
- ADAPT (the one part that does NOT compose): `SheetContent` renders its
  own native `<dialog data-slot="sheet-content" data-gsxui-dialog-content
  data-state="closed" data-side="...">` rather than composing
  `ui.DialogContent` the way `AlertDialogContent` composes it.
  `AlertDialogContent`'s merge-through-`class` trick worked only because its
  content recipe is nearly identical to `DialogContent`'s own (`##
  alert-dialog` ADAPT) — `DialogContent`'s centered-card recipe
  (`top-[50%]`/`left-[50%]`/`translate-*`/`max-w-lg`/`grid`) and `Sheet`'s
  side-anchored recipe (`inset-y-0`-or-`inset-x-0`/`h-full`-or-`h-auto`/
  `w-3/4`/`flex`) target the same CSS properties with materially different
  values on every one of them — there is no single caller class string that
  tailwind-merge could resolve into "dialog centered card, except also a
  side drawer." A second, independent `<dialog>` element is simpler and
  matches the design brief's own call ("the centering/grid class fights
  aren't worth tailwind-merge roulette").
- ADAPT: the base class's unscoped `flex` (from
  `registry/new-york-v4/ui/sheet.tsx`'s `SheetContent`) becomes `open:flex`
  — the identical display-gating fix as `DialogContent`'s own `open:grid`
  (`## dialog` ADAPT): content stays in the DOM while the native `<dialog>`
  is closed (no Radix unmount), so an ungated display utility would defeat
  the UA's closed-dialog `display:none`.
- ADAPT: `text-foreground` is added, the same fix and the same reason as
  `DialogContent`'s own (`## dialog` ADAPT) — native `<dialog>` gets UA
  `color: CanvasText` and does not inherit the themed body color the way a
  plain Radix `<div>` would.
- ADAPT (three-part fix, all found only by rendering `site/examples/sheet`
  in-browser per Task 3's render-and-verify requirement — none of these
  three UA-default fights are visible from the class string alone or from
  `go test`'s render-string pins, since HTML/CSS's own cascade/layout
  algorithm is what's at issue, not the generated markup): a modal
  `<dialog>`'s Chrome UA style is `position:fixed; inset:0; width:fit-
  content; height:fit-content; margin:auto; max-width/max-height:
  calc(100% - 6px - 2em)`. `sheet.tsx`'s ported side classes were written
  for a plain Radix `<div>`, which has none of these UA defaults to fight,
  and each defeats a different part of them:
  1. **Wrong edge (margin):** `sheet.tsx`'s side classes set only THREE of
     the four inset sides explicitly (e.g. the right variant's `inset-y-0
     right-0` sets top/bottom/right, leaving `left` unset) plus an explicit
     `width` (`w-3/4`). Reproduced live: with the UA's `margin:auto` and
     `left:0` (from its base, un-overridden `inset:0`) both left standing,
     the drawer first rendered dead-CENTER instead of flush right —
     `left`/`width`/`right` all specified with `margin` still `auto` hits
     CSS2.1 §10.3.7's "both margins auto → center" branch. Adding `m-0` to
     the base string (unconditionally zeroing margin, side-agnostic since
     every variant leaves exactly one edge unset the same way) forecloses
     that branch — but with `m-0` alone, the drawer then rendered flush
     against the WRONG edge instead (a right-side sheet opened flush
     LEFT): `left`(UA-0)/`width`(author)/`right`(author) all concretely
     specified and margins pinned to `0` is CSS2.1's genuinely
     over-constrained case, whose spec text says `left` should be the
     value ignored/recomputed for `direction:ltr` (i.e. `right`+`width`
     should win) — empirically, in Chrome, `left` won instead. The actual
     fix is each side switch case adding the ONE inset utility on the
     OPPOSITE edge from its anchor, set to `auto` (`left-auto` for right,
     `right-auto` for left, `bottom-auto` for top, `top-auto` for bottom):
     once that edge is genuinely `auto` (not merely absent-and-UA-
     defaulted-to-`0`), the box falls through to the unambiguous "one
     offset auto, the other two given" case and anchors correctly. `m-0`
     alone was necessary but not sufficient.
  2. **Height short by ~38px (max-height):** Chrome's UA `max-height:
     calc(100% - 6px - 2em)` applies to `:modal` dialogs too (confirmed via
     `getComputedStyle` in the live session) and clamps LEFT/RIGHT's
     `h-full` (an explicit `height:100%`, not `auto`) short of the full
     viewport by that same ~38px — measured directly (`getBoundingClientRect
     ().height` 38px less than `window.innerHeight`) before being
     neutralized by adding `max-h-none` to the base string (applies to all
     four sides; harmless for TOP/BOTTOM's small `h-auto` box, load-bearing
     for LEFT/RIGHT's `h-full`).
  3. **Width shrink-to-content instead of full-bleed (TOP/BOTTOM only):**
     `sheet.tsx`'s TOP/BOTTOM variants rely on bare `inset-x-0` with `width`
     left `auto` to stretch edge-to-edge — the plain-`<div>` default. On the
     native `<dialog>`, that did NOT stretch: measured rendered width came
     out shrink-to-fit (content width, e.g. ~447px on an ~1888px-wide test
     viewport) instead of full-bleed — a real Chrome behavior divergence
     from the classical CSS2.1 stretch algorithm this port works around
     rather than fully explains (setting `max-width:none` via inline style
     had zero effect on this specific symptom, ruling out `max-width` as
     ITS cause — it is a separate `auto`-width-resolution quirk). TOP/BOTTOM
     add `w-full` (not in `sheet.tsx`'s own source) to sidestep the `auto`-
     resolution ambiguity entirely — but `width:100%` is then itself
     subject to the SAME UA `max-width` (2) neutralizes on the block axis
     only, so TOP/BOTTOM also add their own explicit `max-w-none` (LEFT/
     RIGHT don't need one: their own `sm:max-w-sm` is already the sole
     author-origin `max-width` declaration and always beats the UA's,
     regardless of which is numerically smaller — cascade origin, not
     specificity or value, decides that tie).
  All four sides confirmed flush-to-edge, full-size, with working slide-in/
  slide-out enter and exit animations, in a real (focused, foregrounded)
  browser tab — see this task's report for the debugging trail, including a
  red herring (an unfocused/backgrounded MCP automation tab pauses CSS
  animations entirely, `document.visibilityState === "hidden"`, which
  briefly looked like a stuck-mid-animation position bug and was not one).
- ADAPT: `backdrop:bg-black/50` is folded onto `SheetContent`'s own class
  string — the same native-`::backdrop`-replaces-`Overlay`/`Portal`
  substitution as `DialogContent`'s own (`## dialog` WIN/GAP): `sheet.tsx`'s
  separate `SheetOverlay` (`fixed inset-0 z-50 bg-black/50` plus its own
  `data-[state=...]:animate-in/out fade-in/out-0` pair) is not ported at
  all — there is no second element for it to be. The backdrop's fade timing
  is not reproduced (a plain, unanimated translucent backdrop, exactly
  `DialogContent`'s own existing behavior) — not a new gap this task
  introduces, the same one `## dialog` already accepts.
- NOTE (signature judgment call): `SheetContent` carries a `hideCloseButton
  bool` param mirroring `DialogContent`'s own convention (zero value keeps
  shadcn's `showCloseButton` default of `true`) even though the task
  brief's own parts list only calls out `SheetContent(side)`.
  `showCloseButton` is a real, pre-existing `sheet.tsx` prop (not something
  newly grown since the brief was written, the category the brief's
  scope-out-via-GAP instruction actually targets), and `DialogContent`
  already established the exact bool-inversion pattern for this prop on the
  sibling component the brief explicitly points at for the close-button
  rendering ("compare dialog.gsx's injected button") — porting it keeps
  `SheetContent` at parity with `DialogContent`'s own capability. Flagged
  here as a judgment call, not a silent addition.
- MECHANISM: `side |> default("right")` both stamps `data-side` on the
  element (consistency with every other data-variant stamp in this
  codebase, plus availability to any downstream selector/JS) and selects
  one of `sheet.tsx`'s four static class blocks (right/left/top/bottom — no
  `data-[side=...]` selectors in the source to preserve, so this is a
  `switch` inside `class={}`, the same idiom as `item.gsx`'s variant/size
  pair, not a CSS-side data-attribute selector). `ui/dialog.js` itself never
  reads `data-side` — the slide direction is fully determined server-side
  by which static class block was selected, not client-side by the stamp.
- MECHANISM (close button structure, not `sheet.tsx`'s own): the injected X
  follows `DialogContent`'s own established structure (inline `<svg>` +
  `aria-label="Close"` on the `<button>`) rather than `sheet.tsx`'s actual
  source shape (`<XIcon className="size-4" />` + a separate `<span
  className="sr-only">Close</span>` for the accessible name) — for
  consistency across the dialog family (`Dialog`, and now `Sheet`), per the
  task brief's own framing ("compare dialog.gsx's injected button"). The
  CLASS STRING on the button and its icon IS `sheet.tsx`'s own, carried
  token-for-token (`data-[state=open]:bg-secondary`, not `DialogContent`'s
  `data-[state=open]:bg-accent data-[state=open]:text-muted-foreground`)
  — only the accessible-name/icon-embedding mechanics are shared with
  `DialogContent` rather than reproduced from `sheet.tsx`'s own JSX shape.
  Because `sheet.tsx`'s own close-button class string has no
  `[&_svg:not([class*='size-'])]:size-4` selector (unlike `DialogContent`'s
  own, which relies on exactly that selector to size an unclassed raw
  `<svg>`), the inlined `<svg>` here carries an explicit `class="size-4"`
  directly — matching what `sheet.tsx`'s own `<XIcon className="size-4" />`
  actually does, token-for-token, even though the icon-*embedding*
  mechanism (inline svg vs. a component) differs.
- GAP: `SheetOverlay`/`SheetPortal` are not ported — same reasoning as
  `## dialog`'s own Portal/Overlay GAP: the native `<dialog>` top layer
  replaces `Portal`, and `::backdrop` (folded onto `SheetContent`, see the
  ADAPT above) replaces `Overlay`. Nothing in this port ever needs a portal
  or a separate overlay element to target.
- GAP (narrow): `asChild` is not ported anywhere in this component — same
  narrow gap as `## dialog`'s own MECHANISM entry; the
  `data-gsxui-dialog-trigger`/`data-gsxui-dialog-close` attribute idiom is
  the replacement throughout.
- Registry: `sheet.gsx` imports nothing from `ui/icon`; `Sheet` calls
  `ui.Dialog` (flat-package intra-package edge, same shape as
  `## alert-dialog`'s own `dialog` dep) — `SheetTrigger`/`SheetContent`'s
  injected close button/`SheetClose` all render their own `<button>`
  rather than composing `ui.Button`, so `dialog` is the ONLY dep —
  `registry.Deps("sheet") == ["dialog"]`, pinned in
  `internal/registry/registry_test.go`. `HasJS("sheet")` is `false` — like
  `alert-dialog`, it has no behavior module of its own, only `ui/dialog.js`,
  pulled in transitively.

## skeleton
- Straight port; no divergences.

## spinner
- MECHANISM: shadcn's `Spinner` renders lucide-react's `Loader2Icon` directly (`registry/new-york-v4/ui/spinner.tsx`). `ui/icon`'s generated set carries the same Lucide glyph as `icon.LoaderCircle` (Lucide icon name `loader-circle` — `Loader2Icon` is lucide-react's own alias for it, same underlying SVG path data), so this port is one `icon.LoaderCircle` call — `spinner → icon`, the dependency `internal/registry` derives and `registry_test.go` pins — instead of a bespoke `<svg>`.
- MECHANISM: `ui/icon/icon.gsx`'s shared `svgIcon` wrapper defaults every icon to `aria-hidden="true"` unless the caller supplies its own `aria-hidden` — correct for the purely decorative chevrons/glyphs every other `icon.*` call in this package renders, but wrong for `Spinner`: shadcn's version is deliberately NOT hidden from assistive tech (`role="status"` and `aria-label="Loading"` announce it as a live loading indicator). `Spinner` supplies `aria-hidden="false"` explicitly alongside `role`/`aria-label`, all three literal attrs authored before `{ attrs... }` — the same overridable-default idiom as button's `type="button"` — so a caller can still override any of the three.

## textarea
- ADAPT: shadcn's Textarea takes its content via a `value` prop forwarded through `...props` onto React's controlled `<textarea value={...}>`. Native HTML `<textarea>` has no `value` attribute — its initial content is a text child. Ported as a `value string` param rendered as the (escaped) text child instead: `textarea.Textarea("initial text", nil)`.

## native-select
- RENAME (2026-07-24): `Select`/`SelectOption`/`SelectGroup` (file `ui/select.gsx`) → `NativeSelect`/`NativeSelectOption`/`NativeSelectGroup` (file `ui/native-select.gsx`); registry key `select` → `native-select`. `data-slot="select"`/`"select-trigger"` → `data-slot="native-select-wrapper"`/`"native-select"`, aligned with shadcn's own `registry/new-york-v4/ui/native-select.tsx` slot names (that file's `data-slot="native-select-wrapper"` div and `data-slot="native-select"` `<select>` map 1:1 onto this port's wrapper/select pair). Mirrors shadcn's own permanently-coexisting `select.tsx`/`native-select.tsx` split — frees the `Select`/`select` identifier space for the Tier 3 custom listbox (`docs/component-roadmap.md`). Pure rename: no class-string changes.
- ADAPT (native-first, prominent): shadcn's Select is Radix's full custom listbox — `SelectPrimitive.Root`/`Trigger`/`Portal`/`Content`/`Viewport`/`Item`/`ItemIndicator`/scroll buttons, driven entirely by client JS. v1 ships a styled **native `<select>`** instead: form-native, mobile-superior (the OS's own picker UI on touch devices, which Radix's custom listbox cannot delegate to), zero JS. shadcn's actual custom-listbox look (floating panel, check-mark indicator, keyboard-typeahead-driven Item highlighting) is **out of scope for v1 and tracked as post-v1 backlog** — what ships here is the shadcn *visual skin* (SelectTrigger's classes) applied to the real control, not a listbox reimplementation. `SelectContent`/`SelectItem`/`SelectValue`/`SelectLabel`/`SelectSeparator`/scroll-buttons are not ported; `NativeSelectGroup`'s label collapses onto native `<optgroup label=...>` (see MECHANISM below).
- ADAPT: SelectTrigger's class string is carried over token-for-token from `registry/new-york-v4/ui/select.tsx` except the following dropped, each dead or meaningless on a native `<select>`:
  - `data-[placeholder]:text-muted-foreground` — keys off Radix `SelectValue`'s `data-placeholder` state, which nothing in this port ever stamps (there is no controlled-value component standing in for the native control's own placeholder-less first-option display).
  - `data-[size=default]:h-9` / `data-[size=sm]:h-8` — Trigger's `size` prop (default `"default"`) is not ported (no size variant asked for in this task); replaced by an unconditional height — the default variant's own value — so the default visual height is preserved instead of shipping dead `data-[size=*]` selectors that nothing ever stamps. Same "bake in the default, drop the dead variant selector" call as avatar's dropped `data-[size=lg]/[size=sm]` pair (kept `size-8` unconditional). UPDATE (2026-07-24, nova density retarget): that baked-in default is now `h-8` (was `h-9`, the pre-nova default-variant value) — `## nova density` roadmap-lists the still-unported `size=sm` axis as a follow-up.
  - `*:data-[slot=select-value]:line-clamp-1 *:data-[slot=select-value]:flex *:data-[slot=select-value]:items-center *:data-[slot=select-value]:gap-2` — target a `data-slot="select-value"` child (Radix's `SelectValue`) that does not exist in this port; a native `<select>`'s displayed value is UA-rendered, not a child element these selectors could ever match.
  - `[&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg:not([class*='text-'])]:text-muted-foreground` — descendant-`svg` selectors; the chevron is a *sibling* of `<select>` in this port (a native `<select>` can only contain `<option>`/`<optgroup>`, so an `<svg>` could never legally be its descendant either), styled instead via its own literal class list.
  - Kept and appended: `appearance-none pr-8` — `appearance-none` suppresses the UA select arrow so the custom chevron is the only one shown (same mechanical necessity as checkbox/switch's `appearance-none`, not present on the Radix version since it renders no native `<select>` at all); `pr-8` reserves the gap the chevron sits in.
- MECHANISM: the chevron renders via `<icon.ChevronDown class="pointer-events-none absolute right-3 top-1/2 size-4 -translate-y-1/2 opacity-50"/>` as a sibling of `<select>` inside the `data-slot="native-select-wrapper" class="relative w-fit"` wrapper div, absolutely positioned over the trigger — this is the `native-select` → `icon` dependency `internal/registry` derives from the `.gsx` import and `registry_test.go` pins (`registry.Deps("native-select") == ["icon"]`).
- FINDING (2026-07-23, user report "style is wrong and very different than shadcn's"): the wrapper originally shipped as bare `class="relative"` (full-width block) with `w-fit` on the `<select>` — so the absolutely-anchored chevron pinned to the *wrapper's* far right edge, visually detached from the select box by the whole container width. shadcn never hits this because its chevron is a flex child *inside* the Trigger. Fix: the **wrapper** carries `w-fit` plus the caller's merged class (width intent like `w-full`/`w-[180px]` — where shadcn callers put it on Trigger — must size the box the chevron anchors to), and the `<select>` fills it with `w-full`; non-class attrs still land on the `<select>` (name/id/aria/disabled are form-control concerns). Same class-split idiom as AccordionContent. Render-string tests could never catch the detachment — found only by looking at the page next to shadcn's.
- MECHANISM: shadcn's separate `SelectGroup` (wrapper) + `SelectLabel` (child text with its own class string) collapse onto the one native `<optgroup label="...">` element — `<optgroup>` has no equivalent of an arbitrary styled label child, only a `label` attribute, so there is no class string of SelectLabel's to port; the attribute carries the text instead.
- WIN: `selected`/`disabled` on `NativeSelectOption` are HTML boolean attributes (`gsx.IsBooleanAttr` classifies both) — the Go `bool` params render bare/absent exactly like checkbox's `checked`, no `data-state` plumbing needed, unlike Radix's `SelectItem` (`data-[disabled]`, `ItemIndicator` mount-gating).

## select
The custom Radix listbox (distinct from `## native-select`, which ships the styled native `<select>`; shadcn keeps both permanently, and so does gsxui). Parts: `Select`/`SelectTrigger`/`SelectValue`/`SelectContent`/`SelectGroup`/`SelectLabel`/`SelectItem`/`SelectSeparator`. Built on `ui/select.js`, which reuses `ui/dropdown.js`'s popover machinery and layers a value model, typeahead, and a form bridge on top.
- WIN (machinery reuse): `SelectContent` rides the exact `DropdownMenuContent` mechanism — `popover="auto"` (top layer, light dismiss, free Esc), `closest("[data-gsxui-select]")` proximity wiring (no ids), the `data-state="open"` stamp BEFORE `showPopover()` (flash fix), the `pointerdown`-records-`wasOpen`/`click`-converges guard (the same outside-pointerdown-vs-light-dismiss race, since `SelectTrigger` is a real `<button>`), and `position: fixed` anchoring below the trigger (`data-side="bottom"` stamped statically). `select.js` additionally floors the content's `min-width` at the trigger's width (popper-equivalent to Radix's `--radix-select-trigger-width`).
- MECHANISM (discrete-transition reuse): the enter/exit class block (`opacity-0 scale-95 … transition-discrete duration-150 open:opacity-100 open:scale-100 starting:open:… data-[side=bottom]:starting:open:-translate-y-2`) ports byte-for-byte from `DropdownMenuContent`, replacing Radix's `data-[state=…]:animate-*`/`slide-in-from-*` keyframe block (same substitution ledgered under `## animations` and `## dropdown`). Exit is inert for the same reason as dropdown's (popover is `display:none` before `data-state="closed"` is stamped) — closing snaps.
- FINDING (`aria-selected` is NOT simply "is this the value"): traced from `@radix-ui/react-select@2.3.1` `dist/index.mjs`, an item's `aria-selected = isSelected && isFocused` — an item that IS the select's current value but is not ALSO the currently highlighted item reports `aria-selected="false"`. `data-state="checked"|"unchecked"` is the SEPARATE, simpler attribute that tracks the value alone (and is what the checkmark's CSS keys off). `select.js` mirrors this exactly: `focusItem()` recomputes `aria-selected` on every focus change (true only when the newly-focused item is also the checked value; every other item forced false), while the value model touches only `data-state`. All items server-render `aria-selected="false"` (nothing is focused at first paint).
- MECHANISM (focus model — real DOM focus, NOT `aria-activedescendant`): items are permanently `tabindex="-1"`; `select.js` moves REAL DOM `.focus()` among them (`focus:bg-accent` gives the highlight), one at a time — `dropdown.js`'s own roving-focus idiom, which the trace confirms is ALSO exactly what Radix's Select source does (`event.currentTarget.focus({preventScroll:true})` on hover, `focusFirst(candidateNodes)` on arrow keys). Deliberately NOT `command.js`'s `aria-activedescendant` model — Select's trigger is a `<button>`, there is no input to pin focus to and no visible-typeahead-box UX to protect. Hover-is-focus is gated on `pointerType==="mouse"` (touch/pen don't focus-on-hover); leaving the listbox parks focus on the content (`tabindex="-1"`, so arrow keys keep working) and clears every `aria-selected`.
- MECHANISM (typeahead — prefix + cycle, NOT command's fuzzy scorer): a bespoke ~20-line helper — a keystroke buffer reset after 1000ms idle, plain `startsWith` prefix matching, candidate list wrapped to start at the current item, and a same-character-repeat rule (a buffer of one repeated char, e.g. "bbb", normalizes to a single char and excludes the current item so repeated presses of the same key step through matches one at a time). Deliberately NOT `command.js`'s `commandScore` fuzzy engine — a materially different algorithm. Runs on BOTH the open listbox (moves focus) AND the CLOSED trigger (selects the value in place, like a native `<select>`). Space is swallowed as a search character only while a search is already in progress (`buffer !== ""`); a Space with an empty buffer is an OPEN request on the closed trigger / a selection on the open listbox. `OPEN_KEYS = [" ", "Enter", "ArrowUp", "ArrowDown"]`; on open, focus goes to the checked item if any, else the first enabled one (the map specifies no ArrowUp→last special case). `Tab` is `preventDefault`'d inside the open listbox (never leaves it, matching a native open `<select>`); `contextmenu` inside the content is suppressed. Arrow/Home/End move among enabled items only, clamped (no wrap).
- MECHANISM + GAP (form bridge — a hidden native `<select>`, not an `<input type="hidden">`): when `name != ""`, `Select` server-renders a real `<select aria-hidden="true" tabindex="-1" class="sr-only" name required disabled form>` sibling so native submission / `FormData` / a `<label>`'s click-through / autofill all see an ordinary working `<select>` (Radix's own `SelectBubbleInput` approach — a second fully-optioned select, not a hidden input). `select.js` keeps its `.value` synced on every selection and dispatches a bubbling `change` event (plain-DOM assignment, no property-descriptor bypass — there is no virtual-DOM layer to fight). **GAP (no-JS form submission carries no options)**: gsx has no context to collect the `SelectItem` values into the bridge at render time, so the server renders it with only a synthetic `<option value="">`; `select.js` populates one `<option>` per DOM item at init (module load, before any interaction — the bridge is complete by first user action). A form submitted with JS disabled therefore carries only the empty value. Radix's own bridge is SSR-populated via React context this port has no equivalent for — stated honestly, ledgered here.
- GAP (scroll up/down buttons dropped): `SelectScrollUpButton`/`DownButton` are not ported. Radix needs them mainly for its default `position="item-aligned"` mode (aligns the selected item under the trigger, can clip at the viewport edge); gsxui's anchoring is structurally always "popper-equivalent" (fixed-below-trigger, the only mode `dropdown.js` implements), so the content's own `overflow-y-auto` (`max-h-96` cap) just scrolls natively — exactly `DropdownMenuContent`'s own precedent. The `scrollable` site example (27 items / 5 groups) exercises this under real long-content conditions.
- MECHANISM (groups/labels are a real structural port, not a collapse): unlike `ui.NativeSelect`'s `SelectGroup` (which collapses onto native `<optgroup label=…>` because an optgroup can't hold an arbitrary styled child), the custom listbox CAN, so `SelectGroup` → `<div role="group">` and `SelectLabel` → its own styled `<div>` both port as real, separate parts. `select.js` wires each group's `aria-labelledby` to its label's generated id at init (no context to pass ids at render time).
- MECHANISM (checkmark visibility — CSS gating stands in for Radix's `ItemIndicator` mount-gating): Radix mounts the `CheckIcon` only when the item is selected; gsxui server-renders it always and gates visibility in CSS — the item carries `group/select-item`, and the indicator span is `hidden group-data-[state=checked]/select-item:flex`, so the check shows only when `data-state="checked"`. Works with JS off (a server-checked item shows its check). `size-3.5 → size-4` is the one place nova's icon size GROWS.
- ADAPT (nova metrics): `SelectTrigger` — `gap-2→gap-1.5`, `rounded-md→rounded-lg`, `px-3→pr-2 pl-2.5` (directional split), `shadow-xs` DROPPED (matches native-select's own approved drop), `data-[size=default]:h-9→h-8` / `data-[size=sm]:h-8→h-7` plus the new `data-[size=sm]:rounded-[min(var(--radius-md),10px)]`, and the value-slot `gap-2→gap-1.5`; `focus-visible:ring-[3px]` kept verbatim (this codebase's universal focus-ring syntax, not nova's `ring-3`). `SelectContent` — `min-w-[8rem]→min-w-36`, `rounded-md→rounded-lg`, `border` KEPT (no `border→ring` swap, house convention), `duration-150` KEPT (popover family standard, not nova's one-off `duration-100`). `SelectLabel` — `px-2→px-1.5`, `py-1.5→py-1`. `SelectItem` — `gap-2→gap-1.5`, `rounded-sm→rounded-md`, `py-1.5→py-1`, `pl-2→pl-1.5`, `pr-8` unchanged (reserves the indicator slot). `SelectSeparator` — byte-identical.
- ADAPT (`data-disabled="true"` vs Radix's bare `data-disabled=""`): a disabled `SelectItem` stamps `data-disabled="true"` (gsx's conditional-attribute form) plus `aria-disabled="true"`; Radix stamps `data-disabled=""`. The `data-[disabled]` presence selector matches either, and `select.js`'s `"disabled" in item.dataset` check works the same — a value difference with no behavioral effect. Enabled items carry neither attribute (correct `data-[disabled]` presence semantics, unlike an always-rendered `data-disabled="false"`).
- ADAPT (pointer-type-aware activation): mouse selects on `pointerup` (Radix's snappier path; a `suppressClick` flag + `setTimeout(…, 0)` reset stops the trailing `click` from double-firing and never leaves the flag stale); touch/pen fall through to `click`; keyboard `Enter`/`Space` select via the content keydown handler. `gsxui:select` (with `{ value }`) is emitted on the root for every selection path.
- Registry: `select.gsx` imports `ui/icon` (SelectTrigger's `ChevronDown`, SelectItem's `Check`) — `registry.Deps("select") == ["icon"]`, the second `select → icon` edge the controls map called out (alongside `native-select`'s pre-existing one). `HasJS("select")` is true (`ui/select.js`); `native-select` remains JS-free.

## table
- NOTE: `Table` renders a scroll-container `<div data-slot="table-container">` wrapping `<table data-slot="table">`, matching shadcn's structure exactly. Fallthrough `attrs` land on the `<table>` element (where shadcn's `{...props}` lands), not the container div — the container is purely structural scroll-wrapping and has no props of its own in the source either.

## tabs
- GAP: Radix's client context (`TabsPrimitive.Root` tracking `value` and broadcasting it to every `Trigger`/`Content` descendant) has no gsx equivalent — parts are plain siblings, no shared tree state. `TabsTrigger`/`TabsContent` each take an explicit `selected bool`; the caller (which already has both its own value and the root's in scope when building the tree) resolves the `value == root's value` comparison and passes the result in. The zero value (`false`) renders inactive, so a caller who forgets the comparison never gets an accidentally-active tab. Same shape as the switch/checkbox/radio explicit-state ADAPTs, applied to a value the server genuinely cannot look up without context.
- WIN: server-rendered initial selection — `aria-selected`, `data-state`, and roving `tabindex` (`0` for the selected trigger, `-1` for the rest, matching the WAI-ARIA APG tablist pattern) are all derived from `selected` at render time, so first paint already matches the state `tabs.js`'s `activate()` would produce after a click — no visible "tabindex flip"/aria flip once the behavior module loads. `TabsContent` stamps the native `hidden` boolean attribute from `!selected`, so an inactive panel is genuinely not rendered/not in the accessibility tree even with JS disabled or not yet loaded.
- ADAPT: shadcn's `orientation` (horizontal/vertical) prop and `TabsList`'s `variant` (default/line) `cva()` axis are both out of task scope — no param for either, only horizontal/default ship. Every class token whose sole purpose was to key off one of those two Radix-only-reachable states is dropped as dead weight (same "drop the selector, don't ship dead CSS" call as avatar's size prop and dialog's close-button `data-[state=open]:` pair, not a new precedent): `Tabs`' root drops `group/tabs` and unwraps `data-[orientation=horizontal]:flex-col` to an unconditional `flex-col`; `TabsList` drops `group/tabs-list`, both `group-data-[orientation=vertical]/tabs:` tokens, and `data-[variant=line]:rounded-none`, folding cva's default-variant `bg-muted` in unconditionally; `TabsTrigger` drops the vertical-orientation `w-full`/`justify-start` pair, drops the entire `group-data-[variant=line]/tabs-list:` family (background/border overrides for a variant that no longer exists), unwraps `group-data-[variant=default]/tabs-list:data-[state=active]:shadow-sm` to an unconditional `data-[state=active]:shadow-sm`, and drops the `after:` line-indicator pseudo-element entirely (its only visibility rule, `group-data-[variant=line]/tabs-list:data-[state=active]:after:opacity-100`, can never fire under the one variant shipped, so the positioned-but-invisible `after:` box is pure dead CSS). Neither `data-orientation` nor `data-variant` is stamped anywhere — nothing reads them. `TabsContent`'s `flex-1 outline-none` has no orientation/variant dependency and ports unchanged.
- MECHANISM: `ui/tabs/tabs.js` (click + roving ArrowLeft/ArrowRight, both ordinary bubbling events — no `{ capture: true }` needed) re-stamps `data-state`/`aria-selected`/`tabindex` on every trigger and `data-state`/`hidden` on every panel by scoping to `closest("[data-gsxui-tabs]")`, then emits `gsxui:change` on the root with `{ value }` — the same closest-root delegation idiom as dialog's `rootOf`, and the same `gsxui:*` CustomEvent-as-API idiom as dialog's `gsxui:open`/`gsxui:close`.

## accordion
- WIN: Radix's `AccordionPrimitive.Root`+`Item` client state machine (`type="single"`/`"multiple"`, `collapsible`, `value`/`onValueChange`, keyboard-nav Item coordination) is replaced by grouped native `<details name="…">`: the browser itself enforces "opening one member of a group closes the rest," zero JS, zero client state to keep in sync with the server-rendered markup. Verified interactively (dev/preview.html): opening Item 2 closes Item 1 with no script involved at all.
- GAP: same shape as tabs' `value` — grouping-by-name is a real browser mechanism, not client-side context, so nothing exists for a root to propagate down to its items. The caller passes the same `name` to `Accordion` and to every `AccordionItem` in the group explicitly. `Accordion`'s own `data-name` stamp is a readability/debugging aid only — nothing reads it back (there is no `accordion.js`).
- WIN: `AccordionItem`'s `open bool` server-renders the initial expanded/collapsed state as the native `open` boolean attribute (zero value `false` = collapsed, matching shadcn's Radix default of nothing expanded until interacted with) — thereafter, opening/closing and the exclusive-group behavior are entirely native, no hydration step to reconcile.
- ADAPT: the Radix `Header`+`Trigger` pair is flattened into a single `<summary>` — a bare block-level `<summary>` already lays out its own clickable row, so the wrapping `<AccordionPrimitive.Header className="flex">` is dropped along with the now-meaningless `flex-1` on the Trigger (there is nothing left to flex within). `list-none` and `[&::-webkit-details-marker]:hidden` are added (not in shadcn's source, which has no native disclosure marker to suppress) to strip the browser-drawn triangle both marker mechanisms produce, leaving only the ported chevron visible.
- ADAPT: the rotate-on-open selector moves off the trigger and onto the chevron icon itself, because `<summary>` carries no `data-state` (Radix's trigger-side hook `[&[data-state=open]>svg]:rotate-180` has nothing to key off here) — native `<details>` truth is the `[open]` attribute on the ancestor `<details>`, not anything on `<summary>`. The chevron's own class carries `[[open]>summary_&]:rotate-180`, an arbitrary-variant selector reading "an ancestor with `[open]`, with a `summary` directly beneath it, with this element somewhere under that summary." Validated two ways in this sandbox (no Tailwind build available to compile it directly, flagged for a real-browser check downstream): (1) against Tailwind v4's documented ancestor-selector idiom for arbitrary variants (`[.dark_&]:…`-style "referencing the parent when nesting," substituting the attribute selector `[open]` plus a child combinator for the documented class-selector case — structurally the same shape, just one hop deeper); (2) against `merge.Merge` (`tailwind-merge-go`, this repo's configured `class_merger`), which round-trips the bracketed token unchanged and correctly buckets it as a `rotate` utility. Note the merge semantics precisely: a bare caller-supplied `rotate-90` does *not* replace this variant-scoped token — an unscoped utility and a `[[open]>summary_&]:`-scoped one aren't merge conflicts (the same non-collision as an unscoped utility next to a `hover:`-scoped one), so both would apply and the CSS cascade decides. A caller who wants to override the rotation must supply the same variant prefix, e.g. `[[open]>summary_&]:rotate-90`. No structurally-implausible fallback was needed; the brief's given selector is used as specified.
- ADAPT: shadcn's `ChevronDownIcon` class carries `pointer-events-none` and `translate-y-0.5` that this port drops. `pointer-events-none` existed so the icon couldn't steal a click from Radix's button-role Trigger; here the click target is the whole `<summary>` (native disclosure semantics, not an attached click handler), so there is nothing for the icon to steal. `translate-y-0.5` nudged the icon against the now-removed `Header` wrapper's row baseline; with the icon sitting directly in `<summary>`'s own flex row there is no longer a baseline mismatch to correct for.
- MECHANISM (CSS-only animation, shipped): shadcn's `data-[state=open]:animate-accordion-down`/`data-[state=closed]:animate-accordion-up` pair is dropped (nothing stamps `data-state` here) and replaced by pure CSS in gsxui.css keyed on the data-slot attributes: `[data-slot="accordion-item"]::details-content` — the pseudo-element wrapping everything but the summary — transitions its grid row `0fr -> 1fr` on `[open]`, which is a height-to-auto animation needing neither `interpolate-size` nor Radix's measured `--radix-accordion-content-height`, plus `content-visibility ... allow-discrete` so the collapsing content stays rendered until the close transition finishes (`min-height: 0` on the content div lets the 0fr row actually collapse). Browser-verified both directions via rAF height sampling (0->36px easing across ~100 frames at a stretched duration). Progressive enhancement: browsers without `::details-content` ignore the rules and toggle instantly, exactly the pre-animation behavior. The block lives in assets/gsxui.css AND web/site.css (the site copies, never imports — `TestAccordionAnimationCSSDriftPin` pins the two byte-identical).
- WIN: shadcn's `AccordionContent` splits its rest-props from its `className` by hand — Radix's outer `Content` element takes `{...props}` and a hardcoded `className="overflow-hidden text-sm"`, while the caller's `className` is redirected onto a separate inner `<div className={cn("pt-0 pb-4", className)}>` wrapper via `cn()`. gsx expresses the same routing declaratively instead of by manual prop-plumbing: `AccordionContent` spreads `attrs.Without("class")...` onto the outer div (rest props, class excluded) and merges `attrs.Class()` into the inner wrapper's `class={ "pt-0 pb-4", attrs.Class() }` — the queryable-bag pattern (`Without`/`Class` as bag accessors) makes the split a one-line read instead of a second component-local variable.
- MECHANISM: `AccordionItem`'s `border-b last:border-b-0` ports unchanged.
- Registry: `accordion.gsx` imports `ui/icon` for the chevron — `registry.Deps("accordion") == ["icon"]`, same dependency-derivation mechanism as Task 6's dropdown/select → icon.

## dropdown
- MECHANISM: hover highlight IS focus — shadcn item classes style `focus:` only because Radix moves roving focus onto hovered items; dropdown.js ports that (pointerover → item.focus(); leaving the menu parks focus on the content via tabindex="-1" so arrow keys keep working).
- ADAPT: Radix runtime CSS vars (`--radix-*-available-height`/`transform-origin`) don't exist without Radix — replaced with static `max-h-96` / `origin-top-left`.
- ADAPT: `dropdown.js` sets `content.style.inset = "auto"` before setting `left`/`top` — the UA popover default rule (`inset: 0`) otherwise leaves `right: 0; bottom: 0` active alongside the hand-set `left`/`top`, collapsing/mispositioning the content since only two of the four inset edges were ever overridden.
- WIN: Radix Portal/Content replaced by the native popover API (`popover="auto"` on `DropdownMenuContent`) — top layer, light dismiss, and Esc are all browser-native, same win as dialog's `<dialog>`. Trigger↔content wiring is `closest("[data-gsxui-dropdown]")` proximity in JS, no ids, same MECHANISM as dialog.
- NOTE: positioning is a hand-rolled `position: fixed` anchor to the trigger's `getBoundingClientRect()` in `dropdown.js` (open below, left-aligned, 4px gap), not Radix's Floating-UI collision-aware placement. CSS anchor positioning (`anchor()`/`position-anchor`) — the eventual native replacement — is not yet Baseline (Chrome-only as of this writing); JS positioning is a stopgap until it is. `data-side="bottom"` is stamped statically (dropdown.js always anchors below), so the ported `data-[side=bottom]:slide-in-from-top-2` enter slide is live; the other three `data-[side=*]` selectors are dead until real placement logic lands.
- GAP: `DropdownMenuCheckboxItem`, `DropdownMenuRadioGroup`/`DropdownMenuRadioItem`, `DropdownMenuSub`/`DropdownMenuSubTrigger`/`DropdownMenuSubContent`, and `DropdownMenuGroup` are not ported — post-v1 backlog per the task brief. Only `DropdownMenu`, `DropdownMenuTrigger`, `DropdownMenuContent`, `DropdownMenuItem`, `DropdownMenuLabel`, `DropdownMenuSeparator`, `DropdownMenuShortcut` ship.
- ADAPT: the `inset` prop on `DropdownMenuItem`/`DropdownMenuLabel` is dropped along with it — `data-[inset]:pl-8` is removed from both classes as dead weight, the same "drop the selector, don't ship dead CSS" call as dialog's close-button ADAPT and avatar's size ADAPT.
- MECHANISM: `DropdownMenuItem` is a real menu item on `<div role="menuitem" tabindex="-1">` (Radix's own `Item` is also a non-native, ARIA-role div, so this isn't a native-first swap like checkbox/radio/switch — it's a straight structural port). `dropdown.js`'s arrow-key roving focus walks `[role="menuitem"]:not([aria-disabled])` within the content; click on an item emits `gsxui:select` then closes via `hidePopover()`.
- A11Y: `aria-haspopup="menu"` and the initial `aria-expanded="false"` are server-rendered on `DropdownMenuTrigger`; `aria-expanded` is synced by `dropdown.js` on the `toggle` event, same non-bubbling-event/capture pattern as dialog. On open, the first enabled menu item receives focus.
- NOTE: state + events ride `ToggleEvent` (`toggle`, capture-delegated) — same rationale as dialog's ledger NOTE (covers every open/close path including light dismiss and Esc, which a `close`-only or click-only wiring would miss).
- ADAPT (verified in-browser): exit animations are inert — the native popover is display:none before data-state="closed" is stamped, so data-[state=closed]:animate-out never runs; closing snaps. Accepted for v1; an animated-close strategy (beforetoggle/allow-discrete or dialog-style requestClose) is in the README backlog.
- MECHANISM (trigger click-close guard): with `popover="auto"`, a real pointerdown on the trigger is an *outside* pointerdown relative to the content (the trigger isn't inside it) — the popover API light-dismisses on that pointerdown before the subsequent click fires, so a bare `togglePopover()` in the click handler would see an already-closed popover and wrongly reopen it. `dropdown.js` guards this by recording, on `pointerdown`, whether the content was open at that instant (`trigger.dataset.gsxuiWasOpen`); the `click` handler reads and clears that flag and no-ops if it was `"true"` (light dismiss already closed it), otherwise falls through to an explicit `hidePopover()`/`showPopover()` pair (never a bare `toggle`) so the keyboard-activation close path (no intervening pointerdown, flag stays `"false"`) still works. Real-pointer verification of the light-dismiss race is unavailable in this session's dev harness (no real pointer device to reproduce the pointerdown/click ordering); the guard is state-preserving either way — with or without the race actually firing, the trigger always converges on toggling from the popover's actual current state rather than an assumed one.

## toggle
- GAP (form: `<input type="checkbox">` rejected): the same "real form control, zero JS" idiom `## checkbox`/`## switch` use was considered here too and rejected — `input` is a VOID element (no children, no closing tag), while `Toggle`'s entire visible surface IS its children (an icon, text, or both — every one of shadcn's own `toggle-*.tsx` examples renders a Lucide icon, sometimes plus a label, as `Toggle`'s children), rendered inside the pressable control itself, not a `<label>` sibling to a hidden control. A `<button>` is the only element shape that can both be the toggle and hold arbitrary child content, so this port needs a small behavior module (`ui/toggle.js`) rather than riding free on browser `:checked` state.
- ADAPT (explicit initial state): `pressed bool` stands in for Radix's `TogglePrimitive.Root` uncontrolled `pressed`/`defaultPressed` — the same shape as `## tabs`' `selected`/`## collapsible`'s `open`. Zero value (`false`) renders unpressed (`aria-pressed="false" data-state="off"`), matching Radix's own default; `pressed={true}` server-renders already-on with no click required.
- WIN: `toggleVariants`' cva map (`variant`: default/outline, `size`: default/sm/lg) is every case a static class block picked by the JS-resolved prop value — no `data-[variant=]:`/`data-[size=]:` selectors exist in `registry/new-york-v4/ui/toggle.tsx` to preserve, so both port as a `switch` inside `class={}`, the same idiom as `## button`'s `variantClass`/`sizeClass` pair (inlined here, like `## item`'s own pair, since only `Toggle` uses it). The one data-keyed selector upstream DOES have, `data-[state=on]:bg-accent data-[state=on]:text-accent-foreground`, carries verbatim in the shared base class — it lights up once `ui/toggle.js` flips `data-state` on click.
- ADAPT: `data-variant`/`data-size` are stamped via the house `|> default(...)` pattern (`## button-group`'s own ADAPT note names this convention explicitly) for consistency with every other cva-backed component in this codebase, even though shadcn's own `Toggle` stamps neither — Radix's `TogglePrimitive.Root` receives only `data-slot`; `toggleVariants` resolves straight to `className`, no data attrs of its own to port.
- MECHANISM (`ui/toggle.js`): one capture-delegated click listener on `[data-gsxui-toggle]` — reads the current `aria-pressed`, flips both `aria-pressed` and `data-state` to the opposite state, then emits a bubbling `gsxui:change` CustomEvent with `{ pressed }`, the same event-shape convention as `## tabs`' own `gsxui:change` (value-flip components emit `gsxui:change`; open/close components emit `gsxui:open`/`gsxui:close`, per `## dialog`). No non-bubbling event to intercept here (unlike dialog's `toggle`/`cancel`), so the listener needs no `{ capture: true }` — a plain `click` delegated the same way as `## tabs`' own trigger click, the simplest existing behavior module and this task's stated template.
- Registry: `toggle.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("toggle")` is empty (the site example's `icon.Bold`/`Italic`/`Underline` imports live in `site/examples/toggle/basic.gsx`, which `internal/registry` never scans — same shape as `## collapsible`'s own deps entry). `HasJS("toggle")` is `true` (this task's own `ui/toggle.js`) — the first Tier 2 component to gain a behavior module of its own rather than reusing another component's (`## alert-dialog`/`## sheet` both reuse `dialog.js` transitively; `toggle` is standalone).

## toggle-group
- FINDING (shadcn upstream, dead selector in `toggle-group.tsx` itself): `ToggleGroup`'s root class carries `data-[spacing=default]:data-[variant=outline]:shadow-xs`, but `data-spacing` is stamped as the literal prop value (`data-spacing={spacing}`, default `0`, a JS **number**) — never the string `"default"`. `data-[spacing=default]` can therefore never match anything the component itself ever renders; it is dead CSS in the shadcn source, not just a port-time casualty. (Circumstantial evidence for what was probably intended: the ITEM's own `data-[spacing=0]:shadow-none` suppresses the per-item outline shadow when items are joined into one pill — the root's shadow was very likely meant to substitute a single shadow around the whole pill in that same `spacing=0` case, i.e. the selector should read `data-[spacing=0]`, not `data-[spacing=default]`.) Unlike this codebase's usual "port dead weight, ledger it" convention (`## dialog`'s dropped `data-[state=open]:` pair is that precedent), this port DROPS the selector outright instead of porting it verbatim — nova's own `.cn-toggle-group` (`style-nova.css` L1371-1373) already drops it rather than fixing the typo, and this port follows nova's own precedent (a deliberate deviation from the house convention, called out here rather than silently mirroring nova).
- GAP (horizontal-only v1): `data-orientation="horizontal"` is stamped on both `ToggleGroup` and `ToggleGroupItem`, but only the horizontal corner-rounding selectors ship (`data-[orientation=horizontal]:data-[spacing=0]:first:rounded-l-lg`/`last:rounded-r-lg`). Vertical orientation is real new functionality nova adds on top of new-york-v4 (`group-data-vertical/toggle-group:data-[spacing=0]:first:rounded-t-lg` etc. — new-york-v4's own markup never varies rounding by orientation at all, so a naive vertical port would round the wrong corners), not just a metric bump — out of v1 scope, same "scope out the axis, ledger it" call as `## tabs`' dropped `orientation` prop. `toggle-group.js`'s arrow-key handling is written for horizontal only (ArrowLeft/Up and ArrowRight/Down both move focus, since `data-orientation` is always `"horizontal"` in this port); wiring vertical would mean gating that on the real attribute the same way the CSS above would need to.
- ADAPT (`groupType`, not `type` — group→item inheritance is explicit params): Radix's `ToggleGroupContext.Provider` broadcasts `type`/`variant`/`size`/`spacing` from `ToggleGroup` to every `ToggleGroupItem` automatically; gsx has no context, so `ToggleGroup(groupType, variant, size, spacing, …)` and `ToggleGroupItem(groupType, variant, size, spacing, pressed, value, …)` both take all four explicitly — the caller (which already has every value in scope when building the tree) repeats them at every item call site, the same "no context" shape as `## tabs`' `selected`. The parameter is spelled `groupType`, not `type` — `type` is a Go keyword and cannot be a component parameter name (confirmed: `component X(type string, …)` fails to parse with `expected ')', found 'type'`), a different flavor of Go-keyword workaround than `select`/`switch`'s own (those needed the workaround at the file/package-name level only; this one is at the parameter level). `groupType == "single"` selects `role="radiogroup"` (root) / `role="radio" aria-checked` (item); anything else, including the Go zero value `""`, renders `role="toolbar"` (root) / `aria-pressed`, no role override (item) — Radix throws if `type` is omitted at mount, which Go has no render-construction-time equivalent for, so this is a doc-comment API contract rather than a runtime check. Radix's group-level `disabled` OR-cascade (`context.disabled || props.disabled`) has the same no-context shape and is NOT ported as special logic: `disabled` flows through each element's own `attrs` bag independently (native `<button disabled>` on the item; an inert-but-present attribute on the root's `<div>`), so a caller disabling a whole group passes `disabled` to every item explicitly, same as the other four inherited values.
- ADAPT (JS-normalized roving tabindex, not SSR-computed): unlike `## tabs`' `TabsTrigger`, which server-renders its roving `tabindex` from a caller-supplied `selected bool` comparison, `ToggleGroupItem` stamps NO `tabindex` attribute at all — every item is a plain, fully-tabbable element until `ui/toggle-group.js` loads (a graceful no-JS fallback: no JS means every item is its own tab stop, not zero of them). At module-init time (same one-time `document.querySelectorAll` scan `## command`'s own initial `filter()` pass uses), and again on every click/arrow-key move, `toggle-group.js` collapses the group to exactly one `tabindex="0"` item — the pressed-and-enabled item if one exists, else the first non-disabled item — and sets every other item to `tabindex="-1"`. This priority is type-agnostic: it applies identically whether `groupType` is `"single"` or `"multiple"`, matching Radix's own traced `RovingFocusGroup` entry-focus chain (the map's behavior contract conditions "active item wins, else first" on nothing but disabled-ness — never on `type`), not a single-type-only rule. This was chosen over SSR-computed tabindex because, unlike `TabsTrigger`'s single `selected` comparison, roving-focus entry priority needs a scan across ALL sibling items (find the pressed one, or fall back to first-enabled) that the item component itself cannot do — it only ever sees its own params, not its siblings'.
- ADAPT (Shift+Tab exit has no dedicated handler): the map's traced Radix behavior intercepts Shift+Tab (`onItemShiftTab`) to force the group container's own `tabIndex` to `-1` for one blur/focus cycle, ensuring Shift+Tab exits the whole group in a single press. This port has no equivalent keydown interception — and doesn't need one: because roving tabindex already collapses the group to exactly one tab stop at all times post-init, a plain Shift+Tab from that one focused item already lands on whatever precedes the group in the page, with no other item in the group reachable via Tab to "catch" it. Radix's `onItemShiftTab` earns its keep inside their FocusScope architecture (which additionally makes the group's own container element a fallback tab stop in some states); a plain-DOM roving-tabindex port has no such fallback container tab stop to guard against, so the one-tab-stop invariant alone is sufficient.
- GAP (RTL arrow-key swap not ported): the map's traced behavior contract includes Radix's `getDirectionAwareKey`, which swaps `ArrowLeft`/`ArrowRight` when `dir === "rtl"`. `toggle-group.js` has no `dir`/RTL awareness at all — ArrowLeft is always "prev" and ArrowRight is always "next", regardless of the page's writing direction. No gsxui component has `dir`/RTL support today (no other ported component's own ledger entry mentions one either), so this isn't a toggle-group-specific gap so much as a codebase-wide one surfacing here first because toggle-group is the first component whose traced Radix source happens to branch on `dir`; revisit if RTL ever becomes a target.
- MECHANISM (single-type replace-on-activate): `ToggleGroupImplSingle`'s `onItemActivate` is `setValue` — clicking a new item simply replaces which one is checked, no "uncheck the others" loop is needed because there's a single value context to update. This port has no shared value state to update (no context), so `toggle-group.js` restates the same *effect* explicitly: on a `groupType="single"` item's activation, it flips that item's own `data-state`/`aria-checked`, then walks every other item under the same root and force-clears theirs — same visible outcome (exactly one `data-state="on"` at a time), different mechanism (an explicit sibling walk standing in for context's implicit fan-out). Clicking the already-pressed item toggles it off, porting Radix's own allow-empty-single-value default.
- MECHANISM (`ui/toggle-group.js`): reuses `## dropdown`'s `closest("[data-gsxui-toggle-group]")` proximity-wiring idiom (not its code) — real DOM focus moves between real `<button>` elements via `.focus()`, no `aria-activedescendant`, the same shape dropdown's roving arrow-key walk uses, except toggle-group's items are ALWAYS real tab stops (roving tabindex, not "only reachable once a popover is open"), so the walk also writes real `tabindex` 0/-1 toggling rather than dropdown's fixed `-1`/`.focus()`-only pattern. Key map: `ArrowLeft`/`ArrowUp` → prev, `ArrowRight`/`ArrowDown` → next, `Home`/`PageUp` → first, `End`/`PageDown` → last, looping at both ends — ported verbatim from the traced `getFocusIntent`. A held modifier key (`metaKey`/`ctrlKey`/`altKey`/`shiftKey`) suppresses the whole handler, porting Radix's own early-return guard. Click activation emits `gsxui:change` on the root with `{ value }` — a string (the newly-pressed item's `value`, or `""` if the click toggled it off) for `groupType="single"`, an array of every currently-pressed item's `value` for `groupType="multiple"` — the same `gsxui:change` event-name convention `## tabs`/`## toggle` both use for value-flip components, with a payload shape chosen to fit a group value rather than a single boolean/string.
- Registry: `toggle-group.gsx` has no `ui/icon` import (the site examples' `icon.Bold`/`Italic`/`Underline` imports live in `site/examples/togglegroup/*.gsx`, which `internal/registry` never scans — same shape as `## toggle`'s own deps entry) but `ToggleGroupItem` calls `toggle.gsx`'s package-private `toggleBase`/`toggleVariantClass`/`toggleSizeClass` directly (extracted from `Toggle`'s own inline `class={}` switch expressions into named helpers for this reuse, the same `## pagination` → `## button` `base`/`variantClass`/`sizeClass` shape) — `registry.Deps("toggle-group")` is `["toggle"]`, resolved via `declIndex` with no import to scan. `HasJS("toggle-group")` is `true` (this task's own `ui/toggle-group.js`, not a reuse of `toggle.js` despite the class dependency — the two components' interaction models, roving-tabindex-across-siblings vs. single-button click-flip, don't overlap enough to share one behavior module).

## tooltip
- ADAPT (2026-07-24): `TooltipPrimitive.Arrow` ports as a static child `<span data-slot="tooltip-arrow">` carrying shadcn's Arrow classes (size-2.5 rotate-45 rounded-[2px] bg-foreground) positioned `top-full left-1/2` with `-translate-x-1/2 -translate-y-[calc(50%+2px)]` — our tooltip is ALWAYS JS-anchored above the trigger, so the diamond always straddles the bubble's bottom-center; Radix's side-tracking arrow slot collapses to static CSS. (Originally dropped with a dangling "see ledger" pointer and no entry — the visual gap was user-reported against shadcn's docs side-by-side.)
- ADAPT: Radix runtime CSS vars (`--radix-*-transform-origin`) don't exist without Radix — replaced with static `origin-bottom` (the tooltip is positioned above and horizontally centered on the trigger via `translate(-50%, -100%)`, so its anchor point is its own bottom-center edge).
- ADAPT: `tooltip.js` sets `content.style.inset = "auto"` before setting `left`/`top`, same rationale as dropdown's ledger ADAPT — the UA popover default `inset: 0` otherwise leaves `right`/`bottom` active alongside the hand-set `left`/`top`.
- WIN: Radix Portal/Content replaced by the native popover API — `popover="manual"` on `TooltipContent` puts it in the top layer without the light-dismiss/Esc behavior `popover="auto"` would add (which would race hover-driven show/hide and dismiss on an unrelated outside click while the pointer is still over the trigger).
- NOTE: positioning is a hand-rolled `position: fixed` anchor centered above the trigger's `getBoundingClientRect()` in `tooltip.js`, not Radix's Floating-UI collision-aware placement — same stopgap-until-CSS-anchor-positioning rationale as dropdown's NOTE. `data-side="top"` is stamped statically (tooltip.js always anchors above), so the ported `data-[side=top]:slide-in-from-bottom-2` enter slide is live; the other three `data-[side=*]` selectors are dead until real placement logic lands.
- GAP: `TooltipProvider` (shared `delayDuration`/skip-delay-group machinery across multiple tooltips) is not ported — `tooltip.js` hard-codes a 300ms open delay per trigger, no cross-tooltip grouping. `TooltipPrimitive.Arrow` (the small pointing triangle) is also not ported — no `Arrow` part or class exists in `tooltip.gsx`.
- MECHANISM: `pointerover`/`pointerout`/`focusin`/`focusout` are used (not `mouseover`/`mouseout`/`focus`/`blur`) specifically because they bubble — dialog/dropdown's non-bubbling events (`toggle`) need `{ capture: true }`, but these delegate the ordinary way via `on()`'s bubble-phase default. `showPopover()`/`hidePopover()` are called directly rather than `togglePopover()` since show and hide are driven by two independent event pairs (pointer vs. focus), not one toggling control.
- A11Y: `TooltipContent` server-renders `role="tooltip"`; no `aria-describedby` is wired from trigger to content (Radix's Tooltip does this internally) — narrow GAP, not ledgered separately since the task brief scoped a11y to `role="tooltip"` only.
- ADAPT (verified in-browser): exit animations are inert — the native popover is display:none before data-state="closed" is stamped, so data-[state=closed]:animate-out never runs; closing snaps. Accepted for v1; an animated-close strategy (beforetoggle/allow-discrete or dialog-style requestClose) is in the README backlog.

## popover
- WIN: this is `dropdown.gsx`/`dropdown.js`'s own popover-API mechanism (native top layer replaces Portal, `popover="auto"` gives free light dismiss and Esc, `closest("[data-gsxui-popover]")` proximity wiring replaces Radix client context) with the menu semantics stripped: no `role="menu"`, no `role="menuitem"` items, no arrow-key roving focus, no click-on-item close. `PopoverContent` holds arbitrary content (a form, free text — shadcn's own demo is a dimensions form), not a list of selectable items, so none of `dropdown.js`'s item-focused machinery has anything to attach to.
- ADAPT: alignment is centered, not left-aligned — Radix's own `Popover` default is `side="bottom" align="center"`, unlike `DropdownMenuContent`'s `align="start"`. `popover.js`'s positioning calc is `r.left + r.width / 2 - content.offsetWidth / 2` (trigger's horizontal midpoint minus half the content's own width) instead of dropdown.js's bare `r.left`; `origin-top` (not dropdown's `origin-top-left`) replaces shadcn's Radix runtime transform-origin var (`--radix-popover-content-transform-origin`) for the same reason — the enter/exit scale always originates from top-center under a centered-below anchor.
- NOTE: same stopgap as dropdown's own NOTE — a hand-rolled `position: fixed` anchor to the trigger's `getBoundingClientRect()`, not Radix's Floating-UI collision-aware placement; no viewport-edge clamping is done (CSS anchor positioning is the eventual native replacement, not yet Baseline). `data-side="bottom"` is stamped statically since `popover.js` always anchors below.
- GAP: `PopoverAnchor` (a separate positioning-reference part, decoupled from the visible trigger) and `PopoverHeader`/`PopoverTitle`/`PopoverDescription` (structural sub-parts `registry/new-york-v4/ui/popover.tsx` has grown at its current HEAD) are not ported — out of scope per the task brief's own parts list (`Popover`/`PopoverTrigger`/`PopoverContent` only), the same source-revision scoping call as `## alert-dialog`'s own NOTE. A caller wanting a header/title/description inside `PopoverContent` composes plain `<div>`/`<h4>`/`<p>` directly, exactly as the `site/examples/popover/basic.gsx` demo does.
- MECHANISM (trigger click-close guard): identical to dropdown.js's own — `popover="auto"` light-dismisses on the trigger's own outside `pointerdown` before the subsequent `click` fires, so `popover.js` records `trigger.dataset.gsxuiWasOpen` on `pointerdown` and the `click` handler reads/clears it, converging on the popover's actual current state rather than assuming it closed (see `## dropdown`'s own MECHANISM entry for the full rationale).
- A11Y: `aria-expanded="false"` is server-rendered on `PopoverTrigger`, synced by `popover.js` on the `toggle` event (capture-delegated, non-bubbling) — same pattern as dropdown's own, minus `aria-haspopup="menu"` (a popover isn't a menu, so `aria-haspopup` is not stamped at all).
- ADAPT (verified in-browser, same as dropdown/tooltip): exit animations are inert — the native popover is `display:none` before `data-state="closed"` is stamped, so `data-[state=closed]:animate-out` never runs; closing snaps. Accepted for v1, same ledgered post-v1 backlog item as dropdown/tooltip.
- Registry: `popover.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("popover")` is empty (the site example composes `ui.Button`/`ui.Label`/`ui.Input`, which `internal/registry` never scans — same shape as `## toggle`'s own deps entry). `HasJS("popover")` is `true` (`ui/popover.js`).

## hover-card
- WIN: this is `tooltip.gsx`/`tooltip.js`'s own popover-API mechanism (`popover="manual"` for a top-layer element with no light dismiss, so hover/focus-driven show/hide never races an outside-click dismissal) with the arrow dropped — hover-card has none, unlike tooltip's diamond `<span>` — and anchored BELOW the trigger instead of above: Radix's own `HoverCard` default side is `bottom`, `Tooltip`'s is `top`. The `top` calc is flipped (`r.bottom + 4` — Radix's own `sideOffset` default for both `popover.tsx` and `hover-card.tsx` is 4 — replacing tooltip's `r.top - 6 - content.offsetHeight`, whose extra 6px cleared room for the arrow this component doesn't have); the `left` calc (already centered in tooltip.js) is reused unchanged.
- MECHANISM (`HoverCardTrigger` is a `<span>`, not a `<button>`): shadcn's own `hover-card-demo.tsx` asChild-wraps a link-styled `<Button variant="link">`, and Radix's `HoverCardTrigger` itself typically renders as an `<a>` — a hover card almost always previews a link's target. A `<span>` imposes none of `DialogTrigger`/`TooltipTrigger`'s button-in-button trap (`## dialog` FINDING): `<button>` is legal phrasing content inside a plain `<span>` (a `<span>` isn't itself interactive content, unlike `<button>`), so `site/examples/hovercard/basic.gsx` composes `<ui.Button variant="link">` as `HoverCardTrigger`'s real child directly — no data-attribute idiom needed on the child at all, since `HoverCardTrigger`'s own root already carries `data-gsxui-hovercard-trigger`. `asChild` itself is still not ported (no dynamic tag) — this is the same narrow gap as button's/breadcrumb's own `asChild` entries, just resolved by direct composition instead of the data-attribute idiom (`## dialog` MECHANISM) since no wrapper/trap exists here to route around.
- ADAPT: Radix HoverCard's own `openDelay`/`closeDelay` (700ms/300ms) replace tooltip.js's flat 300ms-open/immediate-close. Because `HoverCardContent` can hold real interactive content (shadcn's own demo: an avatar plus a bio, and in general a link or button), a `closeDelay` only does its job if hovering onto the content itself also cancels a pending close — `hover-card.js` adds listeners tooltip.js has no equivalent of on `[data-gsxui-hovercard-content]`, canceling/rescheduling the same trigger-keyed timer the trigger's own listeners manage. Opening (`focusin`) stays immediate on both trigger and content, no delay — Radix's `openDelay` is documented as hover-only. See the MECHANISM entry below for the closing (leave) side, which is NOT immediate for either input modality.
- MECHANISM (unified hover/focus close-grace model, 2026-07-24 fix — a task reviewer caught this against the original port, which routed keyboard `focusout` straight to an immediate `hide()`): Radix treats hover and keyboard focus as one "is the user still interacting with this trigger-or-content pair" model on the way OUT, not two independently-timed controls — leaving the trigger by *either* modality (`pointerout` or `focusout`) schedules the same `closeDelay`-gated `hide()` (`scheduleHide`, not a bare `hide()`), and entering the content by *either* modality (`pointerover` or `focusin`) cancels that pending close, symmetric with the trigger's own pair. Without the keyboard leg of this, a Tab press moving focus off the trigger and onto a focusable child of `HoverCardContent` (a link, per the design brief's own "children usually carry the real link" framing) raced the popover being `display:none`'d by `hide()` out from under the very element the user just tabbed onto. `hover-card.js`'s content `pointerout` handler also gained a `relatedTarget` guard — `if (e.relatedTarget instanceof Element && content.contains(e.relatedTarget)) return;`, the same one `dropdown.js`'s own content `pointerout` handler already carries — so a pointer move between two children still inside the content doesn't needlessly churn a schedule/clear pair; leaving the content for anywhere outside both the trigger and the content (the true tab-out-entirely / pointer-truly-left case) has nothing left to cancel the just-scheduled close, so `hide()` fires after `CLOSE_DELAY` exactly as intended.
- ADAPT: `origin-top` (not tooltip's `origin-bottom`) replaces shadcn's Radix runtime transform-origin var (`--radix-hover-card-content-transform-origin`) — the content is centered BELOW the trigger (Radix's own `align="center"` default for both `Popover` and `HoverCard`), so its scale/fade animation originates from top-center, the same substitution as `## popover`'s own `origin-top` entry (not `origin-top-left`, which was dropdown's own align=start case).
- WIN: `popover="manual"` on `HoverCardContent` — same load-bearing reason as `TooltipContent`'s own WIN entry: `"auto"` popovers light-dismiss on outside pointerdown, which would race `hover-card.js`'s own pointerout/focusout hide logic.
- NOTE: no `role` is stamped on `HoverCardContent` (unlike `TooltipContent`'s `role="tooltip"`) — Radix's own `HoverCardContent` sets no ARIA role either; a hover card holds real content (Radix's own docs call it a "preview card"), not a tooltip announcement, so `role="tooltip"` would misrepresent it to assistive tech.
- GAP: no `data-side` value beyond the static `"bottom"` this port always stamps — same stopgap-until-CSS-anchor-positioning NOTE as dropdown/tooltip/popover: a hand-rolled `position: fixed` anchor, not Radix's Floating-UI collision-aware placement, no viewport-edge clamping.
- ADAPT (verified in-browser, same as dropdown/tooltip/popover): exit animations are inert — the native popover is `display:none` before `data-state="closed"` is stamped, so `data-[state=closed]:animate-out` never runs; closing snaps. Accepted for v1, same ledgered post-v1 backlog item.
- NOTE (data attribute naming): `data-gsxui-hovercard`/`-trigger`/`-content` drop the hyphen the `hover-card` component basename and its `.js`/`.gsx` filenames both carry — `HasJS` derives strictly from `<basename>.js` (so the behavior file is `ui/hover-card.js`, hyphen intact), but the in-DOM attribute follows the same no-hyphen-in-the-selector-prefix convention as every other multi-word interactive component's own hook (e.g. the upcoming context-menu task's `data-gsxui-contextmenu-trigger`, not `-context-menu-trigger`) — a JS `dataset` property name can't contain a hyphen without becoming camelCase (`dataset.gsxuiHoverCard`), so the attribute itself is written without one instead.
- Registry: `hover-card.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("hover-card")` is empty (the site example composes `ui.Avatar`/`ui.Button`, which `internal/registry` never scans — same shape as `## popover`'s own deps entry). `HasJS("hover-card")` is `true` (`ui/hover-card.js`).

## context-menu
- WIN: this is `dropdown.gsx`/`dropdown.js`'s own mechanism, reused verbatim for menu semantics (`popover="auto"` gives the top layer/light-dismiss/free-Esc win, `role="menu"`/`role="menuitem"`, arrow-key roving focus, close-on-select, `closest("[data-gsxui-contextmenu]")` proximity wiring) with one structural swap: `ContextMenuTrigger` is an AREA (a plain `<div>`, shadcn's own demo renders a dashed-border drop-zone), not a `<button>` — a context menu opens on a right-click anywhere inside an arbitrary region, not a left-click on a single control.
- MECHANISM (open on `contextmenu`, not `click`): `context-menu.js` listens for the `contextmenu` event (which bubbles, unlike `toggle` — no `{ capture: true }` needed for this one listener) delegated to `[data-gsxui-contextmenu-trigger]` and calls `preventDefault()` to suppress the native browser context menu, in place of dropdown.js's `click` listener on the trigger BUTTON.
- ADAPT (cursor positioning, not trigger-rect anchoring): dropdown.js/popover.js/hover-card.js all anchor to a single element's `getBoundingClientRect()` (the trigger); a context-menu trigger is an AREA, not a point, so there is no one rect to anchor to. `context-menu.js` instead positions at the cursor — `event.clientX`/`event.clientY` off the `contextmenu` event itself — the same numeric-position-after-`showPopover()`-never-`style.transform` rule as every sibling (the enter animation's keyframes animate `transform`, so a positioning translate would be overridden for the animation's duration and the popover would enter at the untranslated spot and snap).
- ADAPT (viewport-edge clamping — the one deviation from the dropdown/tooltip/popover/hover-card family's own documented no-clamp NOTE): those four siblings anchor to a FIXED side of a known trigger element and accept imprecision near viewport edges as a stopgap until CSS anchor positioning (`anchor()`/`position-anchor`) is Baseline — their own NOTE entries say so explicitly. A context menu has no fixed anchor at all: it opens wherever the cursor was, which on an ordinary (not just an edge-case) right-click can be arbitrarily close to the viewport's own right or bottom edge — an unclamped menu would routinely render partially or fully offscreen, not rarely. `context-menu.js` clamps `left`/`top` against `document.documentElement.clientWidth - content.offsetWidth`/`clientHeight - content.offsetHeight` (client metrics, not `window.inner*` — the inner metrics include classic-scrollbar gutters and clamping against them tucks the menu edge under the scrollbar; found and fixed in the Tier 2 browser pass) (read AFTER `showPopover()` — a hidden popover has no layout box, same ordering rule as the positioning ADAPT above), floored at `0`. This is a genuine, deliberate divergence from the sibling precedent, not an omission — ledgered here per the task brief's explicit instruction to comment and ledger the deviation.
- MECHANISM (no trigger-click-close guard, unlike dropdown/popover's own): those two need the `pointerdown`/`click` `wasOpen` flag (`## dropdown`'s own MECHANISM entry) because a left-click on the trigger BUTTON is itself an *outside* pointerdown relative to the content, racing `popover="auto"`'s own light dismiss before the `click` handler runs. `context-menu.js` has no equivalent: a right-click's `pointerdown` ALSO counts as an outside pointerdown and already light-dismisses an open menu before the `contextmenu` event fires, so by the time the single `contextmenu` handler runs the popover has normally already closed on its own — no separate `click` listener, no flag, nothing to converge. The handler's own `if (content.matches(":popover-open")) content.hidePopover();` before re-showing is purely defensive, for a `contextmenu` event dispatched without a preceding `pointerdown` (e.g. the keyboard Menu key) — this is also what makes a second right-click inside the trigger area reposition the menu rather than merely toggling it closed: the first right-click's own `pointerdown` already closed the open popover before the second `contextmenu` fires, so the handler always falls through to the show-and-position path at the new cursor coordinates.
- ADAPT: NO `data-side` is ever stamped on `ContextMenuContent`, unlike `DropdownMenuContent`/`PopoverContent`/`HoverCardContent`'s own static stamps — those three always open toward one fixed side of their trigger (dropdown/popover below, hover-card below), so a static `data-side="bottom"` is accurate for every render. A context menu's effective side varies with both the cursor position and the clamping above; there is no single value that would be accurate, so none is stamped at all, and the class string's `data-[side=*]:slide-in-from-*` selectors are permanently dead weight — more so than dropdown's own (whose `data-[side=bottom]` selector IS live, only the other three are dead) — kept for future-proofing per dropdown's own precedent (`## dropdown`'s own NOTE).
- ADAPT: `origin-top-left` replaces shadcn's Radix runtime transform-origin var (`--radix-context-menu-content-transform-origin`), the same substitution shape as dropdown's own — but here it is a majority-case approximation, not an always-accurate one: Radix computes the real origin dynamically from the actual collision-resolved placement, dropdown's static choice is exact (always-below-left-aligned), while context-menu's is exact only for the common case (menu opens below-right of the cursor, unclamped) and visually approximate once viewport-edge clamping (the ADAPT above) has shifted the menu away from the cursor — accepted for v1, no per-render origin computation implemented.
- GAP: `ContextMenuCheckboxItem`, `ContextMenuRadioGroup`/`ContextMenuRadioItem`, `ContextMenuSub`/`ContextMenuSubTrigger`/`ContextMenuSubContent`, `ContextMenuGroup`, and `ContextMenuPortal` are not ported — same post-v1-backlog scoping call as `## dropdown`'s own GAP entry (its `DropdownMenuCheckboxItem`/`DropdownMenuRadioGroup`&`RadioItem`/`DropdownMenuSub`-family/`DropdownMenuGroup`). Only `ContextMenu`, `ContextMenuTrigger`, `ContextMenuContent`, `ContextMenuItem`, `ContextMenuLabel`, `ContextMenuSeparator`, `ContextMenuShortcut` ship.
- ADAPT: the `inset` prop on `ContextMenuItem`/`ContextMenuLabel` is dropped along with it — `data-[inset]:pl-8` is removed from both classes, the same "drop the selector, don't ship dead CSS" call as dropdown's own identical ADAPT. With `inset` gone, `ContextMenuItem`'s remaining class string is byte-identical, token-for-token, to `DropdownMenuItem`'s own pinned class (`TestDropdownPinned`/`TestContextMenuPinned`) — the two shadcn sources share every other token; a coincidence of upstream authoring, not a copy-paste port of one from the other.
- NOTE (a genuine upstream difference, not a copy error): `ContextMenuLabel`'s class carries `text-foreground`; `DropdownMenuLabel`'s does not (`px-2 py-1.5 text-sm font-medium` only, verified against `registry/new-york-v4/ui/dropdown-menu.tsx` on disk). Both are ported verbatim from their own binding-to-disk source — shadcn's two menu components simply aren't kept in lockstep with each other upstream.
- MECHANISM: `ContextMenuTrigger` carries no `aria-haspopup`/`aria-expanded` of its own, unlike `DropdownMenuTrigger` — it is a passive AREA, not an interactive control (Radix's own `ContextMenuTrigger` renders a plain wrapper, not an ARIA-widget-shaped element either); there is no toggle-on-click affordance to describe to assistive tech the way a dropdown trigger button has. `context-menu.js`'s `toggle` handler correspondingly has no trigger to sync `aria-expanded` onto and no `gsxuiWasOpen` flag to clear — a straight subset of dropdown.js's own toggle handler with the trigger-aria lines removed.
- Registry: `context-menu.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("context-menu")` is empty (the site example composes nothing from another `ui.*` component either — same shape as `## popover`'s own deps entry). `HasJS("context-menu")` is `true` (`ui/context-menu.js`).

## slider
- WIN (ARIA + keyboard entirely for free): the four-part Radix tree (`Root`/`Track`/`Range`/`Thumb`) collapses onto ONE native `<input type="range">`, which already implicitly carries `role="slider"` and auto-derives `aria-valuemin`/`aria-valuenow`/`aria-valuemax` from its own `min`/`value`/`max` attributes — zero markup needed to reproduce the Thumb's ARIA contract. The browser's own keyboard model (Left/Right/Down/Up adjust by `step`, PageUp/PageDown by a bigger jump, Home/End to min/max) already matches Radix's own traced behavior contract (`BACK_KEYS`/`PAGE_KEYS`/`ARROW_KEYS`, Shift-multiplies-step-by-10) closely enough that no keyboard JS is needed either — this is the single biggest reason the controls source map called this ADAPT "styled native `<input type=range>`."
- ADAPT (gradient fill, not accent-color): the Radix Range part (the primary-colored portion from min to the current value) has no free native equivalent — a bare styled range input has one uniform track. Ported as a `linear-gradient()` on `::-webkit-slider-runnable-track`/`::-moz-range-track` (`assets/gsxui.css` + `web/site.css`) keyed off a `--fill` percentage custom property, NOT `accent-color` — `accent-color`'s cross-browser track-fill behavior (beyond coloring the thumb) is unverified across Chromium/Firefox/Safari, while the gradient is guaranteed cross-browser-identical. `ui/slider.gsx` computes the INITIAL `--fill` server-side with exact float64 arithmetic from `value`/`min`/`max` (handles `min != 0` correctly, e.g. `min=20 max=40 value=25` → `--fill: 25%`, not the wrong-if-min-were-ignored 62.5%) — zero JS needed for first paint. `ui/slider.js` is a single delegated `input` listener (via `on()`, `ui/gsxui.js`'s core) that recomputes `--fill` from the range's own live `min`/`max`/`value` on every drag or keystroke thereafter.
- GAP (hit-target regression, size-3 thumb with no compensation): nova's own thumb carries `after:absolute after:-inset-2`, an invisible `::after` enlarging the hoverable/draggable area beyond the `size-4→size-3` shrink. `::-webkit-slider-thumb`/`::-moz-range-thumb` are themselves UA-shadow pseudo-elements; CSS does not allow a pseudo-element on a pseudo-element (`::-webkit-slider-thumb::after` is not a legal/supported selector), so this compensation is categorically unreachable on a native-input port. This port adopts nova's `size-3` visual anyway (per task Decisions) and accepts a smaller real hit target than nova's own React/Radix version has — visual parity without the interaction parity nova's own extra markup buys it. Ledgered explicitly, not silently dropped.
- GAP (single scalar value only): a native range input is always ONE thumb, ONE value. Radix's multi-thumb/range-slider support (two or more thumbs, e.g. a min/max price range — `defaultValue={[25,50]}`), `orientation="vertical"`, `inverted` (flips which visual end is "low"), and `minStepsBetweenThumbs` (prevents one thumb's move from crossing closer than N steps to a neighbor) all have no analog on this ADAPT and are not ported. A real range slider (two thumbs) would need Radix's actual multi-part architecture or a from-scratch custom-JS thumb pair, not this single-`<input>` ADAPT. Only the in-scope `slider-demo.tsx` shape (single thumb, `defaultValue={[50]} max={100} step={1}`) is ported as `site/examples/slider/basic.gsx`.
- MECHANISM (WebKit thumb vertical centering): Firefox's `::-moz-range-track` box model auto-centers `::-moz-range-thumb` on the track; WebKit's `::-webkit-slider-thumb` does not — left alone, the thumb sits flush with the track's top edge instead of centered on it. Compensated with `margin-top: calc((0.25rem - 0.75rem) / 2)` (track height minus thumb height, halved — nova's `h-1`/`size-3` metrics) on the WebKit selector only; Firefox needs no such rule. A well-known cross-browser quirk that looks fine in isolation and visibly misaligned next to a real shadcn slider without it.
- MECHANISM (paired cross-browser pseudo-elements, no single selector reaches every engine): track (`::-webkit-slider-runnable-track` + `::-moz-range-track`) and thumb (`::-webkit-slider-thumb` + `::-moz-range-thumb`) each need both selector families written out, `assets/gsxui.css`/`web/site.css` (the same shared hand-authored CSS location as `## accordion`'s `::details-content` block — nothing here is reachable via a Tailwind utility class). Focus ring is the same story: `input:focus-visible::-webkit-slider-thumb` is valid in Chromium/Safari, but Firefox requires the `::-moz-range-thumb` pairing instead — both are written for hover/focus-visible/active alike (nova's `hover:ring-3 focus-visible:ring-3 active:ring-3`, `active:ring-3` being a wholly new state nova adds), each resolving to `box-shadow: 0 0 0 3px color-mix(in oklab, var(--ring) 50%, transparent)` — the same formula Tailwind itself compiles `ring-[3px] ring-ring/50` to (verified against the site's own compiled `site/dist/assets/main-*.css`).
- ADAPT (color scope, nova's border-ring NOT adopted): nova recolors the thumb border to `border-ring`; this port keeps new-york-v4's `border-primary` — a color substitution, not a metric one, out of this retarget's stated scope (heights/paddings/gaps/text/radii/svg sizes only), same house policy as every other new-port "keep border-primary" call.
- ADAPT (max/step zero-value defaulting): `max`/`step` fall back to shadcn's own defaults (100/1) when left at the Go zero value (`0`); `min`'s zero value needs no such fallback since 0 is already Radix's own default. Same unset-vs-explicit-zero ambiguity every other zero-value-defaulted param in this codebase accepts (e.g. Toggle's `variant`/`size` via `|> default(...)`) — not reachable in practice for a slider (step=0 is invalid on a native range input regardless of this port).
- Registry: `slider.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("slider")` is empty, same shape as `## toggle`'s/`## popover`'s own deps entries. `HasJS("slider")` is `true` (`ui/slider.js`).

## scroll-area
- ADAPT (CSS-first, collapsed single div): the roadmap's own stated preference for this component ("CSS `scrollbar-width`/`scrollbar-color` styling first; Radix-style custom thumbs only if that falls short") is taken all the way — Radix's four-part tree (`Root`/`Viewport`/`ScrollBar`/`Thumb`) collapses onto ONE native `<div class="overflow-auto ...">`, `Root`'s own `relative` and `Viewport`'s own `rounded-[inherit]`/`transition-[color,box-shadow]`/`outline-none`/`focus-visible:ring-[3px] focus-visible:ring-ring/50 focus-visible:outline-1` both porting unchanged onto the same node. A real `overflow: auto` element keeps every native scroll input modality (wheel/trackpad/touch/keyboard) working with zero JS, and RTL scrollbar layout comes free with it too — strictly no worse than Radix's own approach, which keeps native scrolling underneath its custom chrome regardless (non-gap, named for completeness per the source map).
- MECHANISM (dual-CSS-surface, standard + legacy WebKit, layered not either/or): the standardized `scrollbar-width: thin` / `scrollbar-color: var(--border) transparent` pair (CSS Scrollbars Module Level 1) rides in `ui/scroll-area.gsx`'s own class string as Tailwind arbitrary properties — Baseline-available in Firefox/Chromium for years, Safari since 18.2 — and alone gets a themed, thin(ner) scrollbar close to shadcn's visual weight with zero JS. It cannot control thumb SHAPE (`rounded-full`, the `p-px` inset gap off the track edge) — the standard property paints the browser's own squared-off thumb geometry, no radius/inset/border. Full shape fidelity needs the legacy, WebKit-proprietary (Chromium-inherited) `::-webkit-scrollbar`/`-track`/`-thumb` pseudo-elements layered on top, scoped to `[data-slot="scroll-area"]` in `assets/gsxui.css` AND `web/site.css` (dual-home, byte-identical — same `::details-content`/range-input hand-authored-CSS precedent as `## accordion`/`## slider`). Firefox has never implemented the legacy pseudo-elements at all, so both surfaces are required together, roughly doubling the CSS surface for one visual result — ledgered as ongoing maintenance surface, not a one-time cost. `var(--border)` needs no separate `dark:` variant: it is itself a theme token that already changes value under `.dark`, so the thumb color tracks the theme for free.
- GAP (no `ScrollAreaCorner`): `::-webkit-scrollbar-corner` exists in WebKit/Blink only; the standard `scrollbar-*` properties have no corner concept whatsoever; Firefox has no way to style the intersection square under any mechanism, period — the one part with no browser-portable equivalent at all, not just an unstyled fallback. Dropped outright for v1, matching this codebase's "drop the part, don't ship dead CSS" convention (e.g. dropdown's never-ported scroll buttons, `## accordion`'s dropped `Header` wrapper).
- GAP (no `type` visibility-timing prop): Radix's `hover`/`scroll`/`auto`/`always` scrollbar show/hide TIMING is entirely OS/browser policy on a native scrollbar — macOS's own "Show scroll bars: Always / When scrolling / Automatically" System Settings preference overrides any page-level attempt to mimic Radix's show-on-hover/fade-after-600ms behavior; a user with "Always" set never sees this port's scrollbar "hide" no matter what CSS says, and vice versa. Not a fixable styling gap — ledgered as an accepted platform-policy override the same way `## dropdown`'s inert exit-animation ADAPT accepts a browser-imposed limit rather than working around it.
- GAP (no `ScrollBar`/`ScrollAreaThumb` as separate components): both parts collapse into the one div's own class string/hand-authored CSS above — there is nothing left for a standalone `ScrollBar`/`ScrollAreaThumb` gsx component to own, so neither ships. A caller cannot independently restyle "just the thumb" or "just the track" the way a Radix consumer could re-target `ScrollAreaScrollbar`'s own `className`; only the whole scrollbar's color/width (via `class` fall-through onto the collapsed div, which the standard properties read from the same element) and, upstream in this library's own CSS, the WebKit shape.
- WIN (thumb SIZE, non-gap named for completeness): native thumb size genuinely IS the viewport/content ratio, identical to what Radix's own `getThumbSize` computes from `ResizeObserver`-measured viewport/content dimensions — nothing to reproduce, the browser already does the same math for free.
- ADAPT (`orientation`, no second component): `orientation="horizontal"` picks `overflow-x-auto` over the default `overflow-auto` via a plain Go `switch` inside `class={}` (the same idiom `## button-group`'s own orientation switch uses) — Radix needs a second literal `<ScrollBar orientation="horizontal"/>` element opted into explicitly; here it is one different overflow utility on the same collapsed div, no markup shape change. No `data-orientation` attribute is stamped: shadcn's own `Root` never carries one either (only the dropped `ScrollBar` part did), so there is nothing for a CSS selector to key off here — unlike `## separator`, which needs the attribute because its own `data-[orientation=...]:` selectors do the dispatching.
- Registry: `scroll-area.gsx` imports nothing from `ui/icon` and calls no other `ui.*` component — `registry.Deps("scroll-area")` is empty, same shape as `## toggle`'s/`## popover`'s/`## slider`'s own deps entries (the site examples compose `ui.Separator`, but `internal/registry` only scans `ui/*.gsx`, not `site/examples/`). `HasJS("scroll-area")` is `false` — no `ui/scroll-area.js` ships; every behavior this port has (scroll input, thumb color/shape, focus ring) is either native or pure CSS.

## drawer
- WIN: `Drawer`/`DrawerTrigger` compose `ui.Dialog` directly, byte-for-byte
  the same mechanism as `Sheet`/`AlertDialog` (`## sheet` WIN) — `Drawer`
  overrides the composed `data-slot` (`"dialog"`→`"drawer"`), so
  `drawer → dialog` derives and the CLI vendors `ui/dialog.js` unmodified:
  trigger/content proximity wiring, Esc-to-close, `data-state` stamping,
  and the `getAnimations`-gated exit-animation wait are all reused as-is.
- ADAPT (the one part that does NOT compose): `DrawerContent` renders its
  own native `<dialog data-slot="drawer-content" data-gsxui-dialog-content
  data-state="closed" data-side="...">` rather than composing
  `ui.DialogContent`/`ui.SheetContent` — same reasoning as `## sheet`'s own
  ADAPT: the centered-card recipe, `Sheet`'s side-anchored recipe, and
  `Drawer`'s own per-direction recipe all target the same CSS properties
  with materially different values, so no single class string merges them.
  Unlike `DialogContent`/`SheetContent`, `DrawerContent` never injects a
  close X button — upstream `drawer.tsx` has no `showCloseButton`-
  equivalent prop at all (dismissal is backdrop-click/Esc, or an explicit
  `DrawerClose` the caller places wherever the design wants it) — so
  `DrawerContent` has no `hideCloseButton` param, unlike `SheetContent`'s.
- MECHANISM (naming, `direction` not `side`): `Drawer` has no vaul
  underneath to key a client CSS attribute selector off of — direction
  selection happens server-side via a Go `switch` in `class={}`, the same
  idiom as `Sheet`'s own `side` switch. The Go param is named `direction`
  (vaul's own vocabulary, distinct from `Sheet`'s `side`) per the task's
  binding decision, both selecting the class-string case and stamping
  `data-side` — reusing `Sheet`'s own internal attribute name (for any
  future shared tooling/CSS keying off it generically) rather than
  inventing a second attribute. `direction`'s default is `"bottom"` (vaul's
  own conventional anchor, the mobile bottom-sheet pattern) — a real
  behavioral difference from `Sheet`'s own `"right"` default.
- ADAPT (per-direction class strings, Sheet's UA-default fixes carried
  forward): `m-0`, `open:flex`, and each side's opposite-edge `-auto`
  utility are all carried from `## sheet`'s own three-part ADAPT
  (transcribed onto `drawer.tsx`'s values, not independently re-derived).
  TOP/BOTTOM keep their own explicit `max-h-[80vh]` (author-origin, already
  beats Chrome's UA `max-height` safety net) and so do NOT need `Sheet`'s
  `max-h-none` escape hatch; LEFT/RIGHT use `h-full` like `Sheet`'s own
  left/right and DO need `max-h-none`, copied verbatim. TOP/BOTTOM also
  carry `Sheet`'s own `w-full max-w-none` pair (the plain-`<div>`
  auto-width edge-to-edge stretch does not reproduce on a native
  `<dialog>`). These strings are **transcribed-not-verified**: `Sheet`'s
  six ADAPTs were only found by rendering `site/examples/sheet` in a real
  browser tab (`## sheet`'s own three-part fix note) — the controller runs
  that same in-browser verification pass for all four drawer directions
  (`site/examples/drawer`) before this task is considered closed; this
  entry will be updated with any correction that pass finds.
- NOVA (the one deliberate `bg-popover`/`text-popover-foreground`
  exception): every other retargeted component in this codebase kept
  `bg-background`/`text-foreground` per the density retarget's own
  NOT-ADOPTED color-scope stance (`## nova density`) specifically to stay
  consistent with an existing gsxui baseline. `Drawer` has no prior
  new-york-v4-based gsxui version to preserve continuity with — there is no
  baseline to break — so this port adopts nova's `bg-popover
  text-popover-foreground` on `DrawerContent` outright, confirmed per the
  task's own binding decision. `rounded-*-xl` on all four directions is
  also nova (new-york-v4 rounds only top/bottom; left/right get a plain
  border, no rounding at all — nova rounds every direction, including the
  free/non-anchored edge on left/right). Backdrop (`bg-black/10` +
  `supports-backdrop-filter:backdrop:backdrop-blur-xs`, `duration-200` both
  directions) is identical to `Sheet`'s own.
- GAP (handle bar kept as decoration): upstream renders an unconditional
  inline `<div>` (not a named sub-component), visible only for the bottom
  direction via a `group-data-[vaul-drawer-direction=bottom]/drawer-
  content:block` override on a base `hidden` class. `Drawer` has no vaul
  underneath and direction is server-known, so this port replaces the
  client CSS group-data selector with a plain Go `if direction == "" ||
  direction == "bottom"` — the handle (`data-slot="drawer-handle"`, `h-1`
  nova thinner than new-york-v4's `h-2`) is rendered ONLY for the bottom
  direction, absent (not merely hidden) otherwise. v1 ships no drag
  gesture, so a "drag me" affordance that doesn't actually drag is a real,
  if minor, UX mismatch — accepted as a GAP (visual parity chosen over
  silently dropping the handle).
- MECHANISM (`DrawerHeader`'s direction-conditional text-alignment, ported
  faithfully via `data-side` + `group/drawer-content`): upstream's
  `DrawerHeader` class carries direction-conditional text alignment
  (centered for bottom/top at every breakpoint, left-aligned at `md:`+ for
  left/right) via `group-data-[vaul-drawer-direction=...]/drawer-content`
  — the same selector the handle bar's visibility rule is keyed off of.
  `DrawerHeader` takes no `direction` param of its own (matching
  `drawer.tsx`'s own signature), but it doesn't need one: `DrawerContent`
  already stamps `data-side` (always non-empty — `direction |>
  default("bottom")`, `## drawer`'s own MECHANISM entry above) and now
  also carries the named-group class `group/drawer-content` on the same
  element, so the identical selector SHAPE ports directly with only the
  attribute/group name swapped — `data-side` replaces
  `data-vaul-drawer-direction`, `drawer-content` (already
  `DrawerContent`'s own `data-slot`) is the group name.
  `DrawerHeader`'s class is therefore `flex flex-col gap-0.5 p-4
  group-data-[side=bottom]/drawer-content:text-center
  group-data-[side=top]/drawer-content:text-center md:text-left`, byte-for-
  byte the same alignment logic as upstream. This uses the same
  `group/name` + `group-data-[...]/name` idiom already established
  elsewhere in this codebase (`## item`, `## field`, `## input-group`,
  `## tabs`, `## toggle-group`) — not a new mechanism, just its first
  cross-component (ancestor-stamped-attribute rather than
  same-element-state) application. (Coordinator review fix, replacing an
  earlier draft of this entry that had incorrectly GAPed this as dropped —
  the mechanism was available with zero signature change and is now
  ported.)
- GAP (nova `gap-0.5` header at every breakpoint): new-york-v4's own
  `DrawerHeader` bumps to `md:gap-1.5`; nova drops the responsive bump
  entirely, staying `gap-0.5` at every breakpoint — adopted as-is (no
  divergence, just recording the nova delta this port took).
- GAP (drag-to-dismiss / snap points / background scaling not ported): v1
  replaces vaul's live-transform drag physics entirely with the same
  `<dialog>` + Tailwind-keyframe slide-in/out architecture `Sheet` already
  uses — no pointer tracking, no spring-back, no snap-point stops, no
  `shouldScaleBackground` (upstream's own `drawer.tsx` doesn't default it
  on either, so no shadcn docs demo exercises it). Matches the roadmap's
  own Tier 3 listing verbatim ("sheet variant; v1 without vaul's
  drag-to-dismiss gesture, ledger the gap").
- GAP (`asChild` not ported): same narrow gap as `## dialog`'s/`## sheet`'s
  own MECHANISM entries — the `data-gsxui-dialog-trigger`/
  `data-gsxui-dialog-close` attribute idiom is the replacement throughout.
- GAP (demo scope, `drawer-dialog` deferred): upstream ships two demos —
  `drawer-demo.tsx` (adapted here, see `site/examples/drawer/basic.gsx`'s
  own doc comment for the recharts-sparkline-to-static-bars and
  useState-stepper-to-static-markup substitutions, neither of which this
  codebase has an equivalent for) and `drawer-dialog.tsx` (a responsive
  `useMediaQuery`-driven `Dialog`-on-desktop/`Drawer`-on-mobile swap, no
  gsxui equivalent for a JS media-query-driven component swap). The latter
  is deferred to a future patterns/pages-phase example per the roadmap
  (`docs/component-roadmap.md`'s new drawer-dialog note), out of scope for
  this component task.
- Registry: `drawer.gsx` imports nothing from `ui/icon`; `Drawer` calls
  `ui.Dialog` (flat-package intra-package edge, same shape as `## sheet`'s
  own `dialog` dep) — `DrawerTrigger`/`DrawerClose` render their own
  `<button>` rather than composing `ui.Button`, so `dialog` is the ONLY
  dep — `registry.Deps("drawer") == ["dialog"]`, pinned in
  `internal/registry/registry_test.go`. `HasJS("drawer")` is `false` — like
  `sheet`, it has no behavior module of its own, only `ui/dialog.js`,
  pulled in transitively.

## carousel
- ADAPT (scroll-snap substitution, the component's central decision): embla
  (`embla-carousel-react`) is a JS-transform carousel — it applies
  `transform: translate3d(...)` to the flex track div directly, driven by
  its own pointer-drag-physics engine, and uses neither CSS `scroll-snap`
  nor native `overflow: auto` anywhere. This port replaces the whole
  mechanism with real native CSS scroll-snap: `overflow-x-auto snap-x
  snap-mandatory` (or the vertical `-y`/`snap-y` pair) on `CarouselContent`'s
  outer viewport div — REPLACING shadcn's own bare `overflow-hidden`, since
  a scroll container needs `overflow: auto` to scroll at all — plus
  `snap-start` on every `CarouselItem` (both classes are **new**, present in
  neither shadcn's source nor nova's own CSS). `ui/carousel.js` supplies
  only the prev/next scroll-by-one-item calls and disabled-state/
  current-index bookkeeping; there is no drag-physics/momentum engine to
  port because the browser's own native scroll input now supplies it.
- WIN (momentum/rubber-banding, a genuine capability gain, not just a
  substitution): because scrolling is now real native `overflow: auto`,
  touch/trackpad momentum and rubber-banding at the ends come free from the
  browser — embla's transform-based approach never had anything analogous,
  it only ever animated to discrete snap points on release.
- GAP (`loop`, no native-scroll analog): embla's infinite-wraparound `loop`
  option works by literally cloning slides at both ends of the track: real,
  reproducible extra work with no clean scroll-snap equivalent (a native
  scroll container has two hard, non-wrapping ends). Not ported — no docs
  demo (`carousel-demo`/`-api`/`-orientation`/`-size`/`-spacing`/`-plugin`)
  requests it, so this is an accepted v1 gap, not a silent drop.
- GAP (`align`, no `center` mode needed or ported): embla's own alignment
  default is commonly documented as `"start"`; every docs demo that varies
  more than one slide per view (`-orientation`/`-size`/`-spacing`) already
  passes `opts={{align:"start"}}` explicitly, and the single-per-view
  baseline (`-demo`) makes alignment moot (one slide always fills the
  viewport). Native `scroll-snap-align: start` on every item is therefore
  sufficient for every demo in scope; no `center`-alignment mode is ported.
- ADAPT (prev/next scroll amount, one item not one viewport): scrolling by a
  full viewport width would visibly skip slides whenever more than one is
  visible per view (the `-size`/`-orientation` demos) — `ui/carousel.js`
  scrolls by the FIRST item's own measured `getBoundingClientRect()` width
  instead (embla's default `slidesToScroll: 1`). Because `CarouselContent`'s
  flex track lays items out with no `gap` property (spacing comes from each
  item's own `pl-4`/`pt-4` padding plus the track's compensating `-ml-4`/
  `-mt-4`), one item's border-box width already IS the distance to the next
  item's own start — no separate "plus the gap" term needed in the
  arithmetic despite the map's own framing describing it that way.
- ADAPT (prev/next disabled-state, computed from scroll position not an
  embla API): embla's `canScrollPrev()`/`canScrollNext()` read its own
  internal scroll-progress/edge state. This port instead compares the
  viewport's own `scrollLeft`/`scrollTop` against `0` and against
  `scrollWidth - clientWidth`/`scrollHeight - clientHeight`, with a 1px
  epsilon for sub-pixel rounding, recomputed on a rAF-throttled `scroll`
  listener AND a `ResizeObserver` on the viewport (covers a responsive
  `basis-*` breakpoint change flipping disabled state or the resting index
  with no scroll event ever firing) — set directly as the native `disabled`
  attribute on the composed `Button`s, so `disabled:` variant classes apply
  for free.
- ADAPT (initial server-rendered disabled state, an asymmetric default):
  shadcn's own `canScrollPrev`/`canScrollNext` both start `false` (a bare
  `useState(false)`, corrected by a layout effect before paint) —
  unobservable in practice for a hydrating React app, but a real,
  observable server-rendered value here, since there is no synchronous
  hydration step. `CarouselPrevious` renders `disabled` from the start: a
  freshly mounted scroll container always has `scrollLeft`/`scrollTop` `0`,
  a real DOM invariant, not a guess — Previous genuinely has nowhere to go
  yet. `CarouselNext` renders enabled (not disabled) by default: whether
  there is more content to scroll TO depends on rendered content/viewport
  widths this Go component cannot measure at render time, so it takes the
  permissive default — a button that turns out to have nothing to do is
  harmless, one that turns out to be wrongly disabled is not.
  `ui/carousel.js`'s own init pass recomputes and corrects both from the
  real DOM immediately on load either way.
- NEW (scrollbar hiding, present in neither source): `CarouselContent`'s
  outer viewport div also carries `[scrollbar-width:none]
  [&::-webkit-scrollbar]:hidden` — needed only because this port is now on
  a real native scroll container that would otherwise show a visible
  scrollbar; embla's transform-based approach never scrolled a real
  overflow box, so it never had a scrollbar to hide in the first place.
  Paired standard + legacy-WebKit selectors, same dual-surface rationale as
  `## scroll-area`'s own MECHANISM entry (Firefox never implemented the
  legacy `::-webkit-scrollbar` family, so both are required together for
  full cross-browser coverage) — here both are plain Tailwind arbitrary
  properties/variants in the class string, no hand-authored CSS file needed
  (unlike scroll-area's shape-fidelity pseudo-elements, which need real
  `::-webkit-scrollbar-thumb` geometry `scrollbar-width` alone can't
  express; hiding the bar entirely needs no shape at all).
- GAP (no context, `orientation` passed explicitly to every part):
  Radix-less — there is no shared context broadcasting `orientation` from
  `Carousel` down to `CarouselContent`/`CarouselItem`/`CarouselPrevious`/
  `CarouselNext` the way shadcn's own `useCarousel()` hook does. Every part
  that needs it takes its own `orientation` param instead, the same
  explicit-prop-instead-of-context shape `## toggle-group`'s
  groupType/variant/size/spacing already establishes.
- WIN (spacing stays caller-controlled, no separate prop needed): embla's
  default `-ml-4`/`pl-4` gap, overridden to `-ml-1`/`pl-1` by
  `carousel-spacing.tsx`, proves gap is entirely a `className` override in
  shadcn's own source — ported as-is via the ordinary class-merge
  mechanism, no separate `spacing` param invented.
- ADAPT (API-surface reduction, deliberate, ledgered): embla's `CarouselApi`
  (`scrollTo(i)`, `scrollPrev()`, `scrollNext()`, `canScrollPrev()`,
  `canScrollNext()`, `selectedScrollSnap()`, `scrollSnapList()`,
  `on(event, cb)`/`off()`, plus its whole plugin-extension mechanism) is not
  reproduced 1:1 — there is no plugin ecosystem to support here. The
  replacement surface is much smaller: a `gsxui:carousel-select` CustomEvent
  (`{index, count}`, both 0-based, emitted on the carousel root whenever the
  resting index actually changes — covers `carousel-api.tsx`'s "Slide X of
  Y" indicator use case completely, see `site/examples/carousel/api.gsx`)
  plus three methods stashed directly on the DOM node
  (`carouselEl.gsxuiCarousel = {scrollTo(i), next(), prev()}`) for any
  script needing imperative control. `data-current-index` is also stamped
  on the root on every change, a CSS-only hook for a caller-authored
  `[data-index="N"]` dot-indicator list needing no JS of its own. Embla's
  full plugin system (Autoplay, ClassNames, Fade, WheelGesture, …) is not
  ported at all; only the one bespoke `data-gsxui-carousel-autoplay`
  attribute below stands in for the single plugin the docs demos actually
  use.
- ADAPT (autoplay, a bespoke data attribute standing in for
  `embla-carousel-autoplay`): `data-gsxui-carousel-autoplay="<ms>"` on the
  root starts a `setInterval` calling `next()` every `<ms>` milliseconds,
  paused on `pointerenter`/`focusin` anywhere within the carousel and
  resumed on `pointerleave`/`focusout` — reproducing `carousel-plugin.tsx`'s
  ACTUAL demo behavior (its own explicit
  `onMouseEnter={plugin.stop}`/`onMouseLeave={plugin.reset}` hover
  pause/resume), not embla Autoplay's own `stopOnInteraction` semantics
  (which trigger on drag/click, not hover). Attached via direct
  `addEventListener` per carousel root at module-init time, not the
  document-delegated `on()` helper — `pointerenter`/`pointerleave` don't
  bubble and are meant to fire only on the exact listened-to element when
  the pointer crosses ITS boundary, which is precisely the "hovering
  anywhere within the carousel" semantic wanted here; delegating them
  through `on()`'s `closest()`-based dispatch would not reproduce that
  boundary-crossing behavior correctly. No loop mode (see the `loop` GAP
  above) — autoplay simply stops advancing once it reaches the last slide
  rather than wrapping back to the first. This attribute is not yet wired
  into any site example (a stretch 5th example per the map's own demo
  inventory, deferred).
- MECHANISM (`scroll` is non-bubbling, delegated via capture): the viewport
  `scroll` listener is registered `on("scroll", '[data-slot="carousel-
  content"]', ..., { capture: true })` — `scroll` doesn't bubble, the same
  documented reason `ui/gsxui.js`'s own header comment gives for
  `toggle`/`close`/`focus`/`blur`, and rAF-throttled since `scroll` fires
  far faster than layout needs to be re-measured.
- MECHANISM (keyboard, ArrowLeft/ArrowRight regardless of orientation):
  `ui/carousel.js` binds `ArrowLeft`/`ArrowRight` to prev/next within
  `[data-gsxui-carousel]`, mirroring shadcn's own `onKeyDownCapture` exactly
  — including the fact that it is hard-coded to `ArrowLeft`/`ArrowRight`
  UNCONDITIONALLY, never `ArrowUp`/`ArrowDown`, even for
  `orientation="vertical"` carousels (verified against `carousel.tsx`'s own
  `handleKeyDown`, which never branches on orientation).
- GAP (late-added carousels, module-init scan only): the module-level init
  loop (`gsxuiCarousel` handle attachment, initial disabled/index
  recompute, `ResizeObserver` registration, autoplay start) runs once over
  `document.querySelectorAll("[data-gsxui-carousel]")` at module load — the
  same one-time-scan shape `toggle-group.js`'s `normalize()` loop and
  `command.js`'s `filter()` loop already establish. A carousel added later
  via an HTMX swap after this module has already run is not picked up; the
  same accepted limitation those two modules' own init loops carry.
- Registry: `carousel.gsx` imports `ui/icon` (`CarouselPrevious`/
  `CarouselNext`'s `icon.ArrowLeft`/`icon.ArrowRight`) and composes `Button`
  — `registry.Deps("carousel") == ["button", "icon"]`, pinned in
  `internal/registry/registry_test.go`. `HasJS("carousel")` is `true` — real
  new interactive JS (`ui/carousel.js`), unlike `sheet`/`alert-dialog`/
  `drawer`'s own JS-free reuse of `ui/dialog.js`.

## input-otp
- MECHANISM (stated plainly: ONE real input, not N inputs): the entire
  behavior contract is that there is exactly ONE actual `<input>` element
  (`data-gsxui-input-otp-input`), absolutely positioned to cover the whole
  slots row and made visually invisible with `opacity-0` (never
  `sr-only`/`hidden`/`display:none`, which would break focusability) while
  staying clickable (no `pointer-events-none`, kept `z-10` on top). It owns
  focus, caret, and paste. `InputOTPSlot` renders a purely presentational,
  empty `<div>` — every slot's character, `data-active` state, and
  fake-caret visibility are computed entirely by `ui/input-otp.js` from the
  real input's live `value` + `selectionStart`/`selectionEnd` on every
  `input`/`selectionchange`/`focus`/`blur` event. This is **not** N
  separate `<input>` elements with manual per-slot focus-advance-on-keypress
  JS — the common naive OTP reimplementation this port deliberately avoids.
  Native browser text editing (typing advances the cursor, Backspace
  deletes-and-moves-back, arrow keys move the selection) supplies
  focus-advance, backspace-across-slots, and arrow navigation entirely for
  free — `ui/input-otp.js` has zero per-slot keyboard handling. Paste needs
  no special-case code path either: it lands in the real input's `value`
  like any paste and is picked up by the ordinary `input`-event recompute
  (verified: paste fires a native `input` event with
  `inputType: "insertFromPaste"`, which the delegated `input` listener
  handles identically to typing, pattern filter included).
- ADAPT (DOM-order index stamping — a deliberate API departure from
  shadcn's own `InputOTPSlot(index)`): shadcn's `InputOTPSlot` takes an
  explicit `index` prop keying into React context's `slots[index]`; gsx has
  no equivalent shared context, and threading an explicit index through
  every call site is error-prone across `InputOTPSeparator`-split groups (a
  3-group layout needs indices `0,1 | 2,3 | 4,5` spanning group boundaries
  — easy to get wrong by hand if slots are ever reordered). `InputOTPSlot`
  takes **no index param at all**; `ui/input-otp.js` walks each root's
  slots (`querySelectorAll('[data-slot="input-otp-slot"]')`) in DOM order
  and stamps `data-index` positionally, idempotently, on first use (init
  scan, or the first `input`/click event to reach an as-yet-unstamped
  root) — the exact same "stamp source order once, JS-computed identity
  thereafter" idiom `command.js` already establishes for its own
  `data-gsxui-index` (`ui/command.js`'s `filter()`). Recompute itself
  doesn't strictly need the stamp (array position already equals slot
  index within a stable DOM), but the click-to-position handler does need
  a stable per-element identity to resolve "which index was clicked" — the
  stamp is the mechanism for that. This changes `InputOTPSlot`'s call
  signature from shadcn's own (`ui.InputOTPSlot()`, not
  `ui.InputOTPSlot(0)`) — a real, deliberate choice, not an oversight.
- DECISION (`data-gsxui-input-otp-pattern` is an UNANCHORED, per-character
  class — NOT shadcn's whole-string-anchored `REGEXP_ONLY_*` constants):
  `ui/input-otp.js` filters keystrokes (and pastes) by testing ONE
  character at a time against this pattern
  (`[...input.value].filter(c => re.test(c))`). shadcn's own exported
  patterns (e.g. `REGEXP_ONLY_DIGITS_AND_CHARS`,
  `^[0-9A-Za-z]*$`-shaped) are anchored to the WHOLE string for a different
  internal use inside the `input-otp` library; testing a single character
  against a `*`-quantified whole-string anchor would reject every
  non-empty character (a one-char string never satisfies `^...*$` the way
  one might assume from a naive per-character loop). `InputOTP` therefore
  documents — in `ui/input-otp.gsx`'s own doc comment and in
  `site/examples/inputotp/pattern.gsx` — that callers must pass a bare
  per-character class such as `[0-9]`, not `[0-9]*` or `^[0-9]*$`. The
  native HTML5 `pattern` attribute (whole-string-anchored, browser-native
  constraint validation) is a SEPARATE, independently-set attribute that
  still falls through via `attrs` unchanged — the pattern example sets
  both (`pattern="[0-9]*"` for native validity, `data-gsxui-input-otp-
  pattern="[0-9]"` for the live per-character filter) to make the split
  explicit. The pattern RegExp is constructed once per real input (cached
  in a `WeakMap`, not reconstructed every keystroke) and an invalid pattern
  source is caught so it degrades to "no filtering" instead of throwing on
  every `input` event.
- GAP (slots render EMPTY server-side, a first-paint divergence from
  shadcn): gsx has no client-side React context to pre-populate a slot's
  `char`/`isActive` from an initial `value` the way shadcn's SSR-capable
  React version can. Every `InputOTPSlot` renders inert
  (`data-active="false"`, no text, no caret) at server-render time
  regardless of the real input's initial `value` attribute — even a
  server-rendered pre-filled OTP (an initial-`value`-set adaptation of
  shadcn's `input-otp-controlled.tsx`, which itself doesn't translate
  directly since gsx has no two-way client↔server value binding without a
  page reload or HTMX round trip) would show empty slots for one frame
  until `ui/input-otp.js`'s own init scan populates every root present at
  parse time from the real input's actual value. This is a real, ledgered
  divergence, not an oversight — `TestInputOTPSlotNoCaretServerSide` in
  `ui/input-otp_test.go` pins the empty server render explicitly.
- ADAPT (does NOT compose `ui.Input`): the per-slot visual chrome (borders,
  `first:`/`last:` corner rounding, `data-[active=true]` ring) is different
  enough from `Input`'s own single-box recipe that composing would fight
  rather than help — the same reasoning `## sheet` gives for not composing
  `DialogContent`.
- Nova retarget (2026-07-24 controls source map, `## input-otp`, nova
  deltas table, `style-nova.css` lines 682-701): slot `size-8` (was
  new-york-v4's `h-9 w-9`), `shadow-xs` dropped (nova's shadow-presence
  drop, consistent with `input`/`checkbox`'s own nova entries),
  `first:rounded-l-lg`/`last:rounded-r-lg` (was `-md`, matches the global
  radius-scale bump), `ring-3` (was `ring-[3px]`, same value, Tailwind v4
  token spelling only). `data-[active=true]:z-10` is KEPT regardless of
  nova's own excerpt omitting it — functionally necessary (keeps the
  active ring from being occluded by a neighboring slot's border), judged
  a reference-sheet omission rather than a deliberate nova drop.
  `InputOTPGroup` adopts nova's NEW group-level `has-aria-invalid` ring
  block (`has-aria-invalid:ring-destructive/20
  dark:has-aria-invalid:ring-destructive/40 has-aria-invalid:border-
  destructive rounded-lg has-aria-invalid:ring-3`), entirely absent from
  new-york-v4's own `InputOTPGroup` (which only ever styles per-slot
  `aria-invalid:border-destructive`) — adopted as precedent-consistent
  (the destructive-ring-on-invalid color pattern already exists on
  `checkbox`/`radio`/`input`'s own `aria-invalid:ring-destructive*`), not a
  novel color introduction. `InputOTPSeparator` carries nova's
  `[&_svg:not([class*='size-'])]:size-4` safeguard, likely a no-op since
  `icon.Minus` already defaults to `size-4`, kept regardless per the map.
- `animate-caret-blink` needs no new CSS: `tw-animate-css` (already
  imported by both `assets/gsxui.css` and `web/site.css` via `@import
  "tw-animate-css"`) defines `--animate-caret-blink: caret-blink 1.25s
  ease-out infinite` plus the `@keyframes caret-blink` rule inside its own
  `@theme inline` block, and Tailwind v4's `--animate-*` theme namespace
  auto-generates the `animate-caret-blink` utility from that variable —
  the same mechanism that already makes `animate-pulse`/`animate-spin`
  work elsewhere in this codebase with no hand-authored keyframes. Verified
  by grep against `node_modules/tw-animate-css/dist/tw-animate.css`; no
  edit to either CSS file was needed.
- Registry: `input-otp.gsx` imports `ui/icon` (`InputOTPSeparator`'s
  `icon.Minus`) — `registry.Deps("input-otp") == ["icon"]`, pinned in
  `internal/registry/registry_test.go`. `HasJS("input-otp")` is `true`.

## sonner
- MECHANISM (the codebase's FIRST client-constructed-DOM module — stated
  plainly): every other gsxui behavior (`dialog`/`dropdown`/`command`/
  `carousel`/`input-otp`/…) attaches delegated behavior to server-rendered
  markup already in the DOM. Sonner has **no server markup for the toast
  itself** — a toast is definitionally a client-triggered response to a JS
  event (a fetch resolving, a form saving), so `ui/sonner.js` **constructs**
  each toast `<li>` from scratch (`document.createElement` tree) on every
  `toast()` call, appends it into the one static `<ol data-gsxui-toaster>`
  region `ui.Toaster` server-renders, and owns its whole lifecycle: mount →
  stacking → dismiss timer → animated remove. This is the single largest
  new-architecture item in the Tier 3 batch; the server surface is *only*
  the empty positioned region (`ui.Toaster`), pinned by `TestToasterPinned`
  in `ui/sonner_test.go` (`command.js`'s precedent: Go-pin only the server-
  rendered surface, exercise the JS via examples + a browser pass).
- SYNTHESIZED CLASSES (no Tailwind source upstream — the map IS the spec):
  shadcn's own `sonner.tsx` renders nothing but a re-themed `<Sonner>`
  passthrough; the toast library owns 100% of the toast DOM and ships its
  look from a **non-Tailwind** stylesheet (`sonner/dist/styles.css`) re-
  themed through four CSS custom properties. There is therefore no class
  string to transcribe — the toast card's classes in `ui/sonner.js`
  (`rounded-2xl border border-border bg-popover w-[356px] text-popover-
  foreground shadow-lg`, type-tinted icons via `data-[type=…]:[&>[data-
  icon]]:text-{emerald,sky,amber,destructive}`) are the synthesized spec
  from the wrapped source map's `## sonner`, reconstructed to match our
  popover/card surfaces. This is a genuine simplification (WIN): none of
  sonner's `--normal-bg`/`--normal-text`/`--border-radius` indirection is
  reproduced — the surface tokens are hardcoded Tailwind classes.
- DECISION (type-tinted icons are a house choice, not new-york-v4's literal
  default): shadcn's own `sonner.tsx` does NOT pass `richColors`, so strict
  new-york-v4 parity would render monochrome icons. The tinted convention
  (emerald/sky/amber/destructive) tracks ui.shadcn.com's live nova site and
  is more legible; a monochrome `text-foreground` alternative is equally
  defensible. Chosen deliberately per the brief, ledgered here.
- MAINTENANCE SEAM (icon paths hand-copied into JS): the toast `<li>` is
  built by client JS, so `icon.CircleCheck` (a server-side Go call) is
  unreachable. The five type glyphs (`circle-check`/`info`/`triangle-alert`/
  `octagon-x`/`loader-circle`) plus the close `x` are duplicated as literal
  SVG-path strings in `ui/sonner.js`'s `GLYPHS`/`X_GLYPH` (with a provenance
  comment naming `ui/icon/icon_data.go`) — the same "static data ported into
  a JS module" precedent as `command.js`'s `commandScore`. If `ui/icon`'s
  glyphs are ever regenerated from a newer Lucide, these copies do NOT
  update automatically.
- BARREL-EXPORT PRECEDENT (first public imperative API): `ui/index.js` gains
  both `import "./sonner.js"` (side-effect: the declarative trigger + the
  window global) AND `export { toast } from "./sonner.js"`, so
  `import { toast } from "gsxui"` works for page authors. Every prior module
  only exported internals for sibling-module use (`dialog.js`'s
  `requestClose` for `command.js`); `toast` is the first public
  imperative surface re-exported through the barrel. Because an inline
  `<script>` on a demo page cannot import the barrel, `ui/sonner.js` also
  attaches `window.gsxui = Object.assign(window.gsxui ?? {}, { toast })` —
  `site/examples/sonner/types.gsx`'s promise button uses
  `window.gsxui.toast.promise` (documented choice).
- ADAPT (declarative trigger for zero-JS pages): mirroring the
  `data-gsxui-dialog-trigger` idiom, a click-delegated
  `data-gsxui-toast="message"` (+ `-description`, `-type`, `-action` for the
  label) on any element calls the same internal `show()` the imperative API
  uses — the docs demos need no page-specific `<script>`. The constructed
  toast `<li>` ALSO carries `data-gsxui-toast` (an empty slot marker, per
  the map's markup), so the delegated handler guards on a non-empty value —
  that guard is what stops a click INSIDE a toast from spawning a blank one.
- STACKING/TIMERS (fresh design, no source to port): a plain array of toast
  records (`{id, el, type, duration, remaining, timer, startedAt}`), NOT
  sonner's own CSS-custom-property machine (we ship fixed Tailwind classes,
  not a themeable third-party stylesheet). At most `MAX_VISIBLE=3` show
  collapsed; each older toast peeks `COLLAPSE_PEEK=16px` up and shrinks
  `SCALE_STEP=0.05` per level via an inline `transform` recomputed per
  layout (`origin-bottom` keeps bottoms aligned, descending `z-index` paints
  the front on top); a 4th+ toast is hidden until a dismiss promotes it.
  Hovering any toast expands the whole stack (toasts separate upward by
  their measured `offsetHeight` + `EXPAND_GAP=14px`) AND pauses every timer
  (remaining-time resume on leave, debounced `HOVER_LEAVE_MS=80` so crossing
  a gap doesn't collapse). Default duration `4000ms`; a `loading` toast has
  `Infinity` until its promise settles. Enter/exit follow the house
  discrete-transition idiom adapted to a JS-constructed node (set the closed
  visual state, force one frame, flip to open; exit reverses and removes on
  `transitionend`, capped at `600ms` — the same race-against-a-cap idea as
  `dialog.js`'s `requestClose`, for backgrounded tabs with a frozen
  transition clock). `toast.promise` MORPHS THE SAME NODE in place on settle
  (swap icon/type/title, restart the timer, no re-animation, no stack
  reflow) — a naive dismiss-old/spawn-new would visibly jump the stack.
- GAP (positions): v1 ships only the default `bottom-right` region. The
  other five sonner positions (`top-left`/`top-center`/`top-right`/
  `bottom-left`/`bottom-center`) are a real feature reduction, ledgered —
  none of the docs demos exercise anything but the default.
- GAP (swipe-dismiss): sonner is touch swipe-dismissible (drag past a
  threshold). gsxui v1 ships no gesture layer — dismissal is the close
  button, the action button, the auto-dismiss timer, or `toast.dismiss(id?)`.
- Registry: `sonner.gsx` (Toaster) is a plain `<section>`/`<ol>` — no
  `ui/icon` import (icons live as SVG-path strings in `ui/sonner.js`) and no
  intra-package component reference, so `registry.Deps("sonner")` is empty,
  pinned in `internal/registry/registry_test.go`. `HasJS("sonner")` is
  `true` (`ui/sonner.js`, exact-basename match).
