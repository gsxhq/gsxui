import { expect, test } from "../support/fixtures";

async function openFixture(
  page: import("@playwright/test").Page,
  name: string,
  viewport = { width: 1024, height: 768 },
) {
  await page.setViewportSize(viewport);
  const response = await page.goto(`/f/sidebar-contract?case=${name}`);
  expect(response?.status(), `${name} fixture response`).toBe(200);
}

const slot = (name: string) => `[data-gsxui-slot~="${name}"]`;

test("desktop offcanvas state drives gap, fixed edge, rail cursor, event, and cookie ownership", async ({
  page,
}) => {
  await openFixture(page, "offcanvas-left");
  const wrapper = page.locator(slot("sidebar-wrapper"));
  const desktop = page.locator(slot("sidebar-desktop"));
  const gap = page.locator(slot("sidebar-gap"));
  const container = page.locator(slot("sidebar-container"));
  const rail = desktop.locator(slot("sidebar-rail"));

  expect(
    await wrapper.evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        state: element.getAttribute("data-state"),
        width: css.getPropertyValue("--sidebar-width").trim(),
        iconWidth: css.getPropertyValue("--sidebar-width-icon").trim(),
      };
    }),
  ).toEqual({ state: "expanded", width: "16rem", iconWidth: "3rem" });
  expect(await gap.evaluate((element) => getComputedStyle(element).width)).toBe("256px");
  expect(
    await container.evaluate((element) => {
      const css = getComputedStyle(element);
      return { left: css.left, width: css.width, position: css.position };
    }),
  ).toEqual({ left: "0px", width: "256px", position: "fixed" });
  expect(await rail.evaluate((element) => getComputedStyle(element).cursor)).toBe("w-resize");

  await page.evaluate(() => {
    document.cookie = "sidebar_owner=caller; path=/";
    const wrapper = document.querySelector('[data-gsxui-slot~="sidebar-wrapper"]')!;
    (window as typeof window & { sidebarChanges?: boolean[] }).sidebarChanges = [];
    wrapper.addEventListener("gsxui:change", (event) => {
      (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges.push(
        Boolean((event as CustomEvent).detail.open),
      );
    });
  });
  await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  await page.evaluate(() => {
    for (const animation of document.getAnimations()) animation.finish();
  });

  await expect(wrapper).toHaveAttribute("data-state", "collapsed");
  await expect(desktop).toHaveAttribute("data-state", "collapsed");
  await expect(desktop).toHaveAttribute("data-collapsible", "offcanvas");
  expect(await gap.evaluate((element) => getComputedStyle(element).width)).toBe("0px");
  expect(await container.evaluate((element) => getComputedStyle(element).left)).toBe("-256px");
  expect(await rail.evaluate((element) => getComputedStyle(element).cursor)).toBe("e-resize");
  expect(
    await page.evaluate(
      () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
    ),
  ).toEqual([false]);
  expect(await page.evaluate(() => document.cookie)).toContain("sidebar_owner=caller");

  await rail.click({ force: true });
  await expect(wrapper).toHaveAttribute("data-state", "expanded");
  await expect(desktop).toHaveAttribute("data-collapsible", "");
  expect(
    await page.evaluate(
      () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
    ),
  ).toEqual([false, true]);

  await page.getByRole("button", { name: "Custom Sidebar Trigger" }).click();
  await expect(wrapper).toHaveAttribute("data-state", "collapsed");
  expect(
    await page.evaluate(
      () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
    ),
  ).toEqual([false, true, false]);
});

