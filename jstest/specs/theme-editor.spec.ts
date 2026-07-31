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

function picker(page: Page, kind: "baseColor" | "theme" | "radius") {
  return page.locator(`[data-theme-picker="${kind}"]`);
}

async function choose(
  page: Page,
  kind: "baseColor" | "theme" | "radius",
  accessibleName: string,
) {
  const control = picker(page, kind);
  const content = control.locator("[data-gsxui-slot-popover-content]");
  if (!(await content.evaluate((element) => element.matches(":popover-open")))) {
    await control.locator("[data-theme-picker-trigger]").click();
  }
  await control.getByRole("radio", { name: accessibleName, exact: true }).click();
}

async function selectionValue(
  page: Page,
  kind: "baseColor" | "theme" | "radius",
) {
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

  for (const kind of ["baseColor", "theme", "radius"] as const) {
    await expect(picker(page, kind).locator("[data-theme-selection-swatch]")).toBeVisible();
    await expect(picker(page, kind).locator("[data-theme-choice-swatch]")).not.toHaveCount(0);
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
