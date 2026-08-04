// Avatar behavior: native image load/error drives which of image/fallback
// shows. error/load don't bubble — capture-delegated.
import { on, init } from "./gsxui.js";

const sync = (img, ok) => {
  img.style.display = ok ? "" : "none";
  const fallback = img
    .closest("[data-gsxui-slot-avatar]")
    ?.querySelector("[data-gsxui-slot-avatar-fallback]");
  if (fallback) fallback.style.display = ok ? "none" : "";
};

on("error", "[data-gsxui-slot-avatar-image]", (_e, img) => sync(img, false), { capture: true });
on("load", "[data-gsxui-slot-avatar-image]", (_e, img) => sync(img, true), { capture: true });

// Images that settled before this module's init ran already fired
// load/error — init() (self-healing: current matches, later-added matches,
// and matches whose subtree is morphed back to server state) covers the
// rest, plus a window-load sweep for images still loading at import time.
// Pure reflection (img.complete/naturalWidth are read-only browser state),
// so re-running on an already-synced image is a no-op — safe to re-run
// unguarded.
function sweep(img) {
  if (img.complete) sync(img, img.naturalWidth > 0);
}
init("[data-gsxui-slot-avatar-image]", sweep);
window.addEventListener(
  "load",
  () => {
    for (const img of document.querySelectorAll("[data-gsxui-slot-avatar-image]")) sweep(img);
  },
  { once: true },
);
