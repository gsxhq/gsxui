# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> joined ToggleGroup items override composed Toggle sizing and borders
- Location: jstest/specs/style-visual.spec.ts:255:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  158 |       '[data-style-contract-reference="destructive-hover"]',
  159 |       "backgroundColor",
  160 |     ),
  161 |   );
  162 | 
  163 |   const darkInvalidRing = await computed(
  164 |     '[data-style-contract-reference="dark-invalid-ring"]',
  165 |     "boxShadow",
  166 |   );
  167 |   const invalidButton = page.locator('[data-style-contract="dark-button-invalid"]');
  168 |   await invalidButton.focus();
  169 |   await finishTransitions();
  170 |   expect(
  171 |     await computed('[data-style-contract="dark-button-invalid"]', "boxShadow"),
  172 |   ).toBe(darkInvalidRing);
  173 |   const invalidBadge = page.locator('[data-style-contract="dark-badge-invalid"]');
  174 |   await invalidBadge.focus();
  175 |   await finishTransitions();
  176 |   expect(
  177 |     await computed('[data-style-contract="dark-badge-invalid"]', "boxShadow"),
  178 |   ).toBe(darkInvalidRing);
  179 | 
  180 |   const outline = page.locator('[data-style-contract="dark-button-outline"]');
  181 |   await outline.hover();
  182 |   await finishTransitions();
  183 |   expect(await computed('[data-style-contract="dark-button-outline"]', "backgroundColor")).toBe(
  184 |     await computed(
  185 |       '[data-style-contract-reference="dark-outline-hover"]',
  186 |       "backgroundColor",
  187 |     ),
  188 |   );
  189 | 
  190 |   const ghost = page.locator('[data-style-contract="dark-button-ghost"]');
  191 |   await ghost.hover();
  192 |   await finishTransitions();
  193 |   expect(await computed('[data-style-contract="dark-button-ghost"]', "backgroundColor")).toBe(
  194 |     await computed(
  195 |       '[data-style-contract-reference="dark-ghost-hover"]',
  196 |       "backgroundColor",
  197 |     ),
  198 |   );
  199 | });
  200 | 
  201 | test("Pagination edge padding overrides Button defaults and remains caller-overridable", async ({
  202 |   page,
  203 | }) => {
  204 |   const response = await page.goto("/f/style-contract");
  205 |   expect(response?.status(), "style contract fixture response").toBe(200);
  206 | 
  207 |   const padding = async (selector: string) =>
  208 |     page.locator(selector).evaluate((element) => {
  209 |       const css = getComputedStyle(element);
  210 |       return { left: css.paddingLeft, right: css.paddingRight };
  211 |     });
  212 | 
  213 |   expect(await padding('[data-style-contract="pagination-previous"]')).toEqual({
  214 |     left: "6px",
  215 |     right: "8px",
  216 |   });
  217 |   expect(await padding('[data-style-contract="pagination-next"]')).toEqual({
  218 |     left: "8px",
  219 |     right: "6px",
  220 |   });
  221 |   expect(await padding('[data-style-contract="pagination-previous-caller"]')).toEqual({
  222 |     left: "48px",
  223 |     right: "8px",
  224 |   });
  225 | });
  226 | 
  227 | test("an active invalid InputOTP slot keeps destructive border and ring semantics", async ({
  228 |   page,
  229 | }) => {
  230 |   const response = await page.goto("/f/style-contract");
  231 |   expect(response?.status(), "style contract fixture response").toBe(200);
  232 | 
  233 |   const colors = async (selector: string) =>
  234 |     page.locator(selector).evaluate((element) => {
  235 |       const css = getComputedStyle(element);
  236 |       return { borderColor: css.borderColor, boxShadow: css.boxShadow };
  237 |     });
  238 | 
  239 |   const slot = '[data-style-contract="otp-active-invalid"]';
  240 |   expect(await colors(slot)).toEqual(
  241 |     await colors('[data-style-contract-reference="otp-invalid-light"]'),
  242 |   );
  243 | 
  244 |   await page.evaluate(() => document.documentElement.classList.add("dark"));
  245 |   await page.evaluate(() => {
  246 |     for (const animation of document.getAnimations()) {
  247 |       animation.finish();
  248 |     }
  249 |   });
  250 |   expect(await colors(slot)).toEqual(
  251 |     await colors('[data-style-contract-reference="otp-invalid-dark"]'),
  252 |   );
  253 | });
  254 | 
  255 | test("joined ToggleGroup items override composed Toggle sizing and borders", async ({
  256 |   page,
  257 | }) => {
> 258 |   const response = await page.goto("/f/style-contract");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  259 |   expect(response?.status(), "style contract fixture response").toBe(200);
  260 | 
  261 |   const metrics = async (selector: string) =>
  262 |     page.locator(selector).evaluate((element) => {
  263 |       const css = getComputedStyle(element);
  264 |       const box = element.getBoundingClientRect();
  265 |       return {
  266 |         x: box.x,
  267 |         width: box.width,
  268 |         height: css.height,
  269 |         paddingLeft: css.paddingLeft,
  270 |         paddingRight: css.paddingRight,
  271 |         borderLeftWidth: css.borderLeftWidth,
  272 |         borderTopLeftRadius: css.borderTopLeftRadius,
  273 |         borderTopRightRadius: css.borderTopRightRadius,
  274 |       };
  275 |     });
  276 | 
  277 |   const first = await metrics('[data-style-contract="toggle-group-sm-first"]');
  278 |   const iconItem = await metrics('[data-style-contract="toggle-group-sm-icon"]');
  279 |   const last = await metrics('[data-style-contract="toggle-group-sm-last"]');
  280 |   expect(first).toMatchObject({
  281 |     height: "28px",
  282 |     paddingLeft: "12px",
  283 |     paddingRight: "12px",
  284 |     borderLeftWidth: "1px",
  285 |     borderTopRightRadius: "0px",
  286 |   });
  287 |   expect(iconItem).toMatchObject({
  288 |     height: "28px",
  289 |     paddingLeft: "6px",
  290 |     paddingRight: "6px",
  291 |     borderLeftWidth: "0px",
  292 |     borderTopLeftRadius: "0px",
  293 |     borderTopRightRadius: "0px",
  294 |   });
  295 |   expect(last).toMatchObject({
  296 |     height: "28px",
  297 |     paddingLeft: "12px",
  298 |     paddingRight: "12px",
  299 |     borderLeftWidth: "0px",
  300 |     borderTopLeftRadius: "0px",
  301 |   });
  302 |   expect(iconItem.x).toBeCloseTo(first.x + first.width, 5);
  303 |   expect(last.x).toBeCloseTo(iconItem.x + iconItem.width, 5);
  304 | 
  305 |   expect(
  306 |     await metrics('[data-style-contract="toggle-group-default"]'),
  307 |   ).toMatchObject({
  308 |     height: "32px",
  309 |     paddingLeft: "12px",
  310 |     paddingRight: "12px",
  311 |   });
  312 |   expect(
  313 |     await metrics('[data-style-contract="toggle-group-large"]'),
  314 |   ).toMatchObject({
  315 |     height: "36px",
  316 |     paddingLeft: "12px",
  317 |     paddingRight: "12px",
  318 |   });
  319 |   expect(
  320 |     await metrics('[data-style-contract="toggle-group-caller"]'),
  321 |   ).toMatchObject({
  322 |     height: "28px",
  323 |     paddingLeft: "32px",
  324 |     paddingRight: "32px",
  325 |   });
  326 | });
  327 | 
  328 | test("an outline FieldGroup applies the separator overlap", async ({ page }) => {
  329 |   const response = await page.goto("/f/style-contract");
  330 |   expect(response?.status(), "style contract fixture response").toBe(200);
  331 | 
  332 |   expect(
  333 |     await page
  334 |       .locator('[data-style-contract="field-outline-separator"]')
  335 |       .evaluate((element) => getComputedStyle(element).marginBottom),
  336 |   ).toBe("-8px");
  337 | });
  338 | 
  339 | test("Dialog keeps dedicated a11y hooks, semantic backdrop, and caller cascade", async ({
  340 |   page,
  341 | }) => {
  342 |   const response = await page.goto("/f/style-contract");
  343 |   expect(response?.status(), "style contract fixture response").toBe(200);
  344 | 
  345 |   await page.getByRole("button", { name: "Open contract dialog" }).click();
  346 |   const dialog = page.locator('[data-style-contract="dialog-caller"]');
  347 |   await expect(dialog).toBeVisible();
  348 |   await expect(dialog).toHaveAttribute("data-state", "open");
  349 | 
  350 |   expect(
  351 |     await dialog.evaluate((element) => {
  352 |       const title = element.querySelector("[data-gsxui-dialog-title]");
  353 |       const description = element.querySelector("[data-gsxui-dialog-description]");
  354 |       const css = getComputedStyle(element);
  355 |       return {
  356 |         labelledBy: element.getAttribute("aria-labelledby"),
  357 |         titleID: title?.id,
  358 |         describedBy: element.getAttribute("aria-describedby"),
```