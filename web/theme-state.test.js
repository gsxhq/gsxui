import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";

import {
  canonicalJSON,
  clearPalettePreview,
  commandStrings,
  createThemeState,
  decodeShare,
  encodeShare,
  importPresetJSON,
  importThemeCSS,
  loadShareFromURL,
  manualCopyState,
  previewBaseColor,
  previewPreset,
  previewRadius,
  previewTheme,
  replacePreset,
  resetThemeState,
  selectBaseColor,
  selectMode,
  selectRadius,
  selectStyle,
  selectTheme,
  themeCSS,
} from "./theme-state.js";

const schema = {
  schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
  schemaVersion: 1,
  transportPrefix: "gsxui:v1:",
  tokenNames: ["background", "primary"],
  styles: ["nova", "maia"],
  defaults: {
    nova: {
      $schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
      schemaVersion: 1,
      style: "nova",
      radius: "0.625rem",
      theme: {
        light: { background: "white", primary: "black" },
        dark: { background: "black", primary: "white" },
      },
    },
    maia: {
      $schema: "https://ui.gsxhq.dev/schemas/preset-v1.json",
      schemaVersion: 1,
      style: "maia",
      radius: "1rem",
      theme: {
        light: { background: "ivory", primary: "navy" },
        dark: { background: "navy", primary: "ivory" },
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
    resolved: {
      neutral: {
        neutral: {
          light: { background: "white", primary: "black" },
          dark: { background: "black", primary: "white" },
        },
        blue: {
          light: { background: "white", primary: "blue" },
          dark: { background: "black", primary: "skyblue" },
        },
        rose: {
          light: { background: "mistyrose", primary: "crimson" },
          dark: { background: "maroon", primary: "pink" },
        },
      },
      stone: {
        stone: {
          light: { background: "linen", primary: "sienna" },
          dark: { background: "sienna", primary: "linen" },
        },
        blue: {
          light: { background: "linen", primary: "blue" },
          dark: { background: "sienna", primary: "skyblue" },
        },
        rose: {
          light: { background: "linen", primary: "crimson" },
          dark: { background: "sienna", primary: "pink" },
        },
      },
    },
    defaultSelection: { baseColor: "neutral", theme: "neutral", radius: "medium" },
  },
};

const validators = {
  color(value) {
    return typeof value === "string" && value !== "" && value !== "invalid";
  },
  radius(value) {
    return /^(?:0|[0-9.]+(?:px|rem))$/.test(value);
  },
};

test("state starts matched to the default palette selection", () => {
  const state = createThemeState(schema);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "neutral", radius: "medium" });
  assert.equal(state.previewResolved, null);
});

test("palette selections commit colors while preserving the current radius", () => {
  let state = selectStyle(createThemeState(schema), "maia", schema);
  state = selectBaseColor(state, "stone", schema);
  state = selectTheme(state, "blue", schema);
  assert.deepEqual(state.selection, { baseColor: "stone", theme: "blue", radius: "medium" });
  assert.equal(state.resolved.style, "maia");
  assert.equal(state.resolved.theme.light.background, "linen");
  assert.equal(state.resolved.theme.light.primary, "blue");
  assert.equal(state.resolved.radius, "0.625rem");

  const beforeMode = structuredClone(state.resolved);
  state = selectMode(state, "dark");
  assert.equal(state.mode, "dark");
  assert.deepEqual(state.resolved, beforeMode);
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
  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "medium" });

  const changedRadius = structuredClone(schema.defaults.nova);
  changedRadius.radius = "1rem";
  state = replacePreset(state, changedRadius, schema);
  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "neutral", radius: "custom" });
});

test("base selection from custom restores its same-named theme and preserves a custom radius", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.radius = "1rem";
  imported.theme.light.primary = "violet";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectBaseColor(state, "stone", schema);

  assert.deepEqual(state.selection, { baseColor: "stone", theme: "stone", radius: "custom" });
  assert.equal(state.resolved.theme.light.background, "linen");
  assert.equal(state.resolved.radius, "1rem");
});

test("theme selection from custom uses neutral as its base", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.theme.light.primary = "violet";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectTheme(state, "blue", schema);

  assert.deepEqual(state.selection, { baseColor: "neutral", theme: "blue", radius: "medium" });
  assert.equal(state.resolved.theme.light.background, "white");
  assert.equal(state.resolved.theme.light.primary, "blue");
});

test("named radius selection preserves imported custom colors exactly", () => {
  const imported = structuredClone(schema.defaults.nova);
  imported.theme.light.background = "rgb(1 2 3)";
  imported.theme.dark.primary = "hsl(1deg 2% 3%)";
  let state = replacePreset(createThemeState(schema), imported, schema);
  state = selectRadius(state, "large", schema);

  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "large" });
  assert.equal(state.resolved.theme.light.background, "rgb(1 2 3)");
  assert.equal(state.resolved.theme.dark.primary, "hsl(1deg 2% 3%)");
  assert.equal(state.resolved.radius, "0.875rem");
});

test("reset restores the selected style built-in preset and recomputes selection", () => {
  let state = selectStyle(createThemeState(schema), "maia", schema);
  state = resetThemeState(state, schema);
  assert.deepEqual(state.resolved, schema.defaults.maia);
  assert.deepEqual(state.selection, { baseColor: "custom", theme: "custom", radius: "custom" });
});

test("canonical JSON, share code, share URL, CSS, and commands round trip", () => {
  const preset = selectStyle(createThemeState(schema), "maia", schema).resolved;
  const json = canonicalJSON(preset, schema);
  assert.equal(importPresetJSON(json, schema, validators).style, "maia");

  const share = encodeShare(preset, schema);
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

test("clipboard fallback retains selected manual-copy text", () => {
  assert.deepEqual(manualCopyState(true, "value"), { visible: false, text: "" });
  assert.deepEqual(manualCopyState(false, "value"), { visible: true, text: "value" });
});
