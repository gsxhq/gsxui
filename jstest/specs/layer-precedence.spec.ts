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

// DropdownMenu/ContextMenu/Menubar migration pins (wave 3b). Two shapes:
// regular per-component regression pins on their own migrated recipe rules,
// and one §10b caller-override pin per component proving the retained
// data-inset/data-disabled rules (assets/css/styles/default/menu.css) still
// lose to a caller's plain utility — see
// jstest/harness/style_contract.gsx's "dropdown-menu-item-caller" /
// "context-menu-item-caller" / "menubar-item-caller" fixtures.

test("DropdownMenu content keeps its rounded-lg popover chrome", async ({ page }) => {
  const response = await page.goto("/x/dropdown-menu/basic");
  expect(response?.status(), "dropdown-menu/basic fixture response").toBe(200);

  const trigger = page.locator("[data-gsxui-dropdown-trigger]").first();
  await trigger.click();
  const content = page.locator("[data-gsxui-slot-dropdown-menu-content]");
  const borderRadius = await content.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(borderRadius).toBe("10px");
});

test("DropdownMenuItem's inset padding and disabled opacity stay caller-overridable", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="dropdown-menu-item-caller"]');
  const style = await el.evaluate((n) => {
    const css = getComputedStyle(n);
    return { paddingLeft: css.paddingLeft, opacity: css.opacity };
  });
  // pl-2 (8px) must beat the retained [data-inset]:pl-8 (32px) rule, and
  // opacity-100 (1) must beat the retained [data-disabled]:opacity-50 (0.5)
  // rule. pointer-events-none for [data-disabled] is baked onto the recipe
  // on purpose and is NOT exercised here (it is meant to stay unoverridable).
  expect(style.paddingLeft).toBe("8px");
  expect(style.opacity).toBe("1");
});

test("ContextMenuItem's inset padding and disabled opacity stay caller-overridable", async ({
  page,
}) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="context-menu-item-caller"]');
  const style = await el.evaluate((n) => {
    const css = getComputedStyle(n);
    return { paddingLeft: css.paddingLeft, opacity: css.opacity };
  });
  expect(style.paddingLeft).toBe("8px");
  expect(style.opacity).toBe("1");
});

test("Menubar keeps its own bar radius, and its item's inset/disabled stay caller-overridable", async ({
  page,
}) => {
  const response = await page.goto("/x/menubar/basic");
  expect(response?.status(), "menubar/basic fixture response").toBe(200);

  const bar = page.locator("[data-gsxui-slot-menubar]");
  const borderRadius = await bar.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(borderRadius).toBe("10px");

  const contractResponse = await page.goto("/f/style-contract");
  expect(contractResponse?.status(), "style contract fixture response").toBe(200);
  const el = page.locator('[data-style-contract="menubar-item-caller"]');
  const style = await el.evaluate((n) => {
    const css = getComputedStyle(n);
    return { paddingLeft: css.paddingLeft, opacity: css.opacity };
  });
  expect(style.paddingLeft).toBe("8px");
  expect(style.opacity).toBe("1");
});

test("NavigationMenu content chrome renders only when its viewport ancestor is false", async ({
  page,
}) => {
  const response = await page.goto("/x/navigation-menu/basic");
  expect(response?.status(), "navigation-menu/basic fixture response").toBe(200);

  const trigger = page.locator("[data-gsxui-navigation-menu-trigger]").first();
  await trigger.click();
  const content = page.locator("[data-gsxui-slot-navigation-menu-content]");
  await expect(content).toBeVisible();
  const borderRadius = await content.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(borderRadius).toBe("10px");
});

test("PopoverContent keeps its rounded corners", async ({ page }) => {
  await page.goto("/x/popover/basic");
  const el = page.locator("[data-gsxui-slot-popover-content]");
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("10px");
});

test("HoverCardContent keeps its rounded corners", async ({ page }) => {
  await page.goto("/x/hover-card/basic");
  const el = page.locator("[data-gsxui-slot-hover-card-content]");
  const radius = await el.evaluate((n) => getComputedStyle(n).borderRadius);
  expect(radius).toBe("10px");
});

test("TooltipContent keeps its own has-kbd padding when it contains a Kbd", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);
  const withKbd = page.locator('[data-style-contract="tooltip-kbd"]');
  const paddingRightWithKbd = await withKbd.evaluate((n) => getComputedStyle(n).paddingRight);
  // has-[[data-gsxui-slot-kbd]]:pr-1.5 (6px) fires because this fixture's
  // TooltipContent actually renders a Kbd child.
  expect(paddingRightWithKbd).toBe("6px");

  await page.goto("/x/tooltip/basic");
  const withoutKbd = page.locator("[data-gsxui-slot-tooltip-content]").first();
  const paddingRightPlain = await withoutKbd.evaluate((n) => getComputedStyle(n).paddingRight);
  // No Kbd child here, so the base px-3 (12px) applies unchanged.
  expect(paddingRightPlain).toBe("12px");
});

// Input's and Textarea's own text-base/md:text-sm pair is retained the same
// way (see assets/css/styles/default/input.css and
// assets/css/styles/default/textarea.css) — both properties stayed off the
// recipe entirely, not just md:text-sm, because a components-layer rule can
// never beat ANY utilities-layer rule: leaving text-base on the recipe while
// retaining only md:text-sm made the retained rule permanently dead instead
// of merely losing to a caller (make sweep-compare caught this on the first
// attempt). These pins run at the default (desktop, >=48rem) viewport, where
// md:text-sm would otherwise apply, and assert the caller's text-lg wins at
// every property in that group.

test("Input's caller text-lg stays overridable against the retained md:text-sm", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="input-caller-text-size"]');
  const fontSize = await el.evaluate((n) => getComputedStyle(n).fontSize);
  // text-lg (18px) must beat the retained md:text-sm (14px) rule.
  expect(fontSize).toBe("18px");
});

test("Textarea's caller text-lg stays overridable against the retained md:text-sm", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="textarea-caller-text-size"]');
  const fontSize = await el.evaluate((n) => getComputedStyle(n).fontSize);
  // text-lg (18px) must beat the retained md:text-sm (14px) rule.
  expect(fontSize).toBe("18px");
});

// Dialog's content radius, and AlertDialogContent's retained max-w-xs caller
// override (section 10b — see assets/css/styles/default/alert-dialog.css).

test("DialogContent keeps its rounded-xl radius and centered box once open", async ({ page }) => {
  await page.goto("/x/dialog/basic");
  const dialog = page.locator("dialog[data-gsxui-dialog-content]");
  await dialog.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(dialog).toHaveJSProperty("open", true);
  const radius = await dialog.evaluate((n) => getComputedStyle(n).borderRadius);
  const position = await dialog.evaluate((n) => getComputedStyle(n).position);
  expect(radius).toBe("14px");
  expect(position).toBe("fixed");
});

