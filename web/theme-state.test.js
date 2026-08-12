import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";

import {
  HISTORY_LIMIT,
  canRedo,
  canUndo,
  canonicalJSON,
  clearPalettePreview,
  commandStrings,
  createHistory,
  createThemeState,
  decodeShare,
  encodeShare,
  importPresetJSON,
  importThemeCSS,
  loadShareFromURL,
  manualCopyState,
  previewBaseColor,
  previewChartColor,
  previewFont,
  previewFontHeading,
  previewPreset,
  previewRadius,
  previewTheme,
  pushHistory,
  redoHistory,
  replacePreset,
  resetThemeState,
  selectBaseColor,
  selectChartColor,
  selectFont,
  selectFontHeading,
  selectMenuAccent,
  selectMode,
  selectPage,
  selectRadius,
  selectStyle,
  selectTheme,
  themeCSS,
  undoHistory,
} from "./theme-state.js";

// Shared fixtures: the "neutral" base color/chart color baseline, reused
// across defaults, resolved, and chartColors below so every fixture that
// is supposed to start out matching "neutral" actually does, byte for
// byte (matchedPaletteSelection compares structurally).
const neutralBaseValues = { background: "white", primary: "black", "primary-foreground": "white", accent: "gray", "accent-foreground": "black" };
const neutralBaseValuesDark = { background: "black", primary: "white", "primary-foreground": "black", accent: "darkgray", "accent-foreground": "white" };
const neutralChartValues = {
  light: { "chart-1": "gainsboro", "chart-2": "silver", "chart-3": "darkgray", "chart-4": "dimgray", "chart-5": "black" },
  dark: { "chart-1": "black", "chart-2": "dimgray", "chart-3": "darkgray", "chart-4": "silver", "chart-5": "gainsboro" },
};

