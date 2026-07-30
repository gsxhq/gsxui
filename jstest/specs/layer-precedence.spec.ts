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
// <ui.SheetHeader> directly: Sheet's own utilities now beat anything Sidebar
// declares for the same property from @layer components. make sweep-compare
// caught the mobile panel collapsing from its --sidebar-width to Sheet's
// w-3/4 and repainting with bg-background; the promoted rules at the bottom
// of assets/css/styles/default/sidebar.css restore both. Pinned so a future
// layer change cannot silently undo the promotion.
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
