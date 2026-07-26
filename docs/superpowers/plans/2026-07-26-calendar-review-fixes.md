# Calendar Review Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Resolve all six final calendar review findings, prove each with a regression, and merge the verified branch into `main`.

**Architecture:** Keep the fixed 42-cell Go/JavaScript rendering split. Normalize the public month before either renderer sees it; centralize exact-year UTC construction in JavaScript; make the roving tab stop focusable through `aria-disabled`; let only multiple-mode hidden form controls grow and shrink; and record the intentional Calendar/Combobox reset overlap explicitly.

**Tech Stack:** gsx and Go `time`, vanilla ES modules through `ui/gsxui.js`, native `FormData`, Playwright Chromium, Go render tests.

## Global Constraints

- Write and run each regression before changing production code.
- JavaScript never creates or destroys one of the 42 calendar day cells.
- Multiple mode uses repeated hidden inputs with the caller's exact `name`.
- Calendar and Combobox may both handle one mixed form's `reset`; no other selector overlap is added.
- Update generated `.x.go` and highlighted example output through the repository generators.
- `make check` must pass end to end before merge.

---

### Task 1: Server defaults, disabled tab stop, and multiple form values

**Files:**
- Modify: `ui/calendar_test.go`
- Modify: `ui/calendar.gsx`
- Generated: `ui/calendar.x.go`

**Interfaces:**
- Consumes: `Calendar(...)`, `dayOnly`, the existing hidden-input bridge.
- Produces: zero-value month resolution and lossless initial multiple-mode inputs.

- [ ] **Step 1: Write failing render regressions**

Add tests that render:

```go
func TestCalendarZeroMonthDefaultsFromSelection(t *testing.T) {
	selected := time.Date(2031, time.March, 18, 12, 0, 0, 0, time.FixedZone("east", 9*60*60))
	got := render(t, ui.Calendar("", time.Time{}, []time.Time{selected},
		time.Time{}, time.Time{}, time.Sunday, true, "label", 0, 0,
		time.Time{}, time.Time{}, nil, nil, "", nil))
	if !strings.Contains(got, `data-gsxui-calendar-month="2031-03"`) {
		t.Fatalf("zero month did not resolve from selection:\n%s", got)
	}
}

func TestCalendarZeroMonthDefaultsFromRangeStart(t *testing.T) {
	from := time.Date(2032, time.April, 9, 0, 0, 0, 0, time.UTC)
	got := render(t, ui.Calendar("range", time.Time{}, nil, from, time.Time{},
		time.Sunday, true, "label", 0, 0, time.Time{}, time.Time{},
		nil, nil, "", nil))
	if !strings.Contains(got, `data-gsxui-calendar-month="2032-04"`) {
		t.Fatalf("zero month did not resolve from range start:\n%s", got)
	}
}
```

Add a third zero-month case with no selection, accepting the UTC month
captured immediately before or after render. Add a disabled-bound case
asserting the sole `tabindex="0"` button has no native `disabled` and has
`aria-disabled="true"`, while another disabled non-tab-stop remains natively
disabled. Add a multiple-mode `name="dates"` case asserting two selected
dates render two hidden inputs with that name and both ISO values; an empty
selection retains one empty placeholder input.

- [ ] **Step 2: Verify RED**

Run:

```bash
go test ./ui -run 'TestCalendar(ZeroMonth|DisabledInitial|MultipleHidden)' -count=1
```

Expected: failures showing year 1, a native-disabled tab stop, and only one multiple value.

- [ ] **Step 3: Implement the minimal server fixes**

Before reading `month.Year()`:

```go
today := time.Now().UTC()
if month.IsZero() {
	switch {
	case mode == "range" && !from.IsZero():
		month = dayOnly(from)
	case len(selected) > 0:
		month = dayOnly(selected[0])
	default:
		month = today
	}
}
```

For each day, compute `isTabStop`. A disabled tab stop emits
`aria-disabled="true"` and omits native `disabled`; hidden days remain
natively disabled. Split multiple-mode hidden rendering from single mode:
one marked input per selected date, or one empty marked placeholder when the
selection is empty.

- [ ] **Step 4: Generate and verify GREEN**

Run:

```bash
go tool gsx generate
go test ./ui -run 'TestCalendar(ZeroMonth|DisabledInitial|MultipleHidden)' -count=1
gopls check -severity=hint ui/calendar.x.go ui/calendar_test.go
```

Expected: focused tests pass; only reviewed modernization hints may remain.

---

### Task 2: Exact years, dynamic multiple inputs, reset month, and mixed reset

**Files:**
- Modify: `jstest/specs/calendar.spec.ts`
- Modify: `site/examples/calendar/multiple.gsx`
- Modify: `site/examples/calendar/form.gsx`
- Modify: `ui/calendar.js`
- Modify: `ui/combobox.js` only if the mixed reset proves its reflection runs before native reset state is available.
- Generated: matching `site/examples/calendar/*.x.go`

**Interfaces:**
- Consumes: the server markers from Task 1, `goTo`, `syncHiddenInputs`, `liveFocused`, `tabStop`.
- Produces: exact-year client arithmetic, dynamic repeated inputs, reset-to-today, and composed reset behavior.

- [ ] **Step 1: Write failing browser regressions**

Add:

