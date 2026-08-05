import { expect, test } from "@playwright/test";

// Regression: after selecting an item, clearing the input's text must
// stick. The self-healing init pass fires on filter()'s show/hide
// mutations while the user edits; before the live-interaction guard,
// initRoot saw "checked item + empty input", refilled the committed label,
// and re-stamped gsxuiCommitted — the user could never clear the text.
test("cleared input text stays cleared after selecting an item", async ({
  page,
}) => {
  await page.goto("/x/combobox/basic");
  const input = page.locator("[data-gsxui-slot-combobox-input]");

  await input.click();
  await page
    .locator('[data-gsxui-slot-combobox-item]:not([aria-disabled="true"])')
    .first()
    .click();
  await expect(input).toHaveValue("Next.js");

  await input.click();
  await input.press("ControlOrMeta+a");
  await input.press("Delete");
  await expect(input).toHaveValue("");

  // The stomp arrived one microtask after the input event; hold long
  // enough for any scheduled re-init pass (and one more mutation cycle)
  // to have run before asserting the clear survived.
  await page.waitForTimeout(100);
  await expect(input).toHaveValue("");
  await expect(input).not.toHaveAttribute("data-gsxui-committed", "true");
});
