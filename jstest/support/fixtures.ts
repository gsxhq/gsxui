import { test as base } from "@playwright/test";
import type { Registration } from "./globals";

/**
 * `registrations` is the full delegation registry, recorded once per worker
 * by loading /registrations — a blank page whose only script is
 * /shim/index.js, where on() records instead of binding.
 *
 * Worker-scoped rather than collected in globalSetup because globalSetup's
 * ordering against webServer startup is not something to depend on: a
 * fixture runs after the server is definitely up.
 */
export const test = base.extend<{}, { registrations: Registration[] }>({
  registrations: [
    async ({ browser }, use) => {
      const page = await browser.newPage();
      await page.goto("/registrations");
      const recorded = await page.evaluate(() => window.__gsxuiRegistrations);
      await page.close();
      if (!recorded || recorded.length === 0) {
        throw new Error("/registrations recorded nothing — the shim did not run");
      }
      await use(recorded);
    },
    { scope: "worker" },
  ],
});

export { expect } from "@playwright/test";
