# dev/

Buildless smoke pages: `python3 -m http.server -d .. 8080` then open
http://localhost:8080/dev/preview.html — no bundler, no Tailwind; verifies
delegation/behavior JS only.

Checklist per interactive component: open via trigger, close via X, close
via Esc, close via backdrop click; every transition logs its CustomEvent.
