import { expect, test } from "../support/fixtures";

test("sidebar documentation contains each app shell inside its own viewport", async ({ page }) => {
  await page.goto("/components/sidebar");

  const previews = page.locator("iframe[data-site-isolated-preview]");
  await expect(previews).toHaveCount(10);
  await expect(page.locator("[data-gsxui-slot-sidebar-container]")).toHaveCount(0);
  await expect(page.locator('a[href="/components/button"]').first()).toBeVisible();

  const basic = page.frameLocator('iframe[src="/examples/sidebar/basic"]');
  const wrapper = basic.locator("[data-gsxui-sidebar-wrapper]");
  const container = basic.locator("[data-gsxui-slot-sidebar-container]");
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
  await expect(floating.locator("[data-gsxui-sidebar-wrapper]")).toHaveCount(1);
  await expect(floating.locator("[data-gsxui-slot-sidebar-container]")).toHaveCount(1);

  const persistedFrame = page.locator('iframe[src="/examples/sidebar/persisted"]');
  await persistedFrame.scrollIntoViewIfNeeded();
  const persisted = page.frameLocator('iframe[src="/examples/sidebar/persisted"]');
  const persistedHeight = await persisted
    .locator("[data-gsxui-sidebar-wrapper]")
    .evaluate((element) => ({
      wrapper: element.getBoundingClientRect().height,
      viewport: document.documentElement.clientHeight,
    }));
  expect(Math.abs(persistedHeight.wrapper - persistedHeight.viewport)).toBeLessThanOrEqual(1);
});
