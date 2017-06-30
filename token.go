package ligo

import (
	"fmt"
	"strings"
)

func ScanTokens(ltxt string) ([]string, error) {
	strList := make([]string, 0)
	ltxt = strings.TrimSpace(ltxt)
	if len(ltxt) < 2 {
		return nil, LigoError("Expected atleast (), got : " + ltxt)
	}

	if ltxt[0] != '(' {
		return nil, LigoError("Expected '(' at the start of the expression, got : " + string(ltxt[0]) + "\n" +
			ltxt,
		)
	}
	if ltxt[len(ltxt)-1] != ')' {
		return nil, LigoError("Expected ')' at the start of the expression, got : " + string(ltxt[0]))
	}
	inQuotes := false
	inSBkts := false
	current := ""
	for i := 1; i < len(ltxt); i++ {
		c := string(ltxt[i])
		switch c {
		case " ", "\n", "\r", "\t":
			if !inQuotes && !inSBkts {
				if current == "" {
					continue
				}
				strList = append(strList, current)
				current = ""
				continue
			}
			current += c
		case "|":
			if inSBkts {
				current += c
				inSBkts = false
				continue
			}
			if !inQuotes {
				if current != "" {
					return nil, LigoError("Closure not seperated by a space")
				}
				current += c
				inSBkts = true
				continue
			}
			current += c
		case "\"":
			if !inQuotes {
				inQuotes = true
				if current != "" {
					return nil, LigoError("Not seperated by a space")
				}
				current += "\""
				continue
			}
			current += "\""
			strList = append(strList, current)
			inQuotes = false
			current = ""
		case "[":
			if inSBkts {
				return nil, LigoError("'[' not expected inside a closure")
			}
			if inQuotes {
				current += c
				continue
			}
			if current != "" {
				return nil, LigoError("Array not seperated by a space")
			}
			off := MatchChars(ltxt, int64(i), '[', ']') + 1
			current = ltxt[i:off]
			strList = append(strList, current)
			i = int(off)
			if strings.TrimSpace(string(ltxt[i])) != "" && strings.TrimSpace(string(ltxt[i])) != ")" {
				return nil, LigoError("Unexpected character found at array end : " + string(ltxt[i]))
			}
			current = ""
		case "(":
			if inSBkts {
				return nil, LigoError("'(' not expected inside a closure")
			}
			if inQuotes {
				current += c
				continue
			}
			if current != "" {
				return nil, LigoError("Expression not seperated by a space")
			}
			off := MatchChars(ltxt, int64(i), '(', ')') + 1
			current = ltxt[i:off]
			strList = append(strList, current)
			i = int(off)
			if strings.TrimSpace(string(ltxt[i])) != "" && strings.TrimSpace(string(ltxt[i])) != ")" {
				return nil, LigoError("Unexpected character found at expression end : " + string(ltxt[i]))
			}
			current = ""
		case ")":
			if inSBkts {
				return nil, LigoError("')' not expected inside a closure")
			}
			if inQuotes {
				current += c
				continue
			}
			if len(ltxt)-1 != i {
				return nil, LigoError("Expected EOL, got " + string(ltxt[i]) + " at " + fmt.Sprint(i))
			}
			if current != "" {
				strList = append(strList, current)
				current = ""
			}
		case "]":
			if inSBkts {
				return nil, LigoError("']' not expected inside a closure")
			}
			if inQuotes {
				current += c
				continue
			}
			strList = append(strList, current)
			current = ""
		default:
			current += c
		}
	}
	if inQuotes {
		return nil, LigoError("Quote not closed correctly")
	}
	if inSBkts {
		return nil, LigoError("Closure not closed correctly")
	}
	return strList, nil
}

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
			count += 1
		}
		if ltxt[i] == close && !inQuotes {
			count -= 1
		}
		if count == 0 {
			return i
		}
	}
	return -1
}

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
