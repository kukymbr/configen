package version

import "fmt"

// These variables are populated during the build.
var (
	Version  = "unknown"
	Revision = "unknown"
	BuiltAt  = "2025-10-04 00:00:00"
)

func GetVersion() string {
	return fmt.Sprintf("%s (revision %s, built at %s)", Version, Revision, BuiltAt)
}
