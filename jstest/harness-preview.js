// Browser-native equivalent of web/preview.js for the Playwright harness.
// The manifest supplies the already-compiled site stylesheet. Like the
// production entry, the preview loads the full behaviour barrel: the gallery
// it renders is interactive (menus, dialogs, tooltips), and the behaviours
// scope via data-gsxui-slot-* markers, so the twice-rendered gallery is safe.
import "../ui/index.js";
import "../web/theme-preview.js";
