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
  "sidebar/variants?_preview=floating",
  "sidebar/variants?_preview=inset",
  "sidebar/variants?_preview=offcanvas",
  "sidebar/variants?_preview=icon",
  "sidebar/variants?_preview=none",
  "sidebar/variants?_preview=right-collapsed",
  "sidebar/variants?_preview=icon-collapsed",
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

function snapshotSlug(route: VisualRoute) {
  return route.replace("?_preview=", "-").replace("/", "-");
}

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

  if (route === "sidebar/variants" && (page.viewportSize()?.width ?? 0) < 768) {
    await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
    await expect(page.getByRole("dialog")).toBeVisible();
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

test("an active invalid InputOTP slot keeps destructive border and ring semantics", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const colors = async (selector: string) =>
    page.locator(selector).evaluate((element) => {
      const css = getComputedStyle(element);
      return { borderColor: css.borderColor, boxShadow: css.boxShadow };
    });

  const slot = '[data-style-contract="otp-active-invalid"]';
  expect(await colors(slot)).toEqual(
    await colors('[data-style-contract-reference="otp-invalid-light"]'),
  );

  await page.evaluate(() => document.documentElement.classList.add("dark"));
  await page.evaluate(() => {
    for (const animation of document.getAnimations()) {
      animation.finish();
    }
  });
  expect(await colors(slot)).toEqual(
    await colors('[data-style-contract-reference="otp-invalid-dark"]'),
  );
});

test("joined ToggleGroup items override composed Toggle sizing and borders", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const metrics = async (selector: string) =>
    page.locator(selector).evaluate((element) => {
      const css = getComputedStyle(element);
      const box = element.getBoundingClientRect();
      return {
        x: box.x,
        width: box.width,
        height: css.height,
        paddingLeft: css.paddingLeft,
        paddingRight: css.paddingRight,
        borderLeftWidth: css.borderLeftWidth,
        borderTopLeftRadius: css.borderTopLeftRadius,
        borderTopRightRadius: css.borderTopRightRadius,
      };
    });

  const first = await metrics('[data-style-contract="toggle-group-sm-first"]');
  const iconItem = await metrics('[data-style-contract="toggle-group-sm-icon"]');
  const last = await metrics('[data-style-contract="toggle-group-sm-last"]');
  expect(first).toMatchObject({
    height: "28px",
    paddingLeft: "12px",
    paddingRight: "12px",
    borderLeftWidth: "1px",
    borderTopRightRadius: "0px",
  });
  expect(iconItem).toMatchObject({
    height: "28px",
    paddingLeft: "6px",
    paddingRight: "6px",
    borderLeftWidth: "0px",
    borderTopLeftRadius: "0px",
    borderTopRightRadius: "0px",
  });
  expect(last).toMatchObject({
    height: "28px",
    paddingLeft: "12px",
    paddingRight: "12px",
    borderLeftWidth: "0px",
    borderTopLeftRadius: "0px",
  });
  expect(iconItem.x).toBeCloseTo(first.x + first.width, 5);
  expect(last.x).toBeCloseTo(iconItem.x + iconItem.width, 5);

  expect(
    await metrics('[data-style-contract="toggle-group-default"]'),
  ).toMatchObject({
    height: "32px",
    paddingLeft: "12px",
    paddingRight: "12px",
  });
  expect(
    await metrics('[data-style-contract="toggle-group-large"]'),
  ).toMatchObject({
    height: "36px",
    paddingLeft: "12px",
    paddingRight: "12px",
  });
  expect(
    await metrics('[data-style-contract="toggle-group-caller"]'),
  ).toMatchObject({
    height: "28px",
    paddingLeft: "32px",
    paddingRight: "32px",
  });
});

test("an outline FieldGroup applies the separator overlap", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  expect(
    await page
      .locator('[data-style-contract="field-outline-separator"]')
      .evaluate((element) => getComputedStyle(element).marginBottom),
  ).toBe("-8px");
});

test("Dialog keeps dedicated a11y hooks, semantic backdrop, and caller cascade", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  await page.getByRole("button", { name: "Open contract dialog" }).click();
  const dialog = page.locator('[data-style-contract="dialog-caller"]');
  await expect(dialog).toBeVisible();
  await expect(dialog).toHaveAttribute("data-state", "open");

  expect(
    await dialog.evaluate((element) => {
      const title = element.querySelector("[data-gsxui-dialog-title]");
      const description = element.querySelector("[data-gsxui-dialog-description]");
      const css = getComputedStyle(element);
      return {
        labelledBy: element.getAttribute("aria-labelledby"),
        titleID: title?.id,
        describedBy: element.getAttribute("aria-describedby"),
        descriptionID: description?.id,
        borderRadius: css.borderRadius,
        display: css.display,
        backdrop: getComputedStyle(element, "::backdrop").backgroundColor,
      };
    }),
  ).toEqual({
    labelledBy: expect.any(String),
    titleID: expect.any(String),
    describedBy: expect.any(String),
    descriptionID: expect.any(String),
    borderRadius: "0px",
    display: "grid",
    backdrop: await page
      .locator('[data-style-contract-reference="overlay"]')
      .evaluate((element) => getComputedStyle(element).backgroundColor),
  });
  const relationships = await dialog.evaluate((element) => ({
    labelledBy: element.getAttribute("aria-labelledby"),
    titleID: element.querySelector("[data-gsxui-dialog-title]")?.id,
    describedBy: element.getAttribute("aria-describedby"),
    descriptionID: element.querySelector("[data-gsxui-dialog-description]")?.id,
  }));
  expect(relationships.labelledBy).toBe(relationships.titleID);
  expect(relationships.describedBy).toBe(relationships.descriptionID);
});