```ts
test("years below 100 survive client navigation", async ({ page }) => {
  const server = await page.request.get(`${BASIC}?month=0099-02`);
  const expected = [...(await server.text()).matchAll(
    /data-gsxui-calendar-day[^>]*data-date="([^"]+)"/g,
  )].map((match) => match[1]);
  await page.goto(`${BASIC}?month=0099-01`);
  await page.locator("[data-gsxui-calendar-next]").click();
  await expect(page.locator("[data-gsxui-calendar]")).toHaveAttribute(
    "data-gsxui-calendar-month",
    "0099-02",
  );
  expect(await gridDates(page)).toEqual(expected);
});

test("multiple form data contains every selected date", async ({ page }) => {
  await page.goto(MULTIPLE);
  const root = page.locator("[data-gsxui-calendar]");
  await root.evaluate((el) => {
    const form = document.createElement("form");
    el.parentNode!.insertBefore(form, el);
    form.appendChild(el);
  });
  await dayFor(page, "2026-01-05").click();
  await dayFor(page, "2026-01-08").click();
  expect(
    await root.evaluate((el) =>
      new FormData(el.closest("form")!).getAll("dates").map(String),
    ),
  ).toEqual(["2026-01-05", "2026-01-08"]);
});
```

Add a real Tab test on `LOADED` proving the third Tab enters the day grid,
not `document.body`. Extend the form-reset test: navigate away, reset, and
assert `data-gsxui-calendar-month` equals the client's local `YYYY-MM`.
Compose a Combobox into `calendar/form.gsx`, select a calendar day and a
combobox item, reset once, and assert both return to their server values.

- [ ] **Step 2: Verify RED**

Run the calendar spec against a clean sanctioned harness:

```bash
npx playwright test --config jstest/playwright.config.ts jstest/specs/calendar.spec.ts
```

Expected: the new tests fail for the reviewed reasons; existing tests remain green.

- [ ] **Step 3: Implement exact-year UTC construction**

Add and use everywhere calendar arithmetic constructs a date:

```js
function utcDate(year, month, day) {
  const date = new Date(0);
  date.setUTCFullYear(year, month, day);
  return date;
}
```

Replace `new Date(Date.UTC(...))` in `monthGrid`, `parseISO`, `clientToday`,
day/month/year movement, last-day lookup, and nav-bound clamping.

- [ ] **Step 4: Implement multiple-input synchronization**

Use `data-gsxui-calendar-hidden-multiple` to identify only Calendar's
multiple-mode controls. Keep at least one placeholder so the caller's field
name remains available. Create or remove marked hidden inputs to match
`max(1, selected.length)`, then assign each selected ISO value in order.

- [ ] **Step 5: Implement reset and disabled-tab-stop client behavior**

In `repaint`, degrade a disabled day to `aria-disabled` when it is either
live-focused or the current roving tab stop. During reset, restore the
selection snapshot, clear `liveFocused` and `tabStop`, synchronize inputs,
then call `goTo` with `clientToday()`'s year/month.

If the mixed Calendar/Combobox test demonstrates pre-reset bridge state,
defer each module's reflection to a microtask so it runs after native form
control reset while keeping one delegated listener per module.

- [ ] **Step 6: Generate examples and verify GREEN**

Run:

```bash
go tool gsx generate
make highlight
npx playwright test --config jstest/playwright.config.ts jstest/specs/calendar.spec.ts
```

Expected: all calendar tests pass, including the five new browser regressions.

---

### Task 3: Record the intentional overlap, reconcile docs, and merge

**Files:**
- Modify: `jstest/support/selector-allowlist.ts`
- Modify: `docs/jsx-parity.md`
- Modify: `docs/component-roadmap.md`
- Modify: `docs/superpowers/plans/2026-07-25-calendar.md` only where its final-state claims are now stale.

**Interfaces:**
- Consumes: the mixed form example and verified behavior from Task 2.
- Produces: accurate parity ledger, invariant exception, authoritative green gate.

- [ ] **Step 1: Add the reviewed overlap**

```ts
export const allowedOverlaps: AllowedOverlap[] = [
  {
    modules: ["calendar.js", "combobox.js"],
    key: "reset:false",
    reason:
      "A form may compose Calendar and Combobox. Both reset handlers intentionally " +
      "match that form and restore disjoint descendant state; the mixed calendar/form " +
      "example and calendar behavior spec verify both run correctly.",
  },
];
```

- [ ] **Step 2: Remove stale ledger claims**

Replace the multiple-mode data-loss GAP with the repeated-input mechanism.
Record exact-year UTC construction, reset-to-client-month, and the disabled
tab-stop degradation. Remove roadmap text saying multiple form submission
drops values.

- [ ] **Step 3: Run the authoritative gate**

First ensure no stale test harness owns port 7799. Do not kill unrelated user
processes; identify ownership and use a clean run. Then:

```bash
make check
git diff --check
git status --short
```

Expected: full Go and 154+ Playwright suite green, generated output stable,
no formatting drift, only intended files modified.

- [ ] **Step 4: Commit and merge**

Commit the fix set on `tier4-calendar`. Switch the normal checkout to
`main`, fast-forward or merge the verified feature branch as Git permits,
then run `make check` again on merged `main`. Delete the feature branch only
after the merged gate passes.
