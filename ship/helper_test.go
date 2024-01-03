package ship

import (
	"os"
	"testing"
)

func skipCI(t *testing.T) {
	if os.Getenv("ACTION_ENVIRONMENT") == "CI" {
		t.Skip("Skipping testing in CI environment")
	}
}
