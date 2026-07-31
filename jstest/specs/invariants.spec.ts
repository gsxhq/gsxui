import { expect, test } from "../support/fixtures";
import { examples } from "../support/manifest";
import { allowedOverlaps } from "../support/selector-allowlist";

/**
 * One test per example, each asserting every invariant with expect.soft so a
 * single run reports everything an example violates.
 *
 * These four checks are where the leverage is: they cover all 103 examples
 * without a line of per-component test code, and two of them encode defect
 * classes that have actually shipped here.
 */
for (const example of examples()) {
  test(`${example.component}/${example.example}`, async ({ page, registrations }) => {
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

    // Invariant 2: no ghost boxes. This is one defect class with two UA
    // mechanisms: a closed `[popover]` must compute `display: none` (the
    // UA's closed-popover rule), and a `<dialog>` without the `open`
    // attribute must too (the UA's `dialog:not([open])` rule). An author
    // display utility (a bare `block`/`grid`) beats either UA rule and
    // leaves an invisible but hit-testable box — this shipped in both
    // dialog (a `<dialog>`) and sidebar (a `[popover]`), and the fix is to
    // gate the utility on `:open` (`open:grid`).
    const ghosts = await page.evaluate(() => {
      const describe = (el: Element, mechanism: string) => ({
        mechanism,
        slot: (el as HTMLElement).dataset.slot ?? null,
        id: el.id || null,
        display: getComputedStyle(el).display,
        classes: el.className,
      });
      const closedPopovers = [...document.querySelectorAll("[popover]")]
        .filter((el) => !el.matches(":popover-open"))
        .filter((el) => getComputedStyle(el).display !== "none")
        .map((el) => describe(el, "popover"));
      const closedDialogs = [...document.querySelectorAll("dialog:not([open])")]
        .filter((el) => getComputedStyle(el).display !== "none")
        .map((el) => describe(el, "dialog"));
      return [...closedPopovers, ...closedDialogs];
    });
    expect.soft(ghosts, "closed overlays computing a display other than none").toEqual([]);

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

    // Invariant 4: selector disjointness. ui/gsxui.js keys its registry by
    // `${type}:${capture}` alone and dispatches to EVERY handler whose
    // selector matches, so if two modules both match one element for one
    // (type, capture) pair, both handlers run on a single event. That is
    // exactly the hook-prefix collision that shipped in Tier 4 Batch B:
    // dropdown-menu.js and context-menu.js both matched data-gsxui-menu-*, so
    // one click on a checkbox item fired two gsxui:change events and left
    // the state unchanged.
    //
    // The check is same-element, not ancestor-chain, and that is deliberate.
    // gsxui.js dispatches via target.closest(selector); when one element
    // matches two modules' selectors, an event on it resolves to that same
    // element for both. Nested elements matching different modules is
    // ordinary composition (a dialog inside a dropdown) and is not a defect.
    const overlaps = await page.evaluate((regs) => {
      const found: { key: string; tag: string; pairs: [string, string][] }[] = [];
      for (const el of document.querySelectorAll("*")) {
        const byKey = new Map<string, Map<string, string>>();
        for (const reg of regs) {
          let matches = false;
          try {
            matches = el.matches(reg.selector);
          } catch {
            // Backstop only. Unparseable selectors are asserted empty, once
            // per run, by selector-coverage.spec.ts — this catch is what
            // keeps one bad selector from aborting the whole sweep, not the
            // thing that reports it.
            continue;
          }
          if (!matches) continue;
          const key = `${reg.type}:${reg.capture}`;
          if (!byKey.has(key)) byKey.set(key, new Map());
          byKey.get(key)!.set(reg.module, reg.selector);
        }
        for (const [key, mods] of byKey) {
          if (mods.size > 1) {
            // [module, selector] pairs, sorted together. Sorting the two as
            // independent arrays lets the failure output pair the wrong
            // selector with the wrong module — and that output is the entire
            // value of this check when it finally fires.
            found.push({
              key,
              tag: el.tagName.toLowerCase(),
              pairs: [...mods].sort((a, b) =>
                a[0] === b[0] ? a[1].localeCompare(b[1]) : a[0].localeCompare(b[0]),
              ),
            });
          }
        }
      }
      return found;
    }, registrations);

    const unexpected = overlaps.filter((o) => {
      const modules = o.pairs.map(([module]) => module);
      return !allowedOverlaps.some(
        (a) =>
          a.key === o.key &&
          modules.length === 2 &&
          a.modules[0] === modules[0] &&
          a.modules[1] === modules[1],
      );
    });
    expect.soft(unexpected, "elements claimed by two modules for one event").toEqual([]);
  });
}
