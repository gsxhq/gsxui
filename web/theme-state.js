import { getNodeValue, parseTree } from "jsonc-parser";

function clone(value) {
  return structuredClone(value);
}

function cloneState(state) {
  return {
    resolved: clone(state.resolved),
    selection: clone(state.selection),
    previewResolved: state.previewResolved === null ? null : clone(state.previewResolved),
    mode: state.mode,
  };
}

function assertStyle(style, schema) {
  if (!schema.styles.includes(style)) {
    throw new Error(`unsupported style ${String(style)}`);
  }
}

function validatePreset(preset, schema, validators) {
  if (!preset || typeof preset !== "object" || Array.isArray(preset)) {
    throw new Error("preset must be an object");
  }
  if (preset.$schema !== schema.schema) {
    throw new Error(`$schema must be ${schema.schema}`);
  }
  if (preset.schemaVersion !== schema.schemaVersion) {
    throw new Error(`schemaVersion must be ${schema.schemaVersion}`);
  }
  assertStyle(preset.style, schema);
  if (!validators.radius(preset.radius)) {
    throw new Error("radius is not a supported CSS length");
  }
  if (!preset.theme || typeof preset.theme !== "object" || Array.isArray(preset.theme)) {
    throw new Error("theme must be an object");
  }
  for (const mode of ["light", "dark"]) {
    const values = preset.theme[mode];
    if (!values || typeof values !== "object" || Array.isArray(values)) {
      throw new Error(`theme.${mode} must be an object`);
    }
    const keys = Object.keys(values).sort();
    const expected = [...schema.tokenNames].sort();
    if (keys.length !== expected.length || keys.some((key, index) => key !== expected[index])) {
      throw new Error(`theme.${mode} must contain every supported token exactly once`);
    }
    for (const name of schema.tokenNames) {
      if (!validators.color(values[name])) {
        throw new Error(`theme.${mode}.${name} is not a valid CSS color`);
      }
    }
  }
  return preset;
}

export function createThemeState(schema, preset = schema.defaults.nova) {
  return replacePreset(
    {
      resolved: clone(preset),
      selection: clone(schema.palette.defaultSelection),
      previewResolved: null,
      mode: "light",
    },
    preset,
    schema,
  );
}

function paletteResolution(schema, baseColor, theme) {
  const resolution = schema.palette.resolved[baseColor]?.[theme];
  if (!resolution) {
    throw new Error(`unsupported palette selection ${baseColor}/${theme}`);
  }
  return resolution;
}

// accentOverrideTokens mirrors catalog.go's accentOverrideTokens: the two
// tokens applyMenuAccent overwrites when bold, and so the two tokens
// excluded from matchedPaletteSelection's base color/theme comparison
// below (schema.palette.resolved is always baked at the "subtle" default,
// so a bold preset's accent/accent-foreground never equal it on those two
// keys alone, independent of which base color or theme actually matches).
const accentOverrideTokens = new Set(["accent", "accent-foreground"]);

function applyMenuAccent(values, accent) {
  if (accent !== "bold") return;
  values.accent = values.primary;
  values["accent-foreground"] = values["primary-foreground"];
}

function applyPaletteTheme(preset, schema, baseColor, theme, menuAccent) {
  const resolution = paletteResolution(schema, baseColor, theme);
  preset.theme.light = clone(resolution.light);
  preset.theme.dark = clone(resolution.dark);
  applyMenuAccent(preset.theme.light, menuAccent);
  applyMenuAccent(preset.theme.dark, menuAccent);
}

function themeSelectionForBase(state, baseColor, schema) {
  const theme = state.selection.theme;
  if (state.selection.baseColor === "custom" || !schema.palette.resolved[baseColor]?.[theme]) {
    return baseColor;
  }
  return theme;
}

// resolvedBaseColorAndTheme collapses "custom" the same way every committing
// palette reducer already does (see selectTheme/selectBaseColor below and
// the pre-existing tests documenting this): an axis picker always operates
// within the catalog model, so touching one axis while the others are
// custom snaps them to a concrete baseline first.
function resolvedBaseColorAndTheme(state) {
  const baseColor = state.selection.baseColor === "custom" ? "neutral" : state.selection.baseColor;
  const theme = state.selection.theme === "custom" ? baseColor : state.selection.theme;
  return { baseColor, theme };
}

