package diag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
)

// FormatError renders a concise, actionable error message with optional retry hints.
func FormatError(err error) string {
	if err == nil {
		return ""
	}
	var he *httpc.Error
	if errors.As(err, &he) {
		// Provide hints for common cases
		switch he.Status {
		case 429:
			return fmt.Sprintf("Too many requests. Try again later. (%s)", he.Err)
		case 401, 403:
			return fmt.Sprintf("Authentication failed. Check your token or re-login. (%s)", he.Err)
		case 400, 422:
			return fmt.Sprintf("Validation error: %s", he.Err)
		default:
			if he.Code == 4 {
				return fmt.Sprintf("Network error. Check connectivity and server URL. (%s)", he.Err)
			}
			return he.Err.Error()
		}
	}
	// generic fallback
	s := err.Error()
	s = strings.TrimSpace(s)
	if s == "" {
		s = "unknown error"
	}
	return s
}
