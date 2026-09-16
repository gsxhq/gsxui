# Opt-in RTL Transform Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A project that sets `"rtl": true` in `gsxui.json` gets its vendored components rewritten from physical to logical direction classes at install time, with `gsxui migrate rtl` for components installed earlier, and the docs site renders its RTL demos from the same transform. Closes #32.

**Architecture:** The registry source stays physical, exactly like upstream shadcn/ui. A new stateless package `internal/rtl` carries upstream's `transform-rtl.ts` mapping table as `Classes(list)`, plus `GSX(filename, src)` (rewrites class attribute literals in gsx source by byte offset, then formats) and `CSS(src)` (rewrites `@apply` lists). The CLI applies it in `selectedStyledArtifact`, the directory-component branch and `init`'s `style.css` when `cfg.RTL` is set, and `gsxui migrate rtl` applies it in place to already-managed files. stylegen emits `site/uirtl/` from `ui/` through the same function so the RTL demos are real.

**Tech Stack:** Go (`internal/rtl`, `internal/srcedit`, `internal/cli`, `internal/stylegen`), gsx parser/ast/gen packages, Tailwind v4 `@utility`, Playwright jstest, structpages site.

**Spec:** `docs/superpowers/specs/2026-09-15-rtl-opt-in-transform-design.md` (read both addenda).

## Global Constraints

- **The registry stays physical.** No task edits `registry/canonical/*.gsx`, `registry/styles/**` or `ui/*.gsx` for direction. Generated files are never hand-edited: after canonical or example edits run `go run ./cmd/stylegen`, `go tool gsx generate`, and `make highlight` when anything under `site/examples/**` or `site/snippets/**` changed.
- **The mapping table is upstream's**, from `packages/shadcn/src/utils/transformers/transform-rtl.ts` at shadcn-ui `ac60ef5c`. Do not add or drop an entry without saying so in the report.
- **Side-keyed classes stay physical in full** (spec addendum): a token whose variant contains `data-[side=left]`, `data-[side=right]`, `data-side=left]` or `data-side=right]`, or that sits in a value arm of `switch side` / `if side …` in a class list, is left unchanged.
- **Idempotent by construction**: `Classes(Classes(x)) == Classes(x)`; companion classes are appended only when absent.
- **Nothing changes for a project without the flag**: every existing byte-identical vendoring assertion in `internal/cli` tests stays as it is.
- Unexported Go identifiers unless serialised or part of the consumer API. `gsx fmt` formatting is enforced by the pre-commit hook: run `go tool gsx fmt -w` on edited `.gsx` files. Attribute interpolation is `name={expr}` with no inner spaces; content interpolation `{ expr }`.
- Docs prose: state the fact once; one fact per bullet; no elaboration.
- Gates after every task: `go build ./...`, `go test ./... -count=1 -short` (full, without `-short`, in Task 9), `go run ./cmd/stylegen --check`, `--check-authoring`, `make audit`, `make verify-generated`, `gofmt -l .`. Playwright is `npx playwright test --config jstest/playwright.config.ts …` (never without the config).
- Commit per task; never `git add -A`. Commit messages end with `Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv`.

---

### Task 1: Extract the byte-offset edit applier into `internal/srcedit`

**Files:**
- Create: `internal/srcedit/srcedit.go`, `internal/srcedit/srcedit_test.go`
- Modify: `internal/stylegen/resolve.go` (delete `literalEdit` and `applyLiteralEdits`, use `srcedit`)

**Interfaces:**
- Produces: `srcedit.Edit{Start, End int; Value string}` and `srcedit.Apply(src []byte, edits []Edit) ([]byte, error)`. Task 3 consumes them.

- [ ] **Step 1: Write the failing test**

`internal/srcedit/srcedit_test.go`:

```go
package srcedit_test

import (
	"testing"

	"github.com/gsxhq/gsxui/internal/srcedit"
)

func TestApplyReplacesSpansFromTheEndSoOffsetsStayValid(t *testing.T) {
	src := []byte("aaa bbb ccc")
	got, err := srcedit.Apply(src, []srcedit.Edit{
		{Start: 0, End: 3, Value: "A"},
		{Start: 8, End: 11, Value: "CCCC"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "A bbb CCCC" {
		t.Fatalf("got %q", got)
	}
	if string(src) != "aaa bbb ccc" {
		t.Fatal("Apply mutated its input")
	}
}

func TestApplyRejectsOverlapAndOutOfRange(t *testing.T) {
	src := []byte("abcdef")
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 0, End: 3, Value: ""}, {Start: 2, End: 4, Value: ""}}); err == nil {
		t.Fatal("overlapping edits accepted")
	}
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 4, End: 9, Value: ""}}); err == nil {
		t.Fatal("out-of-range edit accepted")
	}
	if _, err := srcedit.Apply(src, []srcedit.Edit{{Start: 3, End: 2, Value: ""}}); err == nil {
		t.Fatal("inverted span accepted")
	}
}

func TestApplyWithNoEditsReturnsAnEqualCopy(t *testing.T) {
	src := []byte("unchanged")
	got, err := srcedit.Apply(src, nil)
	if err != nil || string(got) != "unchanged" {
		t.Fatalf("got %q, %v", got, err)
	}
}
```

- [ ] **Step 2: Run it to verify it fails**

Run: `go test ./internal/srcedit -count=1`
Expected: FAIL, package does not exist.

- [ ] **Step 3: Implement**

`internal/srcedit/srcedit.go`:

```go
// Package srcedit applies byte-offset replacements to a source buffer. It
// is the one edit applier shared by the tools that rewrite gsx source in
// place (internal/stylegen's recipe desugaring, internal/rtl's class
// rewriting): parse, record spans, apply from the end so earlier offsets
// stay valid.
package srcedit

import (
	"bytes"
	"fmt"
	"sort"
)

// Edit replaces src[Start:End] with Value.
type Edit struct {
	Start int
	End   int
	Value string
}

// Apply returns a copy of src with every edit applied. Edits may arrive in
// any order; they must not overlap or leave the buffer.
func Apply(src []byte, edits []Edit) ([]byte, error) {
	sorted := append([]Edit(nil), edits...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Start > sorted[j].Start })
	out := append([]byte(nil), src...)
	lastStart := len(src)
	for _, edit := range sorted {
		if edit.Start < 0 || edit.End < edit.Start || edit.End > len(src) {
			return nil, fmt.Errorf("invalid span [%d:%d] in %d bytes", edit.Start, edit.End, len(src))
		}
		if edit.End > lastStart {
			return nil, fmt.Errorf("overlapping span [%d:%d]", edit.Start, edit.End)
		}
		var next bytes.Buffer
		next.Grow(len(out) - (edit.End - edit.Start) + len(edit.Value))
		next.Write(out[:edit.Start])
		next.WriteString(edit.Value)
		next.Write(out[edit.End:])
		out = next.Bytes()
		lastStart = edit.Start
	}
	return out, nil
}
```

In `internal/stylegen/resolve.go`: delete the `literalEdit` type (line ~32) and `applyLiteralEdits` (line ~951); import `"github.com/gsxhq/gsxui/internal/srcedit"`; change the `edits []literalEdit` field and every `literalEdit{start: …, end: …, value: …}` construction to `srcedit.Edit{Start: …, End: …, Value: …}`; replace the `applyLiteralEdits(src, r.edits)` call with `srcedit.Apply(src, r.edits)`. Use the compiler to find every site (`go build ./internal/stylegen`).

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/srcedit ./internal/stylegen -count=1 && go run ./cmd/stylegen --check`
Expected: PASS and no drift (the extraction is behaviour-preserving).

- [ ] **Step 5: Commit**

```bash
git add internal/srcedit internal/stylegen/resolve.go
git commit -m "refactor(stylegen): extract the byte-offset edit applier into internal/srcedit

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 2: `rtl.Classes` — upstream's mapping table

**Files:**
- Create: `internal/rtl/classes.go`, `internal/rtl/classes_test.go`

**Interfaces:**
- Produces: `rtl.Classes(list string) string` and the unexported `classes(list string, sideKeyed bool) string` that Task 3's gsx walker calls with `sideKeyed=true` for `switch side` arms.

- [ ] **Step 1: Write the failing table test**

`internal/rtl/classes_test.go` — every case is upstream's `transform-rtl.test.ts` at the pin, plus the idempotency and side-keyed cases the spec adds:

