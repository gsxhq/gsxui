package canonical

import (
	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// inputOtp binds InputOTP's shape to the accessor calls input-otp.gsx
// authors. The receiver is camel (inputOtp), derived from the hyphenated
// registry name "input-otp" by componentIdentifier.
var inputOtp = inputOtpRecipe{c: recipe.Component{Shape: shapes.InputOTP}}
