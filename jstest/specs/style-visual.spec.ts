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
  "dropdown-menu/basic",
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
  "toaster/types",
  "tabs/basic",
] as const;

const mobileRoutes = [
  "button/variants",
  "dialog/basic",
  "sidebar/variants",
  "calendar/basic",
] as const;

// Two-tier gate. Baselines are generated locally on macOS, where re-renders
// are pixel-exact (measured 0.0000 against committed baselines), so the
// tight 1% ratio does its real screening at authoring time. CI renders the
// same font bytes through FreeType instead of CoreText, which leaves a
// residual 1-2% glyph-AA noise floor on text-dense cards even with
// --font-render-hinting=none (playwright.config.ts) — measured across three
// CI runs where successive flag combinations only reshuffled WHICH borderline
// cards sat just past 1% (buttons/mobile -> feedback/calendar ->
// notification-prefs/navigation, all at exactly 0.02). CI therefore runs as
// a coarser backstop above that noise floor; it still fails on structural
// breakage, which measures well above 3%.
const screenshotOptions = {
  animations: "disabled" as const,
  caret: "hide" as const,
  maxDiffPixelRatio: process.env.CI ? 0.03 : 0.01,
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

  if (route === "dropdown-menu/basic") {
    await page.getByRole("button", { name: "Options" }).click();
    await expect(page.locator("[data-gsxui-slot-dropdown-menu-content]")).toBeVisible();
  }

  if (route === "toaster/types") {
    const triggers = page.locator("button[data-gsxui-toast]");
    await expect(triggers).toHaveCount(5);
    for (let i = 0; i < 5; i++) {
      await triggers.nth(i).click();
    }
    const toasts = page.locator("#gsxui-toaster > [data-gsxui-slot-toast]");
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

// Before the 8-style port, Button/Badge's destructive variant was a SOLID
// bg-destructive treatment (dark:bg-destructive/60 on hover-darkened text),
// which is what the harness's "dark-destructive"/"destructive-hover"
// reference elements (jstest/harness/style_contract.gsx) originally
// matched. The new base+8-styles design (apps/v4/registry/styles/
// style-*.css MARK: Button / Badge, verified 2026-08-11) redesigned
// destructive into a soft tint: bg-destructive/10 (light) /
// dark:bg-destructive/20 (rest) and dark:hover:bg-destructive/30 — the
// same value on both Button and Badge, faithfully carried by this port.
// The reference elements' alpha was updated to /20 and /30 to match.
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

// Padding values below are Toggle's own real per-size geometry
// (registry/styles/nova/toggle.css), not the flat px-3 every size used to
// share pre-port. They also depend on ToggleGroup's root actually carrying
// `group/toggle-group` — the 8-style port ported upstream's
// group-data-[spacing=…]/toggle-group: selectors verbatim onto
// ToggleGroupItem's recipe but left the marker CLASS that makes them match
// anything off ToggleGroup's own root (registry/canonical/toggle-group.gsx),
// a real bug this test caught: every group-data-…/toggle-group: rule in
// every style's toggle-group.css was dead code, silently, and joined items
// fell back to whatever unconditional padding happened to be next in the
// cascade instead of the spacing=0 joined-pill treatment upstream intends.
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
    paddingLeft: "8px",
    paddingRight: "8px",
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
    paddingLeft: "8px",
    paddingRight: "8px",
    borderLeftWidth: "0px",
    borderTopLeftRadius: "0px",
  });
  expect(iconItem.x).toBeCloseTo(first.x + first.width, 5);
  expect(last.x).toBeCloseTo(iconItem.x + iconItem.width, 5);

  expect(
    await metrics('[data-style-contract="toggle-group-default"]'),
  ).toMatchObject({
    height: "32px",
    paddingLeft: "10px",
    paddingRight: "10px",
  });
  expect(
    await metrics('[data-style-contract="toggle-group-large"]'),
  ).toMatchObject({
    height: "36px",
    paddingLeft: "10px",
    paddingRight: "10px",
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
      const title = element.querySelector("[data-gsxui-slot-dialog-title]");
      const description = element.querySelector("[data-gsxui-slot-dialog-description]");
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
    // Before the 8-style port, Dialog's ::backdrop referenced the shared
    // `--overlay` theme token directly (`backdrop:bg-[var(--overlay)]`),
    // so comparing against the "overlay" reference element was a same-token
    // check. Upstream's new base+8-styles design (apps/v4/registry/styles/
    // style-*.css MARK: Dialog / Alert Dialog, verified 2026-08-11) instead
    // hardcodes a PER-STYLE backdrop darkness (nova/vega/lyra: black/10,
    // maia/mira: black/80, luma/rhea: black/30, sera: black/20) with no
    // shared token at all — a genuine, deliberate design change this port
    // faithfully carries. nova's own black/10 happens to render the same
    // visual darkness as the old shared --overlay token, just through
    // Tailwind's built-in color scale (oklab) rather than --overlay's own
    // oklch literal, hence the differing (but visually identical)
    // getComputedStyle serialization.
    backdrop: "oklab(0 0 0 / 0.1)",
  });
  const relationships = await dialog.evaluate((element) => ({
    labelledBy: element.getAttribute("aria-labelledby"),
    titleID: element.querySelector("[data-gsxui-slot-dialog-title]")?.id,
    describedBy: element.getAttribute("aria-describedby"),
    descriptionID: element.querySelector("[data-gsxui-slot-dialog-description]")?.id,
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

// Per-style visual gate (style porter plan, Task 7; extended for the
// theme-creator-parity plan's Task 4c two-page split). Each style's
// site/stylepreview/<style> gallery composition is the same authored
// document (site/stylepreview/gallery.gsx.src) recompiled with that
// style's ported recipes, so the cards on each page always appear in the
// same DOM order — indexing by position within its own
// [data-theme-preview-page] wrapper, rather than by a per-card selector,
// keeps this gate from needing any change to the gallery source itself.
// This lands deliberately before the bulk port (Task 8): it establishes a
// reviewable baseline against nova before the other six styles arrive.
const galleryStyles = ["vega", "nova", "maia", "lyra", "mira", "luma", "sera", "rhea"] as const;

// componentsPageCards mirrors site/stylepreview/gallery_test.go's
// galleryPageCards["components"] slugged for snapshot filenames — keep
// both in sync with galleryComponentsPage's actual card order.
const componentsPageCards = [
  "buttons",
  "login",
  "settings",
  "calendar",
  "menus",
  "feedback",
  "team",
  "table",
  "chart",
  "tabs",
  "navigation",
  "controls",
  "overlays",
  "empty",
  "media",
  "sidebar",
] as const;

// productPageCards mirrors gallery_test.go's galleryPageCards["product"],
// slugged for snapshot filenames — keep both in sync with
// galleryProductPage's actual card order.
const productPageCards = [
  "pricing",
  "stats",
  "chat",
  "roles",
  "booking",
  "activity",
  "uploads",
  "notifications",
  "search",
  "profile",
  "onboarding",
  "projects-empty",
  "verify",
  "notification-prefs",
  "billing",
] as const;

test.describe("per-style gallery cards", () => {
  for (const style of galleryStyles) {
    test(`${style} gallery cards keep their per-style visual contract`, async ({ page }) => {
      const response = await page.goto(`/stylepreview/${style}`);
      expect(response?.status(), `${style} stylepreview response`).toBe(200);
      // The /stylepreview/<style> harness route has no page-tab JS wiring
      // (it renders Gallery() standalone, the same way it has no style-tab
      // JS either — see jstest/harness/stylepreview.go) — the product page
      // starts `hidden` per gallery.gsx.src's own JS-less default. Reveal
      // both pages so each page's own card loop below can screenshot
      // every card in one navigation.
      await page.evaluate(() => {
        for (const section of document.querySelectorAll<HTMLElement>("[data-theme-preview-page]")) {
          section.hidden = false;
        }
      });
      await page.evaluate(async () => {
        await document.fonts.ready;
        await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
      });

      const componentsCards = page.locator('[data-theme-preview-page="components"] > *');
      await expect(componentsCards).toHaveCount(componentsPageCards.length);
      for (const [index, card] of componentsPageCards.entries()) {
        await expect(componentsCards.nth(index)).toHaveScreenshot(`${style}-${card}.png`, screenshotOptions);
      }

      const productCards = page.locator('[data-theme-preview-page="product"] > *');
      await expect(productCards).toHaveCount(productPageCards.length);
      for (const [index, card] of productPageCards.entries()) {
        await expect(productCards.nth(index)).toHaveScreenshot(`${style}-${card}.png`, screenshotOptions);
      }
    });
  }
});

// The three tests below pin presentation that a component's own CSS supplies
// ON TOP OF the Button it renders through. ui.Button ships concrete Tailwind
// utilities compiled from its style recipe, and those land in @layer
// utilities; the cascade compares whole layers before specificity, so a rule
// in @layer components can never override them. Every property asserted here
// silently reverted to Button's own value when the Button switched from a
// recipe token to concrete utilities, and the 313-test suite did not notice.
// If one of these fails, check which layer the rule lives in before touching
// the expectation.

test("Carousel arrows stay circular over Button's own radius", async ({
  page,
}) => {
  const response = await page.goto("/x/carousel/basic");
  expect(response?.status(), "carousel/basic fixture response").toBe(200);

  for (const marker of [
    "[data-gsxui-slot-carousel-previous]",
    "[data-gsxui-slot-carousel-next]",
  ]) {
    const arrow = page.locator(marker).first();
    const { radius, height } = await arrow.evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        radius: Number.parseFloat(css.borderTopLeftRadius),
        height: Number.parseFloat(css.height),
      };
    });
    // rounded-full, not Button's rounded-lg (10px): a pill radius is clamped
    // to half the box, so anything >= height/2 renders as a circle here.
    expect(
      radius,
      `${marker} must be circular, not Button's rounded-lg`,
    ).toBeGreaterThanOrEqual(height / 2);
  }
});

test("InputGroupButton keeps its own type scale and radius over Button's size ramp", async ({
  page,
}) => {
  const response = await page.goto("/x/input-group/basic");
  expect(response?.status(), "input-group/basic fixture response").toBe(200);

  const button = page.locator("[data-gsxui-slot-input-group-button]").first();
  await expect(button).toHaveAttribute("data-size", "xs");
  expect(
    await button.evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        borderRadius: css.borderTopLeftRadius,
        fontSize: css.fontSize,
        paddingLeft: css.paddingLeft,
      };
    }),
  ).toEqual({
    // InputGroup's own xs ramp: rounded-[calc(var(--radius)-3px)] and text-sm,
    // NOT Button's xs ramp (rounded-[min(var(--radius-md),10px)] = 8px and
    // text-xs = 12px).
    borderRadius: "7px",
    fontSize: "14px",
    paddingLeft: "6px",
  });
});
