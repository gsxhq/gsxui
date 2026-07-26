import type { Page } from "@playwright/test";
import { expect, test } from "../support/fixtures";

const BASIC = "/x/dialog/basic";
const DIALOG = "dialog[data-gsxui-dialog-content]";
const TRIGGER = "[data-gsxui-dialog-trigger]";

async function dispatch(page: Page, selector: string, type: string, detail = {}) {
  return page.locator(selector).evaluate(
    (element, { type, detail }) =>
      element.dispatchEvent(
        new CustomEvent(type, { bubbles: true, cancelable: true, detail }),
      ),
    { type, detail },
  );
}

async function open(page: Page, selector = DIALOG) {
  await dispatch(page, selector, "gsxui:request-open");
  await expect(page.locator(selector)).toHaveJSProperty("open", true);
}

async function close(page: Page, selector = DIALOG) {
  await dispatch(page, selector, "gsxui:request-close");
  await expect(page.locator(selector)).toHaveJSProperty("open", false);
}

test("the proximity trigger requests one targeted open with its stable reason", async ({ page }) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    (window as any).__dialogRequests = [];
    document.addEventListener("gsxui:request-open", (event) => {
      const request = event as CustomEvent;
      (window as any).__dialogRequests.push({
        detail: request.detail,
        target: (request.target as Element).matches(
          "dialog[data-gsxui-dialog-content]",
        ),
      });
    });
  });

  await page.locator(TRIGGER).click();

  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  expect(await page.evaluate(() => (window as any).__dialogRequests)).toEqual([
    { detail: { reason: "trigger" }, target: true },
  ]);
});

test("descendant request events drive the native dialog state", async ({ page }) => {
  await page.goto(BASIC);
  const source = `${DIALOG} [data-slot="dialog-title"]`;

  await dispatch(page, source, "gsxui:request-open", { reason: "application" });
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);

  await dispatch(page, source, "gsxui:request-close", { reason: "application" });
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
});

test("a late document listener can cancel both deferred request defaults", async ({ page }) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    (window as any).__cancelledDialogRequests = [];
    document.addEventListener("gsxui:request-open", (event) => {
      event.preventDefault();
      (window as any).__cancelledDialogRequests.push(event.type);
    });
  });

  await page.locator(TRIGGER).click();
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);

  await page.evaluate(() => {
    document.addEventListener("gsxui:request-open", (event) => event.stopImmediatePropagation(), {
      capture: true,
      once: true,
    });
  });
  await page.locator(DIALOG).evaluate((dialog: HTMLDialogElement) => dialog.showModal());
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);

  await page.evaluate(() => {
    document.addEventListener("gsxui:request-close", (event) => {
      event.preventDefault();
      (window as any).__cancelledDialogRequests.push(event.type);
    });
  });
  await page.locator("[data-gsxui-dialog-close]").first().click();
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  expect(await page.evaluate(() => (window as any).__cancelledDialogRequests)).toEqual([
    "gsxui:request-open",
    "gsxui:request-close",
  ]);
});

test("close controls, Escape, and the backdrop request their stable reasons", async ({ page }) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    (window as any).__dialogCloseReasons = [];
    document.addEventListener("gsxui:request-close", (event) => {
      (window as any).__dialogCloseReasons.push((event as CustomEvent).detail.reason);
    });
  });

  await open(page);
  await page.locator("[data-gsxui-dialog-close]").first().click();
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);

  await open(page);
  await page.keyboard.press("Escape");
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);

  await open(page);
  await page.mouse.click(1, 1);
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);

  expect(await page.evaluate(() => (window as any).__dialogCloseReasons)).toEqual([
    "close-button",
    "cancel",
    "backdrop",
  ]);
});

test("the finite exit stays open while closed and emits one notification after native close", async ({
  page,
}) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    (window as any).__dialogCloses = 0;
    document.addEventListener("gsxui:close", () => ((window as any).__dialogCloses += 1));
  });

  await open(page);
  await dispatch(page, DIALOG, "gsxui:request-close");

  await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  expect(await page.evaluate(() => (window as any).__dialogCloses)).toBe(1);
});

test("an open request during exit aborts the pending close", async ({ page }) => {
  await page.goto(BASIC);

  await open(page);
  await dispatch(page, DIALOG, "gsxui:request-close");
  await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);

  await dispatch(page, DIALOG, "gsxui:request-open");
  await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "open");
  await page.waitForTimeout(350);
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
});

