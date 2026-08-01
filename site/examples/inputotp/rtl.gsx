package inputotp

import "github.com/gsxhq/gsxui/ui"

// Rtl wraps Basic's own 6-digit/2-group OTP layout in an Arabic dir="rtl"
// form context (a label plus the surrounding div), but does NOT force
// dir="rtl" onto the OTP widget itself: ui/input-otp.gsx pins dir="ltr" on
// both InputOTP's own root and InputOTPGroup by design (see its own ADAPT
// doc comments) so the digit slots always read left-to-right regardless of
// document direction, matching shadcn's own input-otp behavior under RTL.
// The label and page chrome flip to RTL; the OTP slots themselves don't.
component Rtl() {
	<div dir="rtl" lang="ar">
		<label class="mb-2 block text-sm font-medium">رمز التحقق</label>
		<ui.InputOTP maxlength="6" aria-invalid="true">
			<ui.InputOTPGroup>
				<ui.InputOTPSlot aria-invalid="true"/>
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
	</div>
}