const schema = {
  schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
  schemaVersion: 1,
  transport: {
    fullPrefix: "gsxui:v1:",
    compactPrefix: "gsxui:p1:",
    compact: {
      styles: ["nova", "maia"],
      baseColors: ["neutral", "stone", "zinc", "mauve", "olive", "mist", "taupe"],
      themes: [
        "neutral",
        "stone",
        "zinc",
        "mauve",
        "olive",
        "mist",
        "taupe",
        "amber",
        "blue",
        "cyan",
        "emerald",
        "fuchsia",
        "green",
        "indigo",
        "lime",
        "orange",
        "pink",
        "purple",
        "red",
        "rose",
        "sky",
        "teal",
        "violet",
        "yellow",
      ],
      radii: ["none", "small", "medium", "large"],
      menuAccents: ["subtle", "bold"],
      fonts: ["geist", "inter"],
    },
  },
  tokenNames: [
    "background",
    "primary",
    "primary-foreground",
    "accent",
    "accent-foreground",
    "chart-1",
    "chart-2",
    "chart-3",
    "chart-4",
    "chart-5",
  ],
  styles: ["nova", "maia"],
  defaults: {
    nova: {
      $schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
      schemaVersion: 1,
      style: "nova",
      radius: "0.625rem",
      fontSans: "Geist",
      fontHeading: "Geist",
      theme: {
        light: { ...neutralBaseValues, ...neutralChartValues.light },
        dark: { ...neutralBaseValuesDark, ...neutralChartValues.dark },
      },
    },
    maia: {
      $schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
      schemaVersion: 1,
      style: "maia",
      radius: "1rem",
      fontSans: "Geist",
      fontHeading: "Geist",
      theme: {
        light: {
          background: "ivory",
          primary: "navy",
          "primary-foreground": "ivory",
          accent: "lightblue",
          "accent-foreground": "navy",
          ...neutralChartValues.light,
        },
        dark: {
          background: "navy",
          primary: "ivory",
          "primary-foreground": "navy",
          accent: "darkblue",
          "accent-foreground": "ivory",
          ...neutralChartValues.dark,
        },
      },
    },
  },
  palette: {
    baseColors: [{ name: "neutral" }, { name: "stone" }],
    themes: {
      neutral: [{ name: "neutral" }, { name: "blue" }, { name: "rose" }],
      stone: [{ name: "stone" }, { name: "blue" }, { name: "rose" }],
    },
    radii: [
      { name: "none", value: "0" },
      { name: "medium", value: "0.625rem" },
      { name: "large", value: "0.875rem" },
    ],
    menuAccents: [
      { name: "subtle", title: "Subtle" },
      { name: "bold", title: "Bold" },
    ],
    fonts: [
      { name: "geist", title: "Geist", stack: "Geist" },
      { name: "inter", title: "Inter", stack: "Inter" },
    ],
    // Flat by name, not nested by base color: matches theme.gsx's own
    // schema.palette.chartColors shape (see catalog.go's ChartColorChoices
    // doc comment — the chart color axis is independent of base color).
    chartColors: {
      neutral: neutralChartValues,
      violet: {
        light: { "chart-1": "lavender", "chart-2": "mediumpurple", "chart-3": "blueviolet", "chart-4": "indigo", "chart-5": "purple" },
        dark: { "chart-1": "purple", "chart-2": "indigo", "chart-3": "blueviolet", "chart-4": "mediumpurple", "chart-5": "lavender" },
      },
    },
    resolved: {
      neutral: {
        neutral: {
          light: { ...neutralBaseValues, ...neutralChartValues.light },
          dark: { ...neutralBaseValuesDark, ...neutralChartValues.dark },
        },
        blue: {
          light: { background: "white", primary: "blue", "primary-foreground": "white", accent: "gray", "accent-foreground": "black", ...neutralChartValues.light },
          dark: { background: "black", primary: "skyblue", "primary-foreground": "black", accent: "darkgray", "accent-foreground": "white", ...neutralChartValues.dark },
        },
        rose: {
          light: { background: "mistyrose", primary: "crimson", "primary-foreground": "white", accent: "gray", "accent-foreground": "black", ...neutralChartValues.light },
          dark: { background: "maroon", primary: "pink", "primary-foreground": "black", accent: "darkgray", "accent-foreground": "white", ...neutralChartValues.dark },
        },
      },
      stone: {
        stone: {
          light: { background: "linen", primary: "sienna", "primary-foreground": "white", accent: "silver", "accent-foreground": "black", ...neutralChartValues.light },
          dark: { background: "sienna", primary: "linen", "primary-foreground": "black", accent: "dimgray", "accent-foreground": "white", ...neutralChartValues.dark },
        },
        blue: {
          light: { background: "linen", primary: "blue", "primary-foreground": "white", accent: "silver", "accent-foreground": "black", ...neutralChartValues.light },
          dark: { background: "sienna", primary: "skyblue", "primary-foreground": "black", accent: "dimgray", "accent-foreground": "white", ...neutralChartValues.dark },
        },
        rose: {
          light: { background: "linen", primary: "crimson", "primary-foreground": "white", accent: "silver", "accent-foreground": "black", ...neutralChartValues.light },
          dark: { background: "sienna", primary: "pink", "primary-foreground": "black", accent: "dimgray", "accent-foreground": "white", ...neutralChartValues.dark },
        },
      },
    },
    defaultSelection: { baseColor: "neutral", theme: "neutral", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" },
  },
};

const validators = {
  color(value) {
    return typeof value === "string" && value !== "" && value !== "invalid";
  },
  radius(value) {
    return /^(?:0|[0-9.]+(?:px|rem))$/.test(value);
  },
  fontFamily(value) {
    return typeof value === "string" && value !== "" && value !== "invalid";
  },
};

test("state starts matched to the default palette selection", () => {
  const state = createThemeState(schema);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "neutral", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.previewResolved, null);
});

test("palette selections commit colors while preserving the current radius", () => {
  let state = selectStyle(createThemeState(schema), "maia", schema);
  state = selectBaseColor(state, "stone", schema);
  state = selectTheme(state, "blue", schema);
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.style, "maia");
  assert.equal(state.resolved.theme.light.background, "linen");
  assert.equal(state.resolved.theme.light.primary, "blue");
  assert.equal(state.resolved.radius, "0.625rem");

  const beforeMode = structuredClone(state.resolved);
  state = selectMode(state, "dark");
  assert.equal(state.mode, "dark");
  assert.deepEqual(state.resolved, beforeMode);
});

