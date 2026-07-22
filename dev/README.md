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

Avatar: the good-image block must end up showing only the image (fallback
`display:none`); the broken-image block must end up showing only the
fallback (image `display:none`). The image covers the fallback via
`absolute inset-0`, so no-JS/pre-load rendering is correct; avatar.js
handles only the error path (hide broken image).

checkbox/radio/switchctl are zero-JS (native `<input>` + `:checked`/
`:disabled` browser truth, no delegation/behavior JS to smoke-test) so they
are intentionally absent from this harness — it verifies buildless
JS-behavior only, and these components' correctness is entirely in their
Tailwind class recipes, which need a Tailwind'd site to render meaningfully
(see Plan 4).
