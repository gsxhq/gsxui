// Site chrome behavior — dogfoods the same delegated-listener core the
// components themselves use (ui/gsxui.js), rather than a separate
// listener pattern just for site pages.
import { on } from "../ui/gsxui.js";

// Component pages wrap each example's source block as:
//   <div data-site-example><pre><code>…</code></pre><button data-site-copy>…</button></div>
// Clicking the copy button copies that block's code text to the clipboard.
on("click", "[data-site-copy]", (_event, el) => {
  const code = el.closest("[data-site-example]")?.querySelector("code");
  if (!code) return;
  navigator.clipboard.writeText(code.textContent ?? "");
});

// Header theme toggle (shadcn's site model: one click flips the resolved
// theme, no system/menu step). The stored choice is what the layout's
// paint-blocking head script applies on the next load; clearing storage
// would fall back to the OS preference there.
on("click", "[data-site-theme-toggle]", () => {
  const dark = document.documentElement.classList.toggle("dark");
  try {
    localStorage.setItem("gsxui-theme", dark ? "dark" : "light");
  } catch {
    // storage unavailable (private mode etc.) — the toggle still works
    // for this page view, it just won't persist.
  }
});
