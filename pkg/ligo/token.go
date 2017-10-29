package ligo

import (
	"strings"
)

// StripComments function is used to strip the comments from the passed ligo source
func StripComments(ltxt string) string {
	lines := strings.Split(ltxt, "\n")

	inQuotes := false
	final := ""
	for _, line := range lines {
		for _, ch := range line {
			if ch == '"' {
				inQuotes = !inQuotes
				final += string(ch)
				continue
			}
			if ch == ';' && !inQuotes {
				break
			}
			final += string(ch)
		}
		final += "\n"
	}
	return final
}

// isValid function is used to check whether a string is a valid lisp expression
func isValid(ltxt string) error {
	if len(ltxt) < 2 {
		return Error("Expected atleast (), got : " + ltxt)
	}
	if ltxt[0] != '(' {
		return Error("Expected '(' at the start of the expression, got : " + string(ltxt[0]) + "\n" + ltxt)
	}
	if ltxt[len(ltxt)-1] != ')' {
		return Error("Expected ')' at the end of the expression, got : " + string(ltxt[0]))
	}

	return nil
}

// ScanTokens is used to get token list from a passed ligo expression.
// This is one of the important functions for parsing a ligo source code.
func ScanTokens(ltxt string) ([]string, error) {
	ltxt = StripComments(ltxt)
	ltxt = strings.TrimSpace(ltxt)

	err := isValid(ltxt)
	if err != nil {
		return nil, err
	}

	p := newParser(ltxt)

	for p.i = 1; p.i < len(p.ltxt); p.i++ {
		c := string(p.ltxt[p.i])
		err := p.runHandler(c)
		if err != nil {
			return nil, err
		}
	}
	if p.inQuotes {
		return nil, Error("Quote not closed correctly")
	}
	if p.inSBkts {
		return nil, Error("Closure not closed correctly")
	}
	return p.strList, nil
}

// MatchChars function is used to return the offset at which the matching character of the passed character
// is found in the passed string. Generally used to match brackets.
func MatchChars(ltxt string, off int64, open byte, close byte) int64 {
	if int64(len(ltxt)) <= off {
		return -1
	}
	if ltxt[off] != open {
		return -1
	}
	count := 1
	inQuotes := false
	for i := off + 1; i < int64(len(ltxt)); i++ {
		if ltxt[i] == '"' {
			inQuotes = !inQuotes
		}
		if ltxt[i] == open && !inQuotes {
			count++
		}
		if ltxt[i] == close && !inQuotes {
			count--
		}
		if count == 0 {
			return i
		}
	}
	return -1
}

// getVarsFromClosure is used to extract all the parameter names from a
// closure of a function definition in ligo
// (ie., "|a b v r|" yields an array containing "a", "b", "v" and "r")
func getVarsFromClosure(cl string) []string {
	current := ""
	retParams := make([]string, 0)
	for _, val := range cl {
		if strings.TrimSpace(string(val)) == "" || val == '|' {
			if current == "" {
				continue
			}
			retParams = append(retParams, current)
			current = ""
			continue
		}
		current += string(val)
	}
	return retParams
}

// isVariate is used to check whether a given token string is passed as a variate parameter.
func isVariate(str string) bool {
	if len(str) > 4 && str[:3] == "..." && str[3] != '.' {
		return true
	}
	return false
}