function matchedPaletteSelection(preset, schema) {
  let baseColor = "custom";
  let theme = "custom";
  for (const [candidateBaseColor, themes] of Object.entries(schema.palette.resolved)) {
    for (const [candidateTheme, resolved] of Object.entries(themes)) {
      if (
        themeValuesEqual(preset.theme.light, resolved.light, schema.tokenNames, accentOverrideTokens) &&
        themeValuesEqual(preset.theme.dark, resolved.dark, schema.tokenNames, accentOverrideTokens)
      ) {
        baseColor = candidateBaseColor;
        theme = candidateTheme;
        break;
      }
    }
    if (baseColor !== "custom") break;
  }
  const radius = schema.palette.radii.find((choice) => choice.value === preset.radius)?.name ?? "custom";
  const menuAccent = matchMenuAccent(preset);
  return { baseColor, theme, radius, menuAccent };
}

function themeValuesEqual(a, b, tokenNames, excluded = null) {
  return tokenNames.every((name) => excluded?.has(name) || a[name] === b[name]);
}

// matchMenuAccent has no "custom" outcome, unlike the other axes: bold is
// defined structurally (accent/accent-foreground literally equal
// primary/primary-foreground in both modes, applyMenuAccent's own
// postcondition), so a preset is either bold or — by default — subtle.
// Mirrors catalog.go's matchMenuAccent exactly.
function matchMenuAccent(preset) {
  const bold =
    preset.theme.light.accent === preset.theme.light.primary &&
    preset.theme.light["accent-foreground"] === preset.theme.light["primary-foreground"] &&
    preset.theme.dark.accent === preset.theme.dark.primary &&
    preset.theme.dark["accent-foreground"] === preset.theme.dark["primary-foreground"];
  return bold ? "bold" : "subtle";
}

export function replacePreset(state, preset, schema) {
  const next = cloneState(state);
  next.resolved = clone(preset);
  next.selection = matchedPaletteSelection(next.resolved, schema);
  next.previewResolved = null;
  return next;
}

export function selectBaseColor(state, baseColor, schema) {
  const next = cloneState(state);
  const theme = themeSelectionForBase(state, baseColor, schema);
  applyPaletteTheme(next.resolved, schema, baseColor, theme, state.selection.menuAccent);
  next.selection.baseColor = baseColor;
  next.selection.theme = theme;
  next.previewResolved = null;
  return next;
}

export function selectTheme(state, theme, schema) {
  const next = cloneState(state);
  const baseColor = state.selection.baseColor === "custom" ? "neutral" : state.selection.baseColor;
  applyPaletteTheme(next.resolved, schema, baseColor, theme, state.selection.menuAccent);
  next.selection.baseColor = baseColor;
  next.selection.theme = theme;
  next.previewResolved = null;
  return next;
}

export function selectMenuAccent(state, menuAccent, schema) {
  if (!schema.palette.menuAccents.some((choice) => choice.name === menuAccent)) {
    throw new Error(`unsupported palette menu accent ${String(menuAccent)}`);
  }
  const next = cloneState(state);
  const { baseColor, theme } = resolvedBaseColorAndTheme(state);
  applyPaletteTheme(next.resolved, schema, baseColor, theme, menuAccent);
  next.selection.baseColor = baseColor;
  next.selection.theme = theme;
  next.selection.menuAccent = menuAccent;
  next.previewResolved = null;
  return next;
}

export function selectRadius(state, radius, schema) {
  const choice = schema.palette.radii.find((candidate) => candidate.name === radius);
  if (!choice) {
    throw new Error(`unsupported palette radius ${String(radius)}`);
  }
  const next = cloneState(state);
  next.resolved.radius = choice.value;
  next.selection.radius = radius;
  next.previewResolved = null;
  return next;
}

export function previewBaseColor(state, baseColor, schema) {
  const next = cloneState(state);
  const theme = themeSelectionForBase(state, baseColor, schema);
  next.previewResolved = clone(next.resolved);
  applyPaletteTheme(next.previewResolved, schema, baseColor, theme, state.selection.menuAccent);
  return next;
}

