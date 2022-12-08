package outputhandler

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"
)

var gotColors = true
var origTerminalMode uint32

// Initialize sets up the output for color text
func Initialize() {
	if runtime.GOOS == "windows" {
		// Enable Virtual Terminal Processing by adding the flag to the current mode
		fd := windows.Handle(os.Stdout.Fd())
		if err := windows.GetConsoleMode(fd, &origTerminalMode); err != nil {
			fmt.Printf("Warning: couldn't enable Virtual Terminal Processing: %v", err)
			gotColors = false
		} else {
			if err = windows.SetConsoleMode(fd, origTerminalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING); err != nil {
				fmt.Printf("Warning: couldn't enable Virtual Terminal Processing: %v", err)
				gotColors = false
			}
		}
	}
}

// Reset sets the terminal mode back how Initialize() found it
func Reset() {
	fmt.Println(ResetColor)
	if runtime.GOOS == "windows" {
		// reset terminal mode
		fd := windows.Handle(os.Stdout.Fd())
		windows.SetConsoleMode(fd, origTerminalMode)
	}
}

// TerminalColor is the actual terminal color values for bash
type TerminalColor string

const (
	ResetColor    TerminalColor = "\033[0m"
	DefaultColor  TerminalColor = "\033[39m"
	Red           TerminalColor = "\033[31m"
	Green         TerminalColor = "\033[32m"
	Yellow        TerminalColor = "\033[33m"
	Blue          TerminalColor = "\033[34m"
	Magenta       TerminalColor = "\033[35m"
	Cyan          TerminalColor = "\033[36m"
	BrightRed     TerminalColor = "\033[91m"
	BrightGreen   TerminalColor = "\033[92m"
	BrightYellow  TerminalColor = "\033[93m"
	BrightBlue    TerminalColor = "\033[94m"
	BrightMagenta TerminalColor = "\033[95m"
	BrightCyan    TerminalColor = "\033[96m"
	White         TerminalColor = "\033[97m"
	Gray          TerminalColor = "\033[37m"
	DarkGray      TerminalColor = "\033[90m"
	Black         TerminalColor = "\033[30m"
)

// GetTerminalColor returns the format string for the requested color.
// Note: some terminals may not make use of / correctly implement CSI
func GetForeground(color TerminalColor) string {
	if !gotColors {
		return ""
	}
	return string(color)
}