test("right offcanvas and icon/floating/inset geometry follow explicit axes", async ({
  page,
}) => {
  await openFixture(page, "offcanvas-right");
  const right = page.locator(slot("sidebar-container"));
  expect(
    await right.evaluate((element) => {
      const css = getComputedStyle(element);
      return { right: css.right, borderLeft: css.borderLeftWidth };
    }),
  ).toEqual({ right: "-256px", borderLeft: "1px" });
  expect(
    await page
      .locator(slot("sidebar-desktop"))
      .locator(slot("sidebar-rail"))
      .evaluate((element) => getComputedStyle(element).cursor),
  ).toBe("w-resize");

  await openFixture(page, "icon-sidebar");
  expect(await page.locator(slot("sidebar-gap")).evaluate((el) => getComputedStyle(el).width)).toBe(
    "48px",
  );
  expect(
    await page.locator(slot("sidebar-container")).evaluate((el) => getComputedStyle(el).width),
  ).toBe("48px");

  await openFixture(page, "icon-floating");
  expect(await page.locator(slot("sidebar-gap")).evaluate((el) => getComputedStyle(el).width)).toBe(
    "64px",
  );
  expect(
    await page.locator(slot("sidebar-container")).evaluate((element) => {
      const css = getComputedStyle(element);
      return { width: css.width, padding: css.padding };
    }),
  ).toEqual({ width: "66px", padding: "8px" });
  const floatingChrome = await page.locator(slot("sidebar-inner")).evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        radius: css.borderRadius,
        border: css.borderWidth,
        shadow: css.boxShadow,
      };
    });
  expect(floatingChrome.radius).toBe("10px");
  expect(floatingChrome.border).toBe("1px");
  expect(floatingChrome.shadow).not.toBe("none");

  await openFixture(page, "inset-collapsed");
  expect(
    await page.locator(slot("sidebar-inset")).evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        marginTop: css.marginTop,
        marginLeft: css.marginLeft,
        radius: css.borderRadius,
      };
    }),
  ).toEqual({ marginTop: "8px", marginLeft: "8px", radius: "14px" });

  await openFixture(page, "none");
  await expect(page.locator(slot("sidebar"))).toHaveCount(1);
  await expect(page.locator(slot("sidebar-desktop"))).toHaveCount(0);
  await expect(page.locator(slot("sidebar-mobile-root"))).toHaveCount(0);
  expect(
    await page.locator(slot("sidebar")).evaluate((element) => {
      const css = getComputedStyle(element);
      return { display: css.display, width: css.width };
    }),
  ).toEqual({ display: "flex", width: "256px" });
});

test("icon collapse hides labelled branches and preserves menu button geometry", async ({
  page,
}) => {
  await openFixture(page, "menu-icon");
  const desktop = page.locator(slot("sidebar-desktop"));

  await expect(desktop.locator(slot("sidebar-group-label"))).toHaveCSS("opacity", "0");
  expect(
    await desktop
      .locator(slot("sidebar-group-label"))
      .evaluate((element) => getComputedStyle(element).marginTop),
  ).toBe("-32px");
  await expect(desktop.locator(slot("sidebar-group-action"))).toBeHidden();
  await expect(desktop.locator(slot("sidebar-menu-action")).first()).toBeHidden();
  await expect(desktop.locator(slot("sidebar-menu-badge")).first()).toBeHidden();
  await expect(desktop.locator(slot("sidebar-menu-sub"))).toBeHidden();

  const small = desktop.locator(`${slot("sidebar-menu-button")}[data-size="sm"]`);
  const large = desktop.locator(`${slot("sidebar-menu-button")}[data-size="lg"]`);
  expect(
    await small.evaluate((element) => {
      const css = getComputedStyle(element);
      return { width: css.width, height: css.height, padding: css.padding };
    }),
  ).toEqual({ width: "32px", height: "32px", padding: "8px" });
  expect(
    await large.evaluate((element) => {
      const css = getComputedStyle(element);
      return { width: css.width, height: css.height, padding: css.padding };
    }),
  ).toEqual({ width: "32px", height: "32px", padding: "0px" });

  const tooltipButton = desktop.getByRole("button", { name: "Tooltip item" });
  const tooltip = desktop.getByRole("tooltip");
  await expect(tooltip).toBeHidden();
  await tooltipButton.hover();
  await expect(tooltip).toBeVisible();
});

test("menu item relations reserve action space and align action/badge by button size", async ({
  page,
}) => {
  await openFixture(page, "menu-expanded");
  const desktop = page.locator(slot("sidebar-desktop"));
  const actionItem = desktop.locator('[data-sidebar-contract="action-item"]');
  const button = actionItem.locator(slot("sidebar-menu-button"));
  const action = actionItem.locator(slot("sidebar-menu-action"));
  const badge = actionItem.locator(slot("sidebar-menu-badge"));

  expect(await button.evaluate((element) => getComputedStyle(element).paddingRight)).toBe("32px");
  expect(await action.evaluate((element) => getComputedStyle(element).top)).toBe("4px");
  expect(await badge.evaluate((element) => getComputedStyle(element).top)).toBe("4px");
  expect(await action.evaluate((element) => getComputedStyle(element).color)).toBe(
    "oklch(0.985 0 0)",
  );
  await expect(action).toHaveCSS("opacity", "0");
  await actionItem.hover();
  await expect(action).toHaveCSS("opacity", "1");

  await page.mouse.move(1000, 700);
  await button.focus();
  await expect(action).toHaveCSS("opacity", "1");

  const defaultAction = desktop.locator('[data-sidebar-contract="default-action"]');
  const largeBadge = desktop.locator('[data-sidebar-contract="large-badge"]');
  expect(await defaultAction.evaluate((element) => getComputedStyle(element).top)).toBe("6px");
  expect(await largeBadge.evaluate((element) => getComputedStyle(element).top)).toBe("10px");

  const expandedTooltipButton = desktop.getByRole("button", { name: "Tooltip item" });
  const expandedTooltip = desktop.getByRole("tooltip");
  await expandedTooltipButton.hover();
  await page.waitForTimeout(350);
  await expect(expandedTooltip).toBeHidden();
});

