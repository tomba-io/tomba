package version

import "fmt"

var (
	app     = "tomba"
	version = "v1.1.0-next"
)

// String returns a string.
func String() string {
	return fmt.Sprintf("%s %s", app, version)
}
