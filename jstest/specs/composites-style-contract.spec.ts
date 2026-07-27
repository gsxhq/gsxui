import { expect, test } from "../support/fixtures";

test("ButtonGroup restores only the visible tail before an aria-hidden select", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status()).toBe(200);

  const group = page.locator('[data-style-contract="button-group-tail"]');
  const radii = await group.evaluate((element) => {
    const corners = (selector: string) => {
      const css = getComputedStyle(element.querySelector(selector)!);
      return {
        topRight: css.borderTopRightRadius,
        bottomRight: css.borderBottomRightRadius,
      };
    };
    return {
      earlier: corners('[data-style-contract="button-group-earlier"]'),
      visibleTail: corners('[data-style-contract="button-group-visible-tail"]'),
    };
  });

  expect(radii).toEqual({
    earlier: { topRight: "0px", bottomRight: "0px" },
    visibleTail: { topRight: "10px", bottomRight: "10px" },
  });
});

test("horizontal ButtonGroup separators stay on the cross-axis", async ({ page }) => {
  const response = await page.goto("/x/button-group/basic");
  expect(response?.status()).toBe(200);

  const groups = page.locator('[data-gsxui-slot-button-group]');
  await expect(groups).toHaveCount(4);
  const clipboard = groups.nth(1);
  const separator = clipboard.locator(
    ':scope > [data-gsxui-slot-button-group-separator]',
  );

  await expect(separator).toHaveAttribute("data-orientation", "vertical");
  expect(
    await clipboard.evaluate((group) => {
      const groupRect = group.getBoundingClientRect();
      const separatorRect = group
        .querySelector('[data-gsxui-slot-button-group-separator]')!
        .getBoundingClientRect();
      return {
        separatorWidth: separatorRect.width,
        separatorHeight: separatorRect.height,
        groupHeight: groupRect.height,
        childrenInside: [...group.children].every((child) => {
          const rect = child.getBoundingClientRect();
          return (
            rect.left >= groupRect.left &&
            rect.right <= groupRect.right &&
            rect.top >= groupRect.top &&
            rect.bottom <= groupRect.bottom
          );
        }),
      };
    }),
  ).toEqual({
    separatorWidth: 1,
    separatorHeight: 32,
    groupHeight: 32,
    childrenInside: true,
  });
});

test("vertical ButtonGroup separators stay on the cross-axis", async ({ page }) => {
  const response = await page.goto("/x/button-group/basic");
  expect(response?.status()).toBe(200);

  const group = page.getByRole("group", { name: "Media controls" });
  const children = group.locator(":scope > *");
  await expect(children).toHaveCount(3);
  const separator = group.locator(
    ':scope > [data-gsxui-slot-button-group-separator]',
  );
  await expect(separator).toHaveAttribute("data-orientation", "horizontal");
  expect(
    await group.evaluate((element) => {
      const groupRect = element.getBoundingClientRect();
      const separatorRect = element
        .querySelector('[data-gsxui-slot-button-group-separator]')!
        .getBoundingClientRect();
      const button = element.lastElementChild!;
      const css = getComputedStyle(button);
      return {
        separatorWidth: separatorRect.width,
        separatorHeight: separatorRect.height,
        groupWidth: groupRect.width,
        bottomLeft: css.borderBottomLeftRadius,
        bottomRight: css.borderBottomRightRadius,
      };
    }),
  ).toEqual({
    separatorWidth: 32,
    separatorHeight: 1,
    groupWidth: 32,
    bottomLeft: "10px",
    bottomRight: "10px",
  });
});

test("Carousel keeps native-scroll mechanics and caller spacing overrides", async ({ page }) => {
  const response = await page.goto("/x/carousel/sizes");
  expect(response?.status()).toBe(200);

  const root = page.locator("[data-gsxui-carousel]");
  const viewport = page.locator("[data-gsxui-carousel-content]");
  const track = page.locator('[data-gsxui-slot-carousel-track]');
  const item = page.locator("[data-gsxui-carousel-item]").first();
  const previous = page.locator("[data-gsxui-carousel-prev]");
  const next = page.locator("[data-gsxui-carousel-next]");

  expect(
    await viewport.evaluate((el) => {
      const css = getComputedStyle(el);
      return { overflowX: css.overflowX, snap: css.scrollSnapType };
    }),
  ).toEqual({ overflowX: "auto", snap: "x mandatory" });
  expect(
    await track.evaluate((el) => {
      const css = getComputedStyle(el);
      return { display: css.display, marginLeft: css.marginLeft };
    }),
  ).toEqual({ display: "flex", marginLeft: "-4px" });
  expect(
    await item.evaluate((el) => {
      const css = getComputedStyle(el);
      return {
        flexBasis: css.flexBasis,
        flexGrow: css.flexGrow,
        flexShrink: css.flexShrink,
        paddingLeft: css.paddingLeft,
      };
    }),
  ).toEqual({
    flexBasis: "33.3333%",
    flexGrow: "0",
    flexShrink: "0",
    paddingLeft: "4px",
  });

  await expect(root).toHaveAttribute("data-current-index", "0");
  await expect(previous).toBeDisabled();
  await expect(next).toBeEnabled();
  await root.evaluate(
    (el: HTMLElement & { gsxuiCarousel: { scrollTo(index: number): void } }) =>
      el.gsxuiCarousel.scrollTo(4),
  );
  await expect(previous).toBeEnabled();
  await expect(next).toBeDisabled();
});

