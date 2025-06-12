package main

import (
	"fmt"
)

type Type struct {
	BarField int
}

// type with docstring
type TypeWithDocstring struct {
	FooField string
}

type typeUnexported struct{}

type TypeWithMethods struct{}

func (m *TypeWithMethods) AMethod() {
	fmt.Println("FooMethod")
}

func Func() {
	fmt.Println("AFunc")
}

// a docstring
func FuncWithDocstring() (int, error) {
	return 0, nil
}

func main() {
	fmt.Println("Hello, world!")
}
