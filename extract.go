package pluck

import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

func newQuery(pattern string, src []byte) (*sitter.QueryCursor, error) {
	lang := golang.GetLanguage()
	n, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, fmt.Errorf("creating parse context: %w", err)
	}

	q, err := sitter.NewQuery([]byte(pattern), lang)
	if err != nil {
		return nil, fmt.Errorf("creating query: %w", err)
	}
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	return qc, nil
}

func Extract(src []byte) ([]*Function, error) {
	eventDefPattern := `
		(function_declaration
			name: (identifier) @fnName
			parameters: (parameter_list)
			body: ( _ ) @fnBody
		) @fnDefinition
	`

	qc, err := newQuery(eventDefPattern, src)
	if err != nil {
		return nil, fmt.Errorf("creating query: %w", err)
	}

	var functions []*Function
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		m = qc.FilterPredicates(m, src)
		if len(m.Captures) != 3 {
			return nil, fmt.Errorf("unexpected number of captures: %w", err)
		}

		f := &Function{
			Definition: m.Captures[0].Node.Content(src),
			Name:       m.Captures[1].Node.Content(src),
			Body:       m.Captures[2].Node.Content(src),
		}

		functions = append(functions, f)
	}

	return functions, nil
}