test("dialog identity and trigger ownership stay within their nearest roots", async ({ page }) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    document.body.insertAdjacentHTML(
      "beforeend",
      `
        <div data-gsxui-dialog id="generated-root">
          <button id="generated-trigger" data-gsxui-dialog-trigger aria-expanded="false">Generated</button>
          <dialog data-gsxui-dialog-content data-state="closed">
            <h2 data-slot="dialog-title">Generated title</h2>
            <p data-slot="dialog-description">Generated description</p>
            <button data-gsxui-dialog-close>Close generated</button>
            <div data-gsxui-dialog id="nested-root">
              <button id="nested-trigger" data-gsxui-dialog-trigger aria-expanded="false">Nested</button>
              <button id="nested-root-close" data-gsxui-dialog-close>Close nested</button>
              <dialog data-gsxui-dialog-content data-state="closed">
                <h2 data-slot="dialog-title">Nested title</h2>
                <p data-slot="dialog-description">Nested description</p>
              </dialog>
            </div>
          </dialog>
        </div>
        <div data-gsxui-dialog id="second-root">
          <button id="second-trigger" data-gsxui-dialog-trigger aria-expanded="false">Second</button>
          <dialog data-gsxui-dialog-content data-state="closed">
            <h2 data-slot="dialog-title">Second title</h2>
            <p data-slot="dialog-description">Second description</p>
            <button data-gsxui-dialog-close>Close second</button>
          </dialog>
        </div>
        <div data-gsxui-dialog id="authored-root">
          <button id="authored-trigger" data-gsxui-dialog-trigger aria-controls="keep-control" aria-expanded="false">Authored</button>
          <dialog id="authored-dialog" data-gsxui-dialog-content data-state="closed"
            aria-labelledby="authored-label" aria-describedby="authored-description">
            <h2 id="authored-label" data-slot="dialog-title">Authored title</h2>
            <p id="authored-description" data-slot="dialog-description">Authored description</p>
          </dialog>
        </div>
      `,
    );
  });

  await page.locator("#generated-trigger").click();
  const generated = page.locator("#generated-root > dialog");
  const nested = page.locator("#nested-root dialog");
  await expect(generated).toHaveJSProperty("open", true);
  await expect(nested).toHaveJSProperty("open", false);
  await expect(page.locator("#generated-trigger")).toHaveAttribute("aria-expanded", "true");
  await expect(page.locator("#nested-trigger")).toHaveAttribute("aria-expanded", "false");

  await page.evaluate(() => {
    document.addEventListener("gsxui:request-close", (event) => {
      const target = event.target as Element;
      (window as any).__nestedCloseTarget = target
        .closest("[data-gsxui-dialog]")
        ?.id;
    });
  });
  await page.locator("#nested-root-close").click();
  await expect(generated).toHaveJSProperty("open", true);
  expect(await page.evaluate(() => (window as any).__nestedCloseTarget)).toBe("nested-root");

  const stableGeneratedID = await generated.getAttribute("id");
  await page.evaluate((generatedID) => {
    document.body.insertAdjacentHTML(
      "beforeend",
      `<button id="generated-id-invoker" commandfor="${generatedID}" command="show-modal" aria-expanded="unchanged">Generated ID invoker</button>`,
    );
  }, stableGeneratedID);
  await page.locator("#generated-root > dialog > [data-gsxui-dialog-close]").click();
  await expect(generated).toHaveJSProperty("open", false);
  await expect(page.locator("#generated-id-invoker")).toHaveAttribute(
    "aria-expanded",
    "unchanged",
  );
  await page.locator("#generated-trigger").click();
  await expect(generated).toHaveAttribute("id", stableGeneratedID!);
  await expect(page.locator("#generated-id-invoker")).toHaveAttribute(
    "aria-expanded",
    "unchanged",
  );
  await page.locator("#generated-root > dialog > [data-gsxui-dialog-close]").click();
  await expect(generated).toHaveJSProperty("open", false);

  await page.locator("#second-trigger").click();
  const second = page.locator("#second-root dialog");
  await expect(second).toHaveJSProperty("open", true);
  await page.locator("#second-root [data-gsxui-dialog-close]").click();
  await expect(second).toHaveJSProperty("open", false);
  await page.locator("#authored-trigger").click();
  await expect(page.locator("#authored-dialog")).toHaveJSProperty("open", true);
  const identity = await page.evaluate(() => {
    const roots = [
      document.querySelector("#generated-root")!,
      document.querySelector("#second-root")!,
      document.querySelector("#authored-root")!,
    ];
    return roots.map((root) => {
      const dialog = [...root.querySelectorAll("dialog[data-gsxui-dialog-content]")].find(
        (element) => element.closest("[data-gsxui-dialog]") === root,
      )!;
      const title = [...root.querySelectorAll('[data-slot="dialog-title"]')].find(
        (element) => element.closest("[data-gsxui-dialog]") === root,
      )!;
      const description = [...root.querySelectorAll('[data-slot="dialog-description"]')].find(
        (element) => element.closest("[data-gsxui-dialog]") === root,
      )!;
      const trigger = root.querySelector("[data-gsxui-dialog-trigger]")!;
      return {
        dialog: dialog.id,
        labelledby: dialog.getAttribute("aria-labelledby"),
        describedby: dialog.getAttribute("aria-describedby"),
        title: title.id,
        description: description.id,
        controls: trigger.getAttribute("aria-controls"),
      };
    });
  });

  expect(identity[0].dialog).toMatch(/^gsxui-dialog-/);
  expect(identity[0].title).toMatch(/^gsxui-title-/);
  expect(identity[0].description).toMatch(/^gsxui-desc-/);
  expect(identity[0].controls).toBe(identity[0].dialog);
  expect(identity[1].dialog).toMatch(/^gsxui-dialog-/);
  expect(identity[1].title).toMatch(/^gsxui-title-/);
  expect(identity[1].description).toMatch(/^gsxui-desc-/);
  expect(identity[1].controls).toBe(identity[1].dialog);
  expect(identity[2]).toEqual({
    dialog: "authored-dialog",
    labelledby: "authored-label",
    describedby: "authored-description",
    title: "authored-label",
    description: "authored-description",
    controls: "keep-control",
  });
  expect(new Set([identity[0].dialog, identity[1].dialog]).size).toBe(2);
  expect(new Set([identity[0].title, identity[1].title]).size).toBe(2);
  expect(new Set([identity[0].description, identity[1].description]).size).toBe(2);
  expect(await page.locator(`#${identity[0].labelledby}`).count()).toBe(1);
  expect(await page.locator(`#${identity[0].describedby}`).count()).toBe(1);
  expect(await page.locator(`#${identity[1].labelledby}`).count()).toBe(1);
  expect(await page.locator(`#${identity[1].describedby}`).count()).toBe(1);
});

