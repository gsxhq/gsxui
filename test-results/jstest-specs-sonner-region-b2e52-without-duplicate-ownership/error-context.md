# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sonner.spec.ts >> region replacement reconciles nested adoption without duplicate ownership
- Location: jstest/specs/sonner.spec.ts:225:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types", waiting until "load"

```

# Test source

```ts
  128 |     sectionLabel: "Notifications",
  129 |     sectionTabIndex: "-1",
  130 |     position: "fixed",
  131 |     right: "0px",
  132 |     bottom: "0px",
  133 |     pointerEvents: "none",
  134 |   });
  135 | });
  136 | 
  137 | test("every toast type uses semantic icon color and dedicated hooks", async ({ page }) => {
  138 |   await page.goto(route);
  139 |   await page.evaluate(() => {
  140 |     const api = window.gsxui.toast;
  141 |     api("Default", { duration: 60_000 });
  142 |     api.success("Success", { duration: 60_000 });
  143 |     api.info("Info", { duration: 60_000 });
  144 |     api.warning("Warning", { duration: 60_000 });
  145 |     api.error("Error", { duration: 60_000 });
  146 |     api.loading("Loading", { duration: Infinity });
  147 |   });
  148 | 
  149 |   const cards = page.locator("li[data-gsxui-toast]");
  150 |   await expect(cards).toHaveCount(6);
  151 |   const result = await cards.evaluateAll((elements) =>
  152 |     elements.map((element) => {
  153 |       const type = (element as HTMLElement).dataset.type!;
  154 |       const icon = element.querySelector("[data-gsxui-toast-icon]");
  155 |       let semanticColor: string | null = null;
  156 |       if (!["default", "loading"].includes(type)) {
  157 |         const probe = document.createElement("span");
  158 |         probe.style.color = `var(--${type === "error" ? "destructive" : type})`;
  159 |         document.body.append(probe);
  160 |         semanticColor = getComputedStyle(probe).color;
  161 |         probe.remove();
  162 |       }
  163 |       return {
  164 |         type,
  165 |         toastSlot: element.hasAttribute("data-gsxui-slot-toast"),
  166 |         className: element.getAttribute("class"),
  167 |         toastIconSlot: icon?.hasAttribute("data-gsxui-slot-toast-icon") ?? false,
  168 |         primitiveIconSlot: icon?.hasAttribute("data-gsxui-slot-icon") ?? false,
  169 |         iconColor: icon ? getComputedStyle(icon).color : null,
  170 |         semanticColor,
  171 |       };
  172 |     }),
  173 |   );
  174 |   expect(result.map(({ type }) => type)).toEqual([
  175 |     "default",
  176 |     "success",
  177 |     "info",
  178 |     "warning",
  179 |     "error",
  180 |     "loading",
  181 |   ]);
  182 |   for (const card of result) {
  183 |     expect(card.toastSlot).toBe(true);
  184 |     expect(card.className).toBeNull();
  185 |     if (card.semanticColor) {
  186 |       expect(card.toastIconSlot).toBe(true);
  187 |       expect(card.primitiveIconSlot).toBe(true);
  188 |       expect(card.iconColor).toBe(card.semanticColor);
  189 |     }
  190 |   }
  191 | });
  192 | 
  193 | test("server rows are adopted through li[data-gsxui-toast]", async ({ page }) => {
  194 |   await page.goto(route);
  195 |   const adopted = page.locator('li[data-gsxui-toast][data-server-probe="true"]');
  196 |   await page.evaluate(() => {
  197 |     const template = document.querySelector<HTMLTemplateElement>(
  198 |       'template[data-gsxui-toast-template="success"]',
  199 |     )!;
  200 |     const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
  201 |     row.dataset.serverProbe = "true";
  202 |     row.dataset.duration = "60000";
  203 |     row.querySelector("[data-gsxui-toast-title]")!.textContent = "Server flash";
  204 |     document.querySelector("#gsxui-toaster")!.append(row);
  205 |   });
  206 | 
  207 |   await expect(adopted).toHaveAttribute("data-state", "open");
  208 |   await expect(adopted).toHaveAttribute("data-visible", "true");
  209 |   await expect(adopted.locator("[data-gsxui-toast-title]")).toHaveText("Server flash");
  210 |   expect(
  211 |     await adopted.evaluate((element) => ({
  212 |       opacity: (element as HTMLElement).style.opacity,
  213 |       transform: (element as HTMLElement).style.transform,
  214 |       zIndex: (element as HTMLElement).style.zIndex,
  215 |       pointerEvents: (element as HTMLElement).style.pointerEvents,
  216 |     })),
  217 |   ).toEqual({
  218 |     opacity: "1",
  219 |     transform: "translateY(0px) scale(1)",
  220 |     zIndex: "100",
  221 |     pointerEvents: "auto",
  222 |   });
  223 | });
  224 | 
  225 | test("region replacement reconciles nested adoption without duplicate ownership", async ({
  226 |   page,
  227 | }) => {
> 228 |   await page.goto(route);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  229 |   await page.evaluate(() => {
  230 |     const nativeSetTimeout = window.setTimeout.bind(window);
  231 |     const nativeClearTimeout = window.clearTimeout.bind(window);
  232 |     const active = new Set<number>();
  233 |     const removalCaps = new Set<number>();
  234 |     let created = 0;
  235 |     window.setTimeout = ((handler: TimerHandler, timeout?: number, ...args: unknown[]) => {
  236 |       const id = nativeSetTimeout(handler, timeout, ...args);
  237 |       if (timeout != null && timeout >= 59_000) {
  238 |         active.add(id);
  239 |         created++;
  240 |       }
  241 |       if (timeout === 600) removalCaps.add(id);
  242 |       return id;
  243 |     }) as typeof window.setTimeout;
  244 |     window.clearTimeout = ((id?: number) => {
  245 |       if (id != null) {
  246 |         active.delete(id);
  247 |         removalCaps.delete(id);
  248 |       }
  249 |       nativeClearTimeout(id);
  250 |     }) as typeof window.clearTimeout;
  251 |     (window as any).__sonnerTimerProbe = {
  252 |       active,
  253 |       removalCaps,
  254 |       get created() {
  255 |         return created;
  256 |       },
  257 |     };
  258 | 
  259 |     const oldSection = document.querySelector<HTMLElement>(
  260 |       '[aria-label="Notifications"]',
  261 |     )!;
  262 |     const templates = [...oldSection.querySelectorAll("template")];
  263 |     for (const template of templates) document.body.append(template);
  264 | 
  265 |     const clone = (label: string) => {
  266 |       const template = document.querySelector<HTMLTemplateElement>(
  267 |         'template[data-gsxui-toast-template="default"]',
  268 |       )!;
  269 |       const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
  270 |       row.dataset.duration = "60000";
  271 |       row.dataset.lifecycleProbe = label;
  272 |       row.querySelector("[data-gsxui-toast-title]")!.textContent = label;
  273 |       return row;
  274 |     };
  275 | 
  276 |     const section = document.createElement("section");
  277 |     section.setAttribute("aria-label", "Notifications");
  278 |     section.tabIndex = -1;
  279 |     const region = document.createElement("ol");
  280 |     region.id = "gsxui-toaster";
  281 |     region.setAttribute("data-gsxui-slot-toaster", "");
  282 |     region.setAttribute("data-gsxui-toaster", "");
  283 |     region.append(clone("preexisting"));
  284 |     section.append(region);
  285 |     oldSection.replaceWith(section);
  286 | 
  287 |     (window as any).__sonnerLifecycleProbe = {
  288 |       clone,
  289 |       oldSection,
  290 |       section,
  291 |       region,
  292 |       rows: [] as HTMLElement[],
  293 |       actionEvents: 0,
  294 |     };
  295 |     document.addEventListener("gsxui:toast-action", () => {
  296 |       (window as any).__sonnerLifecycleProbe.actionEvents++;
  297 |     });
  298 |   });
  299 | 
  300 |   const preexisting = page.locator('[data-lifecycle-probe="preexisting"]');
  301 |   await expect(preexisting).toHaveAttribute("data-state", "open");
  302 |   await expect(preexisting).toHaveAttribute("data-visible", "true");
  303 | 
  304 |   await page.evaluate(() => {
  305 |     const probe = (window as any).__sonnerLifecycleProbe;
  306 |     const fragment = document.createDocumentFragment();
  307 |     const direct = probe.clone("fragment-direct");
  308 |     probe.rows.push(direct);
  309 |     fragment.append(direct);
  310 |     probe.region.append(fragment);
  311 | 
  312 |     const nestedWrapper = document.createElement("div");
  313 |     const nested = probe.clone("wrapper-nested");
  314 |     probe.rows.push(nested);
  315 |     nestedWrapper.append(nested);
  316 |     probe.region.append(nestedWrapper);
  317 | 
  318 |     const laterWrapper = document.createElement("div");
  319 |     probe.region.append(laterWrapper);
  320 |     const later = probe.clone("later-descendant");
  321 |     probe.rows.push(later);
  322 |     laterWrapper.append(later);
  323 |   });
  324 | 
  325 |   for (const label of [
  326 |     "fragment-direct",
  327 |     "wrapper-nested",
  328 |     "later-descendant",
```