test("menu accent overlays accent from primary and round-trips back to subtle", () => {
  let state = selectTheme(createThemeState(schema), "blue", schema);
  const subtleLight = structuredClone(state.resolved.theme.light);
  const subtleDark = structuredClone(state.resolved.theme.dark);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });

  state = selectMenuAccent(state, "bold", schema);
  assert.equal(state.selection.menuAccent, "bold");
  assert.equal(state.resolved.theme.light.accent, state.resolved.theme.light.primary);
  assert.equal(state.resolved.theme.light["accent-foreground"], state.resolved.theme.light["primary-foreground"]);
  assert.equal(state.resolved.theme.dark.accent, state.resolved.theme.dark.primary);
  assert.equal(state.resolved.theme.dark["accent-foreground"], state.resolved.theme.dark["primary-foreground"]);
  // Bold only touches accent/accent-foreground; base color, theme, and
  // radius are otherwise untouched.
  assert.equal(state.resolved.theme.light.background, subtleLight.background);
  assert.equal(state.resolved.theme.light.primary, subtleLight.primary);

  state = selectMenuAccent(state, "subtle", schema);
  assert.equal(state.selection.menuAccent, "subtle");
  assert.deepEqual(state.resolved.theme.light, subtleLight);
  assert.deepEqual(state.resolved.theme.dark, subtleDark);
});

test("menu accent survives a later base color or theme change", () => {
  let state = selectMenuAccent(createThemeState(schema), "bold", schema);
  state = selectBaseColor(state, "stone", schema);
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "stone", radius: "medium", chartColor: "neutral", menuAccent: "bold", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light.accent, state.resolved.theme.light.primary);

  state = selectTheme(state, "blue", schema);
  assert.equal(state.selection.menuAccent, "bold");
  assert.equal(state.resolved.theme.dark.accent, state.resolved.theme.dark.primary);
});

test("menu accent selection rejects unsupported values", () => {
  assert.throws(() => selectMenuAccent(createThemeState(schema), "loud", schema));
});

test("chart color overlays only chart-1..5, independent of base color and theme", () => {
  let state = selectTheme(createThemeState(schema), "blue", schema);
  const beforeLight = structuredClone(state.resolved.theme.light);
  const beforeDark = structuredClone(state.resolved.theme.dark);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });

  state = selectChartColor(state, "violet", schema);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "blue", radius: "medium", chartColor: "violet", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light["chart-1"], "lavender");
  assert.equal(state.resolved.theme.dark["chart-1"], "purple");
  // Only chart-1..5 changed — base color, theme, and menu accent's own
  // tokens are untouched.
  for (const name of ["background", "primary", "accent", "accent-foreground"]) {
    assert.equal(state.resolved.theme.light[name], beforeLight[name]);
    assert.equal(state.resolved.theme.dark[name], beforeDark[name]);
  }

  state = selectChartColor(state, "neutral", schema);
  assert.deepEqual(state.resolved.theme.light, beforeLight);
  assert.deepEqual(state.resolved.theme.dark, beforeDark);
});

test("chart color survives a later base color, theme, or menu accent change", () => {
  let state = selectChartColor(createThemeState(schema), "violet", schema);
  state = selectBaseColor(state, "stone", schema);
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "stone", radius: "medium", chartColor: "violet", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light["chart-1"], "lavender");

  state = selectTheme(state, "blue", schema);
  assert.equal(state.selection.chartColor, "violet");
  assert.equal(state.resolved.theme.dark["chart-1"], "purple");

  state = selectMenuAccent(state, "bold", schema);
  assert.equal(state.selection.chartColor, "violet");
  assert.equal(state.resolved.theme.light["chart-1"], "lavender");
  assert.equal(state.resolved.theme.light.accent, state.resolved.theme.light.primary);
});

test("chart color hover preview never alters committed export state", () => {
  const committed = createThemeState(schema);
  const committedJSON = canonicalJSON(committed.resolved, schema);
  const preview = previewChartColor(committed, "violet", schema);

  assert.equal(previewPreset(preview).theme.light["chart-1"], "lavender");
  assert.equal(canonicalJSON(preview.resolved, schema), committedJSON);

  const cleared = clearPalettePreview(preview);
  assert.equal(cleared.previewResolved, null);
  assert.deepEqual(previewPreset(cleared), committed.resolved);
});

test("chart color selection rejects unsupported values", () => {
  assert.throws(() => selectChartColor(createThemeState(schema), "plaid", schema));
});