test("infinite child and dialog animations do not keep a request close open", async ({ page }) => {
  await page.goto(BASIC);
  const dialog = page.locator(DIALOG);

  await open(page);
  await dialog.evaluate((element) => {
    const child = document.createElement("span");
    child.textContent = "spinning child";
    element.append(child);
    child.animate([{ opacity: 0 }, { opacity: 1 }], { duration: 10_000, iterations: Infinity });
  });
  await close(page);

  await open(page);
  await dialog.evaluate((element) => {
    element.animate([{ opacity: 0 }, { opacity: 1 }], { duration: 10_000, iterations: Infinity });
  });
  await dispatch(page, DIALOG, "gsxui:request-close");
  await expect(dialog).toHaveJSProperty("open", false, { timeout: 350 });
});

test("native invokers address authored dialogs exactly and synchronize their ARIA state", async ({
  page,
}) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    document.body.insertAdjacentHTML(
      "beforeend",
      `
        <button id="open-first" commandfor="native-first" command="show-modal">Open first</button>
        <button id="open-second" commandfor="native-second" command="show-modal">Open second</button>
        <dialog id="native-first" data-gsxui-dialog-content data-state="closed"><p>First</p></dialog>
        <dialog id="native-second" data-gsxui-dialog-content data-state="closed"><p>Second</p></dialog>
      `,
    );
  });

  await page.locator("#open-first").click();
  await expect(page.locator("#native-first")).toHaveJSProperty("open", true);
  await expect(page.locator("#native-second")).toHaveJSProperty("open", false);
  await expect(page.locator("#open-first")).toHaveAttribute("aria-expanded", "true");

  await page.locator("#native-first").evaluate((dialog: HTMLDialogElement) => dialog.close());
  await expect(page.locator("#open-first")).toHaveAttribute("aria-expanded", "false");
  await page.locator("#open-second").click();
  await expect(page.locator("#native-first")).toHaveJSProperty("open", false);
  await expect(page.locator("#native-second")).toHaveJSProperty("open", true);
  await expect(page.locator("#open-second")).toHaveAttribute("aria-expanded", "true");
});

