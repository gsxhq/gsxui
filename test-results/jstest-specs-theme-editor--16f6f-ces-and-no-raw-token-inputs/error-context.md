# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> pickers expose accessible catalog choices and no raw token inputs
- Location: jstest/specs/theme-editor.spec.ts:141:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
  44  | }
  45  | 
  46  | async function iframeVariable(page: Page, name: string) {
  47  |   return page
  48  |     .frameLocator("[data-theme-preview-frame]")
  49  |     .locator("html")
  50  |     .evaluate((element, property) => element.style.getPropertyValue(property).trim(), `--${name}`);
  51  | }
  52  | 
  53  | async function commands(page: Page) {
  54  |   return {
  55  |     init: await page.locator('[data-theme-command="init"]').inputValue(),
  56  |     apply: await page.locator('[data-theme-command="apply"]').inputValue(),
  57  |   };
  58  | }
  59  | 
  60  | function shareFromInit(command: string) {
  61  |   const match = command.match(/^gsxui init --preset '([^']+)'$/);
  62  |   if (!match) throw new Error(`unexpected init command ${JSON.stringify(command)}`);
  63  |   return match[1];
  64  | }
  65  | 
  66  | function picker(page: Page, kind: "baseColor" | "theme" | "radius") {
  67  |   return page.locator(`[data-theme-picker="${kind}"]`);
  68  | }
  69  | 
  70  | async function choose(
  71  |   page: Page,
  72  |   kind: "baseColor" | "theme" | "radius",
  73  |   accessibleName: string,
  74  | ) {
  75  |   const control = picker(page, kind);
  76  |   const content = control.locator("[data-gsxui-popover-content]");
  77  |   if (!(await content.evaluate((element) => element.matches(":popover-open")))) {
  78  |     await control.locator("[data-theme-picker-trigger]").click();
  79  |   }
  80  |   await control.getByRole("radio", { name: accessibleName, exact: true }).click();
  81  | }
  82  | 
  83  | async function selectionValue(
  84  |   page: Page,
  85  |   kind: "baseColor" | "theme" | "radius",
  86  | ) {
  87  |   return picker(page, kind).locator("[data-theme-selection-value]").textContent();
  88  | }
  89  | 
  90  | test("theme editor downloads the exact variables-only theme.css", async ({ page }) => {
  91  |   await page.goto("/theme");
  92  |   const download = await downloadText(page, "css");
  93  |   expect(download.filename).toBe("theme.css");
  94  |   const defaults = readFileSync(`${repoRoot}/assets/css/themes/default.css`, "utf8");
  95  | 
  96  |   expect(variables(download.text, ":root")).toEqual(variables(defaults, ":root"));
  97  |   expect(variables(download.text, ".dark")).toEqual(variables(defaults, ".dark"));
  98  |   expect(download.text).not.toMatch(
  99  |     /@import|@theme|@layer|tailwindcss|foundation|style\.css/,
  100 |   );
  101 | });
  102 | 
  103 | test("style and mode remain independent from the curated palette", async ({ page }) => {
  104 |   await page.goto("/theme");
  105 |   await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  106 | 
  107 |   const preview = page.frameLocator("[data-theme-preview-frame]");
  108 |   const nova = preview.locator('[data-theme-preview-style="nova"]');
  109 |   const maia = preview.locator('[data-theme-preview-style="maia"]');
  110 |   await expect(nova).not.toHaveAttribute("hidden", "");
  111 |   await expect(maia).toHaveAttribute("hidden", "");
  112 | 
  113 |   const novaGeometry = await nova
  114 |     .getByRole("button", { name: "Default" })
  115 |     .first()
  116 |     .evaluate((element) => {
  117 |       const style = getComputedStyle(element);
  118 |       return { height: style.height, radius: style.borderRadius };
  119 |     });
  120 | 
  121 |   await page.locator('[data-theme-style="maia"]').click();
  122 |   await expect(nova).toHaveAttribute("hidden", "");
  123 |   await expect(maia).not.toHaveAttribute("hidden", "");
  124 |   const maiaGeometry = await maia
  125 |     .getByRole("button", { name: "Default" })
  126 |     .first()
  127 |     .evaluate((element) => {
  128 |       const style = getComputedStyle(element);
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
> 144 |   await page.goto("/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  229 |   await page.goto("/theme");
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
```