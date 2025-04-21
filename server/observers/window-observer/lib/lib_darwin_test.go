//go:build darwin
// +build darwin

package lib

import "testing"

func TestGetCurrentWindowMacOS_InvalidStrategy(t *testing.T) {
	_, err := GetCurrentWindowMacOS("not-a-strategy")
	fe, ok := err.(FatalError)
	if !ok {
		t.Fatalf("expected FatalError, got %T", err)
	}
	want := `invalid strategy "not-a-strategy"`
	if fe.Error() != want {
		t.Errorf("Error() = %q; want %q", fe.Error(), want)
	}
}

func TestGetCurrentWindowMacOS_JXA_and_AppleScript(t *testing.T) {
	t.Skip("needs mocking of GetInfo and GetInfoAppleScript")
}
