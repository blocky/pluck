package pluck_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bky.sh/pluck"
)

func TestRenderLineBreak(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := pluck.RenderLineBreak()

		// when
		err := f(buf, nil)
		require.NoError(t, err)

		// then
		assert.Equal(t, "\n", buf.String())
	})
}

func TestRenderFunction(t *testing.T) {
	db := &pluck.DB{
		Functions: []*pluck.Function{
			{Name: "foo", Definition: "fooDef", DocString: "fooDoc"},
		},
	}

	t.Run("happy path", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := pluck.RenderFunction("foo")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "fooDoc\nfooDef\n", buf.String())
	})

	t.Run("not in db", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := pluck.RenderFunction("bar")

		// when
		err := f(buf, db)

		// then
		assert.ErrorIs(t, err, pluck.ErrNotFound)
	})
}

func TestRenderType(t *testing.T) {
	db := &pluck.DB{
		Types: []*pluck.Type{
			{Name: "foo", Definition: "fooDef", DocString: "fooDoc"},
		},
	}

	t.Run("happy path", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := pluck.RenderType("foo")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "fooDoc\nfooDef\n", buf.String())
	})

	t.Run("not in db", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := pluck.RenderFunction("bar")

		// when
		err := f(buf, db)

		// then
		assert.ErrorIs(t, err, pluck.ErrNotFound)
	})
}
