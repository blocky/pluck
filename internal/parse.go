package internal

import (
	"fmt"
	"strings"
)

func ParsePickCmd(s string) (RenderFunc, error) {
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
		return RenderFunction(name), nil
	case "type":
		return RenderType(name), nil
	default:
		return nil, fmt.Errorf("invalid kind %q: must be 'function' or 'type'", kind)
	}
}

func PickCmdsToRenderFuncs(picks []string) ([]RenderFunc, error) {
	var renderFuncs []RenderFunc
	for i, p := range picks {
		f, err := ParsePickCmd(p)
		if err != nil {
			return nil, fmt.Errorf("parsing pick: %w", err)
		}

		if i != 0 {
			renderFuncs = append(renderFuncs, RenderLineBreak())
		}
		renderFuncs = append(renderFuncs, f, RenderLineBreak())
	}
	return renderFuncs, nil
}
