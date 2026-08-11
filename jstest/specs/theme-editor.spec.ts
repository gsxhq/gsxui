import { readFileSync } from "node:fs";

import type { Page } from "@playwright/test";

import { expect, test } from "../support/fixtures";
import { repoRoot } from "../support/paths";

type ThemeSchema = {
  palette: {
    resolved: Record<
      string,
      Record<string, { light: Record<string, string>; dark: Record<string, string> }>
    >;
    chartColors: Record<string, { light: Record<string, string>; dark: Record<string, string> }>;
    fonts: { name: string; title: string; stack: string }[];
  };
};

function variables(css: string, selector: string) {
  const escaped = selector.replace(".", "\\.");
  const block = css.match(new RegExp(`${escaped}\\s*\\{([^}]+)\\}`, "s"))?.[1];
  if (!block) throw new Error(`missing ${selector} block`);
  return Object.fromEntries(
    [...block.matchAll(/(--[a-z0-9-]+)\s*:\s*([^;]+);/g)].map(
      ([, name, value]) => [name, value.trim()],
    ),
  );
}

async function downloadText(page: Page, kind: "json" | "css", pointer = true) {
  const downloadPromise = page.waitForEvent("download");
  const button = page.locator(`[data-theme-download="${kind}"]`);
  if (pointer) {
    await button.click();
  } else {
    await button.evaluate((element: HTMLButtonElement) => element.click());
  }
  const download = await downloadPromise;
  const path = await download.path();
  if (!path) throw new Error("download has no local path");
  return { filename: download.suggestedFilename(), text: readFileSync(path, "utf8") };
}

async function schema(page: Page): Promise<ThemeSchema> {
  return JSON.parse(await page.locator("[data-theme-schema]").textContent());
}

async function iframeVariable(page: Page, name: string) {
  return page
    .frameLocator("[data-theme-preview-frame]")
    .locator("html")
    .evaluate((element, property) => element.style.getPropertyValue(property).trim(), `--${name}`);
}

async function commands(page: Page) {
  return {
    init: await page.locator('[data-theme-command="init"]').inputValue(),
    apply: await page.locator('[data-theme-command="apply"]').inputValue(),
  };
}

function shareFromInit(command: string) {
  const match = command.match(/^gsxui init --preset '([^']+)'$/);
  if (!match) throw new Error(`unexpected init command ${JSON.stringify(command)}`);
  return match[1];
}

type PickerKind = "baseColor" | "theme" | "radius" | "chartColor" | "fontSans" | "fontHeading";

function picker(page: Page, kind: PickerKind) {
  return page.locator(`[data-theme-picker="${kind}"]`);
}

async function choose(page: Page, kind: PickerKind, accessibleName: string) {
  const control = picker(page, kind);
  const content = control.locator("[data-gsxui-slot-popover-content]");
  if (!(await content.evaluate((element) => element.matches(":popover-open")))) {
    await control.locator("[data-theme-picker-trigger]").click();
  }
  await control.getByRole("radio", { name: accessibleName, exact: true }).click();
}

async function selectionValue(page: Page, kind: PickerKind) {
  return picker(page, kind).locator("[data-theme-selection-value]").textContent();
}

test("theme editor downloads the exact variables-only theme.css", async ({ page }) => {
  await page.goto("/theme");
  const download = await downloadText(page, "css");
  expect(download.filename).toBe("theme.css");
  const defaults = readFileSync(`${repoRoot}/assets/css/themes/default.css`, "utf8");

  expect(variables(download.text, ":root")).toEqual(variables(defaults, ":root"));
  expect(variables(download.text, ".dark")).toEqual(variables(defaults, ".dark"));
  expect(download.text).not.toMatch(
    /@import|@theme|@layer|tailwindcss|foundation|style\.css/,
  );
});

test("style and mode remain independent from the curated palette", async ({ page }) => {
  await page.goto("/theme");
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");

  const preview = page.frameLocator("[data-theme-preview-frame]");
  const nova = preview.locator('[data-theme-preview-style="nova"]');
  const maia = preview.locator('[data-theme-preview-style="maia"]');
  await expect(nova).not.toHaveAttribute("hidden", "");
  await expect(maia).toHaveAttribute("hidden", "");

  const novaGeometry = await nova
    .getByRole("button", { name: "Default" })
    .first()
    .evaluate((element) => {
      const style = getComputedStyle(element);
      return { height: style.height, radius: style.borderRadius };
    });

  await page.locator('[data-theme-style="maia"]').click();
  await expect(nova).toHaveAttribute("hidden", "");
  await expect(maia).not.toHaveAttribute("hidden", "");
  const maiaGeometry = await maia
    .getByRole("button", { name: "Default" })
    .first()
    .evaluate((element) => {
      const style = getComputedStyle(element);
      return { height: style.height, radius: style.borderRadius };
    });
  expect(maiaGeometry.height).not.toBe(novaGeometry.height);
  expect(maiaGeometry.radius).not.toBe(novaGeometry.radius);

  await page.locator('[data-theme-mode-tab="dark"]').click();
  await expect(preview.locator("html")).toHaveClass(/\bdark\b/);
  await expect(selectionValue(page, "baseColor")).resolves.toBe("Neutral");
  await expect(selectionValue(page, "theme")).resolves.toBe("Neutral");
  await expect(selectionValue(page, "radius")).resolves.toBe("Medium");
});

