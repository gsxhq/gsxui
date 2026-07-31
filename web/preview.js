// The preview document renders the real component gallery, so it needs the
// real behaviours: menus, dialogs, tooltips and friends all key on
// data-gsxui-slot-* markers and scope via closest() roots, which keeps the
// twice-rendered (nova + maia) gallery safe. The barrel's modules are
// per-marker delegates, so the deliberately absent Toaster region is simply
// never activated rather than an error.
import "../ui/index.js";
import "./theme-preview.js";
import "./site.css";
