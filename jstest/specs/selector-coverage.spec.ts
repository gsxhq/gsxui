import { expect, test } from "../support/fixtures";
import { examples } from "../support/manifest";
import { allowedUnmatched } from "../support/selector-allowlist";

/**
 * Corpus-wide checks over the recorded delegation registry.
 *
 * These live in their own file because they are unions over EVERY example,
 * and a per-example test structurally cannot see a union: a selector that
 * matches nothing on the dropdown page is unremarkable on its own, and only
 * means something once all 103 pages have failed to match it.
 *
 * Scope note: the registry covers delegated registrations only. See the
 * SCOPE comment in jstest/harness/shim.js.
 */

/** One entry per distinct (module, selector) — the same selector bound for
 * several event types is one thing to check, not three. */
function distinctSelectors(regs: { module: string; selector: string }[]) {
  const byKey = new Map<string, { module: string; selector: string }>();
  for (const reg of regs) {
    byKey.set(`${reg.module}\0${reg.selector}`, {
      module: reg.module,
      selector: reg.selector,
    });
  }
  return [...byKey.values()].sort((a, b) =>
    a.module === b.module
      ? a.selector.localeCompare(b.selector)
      : a.module.localeCompare(b.module),
  );
}

/**
 * Every recorded selector must parse.
 *
 * An unparseable selector is not a cosmetic problem: ui/gsxui.js's dispatch
 * calls target.closest(selector) on every event of that type, so a typo like
 * `[data-gsxui-x=]` throws on every such event in production. It is also
 * invisible to the disjointness invariant, whose per-element el.matches()
 * throws and is caught — one bad selector silently drops out of that check
 * for the whole run and the suite stays green. So it is validated once,
 * up front, here.
 */
test("every recorded selector parses", async ({ page, registrations }) => {
  await page.goto("/registrations");

  const unparseable = await page.evaluate(
    (entries) =>
      entries
        .map((entry) => {
          try {
            document.querySelector(entry.selector);
            return null;
          } catch (err) {
            return { ...entry, error: String(err) };
          }
        })
        .filter((x) => x !== null),
    distinctSelectors(registrations),
  );

  expect(unparseable, "recorded selectors the browser cannot parse").toEqual([]);
});

/**
 * Every recorded selector must match at least one element somewhere in the
 * corpus.
 *
 * Without this, the disjointness invariant degrades toward vacuity through
 * exactly the drift the harness exists to catch: if dropdown-menu.gsx starts
 * emitting data-gsxui-dropdown-menuitem while ui/dropdown-menu.js still binds
 * [data-gsxui-dropdown-item], dropdown is dead in production, every
 * registration matches nothing, and the invariant just gets quieter — all
 * 103 tests stay green. A selector that matches nothing is a handler that
 * never runs.
 *
 * A zero-match is not automatically a defect: a component part no example
 * composes will legitimately match nothing. Those are enumerated in
 * allowedUnmatched with a reason each. See that list for why the other two
 * causes (markup drift, dead code) are not allowed to live there.
 */
test("every recorded selector matches an element somewhere in the corpus", async ({
  page,
  registrations,
}) => {
  // 103 navigations in one test — the union is the point, so they cannot be
  // split across parallel per-example tests.
  test.setTimeout(300_000);

  const all = distinctSelectors(registrations);

  // Unparseable selectors would throw inside the match loop below and abort
  // the union. They are reported by the parse test above; drop them here so
  // one bad selector produces one failure, not two.
  const parseable = await page.evaluate((entries) => {
    const ok: number[] = [];
    for (const [i, entry] of entries.entries()) {
      try {
        document.querySelector(entry.selector);
        ok.push(i);
      } catch {
        // reported by "every recorded selector parses"
      }
    }
    return ok;
  }, all);

  let remaining = parseable.map((i) => all[i]);
  for (const example of examples()) {
    if (remaining.length === 0) break;
    await page.goto(example.url);
    const matched = new Set(
      await page.evaluate(
        (selectors) =>
          selectors
            .map((selector, i) => (document.querySelector(selector) ? i : -1))
            .filter((i) => i >= 0),
        remaining.map((entry) => entry.selector),
      ),
    );
    remaining = remaining.filter((_, i) => !matched.has(i));
  }

  const unexpected = remaining.filter(
    (entry) =>
      !allowedUnmatched.some(
        (a) => a.module === entry.module && a.selector === entry.selector,
      ),
  );

  expect(
    unexpected,
    "registered selectors matching no element on any example page — markup drift, " +
      "dead code, or an unexercised part that belongs in allowedUnmatched",
  ).toEqual([]);
});