test("font selection sets fontSans/fontHeading independently and never touches theme colors", () => {
  const state = createThemeState(schema);
  const beforeLight = structuredClone(state.resolved.theme.light);
  const beforeDark = structuredClone(state.resolved.theme.dark);
  assert.equal(state.resolved.fontSans, "Geist");
  assert.equal(state.resolved.fontHeading, "Geist");
  assert.deepEqual(state.selection.fontSans, "geist");
  assert.deepEqual(state.selection.fontHeading, "geist");

  const sans = selectFont(state, "inter", schema);
  assert.equal(sans.resolved.fontSans, "Inter");
  assert.equal(sans.resolved.fontHeading, "Geist");
  assert.equal(sans.selection.fontSans, "inter");
  assert.equal(sans.selection.fontHeading, "geist");
  assert.deepEqual(sans.resolved.theme.light, beforeLight);
  assert.deepEqual(sans.resolved.theme.dark, beforeDark);

  const heading = selectFontHeading(sans, "inter", schema);
  assert.equal(heading.resolved.fontSans, "Inter");
  assert.equal(heading.resolved.fontHeading, "Inter");
  assert.equal(heading.selection.fontHeading, "inter");

  const back = selectFontHeading(heading, "geist", schema);
  assert.equal(back.resolved.fontHeading, "Geist");
  assert.equal(back.selection.fontHeading, "geist");
});

test("font selection survives a later base color, theme, or menu accent change", () => {
  let state = selectFont(createThemeState(schema), "inter", schema);
  state = selectFontHeading(state, "inter", schema);
  state = selectBaseColor(state, "stone", schema);
  assert.equal(state.resolved.fontSans, "Inter");
  assert.equal(state.resolved.fontHeading, "Inter");

  state = selectMenuAccent(state, "bold", schema);
  assert.equal(state.resolved.fontSans, "Inter");
  assert.equal(state.resolved.fontHeading, "Inter");
});

test("font hover preview never alters committed export state", () => {
  const committed = createThemeState(schema);
  const committedJSON = canonicalJSON(committed.resolved, schema);

  const sansPreview = previewFont(committed, "inter", schema);
  assert.equal(previewPreset(sansPreview).fontSans, "Inter");
  assert.equal(canonicalJSON(sansPreview.resolved, schema), committedJSON);

  const headingPreview = previewFontHeading(committed, "inter", schema);
  assert.equal(previewPreset(headingPreview).fontHeading, "Inter");
  assert.equal(canonicalJSON(headingPreview.resolved, schema), committedJSON);

  const cleared = clearPalettePreview(headingPreview);
  assert.equal(cleared.previewResolved, null);
  assert.deepEqual(previewPreset(cleared), committed.resolved);
});

test("font selection rejects unsupported values", () => {
  assert.throws(() => selectFont(createThemeState(schema), "comic-sans", schema));
  assert.throws(() => selectFontHeading(createThemeState(schema), "comic-sans", schema));
});

test("page selection defaults to components and switches independently of the preset", () => {
  const state = createThemeState(schema);
  assert.equal(state.page, "components");
  const beforeJSON = canonicalJSON(state.resolved, schema);

  const onProduct = selectPage(state, "product");
  assert.equal(onProduct.page, "product");
  // Page is a client-only view toggle (mirrors mode) — never part of the
  // exported/committed preset.
  assert.equal(canonicalJSON(onProduct.resolved, schema), beforeJSON);
  assert.deepEqual(onProduct.selection, state.selection);

  const back = selectPage(onProduct, "components");
  assert.equal(back.page, "components");
});

test("page selection rejects unsupported values", () => {
  assert.throws(() => selectPage(createThemeState(schema), "marketing"));
});

test("palette previews never alter committed export state", () => {
  const committed = createThemeState(schema);
  const basePreview = previewBaseColor(committed, "stone", schema);
  assert.equal(previewPreset(basePreview).theme.light.background, "linen");
  assert.equal(canonicalJSON(basePreview.resolved, schema), canonicalJSON(committed.resolved, schema));

  const preview = previewTheme(committed, "rose", schema);
  assert.equal(previewPreset(preview).theme.light.primary, "crimson");
  assert.equal(canonicalJSON(preview.resolved, schema), canonicalJSON(committed.resolved, schema));

  const cleared = clearPalettePreview(preview);
  assert.equal(cleared.previewResolved, null);
  assert.deepEqual(previewPreset(cleared), committed.resolved);
});

