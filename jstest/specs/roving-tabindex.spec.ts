import { expect, test } from "@playwright/test";

// Roving tabindex must survive the shared MutationObserver's re-init: the
// same tabindex attribute writes that MOVE the tab stop (setActiveTrigger /
// setActiveTabStop) are themselves mutations under the bar/root, which
// schedule the module's own normalize() one microtask later. Without a
// guard, normalize() resets the tab stop back to the first (menubar) /
// pressed-or-first (toggle-group) item, undoing the arrow-key move.
test.describe("roving tabindex survives re-init", () => {
  test("menubar: ArrowRight moves the tab stop and it stays moved", async ({ page }) => {
    await page.goto("/x/menubar/basic");
    const triggers = page.locator("[data-gsxui-slot-menubar-trigger]");
    await triggers.first().focus();
    await page.keyboard.press("ArrowRight");
    // Past the microtask flush the (buggy) normalize() would land in.
    await page.waitForTimeout(50);
    await expect(triggers.nth(0)).toHaveAttribute("tabindex", "-1");
    await expect(triggers.nth(1)).toHaveAttribute("tabindex", "0");
    await expect(triggers.nth(1)).toBeFocused();
  });

  test("toggle-group: ArrowRight moves the tab stop and it stays moved", async ({ page }) => {
    await page.goto("/x/toggle-group/basic");
    const items = page.locator("[data-gsxui-slot-toggle-group-item]");
    await items.first().focus();
    await page.keyboard.press("ArrowRight");
    await page.waitForTimeout(50);
    await expect(items.nth(0)).toHaveAttribute("tabindex", "-1");
    await expect(items.nth(1)).toHaveAttribute("tabindex", "0");
    await expect(items.nth(1)).toBeFocused();
  });
});
