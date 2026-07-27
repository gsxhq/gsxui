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
  await page.locator("[data-theme-download]").click();
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