test("radius preview is transient and preserves committed export state", () => {
  const committed = createThemeState(schema);
  const committedJSON = canonicalJSON(committed.resolved, schema);
  const preview = previewRadius(committed, "large", schema);

  assert.equal(previewPreset(preview).radius, "0.875rem");
  assert.equal(canonicalJSON(preview.resolved, schema), committedJSON);

  const cleared = clearPalettePreview(preview);
  assert.equal(cleared.previewResolved, null);
  assert.deepEqual(previewPreset(cleared), committed.resolved);
});

test("replacement derives custom palette selections from exact imported values", () => {
  const changedColor = structuredClone(schema.defaults.nova);
  changedColor.theme.light.primary = "violet";
  let state = replacePreset(createThemeState(schema), changedColor, schema);
  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });

  const changedRadius = structuredClone(schema.defaults.nova);
  changedRadius.radius = "1rem";
  state = replacePreset(state, changedRadius, schema);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "neutral", radius: "custom", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
});

test("base selection from custom restores its same-named theme and preserves a custom radius", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.radius = "1rem";
  imported.theme.light.primary = "violet";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectBaseColor(state, "stone", schema);

  assert.deepEqual(state.selection, { baseColor: "stone", theme: "stone", radius: "custom", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light.background, "linen");
  assert.equal(state.resolved.radius, "1rem");
});

test("theme selection from custom uses neutral as its base", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.theme.light.primary = "violet";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectTheme(state, "blue", schema);

  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light.background, "white");
  assert.equal(state.resolved.theme.light.primary, "blue");
});

test("named radius selection preserves imported custom colors exactly", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.theme.light.background = "rgb(1 2 3)";
  imported.theme.dark.primary = "hsl(1deg 2% 3%)";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectRadius(state, "large", schema);

  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "large", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(state.resolved.theme.light.background, "rgb(1 2 3)");
  assert.equal(state.resolved.theme.dark.primary, "hsl(1deg 2% 3%)");
  assert.equal(state.resolved.radius, "0.875rem");
});

test("reset restores the selected style built-in preset and recomputes selection", () => {
  let state = selectStyle(createThemeState(schema), "maia", schema);
  state = resetThemeState(state, schema);
  assert.deepEqual(state.resolved, schema.defaults.maia);
  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "custom", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
});