test("registered vertical Carousel covers every oriented style slot", async ({ page }) => {
  const response = await page.goto("/x/carousel/vertical");
  expect(response?.status()).toBe(200);

  const orientedSlots = [
    "carousel",
    "carousel-content",
    "carousel-track",
    "carousel-item",
    "carousel-previous",
    "carousel-next",
  ];
  for (const slot of orientedSlots) {
    await expect(
      page.locator(`[data-gsxui-slot-${slot}]`).first(),
      `${slot} must reflect the production example's vertical axis`,
    ).toHaveAttribute("data-orientation", "vertical");
  }

  const viewport = page.locator("[data-gsxui-carousel-content]");
  const track = page.locator('[data-gsxui-slot-carousel-track]');
  expect(
    await viewport.evaluate((el) => {
      const css = getComputedStyle(el);
      return { overflowY: css.overflowY, snap: css.scrollSnapType };
    }),
  ).toEqual({ overflowY: "auto", snap: "y mandatory" });
  await expect
    .poll(() => track.evaluate((el) => getComputedStyle(el).flexDirection))
    .toBe("column");

  const itemGeometry = await page
    .locator("[data-gsxui-carousel-item]")
    .evaluateAll((items) =>
      items.slice(0, 2).map((item) => item.getBoundingClientRect().height),
    );
  expect(itemGeometry).toHaveLength(2);
  expect(itemGeometry[0]).toBeCloseTo(100, 0);
  expect(itemGeometry[1]).toBeCloseTo(100, 0);
});

test("Resizable consumes dynamic flex values and remains keyboard operable", async ({ page }) => {
  const response = await page.goto("/x/resizable/handle");
  expect(response?.status()).toBe(200);

  const group = page.locator("[data-gsxui-resizable]").first();
  const panels = group.locator(":scope > [data-gsxui-resizable-panel]");
  const handle = group.locator(":scope > [data-gsxui-resizable-handle]");
  await expect(panels).toHaveCount(2);

  expect(
    await group.evaluate((el) => getComputedStyle(el).display),
  ).toBe("flex");
  expect(
    await panels.first().evaluate((el) => {
      const css = getComputedStyle(el);
      return { basis: css.flexBasis, grow: css.flexGrow, overflow: css.overflow };
    }),
  ).toEqual({ basis: "0px", grow: "25", overflow: "hidden" });

  await expect.poll(() => handle.evaluate((el) => el.getBoundingClientRect().width)).toBe(1);
  await page.addStyleTag({
    content: `
      @layer components {
        :where(
          [data-gsxui-slot-resizable-handle][aria-orientation="vertical"]
        ) {
          width: 6px;
        }
      }
    `,
  });
  await expect.poll(() => handle.evaluate((el) => el.getBoundingClientRect().width)).toBe(6);
  expect(
    await handle.evaluate((el) => {
      const hitTarget = getComputedStyle(el, "::after");
      return { flex: getComputedStyle(el).flex, hitTargetWidth: hitTarget.width };
    }),
  ).toEqual({ flex: "0 0 auto", hitTargetWidth: "4px" });

  await handle.focus();
  await page.keyboard.press("ArrowRight");
  await expect(handle).toHaveAttribute("aria-valuenow", "35");
  await expect
    .poll(() => panels.first().evaluate((el: HTMLElement) => el.style.flexGrow))
    .toBe("35");
});

test("ScrollArea keeps native overflow while caller utilities win", async ({ page }) => {
  const response = await page.goto("/x/scroll-area/basic");
  expect(response?.status()).toBe(200);

  const area = page.locator('[data-gsxui-slot-scroll-area]');
  expect(
    await area.evaluate((el) => {
      const css = getComputedStyle(el);
      return {
        overflowX: css.overflowX,
        overflowY: css.overflowY,
        radius: css.borderRadius,
      };
    }),
  ).toEqual({
    overflowX: "auto",
    overflowY: "auto",
    radius: "8px",
  });
});
