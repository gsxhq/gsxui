import { execFileSync } from "node:child_process";
import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  statSync,
  symlinkSync,
  unlinkSync,
} from "node:fs";
import path from "node:path";
import { auditCompiledCSSFile } from "./support/compiled-css-audit.ts";
import {
  cssPath,
  foundationCSSPath,
  manifestPath,
  repoRoot,
  tmpDir,
} from "./support/paths";

/**
 * Runs once before any worker imports a spec file, which is what lets specs
 * read the manifest synchronously and still work under a bare
 * `npx playwright test`.
 *
 * The full stylesheet imports web/site.css — the production entry — and adds
 * the harness's own authored fixtures to its explicit source boundary. Real
 * CSS is not optional here: the ghost-box invariant is a computed-style
 * assertion, and it caught shipped defects in both dialog and sidebar.
 */
export default function globalSetup() {
  mkdirSync(tmpDir, { recursive: true });

  execFileSync(
    path.join(repoRoot, "node_modules", ".bin", "esbuild"),
    [
      "node_modules/postcss/lib/parse.js",
      "--bundle",
      "--platform=browser",
      "--format=esm",
      `--outfile=${path.join(tmpDir, "postcss-parse.mjs")}`,
      "--log-level=warning",
    ],
    {
      cwd: repoRoot,
      stdio: "inherit",
    },
  );

  execFileSync("go", ["run", "./jstest/harness", "-manifest", manifestPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });

  execFileSync("npx", ["@tailwindcss/cli", "-i", "jstest/full.css", "-o", cssPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });
  execFileSync(
    "npx",
    ["@tailwindcss/cli", "-i", "jstest/foundation.css", "-o", foundationCSSPath],
    {
      cwd: repoRoot,
      stdio: "inherit",
    },
  );

  auditCompiledCSSFile(cssPath);
  auditCompiledCSSFile(foundationCSSPath);
  linkFontAssets();
}

/**
 * Tailwind's @import bundling carries @font-face url() text over from the
 * source stylesheet (node_modules/@fontsource-variable/*\/wght.css)
 * unchanged — it does not rewrite the path for the new output location. So
 * the compiled CSS ends up with `url(./files/geist-....woff2)`, a path
 * meant to resolve relative to the fontsource package directory, not
 * relative to jstest/.tmp/site.css. The harness serves /static/ from the
 * repo root, so the browser's actual request lands on
 * /static/jstest/.tmp/files/<name>, which does not exist — every page load
 * would 404 on every font file.
 *
 * This populates jstest/.tmp/files/ with a symlink per filename the
 * compiled CSS actually references, pointing at the real file under
 * node_modules/@fontsource-variable/<package>/files/. That makes the
 * relative url()s the browser sends resolve, without touching Tailwind's
 * output or the harness's /static/ mount.
 *
 * Why this function verifies its own output rather than trusting the link
 * step: each Geist/Geist Mono subset in the compiled CSS carries a
 * unicode-range, and Chromium only fetches a @font-face src when the page
 * actually needs a glyph in that subset — it does not eagerly fetch every
 * declared-but-unused source. gsxui's examples are Latin-only, so a page
 * load (and smoke.spec.ts's "no failed subresource requests" check) only
 * ever exercises the latin subset; the other nine files (cyrillic,
 * cyrillic-ext, vietnamese, latin-ext, and the four -mono equivalents) are
 * never requested by any browser navigation in this suite. A broken link
 * among those nine — a typo'd package directory, a filename collision
 * that picks the wrong @fontsource-variable package — would stay green
 * forever if this function only checked the subset a browser happens to
 * ask for. So it checks all of them, at setup time, regardless of what
 * any page requests.
 */
function linkFontAssets() {
  const css = readFileSync(cssPath, "utf8");
  const names = new Set<string>();
  for (const m of css.matchAll(/url\(\.\/files\/([^)'"]+)\)/g)) {
    names.add(m[1]);
  }
  if (names.size === 0) return;

  const filesDir = path.join(tmpDir, "files");
  mkdirSync(filesDir, { recursive: true });

  const fontsourceDir = path.join(repoRoot, "node_modules", "@fontsource-variable");
  const packages = readdirSync(fontsourceDir);

  for (const name of names) {
    const source = packages
      .map((pkg) => path.join(fontsourceDir, pkg, "files", name))
      .find((candidate) => existsSync(candidate));
    if (!source) continue; // reported below, together with any broken link
    const link = path.join(filesDir, name);
    if (existsSync(link)) unlinkSync(link);
    symlinkSync(source, link);
  }

  const unresolved = [...names].filter((name) => {
    try {
      return !statSync(path.join(filesDir, name)).isFile();
    } catch {
      return true; // missing entirely, or a dangling symlink
    }
  });
  if (unresolved.length > 0) {
    throw new Error(
      `compiled stylesheet at ${cssPath} references font files that did not resolve to a ` +
        `readable file under ${fontsourceDir}/*/files: ${unresolved.join(", ")}`,
    );
  }
}
