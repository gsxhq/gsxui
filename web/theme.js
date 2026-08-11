import { init, on } from "../ui/gsxui.js";
import { parse as parseCSSDeclaration, tokenize, tokenTypes } from "css-tree";
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

// The editor wires per pageview. setup() captures this page's elements and
// state in one closure and returns the event handlers; the module-scope
// registrations below forward to the CURRENT pageview's closure. A module
// used to scan for [data-theme-schema] once at load, which left the editor
// dead when boosted navigation morphed the theme page in AFTER that scan —
// init() (the shared MutationObserver contract in ui/gsxui.js) re-runs the
// build for every fresh schema element instead, and the forwards never
// stack because they are registered exactly once here.
let editor = null;

init("[data-theme-schema]", (schemaElement) => {
  if (editor?.schemaElement === schemaElement) return;
  editor?.dispose();
  editor = setup(schemaElement);
});

on("change", "[data-theme-picker] [data-gsxui-slot-radio]", (e, el) => editor?.onPickerChange(e, el));
on("pointerenter", "[data-theme-picker] [data-theme-choice]", (e, el) => editor?.onChoiceEnter(e, el), { capture: true });
on("pointerleave", "[data-theme-picker]", (e, el) => editor?.onPickerLeave(e, el), { capture: true });
on("gsxui:open", "[data-theme-picker] [data-gsxui-slot-popover-content]", (e, el) => editor?.onPickerOpen(e, el));
on("toggle", "[data-theme-picker] [data-gsxui-slot-popover-content]", (e, el) => editor?.onPickerToggle(e, el), { capture: true });
on("click", "[data-theme-style]", (e, el) => editor?.onStyleClick(e, el));
on("click", "[data-theme-mode-tab]", (e, el) => editor?.onModeClick(e, el));
on("click", "[data-theme-reset]", (e, el) => editor?.onReset(e, el));
on("click", "[data-theme-copy]", (e, el) => editor?.onCopy(e, el));
on("click", "[data-theme-download]", (e, el) => editor?.onDownload(e, el));
on("click", "[data-theme-import-apply]", (e, el) => editor?.onImportApply(e, el));
on("click", "[data-theme-preview-retry]", (e, el) => editor?.onPreviewRetry(e, el));
addEventListener("message", (e) => editor?.onMessage(e));

