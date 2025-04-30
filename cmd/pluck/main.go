package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"go.bky.sh/pluck"
)

func parsePick(s string) (pluck.RenderFunc, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format %q: expected 'kind:name'", s)
	}

	name := parts[1]
	if name == "" {
		return nil, fmt.Errorf("empty name in pick %q", s)
	}

	kind := strings.ToLower(parts[0])
	switch kind {
	case "function":
		return pluck.RenderFunction(name), nil
	case "type":
		return pluck.RenderType(name), nil
	default:
		return nil, fmt.Errorf("invalid kind %q: must be 'function' or 'type'", kind)
	}
}

func getSrc(inputFile string) ([]byte, error) {
	if inputFile != "" {
		return os.ReadFile(inputFile)
	} else {
		return io.ReadAll(os.Stdin)
	}
}

func getWriter(outputFile string) (io.Writer, func() error, error) {
	var writer io.Writer = os.Stdout
	closeWriter := func() error { return nil }

	if outputFile != "" {
		f, err := os.Create(outputFile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create output file: %w", err)
		}
		writer = f
		closeWriter = f.Close
	}
	return writer, closeWriter, nil
}

func runPluck(inputFile string, picks []string, outputFile string) error {
	source, err := getSrc(inputFile)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}

	// Parse picks
	var renderFuncs []pluck.RenderFunc
	for i, p := range picks {
		f, err := parsePick(p)
		if err != nil {
			return fmt.Errorf("parsing picks: %w", err)
		}

		if i != 0 {
			renderFuncs = append(renderFuncs, pluck.RenderLineBreak())
		}
		renderFuncs = append(renderFuncs, f, pluck.RenderLineBreak())
	}

	// Output
	writer, closeWriter, err := getWriter(outputFile)
	if err != nil {
		return fmt.Errorf("setting up output: %w", err)
	}
	defer closeWriter()

	db, err := pluck.Extract(source)
	if err != nil {
		return fmt.Errorf("analyzing source code: %w", err)
	}

	err = pluck.Render(writer, db, renderFuncs...)
	if err != nil {
		return fmt.Errorf("rendering: %w", err)
	}

	return nil
}

func main() {
	var inputFile string
	var outputFile string
	var picks []string

	rootCmd := &cobra.Command{
		Use:   "pluck",
		Short: "Extract functions and types from source code",
		Long: `Pluck reads source code from a file or stdin and extracts selected functions and types.

Each --pick must be in the form "kind:name", where kind is "function" or "type".

Examples:
  pluck --input myfile.go --pick function:Foo --pick type:Bar
  cat myfile.go | pluck --pick function:Foo --pick function:Baz
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPluck(inputFile, picks, outputFile)
		},
	}

	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input source file (default: stdin)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write output to file (default: stdout)")
	rootCmd.Flags().StringArrayVar(&picks, "pick", nil, "Item to extract in format kind:name (e.g., function:Foo)")

	rootCmd.MarkFlagRequired("pick")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
