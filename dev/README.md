# dev/

Buildless smoke pages: `python3 -m http.server -d .. 8080` then open
http://localhost:8080/dev/preview.html — no bundler, no Tailwind; verifies
delegation/behavior JS only.

Checklist per interactive component: open via trigger, close via X, close
via Esc, close via backdrop click; every transition logs its CustomEvent.
Programmatic path: dialog.showModal()/close() from the console must also
log the events. Click inside the panel's padding: must NOT close; click the
backdrop: must close.
