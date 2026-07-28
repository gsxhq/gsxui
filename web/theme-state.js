import { getNodeValue, parseTree } from "jsonc-parser";

function clone(value) {
  return structuredClone(value);
}

function cloneState(state) {
  return {
    resolved: clone(state.resolved),
    drafts: new Map(state.drafts),
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
  return {
    resolved: clone(preset),
    drafts: new Map(),
    mode: "light",
  };
}

export function selectStyle(state, style, schema) {
  assertStyle(style, schema);
  const next = cloneState(state);
  next.resolved.style = style;
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

export function applyField(state, path, value, validators) {
  const next = cloneState(state);
  const valid =
    path === "radius"
      ? validators.radius(value)
      : /^(light|dark)\.[a-z0-9-]+$/.test(path) && validators.color(value);
  if (!valid) {
    next.drafts.set(path, value);
    return next;
  }

  if (path === "radius") {
    next.resolved.radius = value;
  } else {
    const [mode, name] = path.split(".");
    if (!(name in next.resolved.theme[mode])) {
      throw new Error(`unknown theme field ${path}`);
    }
    next.resolved.theme[mode][name] = value;
  }
  next.drafts.delete(path);
  return next;
}

export function resetThemeState(state, schema) {
  assertStyle(state.resolved.style, schema);
  return {
    resolved: clone(schema.defaults[state.resolved.style]),
    drafts: new Map(),
    mode: state.mode,
  };
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

export function encodeShare(preset, schema) {
  return `${schema.transportPrefix}${base64URL(new TextEncoder().encode(canonicalJSON(preset, schema)))}`;
}

export function decodeShare(code, schema, validators) {
  if (!code.startsWith(schema.transportPrefix)) {
    throw new Error("unsupported share transport");
  }
  const payload = code.slice(schema.transportPrefix.length);
  const source = new TextDecoder("utf-8", { fatal: true }).decode(fromBase64URL(payload));
  const preset = importPresetJSON(source, schema, validators);
  if (encodeShare(preset, schema) !== code) {
    throw new Error("share code is not canonical");
  }
  return preset;
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
  return {
    resolved: candidate,
    drafts: new Map(),
    mode: state.mode,
  };
}

export function manualCopyState(copied, text) {
  return copied ? { visible: false, text: "" } : { visible: true, text };
}
