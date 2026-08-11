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

  test("select reflects server-checked item after boosted navigation", async ({ page }) => {
    await page.goto("/components/button");
    // The sidebar link's accessible name also matches a hidden mobile-popover
    // duplicate that precedes the visible `aside` link in DOM order; `.last()`
    // picks the visible one (verified against the rendered DOM).
    await page.getByRole("link", { name: "select", exact: true }).last().click();
    await expect(page).toHaveURL(/\/components\/select$/);
    await expect(
      page.locator("[data-gsxui-slot-select-trigger]").first(),
    ).toContainText("Apple");
  });

  test("select reflects server-checked item after same-page boosted navigation", async ({ page }) => {
    await page.goto("/components/select");
    await expect(page.locator("[data-gsxui-slot-select-trigger]").first()).toContainText("Apple");
    await page.getByRole("link", { name: "select", exact: true }).last().click();
    await expect(page.locator("[data-gsxui-slot-select-trigger]").first()).toContainText("Apple");
  });

  test("theme editor initializes after boosted navigation", async ({ page }) => {
    await page.goto("/components/button");
    await page.getByRole("link", { name: "Theme", exact: true }).first().click();
    await expect(page).toHaveURL(/\/theme$/);
    // Both signals stay dead when the editor module's one-shot page scan ran
    // on the previous page: the preview handshake reports Live, and a mode
    // tab click re-renders the status line.
    await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
    await page.locator('[data-theme-mode-tab="dark"]').click();
    await expect(page.locator("[data-theme-status]")).toHaveText("Nova · dark");

    // Round trip: leaving and re-entering must rebuild the editor for the
    // fresh DOM (dispose the stale closure, wire the new one). goBack, not a
    // header link — the harness serves /theme as the shell-less ThemeEditor
    // page, which has no site navigation of its own.
    await page.goBack();
    await expect(page).toHaveURL(/\/components\/button$/);
    await page.getByRole("link", { name: "Theme", exact: true }).first().click();
    await expect(page).toHaveURL(/\/theme$/);
    await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
    await page.locator('[data-theme-mode-tab="dark"]').click();
    await expect(page.locator("[data-theme-status]")).toHaveText("Nova · dark");
  });
});
