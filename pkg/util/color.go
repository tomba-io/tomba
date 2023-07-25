package util

import (
	"github.com/mgutz/ansi"
)

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
	return bold(t)
}

// White text
func White(t string) string {
	return white(t)
}

// Red text
func Red(t string) string {
	return red(t)
}

// Yellow text
func Yellow(t string) string {
	return yellow(t)
}

// Green text
func Green(t string) string {
	return green(t)
}

// Gray text
func Gray(t string) string {
	return gray(t)
}

// Magenta text
func Magenta(t string) string {
	return magenta(t)
}

// Cyan text
func Cyan(t string) string {
	return cyan(t)
}

// Blue text
func Blue(t string) string {
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
