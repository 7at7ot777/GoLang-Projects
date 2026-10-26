package color

import (
	"fmt"
	"strconv"
	"strings"
)

// Color represents a color code for terminal output.
type Color struct {
	Value int
}

var (
	Black   = Color{Value: 30}
	Red     = Color{Value: 31}
	Green   = Color{Value: 32}
	Yellow  = Color{Value: 33}
	Blue    = Color{Value: 34}
	Magenta = Color{Value: 35}
	Cyan    = Color{Value: 36}
	White   = Color{Value: 37}

	Bold      = Color{Value: 1}
	Underline = Color{Value: 4}
)

func Text(text string, colors ...Color) string {
	if len(colors) == 0 {
		return text
	}
	var codes []string
	for _, color := range colors {
		codes = append(codes, strconv.Itoa(color.Value))
	}
	return fmt.Sprintf("\033[%sm%s\033[0m", strings.Join(codes, ";"), text)
}
