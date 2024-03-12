package version

import "fmt"

var (
	app     = "tomba"
	version = "v1.0.6"
)

// String returns a string.
func String() string {
	return fmt.Sprintf("%s %s", app, version)
}