```go
package rtl

import "testing"

func TestClassesMatchesUpstreamTable(t *testing.T) {
	cases := []struct{ in, want string }{
		// margins, padding
		{"ml-2", "ms-2"}, {"mr-4", "me-4"}, {"-ml-2", "-ms-2"}, {"-mr-4", "-me-4"},
		{"pl-2", "ps-2"}, {"pr-4", "pe-4"},
		// positioning
		{"left-0", "start-0"}, {"right-0", "end-0"}, {"right-1", "end-1"}, {"-left-2", "-start-2"}, {"-right-2", "-end-2"},
		{"inset-l-0", "inset-inline-start-0"}, {"inset-r-0", "inset-inline-end-0"},
		// borders, radii
		{"border-l", "border-s"}, {"border-r", "border-e"}, {"border-l-2", "border-s-2"}, {"border-r-2", "border-e-2"},
		{"rounded-l-md", "rounded-s-md"}, {"rounded-r-md", "rounded-e-md"},
		{"rounded-tl-md", "rounded-ss-md"}, {"rounded-tr-md", "rounded-se-md"},
		{"rounded-bl-md", "rounded-es-md"}, {"rounded-br-md", "rounded-ee-md"},
		// text, scroll, float, clear, origin
		{"text-left", "text-start"}, {"text-right", "text-end"},
		{"scroll-ml-2", "scroll-ms-2"}, {"scroll-mr-2", "scroll-me-2"}, {"scroll-pl-2", "scroll-ps-2"}, {"scroll-pr-2", "scroll-pe-2"},
		{"float-left", "float-start"}, {"float-right", "float-end"}, {"clear-left", "clear-start"}, {"clear-right", "clear-end"},
		{"origin-left", "origin-start"}, {"origin-right", "origin-end"},
		{"origin-top-left", "origin-top-start"}, {"origin-top-right", "origin-top-end"},
		{"origin-bottom-left", "origin-bottom-start"}, {"origin-bottom-right", "origin-bottom-end"},
		// variants, groups, arbitrary values, modifiers, lists
		{"hover:ml-2", "hover:ms-2"}, {"focus:pl-4", "focus:ps-4"}, {"sm:md:ml-2", "sm:md:ms-2"},
		{"sm:group-data-[size=default]/alert-dialog-content:text-left", "sm:group-data-[size=default]/alert-dialog-content:text-start"},
		{"ml-[10px]", "ms-[10px]"}, {"left-[50%]", "start-[50%]"}, {"ml-2/50", "ms-2/50"},
		{"ml-2 mr-4 pl-2 pr-4", "ms-2 me-4 ps-2 pe-4"},
		// slides inside logical side variants
		{"data-[side=inline-start]:slide-in-from-right-2", "data-[side=inline-start]:slide-in-from-end-2"},
		{"data-[side=inline-start]:slide-out-to-right-2", "data-[side=inline-start]:slide-out-to-end-2"},
		{"data-[side=inline-end]:slide-in-from-left-2", "data-[side=inline-end]:slide-in-from-start-2"},
		{"data-[side=inline-end]:slide-out-to-left-2", "data-[side=inline-end]:slide-out-to-start-2"},
		// physical side variants: untouched (upstream cases, plus our full-physical widening)
		{"data-[side=left]:slide-in-from-right-2", "data-[side=left]:slide-in-from-right-2"},
		{"data-[side=right]:slide-in-from-left-2", "data-[side=right]:slide-in-from-left-2"},
		{"data-[side=left]:-right-1", "data-[side=left]:-right-1"},
		{"data-[side=right]:-left-1", "data-[side=right]:-left-1"},
		{"data-[side=right]:left-0", "data-[side=right]:left-0"},
		{"data-[side=left]:border-r", "data-[side=left]:border-r"},
		{"[[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=left]>&]:border-r", "[[data-gsxui-slot-sidebar-desktop][data-variant=sidebar][data-side=left]>&]:border-r"},
		// unrelated and near-miss classes
		{"bg-red-500", "bg-red-500"}, {"flex", "flex"}, {"mx-auto", "mx-auto"}, {"px-4", "px-4"},
		{"border-ring", "border-ring"}, {"border-ring/50", "border-ring/50"}, {"border-lime-500", "border-lime-500"},
		{"scroll-m-4", "scroll-m-4"}, {"transition-[left,right,width]", "transition-[left,right,width]"},
		// translate-x companions
		{"-translate-x-1/2", "-translate-x-1/2 rtl:translate-x-1/2"},
		{"translate-x-full", "translate-x-full rtl:-translate-x-full"},
		{"-translate-x-px", "-translate-x-px rtl:translate-x-px"},
		{"after:-translate-x-1/2", "after:-translate-x-1/2 rtl:after:translate-x-1/2"},
		{"group-hover:translate-x-2", "group-hover:translate-x-2 rtl:group-hover:-translate-x-2"},
		{"-translate-y-1/2", "-translate-y-1/2"}, {"translate-y-full", "translate-y-full"},
		// reverse companions
		{"space-x-4", "space-x-4 rtl:space-x-reverse"}, {"space-x-0", "space-x-0 rtl:space-x-reverse"},
		{"divide-x-2", "divide-x-2 rtl:divide-x-reverse"},
		{"md:space-x-4", "md:space-x-4 rtl:md:space-x-reverse"},
		{"hover:divide-x-2", "hover:divide-x-2 rtl:hover:divide-x-reverse"},
		{"space-y-4", "space-y-4"}, {"divide-y-2", "divide-y-2"},
		// cursor swaps
		{"cursor-w-resize", "cursor-w-resize rtl:cursor-e-resize"},
		{"cursor-e-resize", "cursor-e-resize rtl:cursor-w-resize"},
		{"hover:cursor-w-resize", "hover:cursor-w-resize rtl:hover:cursor-e-resize"},
		// already-directional prefixes pass through
		{"rtl:rotate-180 size-4", "rtl:rotate-180 size-4"}, {"ltr:ml-2", "ltr:ml-2"},
		// centering composes to something that centres under both directions
		{"left-1/2 -translate-x-1/2", "start-1/2 -translate-x-1/2 rtl:translate-x-1/2"},
	}
	for _, c := range cases {
		if got := Classes(c.in); got != c.want {
			t.Errorf("Classes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestClassesIsIdempotent(t *testing.T) {
	for _, in := range []string{
		"ml-2 mr-4 pl-2 pr-4", "-translate-x-1/2", "space-x-4", "cursor-w-resize", "left-1/2 -translate-x-1/2",
		"translate-x-2 rtl:-translate-x-2", "md:space-x-4 rtl:md:space-x-reverse", "text-left hover:ml-2",
	} {
		once := Classes(in)
		if twice := Classes(once); twice != once {
			t.Errorf("Classes not a fixed point on %q: once=%q twice=%q", in, once, twice)
		}
	}
}

func TestClassesSideKeyedArmStaysPhysical(t *testing.T) {
	arm := "inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm right-auto data-[state=closed]:slide-out-to-left data-[state=open]:slide-in-from-left"
	if got := classes(arm, true); got != arm {
		t.Errorf("side-keyed arm changed:\n got %q\nwant %q", got, arm)
	}
	if got := classes("text-left ml-2", true); got != "text-left ml-2" {
		t.Errorf("side-keyed arm mapped non-positioning classes: %q", got)
	}
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./internal/rtl -count=1`
Expected: FAIL, `Classes` undefined.

- [ ] **Step 3: Implement**

`internal/rtl/classes.go`:

