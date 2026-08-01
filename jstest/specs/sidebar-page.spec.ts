import { expect, test } from "../support/fixtures";

test("sidebar basic keeps separators and trailing controls inside their rows", async ({
  page,
}) => {
  await page.goto("/components/sidebar");

  const basic = page.frameLocator('iframe[src="/examples/sidebar/basic"]');
  const desktop = basic.locator("[data-gsxui-slot-sidebar-desktop]");
  const inner = desktop.locator("[data-gsxui-slot-sidebar-inner]");
  const separator = desktop.locator("[data-gsxui-slot-sidebar-separator]");

  const separatorGeometry = await Promise.all([
    inner.evaluate((element) => {
      const rect = element.getBoundingClientRect();
      return { left: rect.left, right: rect.right };
    }),
    separator.evaluate((element) => {
      const rect = element.getBoundingClientRect();
      return { left: rect.left, right: rect.right };
    }),
  ]);
  expect(separatorGeometry[1].left - separatorGeometry[0].left).toBeCloseTo(8);
  expect(separatorGeometry[0].right - separatorGeometry[1].right).toBeCloseTo(8);

  const decorationGeometry = await desktop
    .locator("[data-gsxui-slot-sidebar-menu-item]")
    .evaluateAll((items) =>
      items.flatMap((item) => {
        const itemRect = item.getBoundingClientRect();
        if (itemRect.width === 0 || itemRect.height === 0) return [];

        const button = item.querySelector("[data-gsxui-slot-sidebar-menu-button]");
        if (!(button instanceof HTMLElement)) return [];
        const buttonRect = button.getBoundingClientRect();
        const decorations = [
          ...item.querySelectorAll(
            ":scope > [data-gsxui-slot-sidebar-menu-action], " +
              ":scope > [data-gsxui-slot-sidebar-menu-badge]",
          ),
        ].filter((element) => element.getBoundingClientRect().width > 0);

        return decorations.map((decoration, index) => {
          const rect = decoration.getBoundingClientRect();
          const overlapsAnother = decorations.some((other, otherIndex) => {
            if (index === otherIndex) return false;
            const otherRect = other.getBoundingClientRect();
            return (
              rect.left < otherRect.right &&
              rect.right > otherRect.left &&
              rect.top < otherRect.bottom &&
              rect.bottom > otherRect.top
            );
          });
          return {
            withinButtonRow:
              rect.top >= buttonRect.top && rect.bottom <= buttonRect.bottom,
            overlapsAnother,
          };
        });
      }),
    );
  expect(decorationGeometry.length).toBeGreaterThan(0);
  expect(decorationGeometry.every((entry) => entry.withinButtonRow)).toBe(true);
  expect(decorationGeometry.every((entry) => !entry.overlapsAnother)).toBe(true);
});

test("sidebar icon collapse keeps a compact brand mark without leaking its name", async ({
  page,
}) => {
  await page.goto("/components/sidebar");

  const icon = page.frameLocator(
    'iframe[src="/examples/sidebar/variants?_preview=icon-collapsed"]',
  );
  await page
    .locator('iframe[src="/examples/sidebar/variants?_preview=icon-collapsed"]')
    .scrollIntoViewIfNeeded();
  const brand = icon
    .locator("[data-gsxui-slot-sidebar-desktop]")
    .locator("[data-sidebar-example-brand]");
  const mark = brand.locator("[data-sidebar-example-brand-mark]");
  const name = brand.locator("[data-sidebar-example-brand-name]");

  await expect(mark).toHaveText("A");
  await expect(name).toHaveText("Acme Inc");

  const geometry = await brand.evaluate((element) => {
    const mark = element.querySelector("[data-sidebar-example-brand-mark]");
    const name = element.querySelector("[data-sidebar-example-brand-name]");
    if (!(mark instanceof HTMLElement) || !(name instanceof HTMLElement)) {
      throw new Error("sidebar example brand is missing its mark or name");
    }
    const rect = element.getBoundingClientRect();
    const markRect = mark.getBoundingClientRect();
    const nameRect = name.getBoundingClientRect();
    return {
      width: rect.width,
      overflowX: getComputedStyle(element).overflowX,
      markInside: markRect.left >= rect.left && markRect.right <= rect.right,
      nameStartsOutside: nameRect.left >= rect.right,
    };
  });
  expect(geometry.width).toBe(32);
  expect(geometry.overflowX).toBe("hidden");
  expect(geometry.markInside).toBe(true);
  expect(geometry.nameStartsOutside).toBe(true);
});