test("bottom Drawer uses Dialog mechanics and content-side header alignment", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  await page.getByRole("button", { name: "Open contract drawer" }).click();
  const drawer = page.locator('[data-style-contract="drawer-bottom"]');
  await expect(drawer).toBeVisible();
  await expect(drawer).toHaveAttribute("data-state", "open");
  await page.evaluate(() => {
    for (const animation of document.getAnimations()) animation.finish();
  });

  expect(
    await drawer.evaluate((element) => {
      const rect = element.getBoundingClientRect();
      const header = element.querySelector('[data-style-contract="drawer-header"]');
      return {
        display: getComputedStyle(element).display,
        side: element.getAttribute("data-side"),
        left: rect.left,
        right: rect.right,
        bottom: rect.bottom,
        viewportWidth: window.innerWidth,
        viewportHeight: window.innerHeight,
        headerAlign: header ? getComputedStyle(header).textAlign : null,
      };
    }),
  ).toEqual({
    display: "flex",
    side: "bottom",
    left: 0,
    right: 1280,
    bottom: 900,
    viewportWidth: 1280,
    viewportHeight: 900,
    headerAlign: "center",
  });
});

test("Accordion caller padding overrides the inner content default", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const content = page.locator(
    '[data-style-contract="accordion-caller-content"]',
  );
  const inner = content.locator(
    ':scope > [data-gsxui-slot-accordion-content-inner]',
  );

  await expect(content).toHaveAttribute("id", "accordion-caller-content");
  await expect(content).not.toHaveClass(/\bpb-8\b/);
  await expect(inner).toHaveClass(/\bpb-8\b/);
  expect(
    await content.evaluate((element) => getComputedStyle(element).paddingBottom),
  ).toBe("0px");
  expect(
    await inner.evaluate((element) => getComputedStyle(element).paddingBottom),
  ).toBe("32px");
});

const directionalOverlayCases = [
  {
    family: "drawer",
    fixture: "default",
    side: "bottom",
    enterProperty: "--tw-enter-translate-y",
    enterValue: "100%",
  },
  {
    family: "drawer",
    fixture: "top",
    side: "top",
    enterProperty: "--tw-enter-translate-y",
    enterValue: "-100%",
  },
  {
    family: "drawer",
    fixture: "left",
    side: "left",
    enterProperty: "--tw-enter-translate-x",
    enterValue: "-100%",
  },
  {
    family: "drawer",
    fixture: "right",
    side: "right",
    enterProperty: "--tw-enter-translate-x",
    enterValue: "100%",
  },
  {
    family: "sheet",
    fixture: "default",
    side: "right",
    enterProperty: "--tw-enter-translate-x",
    enterValue: "100%",
  },
  {
    family: "sheet",
    fixture: "left",
    side: "left",
    enterProperty: "--tw-enter-translate-x",
    enterValue: "-100%",
  },
  {
    family: "sheet",
    fixture: "top",
    side: "top",
    enterProperty: "--tw-enter-translate-y",
    enterValue: "-100%",
  },
  {
    family: "sheet",
    fixture: "bottom",
    side: "bottom",
    enterProperty: "--tw-enter-translate-y",
    enterValue: "100%",
  },
] as const;

