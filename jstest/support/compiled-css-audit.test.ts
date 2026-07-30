import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { parse } from "postcss";
import selectorParser from "postcss-selector-parser";
import { compiledCSSViolations } from "./compiled-css-audit.ts";

const repositoryRoot = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../..",
);

/**
 * The default style's components layer is authored one file per component under
 * assets/css/styles/default/, so migrating a component is a file deletion rather
 * than a block deleted out of the middle of one 3000-line sheet. Reading
 * default.css alone would now see only @import lines and the @layer utilities
 * escape hatch, and every assertion below about a component's authored rules
 * would pass vacuously.
 *
 * This inlines each relative @import in the aggregator's own order — the same
 * composition internal/cli vendors as a consumer's flat style.css — so these
 * tests keep reading exactly the stylesheet a consumer receives.
 */
function defaultStyleCSS(): string {
  const dir = join(repositoryRoot, "assets/css/styles");
  const aggregator = readFileSync(join(dir, "default.css"), "utf8");
  return aggregator
    .split("\n")
    .map((line) => {
      const match = /^@import "(\.\/[^"]+)";[ \t]*$/.exec(line);
      return match ? readFileSync(join(dir, match[1]), "utf8") : line + "\n";
    })
    .join("");
}

for (const style of ["nova", "maia"]) {
  test(`compiles the authored ${style} Button recipe with Tailwind`, () => {
    const input = [
      '@import "tailwindcss" source(none);',
      '@import "./assets/css/foundation.css";',
      '@import "./assets/css/themes/default.css";',
      `@import "./registry/styles/${style}/button.css";`,
    ].join("\n");
    const result = spawnSync(
      join(repositoryRoot, "node_modules", ".bin", "tailwindcss"),
      ["-i", "-", "--cwd", repositoryRoot, "--silent"],
      {
        cwd: repositoryRoot,
        encoding: "utf8",
        input,
      },
    );

    assert.equal(
      result.status,
      0,
      `Tailwind failed to compile ${style} Button recipe:\n${result.stderr}`,
    );
    assert.deepEqual(compiledCSSViolations(result.stdout), []);
  });
}

test("compiled site Button fallback has the exact normal-site scope", () => {
  const result = spawnSync(
    join(repositoryRoot, "node_modules", ".bin", "tailwindcss"),
    ["-i", "web/site.css", "--cwd", repositoryRoot, "--silent"],
    {
      cwd: repositoryRoot,
      encoding: "utf8",
    },
  );
  assert.equal(
    result.status,
    0,
    `Tailwind failed to compile web/site.css:\n${result.stderr}`,
  );

  const selectors = compiledButtonSelectors(result.stdout);
  assert.notEqual(selectors.length, 0, "compiled site CSS has no Button fallback");
  for (const selector of selectors) {
    assert.equal(
      hasExactNormalSiteScope(selector),
      true,
      `Button fallback can escape the normal-site scope: ${selector.toString()}`,
    );
  }
});

test("consumer composition selectors retain their owning marker specificity", () => {
  const root = parse(defaultStyleCSS());
  const buttonGroupSizes = new Set<string>();
  let calendarNavGeometry = false;
  let calendarDayGeometry = false;

  root.walkRules((rule) => {
    const selectorRoot = selectorParser().astSync(rule.selector);
    for (const selector of selectorRoot.nodes) {
      const attributes: import("postcss-selector-parser").Attribute[] = [];
      selector.walkAttributes((attribute) => attributes.push(attribute));
      const names = new Set(
        attributes.map((attribute) => attribute.attribute.toLowerCase()),
      );

      if (
        names.has("data-gsxui-slot-button-group") &&
        names.has("data-size")
      ) {
        const size = attributes.find(
          (attribute) => attribute.attribute.toLowerCase() === "data-size",
        );
        assert.equal(
          size?.parent,
          selector,
          `ButtonGroup data-size must retain selector specificity: ${selector.toString()}`,
        );
        assert.equal(
          selector.nodes.some(
            (node) => node.type === "combinator" && node.value.trim() === ">",
          ),
          true,
          `ButtonGroup size owner must be a direct child: ${selector.toString()}`,
        );
        buttonGroupSizes.add(size?.value ?? "");
      }

      const declarations = new Set(
        rule.nodes
          .filter((node) => node.type === "decl")
          .map((declaration) => declaration.prop.toLowerCase()),
      );
      const nav = attributes.find(
        (attribute) =>
          attribute.attribute.toLowerCase() ===
          "data-gsxui-slot-calendar-nav-button",
      );
      if (nav && declarations.has("width") && declarations.has("height")) {
        assert.equal(nav.parent, selector);
        calendarNavGeometry = true;
      }
      const day = attributes.find(
        (attribute) =>
          attribute.attribute.toLowerCase() ===
          "data-gsxui-slot-calendar-day-button",
      );
      if (day && declarations.has("min-width")) {
        assert.equal(day.parent, selector);
        calendarDayGeometry = true;
      }
    }
  });

  assert.deepEqual(
    [...buttonGroupSizes].sort(),
    ["icon-sm", "icon-xs", "xs"],
  );
  assert.equal(calendarNavGeometry, true);
  assert.equal(calendarDayGeometry, true);
});