test("native request-close enters the animated cancel path while close remains immediate", async ({
  page,
}) => {
  await page.goto(BASIC);
  await page.evaluate(() => {
    document.body.insertAdjacentHTML(
      "beforeend",
      `
        <button id="native-open" commandfor="native-close" command="show-modal">Open</button>
        <dialog id="native-close" data-gsxui-dialog-content data-state="closed">
          <p>Closable</p>
          <button id="native-request-close" commandfor="native-close" command="request-close">Request close</button>
          <button id="native-close-now" commandfor="native-close" command="close">Close</button>
        </dialog>
      `,
    );
    (window as any).__nativeCloseRequests = [];
    document.addEventListener("gsxui:request-close", (event) => {
      (window as any).__nativeCloseRequests.push((event as CustomEvent).detail);
    });
  });

  const dialog = page.locator("#native-close");
  await page.locator("#native-open").click();
  await dialog.evaluate((element) => {
    element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
  });
  await page.locator("#native-request-close").click();
  await expect(dialog).toHaveAttribute("data-state", "closed");
  await expect(dialog).toHaveJSProperty("open", true);
  await expect(dialog).toHaveJSProperty("open", false);
  expect(await page.evaluate(() => (window as any).__nativeCloseRequests)).toEqual([
    { reason: "cancel" },
  ]);

  await page.locator("#native-open").click();
  await page.locator("#native-close-now").click();
  await expect(dialog).toHaveJSProperty("open", false);
  expect(await page.evaluate(() => (window as any).__nativeCloseRequests)).toEqual([
    { reason: "cancel" },
  ]);
});

test("native lifecycle methods stamp state before toggle and bypass cancelable requests", async ({
  page,
}) => {
  await page.goto(BASIC);
  const observed = await page.locator(DIALOG).evaluate(async (dialog: HTMLDialogElement) => {
    const beforetoggle: Array<{ state: string | undefined; newState: string }> = [];
    let closes = 0;
    const requests: string[] = [];
    dialog.addEventListener("beforetoggle", (event) => {
      beforetoggle.push({
        state: dialog.dataset.state,
        newState: (event as ToggleEvent).newState,
      });
    });
    dialog.addEventListener("gsxui:close", () => closes++);
    document.addEventListener("gsxui:request-open", (event) => {
      requests.push(event.type);
      event.preventDefault();
    });
    document.addEventListener("gsxui:request-close", (event) => {
      requests.push(event.type);
      event.preventDefault();
    });

    const nextTask = () => new Promise((resolve) => setTimeout(resolve, 0));
    dialog.showModal();
    await nextTask();
    dialog.close();
    await nextTask();
    dialog.show();
    await nextTask();
    dialog.close();
    await nextTask();
    return { beforetoggle, closes, requests };
  });

  expect(observed).toEqual({
    beforetoggle: [
      { state: "open", newState: "open" },
      { state: "closed", newState: "closed" },
      { state: "open", newState: "open" },
      { state: "closed", newState: "closed" },
    ],
    closes: 2,
    requests: [],
  });
  await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
});

test("direct requestClose follows the cancelable animated request path", async ({ page }) => {
  await page.goto(BASIC);
  const dialog = page.locator(DIALOG);
  await page.evaluate(() => {
    (window as any).__directRequestClose = [];
    document.addEventListener("gsxui:request-close", (event) => {
      (window as any).__directRequestClose.push((event as CustomEvent).detail);
    });
  });

  await dialog.evaluate((element: HTMLDialogElement) => element.showModal());
  await dialog.evaluate((element) => {
    element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
    (element as HTMLDialogElement).requestClose();
  });
  await expect(dialog).toHaveAttribute("data-state", "closed");
  await expect(dialog).toHaveJSProperty("open", true);
  await expect(dialog).toHaveJSProperty("open", false);
  expect(await page.evaluate(() => (window as any).__directRequestClose)).toEqual([
    { reason: "cancel" },
  ]);
});
