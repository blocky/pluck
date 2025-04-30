package pluck_test

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.bky.sh/pluck/internal/pluck"
)

var update = flag.Bool("update", false, "update golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestGenerateFromPickCmds(t *testing.T) {
	testdataDir := "./testdata"
	p := filepath.Join(testdataDir, "src.go")
	source, err := os.ReadFile(p)
	require.NoError(t, err)

	goldendataDir := filepath.Join(testdataDir, "golden")

	testCases := []struct {
		id     string
		golden string
		cmds   []string
	}{
		{id: "000", cmds: []string{"function:main"}},
		{id: "001", cmds: []string{"type:Type"}},
		{id: "002", cmds: []string{"type:TypeWithDocstring"}},
		{id: "003", cmds: []string{"type:typeUnexported"}},
		{id: "004", cmds: []string{"function:TypeWithMethods.AMethod"}},
		{id: "005", cmds: []string{"function:Func"}},
		{id: "006", cmds: []string{"function:FuncWithDocstring"}},
		{id: "007", cmds: []string{"type:TypeWithMethods", "function:TypeWithMethods.AMethod"}},
		{id: "008", cmds: []string{"type:Type", "function:main"}},
		{id: "009", cmds: []string{"function:main", "type:typeUnexported"}},
		{
			id: "010",
			cmds: []string{
				"function:main",
				"function:FuncWithDocstring",
				"function:TypeWithMethods.AMethod",
				"type:Type",
				"type:TypeWithDocstring",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.id, func(t *testing.T) {
			var got bytes.Buffer
			err = pluck.GenerateFromPickCmds(&got, source, tc.cmds)
			require.NoError(t, err)

			f := fmt.Sprintf("go-%s", tc.id)
			goldenFile := filepath.Join(goldendataDir, f)

			if *update {
				err := os.WriteFile(goldenFile, got.Bytes(), 0644)
				require.NoError(t, err)
			}

			want, err := os.ReadFile(goldenFile)
			require.NoError(t, err)

			assert.Equal(t, string(want), got.String())
		})
	}
}
