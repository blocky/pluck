package pluck

import (
	"fmt"
	"io"
)

type RenderFunc func(io.Writer, *DB) error

func write(w io.Writer, hunks ...string) error {
	for _, hunk := range hunks {
		if _, err := w.Write([]byte(hunk)); err != nil {
			return fmt.Errorf("writing hunk: %w", err)
		}
	}
	return nil
}

func RenderBlankLine() RenderFunc {
	return func(w io.Writer, _ *DB) error {
		_, err := w.Write([]byte("\n"))
		return err
	}
}

func RenderFunction(name string) RenderFunc {
	return func(w io.Writer, db *DB) error {
		f, err := db.FindFunctionByName(name)
		if err != nil {
			return fmt.Errorf("querying db for function '%s': %w", name, err)
		}

		err = write(w, f.DocString, "\n", f.Definition, "\n")
		if err != nil {
			return fmt.Errorf("writing function '%s': %w", name, err)
		}

		return err
	}
}

func RenderType(name string) RenderFunc {
	return func(w io.Writer, db *DB) error {
		t, err := db.FindTypeByName(name)
		if err != nil {
			return fmt.Errorf("querying db for type '%s': %w", name, err)
		}

		err = write(w, t.DocString, "\n", t.Definition, "\n")
		if err != nil {
			return fmt.Errorf("writing type '%s': %w", name, err)
		}

		return nil
	}
}

func Render(w io.Writer, db *DB, rf ...RenderFunc) error {
	for _, f := range rf {
		if err := f(w, db); err != nil {
			return fmt.Errorf("applying render func: %w", err)
		}
	}
	return nil
}
