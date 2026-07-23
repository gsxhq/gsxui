// Site chrome behavior — dogfoods the same delegated-listener core the
// components themselves use (ui/core/gsxui.js), rather than a separate
// listener pattern just for site pages.
import { on } from "../ui/core/gsxui.js";

// Component pages wrap each example's source block as:
//   <div data-site-example><pre><code>…</code></pre><button data-site-copy>…</button></div>
// Clicking the copy button copies that block's code text to the clipboard.
on("click", "[data-site-copy]", (_event, el) => {
  const code = el.closest("[data-site-example]")?.querySelector("code");
  if (!code) return;
  navigator.clipboard.writeText(code.textContent ?? "");
});
