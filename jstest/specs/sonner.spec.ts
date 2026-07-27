import type { Page } from "@playwright/test";

import { expect, test } from "../support/fixtures";

const route = "/x/sonner/types";

async function expectSameTurnToastLifecycle(page: Page, label: string) {
  const card = page.locator(`[data-same-turn-toast="${label}"]`);
  await expect(card).toHaveAttribute("data-state", "open");
  await expect(card).toHaveAttribute("data-visible", "true");
  expect(
    await card.evaluate((element) => {
      const probe = (window as any).__sonnerSameTurnProbe;
      const row = element as HTMLElement;
      return {
        expectedParent: row.parentElement === probe.expectedRegion,
        opacity: row.style.opacity,
        transform: row.style.transform,
        pointerEvents: row.style.pointerEvents,
      };
    }),
  ).toEqual({
    expectedParent: true,
    opacity: "1",
    transform: "translateY(0px) scale(1)",
    pointerEvents: "auto",
  });
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as any).__sonnerSameTurnProbe.activeTimers.size,
      ),
    )
    .toBe(1);
  expect(
    await page.evaluate(
      () => (window as any).__sonnerSameTurnProbe.createdTimers,
    ),
  ).toBe(1);

  await page.evaluate(() => {
    const probe = (window as any).__sonnerSameTurnProbe;
    probe.row.addEventListener("gsxui:toast-action", () => {
      probe.directActionEvents++;
    });
    probe.actionButton.click();
  });
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerSameTurnProbe.actionCalls),
    )
    .toBe(1);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerSameTurnProbe.actionEvents),
    )
    .toBe(1);
  await expect(card).toHaveCount(0);
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as any).__sonnerSameTurnProbe.activeTimers.size,
      ),
    )
    .toBe(0);

  await page.evaluate(() => {
    (window as any).__sonnerSameTurnProbe.actionButton.click();
  });
  expect(
    await page.evaluate(() => {
      const probe = (window as any).__sonnerSameTurnProbe;
      return {
        actionCalls: probe.actionCalls,
        directActionEvents: probe.directActionEvents,
        connected: probe.row.isConnected,
        parent: probe.row.parentNode,
      };
    }),
  ).toEqual({
    actionCalls: 1,
    directActionEvents: 1,
    connected: false,
    parent: null,
  });
}

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
  await expect(region).toHaveAttribute("data-gsxui-slot-toaster", "");
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
        toastSlot: element.hasAttribute("data-gsxui-slot-toast"),
        className: element.getAttribute("class"),
        toastIconSlot: icon?.hasAttribute("data-gsxui-slot-toast-icon") ?? false,
        primitiveIconSlot: icon?.hasAttribute("data-gsxui-slot-icon") ?? false,
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
    expect(card.toastSlot).toBe(true);
    expect(card.className).toBeNull();
    if (card.semanticColor) {
      expect(card.toastIconSlot).toBe(true);
      expect(card.primitiveIconSlot).toBe(true);
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

test("region replacement reconciles nested adoption without duplicate ownership", async ({
  page,
}) => {
  await page.goto(route);
  await page.evaluate(() => {
    const nativeSetTimeout = window.setTimeout.bind(window);
    const nativeClearTimeout = window.clearTimeout.bind(window);
    const active = new Set<number>();
    const removalCaps = new Set<number>();
    let created = 0;
    window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
      const id = nativeSetTimeout(handler, timeout, ...args);
      if (timeout != null && timeout >= 59_000) {
        active.add(id);
        created++;
      }
      if (timeout === 600) removalCaps.add(id);
      return id;
    }) as typeof window.setTimeout;
    window.clearTimeout = ((id?: number) => {
      if (id != null) {
        active.delete(id);
        removalCaps.delete(id);
      }
      nativeClearTimeout(id);
    }) as typeof window.clearTimeout;
    (window as any).__sonnerTimerProbe = {
      active,
      removalCaps,
      get created() {
        return created;
      },
    };

    const oldSection = document.querySelector<HTMLElement>(
      '[aria-label="Notifications"]',
    )!;
    const templates = [...oldSection.querySelectorAll("template")];
    for (const template of templates) document.body.append(template);

    const clone = (label: string) => {
      const template = document.querySelector<HTMLTemplateElement>(
        'template[data-gsxui-toast-template="default"]',
      )!;
      const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
      row.dataset.duration = "60000";
      row.dataset.lifecycleProbe = label;
      row.querySelector("[data-gsxui-toast-title]")!.textContent = label;
      return row;
    };

    const section = document.createElement("section");
    section.setAttribute("aria-label", "Notifications");
    section.tabIndex = -1;
    const region = document.createElement("ol");
    region.id = "gsxui-toaster";
    region.setAttribute("data-gsxui-slot-toaster", "");
    region.setAttribute("data-gsxui-toaster", "");
    region.append(clone("preexisting"));
    section.append(region);
    oldSection.replaceWith(section);

    (window as any).__sonnerLifecycleProbe = {
      clone,
      oldSection,
      section,
      region,
      rows: [] as HTMLElement[],
      actionEvents: 0,
    };
    document.addEventListener("gsxui:toast-action", () => {
      (window as any).__sonnerLifecycleProbe.actionEvents++;
    });
  });

  const preexisting = page.locator('[data-lifecycle-probe="preexisting"]');
  await expect(preexisting).toHaveAttribute("data-state", "open");
  await expect(preexisting).toHaveAttribute("data-visible", "true");

  await page.evaluate(() => {
    const probe = (window as any).__sonnerLifecycleProbe;
    const fragment = document.createDocumentFragment();
    const direct = probe.clone("fragment-direct");
    probe.rows.push(direct);
    fragment.append(direct);
    probe.region.append(fragment);

    const nestedWrapper = document.createElement("div");
    const nested = probe.clone("wrapper-nested");
    probe.rows.push(nested);
    nestedWrapper.append(nested);
    probe.region.append(nestedWrapper);

    const laterWrapper = document.createElement("div");
    probe.region.append(laterWrapper);
    const later = probe.clone("later-descendant");
    probe.rows.push(later);
    laterWrapper.append(later);
  });

  for (const label of [
    "fragment-direct",
    "wrapper-nested",
    "later-descendant",
  ]) {
    await expect(page.locator(`[data-lifecycle-probe="${label}"]`)).toHaveAttribute(
      "data-state",
      "open",
    );
  }
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
    )
    .toBe(3);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.created),
    )
    .toBe(4);

  await page.evaluate(() => {
    const probe = (window as any).__sonnerLifecycleProbe;
    const moved = probe.rows[0] as HTMLElement;
    const wrapper = document.createElement("div");
    wrapper.append(moved);
    probe.region.append(wrapper);
  });
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
    )
    .toBe(3);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.created),
    )
    .toBe(4);

  await page.evaluate(() => {
    const probe = (window as any).__sonnerLifecycleProbe;
    const moved = probe.rows[0] as HTMLElement;
    moved
      .querySelector<HTMLButtonElement>("[data-gsxui-toast-action]")!
      .click();
    probe.oldRows = [
      ...probe.region.querySelectorAll("li[data-gsxui-toast]"),
    ];
    probe.staleRow = moved;
    probe.staleActionEvents = 0;
    moved.addEventListener("gsxui:toast-action", () => {
      probe.staleActionEvents++;
    });
    probe.section.remove();
  });
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerLifecycleProbe.actionEvents),
    )
    .toBe(1);

  const fallback = page.locator("#gsxui-toaster");
  await expect(fallback).toHaveAttribute("data-gsxui-slot-toaster", "");
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
    )
    .toBe(0);
  await expect
    .poll(() =>
      page.evaluate(
        () => (window as any).__sonnerTimerProbe.removalCaps.size,
      ),
    )
    .toBe(0);
  expect(
    await page.evaluate(() => {
      const probe = (window as any).__sonnerLifecycleProbe;
      return probe.oldRows.map((row: HTMLElement) => ({
        connected: row.isConnected,
        parent: row.parentNode,
      }));
    }),
  ).toEqual([
    { connected: false, parent: null },
    { connected: false, parent: null },
    { connected: false, parent: null },
    { connected: false, parent: null },
  ]);
  await page.evaluate(() => {
    const probe = (window as any).__sonnerLifecycleProbe;
    probe.staleRow
      .querySelector<HTMLButtonElement>("[data-gsxui-toast-action]")!
      .click();
  });
  expect(
    await page.evaluate(
      () => (window as any).__sonnerLifecycleProbe.staleActionEvents,
    ),
  ).toBe(0);

  await page.evaluate(() => {
    const probe = (window as any).__sonnerLifecycleProbe;
    const section = document.createElement("section");
    section.setAttribute("aria-label", "Notifications");
    section.tabIndex = -1;
    const region = document.createElement("ol");
    region.id = "gsxui-toaster";
    region.setAttribute("data-gsxui-slot-toaster", "");
    region.setAttribute("data-gsxui-toaster", "");
    region.append(probe.clone("replacement-preexisting"));
    section.append(region);
    document.body.append(section);
    probe.replacementSection = section;
  });

  const replacement = page.locator('[data-lifecycle-probe="replacement-preexisting"]');
  await expect(replacement).toHaveAttribute("data-state", "open");
  await expect(page.locator("#gsxui-toaster")).toHaveCount(1);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
    )
    .toBe(1);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.created),
    )
    .toBe(6);

  await page.evaluate(() => {
    (window as any).__sonnerLifecycleProbe.replacementSection.remove();
  });
  await expect(page.locator("#gsxui-toaster")).toHaveCount(1);
  await expect
    .poll(() =>
      page.evaluate(() => (window as any).__sonnerTimerProbe.active.size),
    )
    .toBe(0);
});

