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
    },
  },
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
    "gsxui:p1:H32",
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
    transport: { ...schema.transport, compact },
    palette: {
      baseColors: compact.baseColors.map((name) => ({ name })),
      themes,
      radii,
      resolved,
      defaultSelection: { baseColor: "neutral", theme: "neutral", radius: "medium" },
    },
  };
}

function expectedCompactCode(valueSchema, style, baseColor, theme, radius) {
  const compact = valueSchema.transport.compact;
  const packed =
    compact.styles.indexOf(style) |
    (compact.baseColors.indexOf(baseColor) << 4) |
    (compact.themes.indexOf(theme) << 8) |
    (compact.radii.indexOf(radius) << 13);
  const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
  let value = packed;
  let payload = "";
  do {
    payload = alphabet[value % 62] + payload;
    value = Math.floor(value / 62);
  } while (value > 0);
  return `${valueSchema.transport.compactPrefix}${payload}`;
}
