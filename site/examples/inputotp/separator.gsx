package inputotp

import "github.com/gsxhq/gsxui/ui"

// Separator mirrors shadcn's own input-otp-separator.tsx: three 2-slot
// groups split by two InputOTPSeparators — the strongest proof case for
// DOM-order index stamping (ui/input-otp.js walks every slot in source
// order across all three group boundaries and stamps data-index 0..5
// positionally at mount; see ui/input-otp.gsx's own ADAPT doc comment).
component Separator() {
	<ui.InputOTP maxlength="6">
		<ui.InputOTPGroup>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
		</ui.InputOTPGroup>
		<ui.InputOTPSeparator/>
		<ui.InputOTPGroup>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
		</ui.InputOTPGroup>
		<ui.InputOTPSeparator/>
		<ui.InputOTPGroup>
			<ui.InputOTPSlot/>
			<ui.InputOTPSlot/>
		</ui.InputOTPGroup>
	</ui.InputOTP>
}
