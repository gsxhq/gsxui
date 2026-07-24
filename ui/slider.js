// Slider behavior: a single delegated `input` listener recomputes --fill
// from the range's own live min/max/value on every drag or keystroke.
// ui/slider.gsx already computes the correct --fill for FIRST paint
// server-side (zero JS needed just to render); this module only keeps the
// gradient-fill track (assets/gsxui.css) in sync once the user actually
// interacts with the thumb — see docs/jsx-parity.md `## slider`.
import { on } from "./gsxui.js";

on("input", "[data-gsxui-slider]", (_event, slider) => {
  const min = Number(slider.min || 0);
  const max = Number(slider.max || 100);
  const value = Number(slider.value);
  const fill = max > min ? ((value - min) / (max - min)) * 100 : 0;
  slider.style.setProperty("--fill", `${fill}%`);
});
