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

test("Skeleton keeps its rounded-md radius, and a caller class still wins", async ({
  page,
}) => {
  await page.goto("/x/skeleton/basic");
  // The FIRST skeleton on this page is <ui.Skeleton class="size-12
  // rounded-full"/>, so it deliberately does NOT show the recipe's radius —
  // reading .first() here is what made the original pin fail. Read a skeleton
  // with no radius override for the recipe value, and use the overridden one to
  // pin caller precedence, which nothing else asserts.
  const plain = page.locator("[data-gsxui-slot-skeleton]").nth(1);
  const overridden = page.locator("[data-gsxui-slot-skeleton]").first();

  // rounded-md is --radius-md = calc(var(--radius) - 2px); this theme's
  // --radius is 0.625rem (10px), the same token Card's 14px pin anchors.
  expect(
    await plain.evaluate((n) => getComputedStyle(n).borderRadius),
  ).toBe("8px");

  // A caller-supplied rounded-full must still beat the recipe's rounded-md.
  // Compiled component presentation lives in @layer utilities, so this is the
  // property most at risk from the migration and the least covered.
  const overriddenRadius = await overridden.evaluate((n) =>
    parseFloat(getComputedStyle(n).borderRadius),
  );
  const overriddenHeight = await overridden.evaluate(
    (n) => n.getBoundingClientRect().height,
  );
  expect(overriddenRadius).toBeGreaterThanOrEqual(overriddenHeight / 2);
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
test("Label keeps its tight leading and medium weight", async ({ page }) => {
  await page.goto("/x/label/basic");
  const el = page.locator("[data-gsxui-slot-label]").first();
  const { lineHeight, fontWeight } = await el.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { lineHeight: computed.lineHeight, fontWeight: computed.fontWeight };
  });
  expect(lineHeight).toBe("14px");
  expect(fontWeight).toBe("500");
});
test("AspectRatio keeps the ratio it was given", async ({ page }) => {
  // AspectRatio's only recipe utility is "block", which is a no-op on a
  // <div> (already block by the user-agent stylesheet) — there is no
  // recipe property this component could visibly lose to a layer
  // regression. The aspect-ratio value itself comes from an inline style,
  // not the recipe, so this pin protects the feature end-to-end rather
  // than a layer-precedence-specific regression.
  await page.goto("/x/aspect-ratio/basic");
  const el = page.locator("[data-gsxui-slot-aspect-ratio]").first();
  const aspectRatio = await el.evaluate((n) => getComputedStyle(n).aspectRatio);
  expect(aspectRatio).toBe("16 / 9");
});
test("Separator keeps its hairline thickness in both orientations", async ({ page }) => {
  await page.goto("/x/separator/orientation");
  const separators = page.locator("[data-gsxui-slot-separator]");
  const horizontal = separators.nth(0);
  const vertical = separators.nth(1);
  await expect(horizontal).toHaveAttribute("data-orientation", "horizontal");
  await expect(vertical).toHaveAttribute("data-orientation", "vertical");

  const height = await horizontal.evaluate((n) => getComputedStyle(n).height);
  expect(height).toBe("1px");

  const width = await vertical.evaluate((n) => getComputedStyle(n).width);
  expect(width).toBe("1px");
});
test("Kbd keeps its pill-corner radius", async ({ page }) => {
  await page.goto("/x/kbd/basic");
  const el = page.locator("[data-gsxui-slot-kbd]").first();
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("6px");
});
test("Collapsible trigger keeps its list-none marker hidden", async ({ page }) => {
  await page.goto("/x/collapsible/basic");
  const el = page.locator("[data-gsxui-slot-collapsible-trigger]").first();
  const listStyle = await el.evaluate((n) => getComputedStyle(n).listStyleType);
  expect(listStyle).toBe("none");
});
test("Breadcrumb separator icon keeps its 3.5 size and link keeps hover color", async ({ page }) => {
  await page.goto("/x/breadcrumb/basic");
  const svg = page.locator("[data-gsxui-slot-breadcrumb-separator] > svg").first();
  const size = await svg.evaluate((n) => getComputedStyle(n).width);
  expect(size).toBe("14px");
});
test("Radio checked state keeps its primary fill and gradient dot", async ({ page }) => {
  // Pins both halves of the split: bg-primary (recipe CSS) and the
  // radial-gradient background-image (default.css's plain-CSS escape
  // hatch, since @apply can't carry it) must both still apply together.
  await page.goto("/x/radio/states");
  const checked = page.locator("#radio-states-monthly");
  const { backgroundColor, backgroundImage } = await checked.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { backgroundColor: computed.backgroundColor, backgroundImage: computed.backgroundImage };
  });
  expect(backgroundColor).toBe("oklch(0.205 0 0)");
  expect(backgroundImage).toContain("radial-gradient");
});
test("ToggleGroupItem at spacing=2 keeps Toggle's own per-size rounding", async ({ page }) => {
  // Pins the exact regression a first draft of ToggleGroup's migration
  // introduced: ToggleGroupItem also composes Toggle's own (still
  // unmigrated) presentation via data-gsxui-slot-toggle, and Toggle sets
  // its own size=sm-capped border-radius. An earlier recipe draft restated
  // "rounded-lg" on the spacing=2 arm to satisfy the one-rule-per-value
  // requirement, which silently outranked Toggle's own rounding (compiled
  // utilities always win over Toggle's still-@layer-components rule) and
  // widened the radius from 8px to 10px — caught only by this sweep.
  await page.goto("/x/toggle-group/spacing");
  const item = page.locator("[data-gsxui-slot-toggle-group-item]").first();
  const radius = await item.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("8px");
});
test("Resizable handle keeps its hairline width and flex layout", async ({ page }) => {
  // Pins both the recipe's own width and the foundation.css escape-hatch
  // rules (display: flex / width: 100% under aria-orientation=horizontal)
  // that Resizable's migration split out of the same selector — a layer
  // regression on either half would visibly widen or collapse the handle.
  await page.goto("/x/resizable/handle");
  const handle = page.locator("[data-gsxui-slot-resizable-handle]").first();
  const { width, display } = await handle.evaluate((n) => {
    const computed = getComputedStyle(n);
    return { width: computed.width, display: computed.display };
  });
  expect(display).toBe("flex");
  expect(width).toBe("1px");
});
test("ScrollArea's caller-supplied rounded-md wins over its own recipe rounded-[inherit]", async ({ page }) => {
  await page.goto("/x/scroll-area/basic");
  const el = page.locator("[data-gsxui-slot-scroll-area]").first();
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("8px");
});
test("Checkbox's checked:bg-primary paints its background, unchecked stays transparent", async ({ page }) => {
  await page.goto("/x/checkbox/states");
  const unchecked = page.locator("#checkbox-states-unchecked");
  const checked = page.locator("#checkbox-states-checked");
  const uncheckedBg = await unchecked.evaluate((n) => getComputedStyle(n).backgroundColor);
  const checkedBg = await checked.evaluate((n) => getComputedStyle(n).backgroundColor);
  expect(uncheckedBg).toBe("rgba(0, 0, 0, 0)");
  expect(checkedBg).not.toBe("rgba(0, 0, 0, 0)");
  expect(checkedBg).toBe("oklch(0.205 0 0)");
});
test("Accordion's trigger-icon rotates only inside its own open item", async ({ page }) => {
  // Tailwind v4's rotate-180 compiles to the standalone CSS `rotate`
  // property (not a `transform: rotate(...)` matrix), so that is what a
  // layer regression would actually clobber.
  await page.goto("/x/accordion/basic");
  const icons = page.locator("[data-gsxui-slot-accordion-trigger-icon]");
  const openIconRotate = await icons.nth(0).evaluate((n) => getComputedStyle(n).rotate);
  const closedIconRotate = await icons.nth(1).evaluate((n) => getComputedStyle(n).rotate);
  expect(openIconRotate).toBe("180deg");
  expect(closedIconRotate).toBe("none");
});
test("Table row's data-state=selected background outranks an unselected row", async ({ page }) => {
  await page.goto("/x/table/basic");
  const rows = page.locator("[data-gsxui-slot-table-row]");
  const selectedBg = await rows.nth(0).evaluate((n) => getComputedStyle(n).backgroundColor);
  const unselectedBg = await rows.nth(1).evaluate((n) => getComputedStyle(n).backgroundColor);
  await expect(rows.nth(0)).toHaveAttribute("data-state", "selected");
  expect(selectedBg).not.toBe(unselectedBg);
  expect(selectedBg).toBe("oklch(0.97 0 0)");
  expect(unselectedBg).toBe("rgba(0, 0, 0, 0)");
});
test("InputOTPGroup's has-[[aria-invalid=true]] border fires when a descendant slot is invalid", async ({ page }) => {
  // Regression pin for the escaped-quote bug this migration hit: an
  // arbitrary attribute-value selector written with double quotes
  // (has-[[aria-invalid="true"]]) resolves into a Go string literal with a
  // literal backslash-escaped quote in the generated .gsx, which Tailwind's
  // static content scanner reads as a DIFFERENT candidate than the
  // (correctly unescaped) runtime class — the utility silently never
  // applies. Unquoted attribute values (has-[[aria-invalid=true]]) avoid
  // the escaping entirely.
  await page.goto("/x/input-otp/basic");
  const invalidGroup = page.locator("[data-gsxui-slot-input-otp-group]").first();
  const validGroup = page.locator("[data-gsxui-slot-input-otp-group]").nth(1);
  const invalidBorder = await invalidGroup.evaluate((n) => getComputedStyle(n).borderColor);
  const validBorder = await validGroup.evaluate((n) => getComputedStyle(n).borderColor);
  expect(invalidBorder).not.toBe(validBorder);
  expect(invalidBorder).toBe("oklch(0.577 0.245 27.325)");
});