for (const overlay of directionalOverlayCases) {
  test(`${overlay.family} ${overlay.fixture} opens on ${overlay.side} and closes through its control`, async ({
    page,
  }) => {
    const response = await page.goto("/f/style-contract");
    expect(response?.status(), "style contract fixture response").toBe(200);

    const dialog = page.locator(
      `[data-style-contract="${overlay.family}-${overlay.fixture}"]`,
    );
    await expect(dialog).toBeHidden();
    expect(
      await dialog.evaluate((element: HTMLDialogElement) => ({
        open: element.open,
        state: element.dataset.state,
        display: getComputedStyle(element).display,
      })),
    ).toEqual({
      open: false,
      state: "closed",
      display: "none",
    });

    await page
      .getByRole("button", {
        name: `Open ${overlay.family} ${overlay.fixture}`,
      })
      .click();
    await expect(dialog).toBeVisible();
    await expect(dialog).toHaveAttribute("open", "");
    await expect(dialog).toHaveAttribute("data-state", "open");
    await dialog.evaluate((element) => {
      for (const animation of element.getAnimations()) animation.finish();
    });

    const geometry = await dialog.evaluate(
      (element, enterProperty) => {
        const rect = element.getBoundingClientRect();
        const css = getComputedStyle(element);
        return {
          side: element.getAttribute("data-side"),
          display: css.display,
          animationName: css.animationName,
          enterValue: css.getPropertyValue(enterProperty).trim(),
          left: rect.left,
          right: rect.right,
          top: rect.top,
          bottom: rect.bottom,
          width: rect.width,
          height: rect.height,
          viewportWidth: window.innerWidth,
          viewportHeight: window.innerHeight,
        };
      },
      overlay.enterProperty,
    );
    expect(geometry.side).toBe(overlay.side);
    expect(geometry.display).toBe("flex");
    expect(geometry.animationName).toContain("enter");
    expect(geometry.enterValue).toBe(overlay.enterValue);

    const subpixelEdgeTolerance = 0.5;
    const expectAtEdge = (actual: number, expected: number) =>
      expect(Math.abs(actual - expected)).toBeLessThanOrEqual(
        subpixelEdgeTolerance,
      );
    if (overlay.side === "top" || overlay.side === "bottom") {
      expectAtEdge(geometry.left, 0);
      expectAtEdge(geometry.right, geometry.viewportWidth);
      expectAtEdge(geometry.width, geometry.viewportWidth);
      expect(geometry.height).toBeGreaterThan(0);
      expect(geometry.height).toBeLessThan(geometry.viewportHeight);
      if (overlay.side === "top") {
        expectAtEdge(geometry.top, 0);
      } else {
        expectAtEdge(geometry.bottom, geometry.viewportHeight);
      }
    } else {
      expectAtEdge(geometry.top, 0);
      expectAtEdge(geometry.bottom, geometry.viewportHeight);
      expectAtEdge(geometry.height, geometry.viewportHeight);
      expectAtEdge(geometry.width, 384);
      if (overlay.side === "left") {
        expectAtEdge(geometry.left, 0);
      } else {
        expectAtEdge(geometry.right, geometry.viewportWidth);
      }
    }

    await page
      .getByRole("button", {
        name: `Close ${overlay.family} ${overlay.fixture}`,
      })
      .click();
    await expect(dialog).toHaveAttribute("data-state", "closed");
    const closing = await dialog.evaluate(
      (element: HTMLDialogElement, exitProperty) => {
        const css = getComputedStyle(element);
        return {
          open: element.open,
          animationName: css.animationName,
          exitValue: css.getPropertyValue(exitProperty).trim(),
        };
      },
      overlay.enterProperty.replace("--tw-enter-", "--tw-exit-"),
    );
    expect(closing.open).toBe(true);
    expect(closing.animationName).toContain("exit");
    expect(closing.exitValue).toBe(overlay.enterValue);
    await expect(dialog).toBeHidden();
    expect(
      await dialog.evaluate((element: HTMLDialogElement) => ({
        open: element.open,
        display: getComputedStyle(element).display,
      })),
    ).toEqual({
      open: false,
      display: "none",
    });
  });
}

test("Tooltip Kbd relationship and Accordion native open mechanics compute", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  await page.getByRole("button", { name: "Show contract tooltip" }).focus();
  const tooltip = page.locator('[data-style-contract="tooltip-kbd"]');
  await expect(tooltip).toBeVisible();
  expect(
    await tooltip.evaluate((element) => getComputedStyle(element).paddingRight),
  ).toBe("6px");

  const trigger = page.getByText("Contract accordion", { exact: true });
  await trigger.click();
  const details = trigger.locator("..");
  await expect(details).toHaveAttribute("open", "");
  await page.evaluate(() => {
    for (const animation of document.getAnimations()) animation.finish();
  });
  expect(
    await details.evaluate((element) => {
      const icon = element.querySelector(
        '[data-gsxui-slot-accordion-trigger-icon]',
      );
      return {
        detailsContentDisplay: getComputedStyle(element, "::details-content").display,
        iconRotate: icon ? getComputedStyle(icon).rotate : null,
      };
    }),
  ).toEqual({
    detailsContentDisplay: "grid",
    iconRotate: "180deg",
  });
});

for (const route of desktopRoutes) {
  test(`${route} keeps its desktop light and dark visual contract`, async ({ page }) => {
    for (const theme of ["light", "dark"] as const) {
      await prepareVisualRoute(page, route, theme);
      await expect(page.locator("body")).toHaveScreenshot(
        `desktop-${theme}-${snapshotSlug(route)}.png`,
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
        `mobile-light-${snapshotSlug(route)}.png`,
        screenshotOptions,
      );
    });
  }
});