```go
// Package rtl rewrites physical direction utilities to logical ones. The
// table and the rule order are shadcn/ui's transform-rtl.ts at the pinned
// upstream commit (see docs/jsx-parity.md ## rtl); gsxui applies it to
// vendored source when a project sets "rtl": true in gsxui.json, and to
// site/uirtl for the docs site's RTL demos.
package rtl

import "strings"

// direct maps a physical prefix (or whole class, for entries without a
// trailing dash) to its logical spelling. Order matters: longer prefixes
// first where one is a prefix of another.
var direct = [][2]string{
	{"-ml-", "-ms-"}, {"-mr-", "-me-"}, {"ml-", "ms-"}, {"mr-", "me-"},
	{"pl-", "ps-"}, {"pr-", "pe-"},
	{"-left-", "-start-"}, {"-right-", "-end-"}, {"left-", "start-"}, {"right-", "end-"},
	{"inset-l-", "inset-inline-start-"}, {"inset-r-", "inset-inline-end-"},
	{"rounded-tl-", "rounded-ss-"}, {"rounded-tr-", "rounded-se-"},
	{"rounded-bl-", "rounded-es-"}, {"rounded-br-", "rounded-ee-"},
	{"rounded-l-", "rounded-s-"}, {"rounded-r-", "rounded-e-"},
	{"border-l-", "border-s-"}, {"border-r-", "border-e-"},
	{"border-l", "border-s"}, {"border-r", "border-e"},
	{"text-left", "text-start"}, {"text-right", "text-end"},
	{"scroll-ml-", "scroll-ms-"}, {"scroll-mr-", "scroll-me-"},
	{"scroll-pl-", "scroll-ps-"}, {"scroll-pr-", "scroll-pe-"},
	{"float-left", "float-start"}, {"float-right", "float-end"},
	{"clear-left", "clear-start"}, {"clear-right", "clear-end"},
	{"origin-top-left", "origin-top-start"}, {"origin-top-right", "origin-top-end"},
	{"origin-bottom-left", "origin-bottom-start"}, {"origin-bottom-right", "origin-bottom-end"},
	{"origin-left", "origin-start"}, {"origin-right", "origin-end"},
}

// translateX keeps the class and adds the sign-flipped form under rtl:.
var translateX = [][2]string{{"-translate-x-", "translate-x-"}, {"translate-x-", "-translate-x-"}}

// reverse keeps the class and adds the *-reverse form under rtl:.
var reverse = [][2]string{{"space-x-", "space-x-reverse"}, {"divide-x-", "divide-x-reverse"}}

// swap keeps the class and adds the swapped value under rtl:.
var swap = [][2]string{{"cursor-w-resize", "cursor-e-resize"}, {"cursor-e-resize", "cursor-w-resize"}}

// logicalSideSlides: [variant fragment, physical prefix, logical prefix].
var logicalSideSlides = [][3]string{
	{"data-[side=inline-start]", "slide-in-from-right", "slide-in-from-end"},
	{"data-[side=inline-start]", "slide-out-to-right", "slide-out-to-end"},
	{"data-[side=inline-end]", "slide-in-from-left", "slide-in-from-start"},
	{"data-[side=inline-end]", "slide-out-to-left", "slide-out-to-start"},
}

// sideKeyedVariants mark a token as placed on a physical side. Upstream
// spells it data-[side=…]; gsxui's sidebar rail spells it inside an
// arbitrary variant as data-side=…]. Both stay physical in full.
var sideKeyedVariants = []string{"data-[side=left]", "data-[side=right]", "data-side=left]", "data-side=right]"}

// Classes rewrites one whitespace-separated class list. Companion classes
// are appended only when absent, so the function is a fixed point after
// one application.
func Classes(list string) string { return classes(list, false) }

func classes(list string, sideKeyed bool) string {
	fields := strings.Fields(list)
	present := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		present[f] = struct{}{}
	}
	out := make([]string, 0, len(fields))
	emit := func(c string) { out = append(out, c); present[c] = struct{}{} }
	companion := func(c string) {
		if _, ok := present[c]; !ok {
			emit(c)
		}
	}
	for _, class := range fields {
		if strings.HasPrefix(class, "rtl:") || strings.HasPrefix(class, "ltr:") {
			emit(class)
			continue
		}
		variant, value, modifier := split(class)
		if value == "" || sideKeyed || isSideKeyed(variant) {
			emit(class)
			continue
		}
		join := func(v string) string {
			if variant != "" {
				v = variant + ":" + v
			}
			return v
		}
		withModifier := func(v string) string {
			if modifier != "" {
				return v + "/" + modifier
			}
			return v
		}
		if p, ok := prefixMatch(value, translateX); ok {
			emit(class)
			companion("rtl:" + join(withModifier(p[1]+value[len(p[0]):])))
			continue
		}
		if p, ok := prefixMatch(value, reverse); ok {
			emit(class)
			companion("rtl:" + join(p[1]))
			continue
		}
		if p, ok := exactMatch(value, swap); ok {
			emit(class)
			companion("rtl:" + join(p[1]))
			continue
		}
		mapped := value
		slid := false
		for _, s := range logicalSideSlides {
			if strings.Contains(variant, s[0]) && strings.HasPrefix(value, s[1]) {
				mapped = s[2] + value[len(s[1]):]
				slid = true
				break
			}
		}
		if !slid {
			for _, p := range direct {
				if !strings.HasPrefix(value, p[0]) {
					continue
				}
				if !strings.HasSuffix(p[0], "-") && value != p[0] {
					continue // border-r must not match border-ring
				}
				mapped = p[1] + value[len(p[0]):]
				break
			}
		}
		emit(join(withModifier(mapped)))
	}
	return strings.Join(out, " ")
}

func isSideKeyed(variant string) bool {
	for _, s := range sideKeyedVariants {
		if strings.Contains(variant, s) {
			return true
		}
	}
	return false
}

func prefixMatch(value string, table [][2]string) ([2]string, bool) {
	for _, p := range table {
		if strings.HasPrefix(value, p[0]) {
			return p, true
		}
	}
	return [2]string{}, false
}

func exactMatch(value string, table [][2]string) ([2]string, bool) {
	for _, p := range table {
		if value == p[0] {
			return p, true
		}
	}
	return [2]string{}, false
}

// split is upstream's splitClassName: the variant ends at the last colon
// outside brackets; the modifier follows the last slash of what remains.
func split(class string) (variant, value, modifier string) {
	if !strings.ContainsAny(class, ":/") {
		return "", class, ""
	}
	depth := 0
	colon := -1
	for i := len(class) - 1; i >= 0; i-- {
		switch class[i] {
		case ']':
			depth++
		case '[':
			depth--
		case ':':
			if depth == 0 {
				colon = i
			}
		}
		if colon >= 0 {
			break
		}
	}
	rest := class
	if colon >= 0 {
		variant, rest = class[:colon], class[colon+1:]
	}
	if slash := strings.LastIndex(rest, "/"); slash >= 0 {
		return variant, rest[:slash], rest[slash+1:]
	}
	return variant, rest, ""
}
```

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/rtl -count=1 -v 2>&1 | tail -5`
Expected: PASS. If a single upstream case fails, the table or rule order is wrong; fix the implementation, never the expectation.

- [ ] **Step 5: Commit**

```bash
git add internal/rtl/classes.go internal/rtl/classes_test.go
git commit -m "feat(rtl): Classes, upstream's physical-to-logical class table

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 3: `rtl.GSX` and `rtl.CSS` — rewriting source in place

**Files:**
- Create: `internal/rtl/gsx.go`, `internal/rtl/css.go`, `internal/rtl/gsx_test.go`, `internal/rtl/testdata/composed.gsx`, `internal/rtl/testdata/composed.want.gsx`

**Interfaces:**
- Consumes: `srcedit.Apply`, `classes(list, sideKeyed)`.
- Produces: `rtl.GSX(filename string, src []byte) ([]byte, error)` and `rtl.CSS(src []byte) []byte`. Tasks 5, 6, 7 consume them.

- [ ] **Step 1: Write the fixtures and failing tests**

`internal/rtl/testdata/composed.gsx`:

```gsx
package ui

import "github.com/gsxhq/gsx"

component Fixture(side string, active string, children gsx.Node, attrs gsx.Attrs) {
	<div class="ml-2 text-left" data-a>
		<span
			class={
				"pl-4 -translate-x-1/2",
				"border-l": active == "left",
				switch side {
				case "left":
					"left-0 border-r slide-in-from-left"
				default:
					"right-0 border-l"
				}
			}
			title="pl-4 stays: not a class"
		>
			{ children }
		</span>
		<i class={ "data-[side=left]:-right-1 hover:mr-2", attrs.Class() }></i>
		<b class={ `mr-1` }></b>
	</div>
}
```

`internal/rtl/testdata/composed.want.gsx`: the same file after transformation, which is what `gsx fmt` produces from these edits:

```gsx
package ui

import "github.com/gsxhq/gsx"

component Fixture(side string, active string, children gsx.Node, attrs gsx.Attrs) {
	<div class="ms-2 text-start" data-a>
		<span
			class={
				"ps-4 -translate-x-1/2 rtl:translate-x-1/2",
				"border-s": active == "left",
				switch side {
				case "left":
					"left-0 border-r slide-in-from-left"
				default:
					"right-0 border-l"
				}
			}
			title="pl-4 stays: not a class"
		>
			{ children }
		</span>
		<i class={ "data-[side=left]:-right-1 hover:me-2", attrs.Class() }></i>
		<b class={ `me-1` }></b>
	</div>
}
```

Before relying on the `.want` file, run `go tool gsx fmt -l internal/rtl/testdata/composed.want.gsx`; if fmt wants a different layout, take fmt's output as the fixture (the transform ends with `gen.Format`, so the expectation must be fmt-stable).

`internal/rtl/gsx_test.go`:

```go
package rtl_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/gsxhq/gsxui/internal/rtl"
)

func TestGSXRewritesOnlyClassLiterals(t *testing.T) {
	src, err := os.ReadFile("testdata/composed.gsx")
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/composed.want.gsx")
	if err != nil {
		t.Fatal(err)
	}
	got, err := rtl.GSX("composed.gsx", src)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("transform mismatch\n--- got ---\n%s\n--- want ---\n%s", got, want)
	}
	again, err := rtl.GSX("composed.gsx", got)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(again, got) {
		t.Fatal("GSX is not idempotent")
	}
}

func TestGSXLeavesAFileWithoutDirectionClassesByteIdentical(t *testing.T) {
	src := []byte("package ui\n\ncomponent Plain(attrs gsx.Attrs) {\n\t<p class={ \"flex gap-2\", attrs.Class() }>x</p>\n}\n")
	got, err := rtl.GSX("plain.gsx", src)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, src) {
		t.Fatalf("unchanged file was rewritten:\n%s", got)
	}
}

func TestCSSRewritesApplyListsOnly(t *testing.T) {
	src := []byte(".a {\n  @apply pl-1.5 text-left;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n")
	want := ".a {\n  @apply ps-1.5 text-start;\n}\n/* pl-1.5 in a comment stays */\n.b { color: left; }\n"
	if got := string(rtl.CSS(src)); got != want {
		t.Fatalf("got\n%s\nwant\n%s", got, want)
	}
	if got := string(rtl.CSS([]byte(want))); got != want {
		t.Fatal("CSS is not idempotent")
	}
}

// Every shipped component transforms, reparses, is a fixed point, and
// carries no physical direction class outside side-keyed arms and variants.
func TestGSXAcrossTheRegistry(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "ui", "*.gsx"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no ui/*.gsx found: %v", err)
	}
	physical := regexp.MustCompile(`^-?(m|p|border|rounded|scroll-m|scroll-p)?[lr]-|^-?(left|right)-|^text-(left|right)$|^border-[lr]$|^float-(left|right)$|^clear-(left|right)$|^origin-(top-|bottom-)?(left|right)$`)
	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		out, err := rtl.GSX(filepath.Base(file), src)
		if err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		again, err := rtl.GSX(filepath.Base(file), out)
		if err != nil || !bytes.Equal(again, out) {
			t.Fatalf("%s: not idempotent (%v)", file, err)
		}
		for _, list := range rtl.ClassListsForTest(filepath.Base(file), out) {
			if list.SideKeyed {
				continue
			}
			for _, class := range strings.Fields(list.Value) {
				variant, value := rtl.SplitForTest(class)
				if rtl.IsSideKeyedForTest(variant) {
					continue
				}
				if physical.MatchString(value) {
					t.Errorf("%s: physical class %q survived in %q", filepath.Base(file), class, list.Value)
				}
			}
		}
	}
}
```