test("canonical JSON, share code, share URL, CSS, and commands round trip", () => {
  const preset = createThemeState(schema).resolved;
  const json = canonicalJSON(preset, schema);
  assert.equal(importPresetJSON(json, schema, validators).style, "nova");

  const share = encodeShare(preset, schema);
  assert.equal(share, "gsxui:p1:4GG");
  assert.deepEqual(decodeShare(share, schema, validators), preset);
  assert.deepEqual(
    loadShareFromURL(`https://ui.example/theme?preset=${encodeURIComponent(share)}`, schema, validators),
    preset,
  );
  assert.match(themeCSS(preset, schema), /:root \{\n  --radius: 0.625rem;/);
  assert.deepEqual(commandStrings(share), {
    init: `gsxui init --preset '${share}'`,
    apply: `gsxui apply --preset '${share}'`,
  });
});

test("bold menu accent still uses the compact transport and decodes back to bold", () => {
  const state = selectMenuAccent(createThemeState(schema), "bold", schema);
  const share = encodeShare(state.resolved, schema);
  assert.match(share, /^gsxui:p1:/);
  const decoded = decodeShare(share, schema, validators);
  assert.deepEqual(decoded, state.resolved);
  assert.equal(decoded.theme.light.accent, decoded.theme.light.primary);
});

test("a non-default chart color still uses the compact transport and decodes back to it", () => {
  const state = selectChartColor(createThemeState(schema), "violet", schema);
  const share = encodeShare(state.resolved, schema);
  assert.match(share, /^gsxui:p1:/);
  const decoded = decodeShare(share, schema, validators);
  assert.deepEqual(decoded, state.resolved);
  assert.equal(decoded.theme.light["chart-1"], "lavender");
});

test("non-default fonts still use the compact transport and decode back to them", () => {
  let state = selectFont(createThemeState(schema), "inter", schema);
  state = selectFontHeading(state, "geist", schema);
  const share = encodeShare(state.resolved, schema);
  assert.match(share, /^gsxui:p1:/);
  const decoded = decodeShare(share, schema, validators);
  assert.deepEqual(decoded, state.resolved);
  assert.equal(decoded.fontSans, "Inter");
  assert.equal(decoded.fontHeading, "Geist");
});

test("compact transport matches the Go bit layout for every catalogue combination", () => {
  const exhaustive = compactTestSchema();
  for (const style of exhaustive.transport.compact.styles) {
    for (const baseColor of exhaustive.transport.compact.baseColors) {
      for (const theme of exhaustive.transport.compact.themes) {
        if (
          exhaustive.transport.compact.baseColors.includes(theme) &&
          theme !== baseColor
        ) {
          continue;
        }
        for (const radius of exhaustive.palette.radii) {
          const resolved = exhaustive.palette.resolved[baseColor][theme];
          const preset = {
            $schema: exhaustive.schema,
            schemaVersion: exhaustive.schemaVersion,
            style,
            radius: radius.value,
            fontSans: exhaustive.palette.fonts[0].stack,
            fontHeading: exhaustive.palette.fonts[0].stack,
            theme: structuredClone(resolved),
          };
          const want = expectedCompactCode(exhaustive, style, baseColor, theme, radius.name);
          assert.equal(encodeShare(preset, exhaustive), want);
          assert.deepEqual(decodeShare(want, exhaustive, validators), preset);
        }
      }
    }
  }
});

test("custom values use full transport and historical full codes remain valid", () => {
  const customTheme = structuredClone(schema.defaults.nova);
  customTheme.theme.light.primary = "violet";
  const customThemeCode = encodeShare(customTheme, schema);
  assert.match(customThemeCode, /^gsxui:v1:/);
  assert.deepEqual(decodeShare(customThemeCode, schema, validators), customTheme);

  const customRadius = structuredClone(schema.defaults.nova);
  customRadius.radius = "1rem";
  const customRadiusCode = encodeShare(customRadius, schema);
  assert.match(customRadiusCode, /^gsxui:v1:/);
  assert.deepEqual(decodeShare(customRadiusCode, schema, validators), customRadius);

  const fullBuiltIn = `gsxui:v1:${Buffer.from(canonicalJSON(schema.defaults.nova, schema)).toString("base64url")}`;
  assert.deepEqual(decodeShare(fullBuiltIn, schema, validators), schema.defaults.nova);
});

test("compact decoder rejects malformed and non-canonical payloads", () => {
  for (const code of [
    "gsxui:p1:",
    "gsxui:p1:!",
    "gsxui:p1:zzzzzzzzzzzz",
    "gsxui:p1:2",
    "gsxui:p1:1o",
    "gsxui:p1:1bc",
    "gsxui:p1:8WW",
    "gsxui:p1:G",
    "gsxui:p1:16C8",
    "gsxui:p1:04GG",
  ]) {
    assert.throws(() => decodeShare(code, schema, validators));
  }
});

test("transient previews never alter compact transport", () => {
  const committed = createThemeState(schema);
  const share = encodeShare(committed.resolved, schema);
  for (const preview of [
    previewBaseColor(committed, "stone", schema),
    previewTheme(committed, "blue", schema),
    previewRadius(committed, "large", schema),
  ]) {
    assert.equal(encodeShare(preview.resolved, schema), share);
  }
});

test("browser canonical JSON matches the Go golden byte for byte", () => {
  const source = readFileSync(
    new URL("../internal/preset/testdata/default-nova.json", import.meta.url),
    "utf8",
  );
  const preset = JSON.parse(source);
  const realSchema = { tokenNames: Object.keys(preset.theme.light) };
  assert.equal(canonicalJSON(preset, realSchema), source);
});

test("JSON import rejects duplicates, unknown fields, comments, and trailing commas atomically", () => {
  const good = canonicalJSON(schema.defaults.nova, schema);
  for (const source of [
    good.replace('"style": "nova",', '"style": "nova",\n  "style": "maia",'),
    good.replace('"style": "nova",', '"style": "nova",\n  "unknown": true,'),
    `// comment\n${good}`,
    good.replace('"radius": "0.625rem",', '"radius": "0.625rem",,'),
  ]) {
    assert.throws(() => importPresetJSON(source, schema, validators));
  }
});

test("CSS import commits only a fully valid parsed candidate", () => {
  const state = createThemeState(schema);
  const parseCSS = () => ({
    radius: "",
    light: { primary: "blue" },
    dark: { primary: "white" },
  });
  const imported = importThemeCSS(state, "ignored by injected parser", schema, validators, parseCSS);
  assert.equal(imported.resolved.radius, "0.625rem");
  assert.equal(imported.resolved.theme.light.background, "white");
  assert.equal(imported.resolved.theme.light.primary, "blue");

  const invalidParser = () => ({
    radius: "",
    light: { primary: "invalid" },
    dark: {},
  });
  assert.throws(() => importThemeCSS(state, "invalid", schema, validators, invalidParser));
  assert.deepEqual(state.resolved, schema.defaults.nova);
});

test("undo/redo replays a bounded history stack of committed states", () => {
  let state = createThemeState(schema);
  let history = createHistory(state);
  assert.equal(canUndo(history), false);
  assert.equal(canRedo(history), false);

  state = selectBaseColor(state, "stone", schema);
  history = pushHistory(history, state);
  state = selectTheme(state, "blue", schema);
  history = pushHistory(history, state);
  state = selectRadius(state, "large", schema);
  history = pushHistory(history, state);
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "blue", radius: "large", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(canUndo(history), true);
  assert.equal(canRedo(history), false);

  let step = undoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });

  step = undoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "stone", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(canRedo(history), true);

  step = redoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "blue", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" });
  assert.equal(canUndo(history), true);
  assert.equal(canRedo(history), true);
});