test("AlertDialogContent's caller max-w-md stays overridable against its own max-w-xs", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="alert-dialog-content-caller"]');
  await el.evaluate((n) =>
    n.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(el).toHaveJSProperty("open", true);
  const maxWidth = await el.evaluate((n) => getComputedStyle(n).maxWidth);
  // max-w-md must beat AlertDialog's own max-w-xs — this theme's
  // --container-md token computes to 384px (not the 448px max-w-md usually
  // reads as), confirmed against the actual computed style rather than
  // guessed.
  expect(maxWidth).toBe("384px");
});

// AlertDialog's migration restored a narrowing that was silently dead.
// Before it, max-w-xs lived as a retained :where()-wrapped rule at (0,0,0)
// in assets/css/styles/default/alert-dialog.css and lost unconditionally to
// Dialog's own (0,1,0) recipe class in the same layer — measured at a 400px
// viewport on the pre-migration tree, an AlertDialogContent computed
// max-width: calc(100% - 32px), i.e. Dialog's box, not the narrower alert
// box. Carrying it as AlertDialog's own slot utilities hands it to
// merge.Merge instead, which is how upstream settles it too (shadcn
// alert-dialog.tsx puts AlertDialogContent's max-w on its own className and
// lets twMerge displace DialogContent's). Below the sm breakpoint the two
// values differ, which is the only width where this is observable — the
// suite's default 1280px viewport cannot see it, which is why it went
// unnoticed.
test("AlertDialogContent's own max-w-xs narrowing beats Dialog's box below the sm breakpoint", async ({
  page,
}) => {
  await page.setViewportSize({ width: 400, height: 800 });
  await page.goto("/x/alert-dialog/basic");
  const content = page.locator("dialog[data-gsxui-slot-alert-dialog-content]").first();
  await content.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(content).toHaveJSProperty("open", true);
  const maxWidth = await content.evaluate((n) => getComputedStyle(n).maxWidth);
  expect(maxWidth).toBe("320px");
});

// Composer fallout verification: DrawerContent/SheetContent stamp
// data-gsxui-slot-dialog-content without ever calling Dialog's content
// accessor (ui/drawer.gsx, ui/sheet.gsx) — they must not inherit Dialog's
// plain-modal box. Both now carry that chrome on their own recipe classes;
// the marker fallback in assets/css/styles/default.css and
// assets/css/styles/default/drawer-sheet-shared.css are both gone. Verified
// empirically here, not just by reading the CSS.
//
// These pins are deliberately stronger than the fallback-era ones they
// replace: the sweep only records RESTING (closed) computed style, so the
// open-state chrome and every per-side rule are invisible to it. Nothing but
// these assertions would notice a side arm resolving to the wrong class.

test("DrawerContent carries the shared fixed/z-50 chrome on its own recipe class", async ({
  page,
}) => {
  await page.goto("/x/drawer/basic");
  const content = page.locator("dialog[data-gsxui-slot-drawer-content]").first();
  await content.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(content).toHaveJSProperty("open", true);
  const style = await content.evaluate((n) => {
    const s = getComputedStyle(n);
    return {
      position: s.position,
      zIndex: s.zIndex,
      display: s.display,
      margin: [s.marginLeft, s.marginRight, s.marginBottom].join(" "),
      flexDirection: s.flexDirection,
      boxShadow: s.boxShadow,
      transitionDuration: s.transitionDuration,
    };
  });
  expect(style.position).toBe("fixed");
  expect(style.zIndex).toBe("50");
  // display:flex is drawer-sheet-shared.css's [open] rule, now open:flex on
  // the recipe; m-0 and flex-col came from the same block.
  expect(style.display).toBe("flex");
  // m-0 clears the UA <dialog> auto margins; the bottom drawer's own mt-24
  // is the only margin left, so check the other three sides.
  expect(style.margin).toBe("0px 0px 0px");
  expect(style.flexDirection).toBe("column");
  expect(style.boxShadow).not.toBe("none");
  // `transition` sets 150ms and `duration-200` overrides it — the one place
  // the recipe's @apply ORDER is load-bearing.
  expect(style.transitionDuration).toBe("0.2s");
});

test("DrawerContent's four side arms each resolve to their own anchoring", async ({ page }) => {
  await page.goto("/x/drawer/directions");
  const expected = {
    bottom: { borderTopWidth: "1px", borderBottomWidth: "0px", marginTop: "96px", textAlign: "center" },
    top: { borderTopWidth: "0px", borderBottomWidth: "1px", marginBottom: "96px", textAlign: "center" },
    left: { borderRightWidth: "1px", borderLeftWidth: "0px", textAlign: "left" },
    right: { borderLeftWidth: "1px", borderRightWidth: "0px", textAlign: "left" },
  } as const;
  // left/right expect "left" rather than the inherited "start": the suite
  // viewport is 1280px, above the md breakpoint the side arms' own
  // md:[&_[data-gsxui-slot-drawer-header]]:text-left fires at.
  for (const [side, want] of Object.entries(expected)) {
    const content = page.locator(`dialog[data-gsxui-slot-drawer-content][data-side="${side}"]`);
    await content.evaluate((el) =>
      el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
    );
    await expect(content).toHaveJSProperty("open", true);
    const got = await content.evaluate((n) => {
      const s = getComputedStyle(n);
      const header = n.querySelector("[data-gsxui-slot-drawer-header]");
      return {
        borderTopWidth: s.borderTopWidth,
        borderBottomWidth: s.borderBottomWidth,
        borderLeftWidth: s.borderLeftWidth,
        borderRightWidth: s.borderRightWidth,
        marginTop: s.marginTop,
        marginBottom: s.marginBottom,
        textAlign: header ? getComputedStyle(header).textAlign : "",
      };
    });
    for (const [property, value] of Object.entries(want)) {
      expect(got[property as keyof typeof got], `${side} drawer ${property}`).toBe(value);
    }
    await content.evaluate((el) =>
      el.dispatchEvent(new CustomEvent("gsxui:request-close", { bubbles: true, cancelable: true })),
    );
  }
});

test("SheetContent carries the shared fixed/z-50 chrome on its own recipe class", async ({
  page,
}) => {
  await page.goto("/x/sheet/basic");
  const content = page.locator("dialog[data-gsxui-slot-sheet-content]").first();
  await content.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(content).toHaveJSProperty("open", true);
  const style = await content.evaluate((n) => {
    const s = getComputedStyle(n);
    const closeButton = n.querySelector("[data-gsxui-slot-sheet-close-button]");
    return {
      position: s.position,
      zIndex: s.zIndex,
      display: s.display,
      margin: s.margin,
      flexDirection: s.flexDirection,
      boxShadow: s.boxShadow,
      transitionDuration: s.transitionDuration,
      closePosition: closeButton ? getComputedStyle(closeButton).position : "",
    };
  });
  expect(style.position).toBe("fixed");
  expect(style.zIndex).toBe("50");
  // display:flex is drawer-sheet-shared.css's [open] rule, now open:flex on
  // the recipe; m-0 and flex-col came from the same block. Sheet sets no
  // per-side margin, so all four sides are 0.
  expect(style.display).toBe("flex");
  expect(style.margin).toBe("0px");
  expect(style.flexDirection).toBe("column");
  expect(style.boxShadow).not.toBe("none");
  // `transition` sets 150ms and `duration-200` overrides it — the one place
  // the recipe's @apply ORDER is load-bearing.
  expect(style.transitionDuration).toBe("0.2s");
  // The injected close button migrated too and must keep its own absolute
  // placement.
  expect(style.closePosition).toBe("absolute");
});

