# Internal messages and calendar locale — design

Closes #31 (hard-coded English labels callers can't override) and #30
(Calendar's English month/weekday names, caption and day labels). Issue #32
(physical direction classes) is separate and out of scope here.

## Decisions

- **No context API, no label props, no struct of labels.** Upstream shadcn
  hardcodes these strings and relies on copy-in; its only locale surface is
  the one react-day-picker and sonner already expose. gsxui follows the gsx
  translation pattern instead (`gsx/docs/guide/patterns/i18n-keys.md`): a
  string-kinded type that renders as written until the consuming module
  registers a renderer for it in its own `gsx.toml`.
- **The English text is the message id** (gettext style). No key registry,
  no fallback table: an unregistered type renders its English value.
- **Calendar's locale data is a prop, not a message.** Month and weekday
  names, the caption and day-label patterns and the digit set are
  formatting policy that Go code composes; renderers fire only at
  interpolation. This mirrors react-day-picker's `locale`/`formatters`
  props and gsxui's own `weekStartsOn` precedent.
- **The pattern is documented as a first-class feature**: a site page, the
  vendored file's own header, a contributor rule, and a test that enforces
  the rule.

## API surface

### `ui.T`

One Go-only canonical file, `registry/canonical/i18n.gsx`, declaring:

```go
// T is a message a component writes itself: an sr-only label, an aria-label,
// a visible "Close". Its value is the English text and doubles as the
// message id. It renders as written until the consuming module registers a
// renderer for it, for example in gsx.toml:
//
//	[renderers]
//	"example.com/app/ui.T" = "example.com/app/i18n.Translate"
//
// with func Translate(ctx context.Context, m ui.T) string looking the
// English text up in the request's locale. Text the caller supplies
// (children, props, attrs) never passes through T; only text the caller
// cannot reach does.
type T string
```

Written in components as `{ T("Close") }` in content and
`aria-label={ T("Close") }` in attributes. A `<Button>` child is
`<Button>{ T("Close") }</Button>`.

`T` is exported because the consumer's renderer function must name it. No
functions accompany it.

### `ui.CalendarLocale`

A new `locale CalendarLocale` parameter on `Calendar`. The zero value
reproduces today's English output byte for byte.

```go
// CalendarLocale is the text Calendar composes itself. The zero value is
// English. Patterns substitute exactly {month}, {year}, {weekday} and {day};
// any other text is literal.
type CalendarLocale struct {
	Months        [12]string // wide names, January first
	Weekdays      [7]string  // wide names, Sunday first; day aria-labels
	WeekdaysShort [7]string  // header row; English is the first two letters
	Caption       string     // default "{month} {year}"
	DayLabel      string     // default "{weekday}, {month} {day}, {year}"
	Digits        string     // ten runes replacing 0-9; "" keeps ASCII
}
```

Rules:

- A partially filled struct is a caller error, not a fallback case: an empty
  `Months[i]`, `Weekdays[i]` or `WeekdaysShort[i]` renders empty. Only an
  all-zero struct means English. The docs say so.
- `Digits` applies to visible text only: day cells, the caption, day
  aria-labels, the year dropdown's option text. Never to `value`
  attributes, `data-date`, hidden inputs or `data-gsxui-calendar-*`.
- `Digits` with a rune count other than 0 or 10 renders ASCII and the Go
  test for `CalendarLocale` pins that.
- The month dropdown's option text uses `Months`; the year dropdown's uses
  `Digits`.

## Rendering pipeline

### Helper files in stylegen

`resolveAll` currently errors on a canonical file with no shape. New rule: a
canonical `.gsx` that declares no `component` is a helper. It is emitted to
the same three destinations as a component (`ui/`,
`registry/generated/<style>/`, `site/stylepreview/<style>/`) with only the
package rewrite, and skips shape lookup, accessor generation, recipes and
CSS. `checkHelperCalls` is not run on it. The authoring check treats it as
unmigrated, which already forbids hand-authored presentation, and it has
none.

Detection is syntactic: the gsx parser's component declarations for the
file are empty. No filename convention.

### Vendoring

`internal/registry` derives dependencies from identifiers in generated
`.x.go`, and `T` is a type declaration there, so every component that
writes `T(...)` gains `i18n` as a dependency automatically. `gsxui add
dialog` vendors `i18n.gsx` with no CLI change. `gsxui list` shows `i18n`
like any entry, and each dependent's line reads `requires i18n`.

`gsx generate` in the consumer project produces `i18n.x.go` from the
Go-only file; the gsx package-renderers pattern relies on the same
behaviour.

### Calendar data flow

The server writes one JSON attribute on the root, always present:

```
data-gsxui-calendar-locale='{"months":[...],"weekdays":[...],"caption":"{month} {year}","dayLabel":"...","digits":""}'
```

Only the parts calendar.js writes after navigation are included: wide
month names, wide weekday names, the two patterns and the digits. The
header row is server-only and is not included.

calendar.js deletes `MONTH_NAMES` and `WEEKDAY_NAMES` and reads the
attribute once per root. `captionText` and `ariaLabel` become a pattern
substitution twin of the Go side, and a `digits` twin maps `0-9`. Server
and client agree by construction because both read the same values.

### The toaster.js fallback region

`createFallbackRegion` runs only when no server `<Toaster/>` is mounted.
Its `aria-label="Notifications"` stays English. The docs name this as the
degraded mode; mounting `Toaster`, which the toaster docs already require,
gives the translated landmark because the server `<section>` writes
`aria-label={ T("Notifications") }`.

## The sweep

Every English string on an element or text node that a caller cannot reach
through children, props or `attrs`. From #31, verified against `f0818be`:

| Component | Text | Position |
|---|---|---|
| dialog | `Close` (sr-only span), `Close` (footer button) | content |
| sheet | `Close` (sr-only span) | content |
| toast | `Close` | `aria-label` |
| toaster | `Notifications` | `aria-label` on `<section>` |
| carousel | `Previous slide`, `Next slide` | content |
| sidebar | `Sidebar`, `Displays the mobile sidebar.`, `Toggle Sidebar` (sr-only) | content |
| breadcrumb | `More` | content |
| pagination | `Previous`, `Next`, `More pages` | content |
| calendar | `Previous month`, `Next month`, `Month`, `Year` | `aria-label` |

Not swept, with the reason recorded in the sweep test:

- Literals that precede `{ attrs... }` and so are caller-overridable
  (`Spinner`'s `aria-label="Loading"`, the `Breadcrumb` and `Pagination`
  nav labels, `PaginationPrevious`/`Next` aria-labels, `SidebarRail`,
  `Carousel`'s `aria-roledescription`). The sweep test asserts each is
  still overridable, so a future reorder cannot silently make one
  unreachable.
- Toaster template placeholder texts (`Title`, `Description`, `Action`,
  `Cancel`): always overwritten or removed on clone.
- `input-otp`'s pinned `dir="ltr"` and all `data-gsxui-*` markers: not
  text.

## Documentation

This pattern is a feature and is documented at four levels.

### Site page: `/i18n`

`site/pages/i18n.gsx`, registered next to `Rtl` in `pages.go`, linked from
the layout nav, the command palette and the RTL page. Sections, in the
site's docs style (state the fact once, one fact per bullet, no repetition
of a snippet):

- **Get started** — the `gsx.toml` line and a minimal `Translate`
  function with `ctx`. The snippet is under the docs snippet drift test.
- **How it works** — every string a component writes itself is a `ui.T`;
  the English text is the id; unregistered it renders as English; text
  you pass in is yours and never passes through `T`.
- **Calendar** — the `locale` prop, the four placeholders, digits, the
  all-or-nothing rule, and that calendar.js reads the same values from the
  root.
- **Finding leftover English** — the pseudo-localization approach: a
  renderer that wraps every message, then grep the rendered page.
- **Not translated** — the toaster.js fallback region, and why.

### The vendored file

`i18n.gsx`'s header comment is the pattern in full, because a consumer
reading their own `ui/` directory should not need the site. It carries the
`gsx.toml` example and the renderer signature.

### Contributor rule

`docs/jsx-parity.md` gains a short section, **Internal messages**: any
English string a caller cannot reach is written as `T(...)`; caller text
never is; the sweep test is the gate. The calendar entry moves from GAP to
ADAPT and names `CalendarLocale`. `docs/component-roadmap.md`'s calendar
row and the `direction` line update. `CHANGELOG.md` gets one entry per
issue under `### Added`. README's feature list gains one line pointing at
the page.

### gsx side

The gsx `i18n-keys.md` pattern doc gains a sentence noting a vendored
component library can expose one such type for its internal strings, with
gsxui as the example. Filed as a gsx issue from this work, not done in
this repo.

## Examples

- `site/examples/calendar/localized.gsx`: a `Calendar` with an Arabic
  `CalendarLocale` (native digits, wide names) rendered under `dir="rtl"`,
  listed on the calendar component page and linked from the i18n page.
- The i18n page shows no live translated component: the site registers no
  renderer for `ui.T`, and adding a translator to the site is out of
  scope.

## Testing

- **Registry**: `TestDeps` expectations gain `i18n` for each swept
  component; a new case asserts `Deps("i18n")` is empty and that
  `Components()` lists it.
- **CLI**: an add test vendors `dialog` into a scratch module and asserts
  `i18n.gsx` lands and `gsx generate` succeeds on the Go-only file.
- **Stylegen**: a helper-file case in the generate tests (emitted to the
  three destinations, no recipe or accessor output) and
  `stylegen --check` clean on the committed tree.
- **Sweep test, end to end** (in `internal/cli`, alongside the existing
  scratch-module e2e test): renderers are compiled in by `gsx generate`,
  so this runs in a scratch consumer module whose `gsx.toml` registers a
  marker translator for `ui.T`. It vendors every swept component, generates,
  renders each one, and asserts the output contains no bare occurrence of
  any string in the sweep table and that every marker-wrapped string is one
  of them. This is the reporter's own pseudo-localization gate.
- **Sweep test, source level** (`ui/i18n_test.go`): for each string in the
  sweep table, the canonical source contains it only as the argument of
  `T(`; and for each entry in the "not swept" list, rendering with an
  `attrs` override replaces the literal. Existing component tests pass
  unchanged, which pins the unregistered output as today's English.
- **Calendar Go**: `CalendarLocale` zero value equals today's output;
  Arabic digits fixture; a year-first Japanese caption pattern; a
  malformed `Digits` length falls back to ASCII; `Digits` never touches
  `data-date`, `value` or hidden inputs; the JSON attribute round-trips.
- **Calendar Playwright**: a localized fixture navigates forward and back
  and the agreement diff asserts caption, grid `aria-label` and every day
  button's `aria-label` equal the server-rendered strings for that month,
  the same mechanism `calendar.spec.ts` uses for the grid.
- **Docs**: snippet drift test covers the `gsx.toml` snippet; layercheck
  registration if the page adds any asset.
- Full gate list per project practice: build, tests, `stylegen --check`,
  `--check-authoring`, `make audit`, `make verify-generated`,
  `make verify-generated-styles`, `gofmt -l`, Playwright with the jstest
  config, `make highlight` after adding the example.

## Out of scope

- #32, physical direction classes.
- Translating the toaster.js fallback region.
- A translator for gsxui's own site.
- Messages with arguments. No swept string has one.
- Any label prop or context value on any component.

## Addendum (2026-09-14, planning)

- `CalendarLocale` gains `MonthsShort [12]string`. The dropdown caption's
  month `<select>` renders three-letter names today (`calendarMonthNames`,
  matching upstream's `formatMonthDropdown`), so a wide-only struct could
  not reproduce the zero-value output byte for byte. `Months` stays the
  wide set for the caption and day labels.
- `registry.Components()` excludes helper files and a new
  `registry.Helpers()` lists them. The site's component index, sidebar and
  command palette all iterate `Components()` and link to
  `/components/<name>`, which 404s for a name with no examples, so a
  helper must not appear there. `Deps`, `Resolve`, `HasJS` and `gsxui add`
  accept both sets; `gsxui list` prints helpers with a `(helper)` tag.
- Patterns substitute all four placeholders on every call, with the
  absent parts as empty strings: a caption pattern containing
  `{weekday}` renders it as nothing on both the server and the client.

## Addendum (2026-09-15, Task 8)

- `T` lives in `ui/i18n`, a leaf directory package vendored like `ui/icon`,
  not in package `ui`. Components write `i18n.T("Close")`.
- The relocation is forced: gsx compiles a registered renderer in as a
  qualified call, so every generated `ui/<c>.x.go` imports the translator's
  package, and the translator must name `T`. A `T` in package `ui` puts the
  translator's package in an import cycle, which the end-to-end
  pseudo-localization gate caught.
- The recommended translator location is `ui/i18n/translate.go`, beside the
  vendored type. `gsxui add --overwrite` rewrites only the files it
  vendored, so that file survives re-vendoring.
- The stylegen helper-file pass-through added for the old
  `registry/canonical/i18n.gsx` is gone; `ui/i18n/i18n.go` is plain Go,
  hand-written and outside the style pipeline.
