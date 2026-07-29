import { expect, test } from "../support/fixtures";

test("compact docs navigation keeps docs and component links reachable below lg", async ({
  page,
}) => {
  await page.setViewportSize({ width: 900, height: 900 });
  await page.goto("/docs/theming");

  const headerPadding = await page
    .locator("header > div")
    .evaluate((element) => getComputedStyle(element).paddingLeft);
  const footerPadding = await page
    .locator("[data-site-footer] > div")
    .evaluate((element) => getComputedStyle(element).paddingLeft);
  expect(headerPadding).toBe("24px");
  expect(footerPadding).toBe(headerPadding);

  await expect(page.locator("[data-site-docs-sidebar]")).toBeHidden();
  const compactNav = page.locator("[data-site-docs-mobile-nav]");
  await expect(compactNav).toBeVisible();

  const trigger = compactNav.getByRole("button", {
    name: "Open documentation navigation",
  });
  await trigger.focus();
  await trigger.press("Enter");
  await expect(trigger).toHaveAttribute("aria-expanded", "true");

  const gettingStarted = compactNav.locator(
    'a[href="/docs/getting-started"]',
  );
  const button = compactNav.locator('a[href="/components/button"]');
  await expect(gettingStarted).toBeVisible();
  await expect(gettingStarted).toBeFocused();
  await expect(button).toBeVisible();

  await gettingStarted.click();
  await expect(page).toHaveURL(/\/docs\/getting-started$/);

  const nextCompactNav = page.locator("[data-site-docs-mobile-nav]");
  await nextCompactNav
    .getByRole("button", { name: "Open documentation navigation" })
    .press("Enter");
  await nextCompactNav.locator('a[href="/components/button"]').click();
  await expect(page).toHaveURL(/\/components\/button$/);
});