test("SheetContent's four side arms each resolve to their own anchoring", async ({ page }) => {
  await page.goto("/x/sheet/directions");
  const expected = {
    right: { borderLeftWidth: "1px", borderRightWidth: "0px" },
    left: { borderRightWidth: "1px", borderLeftWidth: "0px" },
    top: { borderBottomWidth: "1px", borderTopWidth: "0px" },
    bottom: { borderTopWidth: "1px", borderBottomWidth: "0px" },
  } as const;
  for (const [side, want] of Object.entries(expected)) {
    const content = page.locator(`dialog[data-gsxui-slot-sheet-content][data-side="${side}"]`);
    await content.evaluate((el) =>
      el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
    );
    await expect(content).toHaveJSProperty("open", true);
    const got = await content.evaluate((n) => {
      const s = getComputedStyle(n);
      return {
        borderTopWidth: s.borderTopWidth,
        borderBottomWidth: s.borderBottomWidth,
        borderLeftWidth: s.borderLeftWidth,
        borderRightWidth: s.borderRightWidth,
      };
    });
    for (const [property, value] of Object.entries(want)) {
      expect(got[property as keyof typeof got], `${side} sheet ${property}`).toBe(value);
    }
    await content.evaluate((el) =>
      el.dispatchEvent(new CustomEvent("gsxui:request-close", { bubbles: true, cancelable: true })),
    );
  }
});

// Sheet's migration fallout on Sidebar, which composes <ui.SheetContent> and
// <ui.SheetHeader> directly: Sheet's own utilities beat anything Sidebar
// declares for the same property from @layer components. make sweep-compare
// caught the mobile panel collapsing from its --sidebar-width to Sheet's
// w-3/4 and repainting with bg-background.
//
// The @layer utilities promotion that first restored this has since been
// RETIRED by Sidebar's own migration: the same three properties are now
// utilities on the sidebar-mobile-content and -mobile-header slots, so they
// ride into SheetContent's / SheetHeader's own class attribute and
// merge.Merge settles them. sm:max-w-none is load-bearing there — merge
// scopes conflicts per variant, so an unprefixed max-w-none would leave
// Sheet's sm:max-w-sm standing. This pin is what holds all of that.
test("Sidebar's mobile panel keeps its own width and surface against Sheet's", async ({ page }) => {
  await page.setViewportSize({ width: 640, height: 900 });
  await page.goto("/x/sidebar/basic");
  const content = page.locator("dialog[data-gsxui-slot-sidebar-mobile-content]").first();
  await content.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(content).toHaveJSProperty("open", true);
  const got = await content.evaluate((n) => {
    const s = getComputedStyle(n);
    const header = n.querySelector("[data-gsxui-slot-sidebar-mobile-header]");
    return {
      width: s.width,
      backgroundColor: s.backgroundColor,
      headerPadding: header ? getComputedStyle(header).padding : "",
    };
  });
  expect(got.width).toBe("288px");
  expect(got.backgroundColor).not.toBe("rgba(0, 0, 0, 0)");
  // sr-only must survive SheetHeader's p-4.
  expect(got.headerPadding).toBe("0px");
});

// The pre-slot-axis menu CSS muted item icons through a :where() selector list
// that named only plain items and sub-triggers, so CheckboxItem's and
// RadioItem's icons rendered at full foreground. Upstream carries the same
// [&_svg:not([class*='text-'])]:text-muted-foreground on those two items
// (shadcn context-menu.tsx), and the recipe now does too — make sweep-compare
// surfaced the change as 9 elements shifting to the muted token in both themes.

test("Command keeps its palette chrome and stays caller-overridable", async ({
  page,
}) => {
  const response = await page.goto("/x/command/basic");
  expect(response?.status(), "command/basic fixture response").toBe(200);

  const root = page.locator("[data-gsxui-slot-command]").first();
  const rootStyle = await root.evaluate((n) => {
    const css = getComputedStyle(n);
    return { borderRadius: css.borderRadius, backgroundColor: css.backgroundColor };
  });
  // site/examples/command/basic.gsx passes class="max-w-md rounded-lg …": the
  // caller's rounded-lg must replace the recipe's own rounded-xl (10px, not
  // this theme's 14px), and the recipe's bg-popover must still apply.
  expect(rootStyle.borderRadius).toBe("10px");
  expect(rootStyle.backgroundColor).not.toBe("rgba(0, 0, 0, 0)");

  // The input wrapper's `> svg` search icon: a child-combinator rule that
  // became a [&>svg]: arbitrary variant on the wrapper's recipe class.
  const icon = page.locator("[data-gsxui-slot-command-input-wrapper] > svg").first();
  const iconSize = await icon.evaluate((n) => getComputedStyle(n).width);
  expect(iconSize).toBe("16px");

  const list = page.locator("[data-gsxui-slot-command-list]").first();
  expect(await list.evaluate((n) => getComputedStyle(n).maxHeight)).toBe("288px");
});

test("CommandItem's disabled opacity stays caller-overridable", async ({ page }) => {
  const response = await page.goto("/f/style-contract");
  expect(response?.status(), "style contract fixture response").toBe(200);

  const el = page.locator('[data-style-contract="command-item-caller"]');
  const style = await el.evaluate((n) => {
    const css = getComputedStyle(n);
    return { opacity: css.opacity, pointerEvents: css.pointerEvents };
  });
  // opacity-100 must beat the retained [data-disabled=true]:opacity-50 rule;
  // pointer-events-none is baked onto the recipe on purpose and must NOT be
  // overridable.
  expect(style.opacity).toBe("1");
  expect(style.pointerEvents).toBe("none");
});

test("ContextMenu checkbox and radio items mute their icons like plain items", async ({ page }) => {
  await page.goto("/x/context-menu/full");
  const muted = "oklch(0.556 0 0)";
  for (const slot of ["checkbox-item", "radio-item"]) {
    const icon = page.locator(`[data-gsxui-slot-context-menu-${slot}] svg`).first();
    const color = await icon.evaluate((n) => getComputedStyle(n).color);
    expect(color, `${slot} icon colour`).toBe(muted);
  }
});

// NativeSelect's migration (Wave 4d). Two things have to keep holding: the
// component's own chrome on a standalone select, and Calendar's override of
// that chrome for the month/year dropdowns it composes. The latter is the
// transitive-fallout shape — Calendar's rules moved to @layer utilities in
// assets/css/styles/default.css because NativeSelect's own resolved
// utilities would otherwise outrank them.

