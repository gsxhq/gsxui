const PREVIEW_MESSAGE = "gsxui:theme-preview:v1";
const READY_MESSAGE = "gsxui:theme-preview-ready:v1";
const ERROR_MESSAGE = "gsxui:theme-preview-error:v1";

const previewDocument = document.querySelector("[data-theme-button-preview]");

if (previewDocument) {
  const root = document.documentElement;
  const tokenNames = (previewDocument.dataset.themePreviewTokens ?? "")
    .split(",")
    .filter(Boolean);
  const styleSections = new Map(
    [...document.querySelectorAll("[data-theme-preview-style]")].map((section) => [
      section.dataset.themePreviewStyle,
      section,
    ]),
  );
  let hasAppliedState = false;
  let focusedStyle = "";

  function fail(message) {
    if (parent !== window) {
      parent.postMessage({ type: ERROR_MESSAGE, message }, location.origin);
    }
  }

  function ready() {
    if (parent !== window) {
      parent.postMessage({ type: READY_MESSAGE }, location.origin);
    }
  }

  function exactTokenMap(value, path) {
    if (!value || typeof value !== "object" || Array.isArray(value)) {
      throw new Error(`${path} must be an object`);
    }
    const keys = Object.keys(value).sort();
    const expected = [...tokenNames].sort();
    if (keys.length !== expected.length || keys.some((key, index) => key !== expected[index])) {
      throw new Error(`${path} must contain every supported token exactly once`);
    }

    const colorProbe = document.createElement("span").style;
    for (const name of tokenNames) {
      const candidate = value[name];
      if (typeof candidate !== "string" || candidate.trim() === "") {
        throw new Error(`${path}.${name} must be a non-empty CSS color`);
      }
      colorProbe.color = "";
      colorProbe.color = candidate;
      if (colorProbe.color === "") {
        throw new Error(`${path}.${name} is not a valid CSS color`);
      }
    }
    return value;
  }

  function validatedState(message) {
    if (!message || message.type !== PREVIEW_MESSAGE) {
      throw new Error("unsupported preview message");
    }
    if (message.mode !== "light" && message.mode !== "dark") {
      throw new Error("mode must be light or dark");
    }
    const preset = message.preset;
    if (!preset || typeof preset !== "object" || Array.isArray(preset)) {
      throw new Error("preset must be an object");
    }
    if (!styleSections.has(preset.style)) {
      throw new Error(`unsupported style ${String(preset.style)}`);
    }
    if (typeof preset.radius !== "string" || preset.radius.trim() === "") {
      throw new Error("radius must be a non-empty CSS length");
    }
    const radiusProbe = document.createElement("span").style;
    radiusProbe.borderRadius = preset.radius;
    if (radiusProbe.borderRadius === "") {
      throw new Error("radius must be a valid CSS length");
    }
    const theme = preset.theme;
    if (!theme || typeof theme !== "object" || Array.isArray(theme)) {
      throw new Error("preset.theme must be an object");
    }
    return {
      style: preset.style,
      radius: preset.radius,
      mode: message.mode,
      values: exactTokenMap(theme[message.mode], `preset.theme.${message.mode}`),
    };
  }

  function applyState(state) {
    const staged = document.createElement("span").style;
    staged.setProperty("--radius", state.radius);
    for (const name of tokenNames) {
      staged.setProperty(`--${name}`, state.values[name]);
    }

    // One style-attribute replacement prevents a partially updated token set
    // from being observable between declarations.
    root.setAttribute("style", staged.cssText);
    root.classList.toggle("dark", state.mode === "dark");
    for (const [style, section] of styleSections) {
      section.hidden = style !== state.style;
    }
    if (focusedStyle !== state.style) {
      styleSections.get(state.style)
        ?.querySelector('[data-theme-preview-case="focus-visible"]')
        ?.focus();
      focusedStyle = state.style;
    }
  }

  addEventListener("message", (event) => {
    if (event.origin !== location.origin || event.source !== parent) {
      return;
    }
    try {
      applyState(validatedState(event.data));
      if (!hasAppliedState) {
        hasAppliedState = true;
        // Repeat ready after the first accepted state. This closes the race
        // where the initial ready message lands before the parent module has
        // installed its listener, without creating an acknowledgement loop.
        ready();
      }
    } catch (error) {
      fail(error instanceof Error ? error.message : "invalid preview state");
    }
  });

  ready();
}
