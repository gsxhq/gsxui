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
