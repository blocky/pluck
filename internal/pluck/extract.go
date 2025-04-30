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

type treeSupport struct {
	query  *sitter.Query
	match  *sitter.QueryMatch
	src    []byte
	errors error
}

func (s *treeSupport) getContent(id string) string {
	node := findNodeByName(s.query, s.match, id)
	if node == nil {
		s.errors = errors.Join(s.errors, fmt.Errorf("could not find %s", id))
		return ""
	}

	return node.Content(s.src)
}

func (s *treeSupport) collectComments(id string) string {
	start := findNodeByName(s.query, s.match, id)
	if start == nil {
		s.errors = errors.Join(s.errors, fmt.Errorf("could not find %s", id))
		return ""
	}

	fnDocStringLines := []string{}
	for cur := start.PrevSibling(); cur != nil; cur = cur.PrevSibling() {
		if cur.Type() != "comment" {
			break
		}

		fnDocStringLines = append([]string{cur.Content(s.src)}, fnDocStringLines...)
	}

	return strings.Join(fnDocStringLines, "\n")
}

func (s *treeSupport) extractTypeName(id string) string {
	start := findNodeByName(s.query, s.match, id)
	if start == nil {
		s.errors = errors.Join(s.errors, fmt.Errorf("could not find %s", id))
		return ""
	}

	switch {
	case start == nil:
		s.errors = errors.Join(errors.New("extracting type received nil node"))
		return ""
	case start.Type() == "type_identifier":
		return start.Content(s.src)
	case start.Type() != "pointer_type":
		err := fmt.Errorf("extracting type received unexpected type %s", start.Type())
		s.errors = errors.Join(s.errors, err)
		return ""
	}

	child := start.Child(1)
	if child == nil {
		err := fmt.Errorf("extracting type received unexpected missing child")
		s.errors = errors.Join(s.errors, err)
		return ""
	}

	return child.Content(s.src)
}

func (s *treeSupport) Close() error {
	return s.errors
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

		support := treeSupport{query: query, match: match, src: src}
		types = append(types, &Type{
			Name:       support.getContent("typeName"),
			Definition: support.getContent("typeDefinition"),
			DocString:  support.collectComments("typeDefinition"),
		})
		if err := support.Close(); err != nil {
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

		support := treeSupport{query: query, match: match, src: src}
		f := &Function{
			Name:       support.getContent("fnName"),
			Definition: support.getContent("fnDefinition"),
			DocString:  support.collectComments("fnDefinition"),
		}
		if err := support.Close(); err != nil {
			return nil, fmt.Errorf("interpreting query: %w", err)
		}

		functions = append(functions, f)
	}

	return functions, nil
}

func extractMethods(src []byte) ([]*Function, error) {
	pattern := []byte(`
		(method_declaration
			receiver: (parameter_list
				(parameter_declaration
					name: (identifier)
					type: (_) @receiverType
				)
			)
			name: (field_identifier) @fnName
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

		support := treeSupport{query: query, match: match, src: src}
		fnName := support.getContent("fnName")
		typeName := support.extractTypeName("receiverType")
		f := &Function{
			Name:       typeName + "." + fnName,
			Definition: support.getContent("fnDefinition"),
			DocString:  support.collectComments("fnDefinition"),
		}
		if err := support.Close(); err != nil {
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

	methods, err := extractMethods(src)
	if err != nil {
		return nil, fmt.Errorf("extracting methods: %w", err)
	}

	types, err := extractType(src)
	if err != nil {
		return nil, fmt.Errorf("extracting types: %w", err)
	}

	return &DB{
		Functions: append(funcs, methods...),
		Types:     types,
	}, nil
}
