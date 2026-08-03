# gsxui

shadcn-style components for [gsx](https://github.com/gsxhq/gsx) — copy-in,
type-checked, server-rendered.

Browse every component, live, at [ui.gsxhq.dev](https://ui.gsxhq.dev),
along with the theme editor.

**Status: pre-release.** 53 components + icon, covering most of shadcn/ui's
set.

## Install

    go install github.com/gsxhq/gsxui/cmd/gsxui@latest

Create a GSX app, then initialize gsxui:

    gsx init app --yes
    cd app
    gsxui init          # Tailwind/Vite + tokens, JS runtime, class merger
    gsxui add dialog    # vendors dialog + its deps (button)
    gsxui list          # what's available

On the unmodified npm/Vite scaffold produced by `gsx init --yes`, `gsxui init`
installs Tailwind CSS, `@tailwindcss/vite`, and `tw-animate-css`; configures the
Vite plugin; and imports `web/gsxui/index.js` and
`web/gsxui/index.css` from `web/main.js`. It also adds the GSX toolchain and
class merger to your `go.mod`. The CSS entry composes the behavior-critical
foundation, semantic `theme.css`, and replaceable component `style.css`.
Recoloring a project means changing only `theme.css`.

For a custom Vite config, entry file, package manager, or gsxui JS/CSS path,
init refuses automatic rewriting before it changes the project. Integrate those
projects manually:

    npm install --save-dev tailwindcss@^4.3.3 @tailwindcss/vite@^4.3.3 tw-animate-css@^1.4.0
    # vite.config.ts: import tailwindcss and add tailwindcss() to plugins
    # web/main.js: import "./gsxui/index.js" and "./gsxui/index.css"

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
hover-card, menubar, popover, sheet, toast, toaster, tooltip

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

## Styling roles

Each semantic styling role is a bare presence attribute. Roles compose through
ordinary fallthrough attribute forwarding—no token-merging helper is needed:

```gsx
<ui.Button data-gsxui-slot-dialog-trigger>
	Open
</ui.Button>
```

The slot attribute carries both the styling role and the behavior/ARIA
contract. Project selectors use the exact presence selector, for example
`[data-gsxui-slot-button]`.

## Contributing

    npm install                       # once
    npx playwright install chromium   # once, for the browser suite

    make check      # go tests + browser tests + fmt/syntax checks
    make site-dev   # showcase site with live reload

See [`docs/backlog.md`](docs/backlog.md) for what's deferred and
[`docs/component-roadmap.md`](docs/component-roadmap.md) for coverage
plans.
