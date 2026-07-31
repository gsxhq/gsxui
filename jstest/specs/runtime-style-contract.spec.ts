import { readFileSync } from "node:fs";

import type { Locator, Page } from "@playwright/test";

import { expect, test } from "../support/fixtures";
import { jstestDir } from "../support/paths";

type Entry = {
  component: string;
  slot: string;
  attribute: string;
  value: string;
  scenario: string;
};

const manifest = JSON.parse(
  readFileSync(`${jstestDir}/runtime-style-contract.json`, "utf8"),
) as Entry[];

const key = (entry: Entry) =>
  [
    entry.component,
    entry.slot,
    entry.attribute,
    entry.value,
    entry.scenario,
  ].join("\u0000");

test("real interactions cover the exact runtime-owned style contract", async ({
  page,
}) => {
  test.setTimeout(120_000);
  const observed = new Set<string>();

  async function observe(
    scenario: string,
    component: string,
    slot: string,
    locator: Locator,
    attribute: string,
    value: string,
  ) {
    const entry = { component, slot, attribute, value, scenario };
    expect(manifest.map(key)).toContain(key(entry));
    if (value === "") {
      await expect(locator).toHaveAttribute(attribute, "");
    } else {
      await expect(locator).toHaveAttribute(attribute, value);
    }
    observed.add(key(entry));
  }

  async function dialogLifecycle(
    scenario: string,
    component: string,
    route: string,
    triggerName: string,
    triggerSlot: string,
    contentSlot: string,
  ) {
    await page.goto(route);
    const trigger = page.locator(`[data-gsxui-slot-${triggerSlot}]`).filter({
      hasText: triggerName,
    });
    const content = page.locator(`[data-gsxui-slot-${contentSlot}]`).first();
    await trigger.click();
    await observe(
      scenario,
      component,
      triggerSlot,
      trigger,
      "aria-expanded",
      "true",
    );
    await observe(
      scenario,
      component,
      contentSlot,
      content,
      "data-state",
      "open",
    );
  }

  await dialogLifecycle(
    "alert-dialog-lifecycle",
    "alert-dialog",
    "/x/alert-dialog/basic",
    "Show dialog",
    "alert-dialog-trigger",
    "alert-dialog-content",
  );
  await dialogLifecycle(
    "dialog-lifecycle",
    "dialog",
    "/x/dialog/basic",
    "Delete account",
    "dialog-trigger",
    "dialog-content",
  );
  await dialogLifecycle(
    "drawer-lifecycle",
    "drawer",
    "/x/drawer/basic",
    "Open Drawer",
    "drawer-trigger",
    "drawer-content",
  );
  await dialogLifecycle(
    "sheet-lifecycle",
    "sheet",
    "/x/sheet/basic",
    "Edit Profile",
    "sheet-trigger",
    "sheet-content",
  );

  await page.goto("/x/hover-card/basic");
  const hoverTrigger = page.locator("[data-gsxui-slot-hover-card-trigger]");
  const hoverContent = page.locator("[data-gsxui-slot-hover-card-content]");
  await hoverTrigger.hover();
  await observe(
    "hover-card-lifecycle",
    "hover-card",
    "hover-card-content",
    hoverContent,
    "data-state",
    "open",
  );

  await page.goto("/x/popover/basic");
  const popoverTrigger = page
    .locator("[data-gsxui-slot-popover-trigger]")
    .first();
  const popoverContent = page.locator("[data-gsxui-slot-popover-content]");
  await popoverTrigger.click();
  await observe(
    "popover-lifecycle",
    "popover",
    "popover-trigger",
    popoverTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "popover-lifecycle",
    "popover",
    "popover-content",
    popoverContent,
    "data-state",
    "open",
  );

  await page.goto("/x/tooltip/basic");
  const tooltipTrigger = page
    .locator("[data-gsxui-slot-tooltip-trigger]")
    .first();
  const tooltipContent = page.locator("[data-gsxui-slot-tooltip-content]");
  await tooltipTrigger.hover();
  await observe(
    "tooltip-lifecycle",
    "tooltip",
    "tooltip-content",
    tooltipContent,
    "data-state",
    "open",
  );

  await page.goto("/x/command/basic");
  const commandItem = page.locator("[data-gsxui-slot-command-item]").first();
  await observe(
    "command-runtime",
    "command",
    "command-item",
    commandItem,
    "data-selected",
    "true",
  );
  await observe(
    "command-runtime",
    "command",
    "command-item",
    commandItem,
    "aria-selected",
    "true",
  );

  await page.goto("/x/combobox/basic");
  const comboboxInput = page.locator("[data-gsxui-slot-combobox-input]");
  const comboboxContent = page.locator("[data-gsxui-slot-combobox-content]");
  await comboboxInput.click();
  const comboboxItem = page
    .locator("[data-gsxui-slot-combobox-item]")
    .first();
  await comboboxItem.hover();
  await observe(
    "combobox-lifecycle",
    "combobox",
    "combobox-input",
    comboboxInput,
    "aria-expanded",
    "true",
  );
  await observe(
    "combobox-lifecycle",
    "combobox",
    "combobox-content",
    comboboxContent,
    "data-state",
    "open",
  );
  await observe(
    "combobox-lifecycle",
    "combobox",
    "combobox-item",
    comboboxItem,
    "data-highlighted",
    "true",
  );
  await comboboxItem.click();
  await observe(
    "combobox-lifecycle",
    "combobox",
    "combobox-item",
    comboboxItem,
    "aria-selected",
    "true",
  );

  await page.goto("/x/dropdown-menu/submenu");
  const dropdownTrigger = page.locator("[data-gsxui-slot-dropdown-menu-trigger]").first();
  const dropdownContent = page.locator(
    "[data-gsxui-slot-dropdown-menu-content]",
  );
  await dropdownTrigger.click();
  await observe(
    "dropdown-lifecycle",
    "dropdown-menu",
    "dropdown-menu-trigger",
    dropdownTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "dropdown-lifecycle",
    "dropdown-menu",
    "dropdown-menu-content",
    dropdownContent,
    "data-state",
    "open",
  );
  const dropdownSubTrigger = page.locator(
    "[data-gsxui-slot-dropdown-menu-sub-trigger]",
  );
  const dropdownSubContent = page.locator(
    "[data-gsxui-slot-dropdown-menu-sub-content]",
  );
  await dropdownSubTrigger.hover();
  await observe(
    "dropdown-lifecycle",
    "dropdown-menu",
    "dropdown-menu-sub-trigger",
    dropdownSubTrigger,
    "data-state",
    "open",
  );
  await observe(
    "dropdown-lifecycle",
    "dropdown-menu",
    "dropdown-menu-sub-trigger",
    dropdownSubTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "dropdown-lifecycle",
    "dropdown-menu",
    "dropdown-menu-sub-content",
    dropdownSubContent,
    "data-state",
    "open",
  );

  await page.goto("/x/context-menu/full");
  await page.locator("[data-gsxui-slot-context-menu-trigger]").click({
    button: "right",
  });
  const contextContent = page.locator(
    "[data-gsxui-slot-context-menu-content]",
  );
  await observe(
    "context-menu-lifecycle",
    "context-menu",
    "context-menu-content",
    contextContent,
    "data-state",
    "open",
  );
  const contextSubTrigger = page.locator(
    "[data-gsxui-slot-context-menu-sub-trigger]",
  );
  const contextSubContent = page.locator(
    "[data-gsxui-slot-context-menu-sub-content]",
  );
  await contextSubTrigger.hover();
  await observe(
    "context-menu-lifecycle",
    "context-menu",
    "context-menu-sub-trigger",
    contextSubTrigger,
    "data-state",
    "open",
  );
  await observe(
    "context-menu-lifecycle",
    "context-menu",
    "context-menu-sub-trigger",
    contextSubTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "context-menu-lifecycle",
    "context-menu",
    "context-menu-sub-content",
    contextSubContent,
    "data-state",
    "open",
  );

  await page.goto("/x/menubar/full");
  const menubarTrigger = page
    .locator("[data-gsxui-slot-menubar-trigger]")
    .first();
  const menubarContent = page
    .locator("[data-gsxui-slot-menubar-content]")
    .first();
  await menubarTrigger.click();
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-trigger",
    menubarTrigger,
    "data-state",
    "open",
  );
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-trigger",
    menubarTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-content",
    menubarContent,
    "data-state",
    "open",
  );
  const menubarSubTrigger = page
    .locator("[data-gsxui-slot-menubar-sub-trigger]")
    .first();
  const menubarSubContent = page
    .locator("[data-gsxui-slot-menubar-sub-content]")
    .first();
  await menubarSubTrigger.hover();
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-sub-trigger",
    menubarSubTrigger,
    "data-state",
    "open",
  );
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-sub-trigger",
    menubarSubTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "menubar-lifecycle",
    "menubar",
    "menubar-sub-content",
    menubarSubContent,
    "data-state",
    "open",
  );

  await page.goto("/x/navigation-menu/basic");
  const navigationTrigger = page.locator(
    "button[data-gsxui-slot-navigation-menu-trigger]",
  );
  const navigationContent = page.locator(
    "[data-gsxui-slot-navigation-menu-content]",
  );
  const navigationIndicator = page.locator(
    "[data-gsxui-slot-navigation-menu-indicator]",
  );
  await navigationTrigger.click();
  await observe(
    "navigation-menu-lifecycle",
    "navigation-menu",
    "navigation-menu-trigger",
    navigationTrigger,
    "data-state",
    "open",
  );
  await observe(
    "navigation-menu-lifecycle",
    "navigation-menu",
    "navigation-menu-trigger",
    navigationTrigger,
    "aria-expanded",
    "true",
  );
  await observe(
    "navigation-menu-lifecycle",
    "navigation-menu",
    "navigation-menu-content",
    navigationContent,
    "data-state",
    "open",
  );
  await observe(
    "navigation-menu-lifecycle",
    "navigation-menu",
    "navigation-menu-indicator",
    navigationIndicator,
    "data-state",
    "visible",
  );
  await navigationTrigger.click();
  await observe(
    "navigation-menu-lifecycle",
    "navigation-menu",
    "navigation-menu-indicator",
    navigationIndicator,
    "data-state",
    "hidden",
  );

  await page.goto("/x/select/basic");
  const selectTrigger = page.locator("[data-gsxui-slot-select-trigger]");
  const selectContent = page.locator("[data-gsxui-slot-select-content]");
  await selectTrigger.click();
  await observe(
    "select-lifecycle",
    "select",
    "select-trigger",
    selectTrigger,
    "data-state",
    "open",
  );
  await observe(
    "select-lifecycle",
    "select",
    "select-content",
    selectContent,
    "data-state",
    "open",
  );
  await observe(
    "select-lifecycle",
    "select",
    "select-item",
    page.locator("[data-gsxui-slot-select-item]").first(),
    "aria-selected",
    "true",
  );

  await page.goto("/x/input-otp/basic");
  await page.locator("[data-gsxui-slot-input-otp-input]").focus();
  const activeSlot = page.locator(
    "[data-gsxui-slot-input-otp-slot][data-active='true']",
  );
  await observe(
    "input-otp-caret",
    "input-otp",
    "input-otp-slot",
    activeSlot,
    "data-active",
    "true",
  );
  await observe(
    "input-otp-caret",
    "input-otp",
    "input-otp-caret-overlay",
    activeSlot.locator("[data-gsxui-slot-input-otp-caret-overlay]"),
    "data-gsxui-slot-input-otp-caret-overlay",
    "",
  );
  await observe(
    "input-otp-caret",
    "input-otp",
    "input-otp-caret",
    activeSlot.locator("[data-gsxui-slot-input-otp-caret]"),
    "data-gsxui-slot-input-otp-caret",
    "",
  );

  await page.goto("/x/toaster/types");
  await page.evaluate(() => {
    for (let index = 0; index < 4; index++) {
      window.gsxui.toast(`Runtime queue ${index}`, { duration: 60_000 });
    }
  });
  const toaster = page.locator("[data-gsxui-slot-toaster]").first();
  const toast = toaster.locator("[data-gsxui-slot-toast]", {
    hasText: "Runtime queue 3",
  });
  const queuedToast = toaster.locator("[data-gsxui-slot-toast]", {
    hasText: "Runtime queue 0",
  });
  await observe(
    "sonner-lifecycle",
    "toast",
    "toast",
    toast,
    "data-state",
    "open",
  );
  await observe(
    "sonner-lifecycle",
    "toast",
    "toast",
    queuedToast,
    "data-visible",
    "false",
  );
  await observe(
    "sonner-lifecycle",
    "toast",
    "toast",
    toast,
    "data-visible",
    "true",
  );
  await toast.hover();
  await observe(
    "sonner-lifecycle",
    "toaster",
    "toaster",
    toaster,
    "data-expanded",
    "true",
  );
  await page.mouse.move(0, 0);
  await observe(
    "sonner-lifecycle",
    "toaster",
    "toaster",
    toaster,
    "data-expanded",
    "false",
  );
  await toast.locator("[data-gsxui-slot-toast-close]").click();
  await observe(
    "sonner-lifecycle",
    "toast",
    "toast",
    toast,
    "data-state",
    "closed",
  );
  expect([...observed].sort()).toEqual(manifest.map(key).sort());
});
