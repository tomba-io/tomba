package version

import "fmt"

var (
	app     = "email"
	version = "v1.0.1"
)

// String returns a string.
func String() string {
	return fmt.Sprintf("%s %s", app, version)
}
