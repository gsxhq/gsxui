import { expect, test } from "../support/fixtures";

type TriggerDemo = {
  route: string;
  behaviorAttribute: string;
  triggerSlot: string;
  count: number;
};

const triggerDemos: TriggerDemo[] = [
  {
    route: "alert-dialog/basic",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "alert-dialog-trigger",
    count: 1,
  },
  {
    route: "dialog/basic",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "dialog-trigger",
    count: 1,
  },
  {
    route: "dialog/footer",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "dialog-trigger",
    count: 1,
  },
  {
    route: "drawer/basic",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "drawer-trigger",
    count: 1,
  },
  {
    route: "drawer/directions",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "drawer-trigger",
    count: 4,
  },
  {
    route: "dropdown-menu/basic",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/checkboxes",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/destructive",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/radios",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/submenu",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "calendar/datepicker",
    behaviorAttribute: "data-gsxui-popover-trigger",
    triggerSlot: "popover-trigger",
    count: 1,
  },
  {
    route: "popover/basic",
    behaviorAttribute: "data-gsxui-popover-trigger",
    triggerSlot: "popover-trigger",
    count: 1,
  },
  {
    route: "sheet/basic",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "sheet-trigger",
    count: 1,
  },
  {
    route: "sheet/directions",
    behaviorAttribute: "data-gsxui-dialog-trigger",
    triggerSlot: "sheet-trigger",
    count: 4,
  },
  {
    route: "tooltip/basic",
    behaviorAttribute: "data-gsxui-tooltip-trigger",
    triggerSlot: "tooltip-trigger",
    count: 1,
  },
  {
    route: "tooltip/wide",
    behaviorAttribute: "data-gsxui-tooltip-trigger",
    triggerSlot: "tooltip-trigger",
    count: 1,
  },
];

for (const demo of triggerDemos) {
  test(`${demo.route} presents one styled control per trigger`, async ({
    page,
  }) => {
    const response = await page.goto(`/x/${demo.route}`);
    expect(response?.status(), `${demo.route} fixture response`).toBe(200);

    const triggers = page.locator(`[${demo.behaviorAttribute}]`);
    await expect(triggers).toHaveCount(demo.count);
    for (let index = 0; index < demo.count; index++) {
      const trigger = triggers.nth(index);
      await expect(trigger).toHaveAttribute("data-gsxui-slot-button", "");
      await expect(trigger).toHaveAttribute(
        `data-gsxui-slot-${demo.triggerSlot}`,
        "",
      );
    }
  });
}

test("combobox/basic does not render a disconnected value label", async ({
  page,
}) => {
  const response = await page.goto("/x/combobox/basic");
  expect(response?.status(), "combobox/basic fixture response").toBe(200);

  await expect(
    page.locator('[data-gsxui-slot-combobox-value]'),
  ).toHaveCount(0);
  await expect(
    page.getByText("Choose a framework", { exact: true }),
  ).toHaveCount(0);
});

test("Button variants and sizes keep Nova presentation in normal site scope", async ({
  page,
}) => {
  const response = await page.goto("/x/button/variants");
  expect(response?.status(), "button/variants fixture response").toBe(200);
  await expect(page.locator("body")).not.toHaveAttribute(
    "data-theme-button-preview",
  );

  const box = async (name: string) =>
    page.getByRole("button", { name, exact: true }).evaluate((element) => {
      const css = getComputedStyle(element);
      return {
        width: css.width,
        height: css.height,
        backgroundColor: css.backgroundColor,
        borderTopWidth: css.borderTopWidth,
      };
    });

  const defaultButton = await box("Default");
  expect(defaultButton.height).toBe("32px");
  expect(defaultButton.backgroundColor).not.toBe("rgba(0, 0, 0, 0)");
  expect((await box("Large")).height).toBe("36px");
  expect(await box("Small icon")).toMatchObject({
    width: "28px",
    height: "28px",
  });
  expect(await box("Large icon")).toMatchObject({
    width: "36px",
    height: "36px",
  });
  expect((await box("Outline")).borderTopWidth).toBe("1px");
  expect((await box("Secondary")).backgroundColor).not.toBe(
    defaultButton.backgroundColor,
  );
});

test("Button fallback cannot style a theme-preview body", async ({ page }) => {
  // The fallback stylesheet (web/site-button.css) still owns every element
  // that carries Button's marker without rendering through ui.Button —
  // PaginationLink and Calendar's nav/day buttons author no class at all.
  // ui.Button itself is no longer a witness for this: it ships concrete
  // utilities compiled from its style recipe, so its inline-flex is intrinsic
  // and no body-level scope can switch it off. A PaginationLink is the
  // element whose presentation the fallback still supplies.
  const response = await page.goto("/x/pagination/basic");
  expect(response?.status(), "pagination/basic fixture response").toBe(200);

  const link = page.locator("[data-gsxui-slot-pagination-link]").first();
  await expect(link).toHaveCSS("display", "inline-flex");
  await page.locator("body").evaluate((body) => {
    body.setAttribute("data-theme-button-preview", "");
  });
  await expect(link).not.toHaveCSS("display", "inline-flex");
});
