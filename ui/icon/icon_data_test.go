package icon

import "testing"

// TestIconCount is an internal-package test (not icon_test) because it
// asserts directly against the generated, unexported data map — the
// per-icon lookup table has no exported enumeration API.
func TestIconCount(t *testing.T) {
	if len(data) <= 1500 {
		t.Errorf("want > 1500 icons, got %d", len(data))
	}
}