function setup(schemaElement) {
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
      const styleTitle =
        state.resolved.style.charAt(0).toUpperCase() + state.resolved.style.slice(1);
      setStatus(`${styleTitle} · ${state.mode}`);
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

    function themeMode(selector) {
      const tokens = [];
      tokenize(selector, (type, start, end) => {
        if (type === tokenTypes.WhiteSpace || type === tokenTypes.Comment) return;
        tokens.push({ type, value: selector.slice(start, end) });
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
      const value = declaration.value.value.trim();
      if (!declaration.important) return value;
      const suffix = declaration.important === true ? "important" : declaration.important;
      return `${value} !${suffix}`;
    }

    function acceptDeclaration(declaration, mode, nested) {
      const name = recognizedName(declaration.property);
      if (!name) return;
      if (nested) {
        throw new Error(
          `Recognized property ${declaration.property} has nested declaration ownership.`,
        );
      }
      if (!mode) {
        throw new Error(
          `Recognized property ${declaration.property} must belong to :root or .dark.`,
        );
      }
      if (seen[mode].has(name)) {
        throw new Error(`${declaration.property} is duplicated in ${mode}.`);
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

    function parseDeclaration(text) {
      try {
        return parseCSSDeclaration(text, {
          context: "declaration",
          parseValue: false,
          onParseError(error) {
            throw error;
          },
        });
      } catch (error) {
        const message = error instanceof Error ? error.message : "invalid syntax";
        throw new Error(`CSS import: malformed CSS: ${message}`);
      }
    }

    // Splits text into statements per the CSS syntax spec: a statement ends
    // at a top-level semicolon (block: null) or at the close of a top-level
    // {} block. validateCSSStructure has already guaranteed balanced blocks
    // and strings, so the depth counters cannot go negative or stay open.
    function splitStatements(text) {
      const statements = [];
      let group = 0;
      let curly = 0;
      let start = -1;
      let preludeEnd = -1;
      let blockStart = -1;
      tokenize(text, (type, tokenStart, tokenEnd) => {
        if (start === -1) {
          if (
            type === tokenTypes.WhiteSpace ||
            type === tokenTypes.Comment ||
            type === tokenTypes.CDO ||
            type === tokenTypes.CDC
          ) return;
          start = tokenStart;
        }
        switch (type) {
          case tokenTypes.Function:
          case tokenTypes.LeftParenthesis:
          case tokenTypes.LeftSquareBracket:
            group++;
            break;
          case tokenTypes.RightParenthesis:
          case tokenTypes.RightSquareBracket:
            group--;
            break;
          case tokenTypes.LeftCurlyBracket:
            if (curly === 0 && group === 0) {
              preludeEnd = tokenStart;
              blockStart = tokenEnd;
            }
            curly++;
            break;
          case tokenTypes.RightCurlyBracket:
            curly--;
            if (curly === 0 && group === 0) {
              statements.push({
                prelude: text.slice(start, preludeEnd),
                block: text.slice(blockStart, tokenStart),
              });
              start = -1;
            }
            break;
          case tokenTypes.Semicolon:
            if (curly === 0 && group === 0) {
              statements.push({ prelude: text.slice(start, tokenStart), block: null });
              start = -1;
            }
            break;
        }
      });
      if (start !== -1) statements.push({ prelude: text.slice(start), block: null });
      return statements;
    }

    // A block statement whose prelude is `--name :` is a custom property
    // declaration with a {} block in its value, not a nested rule.
    function isCustomPropertyPrelude(prelude) {
      const significant = [];
      tokenize(prelude, (type, start, end) => {
        if (type === tokenTypes.WhiteSpace || type === tokenTypes.Comment) return;
        if (significant.length < 2) significant.push({ type, value: prelude.slice(start, end) });
      });
      return (
        significant.length === 2 &&
        significant[0].type === tokenTypes.Ident &&
        significant[0].value.startsWith("--") &&
        significant[1].type === tokenTypes.Colon
      );
    }

    function walkBlock(text, mode, nested) {
      for (const statement of splitStatements(text)) {
        const prelude = statement.prelude.trim();
        if (statement.block === null) {
          if (prelude === "" || prelude.startsWith("@")) continue;
          acceptDeclaration(parseDeclaration(statement.prelude), mode, nested);
          continue;
        }
        if (!prelude.startsWith("@") && isCustomPropertyPrelude(statement.prelude)) {
          acceptDeclaration(
            parseDeclaration(`${statement.prelude}{${statement.block}}`),
            mode,
            nested,
          );
          continue;
        }
        walkBlock(statement.block, "", true);
      }
    }

    for (const statement of splitStatements(source)) {
      const prelude = statement.prelude.trim();
      if (statement.block === null) {
        if (prelude === "" || prelude.startsWith("@")) continue;
        throw new Error("CSS import: malformed CSS: declaration outside a rule");
      }
      if (prelude.startsWith("@")) {
        walkBlock(statement.block, "", true);
        continue;
      }
      walkBlock(statement.block, themeMode(statement.prelude), false);
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

  function onPickerChange(_event, input) {
    const kind = input.closest("[data-theme-picker]").dataset.themePicker;
    if (kind === "baseColor") state = selectBaseColor(state, input.value, schema);
    else if (kind === "theme") state = selectTheme(state, input.value, schema);
    else state = selectRadius(state, input.value, schema);
    render();
  }

  function onChoiceEnter(event, choice) {
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
  }

  function onPickerLeave(event, picker) {
    if (event.relatedTarget instanceof Node && picker.contains(event.relatedTarget)) return;
    state = clearPalettePreview(state);
    render();
  }

  function onPickerOpen(_event, content) {
    content.querySelector("[data-gsxui-slot-radio]:checked")?.focus();
  }

  function onPickerToggle(event) {
    if (event.newState === "closed") {
      state = clearPalettePreview(state);
      render();
    }
  }

  function onStyleClick(_event, button) {
    state = selectStyle(state, button.dataset.themeStyle, schema);
    render();
  }

  function onModeClick(_event, button) {
    state = selectMode(state, button.dataset.themeModeTab);
    render();
  }

  function onReset() {
    state = resetThemeState(state, schema);
    render();
  }

  async function onCopy(_event, button) {
    const artifacts = currentArtifacts();
    const text = artifacts[button.dataset.themeCopy];
    const copied = await copyText(text);
    showManualCopy(copied, text);
    if (copied) setStatus("Copied.");
  }

  function onDownload(_event, button) {
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
  }

  function onImportApply(_event, button) {
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
  }

  function onPreviewRetry() {
    if (!frame) return;
    startPreviewHandshake("Reconnecting to preview…");
    frame.src = frame.src;
  }

  function onMessage(event) {
    // A stale pageview's closure must ignore messages after boosted
    // navigation away — its frame is disconnected, and a null
    // event.source would otherwise compare equal to its null
    // contentWindow.
    if (!frame?.isConnected) return;
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
  }

  frame?.addEventListener("load", () => {
    startPreviewHandshake();
    syncPreview();
  });

  startPreviewHandshake();
  render();

  return {
    schemaElement,
    // dispose stops the pending handshake timer so a replaced pageview's
    // timeout cannot write "Preview did not respond." into the new page.
    dispose: finishPreviewHandshake,
    onPickerChange,
    onChoiceEnter,
    onPickerLeave,
    onPickerOpen,
    onPickerToggle,
    onStyleClick,
    onModeClick,
    onReset,
    onCopy,
    onDownload,
    onImportApply,
    onPreviewRetry,
    onMessage,
  };
}
