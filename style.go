package goprompt

import (
	"fmt"
	"strconv"
)

type attribute int

const (
	reset attribute = iota
	_
	faint
)
const (
	black attribute = iota + 30
	red
	green
	yellow
	blue
	magenta
	cyan
	white
)

var resetCode = fmt.Sprintf("%s%dm", esc, reset)

var (
	styleBlack   = styler(black)
	styleRed     = styler(red)
	styleGreen   = styler(green)
	styleYellow  = styler(yellow)
	styleBlue    = styler(blue)
	styleMagenta = styler(magenta)
	styleCyan    = styler(cyan)
	styleWhite   = styler(white)
	styleNone    = styler(reset)
	styleFaint   = styler(faint)
)

func styler(attr attribute) func(string) string {
	if attr == reset {
		return func(s string) string {
			return s
		}
	}
	attrStr := strconv.Itoa(int(attr))
	return func(s string) string {
		return fmt.Sprintf("%s%sm%v%s", esc, attrStr, s, resetCode)
	}
}
