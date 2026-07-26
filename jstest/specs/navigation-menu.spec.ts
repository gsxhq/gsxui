import { expect, test } from "../support/fixtures";

const BASIC = "/x/navigation-menu/basic";

async function mouseClick(
  page: import("@playwright/test").Page,
  locator: import("@playwright/test").Locator,
) {
  const box = await locator.boundingBox();
  if (!box) throw new Error("navigation-menu trigger has no bounding box");
  await page.mouse.click(box.x + box.width / 2, box.y + box.height / 2);
}

test("a mouse click opens and the second click closes", async ({ page }) => {
  await page.goto(BASIC);
  const trigger = page.locator("[data-gsxui-navigation-menu-trigger]");
  const content = page.locator("[data-gsxui-navigation-menu-content]");

  await mouseClick(page, trigger);
  await expect(trigger).toHaveAttribute("aria-expanded", "true");
  expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);

  await mouseClick(page, trigger);
  await expect(trigger).toHaveAttribute("aria-expanded", "false");
  expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(false);
});

test("keyboard Tab focus opens the trigger panel", async ({ page }) => {
  await page.goto(BASIC);
  const home = page.locator('[data-gsxui-navigation-menu-link]').first();
  const trigger = page.locator("[data-gsxui-navigation-menu-trigger]");
  const content = page.locator("[data-gsxui-navigation-menu-content]");

  await home.focus();
  await page.keyboard.press("Tab");

  await expect(trigger).toBeFocused();
  await expect(trigger).toHaveAttribute("aria-expanded", "true");
  expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
});
