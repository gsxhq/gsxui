import test from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";

import {
  applyField,
  canonicalJSON,
  commandStrings,
  createThemeState,
  decodeShare,
  encodeShare,
  importPresetJSON,
  importThemeCSS,
  loadShareFromURL,
  manualCopyState,
  resetThemeState,
  selectMode,
  selectStyle,
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
};

const validators = {
  color(value) {
    return typeof value === "string" && value !== "" && value !== "invalid";
  },
  radius(value) {
    return /^(?:0|[0-9.]+(?:px|rem))$/.test(value);
  },
};

test("style and mode changes preserve the edited preset", () => {
  let state = createThemeState(schema);
  state = applyField(state, "light.primary", "rebeccapurple", validators);
  state = applyField(state, "radius", "12px", validators);
  state = selectStyle(state, "maia", schema);
  assert.equal(state.resolved.style, "maia");
  assert.equal(state.resolved.theme.light.primary, "rebeccapurple");
  assert.equal(state.resolved.radius, "12px");

  const beforeMode = structuredClone(state.resolved);
  state = selectMode(state, "dark");
  assert.equal(state.mode, "dark");
  assert.deepEqual(state.resolved, beforeMode);
});

test("invalid drafts do not mutate resolved preview or export state", () => {
  const initial = createThemeState(schema);
  const valid = applyField(initial, "light.primary", "red", validators);
  const invalid = applyField(valid, "light.primary", "invalid", validators);
  assert.equal(invalid.drafts.get("light.primary"), "invalid");
  assert.equal(invalid.resolved.theme.light.primary, "red");
  assert.equal(canonicalJSON(invalid.resolved, schema), canonicalJSON(valid.resolved, schema));
});

test("reset restores the selected style built-in preset", () => {
  let state = selectStyle(createThemeState(schema), "maia", schema);
  state = applyField(state, "dark.primary", "red", validators);
  state = resetThemeState(state, schema);
  assert.deepEqual(state.resolved, schema.defaults.maia);
  assert.equal(state.drafts.size, 0);
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
