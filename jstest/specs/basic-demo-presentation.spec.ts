import { expect, test } from "../support/fixtures";

type TriggerDemo = {
  route: string;
  triggerSlot: string;
  count: number;
};

const triggerDemos: TriggerDemo[] = [
  {
    route: "alert-dialog/basic",
    triggerSlot: "alert-dialog-trigger",
    count: 1,
  },
  {
    route: "dialog/basic",
    triggerSlot: "dialog-trigger",
    count: 1,
  },
  {
    route: "dialog/footer",
    triggerSlot: "dialog-trigger",
    count: 1,
  },
  {
    route: "drawer/basic",
    triggerSlot: "drawer-trigger",
    count: 1,
  },
  {
    route: "drawer/directions",
    triggerSlot: "drawer-trigger",
    count: 4,
  },
  {
    route: "dropdown-menu/basic",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/checkboxes",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/destructive",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/radios",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown-menu/submenu",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "calendar/datepicker",
    triggerSlot: "popover-trigger",
    count: 1,
  },
  {
    route: "popover/basic",
    triggerSlot: "popover-trigger",
    count: 1,
  },
  {
    route: "sheet/basic",
    triggerSlot: "sheet-trigger",
    count: 1,
  },
  {
    route: "sheet/directions",
    triggerSlot: "sheet-trigger",
    count: 4,
  },
  {
    route: "tooltip/basic",
    triggerSlot: "tooltip-trigger",
    count: 1,
  },
  {
    route: "tooltip/wide",
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

    const triggers = page.locator(`[data-gsxui-slot-${demo.triggerSlot}]`);
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