test("pickers expose accessible catalog choices and no raw token inputs", async ({
  page,
}) => {
  await page.goto("/theme");
  await expect(page.locator("[data-theme-var]")).toHaveCount(0);
  await expect(page.locator('[data-theme-field^="light."]')).toHaveCount(0);
  await expect(page.locator('[data-theme-field^="dark."]')).toHaveCount(0);
  await expect(page.locator('[data-theme-field="radius"]')).toHaveCount(0);
  await expect(page.locator("iframe")).toHaveCount(1);

  const base = picker(page, "baseColor");
  await base.locator("[data-theme-picker-trigger]").click();
  await expect(base.locator("[data-gsxui-slot-popover-content]")).toHaveJSProperty(
    "popover",
    "auto",
  );
  expect(
    await base
      .locator("[data-gsxui-slot-popover-content]")
      .evaluate((element) => element.matches(":popover-open")),
  ).toBe(true);
  await expect(base.getByRole("radio")).toHaveCount(7);
  await expect(base.getByRole("radio", { name: "Neutral", exact: true })).toBeChecked();

  const theme = picker(page, "theme");
  await theme.locator("[data-theme-picker-trigger]").click();
  await expect(theme.getByRole("radio")).toHaveCount(18);
  await expect(theme.getByRole("radio", { name: "Neutral", exact: true })).toBeChecked();
  await expect(theme.getByRole("radio", { name: "Blue", exact: true })).toHaveCount(1);

  const radius = picker(page, "radius");
  await radius.locator("[data-theme-picker-trigger]").click();
  await expect(radius.getByRole("radio")).toHaveCount(4);
  await expect(radius.getByRole("radio", { name: "Medium", exact: true })).toBeChecked();

  const chartColor = picker(page, "chartColor");
  await chartColor.locator("[data-theme-picker-trigger]").click();
  await expect(chartColor.getByRole("radio")).toHaveCount(18);
  await expect(chartColor.getByRole("radio", { name: "Neutral", exact: true })).toBeChecked();

  for (const kind of ["baseColor", "theme", "radius", "chartColor"] as const) {
    await expect(picker(page, kind).locator("[data-theme-selection-swatch]")).toBeVisible();
    await expect(picker(page, kind).locator("[data-theme-choice-swatch]")).not.toHaveCount(0);
  }

  // Font pickers are the same catalog-radio shape as every axis above, but
  // deliberately carry no swatch — a font isn't a hue or a radius, so
  // ThemePicker's swatch circle stays hidden for every option (catalog.go's
  // FontChoice has no Swatch field, unlike PaletteChoice).
  const fontSans = picker(page, "fontSans");
  await fontSans.locator("[data-theme-picker-trigger]").click();
  await expect(fontSans.getByRole("radio")).toHaveCount(6);
  await expect(fontSans.getByRole("radio", { name: "Geist", exact: true })).toBeChecked();

  const fontHeading = picker(page, "fontHeading");
  await fontHeading.locator("[data-theme-picker-trigger]").click();
  await expect(fontHeading.getByRole("radio")).toHaveCount(6);
  await expect(fontHeading.getByRole("radio", { name: "Geist", exact: true })).toBeChecked();

  for (const kind of ["fontSans", "fontHeading"] as const) {
    await expect(picker(page, kind).locator("[data-theme-selection-swatch]")).toBeHidden();
  }
});

test("keyboard reopening focuses the checked theme-picker radio", async ({ page }) => {
  await page.goto("/theme");
  const theme = picker(page, "theme");
  const trigger = theme.locator("[data-theme-picker-trigger]");
  const content = theme.locator("[data-gsxui-slot-popover-content]");
  const blue = theme.getByRole("radio", { name: "Blue", exact: true });

  await trigger.click();
  await blue.click();
  await expect(blue).toBeChecked();
  await page.keyboard.press("Escape");
  await expect(trigger).toBeFocused();

  await trigger.press("Enter");
  await expect(content).toHaveAttribute("data-state", "open");
  await expect(blue).toBeFocused();
});

test("unrelated gsxui radios stay outside the theme-picker controller", async ({ page }) => {
  await page.goto("/theme");
  const committedJSON = (await downloadText(page, "json")).text;

  const errorMessage = await page.evaluate(async () => {
    const radio = document.createElement("input");
    radio.type = "radio";
    radio.setAttribute("data-gsxui-slot-radio", "");
    let message = "";
    const recordError = (event: ErrorEvent) => {
      event.preventDefault();
      message = event.message;
    };
    addEventListener("error", recordError);
    document.body.append(radio);
    radio.dispatchEvent(new Event("change", { bubbles: true }));
    await new Promise((resolve) => setTimeout(resolve, 0));
    removeEventListener("error", recordError);
    radio.remove();
    return message;
  });

  expect(errorMessage).toBe("");
  expect((await downloadText(page, "json", false)).text).toBe(committedJSON);
});

