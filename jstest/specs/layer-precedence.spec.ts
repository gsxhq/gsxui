import { mkdirSync, writeFileSync } from "node:fs";
import path from "node:path";
import { expect, test } from "@playwright/test";
import { sweepComputedStyles } from "../support/computed-sweep";
import { examples } from "../support/manifest";
import { tmpDir } from "../support/paths";

// This spec is a baseline producer, not an assertion (aside from the pins
// below): `sweep all fixtures` records the resting computed style of every
// data-gsxui-slot-* element across every fixture in the manifest, in both
// colour schemes, and writes the result to disk. Diffing two such runs is
// the actual assertion — that's what `make sweep-compare` does via
// jstest/support/sweep-diff.mjs. See computed-sweep.ts for why: a CSS-layer
// precedence regression changes rendering without failing any assertion
// that isn't looking at computed style directly.
//
// `SWEEP_OUT` selects the output directory; it defaults to a scratch
// location under jstest/.tmp/ so ad hoc `npx playwright test` runs of this
// file don't require the env var. The Makefile targets set it explicitly to
// jstest/.tmp/sweep-baseline or jstest/.tmp/sweep-current.
const outDir = process.env.SWEEP_OUT ?? path.join(tmpDir, "sweep-current");

const schemes = ["light", "dark"] as const;

test.describe.configure({ mode: "serial" });

for (const scheme of schemes) {
  test(`sweep all fixtures (${scheme})`, async ({ page }) => {
    test.setTimeout(300_000);

    if (scheme === "dark") {
      // Runs before every subsequent navigation on this page, so the dark
      // class survives each fixture's page.goto() inside the loop below —
      // an evaluate() after goto would only apply until the next goto.
      await page.addInitScript(() => {
        document.documentElement.classList.add("dark");
      });
    }

    const combined: Record<string, Record<string, string>> = {};
    for (const entry of examples()) {
      const swept = await sweepComputedStyles(page, entry.url);
      for (const [marker, styles] of Object.entries(swept)) {
        // sidebar-menu-skeleton-text renders with a randomized width by
        // design (ui/sidebar.gsx, "The randomized width is the one dynamic
        // presentation value in this part") — it varies run to run
        // regardless of any CSS change, so it would make every compare
        // report a false-positive diff. Not a layer-precedence concern;
        // drop it from the sweep rather than pin a moving target.
        if (marker.startsWith("data-gsxui-slot-sidebar-menu-skeleton-text")) {
          delete styles.width;
        }
        combined[`${entry.component}/${entry.example}::${marker}`] = styles;
      }
    }

    mkdirSync(outDir, { recursive: true });
    writeFileSync(
      path.join(outDir, `sweep-${scheme}.json`),
      JSON.stringify(combined, null, 2),
    );
  });
}

// --- Regression pins -------------------------------------------------------
//
// These pin cases already known to have broken during the Button migration,
// so a future migration can't silently reintroduce them even if nobody
// remembers to diff a baseline. The sweep above is the general net; these
// are pins for two specific fish that already got away once.

test("carousel arrows stay circular", async ({ page }) => {
  await page.goto("/x/carousel/basic");
  const el = page.locator("[data-gsxui-slot-carousel-previous]").first();
  const radius = await el.evaluate((n) => parseFloat(getComputedStyle(n).borderRadius));
  const height = await el.evaluate((n) => n.getBoundingClientRect().height);
  expect(radius).toBeGreaterThanOrEqual(height / 2);
});

test("InputGroupButton keeps its xs type scale and radius", async ({ page }) => {
  await page.goto("/x/input-group/basic");
  const el = page.locator("[data-gsxui-slot-input-group-button]").first();
  await expect(el).toHaveAttribute("data-size", "xs");
  const { fontSize, borderRadius } = await el.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { fontSize: computed.fontSize, borderRadius: computed.borderRadius };
  });
  expect(fontSize).toBe("14px");
  expect(borderRadius).toBe("7px");
});
