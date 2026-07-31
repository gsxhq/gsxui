import { expect, test } from "../support/fixtures";

const foundation = (path: string) =>
  `${path}${path.includes("?") ? "&" : "?"}css=foundation`;

test("foundation mode excludes the default style and preserves overlay lifecycle", async ({
  page,
}) => {
  await page.goto(foundation("/x/dialog/basic"));
  expect(
    await page.evaluate(() =>
      getComputedStyle(document.documentElement)
        .getPropertyValue("--gsxui-default-style")
        .trim(),
    ),
  ).toBe("");

  const dialog = page.locator("[data-gsxui-slot-dialog-content]");
  await expect(dialog).not.toHaveAttribute("open");
  await expect.poll(() => dialog.evaluate((el) => getComputedStyle(el).display)).toBe("none");

  await page.getByRole("button", { name: "Delete account" }).click();
  await expect(dialog).toHaveAttribute("open");
  await expect(dialog).toHaveAttribute("data-state", "open");
  expect(
    await dialog.evaluate((el) => {
      const rect = el.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    }),
  ).toBe(true);

  await page.locator("[data-gsxui-dialog-close]").first().click();
  await expect(dialog).not.toHaveAttribute("open");
  await expect(dialog).toHaveAttribute("data-state", "closed");

  await page.goto(foundation("/x/popover/basic"));
  const popover = page.locator("[data-gsxui-slot-popover-content]");
  await expect.poll(() => popover.evaluate((el) => getComputedStyle(el).display)).toBe("none");
  await page.locator("[data-gsxui-slot-popover-trigger]").first().click();
  await expect.poll(() => popover.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  expect(
    await popover.evaluate((el) => {
      const rect = el.getBoundingClientRect();
      return rect.width > 0 && rect.height > 0;
    }),
  ).toBe(true);
  await page.keyboard.press("Escape");
  await expect.poll(() => popover.evaluate((el) => getComputedStyle(el).display)).toBe("none");
});

test("foundation keeps Carousel navigation geometrically functional", async ({ page }) => {
  await page.goto(foundation("/x/carousel/basic"));
  const root = page.locator("[data-gsxui-carousel]");
  const viewport = page.locator("[data-gsxui-carousel-content]");
  const items = page.locator("[data-gsxui-carousel-item]");
  await expect(items).toHaveCount(5);
  const geometry = await viewport.evaluate((el) => {
    const track = el.querySelector<HTMLElement>("[data-gsxui-slot-carousel-track]")!;
    const first = el.querySelector<HTMLElement>("[data-gsxui-carousel-item]")!;
    return {
      viewportWidth: el.getBoundingClientRect().width,
      trackWidth: track.getBoundingClientRect().width,
      itemWidth: first.getBoundingClientRect().width,
      scrollWidth: el.scrollWidth,
    };
  });
  expect(geometry.viewportWidth).toBeGreaterThan(0);
  // One item spans exactly one track width (foundation's own
  // `flex: 0 0 100%`). Compared against the TRACK, not the viewport: the
  // track's -ml-4 spacing is Carousel's style, and since Carousel migrated to
  // the slot axis that utility is a literal class on the markup Tailwind
  // scans, so it is present in foundation mode too and makes the track 16px
  // wider than the viewport. The mechanic under test is the flex basis, not
  // the spacing.
  expect(geometry.itemWidth).toBeCloseTo(geometry.trackWidth, 0);
  expect(geometry.scrollWidth).toBeGreaterThan(geometry.viewportWidth * 4);

  await expect(root).toHaveAttribute("data-current-index", "0");
  await page.locator("[data-gsxui-carousel-next]").click();
  await expect.poll(() => viewport.evaluate((el) => el.scrollLeft)).toBeGreaterThan(0);
  await expect(root).toHaveAttribute("data-current-index", "1");
});

for (const width of [640, 900]) {
  test(`foundation keeps Resizable keyboard and pointer geometry at ${width}px`, async ({
    page,
  }) => {
    await page.setViewportSize({ width, height: 700 });
    await page.goto(foundation("/x/resizable/handle"));
    const group = page.locator("[data-gsxui-resizable]").first();
    const panels = group.locator(":scope > [data-gsxui-resizable-panel]");
    const handle = group.locator(":scope > [data-gsxui-resizable-handle]");
    await expect(panels).toHaveCount(2);

    const before = await panels.first().evaluate((el) => el.getBoundingClientRect().width);
    expect(before).toBeGreaterThan(0);
    expect(
      await handle.evaluate((el) => {
        const rect = el.getBoundingClientRect();
        const hit = getComputedStyle(el, "::after");
        return rect.height > 0 && Number.parseFloat(hit.width) >= 4;
      }),
    ).toBe(true);

    await handle.focus();
    await page.keyboard.press("ArrowRight");
    await expect(handle).toHaveAttribute("aria-valuenow", "35");
    const keyboardWidth = await panels.first().evaluate(
      (el) => el.getBoundingClientRect().width,
    );
    expect(keyboardWidth).toBeGreaterThan(before);

    const box = await handle.boundingBox();
    if (!box) throw new Error("Resizable handle has no pointer geometry");
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2);
    await page.mouse.down();
    await page.mouse.move(box.x + box.width / 2 + 30, box.y + box.height / 2);
    await page.mouse.up();
    await expect
      .poll(() => panels.first().evaluate((el) => el.getBoundingClientRect().width))
      .toBeGreaterThan(keyboardWidth);
  });
}

