# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sonner.spec.ts >> same-turn server insertion moves imperative ownership off the fallback
- Location: jstest/specs/sonner.spec.ts:542:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types", waiting until "load"

```

# Test source

```ts
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
  469 |   await page.goto(route);
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
> 545 |   await page.goto(route);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  570 |       nativeClearTimeout(id);
  571 |     }) as typeof window.clearTimeout;
  572 | 
  573 |     const fallbackRegion = document.querySelector<HTMLElement>(
  574 |       "[data-gsxui-toaster]",
  575 |     )!;
  576 |     const section = document.createElement("section");
  577 |     section.setAttribute("aria-label", "Notifications");
  578 |     section.tabIndex = -1;
  579 |     const serverRegion = document.createElement("ol");
  580 |     serverRegion.id = "gsxui-toaster";
  581 |     serverRegion.setAttribute("data-gsxui-slot-toaster", "");
  582 |     serverRegion.setAttribute("data-gsxui-toaster", "");
  583 |     section.append(serverRegion);
  584 |     document.body.append(section);
  585 | 
  586 |     const probe = {
  587 |       activeTimers,
  588 |       get createdTimers() {
  589 |         return createdTimers;
  590 |       },
  591 |       actionCalls: 0,
  592 |       actionEvents: 0,
  593 |       directActionEvents: 0,
  594 |       fallbackRegion,
  595 |       expectedRegion: serverRegion,
  596 |       row: null as HTMLElement | null,
  597 |       actionButton: null as HTMLButtonElement | null,
  598 |     };
  599 |     (window as any).__sonnerSameTurnProbe = probe;
  600 |     document.addEventListener("gsxui:toast-action", () => {
  601 |       probe.actionEvents++;
  602 |     });
  603 | 
  604 |     window.gsxui.toast.success("Same-turn server", {
  605 |       duration: 60_000,
  606 |       action: {
  607 |         label: "Act",
  608 |         onClick: () => {
  609 |           probe.actionCalls++;
  610 |         },
  611 |       },
  612 |     });
  613 |     probe.row = serverRegion.querySelector("li[data-gsxui-toast]")!;
  614 |     probe.row.dataset.sameTurnToast = "server";
  615 |     probe.actionButton = probe.row.querySelector("[data-gsxui-toast-action]")!;
  616 |   });
  617 | 
  618 |   await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);
  619 |   expect(
  620 |     await page.evaluate(() => {
  621 |       const probe = (window as any).__sonnerSameTurnProbe;
  622 |       return {
  623 |         regionConnected: probe.expectedRegion.isConnected,
  624 |         fallbackConnected: probe.fallbackRegion.isConnected,
  625 |       };
  626 |     }),
  627 |   ).toEqual({
  628 |     regionConnected: true,
  629 |     fallbackConnected: false,
  630 |   });
  631 |   await expectSameTurnToastLifecycle(page, "server");
  632 | });
  633 | 
  634 | test("controls, promise morph, queue, expansion, and dismiss keep working", async ({
  635 |   page,
  636 | }) => {
  637 |   await page.goto(route);
  638 |   await page.evaluate(() => {
  639 |     window.__sonnerAction = 0;
  640 |     window.__sonnerCancel = 0;
  641 |     const api = window.gsxui.toast;
  642 |     api("Action", {
  643 |       duration: 60_000,
  644 |       action: { label: "Run", onClick: () => window.__sonnerAction++ },
  645 |     });
```