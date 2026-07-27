import type { Page } from "@playwright/test";
import { expect, test } from "../support/fixtures";

async function recordCommandScrollRequests(page: Page) {
  await page.addInitScript(() => {
    (window as any).__commandScrollRequests = 0;
    const scrollIntoView = Element.prototype.scrollIntoView;
    Element.prototype.scrollIntoView = function (...args) {
      if ((this as Element).matches("[data-gsxui-command-item]")) {
        (window as any).__commandScrollRequests++;
      }
      return scrollIntoView.apply(this, args as [ScrollIntoViewOptions]);
    };
  });
}

test("initial selection does not request a document scroll", async ({ page }) => {
  await recordCommandScrollRequests(page);
  await page.goto("/x/command/basic");

  expect(
    await page.evaluate(() => (window as any).__commandScrollRequests),
  ).toBe(0);
});

test("keyboard selection keeps the selected item in view", async ({ page }) => {
  await recordCommandScrollRequests(page);
  await page.goto("/x/command/basic");
  await page.evaluate(() => ((window as any).__commandScrollRequests = 0));

  await page.locator("[data-gsxui-command-input]").first().press("ArrowDown");

  expect(
    await page.evaluate(() => (window as any).__commandScrollRequests),
  ).toBe(1);
});
