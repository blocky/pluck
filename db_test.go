package pluck_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bky.sh/pluck"
)

func TestDB_FindFunctionByName(t *testing.T) {
	func0 := &pluck.Function{Name: "func0"}
	func1 := &pluck.Function{Name: "func1"}

	t.Run("happy path", func(t *testing.T) {
		// given
		db := &pluck.DB{
			Functions: []*pluck.Function{func0, func1},
		}

		// when
		f, err := db.FindFunctionByName("func0")
		require.NoError(t, err)

		// then
		assert.Equal(t, func0, f)
	})

	t.Run("function not found", func(t *testing.T) {
		// given
		db := &pluck.DB{
			Functions: []*pluck.Function{func0, func1},
		}

		// when
		_, err := db.FindFunctionByName("notAfunc")

		// then
		assert.ErrorIs(t, err, pluck.ErrNotFound)
	})
}

func TestDB_FindTypeByName(t *testing.T) {
	type0 := &pluck.Type{Name: "type0"}
	type1 := &pluck.Type{Name: "type1"}

	t.Run("happy path", func(t *testing.T) {
		// given
		db := &pluck.DB{
			Types: []*pluck.Type{type0, type1},
		}

		// when
		f, err := db.FindTypeByName("type0")
		require.NoError(t, err)

		// then
		assert.Equal(t, type0, f)
	})

	t.Run("type not found", func(t *testing.T) {
		// given
		db := &pluck.DB{
			Types: []*pluck.Type{type0, type1},
		}

		// when
		_, err := db.FindTypeByName("notAfunc")

		// then
		assert.ErrorIs(t, err, pluck.ErrNotFound)
	})
}
