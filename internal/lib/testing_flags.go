package lib

import (
	"os"
	"strings"
)

func isRunningTest() bool {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			return true
		}
	}
	return false
}
