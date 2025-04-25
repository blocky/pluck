package pluck_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.bky.sh/pluck"
)

func TestPluckGo(t *testing.T) {
	t.Run("init testing", func(t *testing.T) {
		r := strings.NewReader(`
package main

import "fmt"

func notMain() error {
	return errors.New("not main")
}

func main() {
	fmt.Println("hello world")
}
`)
		_, err := pluck.FromGoSrc(r, "main")
		assert.Error(t, err)
	})

}
