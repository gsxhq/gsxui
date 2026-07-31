import { expect, test } from "../support/fixtures";

const CHECKBOXES = "/x/dropdown-menu/checkboxes";

test("one checkbox click changes state and emits exactly once", async ({ page }) => {
  await page.goto(CHECKBOXES);
  await page.evaluate(() => {
    (window as any).__dropdownChanges = [];
    document.addEventListener("gsxui:change", (event) => {
      (window as any).__dropdownChanges.push((event as CustomEvent).detail);
    });
  });

  const trigger = page.locator("[data-gsxui-slot-dropdown-menu-trigger]");
  const content = page.locator("[data-gsxui-slot-dropdown-menu-content]");
  const item = page.locator(
    '[data-gsxui-slot-dropdown-menu-checkbox-item][data-value="statusbar"]',
  );

  await trigger.click();
  await expect(trigger).toHaveAttribute("aria-expanded", "true");
  expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  await expect(item).toHaveAttribute("aria-checked", "false");

  await item.click();

  await expect(item).toHaveAttribute("aria-checked", "true");
  expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  expect(await page.evaluate(() => (window as any).__dropdownChanges)).toEqual([
    { checked: true, value: "statusbar" },
  ]);
});
