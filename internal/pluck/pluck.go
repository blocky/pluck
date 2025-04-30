package pluck

import (
	"fmt"
	"io"
)

func GenerateFromPickCmds(w io.Writer, source []byte, picks []string) error {
	renderFuncs, err := PickCmdsToRenderFuncs(picks)
	if err != nil {
		return fmt.Errorf("parsing picks: %w", err)
	}

	db, err := Extract(source)
	if err != nil {
		return fmt.Errorf("analyzing source code: %w", err)
	}

	err = Render(w, db, renderFuncs...)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	return nil
}
