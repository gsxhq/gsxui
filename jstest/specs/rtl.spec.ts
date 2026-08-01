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
    // The click leaves the real mouse cursor resting near the content's
    // first item. menubar.js's HOVERABLE pointerover handler (by design —
    // "hover highlight IS focus", see the file's own header) refocuses
    // whatever item the cursor ends up over once layout reflows (e.g. when
    // the submenu below closes) — a synthetic pointerover Chromium fires on
    // hit-test recompute, no actual pointer movement needed. That raced this
    // otherwise keyboard-only test nondeterministically (menubar.js's File
    // menu's first item sits directly under where the trigger click landed,
    // unlike dropdown/context-menu's own submenu tests). Move the real
    // cursor off the menu so only the keyboard drives focus from here.
    await page.mouse.move(0, 0);
    const subTrigger = page.locator("[data-gsxui-slot-menubar-sub-trigger]").first();
    const subContent = page.locator("[data-gsxui-slot-menubar-sub-content]").first();
    await subTrigger.focus();

    await page.keyboard.press("ArrowLeft");
    await expect(subContent).toHaveAttribute("data-state", "open");

    await page.keyboard.press("ArrowRight");
    await expect(subContent).toHaveAttribute("data-state", "closed");
    await expect(subTrigger).toBeFocused();
  });

  test("tabs: roving focus negates direction under RTL", async ({ page }) => {
    await page.goto("/x/tabs/basic");
    await setRTL(page);

    const triggers = page.locator("[data-gsxui-slot-tabs-trigger]");
    const first = triggers.nth(0);
    const second = triggers.nth(1);

    await first.focus();
    // RTL negates roving direction: visually-left ArrowLeft moves focus
    // FORWARD to the second tab.
    await page.keyboard.press("ArrowLeft");
    await expect(second).toBeFocused();
    await page.keyboard.press("ArrowRight");
    await expect(first).toBeFocused();
  });

  test("carousel: next/prev scroll math is mirrored under RTL", async ({
    page,
  }) => {
    await page.goto("/x/carousel/basic");
    await setRTL(page);

    const root = page.locator("[data-gsxui-slot-carousel]");
    const viewport = page.locator("[data-gsxui-slot-carousel-content]");

    await expect(root).toHaveAttribute("data-current-index", "0");
    await page.locator("[data-gsxui-slot-carousel-next]").click();
    // Under RTL, scrollLeft is 0-or-negative (CSSOM spec) — "next" must move
    // it further from 0, not toward it.
    await expect
      .poll(() => viewport.evaluate((el) => Math.abs(el.scrollLeft)))
      .toBeGreaterThan(0);
    await expect(root).toHaveAttribute("data-current-index", "1");

    await page.locator("[data-gsxui-slot-carousel-previous]").click();
    await expect
      .poll(() => viewport.evaluate((el) => Math.abs(el.scrollLeft)))
      .toBeLessThan(1);
    await expect(root).toHaveAttribute("data-current-index", "0");
  });

  test("resizable: ArrowLeft grows the first (inline-start) panel under RTL", async ({
    page,
  }) => {
    await page.goto("/x/resizable/basic");
    await setRTL(page);

    const group = page.locator("[data-gsxui-slot-resizable-panel-group]").first();
    const panels = group.locator(":scope > [data-gsxui-slot-resizable-panel]");
    const handle = group.locator(":scope > [data-gsxui-slot-resizable-handle]");
    const first = panels.first();

    await handle.focus();
    // Native-direction convention: keys mirror under RTL, so the
    // visually-left ArrowLeft key still GROWS the first (inline-start,
    // visually-right-under-RTL) panel, same as ArrowRight does under LTR.
    await page.keyboard.press("ArrowLeft");
    await expect
      .poll(() => first.evaluate((el: HTMLElement) => parseFloat(el.style.flexGrow)))
      .toBeGreaterThan(50);

    await page.keyboard.press("ArrowRight");
    await expect
      .poll(() => first.evaluate((el: HTMLElement) => parseFloat(el.style.flexGrow)))
      .toBeCloseTo(50, 0);
  });

  test("resizable: dragging the handle toward the viewport-left grows the first panel under RTL", async ({
    page,
  }) => {
    await page.goto("/x/resizable/basic");
    await setRTL(page);

    const group = page.locator("[data-gsxui-slot-resizable-panel-group]").first();
    const panels = group.locator(":scope > [data-gsxui-slot-resizable-panel]");
    const handle = group.locator(":scope > [data-gsxui-slot-resizable-handle]");
    const first = panels.first();

    const box = await handle.boundingBox();
    if (!box) throw new Error("resizable handle has no bounding box");
    const startX = box.x + box.width / 2;
    const startY = box.y + box.height / 2;

    await page.mouse.move(startX, startY);
    await page.mouse.down();
    // Dragging toward the viewport-left (smaller clientX) must GROW the
    // first (inline-start, visually-right-under-RTL) panel.
    await page.mouse.move(startX - 50, startY);
    await page.mouse.up();

    await expect
      .poll(() => first.evaluate((el: HTMLElement) => parseFloat(el.style.flexGrow)))
      .toBeGreaterThan(50);
  });
});
