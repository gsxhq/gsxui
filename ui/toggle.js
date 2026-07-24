// Toggle behavior: click flips aria-pressed + data-state; gsxui:change on
// the toggle itself with { pressed } — same event shape as tabs.js's own
// gsxui:change, house style for state-flip components (see also dialog.js's
// gsxui:open/gsxui:close).
import { on, emit } from "./gsxui.js";

on("click", "[data-gsxui-toggle]", (_event, toggle) => {
  const pressed = toggle.getAttribute("aria-pressed") !== "true";
  toggle.setAttribute("aria-pressed", pressed ? "true" : "false");
  toggle.dataset.state = pressed ? "on" : "off";
  emit(toggle, "gsxui:change", { pressed });
});
