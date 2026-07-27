import { readFileSync } from "node:fs";
import { pathToFileURL } from "node:url";
import { parse, type Declaration } from "postcss";

type CompiledCSSViolation = {
  label: string;
  count: number;
  samples: string[];
};

const tailwindHiddenSelectors = new Set([
  '[hidden]:where(:not([hidden="until-found"]))',
  "[hidden]:where(:not([hidden=until-found]))",
]);

function isTailwindHiddenPreflight(declaration: Declaration): boolean {
  const parent = declaration.parent;
  return (
    parent?.type === "rule" &&
    tailwindHiddenSelectors.has(parent.selector) &&
    declaration.prop === "display" &&
    declaration.value === "none"
  );
}

function selectorUsesAttribute(selector: string, attribute: string): boolean {
  const attributeSelector = new RegExp(
    String.raw`\[\s*${attribute}(?=\s*(?:[~|^$*]?=|\]))`,
    "i",
  );
  return attributeSelector.test(selector);
}

export function compiledCSSViolations(css: string): CompiledCSSViolation[] {
  const root = parse(css);
  const legacySlots: string[] = [];
  const packedGSXUISlots: string[] = [];
  const importantDeclarations: string[] = [];

  root.walkRules((rule) => {
    if (selectorUsesAttribute(rule.selector, "data-slot")) {
      legacySlots.push(rule.selector);
    }
    if (selectorUsesAttribute(rule.selector, "data-gsxui-slot")) {
      packedGSXUISlots.push(rule.selector);
    }
  });
  root.walkDecls((declaration) => {
    if (declaration.important && !isTailwindHiddenPreflight(declaration)) {
      const parent = declaration.parent;
      const context = parent?.type === "rule" ? parent.selector : parent?.type;
      importantDeclarations.push(
        `${context ?? "<unknown>"} { ${declaration.prop}: ${declaration.value} !important }`,
      );
    }
  });

  return [
    {
      label: "legacy [data-slot selector",
      count: legacySlots.length,
      samples: [...new Set(legacySlots)].slice(0, 3),
    },
    {
      label: "obsolete packed [data-gsxui-slot] selector",
      count: packedGSXUISlots.length,
      samples: [...new Set(packedGSXUISlots)].slice(0, 3),
    },
    {
      label: "!important declaration",
      count: importantDeclarations.length,
      samples: [...new Set(importantDeclarations)].slice(0, 3),
    },
  ].filter((violation) => violation.count > 0);
}

export function auditCompiledCSSFile(path: string) {
  const violations = compiledCSSViolations(readFileSync(path, "utf8"));
  if (violations.length === 0) return;

  const detail = violations
    .map(
      ({ label, count, samples }) =>
        `${label}: ${count}${samples.length > 0 ? ` (${samples.join(", ")})` : ""}`,
    )
    .join("; ");
  throw new Error(`compiled CSS boundary violation in ${path}: ${detail}`);
}

if (
  process.argv[1] &&
  import.meta.url === pathToFileURL(process.argv[1]).href
) {
  for (const path of process.argv.slice(2)) {
    auditCompiledCSSFile(path);
  }
}
