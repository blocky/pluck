package internal

import (
	"fmt"
	"io"
	"strings"
)

type RenderFunc func(io.Writer, *DB) error

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func write(w io.Writer, hunks ...string) error {
	for _, hunk := range hunks {
		if _, err := w.Write([]byte(hunk)); err != nil {
			return fmt.Errorf("writing hunk: %w", err)
		}
	}
	return nil
}

func RenderLineBreak() RenderFunc {
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

		docString := ""
		if !isEmpty(f.DocString) {
			docString = f.DocString + "\n"
		}
		err = write(w, docString, f.Definition)
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

		docString := ""
		if !isEmpty(t.DocString) {
			docString = t.DocString + "\n"
		}
		err = write(w, docString, t.Definition)
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
