// Package inputotp holds the site's example gsx components for
// ui/input-otp. "input-otp" can't be a Go package name (hyphen), so the
// directory drops it — same selectbox/switchctl-style workaround as
// select/switch, not a Go-keyword issue this time. The registered example
// key stays the hyphenated "input-otp" (see inputotp.go).
package inputotp

import "github.com/gsxhq/gsxui/ui"

// Basic mirrors shadcn's own input-otp-demo.tsx: a 6-digit code split into
// two 3-slot groups by one InputOTPSeparator. maxlength="6" falls through
// via attrs onto the real hidden input (see ui/input-otp.gsx's own doc
// comment — no explicit Go param for it, same convention as ui/input.gsx).
// No InputOTPSlot here takes an index: ui/input-otp.js stamps data-index
// 0..5 positionally in DOM order at mount, spanning the separator boundary
// — the deliberate API departure from shadcn's own index prop, see
// ui/input-otp.gsx's own ADAPT doc comment.
component Basic() {
	<ui.InputOTP maxlength="6">
		<ui.InputOTPGroup>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
		</ui.InputOTPGroup>
		<ui.InputOTPSeparator/>
		<ui.InputOTPGroup>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
		</ui.InputOTPGroup>
	</ui.InputOTP>
}
