package outputhandler

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

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
	fmt.Println(GetReset())
	if runtime.GOOS == "windows" {
		// reset terminal mode
		fd := windows.Handle(os.Stdout.Fd())
		windows.SetConsoleMode(fd, origTerminalMode)
	}
}

// TerminalColor is the actual terminal color values for bash
type TerminalColor string

//ResetColor    TerminalColor = "\033[0m"

const (
	DefaultColor  TerminalColor = "DefaultColor"
	Red           TerminalColor = "Red"
	Green         TerminalColor = "Green"
	Yellow        TerminalColor = "Yellow"
	Blue          TerminalColor = "Blue"
	Magenta       TerminalColor = "Magenta"
	Cyan          TerminalColor = "Cyan"
	BrightRed     TerminalColor = "BrightRed"
	BrightGreen   TerminalColor = "BrightGreen"
	BrightYellow  TerminalColor = "BrightYellow"
	BrightBlue    TerminalColor = "BrightBlue"
	BrightMagenta TerminalColor = "BrightMagenta"
	BrightCyan    TerminalColor = "BrightCyan"
	White         TerminalColor = "White"
	Gray          TerminalColor = "Gray"
	DarkGray      TerminalColor = "DarkGray"
	Black         TerminalColor = "Black"
)

var colorCodeBases = map[TerminalColor]int{
	DefaultColor:  39,
	Red:           31,
	Green:         32,
	Yellow:        33,
	Blue:          34,
	Magenta:       35,
	Cyan:          36,
	BrightRed:     91,
	BrightGreen:   92,
	BrightYellow:  93,
	BrightBlue:    94,
	BrightMagenta: 95,
	BrightCyan:    96,
	White:         97,
	Gray:          37,
	DarkGray:      90,
	Black:         30,
}

func getFGCode(color TerminalColor) string {
	return strconv.Itoa(colorCodeBases[color])
}

func getBGCode(color TerminalColor) string {
	return strconv.Itoa(colorCodeBases[color] + 10)
}

// GetTerminalColor returns the format string for the requested foreground color.
// Resets background to default color.
// Note: some terminals may not make use of / correctly implement CSI.
func GetForeground(color TerminalColor) string {
	if !gotColors {
		return ""
	}
	return "\033[" + getFGCode(color) + ";" + getBGCode(DefaultColor) + "m"
}

// GetBackground returns the format string for the requested background color.
// Resets foreground to default color.
// Using bright colors might make the foreground color black in some terminals.
// Note: some terminals may not make use of / correctly implement CSI.
func GetBackground(color TerminalColor) string {
	if !gotColors {
		return ""
	}
	return "\033[" + getFGCode(DefaultColor) + ";" + getBGCode(color) + "m"
}

// GetColor returns the format string for the requested foreground and background color.
// Using bright background colors might make the foreground color black in some terminals.
// Note: some terminals may not make use of / correctly implement CSI.
func GetColor(foregroundColor TerminalColor, backgroundColor TerminalColor) string {
	if !gotColors {
		return ""
	}
	return "\033[" + getFGCode(foregroundColor) + ";" + getBGCode(backgroundColor) + "m"
}

// GetReset returns the format string that resets the every format setting.
// Note: some terminals may not make use of / correctly implement CSI.
func GetReset() string {
	if !gotColors {
		return ""
	}
	return "\033[0m"
}
