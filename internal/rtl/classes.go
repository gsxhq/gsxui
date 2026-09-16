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