test("NativeSelect keeps its own height, radius and chevron placement", async ({ page }) => {
  await page.goto("/x/native-select/basic");
  const select = page.locator("[data-gsxui-slot-native-select]").first();
  expect(await select.evaluate((n) => getComputedStyle(n).height)).toBe("32px");
  expect(await select.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("10px");
  const chevron = page.locator("[data-gsxui-slot-native-select-wrapper] > svg").first();
  expect(await chevron.evaluate((n) => getComputedStyle(n).position)).toBe("absolute");
  expect(await chevron.evaluate((n) => getComputedStyle(n).width)).toBe("16px");
});

test("Calendar's dropdowns still override NativeSelect's own chrome", async ({ page }) => {
  await page.goto("/x/calendar/dropdown");
  const wrapper = page.locator("[data-gsxui-slot-calendar-dropdowns] [data-gsxui-slot-native-select-wrapper]").first();
  const select = wrapper.locator("> [data-gsxui-slot-native-select]");
  const chevron = wrapper.locator("> svg");
  expect(await wrapper.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("8px");
  expect(await select.evaluate((n) => getComputedStyle(n).height)).toBe("24px");
  expect(await select.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("0px");
  expect(await select.evaluate((n) => getComputedStyle(n).borderTopWidth)).toBe("0px");
  expect(await select.evaluate((n) => getComputedStyle(n).paddingLeft)).toBe("6px");
  expect(await chevron.evaluate((n) => getComputedStyle(n).width)).toBe("14px");
});

// Carousel's migration (Wave 4d). The arrows compose <Button variant="outline"
// size="icon">, and their round shape used to need the @layer utilities escape
// hatch in assets/css/styles/default.css to beat Button's own rounded-lg. On
// the recipe it is a plain utility merged onto Button's own class list, so
// merge.Merge drops rounded-lg — this pins that the arrows are still round and
// still positioned, and that the caller's spacing override still wins on the
// track/item, which is the property pair site/examples/carousel/sizes.gsx
// overrides.

test("Carousel arrows keep their round shape and outboard positioning", async ({ page }) => {
  await page.goto("/x/carousel/basic");
  const prev = page.locator("[data-gsxui-slot-carousel-previous]").first();
  expect(await prev.evaluate((n) => getComputedStyle(n).position)).toBe("absolute");
  expect(await prev.evaluate((n) => getComputedStyle(n).width)).toBe("32px");
  const radius = await prev.evaluate((n) => parseFloat(getComputedStyle(n).borderTopLeftRadius));
  expect(radius, "rounded-full must beat Button's rounded-lg").toBeGreaterThan(100);
});

test("Carousel's caller spacing override still displaces the track and item defaults", async ({ page }) => {
  await page.goto("/x/carousel/sizes");
  const track = page.locator("[data-gsxui-slot-carousel-track]").first();
  const item = page.locator("[data-gsxui-slot-carousel-item]").first();
  // -ml-1 / pl-1 / -scroll-ml-1, not the recipe's own -ml-4 / pl-4.
  expect(await track.evaluate((n) => getComputedStyle(n).marginLeft)).toBe("-4px");
  expect(await item.evaluate((n) => getComputedStyle(n).paddingLeft)).toBe("4px");
  expect(await item.evaluate((n) => getComputedStyle(n).scrollMarginLeft)).toBe("-4px");
});

// Calendar's migration (Wave 4d). CONFORMANCE FIX, pinned so it cannot
// silently revert: a range-middle day that sits at a week-row boundary now
// rounds its outer corner, where before it stayed square. The old
// default.css rule wrapped the whole selector in :where()
// (`:where([...day]:last-child[data-selected=true]) > :where([...day-button])`,
// (0,0,0)), so the day button's own [data-range-middle="true"] rule at
// (0,1,0) beat it and rounded-none won. Upstream shadcn puts the same rule on
// the DAY CELL's class list
// (`[&:last-child[data-selected=true]_button]:rounded-r-md`, calendar.tsx
// line 107), which outranks its own data-[range-middle=true]:rounded-none —
// so upstream caps a range at the row edge and gsxui now matches. make
// sweep-compare surfaced this as four differences (two buttons x two themes)
// on the calendar/loadedrange fixture.

test("Calendar caps a range at the week-row boundary, like upstream", async ({ page }) => {
  await page.goto("/x/calendar/loadedrange");
  const rowEnd = page.locator(
    '[data-gsxui-slot-calendar-day]:last-child[data-selected="true"] > [data-gsxui-slot-calendar-day-button][data-range-middle="true"]',
  );
  const rowStart = page.locator(
    '[data-gsxui-slot-calendar-day]:first-child[data-selected="true"] > [data-gsxui-slot-calendar-day-button][data-range-middle="true"]',
  );
  await expect(rowEnd).toHaveCount(1);
  await expect(rowStart).toHaveCount(1);
  expect(await rowEnd.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("0px 8px 8px 0px");
  expect(await rowStart.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("8px 0px 0px 8px");
});

test("Calendar keeps its own chrome and cell geometry", async ({ page }) => {
  await page.goto("/x/calendar/basic");
  const root = page.locator("[data-gsxui-slot-calendar]").first();
  const nav = page.locator("[data-gsxui-slot-calendar-nav]").first();
  const navButton = page.locator("[data-gsxui-slot-calendar-nav-button]").first();
  expect(await root.evaluate((n) => getComputedStyle(n).padding)).toBe("8px");
  // --cell-size is calc(var(--spacing) * 7) == 28px, carried as an arbitrary
  // property on the root and consumed by the nav/day geometry.
  expect(await root.evaluate((n) => getComputedStyle(n).getPropertyValue("--cell-size").trim())).toBe(
    "calc(0.25rem * 7)",
  );
  expect(await nav.evaluate((n) => getComputedStyle(n).position)).toBe("absolute");
  expect(await navButton.evaluate((n) => getComputedStyle(n).width)).toBe("28px");
  expect(await navButton.evaluate((n) => getComputedStyle(n).height)).toBe("28px");
});

test("Calendar inside a Popover keeps its transparent background", async ({ page }) => {
  await page.goto("/x/calendar/datepicker");
  await page.locator("[data-gsxui-slot-popover-trigger]").first().click();
  const calendar = page.locator("[data-gsxui-slot-popover-content] [data-gsxui-slot-calendar]").first();
  expect(await calendar.evaluate((n) => getComputedStyle(n).backgroundColor)).toBe("rgba(0, 0, 0, 0)");
});


// ButtonGroup's migration (Wave 4c). Three things it must not lose:
//
//  1. ButtonGroupSeparator's bg-input/h-auto. Before the migration these
//     lived UNWRAPPED in assets/css/styles/default.css's @layer utilities
//     escape hatch, promoted there so they could outrank Separator's own
//     migrated bg-border/h-full on the same element. They now ride the
//     ordinary caller-class merge onto <Separator>, where tailwind-merge
//     drops bg-border outright and data-[orientation=vertical]:h-auto
//     outranks h-full by specificity. Same rendering, different mechanism —
//     which is exactly the kind of change that regresses silently.
//  2. The inner/outer corner collapse, which still lives in default.css's
//     @layer utilities block and paints ButtonGroup's CHILDREN. It has to
//     keep beating Button's own migrated rounded-lg.
//  3. The nested-group gap-2, RETAINED in @layer components under § 10b so a
//     caller's own gap utility still wins the layer contest.

test("ButtonGroupSeparator keeps bg-input and its vertical h-auto after migration", async ({ page }) => {
  await page.goto("/x/button-group/basic");
  const vertical = page
    .locator('[data-gsxui-slot-button-group-separator][data-orientation="vertical"]')
    .first();
  // bg-input, not Separator's own bg-border.
  expect(await vertical.evaluate((n) => getComputedStyle(n).backgroundColor)).toBe(
    "oklch(0.922 0 0)",
  );
  // h-auto + self-stretch: the separator stretches to the group, rather than
  // collapsing to Separator's own h-full against an auto-height parent.
  expect(await vertical.evaluate((n) => getComputedStyle(n).width)).toBe("1px");
  expect(await vertical.evaluate((n) => n.getBoundingClientRect().height > 0)).toBe(true);
});

test("ButtonGroup still collapses its children's inner corners over Button's own radius", async ({
  page,
}) => {
  await page.goto("/x/button-group/basic");
  const group = page.locator("[data-gsxui-slot-button-group]").first();
  const first = group.locator("[data-gsxui-slot-button]").first();
  const last = group.locator("[data-gsxui-slot-button]").last();
  const radius = (locator: typeof first) =>
    locator.evaluate((n) => {
      const style = getComputedStyle(n);
      return [
        style.borderTopLeftRadius,
        style.borderTopRightRadius,
        style.borderBottomRightRadius,
        style.borderBottomLeftRadius,
      ].join(" ");
    });
  // Outer corners keep Button's own radius (10px at this size ramp; read from
  // the real computed style, not guessed from the utility name); the seam
  // between them is square.
  expect(await radius(first)).toBe("10px 0px 0px 10px");
  expect(await radius(last)).toBe("0px 10px 10px 0px");
});

test("ButtonGroup's nested-group gap stays caller-overridable", async ({ page }) => {
  await page.goto("/f/style-contract");
  // The retained rule is :has(> [data-gsxui-slot-button-group]) { gap-2 }, in
  // @layer components. A caller's plain gap utility sits in @layer utilities
  // and must win on the layer boundary — which it only can while the rule
  // stays behind. As has-[...]:gap-2 on the recipe it would compile to a
  // (0,2,0) selector and win unconditionally.
  const el = page.locator('[data-style-contract="button-group-nested-caller-gap"]');
  expect(await el.evaluate((n) => getComputedStyle(n).gap)).toBe("32px");
});

// Item's migration (Wave 4c).
//
// The media offset — `item:has(item-description) item-media { translate-y-0.5
// self-start }` — is RETAINED in @layer components under § 10b: as an
// ancestor-scoped arbitrary variant on the media slot's recipe it would compile
// to (0,3,0) and a caller's own self-* could never beat it. These two pins
// assert both halves of that: the rule still applies, AND a caller still wins.

test("Item's media offset applies through the retained rule", async ({ page }) => {
  await page.goto("/f/style-contract");
  const media = page.locator('[data-style-contract="item-media-offset"]');
  expect(await media.evaluate((n) => getComputedStyle(n).alignSelf)).toBe("flex-start");
  expect(await media.evaluate((n) => getComputedStyle(n).translate)).toBe("0px 2px");
});

test("Item's media offset stays caller-overridable", async ({ page }) => {
  await page.goto("/f/style-contract");
  const media = page.locator('[data-style-contract="item-media-offset-caller"]');
  // self-center is a plain (0,1,0) utility in @layer utilities; it can only
  // beat the retained rule while that rule stays in @layer components.
  expect(await media.evaluate((n) => getComputedStyle(n).alignSelf)).toBe("center");
});

// The content flex group DID migrate, in full (see registry/styles/nova/item.css's
// header). This pins the pair that forced it: the first content grows, the
// second — selected by the sibling variant on the first one's own class —
// does not.
test("Item's second adjacent content does not grow", async ({ page }) => {
  await page.goto("/f/style-contract");
  const first = page.locator('[data-style-contract="item-content-first"]');
  const second = page.locator('[data-style-contract="item-content-second"]');
  expect(await first.evaluate((n) => getComputedStyle(n).flexGrow)).toBe("1");
  expect(await second.evaluate((n) => getComputedStyle(n).flexGrow)).toBe("0");
});

// ItemSeparator's my-2 now arrives as an ordinary caller class on <Separator>
// rather than as a marker-keyed default.css rule.
test("ItemSeparator keeps its my-2 margin after migration", async ({ page }) => {
  await page.goto("/x/item/basic");
  const separator = page.locator("[data-gsxui-slot-item-separator]").first();
  expect(await separator.evaluate((n) => getComputedStyle(n).marginTop)).toBe("8px");
  expect(await separator.evaluate((n) => getComputedStyle(n).marginBottom)).toBe("8px");
});

// Field's migration (Wave 4c). Field is the most relational component in the
// catalogue, so these pins cover both sides of the § 10b split.

test("Field's responsive orientation still resolves its container query", async ({ page }) => {
  await page.goto("/x/field/basic");
  // The @container field-group (width >= 28rem) block became Tailwind's
  // @min-[28rem]/field-group: variant on the responsive arm. The container
  // itself is still declared in plain CSS in default.css's escape hatch — this
  // pin is what proves Tailwind can query a container it does not own.
  const field = page.locator('[data-gsxui-slot-field][data-orientation="responsive"]').first();
  expect(await field.evaluate((n) => getComputedStyle(n).flexDirection)).toBe("row");
  // flex-start, not center: this fixture's responsive field contains a
  // FieldContent, so the container query's own
  // has-[>[data-gsxui-slot-field-content]]:items-start arm wins over its
  // items-center. Both live inside the same @min-[28rem]/field-group: block,
  // so this reads the container query twice over.
  expect(await field.evaluate((n) => getComputedStyle(n).alignItems)).toBe("flex-start");
});

test("Field's disabled dimming applies through the retained rule, and stays caller-overridable", async ({
  page,
}) => {
  await page.goto("/f/style-contract");
  const dimmed = page.locator('[data-style-contract="field-title-disabled"]');
  const overridden = page.locator('[data-style-contract="field-title-disabled-caller"]');
  expect(await dimmed.evaluate((n) => getComputedStyle(n).opacity)).toBe("0.5");
  // opacity-100 is a plain (0,1,0) utility; it can only beat the ancestor-scoped
  // rule while that rule stays in @layer components.
  expect(await overridden.evaluate((n) => getComputedStyle(n).opacity)).toBe("1");
});

test("Field's invalid tint applies from the recipe", async ({ page }) => {
  await page.goto("/x/field/invalid");
  // data-[invalid=true]:text-destructive moved ONTO the recipe (the layer gate
  // rejected retaining it — see assets/css/styles/default/field.css), so this
  // pin is what catches the variant silently failing to compile.
  const field = page.locator('[data-gsxui-slot-field][data-invalid="true"]').first();
  expect(await field.evaluate((n) => getComputedStyle(n).color)).toBe("oklch(0.577 0.245 27.325)");
});

test("FieldLabel's leading-snug beats Label's own leading-none", async ({ page }) => {
  await page.goto("/x/field/basic");
  // Before Field migrated this needed an unwrapped @layer utilities promotion
  // in default.css. It is now an ordinary tailwind-merge outcome on the
  // composed <Label>; the promotion is retired.
  const label = page.locator("[data-gsxui-slot-field-label]").first();
  const fontSize = await label.evaluate((n) => parseFloat(getComputedStyle(n).fontSize));
  const lineHeight = await label.evaluate((n) => parseFloat(getComputedStyle(n).lineHeight));
  // leading-snug is 1.375em, leading-none is 1em.
  expect(lineHeight / fontSize).toBeCloseTo(1.375, 3);
});

test("FieldSeparator keeps its outline-group bottom margin", async ({ page }) => {
  await page.goto("/x/field/basic");
  // group-data-[variant=outline] reaching down onto the wrapper, translated as
  // an ancestor-scoped arbitrary variant on the wrapper's own recipe class.
  const wrapper = page.locator("[data-gsxui-slot-field-separator-wrapper]").first();
  expect(await wrapper.evaluate((n) => getComputedStyle(n).marginBottom)).toBe("-8px");
});

test("A nested FieldGroup tightens its gap", async ({ page }) => {
  await page.goto("/f/style-contract");
  const outer = page.locator('[data-style-contract="field-group-outer"]');
  const nested = page.locator('[data-style-contract="field-group-nested"]');
  expect(await outer.evaluate((n) => getComputedStyle(n).rowGap)).toBe("20px");
  expect(await nested.evaluate((n) => getComputedStyle(n).rowGap)).toBe("16px");
});

test("FieldDescription's retained margin-top group survives", async ({ page }) => {
  await page.goto("/f/style-contract");
  const description = page.locator('[data-style-contract="field-description-after-legend"]');
  // The legend-sibling rule and the :last-child rule are ONE property group and
  // both stayed in @layer components; the sibling rule outranks :last-child
  // inside that layer.
  expect(await description.evaluate((n) => getComputedStyle(n).marginTop)).toBe("-6px");
});

// InputGroup's migration (Wave 4c).

test("InputGroup keeps its own chrome and its block-align column layout", async ({ page }) => {
  await page.goto("/x/input-group/basic");
  const group = page.locator("[data-gsxui-slot-input-group]").first();
  expect(await group.evaluate((n) => getComputedStyle(n).borderRadius)).toBe("10px");
  // The whole has-() state block moved onto the recipe (every rule in it
  // contested a property the base rule sets — h-8, border-input, the two ring
  // colours), so these are the pin that catches the variants failing to
  // compile.
  const stacked = page.locator('[data-gsxui-slot-input-group]:has(> [data-align="block-start"])').first();
  expect(await stacked.evaluate((n) => getComputedStyle(n).flexDirection)).toBe("column");
  expect(await stacked.evaluate((n) => getComputedStyle(n).height !== "32px")).toBe(true);
});

test("InputGroup's invalid state borders red", async ({ page }) => {
  await page.goto("/x/input-group/basic");
  const invalid = page.locator('[data-gsxui-slot-input-group]:has([aria-invalid="true"])').first();
  expect(await invalid.evaluate((n) => getComputedStyle(n).borderTopColor)).toBe(
    "oklch(0.577 0.245 27.325)",
  );
});

test("InputGroup's addon border-b padding applies", async ({ page }) => {
  await page.goto("/f/style-contract");
  // `[data-align=block-start].border-b -> pb-2` is § 10b's own worst case
  // (Card's [&.border-b]:pb-4), and it is on the recipe anyway: retaining it in
  // @layer components left it DEAD against the addon's own py-1.5, because
  // padding-bottom is a sub-property of the py-* shorthand. That reads as 6px
  // here, not 8px — this pin is what caught it, since --check-layers does not
  // model py-* vs pb-* and no fixture passes border-b.
  //
  // The second element additionally passes pb-8 and DOES NOT win: at (0,2,0)
  // the variant outranks a caller's plain pb-*. That is the § 10b cost, taken
  // deliberately because the alternative is the rule never applying. Upstream
  // shadcn has the same behaviour.
  const plain = page.locator('[data-style-contract="input-group-addon-border-b"]');
  const overridden = page.locator('[data-style-contract="input-group-addon-border-b-caller"]');
  expect(await plain.evaluate((n) => getComputedStyle(n).paddingBottom)).toBe("8px");
  expect(await overridden.evaluate((n) => getComputedStyle(n).paddingBottom)).toBe("8px");
});

// CONFORMANCE FIX, not a regression. Before InputGroup migrated, its control's
// `focus-visible:ring-0` and `dark:bg-transparent` sat in @layer components and
// lost to Input's/Textarea's own utilities-layer classes, so a focused
// InputGroupInput drew its OWN 3px ring inside the group's ring. They are now
// ordinary utilities merged onto the control by tailwind-merge, which drops
// Input's focus-visible:ring-[3px] and dark:bg-input/30 outright — exactly what
// upstream shadcn's cn() does with the identical class list
// (registry/new-york-v4/ui/input-group.tsx's InputGroupInput). The resting
// computed-style sweep cannot see a focus state, so this pin is the only guard.
test("InputGroupInput draws no ring of its own when focused", async ({ page }) => {
  await page.goto("/x/input-group/basic");
  const input = page.locator("input[data-gsxui-slot-input-group-control]").first();
  await input.focus();
  expect(await input.evaluate((n) => getComputedStyle(n).outlineStyle)).toBe("none");
  // Input's own focus-visible:ring-[3px] is gone; the group draws the ring.
  expect(await input.evaluate((n) => n.className.includes("focus-visible:ring-[3px]"))).toBe(false);
  const group = page.locator("[data-gsxui-slot-input-group]").first();
  expect(await group.evaluate((n) => getComputedStyle(n).boxShadow)).not.toBe("none");
});

test("Combobox's InputGroup keeps its w-auto default and still yields to a caller width", async ({
  page,
}) => {
  // InputGroup's migration killed combobox.css's marker-keyed w-auto rule
  // (a components-layer rule cannot beat InputGroup's own w-full utility), and
  // promoting it to @layer utilities overshot and beat the caller's own
  // w-[220px]. Combobox's own migration made it an ordinary recipe utility on
  // its input-group slot — upstream's own cn("w-auto", className) form.
  await page.goto("/x/combobox/basic");
  const sized = page.locator("[data-gsxui-slot-combobox-input-group]").first();
  expect(await sized.evaluate((n) => getComputedStyle(n).width)).toBe("220px");

  await page.goto("/x/combobox/form");
  const unsized = page.locator("[data-gsxui-slot-combobox-input-group]").first();
  expect(await unsized.evaluate((n) => n.className.includes("w-full"))).toBe(false);
  expect(await unsized.evaluate((n) => n.className.includes("w-auto"))).toBe(true);
});

// Combobox's trigger is hidden whenever a clear button is also present —
// ComboboxInput's own documented contract ("clear wins visually"), matching
// shadcn's composition table. Until Combobox migrated, the rule that enforces
// it was `:where([data-gsxui-slot-input-group]):has([…combobox-clear])
// :where([…combobox-trigger]) { display: none }` at (0,0,0), while the trigger
// renders through InputGroupButton -> Button, whose migrated recipe emits a
// real `inline-flex` at (0,1,0) in the same layer — so the trigger was NOT
// hidden, and no fixture rendered both buttons for the sweep to notice
// (combobox/clear passes showClear without showTrigger). site/examples/
// combobox/trigger-clear.gsx now does, and the migrated form is an arbitrary
// variant on the trigger's own recipe class at (0,3,0).
test("Combobox hides its trigger when a clear button is present", async ({ page }) => {
  await page.goto("/x/combobox/trigger-clear");
  const trigger = page.locator("[data-gsxui-slot-combobox-trigger]").first();
  expect(await trigger.evaluate((n) => getComputedStyle(n).display)).toBe("none");
  // The clear button is the one that stays. Its resting display is `flex`,
  // not Button's own inline-flex: InputGroupButton's own recipe restates
  // display for every button inside an input group.
  const clear = page.locator("[data-gsxui-slot-combobox-clear]").first();
  expect(await clear.evaluate((n) => getComputedStyle(n).display)).toBe("flex");
});

// Conformance fix from Combobox's migration. Upstream's ComboboxTrigger renders
// its chevron with an explicit `className="pointer-events-none size-4
// text-muted-foreground"` (registry/new-york-v4/ui/combobox.tsx:36), and
// upstream's trigger itself is `cn("[&_svg:not([class*='size-'])]:size-4",
// className)` on InputGroupButton — 16px both ways. gsxui rendered 12px before
// the migration: the icon's own marker rule was a (0,0,0) @layer components
// selector and lost to Button's icon-xs `[&_svg:not([class*='size-'])]:size-3`
// compiled utility. Both are ordinary utilities now, so tailwind-merge drops
// the size-3 form and the icon carries size-4 outright. make sweep-compare
// caught this as 12px -> 16px on four fixtures.
test("Combobox's trigger chevron is size-4, not Button's icon-xs size-3", async ({ page }) => {
  await page.goto("/x/combobox/basic");
  const icon = page.locator("[data-gsxui-slot-combobox-trigger-icon]").first();
  expect(await icon.evaluate((n) => getComputedStyle(n).width)).toBe("16px");
  expect(await icon.evaluate((n) => getComputedStyle(n).height)).toBe("16px");
});


// ---------------------------------------------------------------------------
// Sidebar's own migration to the slot axis.
//
// Sidebar composes four already-migrated components — Button, Input,
// Separator and Skeleton — and before this migration each contest was settled
// by an unwrapped rule promoted into assets/css/styles/default.css's
// @layer utilities block. All three promotions are retired: Sidebar's slot
// utilities now travel in the SAME class attribute as the composed
// component's, where merge.Merge decides. These pins are what stops that
// quietly reverting to the composed component's value.
// ---------------------------------------------------------------------------

test("SidebarTrigger stays 28px against Button's size-8 icon arm", async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 768 });
  await page.goto("/f/sidebar-contract?case=menu-expanded");
  const trigger = page.locator("[data-gsxui-slot-sidebar-trigger]").first();
  const box = await trigger.evaluate((n) => {
    const s = getComputedStyle(n);
    return { width: s.width, height: s.height };
  });
  expect(box).toEqual({ width: "28px", height: "28px" });
});

test("SidebarSeparator keeps its own width, inset and colour against Separator's", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1024, height: 768 });
  await page.goto("/f/sidebar-contract?case=tokens");
  const separator = page.locator("[data-gsxui-slot-sidebar-separator]").first();
  const got = await separator.evaluate((n) => {
    const s = getComputedStyle(n);
    return {
      marginLeft: s.marginLeft,
      backgroundColor: s.backgroundColor,
      // Separator's own orientation arm is w-full; the sidebar slot's
      // data-[orientation=horizontal]:w-auto has to outrank it. This fixture
      // renders the desktop tree, which is display:none below the md
      // breakpoint's layout pass, so the computed value is the specified
      // keyword rather than a used length — which is exactly the fact under
      // test: "auto" means w-auto won, "100%" would mean w-full did.
      width: s.width,
    };
  });
  expect(got.marginLeft).toBe("8px");
  // --sidebar-border is rgb(19,20,21) in this fixture; Separator's own
  // bg-border would paint something else entirely.
  expect(got.backgroundColor).toBe("rgb(19, 20, 21)");
  expect(got.width).toBe("auto");
});