Add `internal/rtl/export_test.go` exposing `ClassListsForTest`, `SplitForTest` and `IsSideKeyedForTest` over the unexported walker (see Step 3) — a test-only file in `package rtl`.

- [ ] **Step 2: Run to verify it fails**

Run: `go test ./internal/rtl -count=1`
Expected: FAIL, `GSX`/`CSS` undefined.

- [ ] **Step 3: Implement the gsx walker**

`internal/rtl/gsx.go`:

```go
package rtl

import (
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"strconv"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/srcedit"
)

// classList is one class-carrying string literal found in a gsx file:
// its unquoted value, its byte span in src, whether it was a raw (backtick)
// literal, and whether it sits in a side-keyed arm.
type classList struct {
	Value     string
	Start     int
	End       int
	Raw       bool
	SideKeyed bool
}

// GSX rewrites every class attribute literal in one gsx source file and
// returns the gsx-formatted result. A file with nothing to rewrite comes
// back byte-identical. Only string literals that are class-list elements
// or pair keys change; a literal that is a condition operand, another
// attribute, or text is untouched.
func GSX(filename string, src []byte) ([]byte, error) {
	lists, err := classLists(filename, src)
	if err != nil {
		return nil, err
	}
	var edits []srcedit.Edit
	for _, l := range lists {
		mapped := classes(l.Value, l.SideKeyed)
		if mapped == l.Value {
			continue
		}
		quoted := strconv.Quote(mapped)
		if l.Raw {
			quoted = "`" + mapped + "`"
		}
		edits = append(edits, srcedit.Edit{Start: l.Start, End: l.End, Value: quoted})
	}
	if len(edits) == 0 {
		return src, nil
	}
	edited, err := srcedit.Apply(src, edits)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	formatted, err := gen.Format(filename, edited)
	if err != nil {
		return nil, fmt.Errorf("%s: format after rtl rewrite: %w", filename, err)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, formatted, 0); err != nil {
		return nil, fmt.Errorf("%s: reparse after rtl rewrite: %w", filename, err)
	}
	return formatted, nil
}

// classLists finds every class literal: class="…" static attributes, and
// inside class={ … } every element literal, pair key, and value-form
// if/switch arm literal. Arms of `switch side` / `if side …` are side-keyed.
func classLists(filename string, src []byte) ([]classList, error) {
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}
	var lists []classList
	var walkErr error
	gsxast.Inspect(file, func(n gsxast.Node) bool {
		if walkErr != nil {
			return false
		}
		switch a := n.(type) {
		case *gsxast.StaticAttr:
			if a.Name != "class" {
				return true
			}
			end := fset.Position(a.End()).Offset - 1 // closing quote
			start := end - len(a.Value)
			if start < 0 || string(src[start:end]) != a.Value {
				walkErr = fmt.Errorf("%s: static class span mismatch at %v", filename, fset.Position(a.Pos()))
				return false
			}
			lists = append(lists, classList{Value: a.Value, Start: start, End: end})
		case *gsxast.ComposedAttr:
			if a.Name != "class" {
				return true
			}
			for i := range a.Parts {
				part := &a.Parts[i]
				if part.Expr != "" {
					if l, ok, err := literalIn(fset, src, part.Expr, part.ExprPos, false); err != nil {
						walkErr = err
						return false
					} else if ok {
						lists = append(lists, l)
					}
				}
				if part.CF == nil {
					continue
				}
				if part.CF.If != nil {
					sideKeyed := mentionsSide(part.CF.If.Cond)
					for vi := part.CF.If; vi != nil; vi = vi.ElseIf {
						sideKeyed = sideKeyed || mentionsSide(vi.Cond)
						if vi.Then != nil {
							if l, ok, err := armLiteral(fset, src, vi.Then, sideKeyed); err != nil {
								walkErr = err
								return false
							} else if ok {
								lists = append(lists, l)
							}
						}
						if vi.ElseIf == nil && vi.Else != nil {
							if l, ok, err := armLiteral(fset, src, vi.Else, sideKeyed); err != nil {
								walkErr = err
								return false
							} else if ok {
								lists = append(lists, l)
							}
						}
					}
				}
				if part.CF.Switch != nil {
					sideKeyed := mentionsSide(part.CF.Switch.Tag)
					for _, c := range part.CF.Switch.Cases {
						if c.Value == nil {
							continue
						}
						if l, ok, err := armLiteral(fset, src, c.Value, sideKeyed); err != nil {
							walkErr = err
							return false
						} else if ok {
							lists = append(lists, l)
						}
					}
				}
			}
		}
		return true
	})
	return lists, walkErr
}

func armLiteral(fset *token.FileSet, src []byte, arm *gsxast.ValueArm, sideKeyed bool) (classList, bool, error) {
	if arm.Segments != nil || arm.Expr == "" {
		return classList{}, false, nil // f`…` literal arms carry no class table entries in shipped source
	}
	return literalIn(fset, src, arm.Expr, arm.ExprPos, sideKeyed)
}

// literalIn parses one Go expression from a class list and, when it is a
// plain string literal (parentheses allowed), returns its span in src.
// go/parser.ParseExpr positions start at 1, so offset = exprOffset + pos - 1.
func literalIn(fset *token.FileSet, src []byte, expr string, exprPos token.Pos, sideKeyed bool) (classList, bool, error) {
	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		return classList{}, false, nil // not a literal (a helper call, a variable): nothing to rewrite
	}
	for {
		p, ok := parsed.(*goast.ParenExpr)
		if !ok {
			break
		}
		parsed = p.X
	}
	lit, ok := parsed.(*goast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return classList{}, false, nil
	}
	base := fset.Position(exprPos).Offset
	start := base + int(lit.Pos()) - 1
	end := base + int(lit.End()) - 1
	if start < 0 || end > len(src) || string(src[start:end]) != lit.Value {
		return classList{}, false, fmt.Errorf("class literal span mismatch at offset %d", base)
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return classList{}, false, fmt.Errorf("class literal %s: %w", lit.Value, err)
	}
	return classList{Value: value, Start: start, End: end, Raw: lit.Value[0] == '`', SideKeyed: sideKeyed}, true, nil
}

// mentionsSide reports whether a switch tag or if condition is keyed on the
// component's side parameter. Sheet, Drawer and Sidebar all name it `side`.
func mentionsSide(goExpr string) bool {
	parsed, err := goparser.ParseExpr(goExpr)
	if err != nil {
		return false
	}
	found := false
	goast.Inspect(parsed, func(n goast.Node) bool {
		if id, ok := n.(*goast.Ident); ok && id.Name == "side" {
			found = true
		}
		return !found
	})
	return found
}
```

`internal/rtl/css.go`:

```go
package rtl

import "regexp"

var applyList = regexp.MustCompile(`(@apply\s+)([^;{}]+?)(\s*;)`)

// CSS rewrites the utility list of every @apply. Selectors, properties and
// comments are untouched; the pass is a fixed point like Classes.
func CSS(src []byte) []byte {
	return applyList.ReplaceAllFunc(src, func(m []byte) []byte {
		parts := applyList.FindSubmatch(m)
		return append(append(append([]byte{}, parts[1]...), Classes(string(parts[2]))...), parts[3]...)
	})
}
```

`internal/rtl/export_test.go`:

```go
package rtl

type ClassListForTest = classList

func ClassListsForTest(filename string, src []byte) []classList {
	lists, err := classLists(filename, src)
	if err != nil {
		panic(err)
	}
	return lists
}

func SplitForTest(class string) (variant, value string) {
	variant, value, _ = split(class)
	return variant, value
}

func IsSideKeyedForTest(variant string) bool { return isSideKeyed(variant) }
```

Note `TestGSXAcrossTheRegistry` is in `package rtl_test` and uses these exported test hooks.

The static-attribute branch assumes `StaticAttr`'s span ends at the closing quote and asserts it against `src`. If that assertion trips on the first real file, the span excludes the quotes: derive `start` instead by scanning forward from `fset.Position(a.Pos()).Offset` to the first `"` after `=` and taking `start+1`, keep the `src[start:end] == a.Value` assertion, and say so in the report.

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/rtl -count=1 -v 2>&1 | tail -20`
Expected: PASS. If `TestGSXAcrossTheRegistry` reports a survivor, read it: a class inside a `switch side` arm or a `data-side` variant is expected to survive and the test skips those, so any report is either a missing table entry (add it only if upstream has it, and say so) or a component keying on side by another name (report it; do not widen `mentionsSide` without saying so).

- [ ] **Step 5: Commit**

```bash
git add internal/rtl
git commit -m "feat(rtl): GSX and CSS rewrite class literals and @apply lists in place

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 4: The `rtl` flag in `gsxui.json` and `gsxui init --rtl`