test("foundation keeps InputOTP entry and client caret functional", async ({ page }) => {
  await page.goto(foundation("/x/input-otp/basic"));
  const input = page.locator("[data-gsxui-input-otp-input]");
  const slots = page.locator("[data-gsxui-input-otp-slot]");
  await input.fill("123");
  await expect(slots.nth(0)).toHaveText("1");
  await expect(slots.nth(2)).toHaveText("3");
  await expect(slots.nth(3)).toHaveAttribute("data-active", "true");
  await expect(slots.nth(3).locator("[data-gsxui-input-otp-caret-overlay]")).toHaveCount(1);
  await expect(slots.nth(3).locator("[data-gsxui-input-otp-caret]")).toHaveCount(1);
  expect(
    await input.evaluate((el) => {
      const inputRect = el.getBoundingClientRect();
      const rootRect = el.parentElement!.getBoundingClientRect();
      return inputRect.width === rootRect.width && inputRect.height === rootRect.height;
    }),
  ).toBe(true);
});

test("foundation keeps Sonner fallback creation and server-row adoption functional", async ({
  page,
}) => {
  await page.goto(foundation("/x/toaster/types"));
  await page.evaluate(() => {
    const section = document.querySelector('[aria-label="Notifications"]')!;
    for (const template of [...section.querySelectorAll("template")]) {
      document.body.append(template);
    }
    section.remove();
    window.gsxui.toast.success("Fallback", { duration: 60_000 });
    const template = document.querySelector<HTMLTemplateElement>(
      'template[data-gsxui-toast-template="info"]',
    )!;
    const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
    row.dataset.foundationServer = "true";
    row.dataset.duration = "60000";
    document.querySelector("[data-gsxui-toaster]")!.append(row);
  });

  const region = page.locator("[data-gsxui-toaster]");
  await expect(region).toHaveCount(1);
  await expect(region.locator("li[data-gsxui-toast]")).toHaveCount(2);
  await expect(region.locator('[data-foundation-server="true"]')).toHaveAttribute(
    "data-state",
    "open",
  );
  expect(
    await region.evaluate((el) => {
      const css = getComputedStyle(el);
      return {
        position: css.position,
        pointerEvents: css.pointerEvents,
        inViewport:
          el.getBoundingClientRect().right <= innerWidth &&
          el.getBoundingClientRect().bottom <= innerHeight,
      };
    }),
  ).toEqual({ position: "fixed", pointerEvents: "none", inViewport: true });
});

test("full style leaves caller utilities later in the cascade", async ({ page }) => {
  await page.goto("/f/style-contract");
  expect(
    await page.evaluate(() =>
      getComputedStyle(document.documentElement)
        .getPropertyValue("--gsxui-default-style")
        .trim(),
    ),
  ).toBe("default");
  const button = page.getByRole("button", { name: "Caller override" });
  expect(
    await button.evaluate((el) => {
      const css = getComputedStyle(el);
      return { height: css.height, radius: css.borderRadius };
    }),
  ).toEqual({ height: "48px", radius: "0px" });
});
