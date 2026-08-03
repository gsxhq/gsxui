import { expect, test } from "../support/fixtures";

test.describe("home showcase", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/");
  });

  test("renders all four showcase cards", async ({ page }) => {
    const section = page.locator("#components");
    await expect(section.getByText("Sign in", { exact: true }).first()).toBeVisible();
    await expect(section.getByText("Preferences")).toBeVisible();
    await expect(section.getByText("Usage")).toBeVisible();
    await expect(section.getByText("Interactive, no framework")).toBeVisible();
  });

  test("dialog opens and closes from the overlays card", async ({ page }) => {
    const trigger = page.locator("#components [data-gsxui-slot-dialog-trigger]");
    await trigger.click();
    const dialog = page.locator("#components dialog[open]");
    await expect(dialog).toBeVisible();
    await expect(dialog.getByText("Edit profile")).toBeVisible();
    await page.keyboard.press("Escape");
    await expect(dialog).toBeHidden();
  });

  test("dropdown menu opens from the overlays card", async ({ page }) => {
    const trigger = page.locator("#components [data-gsxui-slot-dropdown-menu-trigger]");
    await trigger.click();
    await expect(page.getByText("My Account")).toBeVisible();
  });

  test("toast appends into the toaster viewport", async ({ page }) => {
    await page.getByRole("tab", { name: "Feedback" }).click();
    await page.locator("#home-showcase-toast-btn").click();
    await expect(page.locator("#gsxui-toaster [data-gsxui-slot-toast]").first()).toBeVisible();
  });
});