test("Base Color and Theme choices update both iframe modes and the same-named option", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);

  await choose(page, "baseColor", "Stone");
  await expect(selectionValue(page, "baseColor")).resolves.toBe("Stone");
  await expect(selectionValue(page, "theme")).resolves.toBe("Stone");
  await expect
    .poll(() => iframeVariable(page, "foreground"))
    .toBe(catalog.palette.resolved.stone.stone.light.foreground);

  const theme = picker(page, "theme");
  await theme.locator("[data-theme-picker-trigger]").click();
  await expect(theme.getByRole("radio", { name: "Stone", exact: true })).toBeChecked();
  await expect(theme.getByRole("radio", { name: "Neutral", exact: true })).toHaveCount(0);
  await expect(theme.getByRole("radio")).toHaveCount(18);

  await theme.getByRole("radio", { name: "Blue", exact: true }).click();
  await expect(selectionValue(page, "theme")).resolves.toBe("Blue");
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.stone.blue.light.primary);

  await page.locator('[data-theme-mode-tab="dark"]').click();
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.stone.blue.dark.primary);
});

test("menu accent tabs overlay accent from primary and survive a theme change", async ({
  page,
}) => {
  await page.goto("/theme");
  const subtleTab = page.locator('[data-theme-menu-accent-tab="subtle"]');
  const boldTab = page.locator('[data-theme-menu-accent-tab="bold"]');
  await expect(subtleTab).toHaveAttribute("aria-pressed", "true");
  await expect(boldTab).toHaveAttribute("aria-pressed", "false");

  const subtleAccent = await iframeVariable(page, "accent");
  const primary = await iframeVariable(page, "primary");
  expect(subtleAccent).not.toBe(primary);

  await boldTab.click();
  await expect(boldTab).toHaveAttribute("aria-pressed", "true");
  await expect(subtleTab).toHaveAttribute("aria-pressed", "false");
  await expect.poll(() => iframeVariable(page, "accent")).toBe(primary);

  // The override survives an unrelated theme change instead of being wiped
  // by the fresh base color/theme resolution.
  await choose(page, "theme", "Blue");
  await expect.poll(() => iframeVariable(page, "primary")).not.toBe(primary);
  await expect.poll(() => iframeVariable(page, "accent")).toBe(await iframeVariable(page, "primary"));

  await subtleTab.click();
  await expect.poll(() => iframeVariable(page, "accent")).not.toBe(await iframeVariable(page, "primary"));
});

test("chart color picker overlays chart-1..5 independent of theme and survives a theme change", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);

  const beforeChart1 = await iframeVariable(page, "chart-1");
  const beforePrimary = await iframeVariable(page, "primary");
  await expect(selectionValue(page, "chartColor")).resolves.toBe("Neutral");

  await choose(page, "chartColor", "Violet");
  await expect(selectionValue(page, "chartColor")).resolves.toBe("Violet");
  // Base color/theme picker selections are untouched by the chart color
  // axis — it is independent, not a collapse to a matching theme.
  await expect(selectionValue(page, "baseColor")).resolves.toBe("Neutral");
  await expect(selectionValue(page, "theme")).resolves.toBe("Neutral");
  await expect.poll(() => iframeVariable(page, "primary")).toBe(beforePrimary);
  await expect
    .poll(() => iframeVariable(page, "chart-1"))
    .toBe(catalog.palette.chartColors.violet.light["chart-1"]);
  expect(await iframeVariable(page, "chart-1")).not.toBe(beforeChart1);

  // The override survives an unrelated theme change instead of being wiped
  // by the fresh base color/theme resolution.
  await choose(page, "theme", "Blue");
  await expect.poll(() => iframeVariable(page, "primary")).not.toBe(beforePrimary);
  await expect
    .poll(() => iframeVariable(page, "chart-1"))
    .toBe(catalog.palette.chartColors.violet.light["chart-1"]);
});

test("font pickers set --font-sans/--font-heading independently and survive a theme change", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);
  const interStack = catalog.palette.fonts.find((choice) => choice.name === "inter")!.stack;
  const playfairStack = catalog.palette.fonts.find(
    (choice) => choice.name === "playfair-display",
  )!.stack;

  const beforeFontSans = await iframeVariable(page, "font-sans");
  const beforeFontHeading = await iframeVariable(page, "font-heading");
  const beforePrimary = await iframeVariable(page, "primary");
  await expect(selectionValue(page, "fontSans")).resolves.toBe("Geist");
  await expect(selectionValue(page, "fontHeading")).resolves.toBe("Geist");
  expect(beforeFontSans).toBe(beforeFontHeading);

  await choose(page, "fontSans", "Inter");
  await expect(selectionValue(page, "fontSans")).resolves.toBe("Inter");
  // Body and heading fonts are independent axes: changing one never touches
  // the other, and neither touches unrelated theme colors.
  await expect(selectionValue(page, "fontHeading")).resolves.toBe("Geist");
  await expect.poll(() => iframeVariable(page, "font-sans")).toBe(interStack);
  await expect.poll(() => iframeVariable(page, "font-heading")).toBe(beforeFontHeading);
  await expect.poll(() => iframeVariable(page, "primary")).toBe(beforePrimary);

  await choose(page, "fontHeading", "Playfair Display");
  await expect(selectionValue(page, "fontHeading")).resolves.toBe("Playfair Display");
  await expect.poll(() => iframeVariable(page, "font-heading")).toBe(playfairStack);
  // Body font untouched by the heading font change.
  await expect.poll(() => iframeVariable(page, "font-sans")).toBe(interStack);

  // Both overrides survive an unrelated theme change.
  await choose(page, "theme", "Blue");
  await expect.poll(() => iframeVariable(page, "primary")).not.toBe(beforePrimary);
  await expect.poll(() => iframeVariable(page, "font-sans")).toBe(interStack);
  await expect.poll(() => iframeVariable(page, "font-heading")).toBe(playfairStack);

  // The picked heading font is visibly applied in the preview, not just
  // set as an inert CSS variable — gallery.gsx.src's CardTitle elements
  // read var(--font-heading) directly.
  const cardTitleFontFamily = await page
    .frameLocator("[data-theme-preview-frame]")
    .locator('[data-theme-preview-style="nova"] [data-gsxui-slot-card-title]')
    .first()
    .evaluate((element) => getComputedStyle(element).fontFamily);
  expect(cardTitleFontFamily).toContain("Playfair Display Variable");
});

