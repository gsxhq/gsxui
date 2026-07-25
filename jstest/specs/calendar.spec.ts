import { expect, test } from "../support/fixtures";

const BASIC = "/x/calendar/basic";

// The grid's 42 data-date values, in DOM order.
async function gridDates(page: import("@playwright/test").Page) {
  return page.$$eval("[data-gsxui-calendar-day]", (els) =>
    els.map((el) => el.getAttribute("data-date")),
  );
}

test("the server renders a complete grid before any JS runs", async ({ page }) => {
  await page.goto(BASIC);
  const dates = await gridDates(page);
  expect(dates).toHaveLength(42);
  expect(dates[0]).toBe("2025-12-28");
  expect(dates[41]).toBe("2026-02-07");
});

test("next and previous navigate a month at a time", async ({ page }) => {
  await page.goto(BASIC);

  await page.click("[data-gsxui-calendar-next]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-02",
  );

  await page.click("[data-gsxui-calendar-prev]");
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "2026-01",
  );
});

test("navigation never creates or destroys a cell", async ({ page }) => {
  await page.goto(BASIC);
  const before = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);

  const first = await page.$("[data-gsxui-calendar-day]");
  await page.click("[data-gsxui-calendar-next]");

  const after = await page.$$eval("[data-gsxui-calendar-day]", (els) => els.length);
  expect(after).toBe(before);

  // The same element object is still there — proof it was updated, not replaced.
  const stillAttached = await first!.evaluate((el) => el.isConnected);
  expect(stillAttached).toBe(true);
});

// The single largest risk in the design: monthGrid exists in Go and in JS and
// the two must agree. Navigate client-side to a month, then compare against
// what the server renders for that month directly. Both pages come from the
// same harness, so this is a real cross-implementation diff.
const AGREEMENT_MONTHS = [
  { clicks: 1, month: "2026-02", why: "28-day February" },
  { clicks: 11, month: "2026-12", why: "year boundary ahead" },
  { clicks: -1, month: "2025-12", why: "year boundary behind" },
  { clicks: -23, month: "2024-02", why: "leap February" },
];

for (const { clicks, month, why } of AGREEMENT_MONTHS) {
  test(`Go and JS agree on ${month} (${why})`, async ({ page }) => {
    await page.goto(BASIC);
    const button = clicks > 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
    for (let i = 0; i < Math.abs(clicks); i++) await page.click(button);

    await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
      "data-gsxui-calendar-month",
      month,
    );
    const clientDates = await gridDates(page);

    // Same month, rendered by Go.
    await page.goto(`${BASIC}?month=${month}`);
    const serverDates = await gridDates(page);

    expect(clientDates).toEqual(serverDates);
  });
}
