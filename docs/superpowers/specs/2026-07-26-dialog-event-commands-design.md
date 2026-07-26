# Dialog Event Commands

*2026-07-26 · Focused hardening of the existing native-dialog behavior.*

## Context

The original Phase 1 review findings have already been addressed on current
`main`:

- `DialogTrigger` server-renders its initial ARIA state and `dialog.js`
  lazily wires `aria-controls`, `aria-labelledby`, and `aria-describedby`.
- component-driven close paths stamp `data-state="closed"` and wait for the
  dialog's exit animations before calling native `close()`;
  `tw-animate-css` is included.
- the four original components have exact render pins.
- `make check` fails on generated `.x.go` drift.

The remaining hardening is architectural. gsxui components are intended to
compose with Alpine, HTMX, and plain DOM code. Commands should therefore cross
the component boundary as bubbling `CustomEvent`s, and component state must
remain inspectable in DOM attributes rather than private JavaScript state.

## Goals

1. Make dialog open and close requests public DOM-event commands.
2. Let Alpine, HTMX, and plain JavaScript dispatch or cancel those commands.
3. Keep `data-state` and native `<dialog open>` as the complete state model.
4. Remove the module-level ID counter and timer-based close cap.
5. Preserve the current accessibility, animation, rapid-reopen, AlertDialog,
   and post-transition notification behavior.

## Non-goals

- No `DialogFooter` API change.
- No framework dependency or framework-specific adapter.
- No interception or monkey-patching of direct `dialog.showModal()` or
  `dialog.close()` calls. Those remain native immediate escape hatches.
- No registry, initialization scan, `WeakMap`, `WeakSet`, or hidden per-instance
  JavaScript state.
- No renaming of the existing `gsxui:open` and `gsxui:close` notifications.

## Public event contract

Two new command events are handled on a dialog content element:

| Event | Meaning | Default action |
|---|---|---|
| `gsxui:open-dialog` | Command the dialog to enter its open state | stamp `data-state="open"`, wire ARIA, call `showModal()` when needed |
| `gsxui:close-dialog` | Command the dialog to enter its closed state | stamp `data-state="closed"`, wait for finite exit animations, call `close()` |

Both events are ordinary bubbling, cancelable `CustomEvent`s. `detail` is an
optional plain object. Built-in controls supply a stable `reason`:

- trigger: `{ reason: "trigger" }`
- injected or authored close control: `{ reason: "close-button" }`
- backdrop: `{ reason: "backdrop" }`
- native cancel/Escape: `{ reason: "escape" }`

Application code may supply any string or omit `reason`; behavior must not
branch on application-provided reason values.

The default action is deferred to a microtask after event propagation
finishes. This is load-bearing: an Alpine/HTMX/plain listener on an ancestor,
`document`, or `window` must be able to call `preventDefault()` even when that
listener was registered after gsxui's delegated listener.

After the native `toggle` event confirms the transition:

- `data-state` is reconciled with `event.newState`;
- every trigger in the owning dialog root receives the matching
  `aria-expanded`;
- `gsxui:open` or `gsxui:close` is emitted as the post-transition
  notification.

## Framework usage

The contract requires no special integration:

```html
<button
  x-data
  @click="$dispatch('gsxui:close-dialog', { reason: 'cancel' })"
>
  Cancel
</button>

<div hx-trigger="gsxui:close from:body">...</div>
```

```js
htmx.trigger(dialog, "gsxui:open-dialog", { reason: "server-response" })
htmx.trigger(dialog, "gsxui:close-dialog", { reason: "saved" })
```

Events dispatched from descendants are resolved to their closest owning
`dialog[data-gsxui-dialog-content]`, matching the existing delegated-listener
model.

## Target identity and native controls

Proximity remains the zero-configuration path: a
`data-gsxui-dialog-trigger` opens the dialog owned by its nearest
`[data-gsxui-dialog]` root. Resolution must exclude dialogs and triggers owned
by nested roots.

When a trigger lives outside that root, or application code controls one of
many dialogs, the dialog content must have an authored stable `id`. Native
Invoker Commands are the primary declarative form:

```html
<button commandfor="delete-dialog" command="show-modal">Open</button>
<dialog id="delete-dialog" data-gsxui-dialog-content>
  <button commandfor="delete-dialog" command="request-close">Cancel</button>
</dialog>
```

`show-modal` and `request-close` are standards-based commands targeted by
exact ID. The latter enters the dialog's cancel path, which gsxui translates
to the animated, cancelable `gsxui:close-dialog` command. Native
`command="close"` remains available when an immediate close is intended.

