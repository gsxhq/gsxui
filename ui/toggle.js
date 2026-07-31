// Toggle behavior: click flips aria-pressed + data-state; gsxui:change on
// the toggle itself with { pressed } — same event shape as tabs.js's own
// gsxui:change, house style for state-flip components (see also dialog.js's
// gsxui:open/gsxui:close).
import { on, emit } from "./gsxui.js";

// ToggleGroupItem stamps data-gsxui-slot-toggle too — purely to share
// Toggle's CSS — while toggle-group.js owns its behavior; exclude it here so
// one click never fires both modules' handlers (the delegation registry
// dispatches to every matching selector).
on("click", "[data-gsxui-slot-toggle]:not([data-gsxui-slot-toggle-group-item])", (_event, toggle) => {
  const pressed = toggle.getAttribute("aria-pressed") !== "true";
  toggle.setAttribute("aria-pressed", pressed ? "true" : "false");
  toggle.dataset.state = pressed ? "on" : "off";
  emit(toggle, "gsxui:change", { pressed });
});
