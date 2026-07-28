import { on } from "../ui/gsxui.js";
import {
  applyField,
  canonicalJSON,
  commandStrings,
  createThemeState,
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

const PREVIEW_MESSAGE = "gsxui:theme-preview:v1";
const READY_MESSAGE = "gsxui:theme-preview-ready:v1";
const ERROR_MESSAGE = "gsxui:theme-preview-error:v1";

const schemaElement = document.querySelector("[data-theme-schema]");

if (schemaElement) {
  const schema = JSON.parse(schemaElement.textContent);
  const frame = document.querySelector("[data-theme-preview-frame]");
  const status = document.querySelector("[data-theme-status]");
  const previewStatus = document.querySelector("[data-theme-preview-status]");
  const previewRetry = document.querySelector("[data-theme-preview-retry]");
  const manualCopy = document.querySelector("[data-theme-manual-copy]");
  const colorProbe = document.createElement("span").style;

  const validators = {
    color(value) {
      if (typeof value !== "string" || value.trim() === "") return false;
      colorProbe.color = "";
      colorProbe.color = value;
      return colorProbe.color !== "";
    },
    radius(value) {
      if (typeof value !== "string" || value.trim() === "") return false;
      if (!globalThis.CSSNumericValue || !globalThis.CSSUnitValue) return false;
      try {
        const parsed = CSSNumericValue.parse(value);
        if (!(parsed instanceof CSSUnitValue) || parsed.value < 0) return false;
        if (parsed.unit === "number") return parsed.value === 0;
        return schema.radiusUnits.includes(parsed.unit);
      } catch {
        return false;
      }
    },
  };

  let state;
  let initialMessage = "";
  try {
    const shared = loadShareFromURL(location.href, schema, validators);
    state = createThemeState(schema, shared ?? schema.defaults.nova);
    if (shared) initialMessage = "Loaded shared preset.";
  } catch (error) {
    state = createThemeState(schema);
    initialMessage = error instanceof Error ? `Share link ignored: ${error.message}` : "Share link ignored.";
  }

  function fieldValue(path) {
    if (state.drafts.has(path)) return state.drafts.get(path);
    if (path === "radius") return state.resolved.radius;
    const [mode, name] = path.split(".");
    return state.resolved.theme[mode][name];
  }

  function setStatus(message, error = false) {
    if (!status) return;
    status.textContent = message;
    status.classList.toggle("text-destructive", error);
    status.classList.toggle("text-muted-foreground", !error);
  }

  function previewPayload() {
    return {
      type: PREVIEW_MESSAGE,
      preset: state.resolved,
      mode: state.mode,
    };
  }

  function syncPreview() {
    frame?.contentWindow?.postMessage(previewPayload(), location.origin);
  }

  function currentArtifacts() {
    const json = canonicalJSON(state.resolved, schema);
    const css = themeCSS(state.resolved, schema);
    const share = encodeShare(state.resolved, schema);
    const shareURL = new URL(location.href);
    shareURL.search = "";
    shareURL.hash = "";
    shareURL.searchParams.set("preset", share);
    return {
      json,
      css,
      share,
      url: shareURL.toString(),
      ...commandStrings(share),
    };
  }

  function render() {
    for (const button of document.querySelectorAll("[data-theme-style]")) {
      const active = button.dataset.themeStyle === state.resolved.style;
      button.setAttribute("aria-pressed", String(active));
      button.classList.toggle("border-primary", active);
      button.classList.toggle("border-border", !active);
      button.classList.toggle("bg-accent/50", active);
    }
    for (const button of document.querySelectorAll("[data-theme-mode-tab]")) {
      const active = button.dataset.themeModeTab === state.mode;
      button.setAttribute("aria-pressed", String(active));
      button.classList.toggle("bg-accent", active);
      button.classList.toggle("text-accent-foreground", active);
      button.classList.toggle("text-muted-foreground", !active);
    }
    for (const panel of document.querySelectorAll("[data-theme-mode-field]")) {
      panel.hidden = panel.dataset.themeModeField !== state.mode;
    }
    for (const input of document.querySelectorAll("[data-theme-field]")) {
      const path = input.dataset.themeField;
      input.value = fieldValue(path);
      const invalid = state.drafts.has(path);
      input.setAttribute("aria-invalid", String(invalid));
      const error = document.querySelector(`[data-theme-field-error="${CSS.escape(path)}"]`);
      if (error) {
        error.hidden = !invalid;
        error.classList.toggle("hidden", !invalid);
        error.textContent = invalid
          ? path === "radius"
            ? "Use zero or a non-negative CSS length."
            : "Enter a valid CSS color."
          : "";
      }
    }

    const artifacts = currentArtifacts();
    for (const output of document.querySelectorAll("[data-theme-command]")) {
      output.value = artifacts[output.dataset.themeCommand];
    }
    if (state.drafts.size > 0) {
      setStatus(`${state.drafts.size} invalid draft${state.drafts.size === 1 ? "" : "s"} — preview and exports keep the last valid values.`, true);
    } else if (initialMessage) {
      setStatus(initialMessage);
      initialMessage = "";
    } else {
      setStatus(`${state.resolved.style === "maia" ? "Maia" : "Nova"} · ${state.mode}`);
    }
    syncPreview();
  }

  function parseCompatibleThemeCSS(source) {
    const sheet = new CSSStyleSheet();
    sheet.replaceSync(source);
    const result = { light: {}, dark: {}, radius: "" };
    const seen = { light: new Set(), dark: new Set(), radius: false };
    const tokenSet = new Set(schema.tokenNames);
    let recognized = 0;

    for (const rule of sheet.cssRules) {
      if (!(rule instanceof CSSStyleRule)) {
        throw new Error("Only plain :root and .dark rules are supported.");
      }
      const mode = rule.selectorText === ":root" ? "light" : rule.selectorText === ".dark" ? "dark" : "";
      for (const property of rule.style) {
        if (!property.startsWith("--")) continue;
        const name = property.slice(2);
        if (name !== "radius" && !tokenSet.has(name)) continue;
        if (!mode) {
          throw new Error(`Recognized property ${property} must belong to :root or .dark.`);
        }
        const value = rule.style.getPropertyValue(property).trim();
        recognized++;
        if (name === "radius") {
          if (mode !== "light") throw new Error("--radius belongs only to :root.");
          if (seen.radius) throw new Error("--radius is duplicated.");
          seen.radius = true;
          result.radius = value;
          continue;
        }
        if (seen[mode].has(name)) throw new Error(`${property} is duplicated in ${rule.selectorText}.`);
        seen[mode].add(name);
        result[mode][name] = value;
      }
    }

    if (recognized === 0) throw new Error("No supported theme properties were found.");
    return result;
  }

  async function copyText(text) {
    try {
      await navigator.clipboard.writeText(text);
      return true;
    } catch {
      return false;
    }
  }

  function showManualCopy(copied, text) {
    const fallback = manualCopyState(copied, text);
    if (!manualCopy) return;
    manualCopy.hidden = !fallback.visible;
    manualCopy.classList.toggle("hidden", !fallback.visible);
    manualCopy.value = fallback.text;
    if (fallback.visible) {
      manualCopy.focus();
      manualCopy.select();
      setStatus("Clipboard unavailable. Copy the selected text manually.", true);
    }
  }

  on("input", "[data-theme-field]", (_event, input) => {
    state = applyField(state, input.dataset.themeField, input.value, validators);
    render();
  });

  on("click", "[data-theme-style]", (_event, button) => {
    state = selectStyle(state, button.dataset.themeStyle, schema);
    render();
  });

  on("click", "[data-theme-mode-tab]", (_event, button) => {
    state = selectMode(state, button.dataset.themeModeTab);
    render();
  });

  on("click", "[data-theme-reset]", () => {
    state = resetThemeState(state, schema);
    render();
  });

  on("click", "[data-theme-copy]", async (_event, button) => {
    const artifacts = currentArtifacts();
    const text = artifacts[button.dataset.themeCopy];
    const copied = await copyText(text);
    showManualCopy(copied, text);
    if (copied) setStatus("Copied.");
  });

  on("click", "[data-theme-download]", (_event, button) => {
    const kind = button.dataset.themeDownload;
    const artifacts = currentArtifacts();
    const text = artifacts[kind];
    const blob = new Blob([text], {
      type: kind === "json" ? "application/json" : "text/css",
    });
    const href = URL.createObjectURL(blob);
    const anchor = document.createElement("a");
    anchor.href = href;
    anchor.download = kind === "json" ? "preset.json" : "theme.css";
    anchor.click();
    URL.revokeObjectURL(href);
  });

  on("click", "[data-theme-import-apply]", (_event, button) => {
    const kind = button.dataset.themeImportApply;
    const input = document.querySelector(`[data-theme-import="${kind}"]`);
    try {
      if (kind === "json") {
        const preset = importPresetJSON(input.value, schema, validators);
        state = { resolved: preset, drafts: new Map(), mode: state.mode };
      } else {
        state = importThemeCSS(state, input.value, schema, validators, parseCompatibleThemeCSS);
      }
      setStatus(`Applied ${kind.toUpperCase()} import.`);
      render();
    } catch (error) {
      setStatus(error instanceof Error ? error.message : `Could not import ${kind}.`, true);
    }
  });

  on("click", "[data-theme-preview-retry]", () => {
    if (!frame) return;
    previewRetry?.classList.add("hidden");
    if (previewStatus) previewStatus.textContent = "Reconnecting to preview…";
    frame.src = frame.src;
  });

  addEventListener("message", (event) => {
    if (event.origin !== location.origin || event.source !== frame?.contentWindow) return;
    if (event.data?.type === READY_MESSAGE) {
      if (previewStatus) previewStatus.textContent = "Live";
      previewRetry?.classList.add("hidden");
      syncPreview();
    } else if (event.data?.type === ERROR_MESSAGE) {
      if (previewStatus) previewStatus.textContent = event.data.message || "Preview rejected the current state.";
      previewRetry?.classList.remove("hidden");
    }
  });

  frame?.addEventListener("load", () => {
    if (previewStatus) previewStatus.textContent = "Connecting to preview…";
    syncPreview();
  });

  render();
}