test("sidebar documentation contains each app shell inside its own viewport", async ({ page }) => {
  await page.goto("/components/sidebar");

  const previews = page.locator("iframe[data-site-isolated-preview]");
  await expect(previews).toHaveCount(11);
  await expect(page.locator("[data-gsxui-slot-sidebar-container]")).toHaveCount(0);
  await expect(
    page.locator('[data-site-docs-sidebar] a[href="/components/button"]'),
  ).toBeVisible();

  const basicFrame = page.locator('iframe[src="/examples/sidebar/basic"]');
  const basic = page.frameLocator('iframe[src="/examples/sidebar/basic"]');
  const wrapper = basic.locator("[data-gsxui-slot-sidebar-wrapper]");
  const container = basic.locator("[data-gsxui-slot-sidebar-container]");
  await expect(basicFrame).toHaveAttribute("width", "1024");
  await expect
    .poll(() =>
      basic
        .locator("html")
        .evaluate(() => document.documentElement.clientWidth),
    )
    .toBe(1024);
  const previewGeometry = await basicFrame
    .locator("..")
    .evaluate((surface) => {
      const frame = surface.querySelector("iframe");
      if (!(frame instanceof HTMLIFrameElement)) {
        throw new Error("isolated preview surface is missing its iframe");
      }
      return {
        surfaceWidth: surface.getBoundingClientRect().width,
        frameWidth: frame.getBoundingClientRect().width,
        overflowX: getComputedStyle(surface).overflowX,
      };
    });
  expect(previewGeometry.surfaceWidth).toBe(640);
  expect(previewGeometry.frameWidth).toBe(1024);
  expect(previewGeometry.overflowX).toBe("hidden");
  const outerDocumentWidth = await page.locator("html").evaluate(() => ({
    client: document.documentElement.clientWidth,
    scroll: document.documentElement.scrollWidth,
  }));
  expect(outerDocumentWidth.scroll).toBe(outerDocumentWidth.client);
  await expect(wrapper).toHaveAttribute("data-state", "expanded");
  await expect(container).toBeVisible();

  const geometry = await container.evaluate((element) => {
    const rect = element.getBoundingClientRect();
    return {
      position: getComputedStyle(element).position,
      left: rect.left,
      right: rect.right,
      viewportWidth: document.documentElement.clientWidth,
    };
  });
  expect(geometry.position).toBe("fixed");
  expect(geometry.left).toBeGreaterThanOrEqual(0);
  expect(geometry.right).toBeLessThanOrEqual(geometry.viewportWidth);

  await basic.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  await expect(wrapper).toHaveAttribute("data-state", "collapsed");

  const parentIsDark = await page.locator("html").evaluate((element) => element.classList.contains("dark"));
  await page.getByRole("button", { name: "Toggle theme" }).click();
  await expect(page.locator("html")).toHaveClass(parentIsDark ? /^(?!.*\bdark\b)/ : /\bdark\b/);
  await expect(basic.locator("html")).toHaveClass(parentIsDark ? /^(?!.*\bdark\b)/ : /\bdark\b/);

  const floating = page.frameLocator(
    'iframe[src="/examples/sidebar/variants?_preview=floating"]',
  );
  await expect(floating.locator("[data-gsxui-slot-sidebar-wrapper]")).toHaveCount(1);
  await expect(floating.locator("[data-gsxui-slot-sidebar-container]")).toHaveCount(1);

  const persistedFrame = page.locator('iframe[src="/examples/sidebar/persisted"]');
  await persistedFrame.scrollIntoViewIfNeeded();
  const persisted = page.frameLocator('iframe[src="/examples/sidebar/persisted"]');
  const persistedHeight = await persisted
    .locator("[data-gsxui-slot-sidebar-wrapper]")
    .evaluate((element) => ({
      wrapper: element.getBoundingClientRect().height,
      viewport: document.documentElement.clientHeight,
    }));
  expect(Math.abs(persistedHeight.wrapper - persistedHeight.viewport)).toBeLessThanOrEqual(1);
});
