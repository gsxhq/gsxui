// gsxui delegation core. One document-level listener per (event type, phase);
// behaviors register selector→handler pairs. No per-instance listeners, no
// init scan — elements added later (HTMX swaps, innerHTML) just work.
//
// Non-bubbling events (toggle, close, focus, blur, …) must be registered with
// { capture: true } — the document-level listener only sees them during the
// capture descent.

const registry = new Map(); // "type:capture" -> [{ selector, handler }]

export function on(type, selector, handler, { capture = false } = {}) {
  const key = `${type}:${capture}`;
  let handlers = registry.get(key);
  if (!handlers) {
    handlers = [];
    registry.set(key, handlers);
    document.addEventListener(
      type,
      (event) => dispatch(handlers, event),
      capture,
    );
  }
  handlers.push({ selector, handler });
}

function dispatch(handlers, event) {
  const target = event.target instanceof Element ? event.target : null;
  if (!target) return;
  for (const { selector, handler } of handlers) {
    const el = target.closest(selector);
    if (el) handler(event, el);
  }
}

export function emit(el, name, detail) {
  return el.dispatchEvent(
    new CustomEvent(name, { bubbles: true, cancelable: true, detail }),
  );
}

// clampToViewport clamps a desired fixed-position (left, top) so a box of
// the element's own offset size stays inside the visual viewport. The
// siblings' documented no-clamp NOTE accepted edge imprecision as a stopgap
// until CSS anchor positioning is Baseline; the theme editor's pickers
// disproved the acceptance — an ordinary panel near the viewport edge
// rendered its popover flush offscreen-left. Shift-only (no side flip), the
// same semantics context-menu.js already ships. clientWidth/Height rather
// than innerWidth/Height: the inner metrics include classic-scrollbar
// gutters, and clamping against them tucks an edge under the scrollbar.
// Call AFTER showPopover() — a hidden popover has no layout box.
export function clampToViewport(content, left, top) {
  const maxLeft = document.documentElement.clientWidth - content.offsetWidth;
  const maxTop = document.documentElement.clientHeight - content.offsetHeight;
  content.style.left = `${Math.max(0, Math.min(left, maxLeft))}px`;
  content.style.top = `${Math.max(0, Math.min(top, maxTop))}px`;
}
