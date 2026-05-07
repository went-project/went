package output

import (
	"fmt"
	"os"
	"strings"
)

const (
	resetColor   = "\033[0m"
	bold         = "\033[1m"
	infoColor    = "\033[36m"
	successColor = "\033[32m"
	warningColor = "\033[33m"
	errorColor   = "\033[31m"
	accentColor  = "\033[35m"

	infoIcon    = "ℹ"
	successIcon = "✓"
	warningIcon = "⚠"
	errorIcon   = "✖"
)

func colorize(color, message string) string {
	return color + message + resetColor
}

func boldText(message string) string {
	return bold + message + resetColor
}

func PrintHeader(title string) {
	border := strings.Repeat("─", len(title)+4)
	fmt.Fprintln(os.Stdout, accentColor+"╭"+border+"╮"+resetColor)
	fmt.Fprintf(os.Stdout, "%s│  %s  │%s\n", accentColor, boldText(title), resetColor)
	fmt.Fprintln(os.Stdout, accentColor+"╰"+border+"╯"+resetColor)
}

func PrintLine(label, value string) {
	fmt.Fprintf(os.Stdout, "%s%s:%s %s\n", accentColor, label, resetColor, value)
}

func PrintInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", infoColor, infoIcon, resetColor, fmt.Sprintf(format, args...))
}

func PrintSuccess(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", successColor, successIcon, resetColor, fmt.Sprintf(format, args...))
}

func PrintWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s%s%s %s\n", warningColor, warningIcon, resetColor, fmt.Sprintf(format, args...))
}

func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "%s%s%s %s\n", errorColor, errorIcon, resetColor, fmt.Sprintf(format, args...))
}

func PrintStep(icon, format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, "  %s %s\n", icon, fmt.Sprintf(format, args...))
}

func PrintSpacing() {
	fmt.Fprintln(os.Stdout)
}
