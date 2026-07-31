# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/calendar.spec.ts >> dropdown caption does not crowd the nav buttons
- Location: jstest/specs/calendar.spec.ts:1044:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/calendar/dropdown", waiting until "load"

```

# Test source

```ts
  945  | // controls only, per the Task 6 brief's own ambiguity resolution, rather
  946  | // than reaching for a private hook — is what actually exercises the client
  947  | // path, since basic.gsx's own server render can never contain the real
  948  | // "today" for a fixed, unrelated month.
  949  | test("today is marked from the client's date", async ({ page }) => {
  950  |   await page.goto(BASIC);
  951  | 
  952  |   const { today, delta } = await page.evaluate(() => {
  953  |     const now = new Date();
  954  |     const iso = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, "0")}-${String(now.getDate()).padStart(2, "0")}`;
  955  |     // Months from January 2026 to the client's current month.
  956  |     const months = (now.getFullYear() - 2026) * 12 + now.getMonth();
  957  |     return { today: iso, delta: months };
  958  |   });
  959  | 
  960  |   const button = delta >= 0 ? "[data-gsxui-calendar-next]" : "[data-gsxui-calendar-prev]";
  961  |   for (let i = 0; i < Math.abs(delta); i++) await page.click(button);
  962  | 
  963  |   const marked = await page.$$eval("[data-today]", (els) =>
  964  |     els.map((el) => el.getAttribute("data-date")),
  965  |   );
  966  |   expect(marked).toEqual([today]);
  967  | });
  968  | 
  969  | // The dropdown caption is the one place this port renders a real, visible
  970  | // <select> where upstream renders an invisible one under a styled label. Two
  971  | // things went wrong there and neither is visible to any attribute assertion:
  972  | // ui.NativeSelect's own border painted a second concentric rounded rect inside
  973  | // dropdown_root's, and its hardcoded h-8 overflowed the caption row. A third
  974  | // came from `relative` living on the padded root, so `nav`'s `top-0` resolved
  975  | // against the padding box and pinned the buttons above the row they label.
  976  | //
  977  | // Geometry is the only way to catch any of them, so assert it.
  978  | test("dropdown caption is one row, one border", async ({ page }) => {
  979  |   await page.goto("/x/calendar/dropdown");
  980  | 
  981  |   const geom = await page.evaluate(() => {
  982  |     const mid = (el: Element | null) =>
  983  |       el ? Math.round(el.getBoundingClientRect().top + el.getBoundingClientRect().height / 2) : null;
  984  |     const select = document.querySelector("[data-gsxui-calendar-month-select]")!;
  985  |     const wrapper = select.closest('[data-gsxui-slot-native-select-wrapper]')!;
  986  |     const row = wrapper.parentElement!;
  987  |     return {
  988  |       rowMid: mid(row),
  989  |       wrapperMid: mid(wrapper),
  990  |       prevMid: mid(document.querySelector("[data-gsxui-calendar-prev]")),
  991  |       nextMid: mid(document.querySelector("[data-gsxui-calendar-next]")),
  992  |       selectBorder: getComputedStyle(select).borderTopWidth,
  993  |       wrapperBorder: getComputedStyle(wrapper).borderTopWidth,
  994  |       selectHeight: Math.round(select.getBoundingClientRect().height),
  995  |       rowHeight: Math.round(row.getBoundingClientRect().height),
  996  |     };
  997  |   });
  998  | 
  999  |   // Everything in the caption shares one vertical centre line.
  1000 |   expect(geom.wrapperMid).toBe(geom.rowMid);
  1001 |   expect(geom.prevMid).toBe(geom.rowMid);
  1002 |   expect(geom.nextMid).toBe(geom.rowMid);
  1003 | 
  1004 |   // Exactly one border: dropdown_root's, on the wrapper.
  1005 |   expect(geom.wrapperBorder).toBe("1px");
  1006 |   expect(geom.selectBorder).toBe("0px");
  1007 | 
  1008 |   // The control fits the row rather than overflowing it.
  1009 |   expect(geom.selectHeight).toBeLessThanOrEqual(geom.rowHeight);
  1010 | });
  1011 | 
  1012 | // ui.NativeSelect ships py-1 for its h-8 standalone form-control shape. At the
  1013 | // caption's h-6 that leaves a 16px content box for a 20px line — the line
  1014 | // cannot fit, so no browser can centre it and the text renders visibly low.
  1015 | // The invariant is that the line fits and the padding is symmetric; the
  1016 | // browser centres it from there.
  1017 | test("dropdown caption text can be vertically centred", async ({ page }) => {
  1018 |   await page.goto("/x/calendar/dropdown");
  1019 | 
  1020 |   const m = await page.evaluate(() => {
  1021 |     const sel = document.querySelector("[data-gsxui-calendar-month-select]")!;
  1022 |     const cs = getComputedStyle(sel);
  1023 |     const px = (v: string) => parseFloat(v);
  1024 |     return {
  1025 |       height: sel.getBoundingClientRect().height,
  1026 |       padTop: px(cs.paddingTop),
  1027 |       padBottom: px(cs.paddingBottom),
  1028 |       lineHeight: px(cs.lineHeight),
  1029 |       borderTop: px(cs.borderTopWidth),
  1030 |       borderBottom: px(cs.borderBottomWidth),
  1031 |     };
  1032 |   });
  1033 | 
  1034 |   expect(m.padTop).toBe(m.padBottom);
  1035 |   const contentBox = m.height - m.padTop - m.padBottom - m.borderTop - m.borderBottom;
  1036 |   expect(contentBox).toBeGreaterThanOrEqual(m.lineHeight);
  1037 | });
  1038 | 
  1039 | // The chevron is absolutely positioned inside the select, so its clearance is
  1040 | // reserved with right padding rather than by flex flow the way upstream's
  1041 | // caption label does it. Over-reserve and the controls crowd the nav buttons
  1042 | // until they nearly touch — 2px of gap at pr-6, against the 6px that separates
  1043 | // the two controls from each other. Assert the caption breathes evenly.
  1044 | test("dropdown caption does not crowd the nav buttons", async ({ page }) => {
> 1045 |   await page.goto("/x/calendar/dropdown");
       |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  1046 | 
  1047 |   const gaps = await page.evaluate(() => {
  1048 |     const r = (el: Element) => el.getBoundingClientRect();
  1049 |     const wraps = [
  1050 |       ...document.querySelectorAll('[data-gsxui-slot-native-select-wrapper]'),
  1051 |     ];
  1052 |     const prev = document.querySelector("[data-gsxui-calendar-prev]")!;
  1053 |     const next = document.querySelector("[data-gsxui-calendar-next]")!;
  1054 |     return {
  1055 |       beforeFirst: Math.round(r(wraps[0]).left - r(prev).right),
  1056 |       between: Math.round(r(wraps[1]).left - r(wraps[0]).right),
  1057 |       afterLast: Math.round(r(next).left - r(wraps[wraps.length - 1]).right),
  1058 |     };
  1059 |   });
  1060 | 
  1061 |   // gap-1.5 between the controls is the reference spacing; the outer gaps must
  1062 |   // be at least as generous, and equal to each other.
  1063 |   expect(gaps.beforeFirst).toBeGreaterThanOrEqual(gaps.between);
  1064 |   expect(gaps.afterLast).toBeGreaterThanOrEqual(gaps.between);
  1065 |   expect(gaps.beforeFirst).toBe(gaps.afterLast);
  1066 | });
  1067 | 
```