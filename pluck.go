package pluck

import (
	"io"

	"github.com/blocky/pluck/internal"
)

func GenerateFromPickCmds(w io.Writer, source []byte, picks []string) error {
	return internal.GenerateFromPickCmds(w, source, picks)
}
