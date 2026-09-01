package util

import (
	"os"

	"github.com/mgutz/ansi"
)

// NoColor disables colored output when set to true.
var NoColor bool

func init() {
	if os.Getenv("NO_COLOR") != "" {
		NoColor = true
	}
}

var (
	magenta = ansi.ColorFunc("magenta")
	cyan    = ansi.ColorFunc("cyan")
	red     = ansi.ColorFunc("red")
	yellow  = ansi.ColorFunc("yellow")
	blue    = ansi.ColorFunc("blue")
	green   = ansi.ColorFunc("green")
	gray    = ansi.ColorFunc("black+h")
	white   = ansi.ColorFunc("white")
	bold    = ansi.ColorFunc("default+b")
)

// Bold text
func Bold(t string) string {
	if NoColor {
		return t
	}
	return bold(t)
}

// White text
func White(t string) string {
	if NoColor {
		return t
	}
	return white(t)
}

// Red text
func Red(t string) string {
	if NoColor {
		return t
	}
	return red(t)
}

// Yellow text
func Yellow(t string) string {
	if NoColor {
		return t
	}
	return yellow(t)
}

// Green text
func Green(t string) string {
	if NoColor {
		return t
	}
	return green(t)
}

// Gray text
func Gray(t string) string {
	if NoColor {
		return t
	}
	return gray(t)
}

// Magenta text
func Magenta(t string) string {
	if NoColor {
		return t
	}
	return magenta(t)
}

// Cyan text
func Cyan(t string) string {
	if NoColor {
		return t
	}
	return cyan(t)
}

// Blue text
func Blue(t string) string {
	if NoColor {
		return t
	}
	return blue(t)
}

// SuccessIcon icon
func SuccessIcon() string {
	return Green("✓")
}

// WarningIcon icon
func WarningIcon() string {
	return Yellow("!")
}

// ErrorIcon icon
func ErrorIcon() string {
	return Red("❌")
}

// InfoIcon icon
func InfoIcon() string {
	return Blue("ℹ️ ")
}