// SidebarInput's bg-background contests BOTH of Input's background utilities.
// bg-transparent is an ordinary merge conflict; dark:bg-input/30 is not — it
// survives the merge as a different variant group and would win the dark
// theme on Tailwind's own variant ordering. registry/styles/*/sidebar.css
// carries an explicit dark:bg-background for exactly that reason, and this is
// the pin that justifies a utility the pre-migration CSS did not have.
for (const scheme of ["light", "dark"] as const) {
  test(`SidebarInput paints the background token in the ${scheme} theme`, async ({ page }) => {
    await page.setViewportSize({ width: 1024, height: 768 });
    if (scheme === "dark") {
      await page.addInitScript(() => document.documentElement.classList.add("dark"));
    }
    await page.goto("/f/sidebar-contract?case=menu-expanded");
    const input = page.locator("[data-gsxui-slot-sidebar-input]").first();
    const got = await input.evaluate((n) => {
      const s = getComputedStyle(n);
      return {
        backgroundColor: s.backgroundColor,
        background: getComputedStyle(document.documentElement)
          .getPropertyValue("--background")
          .trim(),
      };
    });
    expect(got.backgroundColor).not.toBe("rgba(0, 0, 0, 0)");
    expect(got.background).not.toBe("");
  });
}

// Sidebar's own two same-specificity ties, both of which the pre-migration
// CSS resolved by source order at (0,0,0) and the recipe now resolves by
// arbitrary-variant specificity.
test("An active menu button's badge takes the primary foreground, hovered or not", async ({
  page,
}) => {
  await page.setViewportSize({ width: 1024, height: 768 });
  await page.goto("/f/sidebar-contract?case=menu-expanded");
  // The fixture renders both responsive trees; scope to the desktop one.
  const item = page.locator(
    '[data-gsxui-slot-sidebar-desktop] [data-sidebar-contract="action-item"]',
  );
  const badge = item.locator("[data-gsxui-slot-sidebar-menu-badge]");
  const button = item.locator("[data-gsxui-slot-sidebar-menu-button]");
  const resting = await badge.evaluate((n) => getComputedStyle(n).color);
  await button.hover();
  const hovered = await badge.evaluate((n) => getComputedStyle(n).color);
  // The active-sibling rule and the hovered-sibling rule are both (0,3,0);
  // active must win, as it did by source order before the migration.
  expect(hovered).toBe(resting);
});

