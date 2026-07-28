# Isolated Component Previews

## Problem

The component gallery renders every registered example inline inside the
documentation page. That works for ordinary components, but `Sidebar` is an
application-shell block: its desktop container is intentionally
`position: fixed` against the viewport. Rendering it inline makes the preview
escape its card and cover the documentation sidebar. Clipping the card cannot
contain fixed-position descendants.

Shadcn handles its Sidebar demo as a block preview. On desktop the docs page
loads the live block in an iframe, giving the fixed sidebar its own viewport;
on smaller screens it uses a captured image. Gsxui will use a live iframe at
all viewport sizes because the server-rendered component and its vanilla
JavaScript already work responsively without a client framework.

## Goals

- Keep application-shell examples live and interactive without letting their
  layout, dialogs, keyboard shortcuts, or CSS escape into the docs page.
- Keep the exact copyable GSX source on the parent component page.
- Make isolation a generic example presentation capability rather than a
  Sidebar-specific CSS exception.
- Preserve ordinary inline examples unchanged.
- Keep the isolated document on the same theme and asset bundle as the parent
  site.

## Non-goals

- Do not add an embedded or preview mode to the vendored `ui.Sidebar` API.
- Do not change Sidebar's production fixed-position semantics.
- Do not make every component example an iframe.
- Do not introduce screenshots or a second frontend bundle.
- Do not allow arbitrary source paths or unregistered nodes through the
  preview route.

## Architecture

### Example metadata

`examples.Example` gains an `Isolated bool` presentation field. Sidebar's
registered examples set it to true. The field is exported because the pages
package consumes the registry contract; it is not part of the vendored
component API.

The registry exposes an exact `(component, example)` lookup. It returns only
registered examples and preserves each example's optional query-driven render
hook.

### Isolated route

The site adds:

```text
GET /examples/{component}/{example}
```

The route resolves the two exact registry keys. Unknown pairs return 404. It
renders the example node in a minimal HTML document containing:

- the same theme bootstrap used by the main layout;
- the same Vite `web/main.js` entry, CSS, and preloads;
- a viewport-sized, overflow-contained body;
- no documentation header, navigation, footer, source block, or nested
  iframe.

The preview is same-origin, so component behavior and browser tests use the
same production assets without a special harness.

### Shared document head

The existing layout's document head becomes a package-private GSX component
used by both the full site layout and the isolated preview document. This
prevents theme/bootstrap/asset wiring from drifting between the two document
types.

### Parent presentation

`/components/{name}` continues to render ordinary examples inline. For an
isolated example it renders an iframe instead of `ex.Node`:

```html
<iframe
  data-site-isolated-preview
  title="Basic preview"
  src="/examples/sidebar/basic"
></iframe>
```

The iframe has a standard block-preview height, full width, a border, and a
background. Its document may scroll internally when an example is taller
than the preview viewport. The exact highlighted source and Copy control stay
immediately below the iframe.

### Theme synchronization

The isolated document reads the persisted theme before first paint. When the
parent page theme toggle changes during the current view, `web/site.js`
updates the `dark` class on each same-origin isolated document. An inaccessible
or not-yet-loaded frame is ignored; its own head bootstrap applies the stored
theme when it loads.

## Testing

### Go route tests

- `/components/sidebar` contains three isolated iframes and does not render a
  live `data-gsxui-slot-sidebar-container` in the parent document.
- Ordinary component pages still contain inline rendered component markers.
- `/examples/sidebar/basic` renders a minimal document with Sidebar markers
  and the Vite asset entry, without docs navigation or another iframe.
- Unknown component/example pairs return 404.

### Browser regression

- On `/components/sidebar`, every desktop Sidebar container belongs to an
  iframe document; none appears in the parent page.
- The parent docs sidebar retains its normal position and content.
- The Basic iframe's desktop Sidebar container is fixed to the iframe
  viewport, not the parent viewport.
- Its trigger changes state inside that iframe.
- Parent theme changes propagate to loaded Sidebar iframe documents.

## Success criteria

The live Sidebar page visually contains every example, the docs navigation
never gets covered, Sidebar keeps its real fixed app-shell behavior, copied
source remains unchanged, and the full authoritative CI gate passes.