**Files:**
- Modify: `internal/cli/config.go` (`Config`), `internal/cli/init.go` (flag, usage, propagation), `internal/cli/config_test.go`, `internal/cli/init_test.go`

**Interfaces:**
- Produces: `Config.RTL bool` (json `rtl,omitempty`); `gsxui init --rtl`.

- [ ] **Step 1: Write the failing tests**

Append to `internal/cli/config_test.go`:

```go
func TestConfigRTLRoundTripsAndDefaultsFalse(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig()
	if cfg.RTL {
		t.Fatal("RTL must default to false")
	}
	cfg.RTL = true
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	raw, _ := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if !strings.Contains(string(raw), `"rtl": true`) {
		t.Fatalf("gsxui.json missing rtl: %s", raw)
	}
	loaded, err := LoadConfig(dir)
	if err != nil || !loaded.RTL {
		t.Fatalf("RTL not round-tripped: %+v %v", loaded, err)
	}
	cfg.RTL = false
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	raw, _ = os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if strings.Contains(string(raw), "rtl") {
		t.Fatalf("false rtl must be omitted: %s", raw)
	}
}
```

Append to `internal/cli/init_test.go`:

```go
func TestInitRTLFlagPersistsAndSurvivesRerun(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init", "--rtl"}); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadConfig(dir)
	if err != nil || !cfg.RTL {
		t.Fatalf("init --rtl did not persist rtl: %+v %v", cfg, err)
	}
	// style.css is vendored logical: the default sheet carries pl-1.5/pr-1.5.
	style := readFile(t, dir, "web/gsxui/style.css")
	if strings.Contains(style, "@apply pl-1.5") || !strings.Contains(style, "ps-1.5") {
		t.Fatalf("style.css was not transformed:\n%s", style)
	}
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	cfg, err = LoadConfig(dir)
	if err != nil || !cfg.RTL {
		t.Fatalf("a rerun without --rtl reset the flag: %+v %v", cfg, err)
	}
}

func TestInitWithoutRTLLeavesStyleCSSPhysical(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(readFile(t, dir, "web/gsxui/style.css"), "@apply pl-1.5") {
		t.Fatal("style.css changed for a project without the flag")
	}
}
```

(`readFile(t, dir, rel)` already exists in the cli test helpers; if its name differs, use the existing helper.)

- [ ] **Step 2: Run to verify they fail**

Run: `go test ./internal/cli -run 'TestConfigRTL|TestInitRTL|TestInitWithoutRTL' -count=1`
Expected: FAIL, `cfg.RTL` undefined / flag unknown.

- [ ] **Step 3: Implement**

`internal/cli/config.go`: add to `Config`:

```go
	// RTL vendors components with logical direction classes (see
	// internal/rtl). Off by default; set by `gsxui init --rtl`.
	RTL bool `json:"rtl,omitempty"`
```

`internal/cli/init.go` in `runInit`: after the `overwrite` flag add
`rtlFlag := flags.Bool("rtl", false, "vendor components with logical (direction-aware) classes; see /docs/rtl")`, update the usage string to `usage: gsxui init [--preset <file|code|->] [--overwrite] [--rtl]`, and after `cfg` is loaded or defaulted: `if *rtlFlag { cfg.RTL = true }`. In `initArtifacts`, where the `defaultStyleCSS` case composes `content`, add after composition:

```go
			if cfg.RTL {
				content = rtl.CSS(content)
			}
```

Import `"github.com/gsxhq/gsxui/internal/rtl"`.

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/cli -count=1 -short`
Expected: PASS, including every pre-existing init test unchanged.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/config.go internal/cli/init.go internal/cli/config_test.go internal/cli/init_test.go
git commit -m "feat(cli): gsxui.json rtl flag and init --rtl

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 5: `add` and `apply` vendor through the transform when the flag is set

**Files:**
- Modify: `internal/cli/apply.go` (`selectedStyledArtifact`), `internal/cli/add.go` (directory-component `.gsx` branch)
- Test: `internal/cli/add_test.go`, `internal/cli/apply_test.go`

- [ ] **Step 1: Write the failing tests**

Append to `internal/cli/add_test.go`:

```go
func TestAddWithRTLVendorsLogicalClasses(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init", "--rtl"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "native-select", "table"}); err != nil {
		t.Fatal(err)
	}
	ns := readFile(t, dir, "ui/native-select.gsx")
	for _, want := range []string{"end-2.5", "pe-8"} {
		if !strings.Contains(ns, want) {
			t.Errorf("native-select missing %q", want)
		}
	}
	for _, bad := range []string{"right-2.5", "pr-8"} {
		if strings.Contains(ns, bad) {
			t.Errorf("native-select still carries %q", bad)
		}
	}
	if tbl := readFile(t, dir, "ui/table.gsx"); !strings.Contains(tbl, "text-start") || strings.Contains(tbl, "text-left") {
		t.Errorf("table not transformed")
	}
	// The recorded hash is of the transformed content: a second add is a no-op.
	before := readFile(t, dir, "ui/native-select.gsx")
	if err := Run([]string{"add", "native-select"}); err != nil {
		t.Fatal(err)
	}
	if readFile(t, dir, "ui/native-select.gsx") != before {
		t.Fatal("re-adding under rtl changed the vendored file")
	}
}
```

Append to `internal/cli/apply_test.go` (follow the file's existing style-switch test for setup; the assertion is the point):

```go
func TestApplyStyleSwitchKeepsRTL(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init", "--rtl"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "native-select"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"apply", "--preset", presetCode(t, preset.Default(preset.StyleMaia)), "--yes"}); err != nil {
		t.Fatal(err)
	}
	ns := readFile(t, dir, "ui/native-select.gsx")
	if strings.Contains(ns, "right-2.5") || !strings.Contains(ns, "end-2.5") {
		t.Fatalf("apply re-vendored native-select physical:\n%s", ns)
	}
}
```

- [ ] **Step 2: Run to verify they fail**

Run: `go test ./internal/cli -run 'TestAddWithRTL|TestApplyStyleSwitchKeepsRTL' -count=1`
Expected: FAIL on the `end-2.5` assertions.

- [ ] **Step 3: Implement**

In `internal/cli/apply.go` `selectedStyledArtifact`, after `RewriteGsx` succeeds:

```go
	if cfg.RTL {
		rewritten, err = rtl.GSX(sourcePath, rewritten)
		if err != nil {
			return artifact{}, fmt.Errorf("rtl transform %s: %w", sourcePath, err)
		}
	}
```

In `internal/cli/add.go` `addArtifacts`, in the directory-component branch after the `RewriteGsx` call for `.gsx`/`.go` files, add for `.gsx` only:

```go
				if cfg.RTL && strings.HasSuffix(fname, ".gsx") {
					src, err = rtl.GSX(fname, src)
					if err != nil {
						return nil, fmt.Errorf("rtl transform %s/%s: %w", name, fname, err)
					}
				}
```

Import `"github.com/gsxhq/gsxui/internal/rtl"` in both files. Nothing else changes: the artifact hash is taken from `Content` after this point, so managed-file detection sees the transformed bytes.

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/cli -count=1 -short`
Expected: PASS, including `TestAddVendorsWithDeps`, `TestAddButtonUsesSelectedStyleExactSource` and every byte-identical assertion for projects without the flag.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/apply.go internal/cli/add.go internal/cli/add_test.go internal/cli/apply_test.go
git commit -m "feat(cli): add and apply vendor through the rtl transform when the flag is set

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 6: `gsxui migrate rtl`

**Files:**
- Create: `internal/cli/migrate.go`, `internal/cli/migrate_test.go`
- Modify: `internal/cli/run.go` (dispatch and usage)

**Interfaces:**
- Produces: `gsxui migrate rtl`; `gsxui migrate` with no argument lists migrations.

- [ ] **Step 1: Write the failing tests**

`internal/cli/migrate_test.go`:

