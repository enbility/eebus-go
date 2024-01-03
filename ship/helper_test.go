package ship

import (
	"os"
)

func isRunningOnCI() bool {
	return os.Getenv("ACTION_ENVIRONMENT") == "CI"
}
