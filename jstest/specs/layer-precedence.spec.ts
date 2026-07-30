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

test("Card keeps its rounded corners", async ({ page }) => {
  await page.goto("/x/card/compound");
  const el = page.locator("[data-gsxui-slot-card]").first();
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("14px");
});

test("CardHeader keeps its border-b bottom padding and border", async ({ page }) => {
  await page.goto("/x/card/compound");
  const el = page.locator("[data-gsxui-slot-card-header]").first();
  const { borderBottomWidth, paddingBottom } = await el.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { borderBottomWidth: computed.borderBottomWidth, paddingBottom: computed.paddingBottom };
  });
  expect(borderBottomWidth).toBe("1px");
  expect(paddingBottom).toBe("16px");
});

test("Badge keeps its pill corner radius and per-variant background", async ({ page }) => {
  await page.goto("/x/badge/variants");
  const badges = page.locator("[data-gsxui-slot-badge]");
  const defaultBadge = badges.nth(0);
  const secondaryBadge = badges.nth(1);
  await expect(defaultBadge).toHaveAttribute("data-variant", "default");
  await expect(secondaryBadge).toHaveAttribute("data-variant", "secondary");

  const radius = await defaultBadge.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("32px");

  const defaultBg = await defaultBadge.evaluate((n) => getComputedStyle(n).backgroundColor);
  const secondaryBg = await secondaryBadge.evaluate((n) => getComputedStyle(n).backgroundColor);
  expect(defaultBg).not.toBe(secondaryBg);
});

test("Alert destructive variant recolors text but keeps the default's radius", async ({ page }) => {
  await page.goto("/x/alert/variants");
  const alerts = page.locator("[data-gsxui-slot-alert]");
  const defaultAlert = alerts.nth(0);
  const destructiveAlert = alerts.nth(1);
  await expect(defaultAlert).toHaveAttribute("data-variant", "default");
  await expect(destructiveAlert).toHaveAttribute("data-variant", "destructive");

  const defaultColor = await defaultAlert.evaluate((n) => getComputedStyle(n).color);
  const destructiveColor = await destructiveAlert.evaluate((n) => getComputedStyle(n).color);
  expect(defaultColor).not.toBe(destructiveColor);

  const radius = await defaultAlert.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("10px");
});

test("Skeleton keeps its rounded-md corner radius", async ({ page }) => {
  await page.goto("/x/skeleton/basic");
  const el = page.locator("[data-gsxui-slot-skeleton]").first();
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  // rounded-md resolves to --radius-md = calc(var(--radius) - 2px), and this
  // theme's --radius is 0.625rem (10px) — the same token Card's 14px
  // (--radius + 4px) pin already anchors.
  expect(radius).toBe("8px");
});

test("Spinner keeps its spin animation", async ({ page }) => {
  await page.goto("/x/spinner/basic");
  const el = page.locator("[data-gsxui-slot-spinner]").first();
  const { name, duration } = await el.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { name: computed.animationName, duration: computed.animationDuration };
  });
  expect(name).toBe("spin");
  expect(duration).toBe("1s");
});

test("Progress keeps its track height", async ({ page }) => {
  await page.goto("/x/progress/basic");
  const el = page.locator("[data-gsxui-slot-progress]").first();
  const height = await el.evaluate((n) => getComputedStyle(n).height);
  expect(height).toBe("4px");
});

test("Avatar keeps its fixed size-8 dimensions", async ({ page }) => {
  await page.goto("/x/avatar/basic");
  const el = page.locator("[data-gsxui-slot-avatar]").first();
  const { width, height } = await el.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { width: computed.width, height: computed.height };
  });
  expect(width).toBe("32px");
  expect(height).toBe("32px");
});

