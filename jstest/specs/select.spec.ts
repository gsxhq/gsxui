import { expect, test } from "../support/fixtures";

test("Select trigger state follows open and every native close path", async ({
  page,
}) => {
  const response = await page.goto("/x/select/basic");
  expect(response?.status(), "select example response").toBe(200);

  const trigger = page.locator("[data-gsxui-slot~='select-trigger']");
  const content = page.locator("[data-gsxui-slot~='select-content']");

  await expect(trigger).toHaveAttribute("data-state", "closed");
  await expect(trigger).toHaveAttribute("aria-expanded", "false");

  const open = async () => {
    await trigger.click();
    await expect(content).toHaveAttribute("data-state", "open");
    await expect(trigger).toHaveAttribute("data-state", "open");
    await expect(trigger).toHaveAttribute("aria-expanded", "true");
  };

  const expectClosed = async () => {
    await expect(content).toHaveAttribute("data-state", "closed");
    await expect(trigger).toHaveAttribute("data-state", "closed");
    await expect(trigger).toHaveAttribute("aria-expanded", "false");
  };

  await open();
  await page.locator("[data-gsxui-slot~='select-item']").first().click();
  await expectClosed();

  await open();
  await page.keyboard.press("Escape");
  await expectClosed();

  await open();
  await page.locator("[data-harness-root]").click({ position: { x: 1, y: 1 } });
  await expectClosed();

  await open();
  await content.evaluate((element: HTMLElement) => element.hidePopover());
  await expectClosed();
});
