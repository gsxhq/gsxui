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

    // Invariant 2: no ghost boxes. A closed popover must compute
    // `display: none`. An author display utility (a bare `block`/`grid`)
    // beats the UA stylesheet's closed-popover rule and leaves an invisible
    // but hit-testable box — this shipped in both dialog and sidebar, and
    // the fix is to gate the utility on `:open` (`open:grid`).
    const ghosts = await page.evaluate(() =>
      [...document.querySelectorAll("[popover]")]
        .filter((el) => !el.matches(":popover-open"))
        .filter((el) => getComputedStyle(el).display !== "none")
        .map((el) => ({
          slot: (el as HTMLElement).dataset.slot ?? null,
          id: el.id || null,
          display: getComputedStyle(el).display,
          classes: el.className,
        })),
    );
    expect.soft(ghosts, "closed popovers computing a display other than none").toEqual([]);

    // Invariant 3: no duplicate ids. One example renders alone on the page,
    // so any collision is within a single example's own markup — either the
    // example reuses an id, or a component generates one non-uniquely.
    const duplicateIds = await page.evaluate(() => {
      const counts = new Map<string, number>();
      for (const el of document.querySelectorAll("[id]")) {
        counts.set(el.id, (counts.get(el.id) ?? 0) + 1);
      }
      return [...counts].filter(([, n]) => n > 1).map(([id, n]) => ({ id, count: n }));
    });
    expect.soft(duplicateIds, "duplicate element ids").toEqual([]);
  });
}
