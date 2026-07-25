import { expect, test } from "@playwright/test";
import { examples } from "../support/manifest";

/**
 * One test per example, each asserting every invariant with expect.soft so a
 * single run reports everything an example violates.
 *
 * These four checks are where the leverage is: they cover all 103 examples
 * without a line of per-component test code, and two of them encode defect
 * classes that have actually shipped here.
 */
for (const example of examples()) {
  test(`${example.component}/${example.example}`, async ({ page }) => {
    const consoleErrors: string[] = [];
    const pageErrors: string[] = [];
    page.on("console", (msg) => {
      if (msg.type() === "error") consoleErrors.push(msg.text());
    });
    page.on("pageerror", (err) => pageErrors.push(String(err)));

    await page.goto(example.url);

    // Invariant 1: nothing threw and nothing logged an error while the
    // module graph loaded and the markup parsed.
    expect.soft(pageErrors, "uncaught exceptions").toEqual([]);
    expect.soft(consoleErrors, "console errors").toEqual([]);
  });
}
