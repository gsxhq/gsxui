# gsxui

shadcn-style components for [gsx](https://github.com/gsxhq/gsx) — copy-in,
type-checked, server-rendered.

Browse every component, live, at [ui.gsxhq.dev](https://ui.gsxhq.dev),
along with the theme editor.

**Status: pre-release.** 53 components + icon, covering most of shadcn/ui's
set.

## Install

    go install github.com/gsxhq/gsxui/cmd/gsxui@latest

Then, from your Go module:

    gsxui init          # tokens css, js runtime, class merger, gsx.toml wiring
    gsxui add dialog    # vendors dialog + its deps (button)
    gsxui list          # what's available

`gsxui init` prints the CSS and JS entry points to import from your
bundler. It also adds the gsx toolchain and class merger to your `go.mod`.

## Use

Components land in your module as one flat `package ui`:

```go
import "yourmodule/ui"

component Actions() {
	<div class="flex gap-3">
		<ui.Button>Save</ui.Button>
		<ui.Button variant="outline">Cancel</ui.Button>
	</div>
}
```

Each component is a `<name>.gsx` source — JSX-style, named parameters,
fallthrough attrs — plus a `<name>.js` behavior file when it's
interactive. `icon` is the one exception, vendored as its own `ui/icon`
package (generated Lucide data).

You own the vendored code. `gsxui add` never touches a modified file
unless you pass `--overwrite`. After upgrading the `gsxui` binary, re-run
`gsxui add <name> --overwrite` to refresh components — that discards your
local edits to those files.

## Components

**Form controls:** button, button-group, calendar, checkbox, combobox, field,
input, input-group, input-otp, label, native-select, radio, select, slider,
switch, textarea, toggle, toggle-group

**Display:** alert, avatar, badge, card, empty, item, kbd, progress,
separator, skeleton, spinner, table

**Overlay:** alert-dialog, context-menu, dialog, drawer, dropdown,
hover-card, menubar, popover, sheet, sonner, tooltip

**Navigation:** accordion, breadcrumb, command, navigation-menu,
pagination, sidebar, tabs

**Layout:** aspect-ratio, carousel, collapsible, resizable, scroll-area

**Primitives:** icon (Lucide, generated — pulled in as a dependency of
other components, rarely added directly)

## Differences from shadcn/ui

Some components are native-first: checkbox, radio, switch, native-select
and accordion trade a slice of shadcn's Radix-driven behavior for a real
`<input>`/`<select>`/`<details>` element — zero client JS, browser-native
`:checked`/`:disabled`/exclusivity semantics. dropdown and tooltip trade
Radix's Portal for the native popover API.

Every divergence, with its rationale, is ledgered in
[`docs/jsx-parity.md`](docs/jsx-parity.md).

## Upgrading from v1.0

v1.0 vendored one package per component (`ui/button/button.gsx` as
`package button`); current versions vendor flat into `ui/<name>.gsx`, one
`package ui`. The two layouts don't coexist and there's no in-place
migration: delete your vendored `ui/` directory and re-run
`gsxui add <name>...` with a current binary.

## Contributing

    npm install                       # once
    npx playwright install chromium   # once, for the browser suite

    make check      # go tests + browser tests + fmt/syntax checks
    make site-dev   # showcase site with live reload

See [`docs/backlog.md`](docs/backlog.md) for what's deferred and
[`docs/component-roadmap.md`](docs/component-roadmap.md) for coverage
plans.
