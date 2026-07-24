# dev/

Buildless smoke pages: `python3 -m http.server -d .. 8080` then open
http://localhost:8080/dev/preview.html — no bundler, no Tailwind; verifies
delegation/behavior JS only.

Checklist per interactive component: open via trigger, close via X, close
via Esc, close via backdrop click; every transition logs its CustomEvent.
Programmatic path: dialog.showModal()/close() from the console must also
log the events. Click inside the panel's padding: must NOT close; click the
backdrop: must close.

After opening: dialog must have aria-labelledby and aria-describedby;
trigger aria-expanded flips true/false with open state.

Dropdown: click the trigger to open (positioned under it); trigger
aria-expanded flips true/false with open state; ArrowDown/ArrowUp roves
focus between menu items, wrapping at both ends; clicking a menu item logs
`gsxui:select <item text>` and closes the menu; clicking anywhere outside
the open menu light-dismisses it (no JS for this — native `popover="auto"`)
and logs `gsxui:close`; Esc also closes it for free. Every open/close
transition logs its CustomEvent.

Tooltip: hovering the trigger and holding still for ~300ms shows the
tooltip positioned above it; moving the pointer off the trigger hides it
immediately — this real-pointer path needs the controller's visual
verification (`pointerover`/`pointerout` don't fire from synthetic
`dispatchEvent` in the same reliable way `click` does). Synthetic fallback:
Tab to focus the trigger fires `focusin` and shows the tooltip with no
delay; Tab away (or Shift+Tab back) fires `focusout` and hides it — both
verifiable via `gsxui:open`/`gsxui:close` in the log without a real mouse.

Avatar: the good-image block must end up showing only the image (fallback
`display:none`); the broken-image block must end up showing only the
fallback (image `display:none`). The image covers the fallback via
`absolute inset-0`, so no-JS/pre-load rendering is correct; avatar.js
handles only the error path (hide broken image).

checkbox/radio/switch/native-select are zero-JS (native `<input>`/`<select>`
+ browser truth — `:checked`/`:disabled` for the former three, native
picker/`:disabled` for native-select — no delegation/behavior JS to smoke-test) so
they are intentionally absent from this harness — it verifies buildless
JS-behavior only, and these components' correctness is entirely in their
Tailwind class recipes, which need a Tailwind'd site to render meaningfully
(see Plan 4).

Tabs: click Tab B — Tab A must go `data-state="inactive"`/`aria-selected=
"false"`/`tabindex="-1"`, Content A must gain `hidden`, Tab B/Content B
must flip to active/visible; the log must gain one `gsxui:change
{"value":"b"}` line. Focus a trigger and press ArrowRight/ArrowLeft:
focus AND activation both move, wrapping at the ends; each press logs
another `gsxui:change`.

Accordion is zero-JS (native `<details name>` grouping — no
delegation/behavior JS to smoke-test), included here anyway as a markup
sanity check: opening Item 2 by clicking its summary must close Item 1
(same `name` attribute on both `<details>`), with no script involved at
all — reload with JS disabled and the exclusivity still holds.
