package main

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"go.bky.sh/pluck"
)

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

	writer, closeWriter, err := getWriter(outputFile)
	if err != nil {
		return fmt.Errorf("setting up output: %w", err)
	}
	defer closeWriter()

	err = pluck.GenerateFromPickCmds(writer, source, picks)
	if err != nil {
		return fmt.Errorf("generating output: %w", err)
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
