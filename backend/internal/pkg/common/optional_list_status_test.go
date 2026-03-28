package common

import (
	"encoding/json"
	"testing"
)

func TestOptionalListStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		raw  string
		want *int
	}{
		{`null`, nil},
		{`""`, nil},
		{`1`, intPtr(1)},
		{`"1"`, intPtr(1)},
		{`" "`, nil},
	}
	for _, tt := range tests {
		var o OptionalListStatus
		if err := json.Unmarshal([]byte(tt.raw), &o); err != nil {
			t.Fatalf("raw=%s err=%v", tt.raw, err)
		}
		got := o.Ptr()
		if tt.want == nil {
			if got != nil {
				t.Fatalf("raw=%s want nil got %v", tt.raw, *got)
			}
			continue
		}
		if got == nil || *got != *tt.want {
			t.Fatalf("raw=%s want %v got %v", tt.raw, *tt.want, got)
		}
	}
}

func intPtr(n int) *int { return &n }
