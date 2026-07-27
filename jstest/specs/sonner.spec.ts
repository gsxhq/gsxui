import { expect, test } from "../support/fixtures";

const route = "/x/sonner/types";

test("fallback toaster matches the server region contract", async ({ page }) => {
  await page.goto(route);
  await page.evaluate(() => {
    const section = document.querySelector('[aria-label="Notifications"]')!;
    for (const template of [...section.querySelectorAll("template")]) {
      document.body.append(template);
    }
    section.remove();
    window.gsxui.toast.success("Fallback toast", { duration: 60_000 });
    const template = document.querySelector<HTMLTemplateElement>(
      'template[data-gsxui-toast-template="info"]',
    )!;
    const serverRow = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
    serverRow.dataset.fallbackServer = "true";
    serverRow.dataset.duration = "60000";
    document.querySelector("#gsxui-toaster")!.append(serverRow);
  });

  const region = page.locator("#gsxui-toaster");
  await expect(region).toHaveAttribute("data-gsxui-toaster", "");
  await expect(region).toHaveAttribute("data-gsxui-slot", "toaster");
  await expect(region).not.toHaveAttribute("class", /.+/);
  await expect(region.locator("li[data-gsxui-toast]")).toHaveCount(2);
  await expect(
    region.locator('li[data-gsxui-toast][data-fallback-server="true"]'),
  ).toHaveAttribute("data-state", "open");
  expect(
    await region.evaluate((element) => {
      const section = element.parentElement!;
      const css = getComputedStyle(element);
      return {
        sectionLabel: section.getAttribute("aria-label"),
        sectionTabIndex: section.getAttribute("tabindex"),
        position: css.position,
        right: css.right,
        bottom: css.bottom,
        pointerEvents: css.pointerEvents,
      };
    }),
  ).toEqual({
    sectionLabel: "Notifications",
    sectionTabIndex: "-1",
    position: "fixed",
    right: "0px",
    bottom: "0px",
    pointerEvents: "none",
  });
});

test("every toast type uses semantic icon color and dedicated hooks", async ({ page }) => {
  await page.goto(route);
  await page.evaluate(() => {
    const api = window.gsxui.toast;
    api("Default", { duration: 60_000 });
    api.success("Success", { duration: 60_000 });
    api.info("Info", { duration: 60_000 });
    api.warning("Warning", { duration: 60_000 });
    api.error("Error", { duration: 60_000 });
    api.loading("Loading", { duration: Infinity });
  });

  const cards = page.locator("li[data-gsxui-toast]");
  await expect(cards).toHaveCount(6);
  const result = await cards.evaluateAll((elements) =>
    elements.map((element) => {
      const type = (element as HTMLElement).dataset.type!;
      const icon = element.querySelector("[data-gsxui-toast-icon]");
      let semanticColor: string | null = null;
      if (!["default", "loading"].includes(type)) {
        const probe = document.createElement("span");
        probe.style.color = `var(--${type === "error" ? "destructive" : type})`;
        document.body.append(probe);
        semanticColor = getComputedStyle(probe).color;
        probe.remove();
      }
      return {
        type,
        slot: element.getAttribute("data-gsxui-slot"),
        className: element.getAttribute("class"),
        iconSlot: icon?.getAttribute("data-gsxui-slot") ?? null,
        iconColor: icon ? getComputedStyle(icon).color : null,
        semanticColor,
      };
    }),
  );
  expect(result.map(({ type }) => type)).toEqual([
    "default",
    "success",
    "info",
    "warning",
    "error",
    "loading",
  ]);
  for (const card of result) {
    expect(card.slot).toBe("toast");
    expect(card.className).toBeNull();
    if (card.semanticColor) {
      expect(card.iconSlot).toBe("icon toast-icon");
      expect(card.iconColor).toBe(card.semanticColor);
    }
  }
});

