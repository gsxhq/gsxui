# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sonner.spec.ts >> same-turn server removal binds imperative ownership to the fallback
- Location: jstest/specs/sonner.spec.ts:466:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types", waiting until "load"

```

# Test source

```ts
  369 |       .click();
  370 |     probe.oldRows = [
  371 |       ...probe.region.querySelectorAll("li[data-gsxui-toast]"),
  372 |     ];
  373 |     probe.staleRow = moved;
  374 |     probe.staleActionEvents = 0;
  375 |     moved.addEventListener("gsxui:toast-action", () => {
  376 |       probe.staleActionEvents++;
  377 |     });
  378 |     probe.section.remove();
  379 |   });
  380 |   await expect
  381 |     .poll(() =>
  382 |       page.evaluate(() => (window as any).__sonnerLifecycleProbe.actionEvents),
  383 |     )
  384 |     .toBe(1);
  385 | 
  386 |   const fallback = page.locator("#gsxui-toaster");
  387 |   await expect(fallback).toHaveAttribute("data-gsxui-slot-toaster", "");
  388 |   await expect
  389 |     .poll(() =>
  390 |       page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
  391 |     )
  392 |     .toBe(0);
  393 |   await expect
  394 |     .poll(() =>
  395 |       page.evaluate(
  396 |         () => (window as any).__sonnerTimerProbe.removalCaps.size,
  397 |       ),
  398 |     )
  399 |     .toBe(0);
  400 |   expect(
  401 |     await page.evaluate(() => {
  402 |       const probe = (window as any).__sonnerLifecycleProbe;
  403 |       return probe.oldRows.map((row: HTMLElement) => ({
  404 |         connected: row.isConnected,
  405 |         parent: row.parentNode,
  406 |       }));
  407 |     }),
  408 |   ).toEqual([
  409 |     { connected: false, parent: null },
  410 |     { connected: false, parent: null },
  411 |     { connected: false, parent: null },
  412 |     { connected: false, parent: null },
  413 |   ]);
  414 |   await page.evaluate(() => {
  415 |     const probe = (window as any).__sonnerLifecycleProbe;
  416 |     probe.staleRow
  417 |       .querySelector<HTMLButtonElement>("[data-gsxui-toast-action]")!
  418 |       .click();
  419 |   });
  420 |   expect(
  421 |     await page.evaluate(
  422 |       () => (window as any).__sonnerLifecycleProbe.staleActionEvents,
  423 |     ),
  424 |   ).toBe(0);
  425 | 
  426 |   await page.evaluate(() => {
  427 |     const probe = (window as any).__sonnerLifecycleProbe;
  428 |     const section = document.createElement("section");
  429 |     section.setAttribute("aria-label", "Notifications");
  430 |     section.tabIndex = -1;
  431 |     const region = document.createElement("ol");
  432 |     region.id = "gsxui-toaster";
  433 |     region.setAttribute("data-gsxui-slot-toaster", "");
  434 |     region.setAttribute("data-gsxui-toaster", "");
  435 |     region.append(probe.clone("replacement-preexisting"));
  436 |     section.append(region);
  437 |     document.body.append(section);
  438 |     probe.replacementSection = section;
  439 |   });
  440 | 
  441 |   const replacement = page.locator('[data-lifecycle-probe="replacement-preexisting"]');
  442 |   await expect(replacement).toHaveAttribute("data-state", "open");
  443 |   await expect(page.locator("#gsxui-toaster")).toHaveCount(1);
  444 |   await expect
  445 |     .poll(() =>
  446 |       page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
  447 |     )
  448 |     .toBe(1);
  449 |   await expect
  450 |     .poll(() =>
  451 |       page.evaluate(() => (window as any).__sonnerTimerProbe.created),
  452 |     )
  453 |     .toBe(6);
  454 | 
  455 |   await page.evaluate(() => {
  456 |     (window as any).__sonnerLifecycleProbe.replacementSection.remove();
  457 |   });
  458 |   await expect(page.locator("#gsxui-toaster")).toHaveCount(1);
  459 |   await expect
  460 |     .poll(() =>
  461 |       page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
  462 |     )
  463 |     .toBe(0);
  464 | });
  465 | 
  466 | test("same-turn server removal binds imperative ownership to the fallback", async ({
  467 |   page,
  468 | }) => {
> 469 |   await page.goto(route);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  470 |   await page.evaluate(() => {
  471 |     const nativeSetTimeout = window.setTimeout.bind(window);
  472 |     const nativeClearTimeout = window.clearTimeout.bind(window);
  473 |     const activeTimers = new Set<number>();
  474 |     let createdTimers = 0;
  475 |     window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
  476 |       const id = nativeSetTimeout(handler, timeout, ...args);
  477 |       if (timeout != null && timeout >= 59_000) {
  478 |         activeTimers.add(id);
  479 |         createdTimers++;
  480 |       }
  481 |       return id;
  482 |     }) as typeof window.setTimeout;
  483 |     window.clearTimeout = ((id?: number) => {
  484 |       if (id != null) activeTimers.delete(id);
  485 |       nativeClearTimeout(id);
  486 |     }) as typeof window.clearTimeout;
  487 | 
  488 |     const retiredRegion = document.querySelector<HTMLElement>(
  489 |       "[data-gsxui-toaster]",
  490 |     )!;
  491 |     retiredRegion.remove();
  492 | 
  493 |     const probe = {
  494 |       activeTimers,
  495 |       get createdTimers() {
  496 |         return createdTimers;
  497 |       },
  498 |       actionCalls: 0,
  499 |       actionEvents: 0,
  500 |       directActionEvents: 0,
  501 |       retiredRegion,
  502 |       expectedRegion: null as HTMLElement | null,
  503 |       row: null as HTMLElement | null,
  504 |       actionButton: null as HTMLButtonElement | null,
  505 |     };
  506 |     (window as any).__sonnerSameTurnProbe = probe;
  507 |     document.addEventListener("gsxui:toast-action", () => {
  508 |       probe.actionEvents++;
  509 |     });
  510 | 
  511 |     window.gsxui.toast.success("Same-turn fallback", {
  512 |       duration: 60_000,
  513 |       action: {
  514 |         label: "Act",
  515 |         onClick: () => {
  516 |           probe.actionCalls++;
  517 |         },
  518 |       },
  519 |     });
  520 |     probe.expectedRegion = document.querySelector("[data-gsxui-toaster]")!;
  521 |     probe.row = probe.expectedRegion.querySelector("li[data-gsxui-toast]")!;
  522 |     probe.row.dataset.sameTurnToast = "fallback";
  523 |     probe.actionButton = probe.row.querySelector("[data-gsxui-toast-action]")!;
  524 |   });
  525 | 
  526 |   await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);
  527 |   expect(
  528 |     await page.evaluate(() => {
  529 |       const probe = (window as any).__sonnerSameTurnProbe;
  530 |       return {
  531 |         regionConnected: probe.expectedRegion.isConnected,
  532 |         retiredConnected: probe.retiredRegion.isConnected,
  533 |       };
  534 |     }),
  535 |   ).toEqual({
  536 |     regionConnected: true,
  537 |     retiredConnected: false,
  538 |   });
  539 |   await expectSameTurnToastLifecycle(page, "fallback");
  540 | });
  541 | 
  542 | test("same-turn server insertion moves imperative ownership off the fallback", async ({
  543 |   page,
  544 | }) => {
  545 |   await page.goto(route);
  546 |   await page.evaluate(() => {
  547 |     const retiredRegion = document.querySelector<HTMLElement>(
  548 |       "[data-gsxui-toaster]",
  549 |     )!;
  550 |     (window as any).__sonnerRetiredRegion = retiredRegion;
  551 |     retiredRegion.remove();
  552 |   });
  553 |   await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);
  554 | 
  555 |   await page.evaluate(() => {
  556 |     const nativeSetTimeout = window.setTimeout.bind(window);
  557 |     const nativeClearTimeout = window.clearTimeout.bind(window);
  558 |     const activeTimers = new Set<number>();
  559 |     let createdTimers = 0;
  560 |     window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
  561 |       const id = nativeSetTimeout(handler, timeout, ...args);
  562 |       if (timeout != null && timeout >= 59_000) {
  563 |         activeTimers.add(id);
  564 |         createdTimers++;
  565 |       }
  566 |       return id;
  567 |     }) as typeof window.setTimeout;
  568 |     window.clearTimeout = ((id?: number) => {
  569 |       if (id != null) activeTimers.delete(id);
```