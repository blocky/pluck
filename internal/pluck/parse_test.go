package pluck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.bky.sh/pluck"
)

func TestParsePick(t *testing.T) {
	t.Run("invalid format", func(t *testing.T) {
		_, err := pluck.ParsePick("invalid-format")
		assert.ErrorContains(t, err, "invalid format")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := pluck.ParsePick("function:")
		assert.ErrorContains(t, err, "empty name")
	})

	t.Run("invalid kind", func(t *testing.T) {
		_, err := pluck.ParsePick("foo:bar")
		assert.ErrorContains(t, err, "invalid kind")
	})
}
