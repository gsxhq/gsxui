# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dialog.spec.ts >> native lifecycle methods stamp state before toggle and bypass cancelable requests
- Location: jstest/specs/dialog.spec.ts:546:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dialog/basic", waiting until "load"

```

# Test source

```ts
  449 |   expect(await page.locator(`#${identity[1].labelledby}`).count()).toBe(1);
  450 |   expect(await page.locator(`#${identity[1].describedby}`).count()).toBe(1);
  451 | });
  452 | 
  453 | test("infinite child and dialog animations do not keep a request close open", async ({ page }) => {
  454 |   await page.goto(BASIC);
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
> 549 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  550 |   const observed = await page.locator(DIALOG).evaluate(async (dialog: HTMLDialogElement) => {
  551 |     const beforetoggle: Array<{ state: string | undefined; newState: string }> = [];
  552 |     let closes = 0;
  553 |     const requests: string[] = [];
  554 |     dialog.addEventListener("beforetoggle", (event) => {
  555 |       beforetoggle.push({
  556 |         state: dialog.dataset.state,
  557 |         newState: (event as ToggleEvent).newState,
  558 |       });
  559 |     });
  560 |     dialog.addEventListener("gsxui:close", () => closes++);
  561 |     document.addEventListener("gsxui:request-open", (event) => {
  562 |       requests.push(event.type);
  563 |       event.preventDefault();
  564 |     });
  565 |     document.addEventListener("gsxui:request-close", (event) => {
  566 |       requests.push(event.type);
  567 |       event.preventDefault();
  568 |     });
  569 | 
  570 |     const nextTask = () => new Promise((resolve) => setTimeout(resolve, 0));
  571 |     dialog.showModal();
  572 |     await nextTask();
  573 |     dialog.close();
  574 |     await nextTask();
  575 |     dialog.show();
  576 |     await nextTask();
  577 |     dialog.close();
  578 |     await nextTask();
  579 |     return { beforetoggle, closes, requests };
  580 |   });
  581 | 
  582 |   expect(observed).toEqual({
  583 |     beforetoggle: [
  584 |       { state: "open", newState: "open" },
  585 |       { state: "closed", newState: "closed" },
  586 |       { state: "open", newState: "open" },
  587 |       { state: "closed", newState: "closed" },
  588 |     ],
  589 |     closes: 2,
  590 |     requests: [],
  591 |   });
  592 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  593 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  594 | });
  595 | 
  596 | test("direct requestClose follows the cancelable animated request path", async ({ page }) => {
  597 |   await page.goto(BASIC);
  598 |   const dialog = page.locator(DIALOG);
  599 |   await page.evaluate(() => {
  600 |     (window as any).__directRequestClose = [];
  601 |     document.addEventListener("gsxui:request-close", (event) => {
  602 |       (window as any).__directRequestClose.push((event as CustomEvent).detail);
  603 |     });
  604 |   });
  605 | 
  606 |   await dialog.evaluate((element: HTMLDialogElement) => element.showModal());
  607 |   await dialog.evaluate((element) => {
  608 |     element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
  609 |     (element as HTMLDialogElement).requestClose();
  610 |   });
  611 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  612 |   await expect(dialog).toHaveJSProperty("open", true);
  613 |   await expect(dialog).toHaveJSProperty("open", false);
  614 |   expect(await page.evaluate(() => (window as any).__directRequestClose)).toEqual([
  615 |     { reason: "cancel" },
  616 |   ]);
  617 | });
  618 | 
```