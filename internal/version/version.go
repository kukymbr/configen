package version

import "fmt"

// These variables are populated during the build.
var (
	Version  = "unknown"
	Revision = "unknown"
	BuiltAt  = "{% now 'utc', '%Y-%m-%d %H:%M:%S' %}"
)

func GetVersion() string {
	return fmt.Sprintf("%s (revision %s, built at %s)", Version, Revision, BuiltAt)
}
