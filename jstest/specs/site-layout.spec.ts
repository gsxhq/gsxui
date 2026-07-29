import { expect, test } from "../support/fixtures";

test("documentation rails respond around the fixed 640px article", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  await page.goto("/docs/theming");

  await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  await expect(page.locator("[data-site-docs-toc]")).toBeVisible();
  expect(
    await page
      .locator("[data-site-docs-article]")
      .evaluate((element) => element.getBoundingClientRect().width),
  ).toBe(640);

  await page.setViewportSize({ width: 1100, height: 900 });
  await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  await expect(page.locator("[data-site-docs-toc]")).toBeHidden();

  await page.setViewportSize({ width: 900, height: 900 });
  await expect(page.locator("[data-site-docs-sidebar]")).toBeHidden();
  await expect(page.locator("[data-site-docs-toc]")).toBeHidden();

  await page.setViewportSize({ width: 1280, height: 900 });
  await page.goto("/components");
  await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  await expect(page.locator("[data-site-docs-toc]")).toHaveCount(0);
});

test("component table of contents follows existing hashes and observed headings", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  await page.goto("/components/button#example-variants");

  const basicLink = page.locator(
    '[data-site-toc-link][href="#example-basic"]',
  );
  const variantsLink = page.locator(
    '[data-site-toc-link][href="#example-variants"]',
  );
  await expect(page.locator("#example-variants")).toHaveText("Variants");
  await expect(variantsLink).toHaveAttribute("data-active", "");
  await expect(variantsLink).toHaveAttribute("aria-current", "location");
  await expect(basicLink).not.toHaveAttribute("data-active", "");
  await expect(basicLink).not.toHaveAttribute("aria-current", "location");

  await basicLink.click();
  await expect(page).toHaveURL(/#example-basic$/);
  await expect(basicLink).toHaveAttribute("data-active", "");
  await expect(basicLink).toHaveAttribute("aria-current", "location");
  await expect(variantsLink).not.toHaveAttribute("data-active", "");

  await page.locator("#example-variants").evaluate((heading) => {
    window.scrollTo({
      top: heading.getBoundingClientRect().top + window.scrollY - 100,
    });
  });
  await expect(variantsLink).toHaveAttribute("data-active", "");
  await expect(variantsLink).toHaveAttribute("aria-current", "location");
  await expect(basicLink).not.toHaveAttribute("data-active", "");
});

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
  await expect(page.locator("[data-site-docs-toc]")).toBeHidden();
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
