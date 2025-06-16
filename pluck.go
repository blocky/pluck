package pluck

import (
	"io"

	pluckInternal "github.com/blocky/pluck/internal/pluck"
)

func GenerateFromPickCmds(w io.Writer, source []byte, picks []string) error {
	return pluckInternal.GenerateFromPickCmds(w, source, picks)
}