test("An offcanvas rail sits flush and reaches past its own edge", async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 768 });
  // offcanvas-right is the collapsed offcanvas case — the desktop element
  // only carries data-collapsible="offcanvas" while collapsed. Scope to the
  // desktop tree; the mobile one renders a rail too.
  await page.goto("/f/sidebar-contract?case=offcanvas-right");
  const rail = page.locator(
    "[data-gsxui-slot-sidebar-desktop] [data-gsxui-slot-sidebar-rail]",
  );
  const got = await rail.evaluate((n) => {
    const s = getComputedStyle(n);
    return { transform: s.transform, after: getComputedStyle(n, "::after").left };
  });
  // foundation.css's rail base is translateX(-50%); the offcanvas variant
  // must still replace it outright rather than compose with it, which is why
  // the recipe writes [transform:translateX(0)] and not translate-x-0.
  expect(got.transform).toBe("matrix(1, 0, 0, 1, 0, 0)");
  // foundation.css's ::after base is left: 50%.
  expect(got.after).toBe("16px");
});


// Toast's migration. Two things it must not lose, neither of which the
// resting sweep can prove on its own.
test("Toast keeps its card chrome and its close button's own anchor", async ({
  page,
}) => {
  await page.goto("/x/toast/server");
  const card = page.locator("[data-gsxui-slot-toast]").first();
  expect(await card.evaluate((n) => getComputedStyle(n).width)).toBe("356px");
  expect(await card.evaluate((n) => getComputedStyle(n).borderTopLeftRadius)).toBe(
    "16px",
  );
  const close = page.locator("[data-gsxui-slot-toast-close]").first();
  expect(await close.evaluate((n) => getComputedStyle(n).position)).toBe(
    "absolute",
  );
});

