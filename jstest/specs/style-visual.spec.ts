import { expect, test } from "../support/fixtures";

const desktopRoutes = [
  "accordion/basic",
  "alert/variants",
  "button/variants",
  "card/compound",
  "calendar/loadedrange",
  "checkbox/states",
  "combobox/basic",
  "dialog/basic",
  "dropdown/basic",
  "field/invalid",
  "navigation-menu/mega",
  "sidebar/variants",
  "sonner/types",
  "tabs/basic",
] as const;

const mobileRoutes = [
  "button/variants",
  "dialog/basic",
  "sidebar/variants",
  "calendar/basic",
] as const;

const screenshotOptions = {
  animations: "disabled" as const,
  caret: "hide" as const,
  maxDiffPixelRatio: 0.01,
};

type VisualRoute = (typeof desktopRoutes)[number] | (typeof mobileRoutes)[number];

async function prepareVisualRoute(
  page: import("@playwright/test").Page,
  route: VisualRoute,
  theme: "light" | "dark",
) {
  const response = await page.goto(`/x/${route}`);
  expect(response?.status(), `${route} fixture response`).toBe(200);

  await page.evaluate(async (isDark) => {
    document.documentElement.classList.toggle("dark", isDark);
    await document.fonts.ready;
    await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
  }, theme === "dark");

  if (route === "dialog/basic") {
    await page.getByRole("button", { name: "Delete account" }).click();
    await expect(page.getByRole("dialog")).toBeVisible();
  }

  if (route === "dropdown/basic") {
    await page.getByRole("button", { name: "Options" }).click();
    await expect(page.locator("[data-gsxui-dropdown-content]")).toBeVisible();
  }

  if (route === "sonner/types") {
    const triggers = page.locator("button[data-gsxui-toast]");
    await expect(triggers).toHaveCount(5);
    for (let i = 0; i < 5; i++) {
      await triggers.nth(i).click();
    }
    const toasts = page.locator("#gsxui-toaster > [data-gsxui-toast]");
    await expect(toasts).toHaveCount(5);
    await toasts.last().hover();
    await expect(toasts.nth(2)).toBeVisible();
  }
}

test("caller utilities override the Button defaults", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const override = page.getByRole("button", { name: "Caller override" });
  expect(
    await override.evaluate((el) => {
      const css = getComputedStyle(el);
      return {
        height: css.height,
        borderRadius: css.borderRadius,
        display: css.display,
      };
    }),
  ).toEqual({
    height: "48px",
    borderRadius: "0px",
    display: "inline-flex",
  });
});

test("dark primitive states use their dark semantic colors", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);
  await page.evaluate(() => document.documentElement.classList.add("dark"));
  const finishTransitions = () =>
    page.evaluate(() => {
      for (const animation of document.getAnimations()) {
        animation.finish();
      }
    });
  await finishTransitions();

  const computed = async (
    selector: string,
    property: "backgroundColor" | "boxShadow",
  ) =>
    page.locator(selector).evaluate(
      (element, name) => getComputedStyle(element)[name],
      property,
    );

  const darkDestructive = await computed(
    '[data-style-contract-reference="dark-destructive"]',
    "backgroundColor",
  );
  expect(
    await computed(
      '[data-style-contract="dark-button-destructive"]',
      "backgroundColor",
    ),
  ).toBe(darkDestructive);
  expect(
    await computed(
      '[data-style-contract="dark-badge-destructive"]',
      "backgroundColor",
    ),
  ).toBe(darkDestructive);
  const destructive = page.locator(
    '[data-style-contract="dark-button-destructive"]',
  );
  await destructive.hover();
  await finishTransitions();
  expect(
    await computed(
      '[data-style-contract="dark-button-destructive"]',
      "backgroundColor",
    ),
  ).toBe(
    await computed(
      '[data-style-contract-reference="destructive-hover"]',
      "backgroundColor",
    ),
  );

  const darkInvalidRing = await computed(
    '[data-style-contract-reference="dark-invalid-ring"]',
    "boxShadow",
  );
  const invalidButton = page.locator('[data-style-contract="dark-button-invalid"]');
  await invalidButton.focus();
  await finishTransitions();
  expect(
    await computed('[data-style-contract="dark-button-invalid"]', "boxShadow"),
  ).toBe(darkInvalidRing);
  const invalidBadge = page.locator('[data-style-contract="dark-badge-invalid"]');
  await invalidBadge.focus();
  await finishTransitions();
  expect(
    await computed('[data-style-contract="dark-badge-invalid"]', "boxShadow"),
  ).toBe(darkInvalidRing);

  const outline = page.locator('[data-style-contract="dark-button-outline"]');
  await outline.hover();
  await finishTransitions();
  expect(await computed('[data-style-contract="dark-button-outline"]', "backgroundColor")).toBe(
    await computed(
      '[data-style-contract-reference="dark-outline-hover"]',
      "backgroundColor",
    ),
  );

  const ghost = page.locator('[data-style-contract="dark-button-ghost"]');
  await ghost.hover();
  await finishTransitions();
  expect(await computed('[data-style-contract="dark-button-ghost"]', "backgroundColor")).toBe(
    await computed(
      '[data-style-contract-reference="dark-ghost-hover"]',
      "backgroundColor",
    ),
  );
});

test("Pagination edge padding overrides Button defaults and remains caller-overridable", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const padding = async (selector: string) =>
    page.locator(selector).evaluate((element) => {
      const css = getComputedStyle(element);
      return { left: css.paddingLeft, right: css.paddingRight };
    });

  expect(await padding('[data-style-contract="pagination-previous"]')).toEqual({
    left: "6px",
    right: "8px",
  });
  expect(await padding('[data-style-contract="pagination-next"]')).toEqual({
    left: "8px",
    right: "6px",
  });
  expect(await padding('[data-style-contract="pagination-previous-caller"]')).toEqual({
    left: "48px",
    right: "8px",
  });
});

for (const route of desktopRoutes) {
  test(`${route} keeps its desktop light and dark visual contract`, async ({ page }) => {
    for (const theme of ["light", "dark"] as const) {
      await prepareVisualRoute(page, route, theme);
      await expect(page.locator("body")).toHaveScreenshot(
        `desktop-${theme}-${route.replace("/", "-")}.png`,
        screenshotOptions,
      );
    }
  });
}

test.describe("mobile visual contracts", () => {
  test.use({ viewport: { width: 390, height: 844 } });

  for (const route of mobileRoutes) {
    test(`${route} keeps its mobile visual contract`, async ({ page }) => {
      await prepareVisualRoute(page, route, "light");
      await expect(page.locator("body")).toHaveScreenshot(
        `mobile-light-${route.replace("/", "-")}.png`,
        screenshotOptions,
      );
    });
  }
});
