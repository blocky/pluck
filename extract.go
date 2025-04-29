package pluck

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

// findNodeByName searches the capture for the specified name.
// similar to other sitter tools, returns nil if not found.
func findNodeByName(
	query *sitter.Query,
	match *sitter.QueryMatch,
	s string,
) *sitter.Node {
	for _, capture := range match.Captures {
		captureName := query.CaptureNameForId(capture.Index)
		if captureName == s {
			return capture.Node
		}
	}

	return nil
}

// collectComments stating from the current node,
// walk back though siblings collecting all comments
func collectComments(start *sitter.Node, src []byte) string {
	if start == nil {
		return ""
	}

	fnDocStringLines := []string{}
	for cur := start.PrevSibling(); cur != nil; cur = cur.PrevSibling() {
		if cur.Type() != "comment" {
			break
		}

		fnDocStringLines = append([]string{cur.Content(src)}, fnDocStringLines...)
	}

	return strings.Join(fnDocStringLines, "\n")
}

type matchInterpreter struct {
	query  *sitter.Query
	match  *sitter.QueryMatch
	src    []byte
	errors error
}

func (m *matchInterpreter) getContentAndWitholdError(s string) string {
	node := findNodeByName(m.query, m.match, s)
	if node == nil {
		m.errors = errors.Join(m.errors, fmt.Errorf("could not find %s", s))
		return ""
	}

	return node.Content(m.src)
}

func (m *matchInterpreter) ErrorIfAny() error {
	return m.errors
}

func initQuery(src, pattern []byte) (*sitter.QueryCursor, *sitter.Query, error) {
	lang := golang.GetLanguage()
	root, err := sitter.ParseCtx(context.Background(), src, lang)
	if err != nil {
		return nil, nil, fmt.Errorf("creating parse context: %w", err)
	}

	query, err := sitter.NewQuery(pattern, lang)
	if err != nil {
		return nil, nil, fmt.Errorf("creating query: %w", err)
	}
	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	return cursor, query, nil
}

func extractType(src []byte) ([]*Type, error) {
	pattern := []byte(`
		(type_declaration
			(type_spec
				name: (type_identifier) @typeName
			)
		) @typeDefinition
	`)

	cursor, query, err := initQuery(src, pattern)
	if err != nil {
		return nil, fmt.Errorf("initializing query: %w", err)
	}

	var types []*Type
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		typeDefNode := findNodeByName(query, match, "typeDefinition")
		interpreter := matchInterpreter{query: query, match: match, src: src}
		types = append(types, &Type{
			Name:       interpreter.getContentAndWitholdError("typeName"),
			Definition: interpreter.getContentAndWitholdError("typeDefinition"),
			DocString:  collectComments(typeDefNode, src),
		})
		if err := interpreter.ErrorIfAny(); err != nil {
			return nil, fmt.Errorf("interpreting query: %w", err)
		}

	}
	return types, nil

}

func extractFunc(src []byte) ([]*Function, error) {
	pattern := []byte(`
		(function_declaration
			name: (identifier) @fnName
		) @fnDefinition
	`)

	cursor, query, err := initQuery(src, pattern)
	if err != nil {
		return nil, fmt.Errorf("initializing query: %w", err)
	}

	var functions []*Function
	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		fnDefNode := findNodeByName(query, match, "fnDefinition")
		interpreter := matchInterpreter{query: query, match: match, src: src}
		f := &Function{
			Name:       interpreter.getContentAndWitholdError("fnName"),
			Definition: interpreter.getContentAndWitholdError("fnDefinition"),
			DocString:  collectComments(fnDefNode, src),
		}
		if err := interpreter.ErrorIfAny(); err != nil {
			return nil, fmt.Errorf("interpreting query: %w", err)
		}

		functions = append(functions, f)
	}

	return functions, nil
}

func Extract(src []byte) (*DB, error) {
	funcs, err := extractFunc(src)
	if err != nil {
		return nil, fmt.Errorf("extracting funcs: %w", err)
	}

	types, err := extractType(src)
	if err != nil {
		return nil, fmt.Errorf("extracting types: %w", err)
	}

	return &DB{
		Functions: funcs,
		Types:     types,
	}, nil
}
