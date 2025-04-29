package pluck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bky.sh/pluck"
)

const goSrc = `
package main

import "fmt"

// line a
// line b
type aType struct {
	A int
}

type anotherType struct{ oneLine int }


// line 1
// line 2
func notMain(a, b int) error {
	if a + b < 10 {
		return errors.New("too small")
	}

	f := func() {
		fmt.Println("hi")
	}

	return nil
}

func oneLine() { fmt.Println("hello"); fmt.Println("world") }

func main() {
	fmt.Println("hello")
	fmt.Println("world")
}
`

func TestExtract(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		input := []byte(goSrc)

		// when
		got, err := pluck.Extract(input)
		require.NoError(t, err)

		// then
		assert.Len(t, got.Functions, 3)
		assert.Equal(t, "notMain", got.Functions[0].Name)
		assert.Equal(t, "// line 1\n// line 2", got.Functions[0].DocString)
		assert.Equal(t, "oneLine", got.Functions[1].Name)
		assert.Equal(t, "main", got.Functions[2].Name)

		assert.Len(t, got.Types, 2)
		assert.Equal(t, "aType", got.Types[0].Name)
		assert.Equal(t, "// line a\n// line b", got.Types[0].DocString)
		assert.Equal(t, "anotherType", got.Types[1].Name)
	})

	t.Run("invalid source", func(t *testing.T) {
		// given
		input := []byte("this is not go code")

		// when
		got, err := pluck.Extract(input)
		require.NoError(t, err)

		// then
		assert.Len(t, got.Functions, 0)
		assert.Len(t, got.Types, 0)
	})
}