test("server rows are adopted through li[data-gsxui-toast]", async ({ page }) => {
  await page.goto(route);
  const adopted = page.locator('li[data-gsxui-toast][data-server-probe="true"]');
  await page.evaluate(() => {
    const template = document.querySelector<HTMLTemplateElement>(
      'template[data-gsxui-toast-template="success"]',
    )!;
    const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
    row.dataset.serverProbe = "true";
    row.dataset.duration = "60000";
    row.querySelector("[data-gsxui-toast-title]")!.textContent = "Server flash";
    document.querySelector("#gsxui-toaster")!.append(row);
  });

  await expect(adopted).toHaveAttribute("data-state", "open");
  await expect(adopted).toHaveAttribute("data-visible", "true");
  await expect(adopted.locator("[data-gsxui-toast-title]")).toHaveText("Server flash");
  expect(
    await adopted.evaluate((element) => ({
      opacity: (element as HTMLElement).style.opacity,
      transform: (element as HTMLElement).style.transform,
      zIndex: (element as HTMLElement).style.zIndex,
      pointerEvents: (element as HTMLElement).style.pointerEvents,
    })),
  ).toEqual({
    opacity: "1",
    transform: "translateY(0px) scale(1)",
    zIndex: "100",
    pointerEvents: "auto",
  });
});

test("controls, promise morph, queue, expansion, and dismiss keep working", async ({
  page,
}) => {
  await page.goto(route);
  await page.evaluate(() => {
    window.__sonnerAction = 0;
    window.__sonnerCancel = 0;
    const api = window.gsxui.toast;
    api("Action", {
      duration: 60_000,
      action: { label: "Run", onClick: () => window.__sonnerAction++ },
    });
    api("Cancel", {
      duration: 60_000,
      cancel: { label: "Stop", onClick: () => window.__sonnerCancel++ },
    });
  });
  await page.getByRole("button", { name: "Run" }).evaluate((button: HTMLButtonElement) => button.click());
  await page.getByRole("button", { name: "Stop" }).evaluate((button: HTMLButtonElement) => button.click());
  await expect.poll(() =>
    page.evaluate(() => [window.__sonnerAction, window.__sonnerCancel]),
  ).toEqual([1, 1]);

  await page.evaluate(() => {
    window.gsxui.toast.promise(Promise.resolve("done"), {
      loading: "Loading probe",
      success: (value: string) => `Resolved ${value}`,
      error: "Rejected",
    });
  });
  const promiseToast = page.locator("li[data-gsxui-toast]", {
    hasText: "Resolved done",
  });
  await expect(promiseToast).toHaveAttribute("data-type", "success");
  await expect(promiseToast.locator("[data-gsxui-toast-icon]")).toHaveCount(1);
  await page.evaluate(() => window.gsxui.toast.dismiss());
  await expect.poll(() => page.locator("li[data-gsxui-toast]").count()).toBe(0);

  await page.evaluate(() => {
    for (let i = 0; i < 4; i++) {
      window.gsxui.toast(`Queue ${i}`, { duration: 60_000 });
    }
  });
  const queued = page.locator("li[data-gsxui-toast]", { hasText: "Queue 0" });
  await expect(queued).toHaveAttribute("data-visible", "false");
  const front = page.locator("li[data-gsxui-toast]", { hasText: "Queue 3" });
  await front.hover();
  await expect(page.locator("#gsxui-toaster")).toHaveAttribute("data-expanded", "true");
  await expect(front).toHaveCSS("pointer-events", "auto");
  const measured = await page.locator("li[data-gsxui-toast]").evaluateAll((cards) => {
    const frontCard = cards.at(-1) as HTMLElement;
    const behind = cards.at(-2) as HTMLElement;
    return [
      frontCard.style.transform,
      behind.style.transform,
      `translateY(-${frontCard.offsetHeight + 14}px) scale(1)`,
    ];
  });
  expect(measured[0]).toBe("translateY(0px) scale(1)");
  expect(measured[1]).toBe(measured[2]);

  await front.getByRole("button", { name: "Close" }).click();
  await expect(front).toHaveAttribute("data-state", "closed");
  await expect(front).toHaveCount(0);
  await expect.poll(() => page.locator("li[data-gsxui-toast]").count()).toBe(3);
  await expect(queued).toHaveAttribute("data-visible", "true");
});
