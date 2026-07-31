# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dialog.spec.ts >> dialog identity and trigger ownership stay within their nearest roots
- Location: jstest/specs/dialog.spec.ts:309:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dialog/basic", waiting until "load"

```

# Test source

```ts
  210 |     "close-button",
  211 |     "cancel",
  212 |     "backdrop",
  213 |   ]);
  214 | });
  215 | 
  216 | test("the finite exit stays open while closed and emits one notification after native close", async ({
  217 |   page,
  218 | }) => {
  219 |   await page.goto(BASIC);
  220 |   await page.evaluate(() => {
  221 |     (window as any).__dialogCloses = [];
  222 |     document.addEventListener("gsxui:close", (event) => {
  223 |       const dialog = event.target as HTMLDialogElement;
  224 |       (window as any).__dialogCloses.push({
  225 |         open: dialog.open,
  226 |         state: dialog.dataset.state,
  227 |       });
  228 |     });
  229 |   });
  230 | 
  231 |   await open(page);
  232 |   await dispatch(page, DIALOG, "gsxui:request-close");
  233 | 
  234 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  235 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  236 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  237 |   expect(await page.evaluate(() => (window as any).__dialogCloses)).toEqual([
  238 |     { open: false, state: "closed" },
  239 |   ]);
  240 | });
  241 | 
  242 | test("an open request during exit aborts the pending close", async ({ page }) => {
  243 |   await page.goto(BASIC);
  244 | 
  245 |   await open(page);
  246 |   await dispatch(page, DIALOG, "gsxui:request-close");
  247 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  248 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  249 | 
  250 |   await dispatch(page, DIALOG, "gsxui:request-open");
  251 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "open");
  252 |   await page.waitForTimeout(350);
  253 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  254 | });
  255 | 
  256 | test("overlapping close generations wait for the latest finite dialog animation", async ({
  257 |   page,
  258 | }) => {
  259 |   await page.goto(BASIC);
  260 |   const dialog = page.locator(DIALOG);
  261 | 
  262 |   await open(page);
  263 |   await dialog.evaluate((element) => {
  264 |     const animation = element.animate([{ opacity: 1 }, { opacity: 0 }], {
  265 |       duration: 60_000,
  266 |     });
  267 |     animation.pause();
  268 |     animation.currentTime = 0;
  269 |     (window as any).__dialogAnimationA = animation;
  270 |   });
  271 |   await dispatch(page, DIALOG, "gsxui:request-close");
  272 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  273 |   await expect(dialog).toHaveJSProperty("open", true);
  274 | 
  275 |   await dispatch(page, DIALOG, "gsxui:request-open");
  276 |   await expect(dialog).toHaveAttribute("data-state", "open");
  277 |   await dialog.evaluate((element) => {
  278 |     (window as any).__dialogAnimationB = element.animate(
  279 |       [{ opacity: 1 }, { opacity: 0 }],
  280 |       { duration: 60_000 },
  281 |     );
  282 |   });
  283 |   await dispatch(page, DIALOG, "gsxui:request-close");
  284 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  285 |   await expect(dialog).toHaveJSProperty("open", true);
  286 | 
  287 |   await page.evaluate(async () => {
  288 |     const animation = (window as any).__dialogAnimationA as Animation;
  289 |     animation.finish();
  290 |     await animation.finished;
  291 |     await Promise.resolve();
  292 |   });
  293 |   expect(
  294 |     await dialog.evaluate((element) => ({
  295 |       open: (element as HTMLDialogElement).open,
  296 |       state: (element as HTMLElement).dataset.state,
  297 |       latestAnimation: ((window as any).__dialogAnimationB as Animation).playState,
  298 |     })),
  299 |   ).toEqual({
  300 |     open: true,
  301 |     state: "closed",
  302 |     latestAnimation: "running",
  303 |   });
  304 | 
  305 |   await page.evaluate(() => ((window as any).__dialogAnimationB as Animation).finish());
  306 |   await expect(dialog).toHaveJSProperty("open", false);
  307 | });
  308 | 
  309 | test("dialog identity and trigger ownership stay within their nearest roots", async ({ page }) => {
> 310 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  311 |   await page.evaluate(() => {
  312 |     document.body.insertAdjacentHTML(
  313 |       "beforeend",
  314 |       `
  315 |         <div data-gsxui-dialog id="generated-root">
  316 |           <button id="generated-trigger" data-gsxui-dialog-trigger aria-expanded="false">Generated</button>
  317 |           <dialog data-gsxui-dialog-content data-state="closed">
  318 |             <h2 data-gsxui-dialog-title>Generated title</h2>
  319 |             <p data-gsxui-dialog-description>Generated description</p>
  320 |             <button data-gsxui-dialog-close>Close generated</button>
  321 |             <div data-gsxui-dialog id="nested-root">
  322 |               <button id="nested-trigger" data-gsxui-dialog-trigger aria-expanded="false">Nested</button>
  323 |               <button id="nested-root-close" data-gsxui-dialog-close>Close nested</button>
  324 |               <dialog data-gsxui-dialog-content data-state="closed">
  325 |                 <h2 data-gsxui-dialog-title>Nested title</h2>
  326 |                 <p data-gsxui-dialog-description>Nested description</p>
  327 |               </dialog>
  328 |             </div>
  329 |           </dialog>
  330 |         </div>
  331 |         <div data-gsxui-dialog id="second-root">
  332 |           <button id="second-trigger" data-gsxui-dialog-trigger aria-expanded="false">Second</button>
  333 |           <dialog data-gsxui-dialog-content data-state="closed">
  334 |             <h2 data-gsxui-dialog-title>Second title</h2>
  335 |             <p data-gsxui-dialog-description>Second description</p>
  336 |             <button data-gsxui-dialog-close>Close second</button>
  337 |           </dialog>
  338 |         </div>
  339 |         <div data-gsxui-dialog id="authored-root">
  340 |           <button id="authored-trigger" data-gsxui-dialog-trigger aria-controls="keep-control" aria-expanded="false">Authored</button>
  341 |           <dialog id="authored-dialog" data-gsxui-dialog-content data-state="closed"
  342 |             aria-labelledby="authored-label" aria-describedby="authored-description">
  343 |             <h2 id="authored-label" data-gsxui-dialog-title>Authored title</h2>
  344 |             <p id="authored-description" data-gsxui-dialog-description>Authored description</p>
  345 |           </dialog>
  346 |         </div>
  347 |       `,
  348 |     );
  349 |   });
  350 | 
  351 |   await page.locator("#generated-trigger").click();
  352 |   const generated = page.locator("#generated-root > dialog");
  353 |   const nested = page.locator("#nested-root dialog");
  354 |   await expect(generated).toHaveJSProperty("open", true);
  355 |   await expect(nested).toHaveJSProperty("open", false);
  356 |   await expect(page.locator("#generated-trigger")).toHaveAttribute("aria-expanded", "true");
  357 |   await expect(page.locator("#nested-trigger")).toHaveAttribute("aria-expanded", "false");
  358 | 
  359 |   await page.evaluate(() => {
  360 |     document.addEventListener("gsxui:request-close", (event) => {
  361 |       const target = event.target as Element;
  362 |       (window as any).__nestedCloseTarget = target
  363 |         .closest("[data-gsxui-dialog]")
  364 |         ?.id;
  365 |     });
  366 |   });
  367 |   await page.locator("#nested-root-close").click();
  368 |   await expect(generated).toHaveJSProperty("open", true);
  369 |   expect(await page.evaluate(() => (window as any).__nestedCloseTarget)).toBe("nested-root");
  370 | 
  371 |   const stableGeneratedID = await generated.getAttribute("id");
  372 |   await page.evaluate((generatedID) => {
  373 |     document.body.insertAdjacentHTML(
  374 |       "beforeend",
  375 |       `<button id="generated-id-invoker" commandfor="${generatedID}" command="show-modal" aria-expanded="unchanged">Generated ID invoker</button>`,
  376 |     );
  377 |   }, stableGeneratedID);
  378 |   await page.locator("#generated-root > dialog > [data-gsxui-dialog-close]").click();
  379 |   await expect(generated).toHaveJSProperty("open", false);
  380 |   await expect(page.locator("#generated-id-invoker")).toHaveAttribute(
  381 |     "aria-expanded",
  382 |     "unchanged",
  383 |   );
  384 |   await page.locator("#generated-trigger").click();
  385 |   await expect(generated).toHaveAttribute("id", stableGeneratedID!);
  386 |   await expect(page.locator("#generated-id-invoker")).toHaveAttribute(
  387 |     "aria-expanded",
  388 |     "unchanged",
  389 |   );
  390 |   await page.locator("#generated-root > dialog > [data-gsxui-dialog-close]").click();
  391 |   await expect(generated).toHaveJSProperty("open", false);
  392 | 
  393 |   await page.locator("#second-trigger").click();
  394 |   const second = page.locator("#second-root dialog");
  395 |   await expect(second).toHaveJSProperty("open", true);
  396 |   await page.locator("#second-root [data-gsxui-dialog-close]").click();
  397 |   await expect(second).toHaveJSProperty("open", false);
  398 |   await page.locator("#authored-trigger").click();
  399 |   await expect(page.locator("#authored-dialog")).toHaveJSProperty("open", true);
  400 |   const identity = await page.evaluate(() => {
  401 |     const roots = [
  402 |       document.querySelector("#generated-root")!,
  403 |       document.querySelector("#second-root")!,
  404 |       document.querySelector("#authored-root")!,
  405 |     ];
  406 |     return roots.map((root) => {
  407 |       const dialog = [...root.querySelectorAll("dialog[data-gsxui-dialog-content]")].find(
  408 |         (element) => element.closest("[data-gsxui-dialog]") === root,
  409 |       )!;
  410 |       const title = [...root.querySelectorAll("[data-gsxui-dialog-title]")].find(
```