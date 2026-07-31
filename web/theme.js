import { on } from "../ui/gsxui.js";
import { tokenize, tokenTypes } from "css-tree";
import parseCSS from "postcss/lib/parse";
import {
  canonicalJSON,
  clearPalettePreview,
  commandStrings,
  createThemeState,
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
  selectMode,
  selectBaseColor,
  selectRadius,
  selectStyle,
  selectTheme,
  themeCSS,
} from "./theme-state.js";

const PREVIEW_MESSAGE = "gsxui:theme-preview:v1";
const READY_MESSAGE = "gsxui:theme-preview-ready:v1";
const APPLIED_MESSAGE = "gsxui:theme-preview-applied:v1";
const ERROR_MESSAGE = "gsxui:theme-preview-error:v1";
const PREVIEW_HANDSHAKE_TIMEOUT_MS = 2_000;

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
  let previewHandshakeTimer;
  let previewAttempt = "";
  let initialMessage = "";
  try {
    const shared = loadShareFromURL(location.href, schema, validators);
    state = createThemeState(schema, shared ?? schema.defaults.nova);
    if (shared) initialMessage = "Loaded shared preset.";
  } catch (error) {
    state = createThemeState(schema);
    initialMessage = error instanceof Error ? `Share link ignored: ${error.message}` : "Share link ignored.";
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
      attempt: previewAttempt,
      preset: previewPreset(state),
      mode: state.mode,
    };
  }

  function syncPreview() {
    frame?.contentWindow?.postMessage(previewPayload(), location.origin);
  }

  function finishPreviewHandshake() {
    if (previewHandshakeTimer === undefined) return;
    clearTimeout(previewHandshakeTimer);
    previewHandshakeTimer = undefined;
  }

  function startPreviewHandshake(message = "Connecting to preview…") {
    finishPreviewHandshake();
    previewAttempt = crypto.randomUUID();
    if (previewStatus) previewStatus.textContent = message;
    previewRetry?.classList.add("hidden");
    previewHandshakeTimer = setTimeout(() => {
      previewHandshakeTimer = undefined;
      if (previewStatus) previewStatus.textContent = "Preview did not respond.";
      previewRetry?.classList.remove("hidden");
    }, PREVIEW_HANDSHAKE_TIMEOUT_MS);
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
    renderPickers();

    const artifacts = currentArtifacts();
    for (const output of document.querySelectorAll("[data-theme-command]")) {
      output.value = artifacts[output.dataset.themeCommand];
    }
    if (initialMessage) {
      setStatus(initialMessage);
      initialMessage = "";
    } else {
      setStatus(`${state.resolved.style === "maia" ? "Maia" : "Nova"} · ${state.mode}`);
    }
    syncPreview();
  }

  function pickerChoices(kind) {
    if (kind === "baseColor") return schema.palette.baseColors;
    if (kind === "radius") return schema.palette.radii;
    return schema.palette.themes[state.selection.baseColor === "custom" ? "neutral" : state.selection.baseColor];
  }

  function renderPickers() {
    function renderSwatch(swatch, choice, kind) {
      const radius = kind === "radius";
      swatch.hidden = !choice || (!radius && !choice.swatch);
      swatch.style.backgroundColor = radius ? "" : (choice?.swatch ?? "");
      swatch.style.borderRadius = radius ? (choice?.value ?? "") : "";
    }

    for (const kind of ["baseColor", "theme", "radius"]) {
      const picker = document.querySelector(`[data-theme-picker="${kind}"]`);
      if (!picker) continue;
      const choices = pickerChoices(kind);
      const selected = state.selection[kind];
      const selectedChoice = choices.find((choice) => choice.name === selected);
      picker.querySelector("[data-theme-selection-value]").textContent = selectedChoice?.title ?? "Custom";
      const swatch = picker.querySelector("[data-theme-selection-swatch]");
      renderSwatch(swatch, selectedChoice, kind);
      const options = picker.querySelectorAll("[data-theme-choice]");
      choices.forEach((choice, index) => {
        const option = options[index];
        if (!option) return;
        const input = option.querySelector("input");
        input.value = choice.name;
        input.checked = choice.name === selected;
        input.setAttribute("aria-label", choice.title);
        option.querySelector("[data-theme-choice-label]").textContent = choice.title;
        renderSwatch(option.querySelector("[data-theme-choice-swatch]"), choice, kind);
      });
    }
  }

  function parseCompatibleThemeCSS(source) {
    validateCSSStructure(source);
    const tokenSet = new Set(schema.tokenNames);
    const seen = { light: new Set(), dark: new Set() };
    const result = { light: {}, dark: {}, radius: "" };
    let radiusValue = "";

    let stylesheet;
    try {
      stylesheet = parseCSS(source);
    } catch (error) {
      const message = error instanceof Error ? error.message : "invalid syntax";
      throw new Error(`CSS import: malformed CSS: ${message}`);
    }

    function themeMode(rule) {
      if (typeof rule.selector !== "string") return "";
      const tokens = [];
      tokenize(rule.selector, (type, start, end) => {
        if (type === tokenTypes.WhiteSpace || type === tokenTypes.Comment) return;
        tokens.push({ type, value: rule.selector.slice(start, end) });
      });
      if (
        tokens.length === 2 &&
        tokens[0].type === tokenTypes.Colon &&
        tokens[1].type === tokenTypes.Ident &&
        tokens[1].value.toLowerCase() === "root"
      ) {
        return "light";
      }
      if (
        tokens.length === 2 &&
        tokens[0].type === tokenTypes.Delim &&
        tokens[0].value === "." &&
        tokens[1].type === tokenTypes.Ident &&
        tokens[1].value === "dark"
      ) {
        return "dark";
      }
      return "";
    }

    function recognizedName(property) {
      if (!property.startsWith("--")) return "";
      const name = property.slice(2);
      return name === "radius" || tokenSet.has(name) ? name : "";
    }

    function declarationValue(declaration) {
      const value = declaration.value.trim();
      if (!declaration.important) return value;
      return `${value}${declaration.raws.important || " !important"}`.trim();
    }

    function acceptDeclaration(declaration, mode, nested) {
      const name = recognizedName(declaration.prop);
      if (!name) return;
      if (nested) {
        throw new Error(
          `Recognized property ${declaration.prop} has nested declaration ownership.`,
        );
      }
      if (!mode) {
        throw new Error(
          `Recognized property ${declaration.prop} must belong to :root or .dark.`,
        );
      }
      if (seen[mode].has(name)) {
        throw new Error(`${declaration.prop} is duplicated in ${mode}.`);
      }
      seen[mode].add(name);
      const value = declarationValue(declaration);
      if (name !== "radius") {
        result[mode][name] = value;
        return;
      }
      if (radiusValue && radiusValue !== value) {
        throw new Error(
          `Conflicting radius declarations ${JSON.stringify(radiusValue)} and ${JSON.stringify(value)}.`,
        );
      }
      radiusValue = value;
      result.radius = value;
    }

    function rejectNestedRecognized(node) {
      if (node.type === "decl") {
        acceptDeclaration(node, "", true);
        return;
      }
      if (node.nodes) {
        for (const child of node.nodes) rejectNestedRecognized(child);
      }
    }

    for (const node of stylesheet.nodes) {
      if (node.type !== "rule") {
        rejectNestedRecognized(node);
        continue;
      }
      const mode = themeMode(node);
      for (const child of node.nodes) {
        if (child.type === "decl") {
          acceptDeclaration(child, mode, false);
        } else {
          rejectNestedRecognized(child);
        }
      }
    }

    const sheet = new CSSStyleSheet();
    sheet.replaceSync(source);
    return result;
  }

  function validateCSSStructure(source) {
    const blocks = [];

    function push(expected) {
      blocks.push(expected);
    }

    function pop(expected, name) {
      if (blocks.at(-1) !== expected) {
        throw new Error(`CSS import: malformed CSS: unexpected ${name}`);
      }
      blocks.pop();
    }

    function endsWithUnescaped(sourceText, start, end, character) {
      if (end <= start || sourceText[end - 1] !== character) return false;
      let backslashes = 0;
      for (let offset = end - 2; offset >= start && sourceText[offset] === "\\"; offset--) {
        backslashes++;
      }
      return backslashes % 2 === 0;
    }

    tokenize(source, (type, start, end) => {
      switch (type) {
        case tokenTypes.BadString:
          throw new Error("CSS import: malformed CSS: invalid string");
        case tokenTypes.BadUrl:
          throw new Error("CSS import: malformed CSS: invalid URL");
        case tokenTypes.String:
          if (!endsWithUnescaped(source, start, end, source[start])) {
            throw new Error("CSS import: malformed CSS: unclosed string");
          }
          break;
        case tokenTypes.Url:
          if (!endsWithUnescaped(source, start, end, ")")) {
            throw new Error("CSS import: malformed CSS: unclosed URL");
          }
          break;
        case tokenTypes.Comment:
          if (source.slice(end - 2, end) !== "*/") {
            throw new Error("CSS import: malformed CSS: unclosed comment");
          }
          break;
        case tokenTypes.Function:
        case tokenTypes.LeftParenthesis:
          push(tokenTypes.RightParenthesis);
          break;
        case tokenTypes.LeftSquareBracket:
          push(tokenTypes.RightSquareBracket);
          break;
        case tokenTypes.LeftCurlyBracket:
          push(tokenTypes.RightCurlyBracket);
          break;
        case tokenTypes.RightParenthesis:
          pop(type, "right parenthesis");
          break;
        case tokenTypes.RightSquareBracket:
          pop(type, "right bracket");
          break;
        case tokenTypes.RightCurlyBracket:
          pop(type, "right brace");
          break;
      }
    });

    if (blocks.length > 0) {
      throw new Error("CSS import: malformed CSS: unclosed block");
    }
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

  on("change", "[data-theme-picker] [data-gsxui-slot-radio]", (_event, input) => {
    const kind = input.closest("[data-theme-picker]").dataset.themePicker;
    if (kind === "baseColor") state = selectBaseColor(state, input.value, schema);
    else if (kind === "theme") state = selectTheme(state, input.value, schema);
    else state = selectRadius(state, input.value, schema);
    render();
  });

  on("pointerenter", "[data-theme-picker] [data-theme-choice]", (event, choice) => {
    if (
      event.pointerType === "touch" ||
      !matchMedia("(hover: hover) and (pointer: fine)").matches
    ) return;
    const input = choice.querySelector("[data-gsxui-slot-radio]");
    const kind = input.closest("[data-theme-picker]").dataset.themePicker;
    if (kind === "baseColor") state = previewBaseColor(state, input.value, schema);
    else if (kind === "theme") state = previewTheme(state, input.value, schema);
    else state = previewRadius(state, input.value, schema);
    render();
  }, { capture: true });
  on("pointerleave", "[data-theme-picker]", (event, picker) => {
    if (event.relatedTarget instanceof Node && picker.contains(event.relatedTarget)) return;
    state = clearPalettePreview(state);
    render();
  }, { capture: true });
  on(
    "gsxui:open",
    "[data-theme-picker] [data-gsxui-slot-popover-content]",
    (_event, content) => {
      content.querySelector("[data-gsxui-slot-radio]:checked")?.focus();
    },
  );
  on("toggle", "[data-theme-picker] [data-gsxui-slot-popover-content]", (event) => {
    if (event.newState === "closed") {
      state = clearPalettePreview(state);
      render();
    }
  }, { capture: true });

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
        state = replacePreset(state, preset, schema);
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
    startPreviewHandshake("Reconnecting to preview…");
    frame.src = frame.src;
  });

  addEventListener("message", (event) => {
    if (event.origin !== location.origin || event.source !== frame?.contentWindow) return;
    if (event.data?.type === READY_MESSAGE) {
      syncPreview();
    } else if (event.data?.type === APPLIED_MESSAGE) {
      if (event.data.attempt !== previewAttempt) return;
      finishPreviewHandshake();
      if (previewStatus) previewStatus.textContent = "Live";
      previewRetry?.classList.add("hidden");
    } else if (event.data?.type === ERROR_MESSAGE) {
      if (event.data.attempt !== previewAttempt) return;
      finishPreviewHandshake();
      if (previewStatus) previewStatus.textContent = event.data.message || "Preview rejected the current state.";
      previewRetry?.classList.remove("hidden");
    }
  });

  frame?.addEventListener("load", () => {
    startPreviewHandshake();
    syncPreview();
  });

  startPreviewHandshake();
  render();
}
