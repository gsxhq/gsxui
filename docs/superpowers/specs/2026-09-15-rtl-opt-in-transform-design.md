# Opt-in RTL transform — design

Closes #32 (physical direction classes remain in ~15 components while the RTL
page says logical). Mirrors shadcn/ui's RTL model: the registry stays
physical, and a project that opts in gets its vendored components rewritten
to logical classes at install time, with a migrate command for components
installed earlier.

## Decisions

- **Registry source stays physical, exactly as upstream.** No re-port, no
  edits to `registry/styles` or `registry/canonical` for direction. The port
  drift (a re-port of the pinned upstream changes 366 sheets today) is a
  separate issue and untouched here.
- **RTL is a project-level opt-in**, `"rtl": true` in `gsxui.json`, set by
  `gsxui init --rtl`. Presets stay visual and carry no RTL bit.
- **The transform is upstream's table**, carried over from
  `packages/shadcn/src/utils/transformers/transform-rtl.ts` at the pinned
  shadcn-ui commit. We do not invent mappings; when upstream's table changes
  we re-derive ours from it.
- **The docs site shows the truth.** Like upstream's `ui-rtl` build, stylegen
  emits a transformed copy of `ui/` for the site, and every RTL demo renders
  from it. The RTL page describes the opt-in model instead of claiming
  logical classes by default.
- **`side` values are not rewritten.** Upstream rewrites `side="left"` to
  `inline-start` for Base UI. gsxui's floating positioning already resolves
  placement per direction in JS (`isRTL(el)` at position time), and Sheet,
  Drawer and Sidebar `side` is physical by an earlier ruling. Nothing to map.

## API surface

### `gsxui.json`

```json
{ "ui": "ui", "js": "web/gsxui", "css": "web/gsxui/index.css", "rtl": true }
```

`Config` gains `RTL bool` with json tag `rtl,omitempty`. Absent means false.
`gsxui init --rtl` writes it; `gsxui init` without the flag on an existing
project keeps the current value rather than resetting it to false.

### `gsxui migrate rtl`

