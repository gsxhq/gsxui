# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/basic-demo-presentation.spec.ts >> calendar/datepicker presents one styled control per trigger
- Location: jstest/specs/basic-demo-presentation.spec.ts:110:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/datepicker", waiting until "load"

```

# Test source

```ts
  13  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  14  |     triggerSlot: "alert-dialog-trigger",
  15  |     count: 1,
  16  |   },
  17  |   {
  18  |     route: "dialog/basic",
  19  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  20  |     triggerSlot: "dialog-trigger",
  21  |     count: 1,
  22  |   },
  23  |   {
  24  |     route: "dialog/footer",
  25  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  26  |     triggerSlot: "dialog-trigger",
  27  |     count: 1,
  28  |   },
  29  |   {
  30  |     route: "drawer/basic",
  31  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  32  |     triggerSlot: "drawer-trigger",
  33  |     count: 1,
  34  |   },
  35  |   {
  36  |     route: "drawer/directions",
  37  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  38  |     triggerSlot: "drawer-trigger",
  39  |     count: 4,
  40  |   },
  41  |   {
  42  |     route: "dropdown-menu/basic",
  43  |     behaviorAttribute: "data-gsxui-dropdown-trigger",
  44  |     triggerSlot: "dropdown-menu-trigger",
  45  |     count: 1,
  46  |   },
  47  |   {
  48  |     route: "dropdown-menu/checkboxes",
  49  |     behaviorAttribute: "data-gsxui-dropdown-trigger",
  50  |     triggerSlot: "dropdown-menu-trigger",
  51  |     count: 1,
  52  |   },
  53  |   {
  54  |     route: "dropdown-menu/destructive",
  55  |     behaviorAttribute: "data-gsxui-dropdown-trigger",
  56  |     triggerSlot: "dropdown-menu-trigger",
  57  |     count: 1,
  58  |   },
  59  |   {
  60  |     route: "dropdown-menu/radios",
  61  |     behaviorAttribute: "data-gsxui-dropdown-trigger",
  62  |     triggerSlot: "dropdown-menu-trigger",
  63  |     count: 1,
  64  |   },
  65  |   {
  66  |     route: "dropdown-menu/submenu",
  67  |     behaviorAttribute: "data-gsxui-dropdown-trigger",
  68  |     triggerSlot: "dropdown-menu-trigger",
  69  |     count: 1,
  70  |   },
  71  |   {
  72  |     route: "calendar/datepicker",
  73  |     behaviorAttribute: "data-gsxui-popover-trigger",
  74  |     triggerSlot: "popover-trigger",
  75  |     count: 1,
  76  |   },
  77  |   {
  78  |     route: "popover/basic",
  79  |     behaviorAttribute: "data-gsxui-popover-trigger",
  80  |     triggerSlot: "popover-trigger",
  81  |     count: 1,
  82  |   },
  83  |   {
  84  |     route: "sheet/basic",
  85  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  86  |     triggerSlot: "sheet-trigger",
  87  |     count: 1,
  88  |   },
  89  |   {
  90  |     route: "sheet/directions",
  91  |     behaviorAttribute: "data-gsxui-dialog-trigger",
  92  |     triggerSlot: "sheet-trigger",
  93  |     count: 4,
  94  |   },
  95  |   {
  96  |     route: "tooltip/basic",
  97  |     behaviorAttribute: "data-gsxui-tooltip-trigger",
  98  |     triggerSlot: "tooltip-trigger",
  99  |     count: 1,
  100 |   },
  101 |   {
  102 |     route: "tooltip/wide",
  103 |     behaviorAttribute: "data-gsxui-tooltip-trigger",
  104 |     triggerSlot: "tooltip-trigger",
  105 |     count: 1,
  106 |   },
  107 | ];
  108 | 
  109 | for (const demo of triggerDemos) {
  110 |   test(`${demo.route} presents one styled control per trigger`, async ({
  111 |     page,
  112 |   }) => {
> 113 |     const response = await page.goto(`/x/${demo.route}`);
      |                                 ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  114 |     expect(response?.status(), `${demo.route} fixture response`).toBe(200);
  115 | 
  116 |     const triggers = page.locator(`[${demo.behaviorAttribute}]`);
  117 |     await expect(triggers).toHaveCount(demo.count);
  118 |     for (let index = 0; index < demo.count; index++) {
  119 |       const trigger = triggers.nth(index);
  120 |       await expect(trigger).toHaveAttribute("data-gsxui-slot-button", "");
  121 |       await expect(trigger).toHaveAttribute(
  122 |         `data-gsxui-slot-${demo.triggerSlot}`,
  123 |         "",
  124 |       );
  125 |     }
  126 |   });
  127 | }
  128 | 
  129 | test("combobox/basic does not render a disconnected value label", async ({
  130 |   page,
  131 | }) => {
  132 |   const response = await page.goto("/x/combobox/basic");
  133 |   expect(response?.status(), "combobox/basic fixture response").toBe(200);
  134 | 
  135 |   await expect(
  136 |     page.locator('[data-gsxui-slot-combobox-value]'),
  137 |   ).toHaveCount(0);
  138 |   await expect(
  139 |     page.getByText("Choose a framework", { exact: true }),
  140 |   ).toHaveCount(0);
  141 | });
  142 | 
  143 | test("Button variants and sizes keep Nova presentation in normal site scope", async ({
  144 |   page,
  145 | }) => {
  146 |   const response = await page.goto("/x/button/variants");
  147 |   expect(response?.status(), "button/variants fixture response").toBe(200);
  148 |   await expect(page.locator("body")).not.toHaveAttribute(
  149 |     "data-theme-button-preview",
  150 |   );
  151 | 
  152 |   const box = async (name: string) =>
  153 |     page.getByRole("button", { name, exact: true }).evaluate((element) => {
  154 |       const css = getComputedStyle(element);
  155 |       return {
  156 |         width: css.width,
  157 |         height: css.height,
  158 |         backgroundColor: css.backgroundColor,
  159 |         borderTopWidth: css.borderTopWidth,
  160 |       };
  161 |     });
  162 | 
  163 |   const defaultButton = await box("Default");
  164 |   expect(defaultButton.height).toBe("32px");
  165 |   expect(defaultButton.backgroundColor).not.toBe("rgba(0, 0, 0, 0)");
  166 |   expect((await box("Large")).height).toBe("36px");
  167 |   expect(await box("Small icon")).toMatchObject({
  168 |     width: "28px",
  169 |     height: "28px",
  170 |   });
  171 |   expect(await box("Large icon")).toMatchObject({
  172 |     width: "36px",
  173 |     height: "36px",
  174 |   });
  175 |   expect((await box("Outline")).borderTopWidth).toBe("1px");
  176 |   expect((await box("Secondary")).backgroundColor).not.toBe(
  177 |     defaultButton.backgroundColor,
  178 |   );
  179 | });
  180 | 
  181 | test("Button fallback cannot style a theme-preview body", async ({ page }) => {
  182 |   // The fallback stylesheet (web/site-button.css) still owns every element
  183 |   // that carries Button's marker without rendering through ui.Button —
  184 |   // PaginationLink and Calendar's nav/day buttons author no class at all.
  185 |   // ui.Button itself is no longer a witness for this: it ships concrete
  186 |   // utilities compiled from its style recipe, so its inline-flex is intrinsic
  187 |   // and no body-level scope can switch it off. A PaginationLink is the
  188 |   // element whose presentation the fallback still supplies.
  189 |   const response = await page.goto("/x/pagination/basic");
  190 |   expect(response?.status(), "pagination/basic fixture response").toBe(200);
  191 | 
  192 |   const link = page.locator("[data-gsxui-slot-pagination-link]").first();
  193 |   await expect(link).toHaveCSS("display", "inline-flex");
  194 |   await page.locator("body").evaluate((body) => {
  195 |     body.setAttribute("data-theme-button-preview", "");
  196 |   });
  197 |   await expect(link).not.toHaveCSS("display", "inline-flex");
  198 | });
  199 | 
```