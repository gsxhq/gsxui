import { expect, test } from "../support/fixtures";

// Pins for gsxui.js's shared anchored-positioning engine (position/release):
// main-axis flip with data-side tracking, scroll/resize reposition while
// open, the --gsxui-available-height cap on tall list surfaces, and
// listener hygiene (attach once per open, detach on close). Each measured
// pin was verified to FAIL with the corresponding feature stubbed out of
// ui/gsxui.js before being trusted green.

test("a dropdown opened near the bottom edge flips above its trigger and data-side tracks it", async ({
  page,
}) => {
  await page.goto("/x/dropdown-menu/checkboxes");
  const trigger = page.locator("[data-gsxui-slot-dropdown-menu-trigger]");
  const content = page.locator("[data-gsxui-slot-dropdown-menu-content]");
  await expect(content).toHaveAttribute("data-side", "bottom");

  // Push the trigger near the bottom edge of the 900px viewport: the menu
  // (~200px tall) can't fit below, and the room above is far larger.
  await page.evaluate(() => {
    document.body.style.paddingTop = "800px";
  });
  await trigger.click();
  await expect(content).toBeVisible();

  const t = (await trigger.boundingBox())!;
  const c = (await content.boundingBox())!;
  expect(c.y + c.height, "content bottom at or above trigger top").toBeLessThanOrEqual(
    t.y + 1,
  );
  // The enter/exit animations key on data-[side=...]; a flipped placement
  // must report the placed side, exactly as Radix recomputes placement.
  await expect(content).toHaveAttribute("data-side", "top");

  // Closing restores the authored preference for the next open.
  await page.keyboard.press("Escape");
  await expect(content).toHaveAttribute("data-state", "closed");
  await expect(content).toHaveAttribute("data-side", "bottom");
});

test("an open popover tracks its trigger through page scroll", async ({ page }) => {
  await page.goto("/x/popover/basic");
  // Give the page something to scroll.
  await page.evaluate(() => {
    document.body.style.paddingTop = "200px";
    document.body.style.minHeight = "3000px";
  });
  const trigger = page.locator("[data-gsxui-slot-popover-trigger]");
  const content = page.locator("[data-gsxui-slot-popover-content]");
  await trigger.click();
  await expect(content).toBeVisible();

  // Layout metrics (style.top/left), not boundingBox: mid-enter the
  // scale-95 transition skews the painted box by a few px, while layout
  // position is what the engine computed.
  const contentPos = () =>
    content.evaluate((el: HTMLElement) => ({
      top: parseFloat(el.style.top),
      left: parseFloat(el.style.left),
    }));
  const before = { trigger: (await trigger.boundingBox())!, content: await contentPos() };
  const offsetBefore = before.content.top - before.trigger.y;

  await page.evaluate(() => window.scrollBy(0, 120));
  await expect
    .poll(async () => (await trigger.boundingBox())!.y, {
      message: "trigger moved with the scroll",
    })
    .toBeLessThan(before.trigger.y - 100);

  // The content keeps its relative offset to the trigger — it repositioned
  // rather than staying at its viewport-fixed open coordinates. Polled: the
  // scroll event that drives the reposition dispatches asynchronously after
  // scrollBy, while the trigger's own box reflects the scroll immediately.
  await expect
    .poll(
      async () => {
        const t = (await trigger.boundingBox())!;
        const c = await contentPos();
        return Math.abs(c.top - t.y - offsetBefore);
      },
      { message: "content tracked the trigger" },
    )
    .toBeLessThanOrEqual(1);
  expect(Math.abs((await contentPos()).left - before.content.left)).toBeLessThanOrEqual(1);
});

test("a Select near the bottom edge is capped to the available height and scrolls internally", async ({
  page,
}) => {
  await page.goto("/x/select/scrollable");
  // Shrink the viewport so the room below the trigger (and above it) is far
  // smaller than the list's natural (max-h-96, 384px) height. The room
  // above is smaller than below, so the placement stays side=bottom and the
  // cap applies to the space under the trigger.
  await page.setViewportSize({ width: 1280, height: 360 });
  const trigger = page.locator("[data-gsxui-slot-select-trigger]");
  const content = page.locator("[data-gsxui-slot-select-content]");
  await trigger.click();
  await expect(content).toBeVisible();

  // Layout metrics (style.top / offsetHeight), not boundingBox: mid-enter
  // the scale-95 transition shrinks the painted box by a few px (the same
  // discrimination the theme-editor viewport pin documents), while layout
  // position/size are what the engine actually computed.
  const t = (await trigger.boundingBox())!;
  const c = await content.evaluate((el: HTMLElement) => ({
    top: parseFloat(el.style.top),
    height: el.offsetHeight,
  }));
  expect(c.top, "still placed below the trigger").toBeGreaterThanOrEqual(t.y + t.height);
  expect(c.height, "height capped to the available space").toBeLessThanOrEqual(
    360 - (t.y + t.height),
  );
  expect(c.top + c.height, "bottom edge inside the viewport").toBeLessThanOrEqual(360);
  const scrollable = await content.evaluate(
    (el) => el.scrollHeight > el.clientHeight,
  );
  expect(scrollable, "list scrolls internally").toBe(true);
});

test("reposition listeners attach once per open and detach on close", async ({
  page,
}) => {
  // getEventListeners is a DevTools-console-only API; the CDP
  // DOMDebugger.getEventListeners domain is its programmatic form and the
  // only way page-side code can observe window's listener table.
  await page.goto("/x/popover/basic");
  const client = await page.context().newCDPSession(page);
  const countWindowListeners = async () => {
    const { result } = await client.send("Runtime.evaluate", {
      expression: "window",
    });
    const { listeners } = await client.send("DOMDebugger.getEventListeners", {
      objectId: result.objectId!,
    });
    return listeners.filter((l) => l.type === "scroll" || l.type === "resize")
      .length;
  };

  const baseline = await countWindowListeners();
  const trigger = page.locator("[data-gsxui-slot-popover-trigger]");
  const content = page.locator("[data-gsxui-slot-popover-content]");

  for (let i = 0; i < 3; i++) {
    await trigger.click();
    await expect(content).toBeVisible();
    expect(await countWindowListeners(), "one scroll + one resize while open").toBe(
      baseline + 2,
    );
    await page.keyboard.press("Escape");
    await expect(content).toHaveAttribute("data-state", "closed");
    expect(await countWindowListeners(), "detached on close, no leak").toBe(baseline);
  }
});