// Font is the newest axis (Task 3); this proves it participates in the
// same commitHistory stack as every other picker, exactly like
// TestCompactShareBackwardCompatibleWithPreTask3Codes proves it participates
// in the share codec on the Go side.
test("undo/redo replays font selections exactly like every other axis", () => {
  let state = createThemeState(schema);
  let history = createHistory(state);

  state = selectFont(state, "inter", schema);
  history = pushHistory(history, state);
  state = selectFontHeading(state, "inter", schema);
  history = pushHistory(history, state);
  assert.equal(state.resolved.fontSans, "Inter");
  assert.equal(state.resolved.fontHeading, "Inter");

  let step = undoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.equal(state.resolved.fontSans, "Inter");
  assert.equal(state.resolved.fontHeading, "Geist");

  step = undoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.equal(state.resolved.fontSans, "Geist");
  assert.equal(state.resolved.fontHeading, "Geist");
  assert.equal(canUndo(history), false);

  step = redoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  step = redoHistory(history, state, schema);
  state = step.state;
  assert.equal(state.resolved.fontSans, "Inter");
  assert.equal(state.resolved.fontHeading, "Inter");
});

test("undo/redo never touches page, the same as mode", () => {
  let state = createThemeState(schema);
  let history = createHistory(state);
  state = selectPage(state, "product");
  state = selectMode(state, "dark");

  state = selectTheme(state, "blue", schema);
  history = pushHistory(history, state);

  // Undoing crosses a real commit boundary (back to the pre-"blue"
  // resolved preset), which is what actually exercises replacePreset's
  // clone-then-overwrite path — page and mode must survive it exactly
  // like they already do for every other reducer.
  const step = undoHistory(history, state, schema);
  assert.notEqual(step.state.resolved.theme.light.primary, state.resolved.theme.light.primary);
  assert.equal(step.state.page, "product");
  assert.equal(step.state.mode, "dark");
});

test("undo/redo at the stack's bounds is a no-op", () => {
  const state = createThemeState(schema);
  const history = createHistory(state);

  const undone = undoHistory(history, state, schema);
  assert.equal(undone.history, history);
  assert.equal(undone.state, state);

  const redone = redoHistory(history, state, schema);
  assert.equal(redone.history, history);
  assert.equal(redone.state, state);
});

