/**
 * Reviewed exceptions to the two corpus-wide selector checks.
 *
 * `allowedOverlaps` — selector disjointness (invariants.spec.ts).
 * `allowedUnmatched` — selector coverage (selector-coverage.spec.ts).
 *
 * ---
 *
 * Reviewed exceptions to selector disjointness.
 *
 * An entry says: these two modules may both claim the same element for the
 * same (type, capture) pair, and here is why that is correct. It is a
 * decision someone made and a reviewer can challenge — which is why this is
 * an explicit list rather than a tolerance threshold.
 *
 * Empty is the expected steady state. Adding an entry needs a reason a
 * reviewer would accept, not "the test was failing".
 */
export type AllowedOverlap = {
  /** Module filenames, sorted, e.g. ["dialog.js", "sheet.js"]. */
  modules: [string, string];
  /** "type:capture", e.g. "click:false". */
  key: string;
  reason: string;
};

export const allowedOverlaps: AllowedOverlap[] = [];

/**
 * Reviewed exceptions to selector coverage.
 *
 * An entry says: this module registers a handler on a selector that matches
 * nothing across all 103 example pages, and that is expected — the part is
 * real and correct, no example just happens to compose it.
 *
 * Only that one reason is admissible. The other two things a zero-match can
 * mean are defects, and neither belongs here:
 *
 *   - MARKUP DRIFT — the .gsx stopped emitting the attribute the JS binds.
 *     The component is dead in production. Fix the component.
 *   - DEAD CODE — the part no longer exists at all. Delete the registration.
 *
 * Adding an entry is a claim that a reviewer can check by looking for the
 * part in ui/*.gsx and for its absence in site/examples/. "The test was
 * failing" is not a reason.
 */
export type AllowedUnmatched = {
  /** Module filename, e.g. "navigation-menu.js". */
  module: string;
  /** The recorded selector, verbatim. */
  selector: string;
  reason: string;
};

export const allowedUnmatched: AllowedUnmatched[] = [
  {
    module: "command.js",
    selector: "dialog[data-gsxui-command-dialog]",
    reason:
      "No example composes ui.CommandDialog — site/examples/command/basic.gsx is the " +
      "only command example and it renders an inline palette. The attribute is still " +
      "emitted (ui/command.gsx's CommandDialog, on the DialogContent <dialog>, covered " +
      "by TestCommandDialogComposition), so this is an unexercised part, not drift.",
  },
];
