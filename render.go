package pluck

import (
	"io"
	"strings"
	"unicode"
)

func removeFirst(code string, toRemove string) string {
	return strings.Replace(code, string(toRemove), "", 1)
}

func removeLast(code string, toRemove string) string {
	idx := strings.LastIndex(code, string(toRemove))
	// -1 means that the string was not found
	if idx == -1 {
		return code
	}
	return code[:idx] + code[idx+len(toRemove):]
}

func breakIntoLines(code string) []string {
	return strings.Split(code, "\n")
}

func trimLeadingBlankLines(lines []string) []string {
	firstNonBlankLine := 0
	for ; firstNonBlankLine < len(lines); firstNonBlankLine++ {
		line := lines[firstNonBlankLine]
		if strings.TrimSpace(line) != "" {
			break
		}
	}

	return lines[firstNonBlankLine:]
}

func trimTrailingBlankLines(lines []string) []string {
	lastNonBlankLine := len(lines) - 1
	for ; 0 <= lastNonBlankLine; lastNonBlankLine-- {
		line := lines[lastNonBlankLine]
		if strings.TrimSpace(line) != "" {
			break
		}
	}

	return lines[:lastNonBlankLine+1]
}

func getLinePefix(lines []string) string {
	if len(lines) == 0 {
		return ""
	}

	s := lines[0]
	var b strings.Builder
	for _, r := range s {
		if !unicode.IsSpace(r) {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}

func trimPrefix(lines []string, prefix string) []string {
	trimmedLines := make([]string, len(lines))
	for i, line := range lines {
		trimmedLines[i] = removeFirst(line, prefix)
	}
	return trimmedLines
}

func splitToLines(code string) []string {
	return strings.Split(code, "\n")
}

func joinToCode(lines []string) string {
	return strings.Join(lines, "\n") + "\n"
}

func Render(w io.Writer, fn *Function) error {
	code := removeFirst(fn.Body, "{")
	code = removeLast(code, "}")
	lines := splitToLines(code)
	lines = trimLeadingBlankLines(lines)
	lines = trimTrailingBlankLines(lines)

	prefix := getLinePefix(lines)
	lines = trimPrefix(lines, prefix)

	rendered := joinToCode(lines)

	w.Write([]byte(rendered))
	return nil
}
