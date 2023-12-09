package pkg

import (
	"fmt"
	"runtime"
)

// version must be set from the contents of VERSION file by go build's
// -X main.version= option in the Makefile.
var version = "devel"

// commitSha will be the hash that the binary was built from
// and will be populated by the Makefile
var commitSha = "unknown"

func GetVersion() string {
	return fmt.Sprintf("%s (commit: %s, %s)", version, commitSha, runtime.Version())
}
