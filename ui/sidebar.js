// Sidebar behavior reflects state only. Presentation belongs to CSS, while
// persistence remains caller-owned through the gsxui:change event.
import { on, emit } from "./gsxui.js";

const wrapperOf = (element) => element.closest("[data-gsxui-slot-sidebar-wrapper]");
const desktopOf = (wrapper) => wrapper.querySelector("[data-gsxui-slot-sidebar-desktop]");
const mobileDialogOf = (wrapper) =>
  wrapper.querySelector("[data-gsxui-slot-sidebar-mobile-content]");

function isMobile(desktop) {
  return !desktop || getComputedStyle(desktop).display === "none";
}

function toggleDesktop(wrapper, desktop) {
  const open = wrapper.dataset.state !== "expanded";
  wrapper.dataset.state = open ? "expanded" : "collapsed";
  desktop.dataset.state = wrapper.dataset.state;
  desktop.dataset.collapsible = open
    ? ""
    : desktop.dataset.gsxuiSidebarCollapsible || "offcanvas";
  emit(wrapper, "gsxui:change", { open });
}

function openMobile(wrapper) {
  const dialog = mobileDialogOf(wrapper);
  if (!dialog || dialog.open) return;
  dialog.dataset.state = "open";
  dialog.showModal();
}

function toggle(wrapper) {
  const desktop = desktopOf(wrapper);
  if (isMobile(desktop)) {
    openMobile(wrapper);
  } else {
    toggleDesktop(wrapper, desktop);
  }
}

on("click", "[data-gsxui-slot-sidebar-trigger]", (_event, trigger) => {
  const wrapper = wrapperOf(trigger);
  if (wrapper) toggle(wrapper);
});

on("click", "[data-gsxui-slot-sidebar-rail]", (_event, rail) => {
  const wrapper = wrapperOf(rail);
  if (wrapper) toggle(wrapper);
});

function isTypingTarget(target) {
  if (!(target instanceof Element)) return false;
  if (target.isContentEditable) return true;
  return /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName);
}

addEventListener("keydown", (event) => {
  if (
    event.key.toLowerCase() !== "b" ||
    !(event.metaKey || event.ctrlKey) ||
    event.repeat ||
    isTypingTarget(event.target)
  ) {
    return;
  }
  const wrapper = document.querySelector("[data-gsxui-slot-sidebar-wrapper]");
  if (!wrapper) return;
  event.preventDefault();
  toggle(wrapper);
});