// The standalone showcase row opts out of the stack with a plain `static`
// utility. foundation.css's stack-anchor rule is deliberately :where()-wrapped
// in @layer components so that caller utility wins; promoting it to @layer
// utilities (the layer gate's usual fix for a migrated component) flips this
// to absolute — see toastStackAnchorReason in internal/stylegen/layercheck.go.
test("Toast's caller `static` still beats the stack anchor", async ({ page }) => {
  await page.goto("/x/toast/server");
  const card = page.locator("[data-gsxui-slot-toast]").first();
  expect(await card.evaluate((n) => getComputedStyle(n).position)).toBe("static");
});

// data-type is NOT a recipe dimension: ui/toaster.js mutates it at runtime
// (morph() turns a loading toast into a success one in place), so the type's
// icon colour lives on the icon slot's own recipe rule as an ancestor-keyed
// arbitrary variant. This pin exercises exactly the path a baked-in class
// would break — the icon must recolour after the attribute changes.
test("Toast's type-keyed icon colour follows a runtime data-type change", async ({
  page,
}) => {
  await page.goto("/x/toaster/types");
  // Resolve the token through a probe element so the comparison is against a
  // computed colour in the same notation. A probe carrying `text-success` as a
  // class would not work: that utility is only ever emitted inside the icon
  // rule's arbitrary variant, never as a standalone class.
  const successColor = await page.evaluate(() => {
    const probe = document.createElement("div");
    probe.style.color = "var(--success)";
    document.body.appendChild(probe);
    const color = getComputedStyle(probe).color;
    probe.remove();
    return color;
  });

  await page.getByRole("button", { name: "Promise" }).click();
  const icon = page.locator(
    "#gsxui-toaster [data-gsxui-slot-toast] [data-gsxui-slot-toast-icon]",
  );
  await expect(icon).toHaveCount(1);
  await expect(
    page.locator("#gsxui-toaster [data-gsxui-slot-toast][data-type='loading']"),
  ).toHaveCount(1);

  await expect(
    page.locator("#gsxui-toaster [data-gsxui-slot-toast][data-type='success']"),
  ).toHaveCount(1, { timeout: 10_000 });
  expect(await icon.evaluate((n) => getComputedStyle(n).color)).toBe(
    successColor,
  );
});

