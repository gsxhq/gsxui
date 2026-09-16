// Package kbd holds the site's example gsx components for ui/kbd.
package kbd

import "github.com/gsxhq/gsxui/site/uirtl"

// Rtl mirrors shadcn's own kbd-rtl demo: a KbdGroup of four modifier keys
// followed by a Ctrl+B compound shortcut, wrapped in dir="rtl". Kbd/
// KbdGroup carry no user-facing English text, so nothing needs translating
// — the dir="rtl" wrapper exercises the RTL layout on its own.
component Rtl() {
	<div dir="rtl" lang="ar" class="flex flex-col items-center gap-4">
		<uirtl.KbdGroup>
			<uirtl.Kbd>⌘</uirtl.Kbd>
			<uirtl.Kbd>⇧</uirtl.Kbd>
			<uirtl.Kbd>⌥</uirtl.Kbd>
			<uirtl.Kbd>⌃</uirtl.Kbd>
		</uirtl.KbdGroup>
		<uirtl.KbdGroup>
			<uirtl.Kbd>Ctrl</uirtl.Kbd>
			<span>+</span>
			<uirtl.Kbd>B</uirtl.Kbd>
		</uirtl.KbdGroup>
	</div>
}
