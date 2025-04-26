package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"go.bky.sh/pluck"
)

type config struct {
	name string
	src  []byte
}

func getSrcFromFile(fileName string) ([]byte, error) {
	return os.ReadFile(fileName)
}

func getSrcFromStdin() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

func getSrc(args []string) ([]byte, error) {
	if len(args) > 1 {
		return getSrcFromFile(os.Args[1])
	} else {
		return getSrcFromStdin()
	}
}

func parse(args []string) (*config, error) {
	switch l := len(args); {
	case l <= 1:
		return nil, errors.New("invalid args")
	case l == 2:

		fmt.Println("getting from stdin")
		src, err := getSrcFromStdin()
		if err != nil {
			return nil, fmt.Errorf("getting source from stdin: %w", err)
		}

		return &config{name: args[1], src: src}, nil
	case l > 2:
		fileName := args[2]
		src, err := getSrcFromFile(fileName)
		if err != nil {
			return nil, fmt.Errorf("getting source from a file: %w", err)
		}

		return &config{name: args[1], src: src}, nil
	default:
		return nil, errors.New("input not handles")
	}
}

func exitOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("terminating: %s: %s", msg, err.Error())
	}
}

func main() {
	config, err := parse(os.Args)
	exitOnError("parsing args", err)

	functions, err := pluck.Extract(config.src)
	exitOnError("extracting functions", err)

	f, err := pluck.Filter(functions, config.name)
	exitOnError("filtering for target function", err)

	err = pluck.Render(os.Stdout, f)
	exitOnError("rendering", err)
}
