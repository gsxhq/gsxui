// Avatar behavior: native image load/error drives which of image/fallback
// shows. error/load don't bubble — capture-delegated.
import { on } from "../core/gsxui.js";

const sync = (img, ok) => {
  img.style.display = ok ? "" : "none";
  const fallback = img
    .closest('[data-slot="avatar"]')
    ?.querySelector('[data-slot="avatar-fallback"]');
  if (fallback) fallback.style.display = ok ? "none" : "";
};

on("error", "[data-gsxui-avatar-image]", (_e, img) => sync(img, false), { capture: true });
on("load", "[data-gsxui-avatar-image]", (_e, img) => sync(img, true), { capture: true });
