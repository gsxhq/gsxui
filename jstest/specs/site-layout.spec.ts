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

  const sidebarNav = page.locator("[data-site-docs-sidebar] nav");
  await sidebarNav.evaluate((element) => {
    element.scrollTop = element.scrollHeight;
  });
  const sidebarEnd = await page.evaluate(() => {
    const lastLink = document.querySelector<HTMLElement>(
      '[data-site-docs-sidebar] a[href="/components/tooltip"]',
    );
    if (!lastLink) {
      throw new Error("final documentation sidebar link is missing");
    }
    return {
      documentScrollY: window.scrollY,
      linkBottom: lastLink.getBoundingClientRect().bottom,
      viewportBottom: document.documentElement.clientHeight,
    };
  });
  expect(sidebarEnd.documentScrollY).toBe(0);
  expect(sidebarEnd.linkBottom).toBeLessThanOrEqual(sidebarEnd.viewportBottom);
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

test("theme workspace contains desktop overflow and preserves narrow document order", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  await page.goto("/site/theme");

  const stylePanel = page.locator("[data-theme-style-panel]");
  const previewPanel = page.locator("[data-theme-preview-panel]");
  const controlsPanel = page.locator("[data-theme-controls-panel]");

  await expect(stylePanel).toBeVisible();
  await expect(previewPanel).toBeVisible();
  await expect(controlsPanel).toBeVisible();
  await expect(page.locator("[data-theme-preview-frame]")).toHaveCount(1);

  const desktopGeometry = await page.evaluate(() => {
    const controls = document.querySelector<HTMLElement>(
      "[data-theme-controls-panel]",
    );
    const preview = document.querySelector<HTMLElement>(
      "[data-theme-preview-panel]",
    );
    if (!controls || !preview) {
      throw new Error("theme workspace panels are missing");
    }
    return {
      documentOverflow:
        document.documentElement.scrollHeight >
        document.documentElement.clientHeight,
      controlsWidth: controls.getBoundingClientRect().width,
      previewWidth: preview.getBoundingClientRect().width,
      controlsOverflow: getComputedStyle(controls).overflowY,
      controlsScrollHeight: controls.scrollHeight,
      controlsClientHeight: controls.clientHeight,
    };
  });
  expect(desktopGeometry.documentOverflow).toBe(false);
  expect(desktopGeometry.controlsWidth).toBeLessThan(
    desktopGeometry.previewWidth,
  );
  expect(desktopGeometry.controlsOverflow).toBe("auto");
  expect(desktopGeometry.controlsScrollHeight).toBeGreaterThan(
    desktopGeometry.controlsClientHeight,
  );

  await page.setViewportSize({ width: 900, height: 900 });

  const narrowGeometry = await page.evaluate(() => {
    const top = (selector: string) => {
      const element = document.querySelector<HTMLElement>(selector);
      if (!element) {
        throw new Error(`missing ${selector}`);
      }
      return element.getBoundingClientRect().top;
    };
    return {
      documentOverflow:
        document.documentElement.scrollHeight >
        document.documentElement.clientHeight,
      styleTop: top("[data-theme-style-panel]"),
      previewTop: top("[data-theme-preview-panel]"),
      controlsTop: top("[data-theme-controls-panel]"),
      iframeCount: document.querySelectorAll("[data-theme-preview-frame]")
        .length,
    };
  });
  expect(narrowGeometry.documentOverflow).toBe(true);
  expect(narrowGeometry.styleTop).toBeLessThan(narrowGeometry.previewTop);
  expect(narrowGeometry.previewTop).toBeLessThan(narrowGeometry.controlsTop);
  expect(narrowGeometry.iframeCount).toBe(1);
});