Programmatic helpers require the exact `HTMLDialogElement`; they never select
the nearest or first dialog:

```js
const dialog = document.getElementById("delete-dialog")
openDialog(dialog)
closeDialog(dialog)
```

Passing anything other than an `HTMLDialogElement` throws `TypeError`.
Generated accessibility IDs are not a public lookup contract: external
control requires an authored ID.

The native escape hatches remain fully supported and reconcile through
`beforetoggle`/`toggle`:

```js
dialog.showModal() // immediate native modal open
dialog.show()      // immediate native non-modal open
dialog.close()     // immediate native close
```

`beforetoggle` stamps the impending `data-state` before paint, including for
native `showModal()` and `command="show-modal"` paths. `toggle` remains the
confirmed-state notification source.

## DOM-only state and identity

`data-state` is the idempotency and transition guard:

- `open` means open or reopening;
- `closed` while the native dialog remains open means its exit transition is
  in progress;
- a second close request while `data-state="closed"` is a no-op;
- an open request during exit stamps `data-state="open"` and thereby cancels
  the pending close when its animation promise settles.

Accessibility identity is also stored only in the DOM. When an authored ID is
absent, `wireA11y` creates a prefixed ID from
`crypto.getRandomValues()` and writes it to the element's `id`.
`getRandomValues()` works in non-secure HTTP contexts where
`crypto.randomUUID()` does not, which matters for local-network HTMX apps. The
helper never overwrites authored IDs or authored
`aria-labelledby`, `aria-describedby`, or `aria-controls`.

There is no module-level counter or per-dialog JavaScript object.

## Animation completion

Closing waits only for animations attached to the dialog element itself, never
descendant animations. It filters those animations to effects whose computed
end time is finite, then awaits their `finished` promises with
`Promise.allSettled`.

- no finite animation: close immediately;
- canceled animation: its rejected `finished` promise settles normally;
- background tab: close completes when the real animation clock completes
  after resumption;
- caller-added infinite dialog animation: ignored;
- rapid reopen: the final `data-state` check aborts the pending close.

No timeout estimates animation duration.

## Internal interfaces

`dialog.js` renames its existing `requestClose` export to:

```js
closeDialog(dialog, detail?)
```

It now dispatches `gsxui:close-dialog` and returns the event's boolean
dispatch result. It does not close directly.

It additionally exports:

```js
openDialog(dialog, detail?)
```

with the equivalent `gsxui:open-dialog` behavior. Sibling modules such as
`command.js` use these exports rather than manipulating dialog state directly.
Both helpers require an `HTMLDialogElement` and are re-exported from
`ui/index.js` as the public programmatic API.

Private default-action functions may be asynchronous, but they retain no state
after their invocation settles.

## Tests

Use the existing Playwright layer against real Chromium:

1. trigger dispatches `gsxui:open-dialog` with reason `trigger`;
2. Alpine-style descendant dispatch opens and closes a dialog;
3. HTMX-style dispatch on the dialog opens and closes it;
4. a late ancestor/document listener can cancel either request;
5. close-button, backdrop, and Escape requests carry their stable reasons;
6. `data-state` is `closed` while the exit animation runs, then native
   `open` becomes false and one `gsxui:close` notification fires;
7. an `open-dialog` command during exit leaves the dialog open;
8. generated title, description, and dialog IDs are unique, stable in the
   DOM, and reflected by ARIA attributes;
9. authored IDs and authored ARIA relationships remain unchanged;
10. proximity ignores dialogs and triggers owned by nested dialog roots;
11. authored-ID `commandfor="…" command="show-modal"` opens the exact dialog;
12. authored-ID `commandfor="…" command="request-close"` uses the animated
    close path;
13. direct `showModal()`/`show()`/`close()` still reconcile state and
    notifications through `beforetoggle`/`toggle`;
14. helpers reject non-dialog targets with `TypeError`.

The existing Go render pins and generated-drift check remain unchanged because
their review findings are already fixed on `main`.

## Compatibility

- Standard DOM APIs only.
- No Alpine or HTMX runtime dependency in gsxui.
- Event names and detail objects work with Alpine `$dispatch`, HTMX
  `htmx.trigger`, `HX-Trigger`, `hx-trigger`, and native
  `dispatchEvent(new CustomEvent(...))`.
- Dynamically inserted HTMX content works without initialization because all
  command and native lifecycle handling stays document-delegated.
