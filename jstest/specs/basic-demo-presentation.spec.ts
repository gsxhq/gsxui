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
    route: "dropdown/basic",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown/checkboxes",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown/destructive",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown/radios",
    behaviorAttribute: "data-gsxui-dropdown-trigger",
    triggerSlot: "dropdown-menu-trigger",
    count: 1,
  },
  {
    route: "dropdown/submenu",
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
      await expect(triggers.nth(index)).toHaveAttribute(
        "data-gsxui-slot",
        /(?:^|\s)button(?:\s|$)/,
      );
      await expect(triggers.nth(index)).toHaveAttribute(
        "data-gsxui-slot",
        new RegExp(`(?:^|\\s)${demo.triggerSlot}(?:\\s|$)`),
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
    page.locator('[data-gsxui-slot~="combobox-value"]'),
  ).toHaveCount(0);
  await expect(
    page.getByText("Choose a framework", { exact: true }),
  ).toHaveCount(0);
});
