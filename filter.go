package pluck

import (
	"errors"
	"slices"
)

func Filter(functions []*Function, fnName string) (*Function, error) {
	idx := slices.IndexFunc(functions, func(f *Function) bool {
		return f.Name == fnName
	})

	if idx == -1 {
		return nil, errors.New("not found")
	}

	return functions[idx], nil
}