test("binary Sidebar markers use presence selectors", () => {
  const root = parse(defaultStyleCSS());
  const markers = new Map([
    ["data-active", false],
    ["data-show-on-hover", false],
  ]);

  root.walkRules((rule) => {
    const selectorRoot = selectorParser().astSync(rule.selector);
    for (const selector of selectorRoot.nodes) {
      const attributes: import("postcss-selector-parser").Attribute[] = [];
      selector.walkAttributes((attribute) => attributes.push(attribute));
      if (
        !attributes.some((attribute) =>
          attribute.attribute
            .toLowerCase()
            .startsWith("data-gsxui-slot-sidebar-"),
        )
      ) {
        continue;
      }
      for (const attribute of attributes) {
        const name = attribute.attribute.toLowerCase();
        if (!markers.has(name)) continue;
        assert.equal(
          attribute.operator,
          undefined,
          `${name} must be selected by presence: ${selector.toString()}`,
        );
        markers.set(name, true);
      }
    }
  });

  for (const [marker, found] of markers) {
    assert.equal(found, true, `default CSS has no [${marker}] selector`);
  }
});

test("allows only Tailwind's exact hidden preflight important declaration", () => {
  const css = `
    [hidden]:where(:not([hidden="until-found"])) {
      display: none !important;
    }
  `;

  assert.deepEqual(compiledCSSViolations(css), []);
});

test("allows the same hidden preflight declaration after quote-minification", () => {
  const css =
    "[hidden]:where(:not([hidden=until-found])){display:none!important}";

  assert.deepEqual(compiledCSSViolations(css), []);
});

test("rejects an important utility declaration", () => {
  const violations = compiledCSSViolations(`
    .rounded-b-lg\\! {
      border-bottom-right-radius: 0.5rem !important;
    }
  `);

  assert.equal(violations.length, 1);
  assert.equal(violations[0].label, "!important declaration");
  assert.equal(violations[0].count, 1);
});

test("rejects legacy and packed slot selectors without rejecting presence selectors", () => {
  const violations = compiledCSSViolations(`
    [data-slot="button"] { display: inline-flex; }
    [data-gsxui-slot~="button"] { display: inline-flex; }
    [data-gsxui-slot-button] { display: inline-flex; }
    [data-gsxui-slot-button="x"] { display: inline-flex; }
  `);

  assert.equal(violations.length, 3);
  assert.equal(violations[0].label, "legacy [data-slot selector");
  assert.equal(violations[0].count, 1);
  assert.equal(
    violations[1].label,
    "obsolete packed [data-gsxui-slot] selector",
  );
  assert.equal(violations[1].count, 1);
  assert.equal(violations[2].label, "valued [data-gsxui-slot-*] selector");
  assert.equal(violations[2].count, 1);
});

test("rejects comments and escapes that encode forbidden slot selectors", () => {
  const violations = compiledCSSViolations(`
    [data-gsxui-slot-button/**/="x"] { display: inline-flex; }
    [data-gsxui-slot-button\\2d icon="x"] { display: inline-flex; }
    [data-gsxui\\2d slot~="button"] { display: inline-flex; }
  `);

  assert.deepEqual(
    Object.fromEntries(violations.map(({ label, count }) => [label, count])),
    {
      "obsolete packed [data-gsxui-slot] selector": 1,
      "valued [data-gsxui-slot-*] selector": 2,
    },
  );
});

test("fails the audit when a CSS rule has an invalid selector", () => {
  assert.throws(
    () => compiledCSSViolations(`a:: { color: red; }`),
    /invalid CSS selector/,
  );
});

test("rejects important declarations hidden behind the preflight selector", () => {
  const violations = compiledCSSViolations(`
    [hidden]:where(:not([hidden="until-found"])) {
      visibility: hidden !important;
    }
  `);

  assert.equal(violations.length, 1);
  assert.equal(violations[0].label, "!important declaration");
});

function compiledButtonSelectors(css: string) {
  const selectors: ReturnType<
    ReturnType<typeof selectorParser>["astSync"]
  >["nodes"] = [];
  const root = parse(css);
  root.walkRules((rule) => {
    const selectorRoot = selectorParser().astSync(rule.selector);
    for (const selector of selectorRoot.nodes) {
      let targetsButton = false;
      selector.walkAttributes((attribute) => {
        if (attribute.attribute.toLowerCase() === "data-gsxui-slot-button") {
          targetsButton = true;
        }
      });
      if (targetsButton) selectors.push(selector);
    }
  });
  return selectors;
}

function hasExactNormalSiteScope(
  selector: ReturnType<
    ReturnType<typeof selectorParser>["astSync"]
  >["nodes"][number],
) {
  const nodes = selector.nodes;
  for (let index = 0; index + 1 < nodes.length; index++) {
    const scope = nodes[index];
    const descendant = nodes[index + 1];
    if (
      scope.type !== "pseudo" ||
      scope.value.toLowerCase() !== ":where" ||
      descendant.type !== "combinator" ||
      descendant.value.trim() !== ""
    ) {
      continue;
    }
    const scopeSelectors = scope.nodes;
    if (scopeSelectors?.length !== 1) continue;
    const scopeNodes = scopeSelectors[0].nodes;
    if (scopeNodes.length !== 2) continue;
    const body = scopeNodes[0];
    const exclusion = scopeNodes[1];
    if (
      body.type !== "tag" ||
      body.value.toLowerCase() !== "body" ||
      exclusion.type !== "pseudo" ||
      exclusion.value.toLowerCase() !== ":not"
    ) {
      continue;
    }
    const exclusionSelectors = exclusion.nodes;
    if (exclusionSelectors?.length !== 1) continue;
    const exclusionNodes = exclusionSelectors[0].nodes;
    if (exclusionNodes.length !== 1) continue;
    const attribute = exclusionNodes[0];
    if (
      attribute.type === "attribute" &&
      attribute.attribute.toLowerCase() === "data-theme-button-preview" &&
      attribute.operator === undefined
    ) {
      return true;
    }
  }
  return false;
}
