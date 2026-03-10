package version

import "fmt"

var (
	app     = "tomba"
	version = "v1.0.9-next"
)

// String returns a string.
func String() string {
	return fmt.Sprintf("%s %s", app, version)
}