export function previewTheme(state, theme, schema) {
  const next = cloneState(state);
  const baseColor = state.selection.baseColor === "custom" ? "neutral" : state.selection.baseColor;
  next.previewResolved = clone(next.resolved);
  applyPaletteTheme(next.previewResolved, schema, baseColor, theme, state.selection.menuAccent);
  return next;
}

export function previewRadius(state, radius, schema) {
  const choice = schema.palette.radii.find((candidate) => candidate.name === radius);
  if (!choice) {
    throw new Error(`unsupported palette radius ${String(radius)}`);
  }
  const next = cloneState(state);
  next.previewResolved = clone(next.resolved);
  next.previewResolved.radius = choice.value;
  return next;
}

export function clearPalettePreview(state) {
  const next = cloneState(state);
  next.previewResolved = null;
  return next;
}

export function previewPreset(state) {
  return state.previewResolved ?? state.resolved;
}

export function selectStyle(state, style, schema) {
  assertStyle(style, schema);
  const next = cloneState(state);
  next.resolved.style = style;
  next.previewResolved = null;
  return next;
}

export function selectMode(state, mode) {
  if (mode !== "light" && mode !== "dark") {
    throw new Error("mode must be light or dark");
  }
  const next = cloneState(state);
  next.mode = mode;
  return next;
}

export function resetThemeState(state, schema) {
  assertStyle(state.resolved.style, schema);
  return replacePreset(state, schema.defaults[state.resolved.style], schema);
}

// Undo/redo is a bounded stack of committed `resolved` presets (never
// `previewResolved` — hover previews are never pushed, matching every
// committing reducer above, which already only ever mutates `resolved`).
// Mode (light/dark) intentionally rides along outside history: it is a
// client-only view toggle, not part of the preset schema, so undoing never
// flips it — only theme.gsx's calls decide when a state change is a commit
// worth a history entry, by calling pushHistory once per commit, at the
// same granularity the editor already commits at (click-driven pickers,
// style buttons, reset, and import-apply — there is no per-keystroke commit
// path in this editor to begin with, so no separate debounce is needed).
export const HISTORY_LIMIT = 50;

export function createHistory(state) {
  return { entries: [clone(state.resolved)], cursor: 0 };
}

export function pushHistory(history, state) {
  const entries = history.entries.slice(0, history.cursor + 1);
  entries.push(clone(state.resolved));
  const overflow = entries.length - HISTORY_LIMIT;
  const trimmed = overflow > 0 ? entries.slice(overflow) : entries;
  return { entries: trimmed, cursor: trimmed.length - 1 };
}

export function canUndo(history) {
  return history.cursor > 0;
}

export function canRedo(history) {
  return history.cursor < history.entries.length - 1;
}

function travelHistory(history, state, schema, delta) {
  const cursor = history.cursor + delta;
  if (cursor < 0 || cursor >= history.entries.length) return { history, state };
  const resolved = history.entries[cursor];
  return { history: { ...history, cursor }, state: replacePreset(state, resolved, schema) };
}

export function undoHistory(history, state, schema) {
  return travelHistory(history, state, schema, -1);
}

export function redoHistory(history, state, schema) {
  return travelHistory(history, state, schema, 1);
}

function orderedPreset(preset, schema) {
  const theme = {};
  for (const mode of ["light", "dark"]) {
    theme[mode] = {};
    for (const name of schema.tokenNames) {
      theme[mode][name] = preset.theme[mode][name];
    }
  }
  return {
    $schema: preset.$schema,
    schemaVersion: preset.schemaVersion,
    style: preset.style,
    radius: preset.radius,
    theme,
  };
}

export function canonicalJSON(preset, schema) {
  return `${JSON.stringify(orderedPreset(preset, schema), null, 2)}\n`;
}

export function themeCSS(preset, schema) {
  const root = [
    `  --radius: ${preset.radius};`,
    ...schema.tokenNames.map((name) => `  --${name}: ${preset.theme.light[name]};`),
  ].join("\n");
  const dark = schema.tokenNames
    .map((name) => `  --${name}: ${preset.theme.dark[name]};`)
    .join("\n");
  return `:root {\n${root}\n}\n\n.dark {\n${dark}\n}\n`;
}