test("same-turn server removal binds imperative ownership to the fallback", async ({
  page,
}) => {
  await page.goto(route);
  await page.evaluate(() => {
    const nativeSetTimeout = window.setTimeout.bind(window);
    const nativeClearTimeout = window.clearTimeout.bind(window);
    const activeTimers = new Set<number>();
    let createdTimers = 0;
    window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
      const id = nativeSetTimeout(handler, timeout, ...args);
      if (timeout != null && timeout >= 59_000) {
        activeTimers.add(id);
        createdTimers++;
      }
      return id;
    }) as typeof window.setTimeout;
    window.clearTimeout = ((id?: number) => {
      if (id != null) activeTimers.delete(id);
      nativeClearTimeout(id);
    }) as typeof window.clearTimeout;

    const retiredRegion = document.querySelector<HTMLElement>(
      "[data-gsxui-toaster]",
    )!;
    retiredRegion.remove();

    const probe = {
      activeTimers,
      get createdTimers() {
        return createdTimers;
      },
      actionCalls: 0,
      actionEvents: 0,
      directActionEvents: 0,
      retiredRegion,
      expectedRegion: null as HTMLElement | null,
      row: null as HTMLElement | null,
      actionButton: null as HTMLButtonElement | null,
    };
    (window as any).__sonnerSameTurnProbe = probe;
    document.addEventListener("gsxui:toast-action", () => {
      probe.actionEvents++;
    });

    window.gsxui.toast.success("Same-turn fallback", {
      duration: 60_000,
      action: {
        label: "Act",
        onClick: () => {
          probe.actionCalls++;
        },
      },
    });
    probe.expectedRegion = document.querySelector("[data-gsxui-toaster]")!;
    probe.row = probe.expectedRegion.querySelector("li[data-gsxui-toast]")!;
    probe.row.dataset.sameTurnToast = "fallback";
    probe.actionButton = probe.row.querySelector("[data-gsxui-toast-action]")!;
  });

  await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);
  expect(
    await page.evaluate(() => {
      const probe = (window as any).__sonnerSameTurnProbe;
      return {
        regionConnected: probe.expectedRegion.isConnected,
        retiredConnected: probe.retiredRegion.isConnected,
      };
    }),
  ).toEqual({
    regionConnected: true,
    retiredConnected: false,
  });
  await expectSameTurnToastLifecycle(page, "fallback");
});

