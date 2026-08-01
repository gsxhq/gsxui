// RTL keyboard-semantics coverage: under dir="rtl", the menu-family modules
// (dropdown-menu.js, context-menu.js, menubar.js) swap ArrowLeft/ArrowRight
// meaning per WAI-ARIA's "into the menu" / "out of the menu" convention —
// ArrowLeft opens a submenu, ArrowRight closes it, mirroring LTR's own
// ArrowRight-opens/ArrowLeft-closes. setRTL is exported for later RTL tasks
// (5-8) to reuse against their own components.
import type { Page } from "@playwright/test";
import { expect, test } from "../support/fixtures";

export async function setRTL(page: Page) {
  await page.evaluate(() => document.documentElement.setAttribute("dir", "rtl"));
}

test.describe("rtl", () => {
  test("dropdown submenu: ArrowLeft opens, ArrowRight closes and returns focus", async ({
    page,
  }) => {
    await page.goto("/x/dropdown-menu/submenu");
    await setRTL(page);

    const trigger = page.locator("[data-gsxui-slot-dropdown-menu-trigger]");
    const content = page.locator("[data-gsxui-slot-dropdown-menu-content]");
    const subTrigger = page.locator(
      "[data-gsxui-slot-dropdown-menu-sub-trigger]",
    );
    const subContent = page.locator(
      "[data-gsxui-slot-dropdown-menu-sub-content]",
    );

    await trigger.click();
    await expect(content).toHaveAttribute("data-state", "open");
    await subTrigger.focus();

    await page.keyboard.press("ArrowLeft");
    await expect(subContent).toHaveAttribute("data-state", "open");
    expect(await subContent.evaluate((el) => el.matches(":popover-open"))).toBe(
      true,
    );

    await page.keyboard.press("ArrowRight");
    await expect(subContent).toHaveAttribute("data-state", "closed");
    await expect(subTrigger).toBeFocused();
  });

  test("context-menu submenu: ArrowLeft opens, ArrowRight closes and returns focus", async ({
    page,
  }) => {
    await page.goto("/x/context-menu/full");
    await setRTL(page);

    const trigger = page.locator("[data-gsxui-slot-context-menu-trigger]");
    const content = page.locator("[data-gsxui-slot-context-menu-content]");
    const subTrigger = page.locator(
      "[data-gsxui-slot-context-menu-sub-trigger]",
    );
    const subContent = page.locator(
      "[data-gsxui-slot-context-menu-sub-content]",
    );

    await trigger.click({ button: "right" });
    await expect(content).toHaveAttribute("data-state", "open");
    await subTrigger.focus();

    await page.keyboard.press("ArrowLeft");
    await expect(subContent).toHaveAttribute("data-state", "open");

    await page.keyboard.press("ArrowRight");
    await expect(subContent).toHaveAttribute("data-state", "closed");
    await expect(subTrigger).toBeFocused();
  });

  test("menubar: top-level roving negates direction under RTL", async ({
    page,
  }) => {
    await page.goto("/x/menubar/full");
    await setRTL(page);

    const triggers = page.locator("[data-gsxui-slot-menubar-trigger]");
    const fileTrigger = triggers.first();
    const editTrigger = triggers.nth(1);

    await fileTrigger.focus();
    // RTL negates roving direction: ArrowLeft moves toward the NEXT trigger.
    await page.keyboard.press("ArrowLeft");
    await expect(editTrigger).toBeFocused();
    await page.keyboard.press("ArrowRight");
    await expect(fileTrigger).toBeFocused();
  });

  test("menubar: submenu ArrowLeft opens, ArrowRight closes and returns focus under RTL", async ({
    page,
  }) => {
    await page.goto("/x/menubar/full");
    await setRTL(page);

    const fileTrigger = page.locator("[data-gsxui-slot-menubar-trigger]").first();
    await fileTrigger.click();
    const subTrigger = page.locator("[data-gsxui-slot-menubar-sub-trigger]").first();
    const subContent = page.locator("[data-gsxui-slot-menubar-sub-content]").first();
    await subTrigger.focus();

    await page.keyboard.press("ArrowLeft");
    await expect(subContent).toHaveAttribute("data-state", "open");

    await page.keyboard.press("ArrowRight");
    await expect(subContent).toHaveAttribute("data-state", "closed");
    await expect(subTrigger).toBeFocused();
  });
});
