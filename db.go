package pluck

import (
	"slices"
)

type Function struct {
	Name       string
	Definition string
	DocString  string
}

type Type struct {
	Name       string
	Definition string
}

type DB struct {
	Functions []*Function
	Types     []*Type
}

func (db *DB) FindFunctionByName(name string) (*Function, error) {
	idx := slices.IndexFunc(db.Functions, func(f *Function) bool {
		return f.Name == name
	})

	if idx == -1 {
		return nil, ErrNotFound
	}

	return db.Functions[idx], nil
}

func (db *DB) FindTypeByName(name string) (*Type, error) {
	idx := slices.IndexFunc(db.Types, func(t *Type) bool {
		return t.Name == name
	})

	if idx == -1 {
		return nil, ErrNotFound
	}

	return db.Types[idx], nil
}