function base64URL(bytes) {
  let binary = "";
  for (const byte of bytes) {
    binary += String.fromCharCode(byte);
  }
  return btoa(binary).replaceAll("+", "-").replaceAll("/", "_").replace(/=+$/, "");
}

function fromBase64URL(value) {
  if (value.includes("=")) {
    throw new Error("share code uses non-canonical padded base64");
  }
  const standard = value.replaceAll("-", "+").replaceAll("_", "/");
  const binary = atob(standard.padEnd(Math.ceil(standard.length / 4) * 4, "="));
  return Uint8Array.from(binary, (character) => character.charCodeAt(0));
}

const base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const maxUint64 = (1n << 64n) - 1n;

function encodeBase62(value) {
  if (value === 0n) return "0";
  let encoded = "";
  while (value > 0n) {
    encoded = base62Alphabet[Number(value % 62n)] + encoded;
    value /= 62n;
  }
  return encoded;
}

function decodeBase62(value) {
  if (value === "") throw new Error("compact payload is empty");
  let decoded = 0n;
  for (const character of value) {
    const digit = base62Alphabet.indexOf(character);
    if (digit === -1) throw new Error(`invalid compact base62 character ${character}`);
    decoded = decoded * 62n + BigInt(digit);
    if (decoded > maxUint64) throw new Error("compact base62 integer overflow");
  }
  if (encodeBase62(decoded) !== value) {
    throw new Error("compact base62 spelling is not canonical");
  }
  return decoded;
}

function compactSelection(preset, schema) {
  const selection = matchedPaletteSelection(preset, schema);
  if (selection.baseColor === "custom" || selection.theme === "custom" || selection.radius === "custom") {
    return null;
  }
  const compact = schema.transport.compact;
  const indexes = {
    style: compact.styles.indexOf(preset.style),
    baseColor: compact.baseColors.indexOf(selection.baseColor),
    theme: compact.themes.indexOf(selection.theme),
    radius: compact.radii.indexOf(selection.radius),
    menuAccent: compact.menuAccents.indexOf(selection.menuAccent),
  };
  return Object.values(indexes).includes(-1) ? null : indexes;
}

function encodeFullShare(preset, schema) {
  return `${schema.transport.fullPrefix}${base64URL(new TextEncoder().encode(canonicalJSON(preset, schema)))}`;
}

export function encodeShare(preset, schema) {
  const selection = compactSelection(preset, schema);
  if (!selection) return encodeFullShare(preset, schema);
  const packed =
    BigInt(selection.style) |
    (BigInt(selection.baseColor) << 4n) |
    (BigInt(selection.theme) << 8n) |
    (BigInt(selection.radius) << 13n) |
    (BigInt(selection.menuAccent) << 16n);
  return `${schema.transport.compactPrefix}${encodeBase62(packed)}`;
}

function decodeCompactShare(code, schema, validators) {
  const payload = code.slice(schema.transport.compactPrefix.length);
  const packed = decodeBase62(payload);
  if (packed >> 18n !== 0n) {
    throw new Error("compact payload uses bits beyond the p1 schema");
  }
  const compact = schema.transport.compact;
  const indexes = {
    style: Number(packed & 0xfn),
    baseColor: Number((packed >> 4n) & 0xfn),
    theme: Number((packed >> 8n) & 0x1fn),
    radius: Number((packed >> 13n) & 0x7n),
    menuAccent: Number((packed >> 16n) & 0x3n),
  };
  for (const [name, index] of Object.entries(indexes)) {
    const values =
      name === "style"
        ? compact.styles
        : name === "baseColor"
          ? compact.baseColors
          : name === "theme"
            ? compact.themes
            : name === "radius"
              ? compact.radii
              : compact.menuAccents;
    if (index >= values.length) throw new Error(`compact payload uses unused ${name} index`);
  }

  const style = compact.styles[indexes.style];
  const baseColor = compact.baseColors[indexes.baseColor];
  const theme = compact.themes[indexes.theme];
  const radiusName = compact.radii[indexes.radius];
  const menuAccent = compact.menuAccents[indexes.menuAccent];
  const resolved = schema.palette.resolved[baseColor]?.[theme];
  if (!resolved) {
    throw new Error(`compact payload has invalid base-color/theme combination ${baseColor}/${theme}`);
  }
  const radius = schema.palette.radii.find((choice) => choice.name === radiusName)?.value;
  if (radius === undefined) {
    throw new Error(`compact payload radius ${radiusName} is absent from the palette schema`);
  }
  const light = clone(resolved.light);
  const dark = clone(resolved.dark);
  applyMenuAccent(light, menuAccent);
  applyMenuAccent(dark, menuAccent);
  const preset = validatePreset(
    {
      $schema: schema.schema,
      schemaVersion: schema.schemaVersion,
      style,
      radius,
      theme: { light, dark },
    },
    schema,
    validators,
  );
  if (encodeShare(preset, schema) !== code) {
    throw new Error("compact share code is not canonical");
  }
  return clone(preset);
}