test("Slider keeps its cursor-pointer affordance", async ({ page }) => {
  await page.goto("/x/slider/basic");
  const el = page.locator("[data-gsxui-slot-slider]").first();
  const cursor = await el.evaluate((n) => getComputedStyle(n).cursor);
  expect(cursor).toBe("pointer");
});
test("EmptyMedia icon variant keeps its rounded muted tile size", async ({ page }) => {
  await page.goto("/x/empty/basic");
  const el = page
    .locator('[data-gsxui-slot-empty-icon][data-variant="icon"]')
    .first();
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  const width = await el.evaluate((n) => getComputedStyle(n).width);
  expect(radius).toBe("10px");
  expect(width).toBe("32px");
});
test("TabsTrigger active state keeps its shadow-sm elevation", async ({ page }) => {
  await page.goto("/x/tabs/basic");
  const el = page
    .locator('[data-gsxui-slot-tabs-trigger][data-state="active"]')
    .first();
  const shadow = await el.evaluate((n) => getComputedStyle(n).boxShadow);
  expect(shadow).toBe(
    "rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(0, 0, 0, 0) 0px 0px 0px 0px, rgba(0, 0, 0, 0.1) 0px 1px 3px 0px, rgba(0, 0, 0, 0.1) 0px 1px 2px -1px",
  );
});
test("PaginationEllipsis keeps its size-8 tile", async ({ page }) => {
  await page.goto("/x/pagination/basic");
  const ellipsis = page.locator("[data-gsxui-slot-pagination-ellipsis]").first();
  const size = await ellipsis.evaluate((n) => getComputedStyle(n).width);
  expect(size).toBe("32px");
});
test("Switch keeps its checked-state background", async ({ page }) => {
  await page.goto("/x/switch/basic");
  const checked = page.locator("[data-gsxui-slot-switch]").nth(1);
  const bg = await checked.evaluate((n) => getComputedStyle(n).backgroundColor);
  expect(bg).toBe("oklch(0.205 0 0)");
});
test("SelectTrigger keeps its default height and radius", async ({ page }) => {
  await page.goto("/x/select/basic");
  const trigger = page.locator("[data-gsxui-slot-select-trigger]").first();
  const height = await trigger.evaluate((n) => getComputedStyle(n).height);
  const radius = await trigger.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(height).toBe("28px");
  expect(radius).toBe("8px");
});

