# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sonner.spec.ts >> controls, promise morph, queue, expansion, and dismiss keep working
- Location: jstest/specs/sonner.spec.ts:634:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types", waiting until "load"

```

# Test source

```ts
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
> 637 |   await page.goto(route);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  638 |   await page.evaluate(() => {
  639 |     window.__sonnerAction = 0;
  640 |     window.__sonnerCancel = 0;
  641 |     const api = window.gsxui.toast;
  642 |     api("Action", {
  643 |       duration: 60_000,
  644 |       action: { label: "Run", onClick: () => window.__sonnerAction++ },
  645 |     });
  646 |     api("Cancel", {
  647 |       duration: 60_000,
  648 |       cancel: { label: "Stop", onClick: () => window.__sonnerCancel++ },
  649 |     });
  650 |   });
  651 |   await page.getByRole("button", { name: "Run" }).evaluate((button: HTMLButtonElement) => button.click());
  652 |   await page.getByRole("button", { name: "Stop" }).evaluate((button: HTMLButtonElement) => button.click());
  653 |   await expect.poll(() =>
  654 |     page.evaluate(() => [window.__sonnerAction, window.__sonnerCancel]),
  655 |   ).toEqual([1, 1]);
  656 | 
  657 |   await page.evaluate(() => {
  658 |     window.gsxui.toast.promise(Promise.resolve("done"), {
  659 |       loading: "Loading probe",
  660 |       success: (value: string) => `Resolved ${value}`,
  661 |       error: "Rejected",
  662 |     });
  663 |   });
  664 |   const promiseToast = page.locator("li[data-gsxui-toast]", {
  665 |     hasText: "Resolved done",
  666 |   });
  667 |   await expect(promiseToast).toHaveAttribute("data-type", "success");
  668 |   await expect(promiseToast.locator("[data-gsxui-toast-icon]")).toHaveCount(1);
  669 |   await page.evaluate(() => window.gsxui.toast.dismiss());
  670 |   await expect.poll(() => page.locator("li[data-gsxui-toast]").count()).toBe(0);
  671 | 
  672 |   await page.evaluate(() => {
  673 |     for (let i = 0; i < 4; i++) {
  674 |       window.gsxui.toast(`Queue ${i}`, { duration: 60_000 });
  675 |     }
  676 |   });
  677 |   const queued = page.locator("li[data-gsxui-toast]", { hasText: "Queue 0" });
  678 |   await expect(queued).toHaveAttribute("data-visible", "false");
  679 |   const front = page.locator("li[data-gsxui-toast]", { hasText: "Queue 3" });
  680 |   await front.hover();
  681 |   await expect(page.locator("#gsxui-toaster")).toHaveAttribute("data-expanded", "true");
  682 |   await expect(front).toHaveCSS("pointer-events", "auto");
  683 |   const measured = await page.locator("li[data-gsxui-toast]").evaluateAll((cards) => {
  684 |     const frontCard = cards.at(-1) as HTMLElement;
  685 |     const behind = cards.at(-2) as HTMLElement;
  686 |     return [
  687 |       frontCard.style.transform,
  688 |       behind.style.transform,
  689 |       `translateY(-${frontCard.offsetHeight + 14}px) scale(1)`,
  690 |     ];
  691 |   });
  692 |   expect(measured[0]).toBe("translateY(0px) scale(1)");
  693 |   expect(measured[1]).toBe(measured[2]);
  694 | 
  695 |   await front.getByRole("button", { name: "Close" }).click();
  696 |   await expect(front).toHaveAttribute("data-state", "closed");
  697 |   await expect(front).toHaveCount(0);
  698 |   await expect.poll(() => page.locator("li[data-gsxui-toast]").count()).toBe(3);
  699 |   await expect(queued).toHaveAttribute("data-visible", "true");
  700 | });
  701 | 
```