function decodeFullShare(code, schema, validators) {
  const payload = code.slice(schema.transport.fullPrefix.length);
  const source = new TextDecoder("utf-8", { fatal: true }).decode(fromBase64URL(payload));
  const preset = importPresetJSON(source, schema, validators);
  if (encodeFullShare(preset, schema) !== code) {
    throw new Error("share code is not canonical");
  }
  return preset;
}

export function decodeShare(code, schema, validators) {
  if (code.startsWith(schema.transport.compactPrefix)) {
    return decodeCompactShare(code, schema, validators);
  }
  if (code.startsWith(schema.transport.fullPrefix)) {
    return decodeFullShare(code, schema, validators);
  }
  throw new Error("unsupported share transport");
}

export function loadShareFromURL(url, schema, validators) {
  const code = new URL(url).searchParams.get("preset");
  return code ? decodeShare(code, schema, validators) : null;
}

export function commandStrings(share) {
  return {
    init: `gsxui init --preset '${share}'`,
    apply: `gsxui apply --preset '${share}'`,
  };
}

function objectEntries(node, allowed, required, path) {
  if (!node || node.type !== "object") {
    throw new Error(`${path} must be an object`);
  }
  const entries = new Map();
  for (const property of node.children ?? []) {
    const [keyNode, valueNode] = property.children ?? [];
    const key = keyNode?.value;
    if (typeof key !== "string") {
      throw new Error(`${path} contains an invalid property`);
    }
    if (entries.has(key)) {
      throw new Error(`${path}.${key} is duplicated`);
    }
    if (!allowed.includes(key)) {
      throw new Error(`${path}.${key} is unknown`);
    }
    entries.set(key, valueNode);
  }
  for (const key of required) {
    if (!entries.has(key)) {
      throw new Error(`${path}.${key} is required`);
    }
  }
  return entries;
}

export function importPresetJSON(source, schema, validators) {
  const errors = [];
  const root = parseTree(source, errors, {
    allowTrailingComma: false,
    disallowComments: true,
  });
  if (!root || errors.length !== 0) {
    const offset = errors[0]?.offset ?? 0;
    throw new Error(`invalid preset JSON at offset ${offset}`);
  }

  const presetEntries = objectEntries(
    root,
    ["$schema", "schemaVersion", "style", "radius", "theme"],
    ["$schema", "schemaVersion", "style", "radius", "theme"],
    "preset",
  );
  const themeEntries = objectEntries(
    presetEntries.get("theme"),
    ["light", "dark"],
    ["light", "dark"],
    "preset.theme",
  );
  for (const mode of ["light", "dark"]) {
    objectEntries(
      themeEntries.get(mode),
      schema.tokenNames,
      schema.tokenNames,
      `preset.theme.${mode}`,
    );
  }

  return clone(validatePreset(getNodeValue(root), schema, validators));
}

export function importThemeCSS(state, source, schema, validators, parseCSS) {
  const parsed = parseCSS(source, schema);
  const candidate = clone(state.resolved);
  if (parsed.radius) candidate.radius = parsed.radius;
  Object.assign(candidate.theme.light, parsed.light);
  Object.assign(candidate.theme.dark, parsed.dark);
  validatePreset(candidate, schema, validators);
  return replacePreset(state, candidate, schema);
}

export function manualCopyState(copied, text) {
  return copied ? { visible: false, text: "" } : { visible: true, text };
}
