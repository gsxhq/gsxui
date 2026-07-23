// Theme editor (/theme) — vanilla ESM riding ui/gsxui.js's delegated on().
// No framework, no build-time state: two in-memory var maps (light/dark)
// seeded from the server-rendered inputs are the single source of truth;
// every mutation (typing, tab switch, import) re-derives the preview node's
// inline style from whichever map is active. "Applying" always clears every
// tracked property first, so stale values from the other mode never linger
// (see .superpowers/sdd/task-5-brief.md's "simplest correct model").
import { on } from "../ui/gsxui.js";

const preview = document.querySelector("[data-theme-preview]");

// Canonical var order for export — matches assets/gsxui.css's :root block.
// --radius has no .dark override there (it's not a light/dark switch), so
// it's dropped from the exported .dark block below.
const VAR_ORDER = [
  "--radius",
  "--background",
  "--foreground",
  "--card",
  "--card-foreground",
  "--popover",
  "--popover-foreground",
  "--primary",
  "--primary-foreground",
  "--secondary",
  "--secondary-foreground",
  "--muted",
  "--muted-foreground",
  "--accent",
  "--accent-foreground",
  "--destructive",
  "--destructive-foreground",
  "--border",
  "--input",
  "--ring",
];

if (preview) {
  const vars = { light: {}, dark: {} };

  for (const el of document.querySelectorAll("[data-theme-var]")) {
    const name = el.dataset.themeVar;
    const mode = el.dataset.themeMode;
    if (mode === "light" || mode === "dark") vars[mode][name] = el.value;
  }

  let activeMode = "light";

  function applyActiveMode() {
    for (const name of Object.keys(vars.light)) preview.style.removeProperty(name);
    for (const name of Object.keys(vars.dark)) preview.style.removeProperty(name);
    for (const [name, value] of Object.entries(vars[activeMode])) {
      preview.style.setProperty(name, value);
    }
    preview.classList.toggle("dark", activeMode === "dark");
  }

  applyActiveMode();

  // --- Controls: typing a value updates its map; only reapply to the
  // preview when the edited input belongs to the currently active mode. ---
  on("input", "[data-theme-var]", (_event, el) => {
    const name = el.dataset.themeVar;
    const mode = el.dataset.themeMode;
    if (mode !== "light" && mode !== "dark") return;
    vars[mode][name] = el.value;
    if (mode === activeMode) applyActiveMode();
  });

  // --- Light/Dark tabs: switch which map drives the preview, and toggle
  // the "dark" class so dark:-variant utilities inside the preview engage
  // the same way they do under a real .dark ancestor. ---
  const tabButtons = document.querySelectorAll("[data-theme-tab]");

  function setActiveTab(mode) {
    activeMode = mode;
    for (const btn of tabButtons) {
      const isActive = btn.dataset.themeTab === mode;
      btn.classList.toggle("bg-accent", isActive);
      btn.classList.toggle("text-accent-foreground", isActive);
      btn.classList.toggle("text-muted-foreground", !isActive);
      btn.setAttribute("aria-pressed", String(isActive));
    }
    applyActiveMode();
  }

  on("click", "[data-theme-tab]", (_event, el) => {
    setActiveTab(el.dataset.themeTab);
  });

  // --- Export: assemble the full gsxui.css text (boilerplate skeleton from
  // assets/gsxui.css + the current maps), then copy or download it. ---
  function buildCss() {
    const rootLines = VAR_ORDER.map((name) => `  ${name}: ${vars.light[name]};`).join("\n");
    const darkLines = VAR_ORDER.filter((name) => name !== "--radius")
      .map((name) => `  ${name}: ${vars.dark[name]};`)
      .join("\n");
    return `/* gsxui theme — shadcn-compatible tokens on Tailwind v4.
   Exported from /theme. Drop this in place of assets/gsxui.css (or
   web/site.css, for the site itself). */
@import "tailwindcss";

@custom-variant dark (&:is(.dark *));

@theme inline {
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-destructive-foreground: var(--destructive-foreground);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
}

:root {
${rootLines}
}

.dark {
${darkLines}
}

@layer base {
  * {
    @apply border-border outline-ring/50;
  }
  body {
    @apply bg-background text-foreground;
  }
}
`;
  }

  const exportOutput = document.querySelector("[data-theme-export-output]");

  async function copyCss(css) {
    try {
      await navigator.clipboard.writeText(css);
      return true;
    } catch {
      return false;
    }
  }

  on("click", "[data-theme-copy]", async () => {
    const css = buildCss();
    const copied = await copyCss(css);
    if (!copied && exportOutput) {
      // Clipboard access can fail (permissions, insecure context, older
      // browsers) — fall back to a visible, pre-selected textarea so the
      // user can copy manually instead of silently doing nothing.
      exportOutput.value = css;
      exportOutput.classList.remove("hidden");
      exportOutput.focus();
      exportOutput.select();
    }
  });

  on("click", "[data-theme-download]", () => {
    const css = buildCss();
    const blob = new Blob([css], { type: "text/css" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "gsxui.css";
    document.body.appendChild(a);
    a.click();
    a.remove();
    URL.revokeObjectURL(url);
  });

  // --- Import: parse a pasted tweakcn/shadcn-style ":root { ... }" /
  // ".dark { ... }" block and apply any recognized "--var: value;" pairs. ---
  function extractBlock(text, blockRe) {
    const match = text.match(blockRe);
    return match ? match[1] : null;
  }

  function parseVars(blockText) {
    const result = {};
    if (!blockText) return result;
    const pairRe = /--([a-zA-Z0-9-]+)\s*:\s*([^;}]+)(?:;)?/g;
    let m;
    while ((m = pairRe.exec(blockText))) {
      result[`--${m[1]}`] = m[2].trim();
    }
    return result;
  }

  function applyParsed(parsed, mode) {
    for (const [name, value] of Object.entries(parsed)) {
      if (!(name in vars[mode])) continue; // only the 20 tracked vars
      vars[mode][name] = value;
      const input = document.querySelector(
        `[data-theme-var="${name}"][data-theme-mode="${mode}"]`,
      );
      if (input) input.value = value;
    }
  }

  on("click", "[data-theme-import-apply]", () => {
    const textarea = document.querySelector("[data-theme-import]");
    const text = textarea?.value ?? "";
    if (!text.trim()) return;

    const rootBlock = extractBlock(text, /:root\s*{([^}]*)}/s);
    const darkBlock = extractBlock(text, /\.dark\s*{([^}]*)}/s);

    applyParsed(parseVars(rootBlock), "light");
    applyParsed(parseVars(darkBlock), "dark");
    applyActiveMode();
  });
}