test("short desktop theme workspaces keep every control region reachable", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1280, height: 360 });
  await page.goto("/site/theme");

  const workspaceMain = page.locator(
    'body[data-site-layout="workspace"] > div > main',
  );
  const controls = page.locator("[data-theme-controls-panel]");
  const geometry = await page.evaluate(() => {
    const main = document.querySelector<HTMLElement>(
      'body[data-site-layout="workspace"] > div > main',
    );
    const controlsPanel = document.querySelector<HTMLElement>(
      "[data-theme-controls-panel]",
    );
    if (!main || !controlsPanel) {
      throw new Error("theme workspace scrolling regions are missing");
    }
    return {
      documentWidth: document.documentElement.scrollWidth,
      viewportWidth: document.documentElement.clientWidth,
      mainOverflowY: getComputedStyle(main).overflowY,
      mainScrollHeight: main.scrollHeight,
      mainClientHeight: main.clientHeight,
      controlsClientHeight: controlsPanel.clientHeight,
      controlsScrollHeight: controlsPanel.scrollHeight,
      controlsOverflowY: getComputedStyle(controlsPanel).overflowY,
      controlsScrollWidth: controlsPanel.scrollWidth,
      controlsClientWidth: controlsPanel.clientWidth,
    };
  });

  expect(geometry.documentWidth).toBeLessThanOrEqual(geometry.viewportWidth);
  expect(geometry.mainOverflowY).toBe("auto");
  expect(geometry.mainScrollHeight).toBeGreaterThan(geometry.mainClientHeight);
  expect(geometry.controlsClientHeight).toBeGreaterThan(0);
  expect(geometry.controlsOverflowY).toBe("auto");
  expect(geometry.controlsScrollHeight).toBeGreaterThan(
    geometry.controlsClientHeight,
  );
  expect(geometry.controlsScrollWidth).toBeLessThanOrEqual(
    geometry.controlsClientWidth,
  );

  for (const name of [
    "Mode and palette",
    "Preset JSON",
    "Theme CSS",
    "Share and install",
  ]) {
    const heading = controls.getByRole("heading", { name, exact: true });
    await heading.evaluate((element) =>
      element.scrollIntoView({ block: "nearest" }),
    );
    expect(
      await heading.evaluate((element) => {
        const rect = element.getBoundingClientRect();
        const scroller = element
          .closest("[data-theme-controls-panel]")
          ?.getBoundingClientRect();
        return (
          !!scroller &&
          rect.bottom > Math.max(0, scroller.top) &&
          rect.top < Math.min(innerHeight, scroller.bottom)
        );
      }),
      `${name} should be reachable in the short workspace`,
    ).toBe(true);
  }

  await expect(workspaceMain).toBeVisible();
  await expect(page.locator("[data-theme-preview-frame]")).toHaveCount(1);
});

test("production-composed theme workspace loads its controller and preview cleanly", async ({
  page,
}) => {
  const consoleErrors: string[] = [];
  const pageErrors: string[] = [];
  const failedRequests: string[] = [];
  const failedResponses: string[] = [];
  page.on("console", (message) => {
    if (message.type() === "error") consoleErrors.push(message.text());
  });
  page.on("pageerror", (error) => pageErrors.push(String(error)));
  page.on("requestfailed", (request) => {
    failedRequests.push(
      `${request.method()} ${request.url()} ${request.failure()?.errorText ?? ""}`,
    );
  });
  page.on("response", (response) => {
    if (response.status() >= 400) {
      failedResponses.push(`${response.status()} ${response.url()}`);
    }
  });

  await page.setViewportSize({ width: 1280, height: 900 });
  await page.goto("/site/theme");
  await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  await expect(page.locator("[data-theme-preview-retry]")).toBeHidden();

  const theme = page.locator('[data-theme-picker="theme"]');
  await theme.locator("[data-theme-picker-trigger]").click();
  await theme.getByRole("radio", { name: "Blue", exact: true }).click();
  await expect(
    theme.locator("[data-theme-selection-value]"),
  ).toHaveText("Blue");

  const catalog = JSON.parse(
    (await page.locator("[data-theme-schema]").textContent()) ?? "{}",
  );
  await expect
    .poll(() =>
      page
        .frameLocator("[data-theme-preview-frame]")
        .locator("html")
        .evaluate((element) =>
          element.style.getPropertyValue("--primary").trim(),
        ),
    )
    .toBe(catalog.palette.resolved.neutral.blue.light.primary);

  expect(consoleErrors).toEqual([]);
  expect(pageErrors).toEqual([]);
  expect(failedRequests).toEqual([]);
  expect(failedResponses).toEqual([]);
});
