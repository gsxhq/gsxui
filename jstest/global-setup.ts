import { execFileSync } from "node:child_process";
import { existsSync, mkdirSync, readdirSync, readFileSync, symlinkSync, unlinkSync } from "node:fs";
import path from "node:path";
import { cssPath, manifestPath, repoRoot, tmpDir } from "./support/paths";

/**
 * Runs once before any worker imports a spec file, which is what lets specs
 * read the manifest synchronously and still work under a bare
 * `npx playwright test`.
 *
 * The stylesheet is compiled from web/site.css — the production entry — so
 * every class a component can emit is present (site.css already declares
 * `@source "../ui/**\/*.gsx"`). Real CSS is not optional here: the ghost-box
 * invariant is a computed-style assertion, and it caught shipped defects in
 * both dialog and sidebar.
 */
export default function globalSetup() {
  mkdirSync(tmpDir, { recursive: true });

  execFileSync("go", ["run", "./jstest/harness", "-manifest", manifestPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });

  execFileSync("npx", ["@tailwindcss/cli", "-i", "web/site.css", "-o", cssPath], {
    cwd: repoRoot,
    stdio: "inherit",
  });

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
    if (!source) {
      throw new Error(
        `compiled stylesheet at ${cssPath} references font file "${name}" but it was ` +
          `not found under any package in ${fontsourceDir}`,
      );
    }
    const link = path.join(filesDir, name);
    if (existsSync(link)) unlinkSync(link);
    symlinkSync(source, link);
  }
}