test("all eight sidebar theme tokens feed live component declarations", async ({ page }) => {
  await openFixture(page, "tokens");
  const tree = page.locator(slot("sidebar-desktop"));

  const inner = tree.locator(slot("sidebar-inner"));
  const active = tree.locator('[data-sidebar-contract="active"]');
  const hover = tree.locator('[data-sidebar-contract="hover"]');
  const separator = tree.locator(slot("sidebar-separator"));
  const focus = tree.locator(slot("sidebar-group-label"));

  expect(await inner.evaluate((element) => getComputedStyle(element).backgroundColor)).toBe(
    "rgb(1, 2, 3)",
  );
  expect(await tree.evaluate((element) => getComputedStyle(element).color)).toBe(
    "rgb(4, 5, 6)",
  );
  expect(
    await active.evaluate((element) => {
      const css = getComputedStyle(element);
      return { background: css.backgroundColor, color: css.color };
    }),
  ).toEqual({ background: "rgb(7, 8, 9)", color: "rgb(10, 11, 12)" });
  await hover.hover();
  expect(
    await hover.evaluate((element) => {
      const css = getComputedStyle(element);
      return { background: css.backgroundColor, color: css.color };
    }),
  ).toEqual({ background: "rgb(13, 14, 15)", color: "rgb(16, 17, 18)" });
  expect(await separator.evaluate((element) => getComputedStyle(element).backgroundColor)).toBe(
    "rgb(19, 20, 21)",
  );
  await focus.focus();
  expect(await focus.evaluate((element) => getComputedStyle(element).boxShadow)).toContain(
    "rgb(22, 23, 24)",
  );
});

test("caller utility layer wins over sidebar and composed primitive defaults", async ({
  page,
}) => {
  await openFixture(page, "caller");
  const desktop = page.locator(slot("sidebar-desktop"));
  expect(
    await page.locator(slot("sidebar-container")).evaluate((element) => getComputedStyle(element).width),
  ).toBe("320px");
  expect(
    await desktop.locator(slot("sidebar-menu-button")).evaluate((element) => {
      const css = getComputedStyle(element);
      return { height: css.height, radius: css.borderRadius };
    }),
  ).toEqual({ height: "64px", radius: "0px" });
  expect(
    await desktop.locator(slot("sidebar-input")).evaluate((element) => getComputedStyle(element).height),
  ).toBe("48px");
});

test("mobile trigger opens the native Sheet tree and desktop resize keeps one visible tree", async ({
  page,
}) => {
  await openFixture(page, "mobile", { width: 390, height: 844 });
  const desktop = page.locator(slot("sidebar-desktop"));
  const mobileRoot = page.locator(slot("sidebar-mobile-root"));
  const dialog = page.locator("dialog[data-gsxui-sidebar-mobile-dialog]");

  await expect(desktop).toBeHidden();
  await expect(dialog).toBeHidden();
  await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  await expect(dialog).toBeVisible();
  await expect(dialog).toHaveAttribute("open", "");
  await expect(dialog).toHaveAttribute("data-state", "open");
  expect(await dialog.evaluate((element) => getComputedStyle(element).width)).toBe("288px");

  await page.setViewportSize({ width: 1024, height: 768 });
  await expect(mobileRoot).toBeHidden();
  await expect(desktop).toBeVisible();
  await expect(dialog).toBeHidden();
});

test("mobile rail stays non-interactive through the sm breakpoint and only the md tree operates", async ({
  page,
}) => {
  for (const width of [639, 640, 700, 767]) {
    await openFixture(page, "mobile", { width, height: 844 });
    const desktopContainer = page.locator(slot("sidebar-container"));
    const dialog = page.locator("dialog[data-gsxui-sidebar-mobile-dialog]");
    const mobileRail = dialog.locator(slot("sidebar-rail"));

    await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
    await expect(dialog, `mobile dialog at ${width}px`).toBeVisible();
    await expect(desktopContainer, `desktop container at ${width}px`).toBeHidden();
    await expect(mobileRail, `mobile rail at ${width}px`).toBeHidden();
    expect(
      await mobileRail.evaluate((rail) => {
        const rect = rail.getBoundingClientRect();
        const hit = document.elementFromPoint(0, window.innerHeight / 2);
        return {
          display: getComputedStyle(rail).display,
          width: rect.width,
          height: rect.height,
          ownsViewportEdge: hit === rail || rail.contains(hit),
        };
      }),
    ).toEqual({ display: "none", width: 0, height: 0, ownsViewportEdge: false });
  }

  await openFixture(page, "mobile", { width: 768, height: 844 });
  const wrapper = page.locator(slot("sidebar-wrapper"));
  const desktopContainer = page.locator(slot("sidebar-container"));
  const desktopRail = page.locator(slot("sidebar-desktop")).locator(slot("sidebar-rail"));
  const dialog = page.locator("dialog[data-gsxui-sidebar-mobile-dialog]");
  const mobileRail = dialog.locator(slot("sidebar-rail"));

  await expect(dialog).toBeHidden();
  await expect(mobileRail).toBeHidden();
  await expect(desktopContainer).toBeVisible();
  await expect(desktopRail).toBeVisible();
  await desktopRail.click();
  await expect(wrapper).toHaveAttribute("data-state", "collapsed");
});