```go
package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMigrateRTLTransformsManagedFilesInPlaceAndIsIdempotent(t *testing.T) {
	dir, _ := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "native-select", "table"}); err != nil {
		t.Fatal(err)
	}
	// The consumer edited a vendored file; migrate must keep that edit.
	tablePath := filepath.Join(dir, "ui", "table.gsx")
	table, _ := os.ReadFile(tablePath)
	edited := strings.Replace(string(table), "package ui\n", "package ui\n\n// consumer note\n", 1)
	if err := os.WriteFile(tablePath, []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Run([]string{"migrate", "rtl"}); err == nil || !strings.Contains(err.Error(), `"rtl": true`) {
		t.Fatalf("migrate rtl without the flag must refuse and say how to set it, got %v", err)
	}
	cfg, _ := LoadConfig(dir)
	cfg.RTL = true
	if err := cfg.Save(dir); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"migrate", "rtl"}); err != nil {
		t.Fatal(err)
	}
	ns := readFile(t, dir, "ui/native-select.gsx")
	if strings.Contains(ns, "right-2.5") || !strings.Contains(ns, "end-2.5") {
		t.Fatalf("native-select not migrated:\n%s", ns)
	}
	got := readFile(t, dir, "ui/table.gsx")
	if !strings.Contains(got, "// consumer note") {
		t.Fatal("migrate discarded the consumer's edit")
	}
	if strings.Contains(got, "text-left") || !strings.Contains(got, "text-start") {
		t.Fatal("table not migrated")
	}
	style := readFile(t, dir, "web/gsxui/style.css")
	if strings.Contains(style, "@apply pl-1.5") || !strings.Contains(style, "ps-1.5") {
		t.Fatal("style.css not migrated")
	}
	// Idempotent: a second run changes nothing and reports nothing to do.
	snapshot := map[string]string{}
	for _, rel := range []string{"ui/native-select.gsx", "ui/table.gsx", "web/gsxui/style.css", "gsxui.json"} {
		snapshot[rel] = readFile(t, dir, rel)
	}
	if err := Run([]string{"migrate", "rtl"}); err != nil {
		t.Fatal(err)
	}
	for rel, before := range snapshot {
		if readFile(t, dir, rel) != before {
			t.Errorf("second migrate changed %s", rel)
		}
	}
	// The managed hash now matches the migrated content: add --overwrite is a no-op on native-select.
	if err := Run([]string{"add", "native-select"}); err != nil {
		t.Fatal(err)
	}
	if readFile(t, dir, "ui/native-select.gsx") != snapshot["ui/native-select.gsx"] {
		t.Fatal("add after migrate re-vendored a file whose hash should already match")
	}
}

func TestMigrateListsMigrations(t *testing.T) {
	err := Run([]string{"migrate"})
	if err == nil || !strings.Contains(err.Error(), "rtl") {
		t.Fatalf("bare migrate must name the available migrations, got %v", err)
	}
	if err := Run([]string{"migrate", "nope"}); err == nil {
		t.Fatal("unknown migration accepted")
	}
}
```

- [ ] **Step 2: Run to verify they fail**

Run: `go test ./internal/cli -run TestMigrate -count=1`
Expected: FAIL, `unknown command "migrate"`.

- [ ] **Step 3: Implement**

`internal/cli/run.go`: add `case "migrate": return runMigrate(args[1:])`, and extend both usage strings to `init|add|apply|migrate|list`.

`internal/cli/migrate.go`:

```go
package cli

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gsxhq/gsxui/internal/rtl"
)

// runMigrate dispatches `gsxui migrate <name>`. Migrations rewrite files a
// consumer already owns, in place and preserving their edits, which `add
// --overwrite` cannot do.
func runMigrate(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: gsxui migrate <rtl>\n\n  rtl  rewrite vendored components to logical direction classes (requires \"rtl\": true in gsxui.json)")
	}
	switch args[0] {
	case "rtl":
		return runMigrateRTL(args[1:])
	default:
		return fmt.Errorf("unknown migration %q (want rtl)", args[0])
	}
}

// runMigrateRTL applies internal/rtl to every managed .gsx under cfg.UI and
// to the vendored style.css, through the same transaction add uses, so the
// managed hashes move to the migrated content and rollback stays possible.
func runMigrateRTL(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: gsxui migrate rtl")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := recoverArtifactTransaction(dir); err != nil {
		return err
	}
	cfg, err := LoadConfig(dir)
	if err != nil {
		return err
	}
	if !cfg.RTL {
		return fmt.Errorf(`gsxui.json does not set "rtl": true — run 'gsxui init --rtl' first, then 'gsxui migrate rtl'`)
	}
	styleCSS := filepath.ToSlash(filepath.Join(filepath.Dir(cfg.CSS), "style.css"))
	paths := make([]string, 0, len(cfg.Managed))
	for rel := range cfg.Managed {
		paths = append(paths, rel)
	}
	sort.Strings(paths)

	var artifacts []artifact
	for _, rel := range paths {
		isGSX := strings.HasSuffix(rel, ".gsx")
		if !isGSX && rel != styleCSS {
			continue
		}
		path, err := artifactPath(dir, rel)
		if err != nil {
			return err
		}
		current, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue // a managed file the consumer deleted is theirs to have deleted
			}
			return err
		}
		var next []byte
		if isGSX {
			next, err = rtl.GSX(rel, current)
			if err != nil {
				return err
			}
		} else {
			next = rtl.CSS(current)
		}
		if bytes.Equal(next, current) {
			continue
		}
		artifacts = append(artifacts, artifact{RelativePath: rel, Content: next, Managed: true})
	}
	if len(artifacts) == 0 {
		fmt.Println("migrate rtl: nothing to do")
		return nil
	}
	_, plan, err := artifactPlanWithConfig(cfg, artifacts)
	if err != nil {
		return err
	}
	// overwrite=true: rewriting files the consumer has edited is the point.
	if err := validateArtifactPlan(dir, cfg, plan, true); err != nil {
		return err
	}
	fmt.Printf("migrate rtl: %d file(s)\n", len(artifacts))
	return executeArtifactTransaction(
		dir,
		plan,
		func() error { return generateProject(dir) },
		func() error { return generateProject(dir) },
	)
}
```

If `artifactPath` or `validateArtifactPlan` have different signatures than shown (check `internal/cli/managed.go`), adapt the calls; the behaviour must stay: transactional, hashes updated, `overwrite` semantics.

- [ ] **Step 4: Run the tests**

Run: `go test ./internal/cli -count=1 -short`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/cli/migrate.go internal/cli/migrate_test.go internal/cli/run.go
git commit -m "feat(cli): gsxui migrate rtl rewrites vendored files in place

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 7: `origin-*` utilities, `site/uirtl`, RTL demos on it, and the chevron test

**Files:**
- Modify: `assets/css/foundation.css` (six `@utility` rules)
- Modify: `internal/stylegen/generate.go` (emit `site/uirtl/<c>.gsx`), `internal/stylegen/generate_test.go`
- Modify: `site/examples/rtl/login.gsx`, `site/examples/calendar/rtl.gsx`, `site/examples/pagination/rtl.gsx`, `site/examples/sidebar/rtl.gsx` (import `uirtl`)
- Modify: `jstest/specs/rtl.spec.ts`
- Generated: `site/uirtl/*.gsx`, `site/uirtl/*.x.go`, example `.x.go`, `site/hl/blocks.gen.go`

**Interfaces:**
- Produces: Go package `github.com/gsxhq/gsxui/site/uirtl` with every component of `ui/`, transformed.

- [ ] **Step 1: Add the utilities**

Append to `assets/css/foundation.css`, at top level (Tailwind requires `@utility` outside `@layer`), with a comment naming the reason (Tailwind has no logical `transform-origin`; upstream's transform emits these names):

```css
/* Logical transform-origin utilities. Tailwind ships none; shadcn's RTL
   transform (and gsxui's internal/rtl) rewrites origin-top-left to
   origin-top-start and so on, so the names must exist. */
@utility origin-top-start { &:dir(ltr) { transform-origin: top left; } &:dir(rtl) { transform-origin: top right; } }
@utility origin-top-end { &:dir(ltr) { transform-origin: top right; } &:dir(rtl) { transform-origin: top left; } }
@utility origin-bottom-start { &:dir(ltr) { transform-origin: bottom left; } &:dir(rtl) { transform-origin: bottom right; } }
@utility origin-bottom-end { &:dir(ltr) { transform-origin: bottom right; } &:dir(rtl) { transform-origin: bottom left; } }
@utility origin-start { &:dir(ltr) { transform-origin: left; } &:dir(rtl) { transform-origin: right; } }
@utility origin-end { &:dir(ltr) { transform-origin: right; } &:dir(rtl) { transform-origin: left; } }
```

Run `make audit` (it runs `stylegen --check-layers`, which compiles foundation.css); if the layer check needs the utilities registered somewhere, follow its error message and record what you registered.

- [ ] **Step 2: Write the failing stylegen test**

Append to `internal/stylegen/generate_test.go`:

```go
func TestGenerateAllEmitsTheSiteRTLPackage(t *testing.T) {
	root := t.TempDir()
	copyRepoFixture(t, root)
	if err := GenerateAll(root, false); err != nil {
		t.Fatal(err)
	}
	ns, err := os.ReadFile(filepath.Join(root, "site", "uirtl", "native-select.gsx"))
	if err != nil {
		t.Fatalf("missing site/uirtl/native-select.gsx: %v", err)
	}
	if !strings.HasPrefix(string(ns), "package uirtl\n") {
		t.Errorf("package clause = %q", firstLineOf(ns))
	}
	if strings.Contains(string(ns), "right-2.5") || !strings.Contains(string(ns), "end-2.5") {
		t.Errorf("site/uirtl/native-select.gsx is not transformed")
	}
	sheet, _ := os.ReadFile(filepath.Join(root, "site", "uirtl", "sheet.gsx"))
	if !strings.Contains(string(sheet), `"inset-y-0 left-0 h-full`) {
		t.Errorf("sheet side arm must stay physical in site/uirtl")
	}
	if err := GenerateAll(root, true); err != nil {
		t.Fatalf("GenerateAll(check) after write = %v", err)
	}
}

