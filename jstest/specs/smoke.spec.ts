import { expect, test } from "../support/fixtures";
import { examples } from "../support/manifest";

test("the manifest reaches the specs", () => {
  const all = examples();
  expect(all.length).toBeGreaterThan(100);
  expect(all).toContainEqual({
    component: "dropdown-menu",
    example: "checkboxes",
    url: "/x/dropdown-menu/checkboxes",
  });
});

test("markup, stylesheet and behavior JS all arrive", async ({ page }) => {
  await page.goto("/x/toggle/basic");

  // Markup: rendered through the real gsx component.
  const toggle = page.locator("[data-gsxui-slot-toggle]").first();
  await expect(toggle).toBeVisible();

  // CSS: bg-background resolves to an opaque colour. Without the stylesheet
  // the body background is rgba(0, 0, 0, 0).
  const background = await page.evaluate(
    () => getComputedStyle(document.body).backgroundColor,
  );
  expect(background).not.toBe("rgba(0, 0, 0, 0)");

  // JS: ui/toggle.js is bound and flips both attributes on click.
  await expect(toggle).toHaveAttribute("aria-pressed", "false");
  await toggle.click();
  await expect(toggle).toHaveAttribute("aria-pressed", "true");
  await expect(toggle).toHaveAttribute("data-state", "on");
});

test("the shim records what the real modules would have registered", async ({ page }) => {
  await page.goto("/registrations");

  const registrations = await page.evaluate(() => window.__gsxuiRegistrations);
  expect(registrations.length).toBeGreaterThan(50);

  // Every entry is attributed to a real module, not "unknown" — the stack
  // walk in shim.js is the only thing that can break this.
  const modules = new Set(registrations.map((r) => r.module));
  expect(modules).not.toContain("unknown");
  expect(modules).toContain("dropdown-menu.js");
  expect(modules).toContain("toggle.js");
});

// Confirms the /static/ root choice in Task 1's harness: Tailwind bundles
// @fontsource-variable/geist, whose CSS carries url() references, and every
// page's shell links the compiled stylesheet. If those font requests (or any
// other subresource) 404, every page load would log a failed request and
// Task 4's clean-load invariant would be noisy — this keeps that invariant
// meaningful by catching real breakage here instead of filtering it out.
test("no failed subresource requests", async ({ page }) => {
  const failed: string[] = [];
  page.on("requestfailed", (r) => failed.push(r.url()));
  page.on("response", (r) => {
    if (r.status() >= 400) failed.push(`${r.status()} ${r.url()}`);
  });
  await page.goto("/x/toggle/basic");
  expect(failed).toEqual([]);
});