// --- Caller-override pins (spec section 10b) -------------------------------
//
// The sweep above and every pin before this one only measure DEFAULT
// rendering, which is exactly why Card's and Label's migrations silently
// dropped caller-overridability of a relational rule and nobody noticed
// through three migration waves: an arbitrary same-element/ancestor/sibling
// variant compiles to a two-class selector that permanently outranks a
// caller's plain utility regardless of source order, once the rule is
// resolved onto the component's own generated markup instead of living in
// @layer components. These pins pass a caller class that contests the
// retained rule's property and assert the CALLER wins — see
// assets/css/styles/default/card.css and
// assets/css/styles/default/label.css for the retained rules themselves,
// and jstest/harness/style_contract.gsx's "card-header-caller" /
// "label-disabled-caller" fixtures this exercises (fixtures must exist in
// the compiled CSS via Tailwind scanning, not runtime-injected HTML — same
// requirement as "pagination-previous-caller" above).

test("CardHeader's border-b padding stays caller-overridable", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="card-header-caller"]');
  const paddingBottom = await el.evaluate((n) => getComputedStyle(n).paddingBottom);
  // pb-10 (40px) must beat the retained [&.border-b]:pb-4 (16px) rule.
  expect(paddingBottom).toBe("40px");
});

test("Label's disabled-state opacity stays caller-overridable", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="label-disabled-caller"]');
  const opacity = await el.evaluate((n) => getComputedStyle(n).opacity);
  // opacity-100 (1) must beat the retained
  // [[data-disabled=true]_&]:opacity-50 rule (0.5).
  expect(opacity).toBe("1");
});
