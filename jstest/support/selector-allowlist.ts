/**
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