test("desktop hover previews only the iframe and restores on dismissal or pointer leave", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);
  const committedPrimary = await iframeVariable(page, "primary");
  const committedCommands = await commands(page);
  const committedJSON = (await downloadText(page, "json")).text;

  const theme = picker(page, "theme");
  await theme.locator("[data-theme-picker-trigger]").click();
  const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");
  await rose.locator("[data-theme-choice-swatch]").hover();

  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  const roseBox = await rose.boundingBox();
  if (!roseBox) throw new Error("Rose choice has no layout box");
  await page.mouse.move(
    roseBox.x + roseBox.width - 2,
    roseBox.y + roseBox.height / 2,
  );
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  await expect(
    theme.getByRole("radio", { name: "Neutral", exact: true }),
  ).toBeFocused();
  expect(await commands(page)).toEqual(committedCommands);
  expect((await downloadText(page, "json", false)).text).toBe(committedJSON);

  await page.keyboard.press("Escape");
  await expect
    .poll(() =>
      theme
        .locator("[data-gsxui-slot-popover-content]")
        .evaluate((element) => element.matches(":popover-open")),
    )
    .toBe(false);
  await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);

  await theme.locator("[data-theme-picker-trigger]").click();
  await rose.hover();
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  await page.mouse.move(1, 1);
  await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  expect(await commands(page)).toEqual(committedCommands);
});

test("touch pointerenter does not preview on a fine hover-capable device", async ({ page }) => {
  await page.goto("/theme");
  const committedPrimary = await iframeVariable(page, "primary");
  const theme = picker(page, "theme");
  await theme.locator("[data-theme-picker-trigger]").click();
  const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");

  expect(
    await page.evaluate(() =>
      matchMedia("(hover: hover) and (pointer: fine)").matches,
    ),
  ).toBe(true);
  const pointerType = await rose.evaluate((element) => {
    const event = new PointerEvent("pointerenter", {
      bubbles: false,
      composed: true,
      isPrimary: true,
      pointerId: 17,
      pointerType: "touch",
    });
    element.dispatchEvent(event);
    return event.pointerType;
  });
  expect(pointerType).toBe("touch");

  await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
});

test("radius hover previews only the iframe and restores on pointer leave", async ({ page }) => {
  await page.goto("/theme");
  const committedRadius = await iframeVariable(page, "radius");
  const committedCommands = await commands(page);
  const committedJSON = (await downloadText(page, "json")).text;
  const radius = picker(page, "radius");
  await radius.locator("[data-theme-picker-trigger]").click();

  await radius.getByRole("radio", { name: "Large", exact: true }).locator("..").hover();
  await expect.poll(() => iframeVariable(page, "radius")).toBe("0.875rem");
  expect(await commands(page)).toEqual(committedCommands);
  expect((await downloadText(page, "json", false)).text).toBe(committedJSON);

  await page.mouse.move(1, 1);
  await expect.poll(() => iframeVariable(page, "radius")).toBe(committedRadius);
  expect(await commands(page)).toEqual(committedCommands);
});

test("click commits palette state to JSON, share code, commands, and iframe", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);
  const beforeJSON = (await downloadText(page, "json")).text;
  const beforeCommands = await commands(page);
  const beforeShare = shareFromInit(beforeCommands.init);

  await choose(page, "theme", "Blue");

  const afterJSON = (await downloadText(page, "json")).text;
  const afterCommands = await commands(page);
  expect(afterJSON).not.toBe(beforeJSON);
  expect(JSON.parse(afterJSON).theme.light.primary).toBe(
    catalog.palette.resolved.neutral.blue.light.primary,
  );
  expect(shareFromInit(afterCommands.init)).not.toBe(beforeShare);
  expect(afterCommands.init).not.toBe(beforeCommands.init);
  expect(afterCommands.apply).not.toBe(beforeCommands.apply);
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.neutral.blue.light.primary);
});

test("committed changes sync the address bar without growing history", async ({ page }) => {
  await page.goto("/theme");
  const beforeSearch = await page.evaluate(() => location.search);
  expect(beforeSearch).toMatch(/^\?preset=gsxui%3Ap1%3A/);

  await choose(page, "theme", "Blue");
  const afterBlueSearch = await page.evaluate(() => location.search);
  expect(afterBlueSearch).not.toBe(beforeSearch);
  expect(afterBlueSearch).toMatch(/^\?preset=gsxui%3Ap1%3A/);
  expect(shareFromInit((await commands(page)).init)).toBe(
    decodeURIComponent(afterBlueSearch.slice("?preset=".length)),
  );

  await choose(page, "radius", "Large");
  const afterRadiusSearch = await page.evaluate(() => location.search);
  expect(afterRadiusSearch).not.toBe(afterBlueSearch);

  // Both commits used replaceState, not pushState: a single back navigation
  // leaves /theme entirely rather than stepping through each change.
  await page.goBack();
  await expect(page).not.toHaveURL(/\/theme/);
});

