// Browser-native equivalent of web/main.js for the Playwright harness.
// The harness serves real source modules directly, so it omits only Vite's
// virtual dev-panel import and lets the test manifest provide compiled CSS.
import "../ui/index.js";
import "../web/site.js";
import "../web/theme.js";
