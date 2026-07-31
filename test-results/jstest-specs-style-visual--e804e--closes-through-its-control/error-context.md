# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> sheet bottom opens on bottom and closes through its control
- Location: jstest/specs/style-visual.spec.ts:511:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  414 |     }),
  415 |   ).toEqual({
  416 |     display: "flex",
  417 |     side: "bottom",
  418 |     left: 0,
  419 |     right: 1280,
  420 |     bottom: 900,
  421 |     viewportWidth: 1280,
  422 |     viewportHeight: 900,
  423 |     headerAlign: "center",
  424 |   });
  425 | });
  426 | 
  427 | test("Accordion caller padding overrides the inner content default", async ({
  428 |   page,
  429 | }) => {
  430 |   const response = await page.goto("/f/style-contract");
  431 |   expect(response?.status(), "style contract fixture response").toBe(200);
  432 | 
  433 |   const content = page.locator(
  434 |     '[data-style-contract="accordion-caller-content"]',
  435 |   );
  436 |   const inner = content.locator(
  437 |     ':scope > [data-gsxui-slot-accordion-content-inner]',
  438 |   );
  439 | 
  440 |   await expect(content).toHaveAttribute("id", "accordion-caller-content");
  441 |   await expect(content).not.toHaveClass(/\bpb-8\b/);
  442 |   await expect(inner).toHaveClass(/\bpb-8\b/);
  443 |   expect(
  444 |     await content.evaluate((element) => getComputedStyle(element).paddingBottom),
  445 |   ).toBe("0px");
  446 |   expect(
  447 |     await inner.evaluate((element) => getComputedStyle(element).paddingBottom),
  448 |   ).toBe("32px");
  449 | });
  450 | 
  451 | const directionalOverlayCases = [
  452 |   {
  453 |     family: "drawer",
  454 |     fixture: "default",
  455 |     side: "bottom",
  456 |     enterProperty: "--tw-enter-translate-y",
  457 |     enterValue: "100%",
  458 |   },
  459 |   {
  460 |     family: "drawer",
  461 |     fixture: "top",
  462 |     side: "top",
  463 |     enterProperty: "--tw-enter-translate-y",
  464 |     enterValue: "-100%",
  465 |   },
  466 |   {
  467 |     family: "drawer",
  468 |     fixture: "left",
  469 |     side: "left",
  470 |     enterProperty: "--tw-enter-translate-x",
  471 |     enterValue: "-100%",
  472 |   },
  473 |   {
  474 |     family: "drawer",
  475 |     fixture: "right",
  476 |     side: "right",
  477 |     enterProperty: "--tw-enter-translate-x",
  478 |     enterValue: "100%",
  479 |   },
  480 |   {
  481 |     family: "sheet",
  482 |     fixture: "default",
  483 |     side: "right",
  484 |     enterProperty: "--tw-enter-translate-x",
  485 |     enterValue: "100%",
  486 |   },
  487 |   {
  488 |     family: "sheet",
  489 |     fixture: "left",
  490 |     side: "left",
  491 |     enterProperty: "--tw-enter-translate-x",
  492 |     enterValue: "-100%",
  493 |   },
  494 |   {
  495 |     family: "sheet",
  496 |     fixture: "top",
  497 |     side: "top",
  498 |     enterProperty: "--tw-enter-translate-y",
  499 |     enterValue: "-100%",
  500 |   },
  501 |   {
  502 |     family: "sheet",
  503 |     fixture: "bottom",
  504 |     side: "bottom",
  505 |     enterProperty: "--tw-enter-translate-y",
  506 |     enterValue: "100%",
  507 |   },
  508 | ] as const;
  509 | 
  510 | for (const overlay of directionalOverlayCases) {
  511 |   test(`${overlay.family} ${overlay.fixture} opens on ${overlay.side} and closes through its control`, async ({
  512 |     page,
  513 |   }) => {
> 514 |     const response = await page.goto("/f/style-contract");
      |                                 ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  515 |     expect(response?.status(), "style contract fixture response").toBe(200);
  516 | 
  517 |     const dialog = page.locator(
  518 |       `[data-style-contract="${overlay.family}-${overlay.fixture}"]`,
  519 |     );
  520 |     await expect(dialog).toBeHidden();
  521 |     expect(
  522 |       await dialog.evaluate((element: HTMLDialogElement) => ({
  523 |         open: element.open,
  524 |         state: element.dataset.state,
  525 |         display: getComputedStyle(element).display,
  526 |       })),
  527 |     ).toEqual({
  528 |       open: false,
  529 |       state: "closed",
  530 |       display: "none",
  531 |     });
  532 | 
  533 |     await page
  534 |       .getByRole("button", {
  535 |         name: `Open ${overlay.family} ${overlay.fixture}`,
  536 |       })
  537 |       .click();
  538 |     await expect(dialog).toBeVisible();
  539 |     await expect(dialog).toHaveAttribute("open", "");
  540 |     await expect(dialog).toHaveAttribute("data-state", "open");
  541 |     await dialog.evaluate((element) => {
  542 |       for (const animation of element.getAnimations()) animation.finish();
  543 |     });
  544 | 
  545 |     const geometry = await dialog.evaluate(
  546 |       (element, enterProperty) => {
  547 |         const rect = element.getBoundingClientRect();
  548 |         const css = getComputedStyle(element);
  549 |         return {
  550 |           side: element.getAttribute("data-side"),
  551 |           display: css.display,
  552 |           animationName: css.animationName,
  553 |           enterValue: css.getPropertyValue(enterProperty).trim(),
  554 |           left: rect.left,
  555 |           right: rect.right,
  556 |           top: rect.top,
  557 |           bottom: rect.bottom,
  558 |           width: rect.width,
  559 |           height: rect.height,
  560 |           viewportWidth: window.innerWidth,
  561 |           viewportHeight: window.innerHeight,
  562 |         };
  563 |       },
  564 |       overlay.enterProperty,
  565 |     );
  566 |     expect(geometry.side).toBe(overlay.side);
  567 |     expect(geometry.display).toBe("flex");
  568 |     expect(geometry.animationName).toContain("enter");
  569 |     expect(geometry.enterValue).toBe(overlay.enterValue);
  570 | 
  571 |     const subpixelEdgeTolerance = 0.5;
  572 |     const expectAtEdge = (actual: number, expected: number) =>
  573 |       expect(Math.abs(actual - expected)).toBeLessThanOrEqual(
  574 |         subpixelEdgeTolerance,
  575 |       );
  576 |     if (overlay.side === "top" || overlay.side === "bottom") {
  577 |       expectAtEdge(geometry.left, 0);
  578 |       expectAtEdge(geometry.right, geometry.viewportWidth);
  579 |       expectAtEdge(geometry.width, geometry.viewportWidth);
  580 |       expect(geometry.height).toBeGreaterThan(0);
  581 |       expect(geometry.height).toBeLessThan(geometry.viewportHeight);
  582 |       if (overlay.side === "top") {
  583 |         expectAtEdge(geometry.top, 0);
  584 |       } else {
  585 |         expectAtEdge(geometry.bottom, geometry.viewportHeight);
  586 |       }
  587 |     } else {
  588 |       expectAtEdge(geometry.top, 0);
  589 |       expectAtEdge(geometry.bottom, geometry.viewportHeight);
  590 |       expectAtEdge(geometry.height, geometry.viewportHeight);
  591 |       expectAtEdge(geometry.width, 384);
  592 |       if (overlay.side === "left") {
  593 |         expectAtEdge(geometry.left, 0);
  594 |       } else {
  595 |         expectAtEdge(geometry.right, geometry.viewportWidth);
  596 |       }
  597 |     }
  598 | 
  599 |     await page
  600 |       .getByRole("button", {
  601 |         name: `Close ${overlay.family} ${overlay.fixture}`,
  602 |       })
  603 |       .click();
  604 |     await expect(dialog).toHaveAttribute("data-state", "closed");
  605 |     const closing = await dialog.evaluate(
  606 |       (element: HTMLDialogElement, exitProperty) => {
  607 |         const css = getComputedStyle(element);
  608 |         return {
  609 |           open: element.open,
  610 |           animationName: css.animationName,
  611 |           exitValue: css.getPropertyValue(exitProperty).trim(),
  612 |         };
  613 |       },
  614 |       overlay.enterProperty.replace("--tw-enter-", "--tw-exit-"),
```