import type { Page } from "@playwright/test";
import { expect, test } from "../support/fixtures";

async function recordCommandScrollRequests(page: Page) {
  await page.addInitScript(() => {
    (window as any).__commandScrollRequests = 0;
    const scrollIntoView = Element.prototype.scrollIntoView;
    Element.prototype.scrollIntoView = function (...args) {
      if ((this as Element).matches("[data-gsxui-slot-command-item]")) {
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

  await page.locator("[data-gsxui-slot-command-input]").first().press("ArrowDown");

  expect(
    await page.evaluate(() => (window as any).__commandScrollRequests),
  ).toBe(1);
});

// Guards the filter() top-level reorder: a separator contains no items, so it
// scores 0, and an unconditional rank sort sank it below every rank-1 group to
// the very bottom of the list — stranding a between-groups divider at the end.
test("a separator stays between its groups when unfiltered", async ({ page }) => {
  await page.goto("/x/command/basic");

  const kinds = await page
    .locator("[data-gsxui-slot-command-list]")
    .first()
    .evaluate((list) =>
      [...list.children]
        .filter((child) => !(child as HTMLElement).hidden)
        .map((child) =>
          child.matches("[data-gsxui-slot-command-separator]")
            ? "separator"
            : child.matches("[data-gsxui-slot-command-group]")
              ? "group"
              : "other",
        ),
    );

  expect(kinds).toEqual(["group", "separator", "group"]);
});

test("the separator hides while a query filters and restores when it clears", async ({
  page,
}) => {
  await page.goto("/x/command/basic");
  const input = page.locator("[data-gsxui-slot-command-input]").first();
  const separator = page.locator("[data-gsxui-slot-command-separator]").first();

  await expect(separator).toBeVisible();
  await input.fill("cal");
  await expect(separator).toBeHidden();
  await input.fill("");
  await expect(separator).toBeVisible();
});
