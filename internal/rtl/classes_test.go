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
