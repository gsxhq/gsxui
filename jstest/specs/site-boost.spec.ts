import { expect, test } from "@playwright/test";

test.describe("boosted site navigation", () => {
  test("htmx is loaded on site pages", async ({ page }) => {
    await page.goto("/components/button");
    await expect
      .poll(async () => page.evaluate(() => typeof (window as any).htmx))
      .toBe("object");
  });
});