test("the address bar's share code round-trips to the exact committed preset", async ({
  page,
}) => {
  await page.goto("/theme");
  await choose(page, "theme", "Blue");
  const search = await page.evaluate(() => location.search);
  const committedJSON = (await downloadText(page, "json")).text;

  await page.goto(`/theme${search}`);
  expect((await downloadText(page, "json")).text).toBe(committedJSON);
});

test("undo/redo buttons and keyboard shortcuts replay committed changes", async ({ page }) => {
  await page.goto("/theme");
  const undo = page.locator("[data-theme-undo]");
  const redo = page.locator("[data-theme-redo]");
  await expect(undo).toBeDisabled();
  await expect(redo).toBeDisabled();

  await choose(page, "theme", "Blue");
  await expect(selectionValue(page, "theme")).resolves.toBe("Blue");
  await expect(undo).toBeEnabled();
  await expect(redo).toBeDisabled();

  await choose(page, "theme", "Rose");
  await expect(selectionValue(page, "theme")).resolves.toBe("Rose");

  await undo.click();
  await expect(selectionValue(page, "theme")).resolves.toBe("Blue");
  await expect(redo).toBeEnabled();

  // Keyboard shortcuts reach the editor even though focus is not on the
  // undo/redo buttons themselves.
  await page.locator("[data-theme-style-panel] h2").click();
  await page.keyboard.press("ControlOrMeta+z");
  await expect(selectionValue(page, "theme")).resolves.toBe("Neutral");
  await expect(undo).toBeDisabled();

  await page.keyboard.press("ControlOrMeta+Shift+z");
  await expect(selectionValue(page, "theme")).resolves.toBe("Blue");

  await page.keyboard.press("ControlOrMeta+Shift+z");
  await expect(selectionValue(page, "theme")).resolves.toBe("Rose");
  await expect(redo).toBeDisabled();

  // A new commit after undoing discards the redo branch.
  await undo.click();
  await choose(page, "radius", "Large");
  await expect(redo).toBeDisabled();

  // The keyboard shortcut is inert while typing in an import textarea, so
  // native text-undo keeps working there instead.
  const jsonImport = page.locator('[data-theme-import="json"]');
  await jsonImport.fill("hello");
  await jsonImport.press("ControlOrMeta+z");
  await expect(selectionValue(page, "radius")).resolves.toBe("Large");
});

test("share commands use compact built-ins and full custom imports", async ({ page }) => {
  await page.goto("/theme");

  const initialShare = shareFromInit((await commands(page)).init);
  expect(initialShare).toMatch(/^gsxui:p1:/);
  expect(initialShare.length).toBeLessThanOrEqual(12);

  await choose(page, "theme", "Blue");
  const selectedShare = shareFromInit((await commands(page)).init);
  expect(selectedShare).toMatch(/^gsxui:p1:/);
  expect(selectedShare.length).toBeLessThanOrEqual(12);
  expect(selectedShare).not.toBe(initialShare);

  const imported = JSON.parse((await downloadText(page, "json")).text);
  imported.theme.light.primary = "rgb(1 2 3)";
  await page.locator('[data-theme-import="json"]').fill(`${JSON.stringify(imported)}\n`);
  await page.locator('[data-theme-import-apply="json"]').click();

  const customShare = shareFromInit((await commands(page)).init);
  expect(customShare).toMatch(/^gsxui:v1:/);
  expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(imported);
});

test("historical full built-in share URLs still load", async ({ page }) => {
  const source = readFileSync(
    `${repoRoot}/internal/preset/testdata/default-nova.json`,
    "utf8",
  );
  const fullCode = `gsxui:v1:${Buffer.from(source).toString("base64url")}`;
  await page.goto(`/theme?preset=${encodeURIComponent(fullCode)}`);

  await expect(page.locator("[data-theme-status]")).toHaveText("Loaded shared preset.");
  expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(JSON.parse(source));
  expect(shareFromInit((await commands(page)).init)).toMatch(/^gsxui:p1:/);
});

test("custom JSON import is lossless and a built-in base replaces it atomically", async ({
  page,
}) => {
  await page.goto("/theme");
  const catalog = await schema(page);
  const imported = JSON.parse((await downloadText(page, "json")).text);
  imported.radius = "1rem";
  imported.theme.light.primary = "rgb(1 2 3)";
  imported.theme.dark.primary = "rgb(4 5 6)";

  await page.locator('[data-theme-import="json"]').fill(`${JSON.stringify(imported)}\n`);
  await page.locator('[data-theme-import-apply="json"]').click();

  await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  await expect.poll(() => iframeVariable(page, "primary")).toBe("rgb(1 2 3)");
  expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(imported);

  await choose(page, "baseColor", "Stone");
  await expect(selectionValue(page, "baseColor")).resolves.toBe("Stone");
  await expect(selectionValue(page, "theme")).resolves.toBe("Stone");
  await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  await expect
    .poll(() => iframeVariable(page, "primary"))
    .toBe(catalog.palette.resolved.stone.stone.light.primary);
  const replaced = JSON.parse((await downloadText(page, "json")).text);
  expect(replaced.theme.light).toEqual(catalog.palette.resolved.stone.stone.light);
  expect(replaced.theme.dark).toEqual(catalog.palette.resolved.stone.stone.dark);
  expect(replaced.radius).toBe("1rem");
});

