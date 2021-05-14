package goprompt

import (
	"os"

	"github.com/chzyer/readline"
)

// This is a little hack to prevent bells from being emitted to the terminal.
// We do it by checking whether there's a single bell character being written,
// and strip it. That's it.

func suppressBell(f *os.File) *bellSuppressor {
	return &bellSuppressor{f: f}
}

type bellSuppressor struct {
	f *os.File
}

func (s *bellSuppressor) Write(b []byte) (int, error) {

	if len(b) == 1 && b[0] == readline.CharBell {
		return 0, nil
	}
	return s.f.Write(b)
}

func (s *bellSuppressor) Close() error {
	return s.f.Close()
}
