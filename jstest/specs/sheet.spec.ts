import { expect, test } from "../support/fixtures";

for (const side of ["top", "bottom"]) {
  test(`${side} direction demo renders a complete, padded sheet`, async ({
    page,
  }) => {
    const response = await page.goto("/x/sheet/directions");
    expect(response?.status(), "sheet directions fixture response").toBe(200);

    await page.getByRole("button", { name: side, exact: true }).click();
    const content = page.locator(
      `[data-gsxui-slot~="sheet-content"][data-side="${side}"]`,
    );
    await expect(content).toBeVisible();
    await expect(
      content.locator('[data-gsxui-slot~="sheet-header"]'),
    ).toHaveCount(1);
    await expect(
      content.locator('[data-gsxui-slot~="sheet-footer"]'),
    ).toHaveCount(1);
    await expect(content.locator("input")).toHaveCount(2);

    const geometry = await content.evaluate((element) => {
      const panel = element.getBoundingClientRect();
      const title = element
        .querySelector('[data-gsxui-slot~="sheet-title"]')
        ?.getBoundingClientRect();
      return {
        height: panel.height,
        titleInsidePanel:
          !!title && title.top >= panel.top && title.bottom <= panel.bottom,
      };
    });
    expect(geometry.height).toBeGreaterThan(160);
    expect(geometry.titleInsidePanel).toBe(true);
  });
}
