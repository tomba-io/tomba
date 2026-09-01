package version

import (
	"fmt"
	"runtime"
)

// These variables are set at build time via ldflags.
var (
	Version   = "dev"
	BuildDate = "unknown"
	Commit    = "unknown"
)

// String returns full version info with build metadata.
func String() string {
	return fmt.Sprintf("tomba %s\n  built:    %s\n  commit:   %s\n  go:       %s\n  os/arch:  %s/%s",
		Version, BuildDate, Commit, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