test("same-turn server insertion moves imperative ownership off the fallback", async ({
  page,
}) => {
  await page.goto(route);
  await page.evaluate(() => {
    const retiredRegion = document.querySelector<HTMLElement>(
      "[data-gsxui-toaster]",
    )!;
    (window as any).__sonnerRetiredRegion = retiredRegion;
    retiredRegion.remove();
  });
  await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);

  await page.evaluate(() => {
    const nativeSetTimeout = window.setTimeout.bind(window);
    const nativeClearTimeout = window.clearTimeout.bind(window);
    const activeTimers = new Set<number>();
    let createdTimers = 0;
    window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
      const id = nativeSetTimeout(handler, timeout, ...args);
      if (timeout != null && timeout >= 59_000) {
        activeTimers.add(id);
        createdTimers++;
      }
      return id;
    }) as typeof window.setTimeout;
    window.clearTimeout = ((id?: number) => {
      if (id != null) activeTimers.delete(id);
      nativeClearTimeout(id);
    }) as typeof window.clearTimeout;

    const fallbackRegion = document.querySelector<HTMLElement>(
      "[data-gsxui-toaster]",
    )!;
    const section = document.createElement("section");
    section.setAttribute("aria-label", "Notifications");
    section.tabIndex = -1;
    const serverRegion = document.createElement("ol");
    serverRegion.id = "gsxui-toaster";
    serverRegion.setAttribute("data-gsxui-slot-toaster", "");
    serverRegion.setAttribute("data-gsxui-toaster", "");
    section.append(serverRegion);
    document.body.append(section);

    const probe = {
      activeTimers,
      get createdTimers() {
        return createdTimers;
      },
      actionCalls: 0,
      actionEvents: 0,
      directActionEvents: 0,
      fallbackRegion,
      expectedRegion: serverRegion,
      row: null as HTMLElement | null,
      actionButton: null as HTMLButtonElement | null,
    };
    (window as any).__sonnerSameTurnProbe = probe;
    document.addEventListener("gsxui:toast-action", () => {
      probe.actionEvents++;
    });

    window.gsxui.toast.success("Same-turn server", {
      duration: 60_000,
      action: {
        label: "Act",
        onClick: () => {
          probe.actionCalls++;
        },
      },
    });
    probe.row = serverRegion.querySelector("li[data-gsxui-toast]")!;
    probe.row.dataset.sameTurnToast = "server";
    probe.actionButton = probe.row.querySelector("[data-gsxui-toast-action]")!;
  });

  await expect(page.locator("[data-gsxui-toaster]")).toHaveCount(1);
  expect(
    await page.evaluate(() => {
      const probe = (window as any).__sonnerSameTurnProbe;
      return {
        regionConnected: probe.expectedRegion.isConnected,
        fallbackConnected: probe.fallbackRegion.isConnected,
      };
    }),
  ).toEqual({
    regionConnected: true,
    fallbackConnected: false,
  });
  await expectSameTurnToastLifecycle(page, "server");
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
