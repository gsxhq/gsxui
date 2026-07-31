# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> Base Color and Theme choices update both iframe modes and the same-named option
- Location: jstest/specs/theme-editor.spec.ts:226:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
  129 |       return { height: style.height, radius: style.borderRadius };
  130 |     });
  131 |   expect(maiaGeometry.height).not.toBe(novaGeometry.height);
  132 |   expect(maiaGeometry.radius).not.toBe(novaGeometry.radius);
  133 | 
  134 |   await page.locator('[data-theme-mode-tab="dark"]').click();
  135 |   await expect(preview.locator("html")).toHaveClass(/\bdark\b/);
  136 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Neutral");
  137 |   await expect(selectionValue(page, "theme")).resolves.toBe("Neutral");
  138 |   await expect(selectionValue(page, "radius")).resolves.toBe("Medium");
  139 | });
  140 | 
  141 | test("pickers expose accessible catalog choices and no raw token inputs", async ({
  142 |   page,
  143 | }) => {
  144 |   await page.goto("/theme");
  145 |   await expect(page.locator("[data-theme-var]")).toHaveCount(0);
  146 |   await expect(page.locator('[data-theme-field^="light."]')).toHaveCount(0);
  147 |   await expect(page.locator('[data-theme-field^="dark."]')).toHaveCount(0);
  148 |   await expect(page.locator('[data-theme-field="radius"]')).toHaveCount(0);
  149 |   await expect(page.locator("iframe")).toHaveCount(1);
  150 | 
  151 |   const base = picker(page, "baseColor");
  152 |   await base.locator("[data-theme-picker-trigger]").click();
  153 |   await expect(base.locator("[data-gsxui-popover-content]")).toHaveJSProperty(
  154 |     "popover",
  155 |     "auto",
  156 |   );
  157 |   expect(
  158 |     await base
  159 |       .locator("[data-gsxui-popover-content]")
  160 |       .evaluate((element) => element.matches(":popover-open")),
  161 |   ).toBe(true);
  162 |   await expect(base.getByRole("radio")).toHaveCount(7);
  163 |   await expect(base.getByRole("radio", { name: "Neutral", exact: true })).toBeChecked();
  164 | 
  165 |   const theme = picker(page, "theme");
  166 |   await theme.locator("[data-theme-picker-trigger]").click();
  167 |   await expect(theme.getByRole("radio")).toHaveCount(18);
  168 |   await expect(theme.getByRole("radio", { name: "Neutral", exact: true })).toBeChecked();
  169 |   await expect(theme.getByRole("radio", { name: "Blue", exact: true })).toHaveCount(1);
  170 | 
  171 |   const radius = picker(page, "radius");
  172 |   await radius.locator("[data-theme-picker-trigger]").click();
  173 |   await expect(radius.getByRole("radio")).toHaveCount(4);
  174 |   await expect(radius.getByRole("radio", { name: "Medium", exact: true })).toBeChecked();
  175 | 
  176 |   for (const kind of ["baseColor", "theme", "radius"] as const) {
  177 |     await expect(picker(page, kind).locator("[data-theme-selection-swatch]")).toBeVisible();
  178 |     await expect(picker(page, kind).locator("[data-theme-choice-swatch]")).not.toHaveCount(0);
  179 |   }
  180 | });
  181 | 
  182 | test("keyboard reopening focuses the checked theme-picker radio", async ({ page }) => {
  183 |   await page.goto("/theme");
  184 |   const theme = picker(page, "theme");
  185 |   const trigger = theme.locator("[data-theme-picker-trigger]");
  186 |   const content = theme.locator("[data-gsxui-popover-content]");
  187 |   const blue = theme.getByRole("radio", { name: "Blue", exact: true });
  188 | 
  189 |   await trigger.click();
  190 |   await blue.click();
  191 |   await expect(blue).toBeChecked();
  192 |   await page.keyboard.press("Escape");
  193 |   await expect(trigger).toBeFocused();
  194 | 
  195 |   await trigger.press("Enter");
  196 |   await expect(content).toHaveAttribute("data-state", "open");
  197 |   await expect(blue).toBeFocused();
  198 | });
  199 | 
  200 | test("unrelated gsxui radios stay outside the theme-picker controller", async ({ page }) => {
  201 |   await page.goto("/theme");
  202 |   const committedJSON = (await downloadText(page, "json")).text;
  203 | 
  204 |   const errorMessage = await page.evaluate(async () => {
  205 |     const radio = document.createElement("input");
  206 |     radio.type = "radio";
  207 |     radio.setAttribute("data-gsxui-slot-radio", "");
  208 |     let message = "";
  209 |     const recordError = (event: ErrorEvent) => {
  210 |       event.preventDefault();
  211 |       message = event.message;
  212 |     };
  213 |     addEventListener("error", recordError);
  214 |     document.body.append(radio);
  215 |     radio.dispatchEvent(new Event("change", { bubbles: true }));
  216 |     await new Promise((resolve) => setTimeout(resolve, 0));
  217 |     removeEventListener("error", recordError);
  218 |     radio.remove();
  219 |     return message;
  220 |   });
  221 | 
  222 |   expect(errorMessage).toBe("");
  223 |   expect((await downloadText(page, "json", false)).text).toBe(committedJSON);
  224 | });
  225 | 
  226 | test("Base Color and Theme choices update both iframe modes and the same-named option", async ({
  227 |   page,
  228 | }) => {
> 229 |   await page.goto("/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  230 |   const catalog = await schema(page);
  231 | 
  232 |   await choose(page, "baseColor", "Stone");
  233 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Stone");
  234 |   await expect(selectionValue(page, "theme")).resolves.toBe("Stone");
  235 |   await expect
  236 |     .poll(() => iframeVariable(page, "foreground"))
  237 |     .toBe(catalog.palette.resolved.stone.stone.light.foreground);
  238 | 
  239 |   const theme = picker(page, "theme");
  240 |   await theme.locator("[data-theme-picker-trigger]").click();
  241 |   await expect(theme.getByRole("radio", { name: "Stone", exact: true })).toBeChecked();
  242 |   await expect(theme.getByRole("radio", { name: "Neutral", exact: true })).toHaveCount(0);
  243 |   await expect(theme.getByRole("radio")).toHaveCount(18);
  244 | 
  245 |   await theme.getByRole("radio", { name: "Blue", exact: true }).click();
  246 |   await expect(selectionValue(page, "theme")).resolves.toBe("Blue");
  247 |   await expect
  248 |     .poll(() => iframeVariable(page, "primary"))
  249 |     .toBe(catalog.palette.resolved.stone.blue.light.primary);
  250 | 
  251 |   await page.locator('[data-theme-mode-tab="dark"]').click();
  252 |   await expect
  253 |     .poll(() => iframeVariable(page, "primary"))
  254 |     .toBe(catalog.palette.resolved.stone.blue.dark.primary);
  255 | });
  256 | 
  257 | test("desktop hover previews only the iframe and restores on dismissal or pointer leave", async ({
  258 |   page,
  259 | }) => {
  260 |   await page.goto("/theme");
  261 |   const catalog = await schema(page);
  262 |   const committedPrimary = await iframeVariable(page, "primary");
  263 |   const committedCommands = await commands(page);
  264 |   const committedJSON = (await downloadText(page, "json")).text;
  265 | 
  266 |   const theme = picker(page, "theme");
  267 |   await theme.locator("[data-theme-picker-trigger]").click();
  268 |   const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");
  269 |   await rose.locator("[data-theme-choice-swatch]").hover();
  270 | 
  271 |   await expect
  272 |     .poll(() => iframeVariable(page, "primary"))
  273 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  274 |   const roseBox = await rose.boundingBox();
  275 |   if (!roseBox) throw new Error("Rose choice has no layout box");
  276 |   await page.mouse.move(
  277 |     roseBox.x + roseBox.width - 2,
  278 |     roseBox.y + roseBox.height / 2,
  279 |   );
  280 |   await expect
  281 |     .poll(() => iframeVariable(page, "primary"))
  282 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  283 |   await expect(
  284 |     theme.getByRole("radio", { name: "Neutral", exact: true }),
  285 |   ).toBeFocused();
  286 |   expect(await commands(page)).toEqual(committedCommands);
  287 |   expect((await downloadText(page, "json", false)).text).toBe(committedJSON);
  288 | 
  289 |   await page.keyboard.press("Escape");
  290 |   await expect
  291 |     .poll(() =>
  292 |       theme
  293 |         .locator("[data-gsxui-popover-content]")
  294 |         .evaluate((element) => element.matches(":popover-open")),
  295 |     )
  296 |     .toBe(false);
  297 |   await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  298 | 
  299 |   await theme.locator("[data-theme-picker-trigger]").click();
  300 |   await rose.hover();
  301 |   await expect
  302 |     .poll(() => iframeVariable(page, "primary"))
  303 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  304 |   await page.mouse.move(1, 1);
  305 |   await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  306 |   expect(await commands(page)).toEqual(committedCommands);
  307 | });
  308 | 
  309 | test("touch pointerenter does not preview on a fine hover-capable device", async ({ page }) => {
  310 |   await page.goto("/theme");
  311 |   const committedPrimary = await iframeVariable(page, "primary");
  312 |   const theme = picker(page, "theme");
  313 |   await theme.locator("[data-theme-picker-trigger]").click();
  314 |   const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");
  315 | 
  316 |   expect(
  317 |     await page.evaluate(() =>
  318 |       matchMedia("(hover: hover) and (pointer: fine)").matches,
  319 |     ),
  320 |   ).toBe(true);
  321 |   const pointerType = await rose.evaluate((element) => {
  322 |     const event = new PointerEvent("pointerenter", {
  323 |       bubbles: false,
  324 |       composed: true,
  325 |       isPrimary: true,
  326 |       pointerId: 17,
  327 |       pointerType: "touch",
  328 |     });
  329 |     element.dispatchEvent(event);
```