Transforms every managed `.gsx` under `cfg.UI` (recursively, so `ui/icon`
and `ui/i18n` are visited and left unchanged because they carry no classes)
and the composed `style.css` in place, then runs `gsx generate`. It edits
user-modified files too: that is its purpose. It refuses when `rtl` is not
set in `gsxui.json`, with the message to set it first, so the flag and the
tree cannot disagree. It is idempotent: running it twice equals running it
once, pinned by a test, because upstream's own migrate had a duplicate-class
bug (shadcn-ui/ui#9734).

`gsxui migrate` with no subcommand prints the available migrations. `rtl` is
the only one.

### `gsxui add` and `gsxui apply`

When `cfg.RTL` is true, every `.gsx` artifact in the plan and the composed
`style.css` pass through the transform before the artifact hash is recorded.
Managed-file detection keeps working because the recorded hash is of the
transformed content. `apply` re-vendors installed components through the
same path, so a style switch stays RTL. Nothing changes when the flag is
false; the e2e suite's existing byte-identical assertions pin that.

## The transform

Package `internal/rtl`, two entry points and no state:

```go
// GSX rewrites the class attributes of every element in a gsx source file.
func GSX(src []byte) ([]byte, error)

// CSS rewrites the utility lists of every @apply in a stylesheet.
func CSS(src []byte) ([]byte, error)

// Classes rewrites one whitespace-separated class list. It is the unit both
// entry points call and the unit the table tests exercise.
func Classes(list string) string
```

### Token rules (upstream's, in order)

For each class, split into `variant:` prefix (everything up to the last
top-level colon, brackets respected), `value`, and `/modifier`.

1. `translate-x-N` and `-translate-x-N`: keep, and append the sign-flipped
   form under `rtl:` with the same variant (`rtl:-translate-x-N`).
2. `space-x-N`, `divide-x-N`: keep, and append `rtl:space-x-reverse` /
   `rtl:divide-x-reverse` with the same variant.
3. `cursor-w-resize` / `cursor-e-resize`: keep, and append the swapped one
   under `rtl:`.
4. Slide animations inside `data-[side=inline-start]` / `inline-end`
   variants: `slide-in-from-right` → `slide-in-from-end` and the three
   siblings. gsxui emits no `inline-*` sides today; the rule is carried for
   table parity and costs nothing.
5. If the variant contains `data-[side=left]` or `data-[side=right]`, skip
   the positioning prefixes (`left-`, `right-`, `-left-`, `-right-`) and fall
   through to the direct table for everything else.
6. Direct table, first match by prefix, exact match for entries without a
   trailing dash so `border-r` never matches `border-ring`:
   `-ml-`→`-ms-`, `-mr-`→`-me-`, `ml-`→`ms-`, `mr-`→`me-`, `pl-`→`ps-`,
   `pr-`→`pe-`, `-left-`→`-start-`, `-right-`→`-end-`, `left-`→`start-`,
   `right-`→`end-`, `inset-l-`→`inset-inline-start-`,
   `inset-r-`→`inset-inline-end-`, `rounded-tl-`→`rounded-ss-`,
   `rounded-tr-`→`rounded-se-`, `rounded-bl-`→`rounded-es-`,
   `rounded-br-`→`rounded-ee-`, `rounded-l-`→`rounded-s-`,
   `rounded-r-`→`rounded-e-`, `border-l-`→`border-s-`, `border-r-`→`border-e-`,
   `border-l`→`border-s`, `border-r`→`border-e`, `text-left`→`text-start`,
   `text-right`→`text-end`, `scroll-ml-`→`scroll-ms-`, `scroll-mr-`→`scroll-me-`,
   `scroll-pl-`→`scroll-ps-`, `scroll-pr-`→`scroll-pe-`, `float-left`→`float-start`,
   `float-right`→`float-end`, `clear-left`→`clear-start`, `clear-right`→`clear-end`,
   `origin-top-left`→`origin-top-start`, `origin-top-right`→`origin-top-end`,
   `origin-bottom-left`→`origin-bottom-start`,
   `origin-bottom-right`→`origin-bottom-end`, `origin-left`→`origin-start`,
   `origin-right`→`origin-end`.

Companion classes (rules 1-3) are appended only when the exact class is not
already in the list. Mapped outputs are never inputs of any rule, so the
whole pass is a fixed point after one application.

Upstream's `cn-rtl-flip` marker is not carried: gsxui's directional icons
already carry `rtl:rotate-180` in source.

### Where classes live in gsx source

Two attribute shapes carry classes: `class="…"` (`StaticAttr`) and
`class={ … }` (`ExprAttr`). Inside the expression form, the class grammar is
a list of string literals and `"literal": condition` pairs, with recipe
accessor calls already desugared in vendored output. The transform rewrites
every string literal that is a list element or a pair key and never one that
is an operand of a condition (`active == "rtl"` stays). It reuses the
class-expression element walk `internal/stylegen/resolve.go` already
performs for accessor calls rather than a second parser, and edits by byte
offset so comments, formatting and every other literal are untouched, the
same discipline `internal/cli/rewrite.go` uses for import paths. The result
must round-trip through `gsx fmt` unchanged.

### `origin-*` utilities

Tailwind has no logical transform-origin. `assets/css/foundation.css` gains
six `@utility` rules, direction-switched with `:dir()`:

```css
@utility origin-top-start { &:dir(ltr) { transform-origin: top left } &:dir(rtl) { transform-origin: top right } }
```

and likewise `origin-top-end`, `origin-bottom-start`, `origin-bottom-end`,
`origin-start`, `origin-end`. They ship to every project (the flag only
decides whether anything references them) and the layer check registers them.

### tw-animate-css sized logical slides

The minified `tw-animate-css` build drops a descendant-combinator space in
`slide-in-from-start-*` and siblings, so the sized logical slides match only
elements that themselves carry `dir` (Wombosvideo/tw-animate-css#67). gsxui's
sized slides sit under physical `data-[side=…]` variants and stay physical
by rule 5, so the transform never emits a sized logical slide. The RTL page
records the upstream caveat for callers who write their own.

## The site

stylegen emits `site/uirtl/<component>.gsx` for every component in `ui/`,
each the transform of the default style's generated output with the package
clause rewritten to `uirtl`, through the same `--check` drift machinery as
`ui/` itself. `site/examples/rtl/login.gsx` and the per-component RTL
examples (`calendar/rtl.gsx`, `pagination/rtl.gsx`, `sidebar/rtl.gsx`) import
`uirtl` instead of `ui`. Behaviour JS is shared: the transformed markup keeps
every `data-gsxui-slot-*` marker and every `data-gsxui-*` attribute, so
`ui/*.js` binds to it unchanged.

The RTL page's "How it works" becomes:

- Set `rtl: true` in `gsxui.json` (`gsxui init --rtl`); `gsxui add` then
  rewrites physical direction classes to logical ones in what it vendors.
- An `rtl: true` project renders both directions: the document's `dir`
  decides at render time, so `dir="ltr"` looks exactly as before and
  `dir="rtl"` mirrors. Set `dir` from the locale as today.
- `gsxui migrate rtl` rewrites components you installed before opting in.
- Directional icons carry `rtl:rotate-180`; floating placement and keyboard
  semantics mirror in JS (the existing bullets, unchanged).
- Sheet, Drawer and Sidebar `side` stays physical.
- The demos on this page are built from the transformed components.

Plus a short "Not transformed" list: `side` values, behaviour JS, and the
tw-animate-css sized-slide caveat.

## Tests

- `internal/rtl`: table tests for `Classes` ported case by case from
  upstream's `transform-rtl.test.ts` at the pin (variants, modifiers,
  negative values, arbitrary values, the `border-ring` false positive, the
  physical-side skip, each companion rule); idempotency (`Classes(Classes(x))
  == Classes(x)` for every case and for every class list in `ui/`); `GSX`
  leaves a comparison operand alone, edits only class attributes, and
  round-trips `gsx fmt`; `CSS` edits only `@apply` lists.
- Whole-registry: transform every `ui/*.gsx`, write to a temp package, run
  `gsx generate` and `go build`; the count of physical classes in the output
  is zero outside the `data-[side=…]` skip.
- CLI: `TestAddVendorsWithDeps`-style test with `rtl: true` asserting the
  vendored dialog carries `end-2` not `right-2` and that a second `add
  --overwrite` is byte-identical; `migrate rtl` on a modified file preserves
  the modification and is idempotent; `migrate rtl` refuses without the flag;
  `apply` keeps RTL on a style switch.
- e2e: scratch module with `rtl: true`, `add` for the fifteen components in
  #32's table, assert none of the physical classes that table lists survive
  in the vendored files, and `go build` passes.
- Playwright: the RTL login demo gains a `NativeSelect` (a language field),
  and a test asserts its chevron sits at the logical end under `dir="rtl"`,
  the reporter's own symptom; the existing `rtl.spec.ts` stays unchanged.
- Gates: layer check registers the new `@utility` rules; `stylegen --check`
  covers `site/uirtl`; `make highlight` after example edits.

## Docs

- RTL page rewrite as above; `gsxui init --rtl` in the getting-started flags.
- `docs/jsx-parity.md`: a `## rtl` section stating the model, the table's
  provenance (file and pin), the `side` non-mapping and why, and the
  tw-animate caveat.
- CHANGELOG: `### Added` for the flag, the transform and `migrate rtl`;
  `### Changed` for the RTL page's corrected claim.
- README: one line under the vendoring paragraph.

## Out of scope

- Re-porting or hand-editing the recipe sheets to logical classes.
- Rewriting `side` prop values to `inline-start`/`inline-end`.
- The port drift between the committed sheets and a fresh port.
- A reverse migration (upstream has none).

## Addendum (2026-09-16, planning)

- **Side-keyed classes stay physical in full, not positioning-only.** In
  gsx, Sheet's and Drawer's side arms are `switch side { case "left": … }`
  value arms inside the class list, and Sidebar's rail border is keyed by an
  arbitrary variant spelled `[…[data-side=left]>&]:border-r`, not upstream's
  `data-[side=left]:` form. Upstream skips only the positioning prefixes
  inside a physical-side variant and still maps `border-r` to `border-e`,
  which would put a `side="left"` rail's border on the outer edge under
  RTL. gsxui therefore treats a token as side-keyed when its variant
  contains `data-[side=left]`, `data-[side=right]`, `data-side=left]`,
  `data-side=right]`, or Drawer's `vaul-drawer-direction=left]` /
  `=right]`, or when it sits in a value arm of a `switch side`
  or an `if side …` inside a class list (Drawer names the same parameter
  `direction`; both identifiers count), and leaves every such token
  unchanged: positioning, borders, radii and slides alike. This is the
  existing ruling that Sheet, Drawer and Sidebar `side` is physical, and
  `jstest/specs/rtl.spec.ts` already pins it.
- **Centering needs no exception.** `left-1/2 -translate-x-1/2` maps to
  `start-1/2 -translate-x-1/2 rtl:translate-x-1/2`, which centres under
  both directions. Upstream's translate companion rule handles it.
- **`SidebarInset`'s `ml-0`/`ml-2`** are keyed on `variant=inset`, not on
  `side`, so they map to `ms-0`/`ms-2` like upstream. That is right for the
  common RTL layout (a `side="right"` sidebar) and is upstream's behaviour.
- `internal/rtl` and `internal/stylegen` share one byte-offset edit
  applier, extracted to `internal/srcedit`, instead of a second copy.