test("custom style owns floating and inset collapsed density without foundation shrink", async ({
  page,
}) => {
  const customDensity = `
    @layer components {
      :where(
          [data-gsxui-slot~="sidebar-desktop"][data-collapsible="icon"]:is(
              [data-variant="floating"],
              [data-variant="inset"]
            )
        ) {
        --gsxui-sidebar-collapsed-gap-width: calc(var(--sidebar-width-icon) + 1.5rem);
        --gsxui-sidebar-collapsed-container-width: calc(
          var(--sidebar-width-icon) + 1.5rem + 8px
        );
      }
      :where(
          [data-gsxui-slot~="sidebar-desktop"]:is(
              [data-variant="floating"],
              [data-variant="inset"]
            )
        )
        > :where([data-gsxui-slot~="sidebar-container"]) {
        padding: 12px;
      }
      :where(
          [data-gsxui-slot~="sidebar-desktop"]:is(
              [data-variant="floating"],
              [data-variant="inset"]
            )
        )
        :where([data-gsxui-slot~="sidebar-inner"]) {
        border: 4px solid var(--sidebar-border);
      }
    }
  `;

  for (const fixture of ["icon-floating", "inset-collapsed"]) {
    await openFixture(page, fixture);
    await page.addStyleTag({ content: customDensity });
    await page.evaluate(() => {
      for (const animation of document.getAnimations()) animation.finish();
    });

    const geometry = await page.locator(slot("sidebar-desktop")).evaluate((desktop) => {
      const gap = desktop.querySelector<HTMLElement>('[data-gsxui-slot~="sidebar-gap"]')!;
      const container = desktop.querySelector<HTMLElement>(
        '[data-gsxui-slot~="sidebar-container"]',
      )!;
      const inner = desktop.querySelector<HTMLElement>('[data-gsxui-slot~="sidebar-inner"]')!;
      const gapCSS = getComputedStyle(gap);
      const containerCSS = getComputedStyle(container);
      const innerCSS = getComputedStyle(inner);
      return {
        gapWidth: gap.getBoundingClientRect().width,
        containerWidth: container.getBoundingClientRect().width,
        containerPadding: containerCSS.padding,
        innerBorder: innerCSS.borderLeftWidth,
        innerBorderBoxWidth: inner.getBoundingClientRect().width,
        innerContentWidth: inner.clientWidth,
        gapVariable: gapCSS.width,
      };
    });

    expect(geometry, fixture).toEqual({
      gapWidth: 72,
      containerWidth: 80,
      containerPadding: "12px",
      innerBorder: "4px",
      innerBorderBoxWidth: 56,
      innerContentWidth: 48,
      gapVariable: "72px",
    });
  }
});

test("keyboard shortcut uses the first wrapper, ignores typing targets and repeats", async ({
  page,
}) => {
  await openFixture(page, "keyboard");
  const first = page.locator(slot("sidebar-wrapper")).first();
  const second = page.locator(slot("sidebar-wrapper")).nth(1);
  await expect(first).toHaveAttribute("data-state", "expanded");
  await expect(second).toHaveAttribute("data-state", "expanded");

  await page.keyboard.press("Control+b");
  await expect(first).toHaveAttribute("data-state", "collapsed");
  await expect(second).toHaveAttribute("data-state", "expanded");

  await page.locator("input").focus();
  await page.keyboard.press("Control+b");
  await expect(first).toHaveAttribute("data-state", "collapsed");

  await page.evaluate(() => {
    document.dispatchEvent(
      new KeyboardEvent("keydown", {
        key: "b",
        ctrlKey: true,
        repeat: true,
        bubbles: true,
      }),
    );
  });
  await expect(first).toHaveAttribute("data-state", "collapsed");
});
