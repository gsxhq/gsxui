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
    document.addEventListener(type, (event) => dispatch(handlers, event), capture);
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
