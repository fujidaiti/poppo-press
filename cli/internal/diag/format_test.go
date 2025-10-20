package diag

import (
	"errors"
	"testing"

	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
)

func TestFormatError_HTTPAndNetwork(t *testing.T) {
	if got := FormatError(&httpc.Error{Status: 429, Err: errors.New("rate limited")}); got == "" || got[0] == '(' {
		t.Fatalf("expected hint for 429, got %q", got)
	}
	if got := FormatError(&httpc.Error{Status: 401, Err: errors.New("unauthorized"), Code: 3}); got == "" {
		t.Fatalf("expected auth hint, got %q", got)
	}
	if got := FormatError(&httpc.Error{Code: 4, Err: errors.New("conn refused")}); got == "" {
		t.Fatalf("expected network hint, got %q", got)
	}
}