func firstLineOf(b []byte) string {
	line, _, _ := strings.Cut(string(b), "\n")
	return line
}
```

Also extend `copyRepoFixture`'s directory list with `filepath.Join("site", "uirtl")` if the check step needs the directory present (it is created by the write step, so add it only if the test fails on a missing dir).

- [ ] **Step 3: Run to verify it fails**

Run: `go test ./internal/stylegen -run TestGenerateAllEmitsTheSiteRTLPackage -count=1`
Expected: FAIL, file missing.

- [ ] **Step 4: Emit `site/uirtl`**

In `internal/stylegen/generate.go` `resolveAll`, inside `if style == DefaultStyle {` right after the `ui/<component>.gsx` output is appended:

```go
				rtlSource, err := rtl.GSX(canonicalPath, generated)
				if err != nil {
					return nil, fmt.Errorf("derive %s site rtl source: %w", component, err)
				}
				uirtl, err := rewriteGSXPackage(canonicalPath, rtlSource, "uirtl")
				if err != nil {
					return nil, fmt.Errorf("derive %s site rtl package: %w", component, err)
				}
				outputs = append(outputs, generatedSource{
					relativePath: filepath.Join("site", "uirtl", component+".gsx"),
					content:      uirtl,
				})
```

Import `"github.com/gsxhq/gsxui/internal/rtl"`. Then:

```bash
go run ./cmd/stylegen && go tool gsx generate && go build ./... && go test ./internal/stylegen -count=1
```

Expected: `site/uirtl/*.gsx` and `*.x.go` appear, the package compiles, tests pass.

- [ ] **Step 5: Point the RTL demos at `uirtl` and add the NativeSelect**

In `site/examples/rtl/login.gsx`, `site/examples/calendar/rtl.gsx`, `site/examples/pagination/rtl.gsx`, `site/examples/sidebar/rtl.gsx`: replace the import `"github.com/gsxhq/gsxui/ui"` with `"github.com/gsxhq/gsxui/site/uirtl"` and every `ui.` component reference with `uirtl.` (`sed -i '' 's/\bui\./uirtl./g'` on those four files, then inspect the diff: only component references and the import must change). If a file also imports `ui` for a Go type (`ui.CalendarLocale`), keep that import and use `uirtl.CalendarLocale` only if the value is passed to a `uirtl` component.

In `login.gsx`, after the password field, add:

```gsx
					<div class="flex flex-col gap-2">
						<uirtl.Label for="rtl-login-language">اللغة</uirtl.Label>
						<uirtl.NativeSelect id="rtl-login-language" name="language">
							<uirtl.NativeSelectOption value="ar" selected={true}>العربية</uirtl.NativeSelectOption>
							<uirtl.NativeSelectOption value="en">English</uirtl.NativeSelectOption>
						</uirtl.NativeSelect>
					</div>
```

Update the login demo's header comment: the parts adapt because the file renders the transformed `site/uirtl` package, the same output `gsxui add` produces for a project with `"rtl": true`. Then `go tool gsx fmt -w` the four files, `go tool gsx generate`, `go build ./...`, `make highlight`, `go test ./site/... -count=1`.

- [ ] **Step 6: Write the failing Playwright test**

Append inside the `test.describe("rtl", …)` block of `jstest/specs/rtl.spec.ts`:

```ts
  test("native-select chevron sits at the logical end (visual left) in the RTL login demo", async ({ page }) => {
    // The demo renders from site/uirtl, the transform's output. Issue #32's
    // reported symptom was this chevron staying on the physical right.
    await page.goto("/x/rtl/login");
    const wrapper = page.locator("[data-gsxui-slot-native-select-wrapper]").first();
    await expect(wrapper).toBeVisible();
    const chevron = wrapper.locator("> svg");
    const wrapperBox = await wrapper.boundingBox();
    const chevronBox = await chevron.boundingBox();
    if (!wrapperBox || !chevronBox) throw new Error("missing bounding boxes");
    expect(chevronBox.x + chevronBox.width / 2).toBeLessThan(wrapperBox.x + wrapperBox.width / 2);
  });
```

Before Step 5 this fails (physical `right-2.5`); after it, it passes. Run: `npx playwright test --config jstest/playwright.config.ts jstest/specs/rtl.spec.ts`
Expected: all RTL cases pass, including the four pre-existing sidebar/sheet physical-side tests.

- [ ] **Step 7: Gates and commit**

Run: `go run ./cmd/stylegen --check && go run ./cmd/stylegen --check-authoring && make audit && make verify-generated && gofmt -l . && go test ./site/... ./internal/stylegen -count=1`
Expected: clean.

```bash
git add assets/css/foundation.css internal/stylegen/generate.go internal/stylegen/generate_test.go site/uirtl site/examples/rtl site/examples/calendar/rtl.gsx site/examples/calendar/rtl.x.go site/examples/pagination/rtl.gsx site/examples/pagination/rtl.x.go site/examples/sidebar/rtl.gsx site/examples/sidebar/rtl.x.go site/hl/blocks.gen.go jstest/specs/rtl.spec.ts
git commit -m "feat(site): render the RTL demos from a transformed site/uirtl package

Adds the logical transform-origin utilities the transform emits, and a
NativeSelect to the Arabic login demo with a Playwright assertion that its
chevron sits at the logical end.

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

Before adding `site/uirtl`, run `git status --short site/uirtl | head` and confirm it holds only `.gsx` and `.x.go` pairs.

---

### Task 8: Documentation

**Files:**
- Modify: `site/pages/rtl.gsx` (Get started and How it works), `site/pages/getting_started.gsx` (one sentence), `docs/jsx-parity.md` (new `## rtl` section), `CHANGELOG.md`, `README.md`
- Create: `site/snippets/rtl-init.sh.txt`, `site/snippets/rtl-migrate.sh.txt`
- Generated: `site/pages/rtl.x.go`, `site/pages/getting_started.x.go`, `site/hl/blocks.gen.go`

- [ ] **Step 1: Snippets**

`site/snippets/rtl-init.sh.txt`:

```sh
gsxui init --rtl
gsxui add button card
```

`site/snippets/rtl-migrate.sh.txt`:

```sh
gsxui migrate rtl
```

- [ ] **Step 2: Rewrite the RTL page's first two sections**

In `site/pages/rtl.gsx`, replace the intro paragraph and the **Get started** section body with:

```gsx
				<p class="text-muted-foreground">
					Opt in once and the components you vendor use logical direction classes, so one build serves
					<code>dir="ltr"</code> and <code>dir="rtl"</code> alike.
				</p>
```

and, under the Get started heading:

```gsx
				<p>Set the flag when you initialise, then add components as usual:</p>
				<pre><code>{ hl.Node("snippets/rtl-init.sh") }</code></pre>
				<p>Components vendored before the flag was set are rewritten in place, keeping your edits:</p>
				<pre><code>{ hl.Node("snippets/rtl-migrate.sh") }</code></pre>
				<p>
					Set <code>dir="rtl"</code> on <code>&lt;html&gt;</code>, or on any subtree, from the locale. An
					<code>rtl: true</code> project renders both directions; the document's <code>dir</code> decides.
				</p>
				<pre><code>{ `<html lang="ar" dir="rtl">
  ...
</html>` }</code></pre>
```

Replace the **How it works** list and its two paragraphs with:

```gsx
				<ul class="list-disc space-y-2 pl-6">
					<li>
						<strong>Logical classes, at install time.</strong> With <code>"rtl": true</code> in
						<code>gsxui.json</code>, <code>gsxui add</code> rewrites physical classes in what it vendors:
						<code>ml-*</code> to <code>ms-*</code>, <code>left-*</code> to <code>start-*</code>,
						<code>text-left</code> to <code>text-start</code>, and the rest of shadcn's table.
						<code>translate-x-*</code> and <code>space-x-*</code> gain <code>rtl:</code> companions.
					</li>
					<li>
						<strong>Directional icons.</strong> Chevrons and arrows that encode a left/right meaning carry
						<code>rtl:rotate-180</code>.
					</li>
					<li>
						<strong>Direction-aware floating positioning.</strong> Popover, dropdown-menu, select, tooltip and
						the rest of the floating family resolve placement in JS at position time.
					</li>
					<li>
						<strong>Mirrored keyboard semantics.</strong> Arrow keys in menus, tabs, carousel, calendar and the
						other roving-focus components mirror by meaning per WAI-ARIA.
					</li>
					<li>
						<strong>The demos on this page</strong> render from the transformed components, the same output
						<code>gsxui add</code> produces with the flag set.
					</li>
				</ul>
				<p>Not transformed:</p>
				<ul class="list-disc space-y-2 pl-6">
					<li>
						Sheet, Drawer and Sidebar <code>side="left"</code>/<code>side="right"</code> stay physical, matching
						shadcn's <code>data-side</code> contract; their interiors mirror.
					</li>
					<li><code>input-otp</code>'s digit group stays pinned <code>dir="ltr"</code>.</li>
					<li>Behaviour JS carries no classes and is untouched.</li>
					<li>
						Sized logical slide utilities (<code>slide-in-from-start-2</code>) match only elements that
						themselves carry <code>dir</code> in the current tw-animate-css build; gsxui emits none.
					</li>
				</ul>
```

Keep the "Try it out" paragraph that says more RTL variants live on component pages, and the Fonts section. Then `go tool gsx fmt -w site/pages/rtl.gsx`, `go tool gsx generate`, `make highlight`, and check the generated `rtl.x.go` for dropped spaces after closing tags (`grep -nE '</code>[A-Za-z]|[A-Za-z]<code>' site/pages/rtl.x.go`); fix with `{ " " }` where needed.

- [ ] **Step 3: Getting started, parity ledger, changelog, README**

`site/pages/getting_started.gsx`: in the initialize step, after the `init.sh` snippet, add one sentence: `<p>Building for right-to-left languages? Pass <code>--rtl</code>; see <a href={Rtl{} |> url}>RTL</a>.</p>`.

`docs/jsx-parity.md`: append:

```markdown
## rtl

- MECHANISM: RTL is a project-level opt-in, `"rtl": true` in `gsxui.json` (`gsxui init --rtl`), mirroring shadcn's `components.json`. The registry stays physical. `internal/rtl` carries upstream's class table from `packages/shadcn/src/utils/transformers/transform-rtl.ts` at the pinned commit (the same pin as `registry/styles/*/*.css` headers) and rewrites class literals in vendored `.gsx` and `@apply` lists in vendored `style.css` at `gsxui add`, `gsxui apply`, and `gsxui migrate rtl` (in place, edits preserved, idempotent). `site/uirtl/` is stylegen's transformed copy of `ui/` for the site's RTL demos, as upstream's `ui-rtl` is for theirs.
- ADAPT (side stays physical in full): upstream skips only positioning prefixes inside `data-[side=left|right]` variants and still maps `border-r` to `border-e`. gsxui leaves every token in a side-keyed arm or variant untouched (`switch side` arms in Sheet and Drawer, `[data-side=…]` selectors in Sidebar), per the physical-side ruling `jstest/specs/rtl.spec.ts` pins.
- ADAPT (no `side` prop rewrite): upstream maps `side="left"` to `inline-start` for Base UI menu sub-content; gsxui's floating positioning resolves placement per direction in JS, so `side` values are unchanged.
- ADAPT (`origin-*`): `origin-top-start` and siblings are not Tailwind utilities; `assets/css/foundation.css` defines them with `:dir()`.
- GAP: tw-animate-css's minified build drops a descendant space in sized logical slides (`slide-in-from-start-2`), Wombosvideo/tw-animate-css#67; gsxui emits none, callers writing their own must pass `dir` to the element.
```

`CHANGELOG.md`: add a `## 2026-09-16` section at the top:

```markdown
## 2026-09-16

### Added

- **rtl** — `"rtl": true` in `gsxui.json` (`gsxui init --rtl`) makes `gsxui add` and `gsxui apply` vendor components with logical direction classes, using shadcn's `transform-rtl` table; `gsxui migrate rtl` rewrites components vendored earlier in place, keeping your edits. One build then serves `dir="ltr"` and `dir="rtl"` (#32).
- **foundation.css** — `origin-top-start`, `origin-top-end`, `origin-bottom-start`, `origin-bottom-end`, `origin-start`, `origin-end` utilities, which the transform emits.

### Changed

- **site** — the RTL page describes the opt-in model; its demos render from `site/uirtl`, the transform's output. The page previously stated components used logical classes by default, which was not true of the registry.
```

`README.md`: after the vendoring paragraph, add: `Building for right-to-left languages? \`gsxui init --rtl\` vendors logical direction classes; see [RTL](https://ui.gsxhq.dev/docs/rtl).`

- [ ] **Step 4: Gates and commit**

Run: `go test ./site/... -count=1 && gofmt -l . && make audit`
Expected: clean; `TestDocsTableOfContents` for `/docs/rtl` still passes (headings unchanged).

```bash
git add site/pages/rtl.gsx site/pages/rtl.x.go site/pages/getting_started.gsx site/pages/getting_started.x.go site/snippets/rtl-init.sh.txt site/snippets/rtl-migrate.sh.txt site/hl/blocks.gen.go docs/jsx-parity.md CHANGELOG.md README.md
git commit -m "docs(rtl): describe the opt-in transform; parity ledger, changelog, readme

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

---

### Task 9: End-to-end proof and the full gate run

**Files:**
- Modify: `internal/cli/e2e_test.go` (new test)

- [ ] **Step 1: Write the e2e test**

Append to `internal/cli/e2e_test.go`:

```go
// TestE2ERTL is issue #32's table, end to end: a scratch consumer with
// "rtl": true vendors every component the issue lists and none of the
// physical classes it names survive. The module still builds.
func TestE2ERTL(t *testing.T) {
	if testing.Short() {
		t.Skip("network-dependent e2e; run without -short")
	}
	t.Setenv("GOWORK", "off")
	dir := scaffoldGSXProject(t)
	t.Chdir(dir)
	if err := Run([]string{"init", "--rtl"}); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add",
		"native-select", "select", "combobox", "dropdown-menu", "context-menu", "menubar", "command",
		"navigation-menu", "accordion", "table", "alert", "alert-dialog", "field", "item", "input-group",
		"sidebar", "toggle-group", "tooltip",
	}); err != nil {
		t.Fatal(err)
	}
	gone := map[string][]string{
		"native-select.gsx":   {"right-2.5", "pr-8"},
		"select.gsx":          {"text-left", "pr-8 pl-1.5", "right-2"},
		"combobox.gsx":        {"pr-8 pl-1.5", "right-2"},
		"dropdown-menu.gsx":   {"pr-8 pl-1.5", "data-inset:pl-7", "ml-auto"},
		"context-menu.gsx":    {"data-inset:pl-7", "ml-auto"},
		"menubar.gsx":         {"pl-7 pr-1.5", "left-1.5", "ml-auto"},
		"command.gsx":         {"pl-2", "ml-auto"},
		"navigation-menu.gsx": {"ml-1"},
		"accordion.gsx":       {"text-left", "ml-auto"},
		"table.gsx":           {"text-left", "pr-0"},
		"alert.gsx":           {"text-left"},
		"alert-dialog.gsx":    {"text-left"},
		"field.gsx":           {"text-left"},
		"item.gsx":            {"text-left"},
		"input-group.gsx":     {"pr-1.5", "pl-2"},
		"toggle-group.gsx":    {"rounded-l-lg", "rounded-r-lg"},
		"tooltip.gsx":         {"has-[kbd]:pr-1.5"},
	}
	for file, classes := range gone {
		src := readFile(t, dir, filepath.Join("ui", file))
		for _, c := range classes {
			if strings.Contains(src, c) {
				t.Errorf("%s still carries %q", file, c)
			}
		}
	}
	// sidebar: the rail border is side-keyed and stays physical; the menu button text goes logical.
	sidebar := readFile(t, dir, "ui/sidebar.gsx")
	if !strings.Contains(sidebar, "[data-side=left]>&]:border-r") {
		t.Error("sidebar rail border must stay physical")
	}
	if strings.Contains(sidebar, `"text-left`) || strings.Contains(sidebar, ` text-left`) {
		t.Error("sidebar menu button text-left survived")
	}
	if !strings.Contains(readFile(t, dir, "ui/native-select.gsx"), "end-2.5") {
		t.Error("native-select missing end-2.5")
	}
	mustRun(t, dir, "go", "build", "./...")
}
```

If an entry in `gone` is wrong because the shipped nova sheet spells a class differently from issue #32's table (the issue read `main` at `f0818be`), correct the expectation to the class the vendored file actually carried before the transform (check `ui/<file>.gsx` in this repo) and note it in the report. Never delete an entry.

- [ ] **Step 2: Run it**

Run: `go test ./internal/cli -run TestE2ERTL -count=1 -v`
Expected: PASS.

- [ ] **Step 3: Full gates on the committed tree**

```bash
go build ./...
go test ./... -count=1
go run ./cmd/stylegen --check
go run ./cmd/stylegen --check-authoring
go run ./cmd/stylegen --check-layers
make audit
make test-css-audit
make test-theme-state
make verify-generated
make verify-generated-styles
gofmt -l .
npx playwright test --config jstest/playwright.config.ts
```

Expected: all clean. Known parallel-load flakes: `home-showcase.spec.ts:16`, `dialog.spec.ts:453/503`; rerun a failing one alone before treating it as a regression.

- [ ] **Step 4: Manual check**

Run `make site-dev`; open `/docs/rtl`: the login demo's language select shows its chevron on the visual left, the Get started snippets render, and the "Not transformed" list appears. Open `/components/sidebar` RTL example and confirm the rail still sits on the physical side.

- [ ] **Step 5: Commit and finish**

```bash
git add internal/cli/e2e_test.go
git commit -m "test(cli): end-to-end rtl vendoring against issue #32's table

Claude-Session: https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv"
```

Then follow `superpowers:finishing-a-development-branch`. The PR closes #32 and ends with `https://claude.ai/code/session_019u7BSQ2Kt9bkwzAXN9hGwv`.
