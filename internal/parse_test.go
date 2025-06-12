package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/blocky/pluck/internal"
)

func TestParsePick(t *testing.T) {
	t.Run("invalid format", func(t *testing.T) {
		_, err := internal.ParsePickCmd("invalid-format")
		assert.ErrorContains(t, err, "invalid format")
	})

	t.Run("empty name", func(t *testing.T) {
		_, err := internal.ParsePickCmd("function:")
		assert.ErrorContains(t, err, "empty name")
	})

	t.Run("invalid kind", func(t *testing.T) {
		_, err := internal.ParsePickCmd("foo:bar")
		assert.ErrorContains(t, err, "invalid kind")
	})
}
