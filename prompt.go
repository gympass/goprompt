package goprompt

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
)

// Prompt displays a single prompt asking the user to input a given information
type Prompt struct {

	// Label represents the prompt itself. It is presented right before the
	// user's input area
	Label string

	// DefaultValue represents the initial value of the prompt.
	DefaultValue string

	// Description is shown right before the Label line. Usually this should
	// present some context to the user.
	Description string

	// Validation provides a way to ensure the user has provided a valid value
	// to the prompt. When defined, this function receives the value the user
	// is trying to use, and must return whether the value is valid. When this
	// function returns false, the prompt is replaced from "?" to an X,
	// indicating something is wrong, and the prompt is not dismissed until the
	// user input a valid value, or aborts the operation.
	Validation func(s string) bool

	value     []rune
	edited    bool
	rl        *readline.Instance
	lastState int
}

// PromptResult represents the result of a Prompt operation.
type PromptResult struct {
	// Value represents the value input by the user
	Value string

	// Cancelled indicates whether the operation was cancelled by the user;
	// usually when they send an interruption signal to the process.
	Cancelled bool
}

// Run shows the prompt on screen. Returns either a PromptResult, or an error,
// in case the displaying process fails.
func (p *Prompt) Run() (*PromptResult, error) {
	c := &readline.Config{
		Stdin:          os.Stdin,
		Stdout:         os.Stdout,
		EnableMask:     false,
		HistoryLimit:   -1,
		VimMode:        false,
		UniqueEditLine: true,
	}
	if err := c.Init(); err != nil {
		return nil, err
	}
	c.SetListener(p.listener)
	if ex, err := readline.NewEx(c); err != nil {
		return nil, err
	} else {
		defer ex.Close()
		p.rl = ex
	}
	if p.Description != "" {
		_, _ = fmt.Fprintln(p.rl, p.Description)
	}
	p.value = []rune(p.DefaultValue)
	var (
		v   string
		err error
	)
	p.prepare()
	for {
		v, err = p.rl.ReadlineWithDefault(string(p.value))
		if err == readline.ErrInterrupt {
			return &PromptResult{
				Cancelled: true,
			}, nil
		} else if err == io.EOF {
			return &PromptResult{
				Value:     string(p.value),
				Cancelled: len(p.value) > 0,
			}, nil
		} else if err != nil {
			return nil, err
		}

		if p.Validation != nil && !p.Validation(v) {
			continue
		}
		break
	}

	_, _ = fmt.Fprintln(p.rl, v)
	return &PromptResult{
		Value: v,
	}, nil
}

func shouldIgnore(r rune) bool {
	return r == readline.CharEsc || r == readline.MetaBackspace
}

func (p *Prompt) listener(line []rune, newPos int, key rune) ([]rune, int, bool) {
	if shouldIgnore(key) {
		return line, newPos, false
	}

	if key == readline.CharEnter {
		if p.Validation != nil {
			if !p.Validation(string(p.value)) {
				p.setStatus(statusError)
				return p.value, newPos, false
			}
		}
		p.setStatus(statusOK)
	}

	if line != nil {
		p.value = line
	}
	return line, newPos, true
}

const (
	statusPrompt = iota
	statusOK
	statusError
)

func (p *Prompt) setStatus(status int) {
	p.lastState = status
	p.prepare()
}

func (p *Prompt) prepare() {
	style := styleNone
	prefix := glyphPrompt
	switch p.lastState {
	case statusOK:
		prefix = glyphCheck
		style = styleGreen
	case statusError:
		prefix = glyphError
		style = styleMagenta
	}
	p.rl.SetPrompt(fmt.Sprintf("%s %s: ", style(prefix), p.Label))
}
