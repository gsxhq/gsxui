# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dialog.spec.ts >> infinite child and dialog animations do not keep a request close open
- Location: jstest/specs/dialog.spec.ts:453:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dialog/basic", waiting until "load"

```

# Test source

```ts
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
  411 |         (element) => element.closest("[data-gsxui-dialog]") === root,
  412 |       )!;
  413 |       const description = [...root.querySelectorAll("[data-gsxui-dialog-description]")].find(
  414 |         (element) => element.closest("[data-gsxui-dialog]") === root,
  415 |       )!;
  416 |       const trigger = root.querySelector("[data-gsxui-dialog-trigger]")!;
  417 |       return {
  418 |         dialog: dialog.id,
  419 |         labelledby: dialog.getAttribute("aria-labelledby"),
  420 |         describedby: dialog.getAttribute("aria-describedby"),
  421 |         title: title.id,
  422 |         description: description.id,
  423 |         controls: trigger.getAttribute("aria-controls"),
  424 |       };
  425 |     });
  426 |   });
  427 | 
  428 |   expect(identity[0].dialog).toMatch(/^gsxui-dialog-/);
  429 |   expect(identity[0].title).toMatch(/^gsxui-title-/);
  430 |   expect(identity[0].description).toMatch(/^gsxui-desc-/);
  431 |   expect(identity[0].controls).toBe(identity[0].dialog);
  432 |   expect(identity[1].dialog).toMatch(/^gsxui-dialog-/);
  433 |   expect(identity[1].title).toMatch(/^gsxui-title-/);
  434 |   expect(identity[1].description).toMatch(/^gsxui-desc-/);
  435 |   expect(identity[1].controls).toBe(identity[1].dialog);
  436 |   expect(identity[2]).toEqual({
  437 |     dialog: "authored-dialog",
  438 |     labelledby: "authored-label",
  439 |     describedby: "authored-description",
  440 |     title: "authored-label",
  441 |     description: "authored-description",
  442 |     controls: "keep-control",
  443 |   });
  444 |   expect(new Set([identity[0].dialog, identity[1].dialog]).size).toBe(2);
  445 |   expect(new Set([identity[0].title, identity[1].title]).size).toBe(2);
  446 |   expect(new Set([identity[0].description, identity[1].description]).size).toBe(2);
  447 |   expect(await page.locator(`#${identity[0].labelledby}`).count()).toBe(1);
  448 |   expect(await page.locator(`#${identity[0].describedby}`).count()).toBe(1);
  449 |   expect(await page.locator(`#${identity[1].labelledby}`).count()).toBe(1);
  450 |   expect(await page.locator(`#${identity[1].describedby}`).count()).toBe(1);
  451 | });
  452 | 
  453 | test("infinite child and dialog animations do not keep a request close open", async ({ page }) => {
> 454 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  455 |   const dialog = page.locator(DIALOG);
  456 | 
  457 |   await open(page);
  458 |   await dialog.evaluate((element) => {
  459 |     const child = document.createElement("span");
  460 |     child.textContent = "spinning child";
  461 |     element.append(child);
  462 |     child.animate([{ opacity: 0 }, { opacity: 1 }], { duration: 10_000, iterations: Infinity });
  463 |   });
  464 |   await close(page);
  465 | 
  466 |   await open(page);
  467 |   await dialog.evaluate((element) => {
  468 |     element.animate([{ opacity: 0 }, { opacity: 1 }], { duration: 10_000, iterations: Infinity });
  469 |   });
  470 |   await dispatch(page, DIALOG, "gsxui:request-close");
  471 |   await expect(dialog).toHaveJSProperty("open", false, { timeout: 350 });
  472 | });
  473 | 
  474 | test("native invokers address authored dialogs exactly and synchronize their ARIA state", async ({
  475 |   page,
  476 | }) => {
  477 |   await page.goto(BASIC);
  478 |   await page.evaluate(() => {
  479 |     document.body.insertAdjacentHTML(
  480 |       "beforeend",
  481 |       `
  482 |         <button id="open-first" commandfor="native-first" command="show-modal">Open first</button>
  483 |         <button id="open-second" commandfor="native-second" command="show-modal">Open second</button>
  484 |         <dialog id="native-first" data-gsxui-dialog-content data-state="closed"><p>First</p></dialog>
  485 |         <dialog id="native-second" data-gsxui-dialog-content data-state="closed"><p>Second</p></dialog>
  486 |       `,
  487 |     );
  488 |   });
  489 | 
  490 |   await page.locator("#open-first").click();
  491 |   await expect(page.locator("#native-first")).toHaveJSProperty("open", true);
  492 |   await expect(page.locator("#native-second")).toHaveJSProperty("open", false);
  493 |   await expect(page.locator("#open-first")).toHaveAttribute("aria-expanded", "true");
  494 | 
  495 |   await page.locator("#native-first").evaluate((dialog: HTMLDialogElement) => dialog.close());
  496 |   await expect(page.locator("#open-first")).toHaveAttribute("aria-expanded", "false");
  497 |   await page.locator("#open-second").click();
  498 |   await expect(page.locator("#native-first")).toHaveJSProperty("open", false);
  499 |   await expect(page.locator("#native-second")).toHaveJSProperty("open", true);
  500 |   await expect(page.locator("#open-second")).toHaveAttribute("aria-expanded", "true");
  501 | });
  502 | 
  503 | test("native request-close enters the animated cancel path while close remains immediate", async ({
  504 |   page,
  505 | }) => {
  506 |   await page.goto(BASIC);
  507 |   await page.evaluate(() => {
  508 |     document.body.insertAdjacentHTML(
  509 |       "beforeend",
  510 |       `
  511 |         <button id="native-open" commandfor="native-close" command="show-modal">Open</button>
  512 |         <dialog id="native-close" data-gsxui-dialog-content data-state="closed">
  513 |           <p>Closable</p>
  514 |           <button id="native-request-close" commandfor="native-close" command="request-close">Request close</button>
  515 |           <button id="native-close-now" commandfor="native-close" command="close">Close</button>
  516 |         </dialog>
  517 |       `,
  518 |     );
  519 |     (window as any).__nativeCloseRequests = [];
  520 |     document.addEventListener("gsxui:request-close", (event) => {
  521 |       (window as any).__nativeCloseRequests.push((event as CustomEvent).detail);
  522 |     });
  523 |   });
  524 | 
  525 |   const dialog = page.locator("#native-close");
  526 |   await page.locator("#native-open").click();
  527 |   await dialog.evaluate((element) => {
  528 |     element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
  529 |   });
  530 |   await page.locator("#native-request-close").click();
  531 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  532 |   await expect(dialog).toHaveJSProperty("open", true);
  533 |   await expect(dialog).toHaveJSProperty("open", false);
  534 |   expect(await page.evaluate(() => (window as any).__nativeCloseRequests)).toEqual([
  535 |     { reason: "cancel" },
  536 |   ]);
  537 | 
  538 |   await page.locator("#native-open").click();
  539 |   await page.locator("#native-close-now").click();
  540 |   await expect(dialog).toHaveJSProperty("open", false);
  541 |   expect(await page.evaluate(() => (window as any).__nativeCloseRequests)).toEqual([
  542 |     { reason: "cancel" },
  543 |   ]);
  544 | });
  545 | 
  546 | test("native lifecycle methods stamp state before toggle and bypass cancelable requests", async ({
  547 |   page,
  548 | }) => {
  549 |   await page.goto(BASIC);
  550 |   const observed = await page.locator(DIALOG).evaluate(async (dialog: HTMLDialogElement) => {
  551 |     const beforetoggle: Array<{ state: string | undefined; newState: string }> = [];
  552 |     let closes = 0;
  553 |     const requests: string[] = [];
  554 |     dialog.addEventListener("beforetoggle", (event) => {
```