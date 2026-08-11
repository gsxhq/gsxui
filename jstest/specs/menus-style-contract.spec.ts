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
    ["dropdown", "[data-gsxui-slot-dropdown-menu-checkbox-item-indicator]"],
    ["context", "[data-gsxui-slot-context-menu-radio-item-indicator]"],
    ["menubar", "[data-gsxui-slot-menubar-checkbox-item-indicator]"],
    ["combobox", "[data-gsxui-slot-combobox-item-indicator]"],
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

  // Before the 8-style port, DropdownMenuSubContent's enter transition was a
  // hand-authored `@starting-style` CSS-transition keyed on the `:popover-open`
  // pseudo-class (assets/css/styles/default/dropdown-menu.css's old
  // `starting:[&:popover-open]:-translate-x-2`), which used the standalone
  // `translate` property and needed no JS-set attribute — pure
  // showPopover() was enough to engage it. The new base+8-styles design
  // (apps/v4/registry/styles/style-*.css MARK: Dropdown Menu / Context Menu /
  // Menubar sections, verified 2026-08-11) uniformly replaced this with
  // Tailwind's `data-[state=open]:animate-in ... slide-in-from-left-2`
  // @keyframes mechanism, which (a) renders the offset through the `transform`
  // property instead of `translate`, and (b) is gated on `data-state=open`
  // rather than the `:popover-open` pseudo-class alone. dropdown.js already
  // stamps data-state="open" synchronously before showPopover() on every real
  // open (see dropdown-menu.gsx's own SUBMENUS doc comment) — this fixture
  // bypasses dropdown.js, so it must set the attribute itself to reproduce a
  // real open.
  const content = page.locator('[data-style-contract="submenu-right"]');
  const keyframes = await content.evaluate((element) => {
    element.setAttribute("data-state", "open");
    (element as HTMLElement).showPopover();
    return element.getAnimations().flatMap((animation) => {
      const effect = animation.effect as KeyframeEffect | null;
      return effect?.getKeyframes().map((frame) => frame.transform) ?? [];
    });
  });
  expect(keyframes.some((transform) => transform?.includes("-8px"))).toBe(true);
  await content.evaluate((element) => {
    element.setAttribute("data-state", "closed");
    (element as HTMLElement).hidePopover();
  });
});

// Before the 8-style port, CommandDialog bumped its input wrapper/input to
// h-12 (48px) once inside the dialog (assets/css/styles/default/command.css's
// old `[[data-gsxui-slot-command-dialog-content]_&]:h-12`, ultimately traced
// to upstream's OLD new-york-v4 CommandDialog wrapper className). The new
// base+8-styles design (apps/v4/registry/bases/base/ui/command.tsx +
// apps/v4/registry/styles/style-*.css MARK: Command sections, verified
// 2026-08-11) dropped that bump entirely: CommandInput's wrapper carries an
// unconditional h-8 (32px) from cn-command-input-group with no in-dialog
// variant, same genuine upstream removal already established for
// CommandItem's in-dialog padding bump (layer-precedence.spec.ts's "CommandDialog's
// in-dialog rounding" test). wrapperHeight/inputHeight below reflect that
// real, unbumped geometry; only `overflow: hidden` on the dialog itself is
// still real (it's baked into CommandDialog's shared DialogContent className
// in the .tsx, outside any per-style CSS the porter reads, so it silently
// dropped and is restored on command-dialog-content's recipe).
test("CommandDialog sizing is supplied by CSS ancestry", async ({ page }) => {
  await openStyleContract(page);
  await page.keyboard.press(process.platform === "darwin" ? "Meta+K" : "Control+K");

  const dialog = page.locator('[data-style-contract="command-dialog"]');
  await expect(dialog).toBeVisible();
  expect(
    await dialog.evaluate((element) => {
      const command = element.querySelector("[data-gsxui-slot-command]");
      const wrapper = element.querySelector("[data-gsxui-slot-command-input-wrapper]");
      const input = element.querySelector("[data-gsxui-slot-command-input]");
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
    wrapperHeight: "32px",
    inputHeight: "28px",
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

  // Before the 8-style port, the floating (viewport=false) content chrome was
  // a plain `border` plus `overflow-hidden` (assets/css/styles/default/
  // navigation-menu.css's old `[[data-gsxui-slot-navigation-menu]
  // [data-viewport=false]_&]:border`/`:overflow-hidden`). The new
  // base+8-styles design (apps/v4/registry/styles/style-*.css MARK:
  // Navigation Menu, verified 2026-08-11) replaced the border with a
  // box-shadow ring (`ring-1 ring-foreground/10`, the same border->ring
  // redesign already established for Command/Context Menu chrome) and drops
  // overflow-hidden entirely — this repo's navigation-menu-content recipe
  // faithfully carries that real upstream content, so the assertions below
  // check for the ring (via boxShadow) instead of borderWidth, and the real
  // overflow: auto.
  expect(
    await content.evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        borderWidth: css.borderWidth,
        hasRing: css.boxShadow.includes("0px 0px 0px 1px"),
        borderRadius: css.borderRadius,
        overflow: css.overflow,
        position: css.position,
      };
    }),
  ).toEqual({
    borderWidth: "0px",
    hasRing: true,
    borderRadius: "10px",
    overflow: "auto",
    position: "fixed",
  });
  expect(
    await trigger
      .locator("[data-gsxui-slot-navigation-menu-trigger-icon]")
      .evaluate((element) => getComputedStyle(element).rotate),
  ).toBe("180deg");
});