test("a new commit after undo discards the redo branch", () => {
  let state = createThemeState(schema);
  let history = createHistory(state);

  state = selectTheme(state, "blue", schema);
  history = pushHistory(history, state);
  state = selectTheme(state, "rose", schema);
  history = pushHistory(history, state);

  let step = undoHistory(history, state, schema);
  history = step.history;
  state = step.state;
  assert.equal(state.selection.theme, "blue");

  state = selectRadius(state, "large", schema);
  history = pushHistory(history, state);
  assert.equal(canRedo(history), false);
  assert.equal(history.entries.length, 3);
  assert.equal(history.cursor, 2);
});

test("history is bounded to HISTORY_LIMIT entries, dropping the oldest", () => {
  let state = createThemeState(schema);
  let history = createHistory(state);
  const names = ["neutral", "blue", "rose"];

  for (let i = 0; i < HISTORY_LIMIT + 10; i++) {
    state = selectTheme(state, names[i % names.length], schema);
    history = pushHistory(history, state);
  }

  assert.equal(history.entries.length, HISTORY_LIMIT);
  assert.equal(history.cursor, HISTORY_LIMIT - 1);
  assert.deepEqual(history.entries.at(-1), state.resolved);
});

test("clipboard fallback retains selected manual-copy text", () => {
  assert.deepEqual(manualCopyState(true, "value"), { visible: false, text: "" });
  assert.deepEqual(manualCopyState(false, "value"), { visible: true, text: "value" });
});

function compactTestSchema() {
  const compact = structuredClone(schema.transport.compact);
  const radii = [
    { name: "none", value: "0" },
    { name: "small", value: "0.45rem" },
    { name: "medium", value: "0.625rem" },
    { name: "large", value: "0.875rem" },
  ];
  const resolved = {};
  const themes = {};
  for (const [baseIndex, baseColor] of compact.baseColors.entries()) {
    resolved[baseColor] = {};
    themes[baseColor] = [];
    for (const [themeIndex, theme] of compact.themes.entries()) {
      if (themeIndex < compact.baseColors.length && theme !== baseColor) continue;
      themes[baseColor].push({ name: theme });
      resolved[baseColor][theme] = {
        light: {
          background: `oklch(0.${baseIndex + 1} 0 0)`,
          primary: `oklch(0.${themeIndex + 1} 0.1 20)`,
        },
        dark: {
          background: `oklch(0.${baseIndex + 2} 0 0)`,
          primary: `oklch(0.${themeIndex + 2} 0.1 200)`,
        },
      };
    }
  }
  return {
    ...schema,
    // Scoped independently of the outer mock's tokenNames: this test's
    // hand-built resolved matrix only ever carries background/primary, so
    // it must not inherit the outer schema's accent-bearing token set.
    tokenNames: ["background", "primary"],
    transport: { ...schema.transport, compact },
    palette: {
      baseColors: compact.baseColors.map((name) => ({ name })),
      themes,
      radii,
      menuAccents: schema.palette.menuAccents,
      // This exhaustive sweep never carries chart-1..5 tokens at all (see
      // the tokenNames override above), so "neutral" chart color is
      // modeled as an empty overlay — containsThemeValues's "every key in
      // subset" check is vacuously true for an empty subset, so every
      // preset here matches it, which is exactly right: they are all
      // fixed at the chart color default throughout this sweep.
      chartColors: { neutral: { light: {}, dark: {} } },
      resolved,
      fonts: schema.palette.fonts,
      defaultSelection: { baseColor: "neutral", theme: "neutral", radius: "medium", chartColor: "neutral", menuAccent: "subtle", fontSans: "geist", fontHeading: "geist" },
    },
  };
}

function expectedCompactCode(valueSchema, style, baseColor, theme, radius) {
  const compact = valueSchema.transport.compact;
  const packed =
    compact.styles.indexOf(style) |
    (compact.baseColors.indexOf(baseColor) << 4) |
    (compact.themes.indexOf(theme) << 8) |
    (compact.radii.indexOf(radius) << 13) |
    (compact.menuAccents.indexOf("subtle") << 16) |
    (compact.themes.indexOf("neutral") << 18) |
    (compact.fonts.indexOf("geist") << 23) |
    (compact.fonts.indexOf("geist") << 26);
  const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  let value = packed;
  let payload = "";
  do {
    payload = alphabet[value % 62] + payload;
    value = Math.floor(value / 62);
  } while (value > 0);
  return `${valueSchema.transport.compactPrefix}${payload}`;
}
