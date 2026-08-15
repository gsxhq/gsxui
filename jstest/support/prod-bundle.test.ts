// Production-bundle gate: the browser suite runs against dev-served,
// unbundled ui/*.js, so it cannot see what Vite does to a module at build
// time. This test builds the real production bundle once and asserts the
// contract prod depends on: every dynamically imported ui module the
// runtime needs must be EMITTED as a chunk and REFERENCED from the bundle
// that imports it. Regression pin: chart.js's renderer import used a
// computed specifier, which Rollup cannot follow — the chunk was never
// emitted and prod 404'd on ./chart.render.js while every dev test passed.
import { execFileSync } from "node:child_process";
import { readdirSync, readFileSync, existsSync } from "node:fs";
import { join } from "node:path";
import { test } from "node:test";
import assert from "node:assert/strict";

const root = new URL("../../", import.meta.url).pathname;
const assets = join(root, "site", "dist", "assets");

test("vite build emits and references the chart renderer chunk", () => {
  execFileSync("npx", ["vite", "build"], { cwd: root, stdio: "pipe" });
  const files = readdirSync(assets);
  const renderer = files.find((f) => /^chart\.render-[\w-]+\.js$/.test(f));
  assert.ok(renderer, `no chart.render-*.js chunk in dist/assets: ${files.join(", ")}`);
  // The chunk must be reachable from the stub's own bundled code: some
  // non-renderer bundle carries the stub (its failure log string is a
  // stable marker) AND references the emitted chunk by its hashed name.
  // Vite may hoist ui/index.js into whichever entry chunk it likes (today
  // main-*.js imports it from preview-*.js), so look at all of them.
  const bundles = files.filter((f) => f.endsWith(".js") && !f.startsWith("chart.render-"));
  const stubBundles = bundles.filter((f) =>
    readFileSync(join(assets, f), "utf8").includes("gsxui: failed to load chart.render.js"),
  );
  assert.ok(stubBundles.length > 0, `no bundle contains the chart stub: ${bundles.join(", ")}`);
  const wired = stubBundles.filter((f) => readFileSync(join(assets, f), "utf8").includes(renderer!));
  assert.ok(
    wired.length > 0,
    `bundle(s) ${stubBundles.join(", ")} carry the stub but never reference ${renderer} — the dynamic import was not followed`,
  );
  // vite build deletes the placeholder the repo keeps; restore it so the
  // tree stays clean for the commit gates that check it.
  const keep = join(root, "site", "dist", ".gitkeep");
  if (!existsSync(keep)) execFileSync("touch", [keep]);
});