test("custom CSS import is lossless while a named radius preserves its colors", async ({
  page,
}) => {
  await page.goto("/theme");
  await page
    .locator('[data-theme-import="css"]')
    .fill(
      ":root { --radius: 1rem; --primary: rgb(10 20 30); }\n" +
        ".dark { --primary: rgb(40 50 60); }",
    );
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  await expect.poll(() => iframeVariable(page, "primary")).toBe("rgb(10 20 30)");

  await choose(page, "radius", "Large");
  await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  await expect(selectionValue(page, "radius")).resolves.toBe("Large");
  const exported = JSON.parse((await downloadText(page, "json")).text);
  expect(exported.theme.light.primary).toBe("rgb(10 20 30)");
  expect(exported.theme.dark.primary).toBe("rgb(40 50 60)");
  expect(exported.radius).toBe("0.875rem");
});

for (const rejection of [
  {
    name: "duplicate recognized declarations",
    css: ":root { --primary: red; --primary: blue; }",
    message: "duplicated",
  },
  {
    name: "malformed unrelated syntax",
    css: "body { color red; } :root { --primary: green; }",
    message: "malformed",
  },
  {
    name: "selector identifiers split across comments",
    css: ":r/**/oot { --primary: red; }",
    message: "must belong to :root or .dark",
  },
  {
    name: "important recognized declarations",
    css: ":root { --primary: red !important; }",
    message: "theme.light.primary",
  },
]) {
  test(`CSS import rejects ${rejection.name} atomically`, async ({ page }) => {
    await page.goto("/theme");
    const before = await commands(page);
    await page.locator('[data-theme-import="css"]').fill(rejection.css);
    await page.locator('[data-theme-import-apply="css"]').click();

    await expect(page.locator("[data-theme-status]")).toContainText(rejection.message);
    expect(await commands(page)).toEqual(before);
  });
}

