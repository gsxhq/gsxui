import { readFileSync } from "node:fs";

import { expect, test } from "../support/fixtures";
import { repoRoot } from "../support/paths";

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

test("theme editor downloads the exact variables-only theme.css", async ({
  page,
}) => {
  await page.goto("/theme");
  const downloadPromise = page.waitForEvent("download");
  await page.locator('[data-theme-download="css"]').click();
  const download = await downloadPromise;
  expect(download.suggestedFilename()).toBe("theme.css");
  const path = await download.path();
  if (!path) throw new Error("download has no local path");
  const exported = readFileSync(path, "utf8");
  const defaults = readFileSync(
    `${repoRoot}/assets/css/themes/default.css`,
    "utf8",
  );

  expect(variables(exported, ":root")).toEqual(variables(defaults, ":root"));
  expect(variables(exported, ".dark")).toEqual(variables(defaults, ".dark"));
  expect(exported).not.toMatch(/@import|@theme|@layer|tailwindcss|foundation|style\.css/);
});

test("style and valid theme state drive the isolated Button preview", async ({
  page,
}) => {
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

  const primary = page.locator('[data-theme-field="light.primary"]');
  await primary.fill("oklch(0.65 0.2 30)");
  await expect(primary).toHaveAttribute("aria-invalid", "false");
  await expect
    .poll(() =>
      preview.locator("html").evaluate((element) =>
        element.style.getPropertyValue("--primary").trim(),
      ),
    )
    .toBe("oklch(0.65 0.2 30)");

  await primary.fill("not-a-color");
  await expect(primary).toHaveAttribute("aria-invalid", "true");
  await expect(page.locator("[data-theme-status]")).toContainText("last valid values");
  expect(
    await preview.locator("html").evaluate((element) =>
      element.style.getPropertyValue("--primary").trim(),
    ),
  ).toBe("oklch(0.65 0.2 30)");

  await page.locator('[data-theme-mode-tab="dark"]').click();
  await expect(preview.locator("html")).toHaveClass(/\bdark\b/);
});

test("CSS import rejects duplicate recognized declarations atomically", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');
  const original = await primary.inputValue();

  await page
    .locator('[data-theme-import="css"]')
    .fill(":root { --primary: red; --primary: blue; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(page.locator("[data-theme-status]")).toContainText("duplicated");
  await expect(primary).toHaveValue(original);
});

test("CSS import rejects malformed unrelated syntax atomically", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');
  const original = await primary.inputValue();

  await page
    .locator('[data-theme-import="css"]')
    .fill("body { color red; } :root { --primary: green; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(page.locator("[data-theme-status]")).toContainText("malformed");
  await expect(primary).toHaveValue(original);
});

test("CSS import ignores valid unrelated CSS without mutating the preset", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');
  const original = await primary.inputValue();

  await page
    .locator('[data-theme-import="css"]')
    .fill("body { color: red; --unrelated: blue; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(page.locator("[data-theme-status]")).toHaveText("Nova · light");
  await expect(primary).toHaveValue(original);
});

test("CSS import accepts comments between supported selector tokens", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');

  await page
    .locator('[data-theme-import="css"]')
    .fill(":/* comment */root { --primary: red; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(primary).toHaveValue("red");
});

test("CSS import does not merge selector identifiers across comments", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');
  const original = await primary.inputValue();

  await page
    .locator('[data-theme-import="css"]')
    .fill(":r/**/oot { --primary: red; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(page.locator("[data-theme-status]")).toContainText(
    "must belong to :root or .dark",
  );
  await expect(primary).toHaveValue(original);
});

test("CSS import rejects important recognized declarations atomically", async ({
  page,
}) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');
  const original = await primary.inputValue();

  await page
    .locator('[data-theme-import="css"]')
    .fill(":root { --primary: red !important; }");
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(page.locator("[data-theme-status]")).toContainText(
    "theme.light.primary",
  );
  await expect(primary).toHaveValue(original);
});

test("CSS import ignores valid unrelated implicit nesting", async ({ page }) => {
  await page.goto("/theme");
  const primary = page.locator('[data-theme-field="light.primary"]');

  await page
    .locator('[data-theme-import="css"]')
    .fill(
      "body { .nested { color: red; --unrelated: blue; } } :root { --primary: green; }",
    );
  await page.locator('[data-theme-import-apply="css"]').click();

  await expect(primary).toHaveValue("green");
});

test("theme editor exposes Retry when the preview never handshakes", async ({
  page,
}) => {
  await page.route("**/theme/preview/button", async (route) => {
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
});
