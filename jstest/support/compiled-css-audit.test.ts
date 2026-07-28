import assert from "node:assert/strict";
import test from "node:test";
import { compiledCSSViolations } from "./compiled-css-audit.ts";

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
