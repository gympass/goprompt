package goprompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

const esc = "\033["

const (
	moveUp     = esc + "1A"
	hideCursor = esc + "?25l"
	showCursor = esc + "?25h"
	clearLine  = esc + "2K"
)

// Select displays a list composed of provided items, and allows users to select
// a single item using arrow keys.
type Select struct {
	// Description to be shown before the list itself.
	Description string
	// Options represent items to be presented to the user.
	Options []string
	// Label is displayed along the instructions of using arrows to select,
	// enter to confirm. Usually this would be a prompt like "Select an item".
	Label        string
	rl           *readline.Instance
	optLen       int
	currentIndex int
	rendered     bool
	done         bool
}

// SelectResult represents the result of a Select operation.
type SelectResult struct {
	// SelectedIndex indicates which index of provided Select.Options was
	// selected by the user.
	SelectedIndex int

	// SelectedValue indicates which value of provided Select.Options was
	// selected by the user.
	SelectedValue string

	// Cancelled indicates whether the operation was cancelled by the user;
	// usually when they send an interruption signal to the process.
	Cancelled bool
}

// Run shows the select on screen. Returns either a SelectResult, or an error,
// in case the displaying process fails.
func (s *Select) Run() (*SelectResult, error) {
	c := &readline.Config{
		Stdin:          os.Stdin,
		Stdout:         suppressBell(os.Stdout),
		Stderr:         suppressBell(os.Stderr),
		EnableMask:     false,
		HistoryLimit:   -1,
		VimMode:        false,
		UniqueEditLine: true,
	}

	s.optLen = len(s.Options)
	if err := c.Init(); err != nil {
		return nil, err
	}
	c.SetListener(s.listener)
	if ex, err := readline.NewEx(c); err != nil {
		return nil, err
	} else {
		s.rl = ex
	}

	if s.Description != "" {
		_, _ = fmt.Fprintln(s.rl, s.Description)
	}

	_, _ = fmt.Fprintf(s.rl, "%s [Use %s to select, enter to confirm]\n", s.Label, glyphUpDown)

	_, _ = fmt.Fprint(s.rl, hideCursor)
	defer func() {
		fmt.Print(showCursor)
	}()

	for !s.done {
		_, err := s.rl.Readline()
		if err != nil {
			return nil, err
		}
	}

	s.cleanAndCommit()

	return &SelectResult{
		SelectedIndex: s.currentIndex,
		SelectedValue: s.Options[s.currentIndex],
		Cancelled:     false,
	}, nil
}

func (s *Select) cleanAndCommit() {
	for i := 0; i < s.optLen+1; i++ {
		_, _ = fmt.Fprint(s.rl, moveUp)
		_, _ = fmt.Fprint(s.rl, clearLine)
	}
	colon := ""
	if !strings.HasSuffix(s.Label, "?") {
		colon = ":"
	}
	_, _ = fmt.Fprint(s.rl, fmt.Sprintf("%s %s%s %s\n", styleGreen(glyphCheck), s.Label, colon, s.Options[s.currentIndex]))

}

func (s *Select) render() {
	if s.rendered {
		_, _ = fmt.Fprint(s.rl, strings.Repeat(moveUp, s.optLen))
	}

	for i, v := range s.Options {
		if i == s.currentIndex {
			_, _ = fmt.Fprintf(s.rl, "  %s %s\n", styleCyan(glyphSelection), styleCyan(v))
		} else {
			_, _ = fmt.Fprintln(s.rl, "    "+v)
		}
	}
	s.rendered = true
}

func (s *Select) listener(line []rune, newPos int, key rune) ([]rune, int, bool) {
	switch key {
	case readline.CharNext: // down
		if s.currentIndex+1 >= s.optLen {
			s.currentIndex = 0
		} else {
			s.currentIndex++
		}

	case readline.CharPrev: // up
		if s.currentIndex-1 == -1 {
			s.currentIndex = s.optLen - 1
		} else {
			s.currentIndex--
		}

	case readline.CharEnter: // confirm
		s.done = true
		return line, newPos, false
	}
	s.render()
	return line, newPos, false
}
