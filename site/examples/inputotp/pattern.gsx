package inputotp

import "github.com/gsxhq/gsxui/ui"

// Pattern mirrors shadcn's own input-otp-pattern.tsx, which feeds the
// input-otp library's REGEXP_ONLY_DIGITS_AND_CHARS to both native
// validation and live keystroke filtering with a single anchored regex.
// This port needs two DIFFERENT attributes for those two different jobs —
// see ui/input-otp.gsx's own DECISION doc comment: the native HTML5
// pattern attribute stays whole-string-anchored ("[0-9]*", validity only,
// browser-native), while data-gsxui-input-otp-pattern is the UNANCHORED
// per-character class ui/input-otp.js actually tests one keystroke at a
// time against ("[0-9]", not "[0-9]*"). Copying shadcn's own anchored
// pattern constant into the per-character attribute would reject every
// character typed.
component Pattern() {
	<div class="flex flex-col gap-2">
		<ui.InputOTP maxlength="6" pattern="[0-9]*" data-gsxui-input-otp-pattern="[0-9]">
			<ui.InputOTPGroup>
				<ui.InputOTPSlot/>
				<ui.InputOTPSlot/>
				<ui.InputOTPSlot/>
				<ui.InputOTPSlot/>
				<ui.InputOTPSlot/>
				<ui.InputOTPSlot/>
			</ui.InputOTPGroup>
		</ui.InputOTP>
		<p class="text-sm text-muted-foreground">
			Digits only. data-gsxui-input-otp-pattern takes an unanchored per-character class, not a whole-string-anchored regex — see ui/input-otp.gsx's own doc comment for why.
		</p>
	</div>
}
