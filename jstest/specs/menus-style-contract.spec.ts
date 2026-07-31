import { expect, test } from "../support/fixtures";

async function openStyleContract(page: import("@playwright/test").Page) {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);
}

test("menu destructive focus state wins without important specificity", async ({
  page,
}) => {
  await openStyleContract(page);

  const reference = page.locator(
    '[data-style-contract-reference="menu-destructive-focus"]',
  );
  const expected = await reference.evaluate((element) => {
    const css = getComputedStyle(element);
    return { backgroundColor: css.backgroundColor, color: css.color };
  });

  for (const family of ["dropdown", "context", "menubar"]) {
    const item = page.locator(`[data-style-contract="${family}-destructive"]`);
    await item.focus();
    expect(
      await item.evaluate((element) => {
        const css = getComputedStyle(element);
        return { backgroundColor: css.backgroundColor, color: css.color };
      }),
    ).toEqual(expected);
  }
});

test("checked menu and combobox indicators follow semantic owner state", async ({
  page,
}) => {
  await openStyleContract(page);

  const cases = [
    ["dropdown", "[data-gsxui-dropdown-checkbox-indicator]"],
    ["context", "[data-gsxui-contextmenu-radio-indicator]"],
    ["menubar", "[data-gsxui-menubar-checkbox-indicator]"],
    ["combobox", "[data-gsxui-combobox-item-indicator]"],
  ] as const;

  for (const [family, indicator] of cases) {
    const display = async (state: "checked" | "unchecked") =>
      page
        .locator(`[data-style-contract="${family}-${state}"] ${indicator}`)
        .evaluate((element) => getComputedStyle(element).display);
    expect(await display("checked")).toBe("flex");
    expect(await display("unchecked")).toBe("none");
  }
});

test("submenu uses the emitted right-side directional enter offset", async ({
  page,
}) => {
  await openStyleContract(page);

  const content = page.locator('[data-style-contract="submenu-right"]');
  const keyframes = await content.evaluate((element) => {
    (element as HTMLElement).showPopover();
    return element.getAnimations().flatMap((animation) => {
      const effect = animation.effect as KeyframeEffect | null;
      return effect?.getKeyframes().map((frame) => frame.translate) ?? [];
    });
  });
  expect(keyframes).toContain("-8px");
  await content.evaluate((element) => (element as HTMLElement).hidePopover());
});

test("CommandDialog sizing is supplied by CSS ancestry", async ({ page }) => {
  await openStyleContract(page);
  await page.keyboard.press(process.platform === "darwin" ? "Meta+K" : "Control+K");

  const dialog = page.locator('[data-style-contract="command-dialog"]');
  await expect(dialog).toBeVisible();
  expect(
    await dialog.evaluate((element) => {
      const command = element.querySelector("[data-gsxui-command]");
      const wrapper = element.querySelector("[data-gsxui-command-input-wrapper]");
      const input = element.querySelector("[data-gsxui-command-input]");
      const css = getComputedStyle(element);
      return {
        maxWidth: css.maxWidth,
        overflow: css.overflow,
        padding: css.padding,
        commandPadding: command ? getComputedStyle(command).padding : null,
        wrapperHeight: wrapper ? getComputedStyle(wrapper).height : null,
        inputHeight: input ? getComputedStyle(input).height : null,
      };
    }),
  ).toEqual({
    maxWidth: "512px",
    overflow: "hidden",
    padding: "0px",
    commandPadding: "4px",
    wrapperHeight: "48px",
    inputHeight: "48px",
  });
});

test("navigation viewport chrome and trigger rotation follow reflected state", async ({
  page,
}) => {
  await openStyleContract(page);

  const root = page.locator('[data-style-contract="navigation-root"]');
  await expect(root).toHaveAttribute("data-viewport", "false");
  await page.getByRole("button", { name: "Contract products" }).click();

  const trigger = page.locator('[data-style-contract="navigation-trigger"]');
  const content = page.locator('[data-style-contract="navigation-content"]');
  await expect(content).toBeVisible();
  await expect(trigger).toHaveAttribute("data-state", "open");
  await expect(content).toHaveAttribute("data-state", "open");
  await page.evaluate(() => {
    for (const animation of document.getAnimations()) animation.finish();
  });

  expect(
    await content.evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        borderWidth: css.borderWidth,
        borderRadius: css.borderRadius,
        overflow: css.overflow,
        position: css.position,
      };
    }),
  ).toEqual({
    borderWidth: "1px",
    borderRadius: "10px",
    overflow: "hidden",
    position: "fixed",
  });
  expect(
    await trigger
      .locator("[data-gsxui-navigation-menu-trigger-icon]")
      .evaluate((element) => getComputedStyle(element).rotate),
  ).toBe("180deg");
});
