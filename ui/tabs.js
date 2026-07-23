// Tabs behavior: click + roving arrow keys; state stamped on triggers and
// panels; gsxui:change on the root with { value }.
import { on, emit } from "./gsxui.js";

function activate(trigger) {
  const root = trigger.closest("[data-gsxui-tabs]");
  if (!root) return;
  const value = trigger.dataset.value;
  root.dataset.value = value;
  for (const t of root.querySelectorAll('[data-gsxui-tabs-trigger]')) {
    if (t.closest("[data-gsxui-tabs]") !== root) continue;
    const active = t.dataset.value === value;
    t.dataset.state = active ? "active" : "inactive";
    t.setAttribute("aria-selected", active ? "true" : "false");
    t.tabIndex = active ? 0 : -1;
  }
  for (const p of root.querySelectorAll('[role="tabpanel"]')) {
    if (p.closest("[data-gsxui-tabs]") !== root) continue;
    const active = p.dataset.value === value;
    p.dataset.state = active ? "active" : "inactive";
    p.hidden = !active;
  }
  emit(root, "gsxui:change", { value });
}

on("click", "[data-gsxui-tabs-trigger]", (_e, t) => activate(t));

on("keydown", "[data-gsxui-tabs-trigger]", (e, t) => {
  const dir = { ArrowRight: 1, ArrowLeft: -1 }[e.key];
  if (!dir) return;
  const tablist = t.closest('[role="tablist"]');
  if (!tablist) return;
  const list = [...tablist.querySelectorAll('[data-gsxui-tabs-trigger]')];
  const next = list[(list.indexOf(t) + dir + list.length) % list.length];
  next.focus();
  activate(next);
  e.preventDefault();
});
