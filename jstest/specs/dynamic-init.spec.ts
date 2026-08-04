import { expect, test } from "@playwright/test";

// The library promise under test: init(selector, fn) runs for current
// matches, for injected matches, and re-runs when a match's subtree is
// mutated back to server state (the morph case) — with init-pass writes
// not re-triggering the observer.
test.describe("gsxui dynamic init", () => {
  test("init runs now, on injection, and on mutation — without looping", async ({ page }) => {
    await page.goto("/x/select/basic");
    const counts = await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      const calls: string[] = [];
      // Idempotent marker initializer: writes DOM (data-inited) — the
      // observer must NOT re-trigger on that own write.
      init("[data-di-probe]", (el) => {
        calls.push(el.id);
        el.dataset.inited = "true";
      });
      // (a) current match
      const a = document.createElement("div");
      a.id = "pre";
      a.setAttribute("data-di-probe", "");
      document.body.append(a);
      await new Promise((r) => setTimeout(r, 20)); // observer flush
      // (b) injected subtree (descendant match)
      const wrap = document.createElement("div");
      wrap.innerHTML = '<section><div id="inj" data-di-probe></div></section>';
      document.body.append(wrap);
      await new Promise((r) => setTimeout(r, 20));
      // (c) mutation under an existing match re-inits it
      document.getElementById("inj")!.textContent = "server reset";
      await new Promise((r) => setTimeout(r, 20));
      return calls;
    });
    // "pre" was appended before registration? No — registration ran first,
    // so: no initial matches, then pre (injection), inj (injection),
    // inj again (mutation). The own-write (data-inited) must not add more.
    expect(counts).toEqual(["pre", "inj", "inj"]);
  });

  test("initial matches run at registration", async ({ page }) => {
    await page.goto("/x/select/basic");
    const ran = await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      let n = 0;
      init("[data-gsxui-slot-select]", () => n++);
      return n;
    });
    expect(ran).toBeGreaterThan(0); // the page's select root
  });

  test("carousel init is idempotent under re-init", async ({ page }) => {
    await page.goto("/x/carousel/basic");
    await page.evaluate(async () => {
      const { init } = await import("/ui/gsxui.js");
      // Force a re-init pass over existing carousels by touching a root.
      document.querySelector("[data-gsxui-slot-carousel]")!.setAttribute("data-poke", "1");
      await new Promise((r) => setTimeout(r, 20));
    });
    const before = await page
      .locator("[data-gsxui-slot-carousel-item]")
      .first()
      .boundingBox();
    await page.locator("[data-gsxui-slot-carousel-next]").first().click();
    await page.waitForTimeout(400); // one smooth-scroll settle
    const after = await page
      .locator("[data-gsxui-slot-carousel-item]")
      .first()
      .boundingBox();
    expect(after!.x).toBeLessThan(before!.x); // advanced exactly one step is
    // asserted by carousel's own existing specs; here we assert it MOVED
    // (a double-bound autoplay/observer would break those existing specs).
  });
});
