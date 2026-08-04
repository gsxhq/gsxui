import { expect, test } from "@playwright/test";

test.describe("boosted site navigation", () => {
  test("htmx is loaded on site pages", async ({ page }) => {
    await page.goto("/components/button");
    await expect
      .poll(async () => page.evaluate(() => typeof (window as any).htmx))
      .toBe("object");
  });

  test("sidebar scroll survives boosted navigation", async ({ page }) => {
    await page.goto("/components/button");
    const sidebar = page.locator("aside nav").first();
    // Scroll the docs sidebar container to its bottom. The scrollable
    // element is the aside's overflow container; adjust the locator to the
    // actual scrolling element (inspect site/pages/layout.gsx: the left
    // column) and note it in your report.
    const scrolled = await sidebar.evaluate((el) => {
      const scroller = el.closest("[class*='overflow-y']") ?? el;
      scroller.scrollTop = scroller.scrollHeight;
      return scroller.scrollTop;
    });
    expect(scrolled).toBeGreaterThan(0);

    await page.getByRole("link", { name: "tooltip", exact: true }).click();
    await expect(page).toHaveURL(/\/components\/tooltip$/);
    await expect(page.locator("h1")).toHaveText(/tooltip/i);

    const after = await sidebar.evaluate((el) => {
      const scroller = el.closest("[class*='overflow-y']") ?? el;
      return scroller.scrollTop;
    });
    expect(after).toBe(scrolled);
  });

  test("document title updates on boosted navigation", async ({ page }) => {
    await page.goto("/components/button");
    await page.getByRole("link", { name: "Theming", exact: true }).first().click();
    await expect(page).toHaveURL(/\/docs\/theming$/);
    await expect(page).toHaveTitle(/Theming · gsxui/);
  });

  test("theme class survives boosted navigation", async ({ page }) => {
    await page.goto("/components/button");
    await page.locator("[data-site-theme-toggle]").click();
    const dark = await page.evaluate(() => document.documentElement.classList.contains("dark"));
    await page.getByRole("link", { name: "Theming", exact: true }).first().click();
    await expect(page).toHaveURL(/\/docs\/theming$/);
    expect(await page.evaluate(() => document.documentElement.classList.contains("dark"))).toBe(dark);
  });
});
