---
name: gsxui
description: Use when building or restyling UI in a Go project that uses gsxui (shadcn-style gsx components vendored into a `ui` package, `gsxui.json` at the module root) — choosing components, composing pages and forms, styling with Tailwind classes, theming a brand, or running the gsxui CLI.
---

# Designing with gsxui

gsxui is shadcn/ui for gsx. `gsxui add` copies component source into your
module's `ui` package (icons go in `ui/icon`), and you own that code. Every
component reads its look from two places: semantic CSS variables in the theme
and the project's **style**, one of `vega`, `nova`, `maia`, `lyra`, `mira`,
`luma`, `sera` or `rhea`. **Get the look from the theme and the components.
Use classes for layout.**

For gsx syntax itself, use the `gsx` skill.

## CLI

```bash
gsxui init                        # in a Go module: gsxui.json, ui package, web/gsxui/*.css
gsxui list                        # components and what they require
gsxui add button card dialog      # vendor components (dependencies come along)
gsxui add button --diff           # compare the vendored copy with the current release
gsxui apply --preset brand.json   # switch style and theme
```

Run `gsxui help` or `gsxui <command> -h` for the rest. `add` never replaces a
file you have edited unless you pass `--overwrite`.

## Rules

**Compose components, don't rebuild them.** Before writing a styled `div`,
check `gsxui list`. A bordered box with a heading is a `Card`, a coloured
label is a `Badge`, a divider is a `Separator`, a loading block is a
`Skeleton`.

**Variants before classes.** Use `variant` and `size` first. Button has the
variants `default`, `secondary`, `outline`, `ghost`, `destructive` and `link`,
and the sizes `xs`, `sm`, `default` and `lg`, plus `icon`, `icon-xs`,
`icon-sm` and `icon-lg`.

```gsx
// Wrong
<ui.Button class="border bg-transparent text-red-600">Delete</ui.Button>
// Right
<ui.Button variant="destructive">Delete</ui.Button>
```

**Classes on a component are for layout.** Width, margin, grid and flex
placement are fine. A caller class replaces the component's conflicting
utility (`class="h-12"` beats the built-in height), so a colour or radius
class silently forks that one instance from the theme.

**Semantic tokens only.** Use `bg-background`, `text-foreground`, `bg-card`,
`bg-muted`, `text-muted-foreground`, `bg-primary`, `text-primary-foreground`,
`bg-accent`, `text-destructive`, `border-border` and `ring-ring`. Never use raw
palette colours (`bg-blue-500`, `text-gray-600`) and never use `dark:`
overrides. The tokens already switch under `.dark`.

**Spacing:** use `flex` or `grid` with `gap-*`, not `space-y-*` or margins
between siblings. Use `size-*` when width equals height.

**Links that look like buttons** use `<ui.Button href="/new">`, not a styled
`<a>`. Submit buttons need `type="submit"`; the default is `type="button"`.

**Icons** come from `ui/icon` (Lucide). Inside a Button or Badge the component
sizes them, so give no size class. An icon-only button needs `aria-label`.

```gsx
<ui.Button><icon.Plus/>New project</ui.Button>
<ui.Button variant="ghost" size="icon" aria-label="More"><icon.Ellipsis/></ui.Button>
```

**Forms** are built from `FieldSet`, `FieldLegend`, `FieldGroup`, `Field`,
`FieldLabel`, `FieldDescription` and `FieldError`. Don't stack bare labels and
inputs. Show an error with `data-invalid` on the `Field`, `aria-invalid` on the
control, and a `FieldError`.

```gsx
<ui.FieldGroup>
	<ui.Field data-invalid="true">
		<ui.FieldLabel for="email">Email</ui.FieldLabel>
		<ui.Input id="email" type="email" aria-invalid="true"/>
		<ui.FieldError>Enter a valid email address.</ui.FieldError>
	</ui.Field>
</ui.FieldGroup>
```

**Use the full composition.** A Card is `Card` › `CardHeader` (`CardTitle`,
`CardDescription`) › `CardContent` › `CardFooter`. Dialog, Sheet, Empty, Alert
and Item follow the same pattern. Don't put content directly inside the root.

## Choosing a component

| Need | Use |
|---|---|
| Focused task over the page | `Dialog` |
| Confirm a destructive action | `AlertDialog` |
| Side panel, settings, mobile nav | `Sheet` |
| Bottom panel on mobile | `Drawer` |
| Rich content anchored to a trigger | `Popover` |
| Short label for a control | `Tooltip` |
| Preview on hover | `HoverCard` |
| Action list from a button | `DropdownMenu` |
| Right-click actions | `ContextMenu` |
| Command palette, searchable list | `Command` |
| Pick one of a few options | `Radio`, `ToggleGroup` or `Select` |
| Pick from many, with search | `Combobox` |
| Pick one, native control | `NativeSelect` |
| On/off setting | `Switch`; `Checkbox` inside forms |
| Transient feedback | `Toaster` and `Toast` |
| Inline status message | `Alert` |
| Nothing to show yet | `Empty` |
| Loading | `Skeleton` for content, `Spinner` inside a button |
| List rows with media and actions | `Item` |
| App navigation | `Sidebar`; `NavigationMenu` for marketing sites |
| Grouped buttons | `ButtonGroup` |
| Text field with addons | `InputGroup` |
| Keyboard shortcut | `Kbd` |
| Tabular data | `Table` |

## Theming a brand

A theme only changes CSS variables. Component markup never changes for a
brand.

- **Colours, radius, fonts:** `web/gsxui/theme.css` holds `:root` and `.dark`
  blocks of shadcn-compatible variables. Set every token in both blocks, and
  keep each `*-foreground` readable on its pair (aim for WCAG AA).
- **Style:** density, borders, shadows and type scale come from the style. Pick
  one in the theme editor at <https://ui.gsxhq.dev/theme>, which exports a
  preset.
- **Apply a preset:** `gsxui apply --preset <file|code|URL|->` writes the theme
  and swaps the style. Add `--only theme` or `--only style` to change just one.

A preset is JSON:

```json
{
  "$schema": "https://ui.gsxhq.dev/schemas/preset-v1.json",
  "schemaVersion": 1,
  "style": "nova",
  "radius": "0.625rem",
  "fontSans": "\"Geist Variable\", ui-sans-serif, system-ui, sans-serif",
  "fontHeading": "\"Geist Variable\", ui-sans-serif, system-ui, sans-serif",
  "theme": {
    "light": { "background": "oklch(1 0 0)", "primary": "oklch(0.205 0 0)", "...": "..." },
    "dark": { "background": "oklch(0.145 0 0)", "primary": "oklch(0.922 0 0)", "...": "..." }
  }
}
```

To restyle one part everywhere, target its marker from your own CSS instead of
editing vendored files:

```css
[data-gsxui-slot-card] { box-shadow: none; }
```
