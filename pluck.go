package pluck

import (
	"fmt"
	"io"

	"github.com/blocky/pluck/internal"
)

func GenerateFromPickCmds(w io.Writer, source []byte, picks []string) error {
	renderFuncs, err := internal.PickCmdsToRenderFuncs(picks)
	if err != nil {
		return fmt.Errorf("parsing picks: %w", err)
	}

	db, err := internal.Extract(source)
	if err != nil {
		return fmt.Errorf("analyzing source code: %w", err)
	}

	err = internal.Render(w, db, renderFuncs...)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	return nil
}
