// lib_common_test.go
package lib

import "testing"

func ptr(s string) *string { return &s }

func TestStringPtr(t *testing.T) {
	tests := []struct {
		in   string
		want *string
	}{
		{"", nil},
		{"hello", ptr("hello")},
	}
	for _, tt := range tests {
		got := StringPtr(tt.in)
		if tt.want == nil {
			if got != nil {
				t.Errorf("StringPtr(%q) = %v; want nil", tt.in, got)
			}
		} else {
			if got == nil || *got != *tt.want {
				t.Errorf("StringPtr(%q) = %v; want %v", tt.in, got, tt.want)
			}
		}
	}
}

func TestBoolPtr(t *testing.T) {
	b1 := boolPtr(true)
	if b1 == nil || *b1 != true {
		t.Errorf("boolPtr(true) = %v; want pointer to true", b1)
	}
	b2 := boolPtr(false)
	if b2 == nil || *b2 != false {
		t.Errorf("boolPtr(false) = %v; want pointer to false", b2)
	}
}

func TestFatalError_Error(t *testing.T) {
	err := FatalError{"something went wrong"}
	if err.Error() != "something went wrong" {
		t.Errorf("FatalError.Error() = %q; want %q", err.Error(), "something went wrong")
	}
}