// Toaster's migration. Its one authored declaration is not a utility at all
// but a custom property (--gsxui-toast-offset), which rides onto the recipe
// as a Tailwind arbitrary-property utility. foundation.css's mechanics rule
// on the CARD reads it through inheritance, so the region losing it would
// silently pin every toast to the viewport corner instead of insetting it.
test("Toaster's inset custom property still reaches the toast card", async ({
  page,
}) => {
  await page.goto("/x/toaster/types");
  const region = page.locator("[data-gsxui-slot-toaster]");
  expect(await region.evaluate((n) => getComputedStyle(n).padding)).toBe("24px");
  expect(
    await region.evaluate((n) =>
      getComputedStyle(n).getPropertyValue("--gsxui-toast-offset").trim(),
    ),
  ).toBe("1.5rem");

  await page.getByRole("button", { name: "Default" }).click();
  const card = page.locator("#gsxui-toaster [data-gsxui-slot-toast]").first();
  await expect(card).toBeVisible();
  expect(await card.evaluate((n) => getComputedStyle(n).right)).toBe("24px");
});

// Command's five command-dialog-content descendant rules and CommandItem's
// data-selected background had no fixture reaching them when Command migrated:
// site/examples/command/basic.gsx renders its CommandDialog closed, and nothing
// stamps data-selected at rest. Both translate faithfully and compile, but
// "compiles" is not "applies" — these open the dialog and stamp the attribute so
// the rules are actually exercised.

test("CommandDialog's in-dialog density reaches items once the dialog is open", async ({ page }) => {
  await page.goto("/x/command/basic");
  const dialog = page.locator("[data-gsxui-slot-command-dialog-content]");
  await dialog.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(dialog).toHaveJSProperty("open", true);
  const item = dialog.locator("[data-gsxui-slot-command-item]").first();
  // The in-dialog rule loosens the item's vertical padding; a plain CommandItem
  // outside a dialog does not carry it.
  expect(await item.evaluate((n) => getComputedStyle(n).padding)).toBe("12px 8px");
});

test("CommandItem's data-selected background paints the accent", async ({ page }) => {
  await page.goto("/x/command/basic");
  const dialog = page.locator("[data-gsxui-slot-command-dialog-content]");
  await dialog.evaluate((el) =>
    el.dispatchEvent(new CustomEvent("gsxui:request-open", { bubbles: true, cancelable: true })),
  );
  await expect(dialog).toHaveJSProperty("open", true);
  const item = dialog.locator("[data-gsxui-slot-command-item]").first();
  // data-selected is command.js's runtime stamp, so drive it directly.
  await item.evaluate((n) => n.setAttribute("data-selected", "true"));
  expect(await item.evaluate((n) => getComputedStyle(n).backgroundColor)).toBe("oklch(0.97 0 0)");
});
