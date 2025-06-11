package internal_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blocky/pluck/internal"
)

func TestRenderLineBreak(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderLineBreak()

		// when
		err := f(buf, nil)
		require.NoError(t, err)

		// then
		assert.Equal(t, "\n", buf.String())
	})
}

func TestRenderFunction(t *testing.T) {
	db := &internal.DB{
		Functions: []*internal.Function{
			{Name: "foo", Definition: "fooDef", DocString: "fooDoc"},
			{Name: "noDocs", Definition: "noDocsDef", DocString: ""},
		},
	}

	t.Run("happy path - with docstring", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderFunction("foo")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "fooDoc\nfooDef", buf.String())
	})

	t.Run("happy path - no docstring", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderFunction("noDocs")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "noDocsDef", buf.String())
	})

	t.Run("not in db", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderFunction("bar")

		// when
		err := f(buf, db)

		// then
		assert.ErrorIs(t, err, internal.ErrNotFound)
	})
}

func TestRenderType(t *testing.T) {
	db := &internal.DB{
		Types: []*internal.Type{
			{Name: "foo", Definition: "fooDef", DocString: "fooDoc"},
			{Name: "noDocs", Definition: "noDocsDef", DocString: ""},
		},
	}

	t.Run("happy path - with docstring", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderType("foo")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "fooDoc\nfooDef", buf.String())
	})

	t.Run("happy path - no docstring", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderType("noDocs")

		// when
		err := f(buf, db)
		require.NoError(t, err)

		// then
		assert.Equal(t, "noDocsDef", buf.String())
	})

	t.Run("not in db", func(t *testing.T) {
		// given
		buf := &bytes.Buffer{}
		f := internal.RenderFunction("bar")

		// when
		err := f(buf, db)

		// then
		assert.ErrorIs(t, err, internal.ErrNotFound)
	})
}
