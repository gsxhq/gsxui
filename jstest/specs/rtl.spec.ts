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

  test("sidebar trigger icon flips under RTL", async ({ page }) => {
    await page.goto("/f/sidebar-contract?case=offcanvas-left");
    await setRTL(page);

    const trigger = page.locator("[data-gsxui-slot-sidebar-trigger]").first();
    await expect(trigger.locator("svg")).toHaveClass(/rtl:rotate-180/);
  });

  test("mobile sidebar (side=left) opens from the physical left under RTL", async ({
    page,
  }) => {
    await page.setViewportSize({ width: 500, height: 800 });
    await page.goto("/f/sidebar-contract?case=offcanvas-left");
    await setRTL(page);

    const trigger = page.locator("[data-gsxui-slot-sidebar-trigger]").first();
    await trigger.click();

    const content = page.locator("[data-gsxui-slot-sidebar-mobile-content]");
    await expect(content).toBeVisible();
    // side="left" is physical: the sheet must hug the physical left edge of
    // the viewport regardless of dir="rtl" (allow a few px for the dialog's
    // native UA border/backdrop inset). Wait out the 200ms slide-in
    // transition before sampling geometry.
    await expect
      .poll(async () => {
        const box = await content.boundingBox();
        return box ? Math.abs(box.x) : Number.NaN;
      })
      .toBeLessThan(5);
  });

  test("sheet side=left opens from the physical left under RTL, with its close button at a logical (visual-left) position", async ({
    page,
  }) => {
    await page.goto("/x/sheet/directions");
    await setRTL(page);

    await page.getByRole("button", { name: "left", exact: true }).click();
    const content = page.locator(
      '[data-gsxui-slot-sheet-content][data-side="left"]',
    );
    await expect(content).toBeVisible();
    // side="left" is physical: the sheet must hug the physical left edge of
    // the viewport regardless of dir="rtl". Wait out the slide-in transition
    // before sampling geometry.
    await expect
      .poll(async () => {
        const box = await content.boundingBox();
        return box ? Math.abs(box.x) : Number.NaN;
      })
      .toBeLessThan(5);

    const contentBox = await content.boundingBox();
    if (!contentBox) throw new Error("sheet content has no bounding box");

    const closeButton = content.locator("[data-gsxui-slot-sheet-close-button]");
    const closeBox = await closeButton.boundingBox();
    if (!closeBox) throw new Error("close button has no bounding box");
    const contentMidpoint = contentBox.x + contentBox.width / 2;
    // The close button sits at the logical end (interior positioning: end-3
    // in ui/sheet.gsx), which under RTL is the visual LEFT side of the
    // content box — not the physical right edge a plain "right-3" would put
    // it at.
    expect(closeBox.x + closeBox.width / 2).toBeLessThan(contentMidpoint);
  });

  test("input-otp: digit slot order stays LTR under RTL", async ({ page }) => {
    await page.goto("/x/input-otp/basic");
    await setRTL(page);

    const slots = page.locator("[data-gsxui-slot-input-otp-slot]");
    const first = slots.first();
    const others = await slots.all();

    const firstBox = await first.boundingBox();
    if (!firstBox) throw new Error("first OTP slot has no bounding box");

    // OTP codes are numeric tokens, not natural-language RTL text — the
    // slot group forces dir="ltr" (ui/input-otp.gsx's InputOTPGroup) so the
    // first digit's slot stays the leftmost slot even inside a dir="rtl"
    // document, instead of migrating to the visual right the way RTL prose
    // would.
    for (const other of others.slice(1)) {
      const box = await other.boundingBox();
      if (!box) throw new Error("OTP slot has no bounding box");
      expect(firstBox.x).toBeLessThan(box.x);
    }
  });

  test("calendar: arrow-key day navigation mirrors under RTL", async ({ page }) => {
    await page.goto("/x/calendar/basic");
    await setRTL(page);

    const dayFor = (iso: string) =>
      page.locator(`[data-gsxui-slot-calendar-day-button][data-date="${iso}"]`);
    const focusedDate = () =>
      page.evaluate(() => document.activeElement?.getAttribute("data-date") ?? null);

    await dayFor("2026-01-15").focus();

    // RTL negates the day-navigation direction: visually-left ArrowLeft
    // moves focus FORWARD a day (into the next day), and ArrowRight moves
    // it BACKWARD a day — same convention as tabs.js/menubar.js's roving
    // focus under RTL.
    await page.keyboard.press("ArrowLeft");
    expect(await focusedDate()).toBe("2026-01-16");

    await page.keyboard.press("ArrowRight");
    expect(await focusedDate()).toBe("2026-01-15");
  });

  // Slider's fill is a native <input type=range>'s ::-webkit-slider-runnable-
  // track / ::-moz-range-track pseudo-element background-image (a
  // linear-gradient split at var(--fill)), reversed to() left under an
  // unconditional `rtl:[&::-webkit-slider-runnable-track]:bg-[linear-
  // gradient(to_left,...)]` class arm (slider.gsx). Chromium's
  // getComputedStyle does not resolve `::-webkit-slider-runnable-track` (or
  // `::-moz-range-track`, Firefox-only) at all — probed directly, it
  // returns backgroundImage: "none" — so there is no way to sample the
  // rendered gradient direction from script in headless Chromium. This test
  // therefore guarantees only the two preconditions the CSS rule depends on
  // — the rtl: class arm is present in the shipped markup, and the input's
  // resolved direction is actually "rtl" so that Tailwind's :dir(rtl)-based
  // rtl: variant matches — NOT that the fill visually renders on the right.
  test("slider: rtl fill class arm is present and the input resolves rtl direction", async ({
    page,
  }) => {
    await page.goto("/x/slider/basic");
    await setRTL(page);

    const slider = page.locator("[data-gsxui-slot-slider]").first();

    const direction = await slider.evaluate((el) => getComputedStyle(el).direction);
    expect(direction).toBe("rtl");

    const className = await slider.evaluate((el) => el.className);
    expect(className).toMatch(/(?:^|\s)rtl:\[&::-webkit-slider-runnable-track\]/);
    expect(className).toMatch(/(?:^|\s)rtl:\[&::-moz-range-track\]/);
  });
});