test("CSS import ignores valid unrelated CSS without mutating the preset", async ({
  page,
}) => {
  await page.goto("/theme");
  const before = await commands(page);
  await page
    .locator('[data-theme-import="css"]')
    .fill("body { color: red; --unrelated: blue; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  expect(await commands(page)).toEqual(before);
});

test("CSS import accepts comments between supported selector tokens", async ({
  page,
}) => {
  await page.goto("/theme");
  await page
    .locator('[data-theme-import="css"]')
    .fill(":/* comment */root { --primary: red; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  await expect.poll(() => iframeVariable(page, "primary")).toBe("red");
});

test("CSS import ignores valid unrelated implicit nesting", async ({ page }) => {
  await page.goto("/theme");
  await page
    .locator('[data-theme-import="css"]')
    .fill(
      "body { .nested { color: red; --unrelated: blue; } } :root { --primary: green; }",
    );
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  await expect.poll(() => iframeVariable(page, "primary")).toBe("green");
});

// The preview body scope-disables web/site-button.css, so raw slot-button
// markup inside the gallery (Calendar day/nav buttons) gets its Button chrome
// from the generated web/theme-preview-button.css instead — the NOVA section
// from Nova's Button recipe, the MAIA section from Maia's. All of it is
// measured against the /x/calendar/basic fixture, which renders under the
// normal site-button fallback, never assumed.
test("preview calendar day buttons keep per-style Button chrome in both sections", async ({
  page,
}) => {
  const chrome = (element: Element) => {
    const style = getComputedStyle(element);
    return {
      radius: style.borderRadius,
      alignItems: style.alignItems,
      justifyContent: style.justifyContent,
    };
  };

  await page.goto("/x/calendar/basic");
  const fixtureChrome = await page
    .locator('[data-gsxui-slot-calendar-day-button][data-date="2026-01-15"]')
    .evaluate(chrome);
  expect(fixtureChrome.alignItems).toBe("center");
  expect(fixtureChrome.justifyContent).toBe("center");
  expect(fixtureChrome.radius).not.toBe("0px");

  await page.goto("/theme");
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  const preview = page.frameLocator("[data-theme-preview-frame]");
  const nova = preview.locator('[data-theme-preview-style="nova"]');
  const maia = preview.locator('[data-theme-preview-style="maia"]');

  const novaChrome = await nova
    .locator('[data-gsxui-slot-calendar-day-button][data-date="2026-01-15"]')
    .first()
    .evaluate(chrome);
  expect(novaChrome).toEqual(fixtureChrome);

  await page.locator('[data-theme-style="maia"]').click();
  await expect(maia).not.toHaveAttribute("hidden", "");
  const maiaButtonRadius = (
    await maia
      .getByRole("button", { name: "Default" })
      .first()
      .evaluate(chrome)
  ).radius;
  const maiaChrome = await maia
    .locator('[data-gsxui-slot-calendar-day-button][data-date="2026-01-15"]')
    .first()
    .evaluate(chrome);
  expect(maiaChrome.alignItems).toBe("center");
  expect(maiaChrome.justifyContent).toBe("center");
  expect(maiaChrome.radius).toBe(maiaButtonRadius);
  expect(maiaChrome.radius).not.toBe(novaChrome.radius);
});

// The preview entry loads the real behaviour barrel, so the gallery's
// interactive components must actually operate inside the iframe.
test("preview gallery menus and dialogs operate inside the iframe", async ({
  page,
}) => {
  await page.goto("/theme");
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  const nova = page
    .frameLocator("[data-theme-preview-frame]")
    .locator('[data-theme-preview-style="nova"]');

  const dropdownTrigger = nova
    .locator("[data-gsxui-slot-dropdown-menu-trigger]")
    .first();
  await dropdownTrigger.scrollIntoViewIfNeeded();
  await dropdownTrigger.click();
  const dropdownContent = nova
    .locator("[data-gsxui-slot-dropdown-menu-content]")
    .first();
  await expect
    .poll(() =>
      dropdownContent.evaluate((element) => element.matches(":popover-open")),
    )
    .toBe(true);
  await expect(dropdownContent.getByText("My Account")).toBeVisible();
  await page.keyboard.press("Escape");
  await expect
    .poll(() =>
      dropdownContent.evaluate((element) => element.matches(":popover-open")),
    )
    .toBe(false);

  const dialogTrigger = nova.locator("[data-gsxui-slot-dialog-trigger]").first();
  await dialogTrigger.scrollIntoViewIfNeeded();
  await dialogTrigger.click();
  const dialog = nova.locator("dialog[data-gsxui-slot-dialog-content]").first();
  await expect.poll(() => dialog.evaluate((element) => (element as HTMLDialogElement).open)).toBe(
    true,
  );
  await expect(dialog.getByText("Edit profile").first()).toBeVisible();
  await page.keyboard.press("Escape");
  await expect
    .poll(() => dialog.evaluate((element) => (element as HTMLDialogElement).open))
    .toBe(false);
});

test("preview gallery toast keeps its close button anchored to the row", async ({
  page,
}) => {
  // The gallery renders a Toast row in-flow with a position override
  // (foundation absolutely positions rows for the live Toaster stack). The
  // override must stay `relative`, not `static`: the row is the containing
  // block for the absolute close button, and a static row let the button
  // escape to the page's top-right corner.
  await page.goto("/theme/preview");
  const toast = page
    .locator('[data-theme-preview-style="nova"] [data-gsxui-slot-toast]')
    .first();
  const geometry = await toast.evaluate((row) => {
    const button = row.querySelector("[data-gsxui-slot-toast-close]")!;
    const rowRect = row.getBoundingClientRect();
    const buttonRect = button.getBoundingClientRect();
    return {
      position: getComputedStyle(row).position,
      buttonInsideRow:
        buttonRect.left > rowRect.left &&
        buttonRect.left < rowRect.right + 8 &&
        buttonRect.top > rowRect.top - 8 &&
        buttonRect.top < rowRect.bottom,
    };
  });
  expect(geometry.position).toBe("relative");
  expect(geometry.buttonInsideRow).toBe(true);
});

test("theme editor exposes Retry when the preview never handshakes", async ({
  page,
}) => {
  let responsive = false;
  await page.route("**/theme/preview", async (route) => {
    if (responsive) {
      await route.fallback();
      return;
    }
    await route.fulfill({
      contentType: "text/html",
      body: "<!doctype html><title>Unresponsive preview</title>",
    });
  });
  await page.goto("/theme");

  await expect(page.locator("[data-theme-preview-status]")).toHaveText(
    "Preview did not respond.",
  );
  await expect(page.locator("[data-theme-preview-retry]")).toBeVisible();

  responsive = true;
  await page.locator("[data-theme-preview-retry]").click();
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  await expect(page.locator("[data-theme-preview-retry]")).toBeHidden();
});

test("preview acknowledgement survives a late parent-observed iframe load", async ({
  page,
}) => {
  await page.goto("/theme");
  const status = page.locator("[data-theme-preview-status]");
  const retry = page.locator("[data-theme-preview-retry]");
  await expect(status).toHaveText("Live");
  await expect(retry).toBeHidden();

  await page
    .locator("[data-theme-preview-frame]")
    .evaluate((frame) => frame.dispatchEvent(new Event("load")));

  await page.waitForTimeout(2_100);
  await expect(status).toHaveText("Live");
  await expect(retry).toBeHidden();
});

test("stale previous-document responses cannot complete a fresh preview attempt", async ({
  page,
}) => {
  let responsive = true;
  await page.route("**/theme/preview", async (route) => {
    if (responsive) {
      await route.fallback();
      return;
    }
    await route.fulfill({
      contentType: "text/html",
      body: "<!doctype html><body data-unresponsive-preview>Unresponsive preview</body>",
    });
  });

  await page.goto("/theme");
  const frame = page.locator("[data-theme-preview-frame]");
  const preview = page.frameLocator("[data-theme-preview-frame]");
  const status = page.locator("[data-theme-preview-status]");
  const retry = page.locator("[data-theme-preview-retry]");
  await expect(status).toHaveText("Live");

  await preview.locator("html").evaluate(() => {
    addEventListener("message", (event) => {
      if (event.data?.type === "gsxui:theme-preview:v1") {
        (
          globalThis as typeof globalThis & {
            capturedThemePreviewAttempt?: string;
          }
        ).capturedThemePreviewAttempt = event.data.attempt;
      }
    });
  });
  await page.locator('[data-theme-mode-tab="dark"]').click();
  await expect
    .poll(() =>
      preview.locator("html").evaluate(
        () =>
          (
            globalThis as typeof globalThis & {
              capturedThemePreviewAttempt?: string;
            }
          ).capturedThemePreviewAttempt,
      ),
    )
    .toBeTruthy();
  const staleAttempt = await preview.locator("html").evaluate(
    () =>
      (
        globalThis as typeof globalThis & {
          capturedThemePreviewAttempt: string;
        }
      ).capturedThemePreviewAttempt,
  );
  await expect(status).toHaveText("Live");
  await expect(retry).toBeHidden();

  responsive = false;
  await frame.evaluate((element: HTMLIFrameElement) => {
    element.src = element.src;
  });
  await expect(preview.locator("[data-unresponsive-preview]")).toBeVisible();

  await preview.locator("html").evaluate((_, attempt) => {
    parent.postMessage(
      { type: "gsxui:theme-preview-applied:v1", attempt },
      location.origin,
    );
    parent.postMessage(
      {
        type: "gsxui:theme-preview-error:v1",
        attempt,
        message: "stale preview error",
      },
      location.origin,
    );
  }, staleAttempt);

  await expect(status).toHaveText("Preview did not respond.");
  await expect(retry).toBeVisible();

  responsive = true;
  await retry.click();
  await expect(status).toHaveText("Live");
  await expect(retry).toBeHidden();
});

test("preview rejects state messages with an empty attempt identity", async ({
  page,
}) => {
  await page.goto("/theme");
  const frame = page.locator("[data-theme-preview-frame]");
  const preview = page.frameLocator("[data-theme-preview-frame]");
  const status = page.locator("[data-theme-preview-status]");
  const retry = page.locator("[data-theme-preview-retry]");
  await expect(status).toHaveText("Live");

  await page.evaluate(() => {
    addEventListener("message", (event) => {
      if (event.data?.attempt === "") {
        (
          globalThis as typeof globalThis & {
            invalidAttemptResponse?: unknown;
          }
        ).invalidAttemptResponse = event.data;
      }
    });
  });
  await preview.locator("html").evaluate(() => {
    addEventListener("message", (event) => {
      if (event.data?.type === "gsxui:theme-preview:v1") {
        (
          globalThis as typeof globalThis & {
            capturedThemePreviewMessage?: unknown;
          }
        ).capturedThemePreviewMessage = structuredClone(event.data);
      }
    });
  });
  await page.locator('[data-theme-mode-tab="dark"]').click();
  await expect
    .poll(() =>
      preview.locator("html").evaluate(
        () =>
          (
            globalThis as typeof globalThis & {
              capturedThemePreviewMessage?: unknown;
            }
          ).capturedThemePreviewMessage,
      ),
    )
    .not.toBeUndefined();
  const payload = await preview.locator("html").evaluate(
    () =>
      (
        globalThis as typeof globalThis & {
          capturedThemePreviewMessage: unknown;
        }
      ).capturedThemePreviewMessage,
  );

  await frame.evaluate(
    (
      element: HTMLIFrameElement,
      message: Record<string, unknown>,
    ) => {
      element.contentWindow?.postMessage(
        { ...message, attempt: "", mode: "light" },
        location.origin,
      );
    },
    payload,
  );

  await expect
    .poll(() =>
      page.evaluate(
        () =>
          (
            globalThis as typeof globalThis & {
              invalidAttemptResponse?: {
                type?: string;
                message?: string;
              };
            }
          ).invalidAttemptResponse,
      ),
    )
    .toMatchObject({
      type: "gsxui:theme-preview-error:v1",
      message: "attempt must be a non-empty string",
    });
  await expect
    .poll(() =>
      preview
        .locator("html")
        .evaluate((element) => element.classList.contains("dark")),
    )
    .toBe(true);
  await expect(status).toHaveText("Live");
  await expect(retry).toBeHidden();
});

test("the base-color picker's popover stays inside the viewport", async ({ page }) => {
  // At the desktop layout the picker column sits at the left edge with a
  // ~156px trigger; the popover content is 288px wide and centered on the
  // trigger's midpoint, which computes left = -34px — rendered flush
  // offscreen-left before clampToViewport existed (reported from a real
  // session). Asserting x >= 0 discriminates: mid-enter the scale animation
  // shifts the box inward by ~7px, so the unclamped position still measures
  // negative.
  await page.goto("/theme");
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  const trigger = page.locator('[data-theme-picker="baseColor"] [data-gsxui-slot-popover-trigger]');
  await trigger.click();
  const content = page.locator('[data-theme-picker="baseColor"] [data-gsxui-slot-popover-content]');
  await expect(content).toBeVisible();
  const box = await content.boundingBox();
  expect(box, "content has a box").not.toBeNull();
  expect(box!.x, "left edge inside viewport").toBeGreaterThanOrEqual(0);
});
