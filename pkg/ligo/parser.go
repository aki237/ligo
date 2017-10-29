package ligo

import (
	"fmt"
	"strings"
)

type parser struct {
	strList  []string
	inQuotes bool
	inSBkts  bool
	current  string
	ltxt     string
	i        int
}

// newParser method is used to create a new instance of parser and return it
func newParser(ltxt string) *parser {
	p := &parser{}
	p.strList = make([]string, 0)
	p.ltxt = ltxt
	return p
}

// handleWhiteSpace method is used to update the parser state when a whitespace character is passed
// " ","\n","\r","\t"
func (p *parser) handleWhiteSpace(c string) error {
	if !p.inQuotes && !p.inSBkts {
		if p.current == "" {
			return nil
		}
		p.strList = append(p.strList, p.current)
		p.current = ""
		return nil
	}
	p.current += c
	return nil
}

// handleClosure method is used to update the parser state when a bar is passed (|)
func (p *parser) handleClosure(c string) error {
	if p.inSBkts {
		p.current += c
		p.inSBkts = false
		return nil
	}
	if !p.inQuotes {
		if p.current != "" {
			return Error("Closure not separated by a space")
		}
		p.current += c
		p.inSBkts = true
		return nil
	}
	p.current += c
	return nil
}

// handleQuote method is used to update the parser state when a double-quote is passed
func (p *parser) handleQuote(c string) error {
	if !p.inQuotes {
		p.inQuotes = true
		if p.current != "" {
			return Error("Not separated by a space")
		}
		p.current += "\""
		return nil
	}
	p.current += "\""
	p.strList = append(p.strList, p.current)
	p.inQuotes = false
	p.current = ""
	return nil
}

// handleOpenSquareBracket method is used to update the parser state when a open square bracket is passed
func (p *parser) handleOpenSquareBracket(c string) error {
	if p.inSBkts {
		return Error("'[' not expected inside a closure")
	}
	if p.inQuotes {
		p.current += c
		return nil
	}
	if p.current != "" {
		return Error("Array not separated by a space")
	}
	off := MatchChars(p.ltxt, int64(p.i), '[', ']') + 1
	p.current = p.ltxt[p.i:off]
	p.strList = append(p.strList, p.current)
	p.i = int(off)
	if strings.TrimSpace(string(p.ltxt[p.i])) != "" && strings.TrimSpace(string(p.ltxt[p.i])) != ")" {
		return Error("Unexpected character found at array end : " + string(p.ltxt[p.i]))
	}
	p.current = ""
	return nil
}

// handleOpenParen method is used to update the parser state when a open parenthesis ("(") is passed
func (p *parser) handleOpenParen(c string) error {
	if p.inSBkts {
		return Error("'(' not expected inside a closure")
	}
	if p.inQuotes {
		p.current += c
		return nil
	}
	if p.current != "" {
		return Error("Expression not separated by a space")
	}
	off := MatchChars(p.ltxt, int64(p.i), '(', ')') + 1
	p.current = p.ltxt[p.i:off]
	p.strList = append(p.strList, p.current)
	p.i = int(off)
	if strings.TrimSpace(string(p.ltxt[p.i])) != "" && strings.TrimSpace(string(p.ltxt[p.i])) != ")" {
		return Error("Unexpected character found at expression end : " + string(p.ltxt[p.i]))
	}
	p.current = ""
	return nil
}

// handleCloseParen method is used to update the parser state when a close parenthesis (")") is passed
func (p *parser) handleCloseParen(c string) error {
	if p.inSBkts {
		return Error("')' not expected inside a closure")
	}
	if p.inQuotes {
		p.current += c
		return nil
	}
	if len(p.ltxt)-1 != p.i {
		return Error("Expected EOL, got " + string(p.ltxt[p.i]) + " at " + fmt.Sprint(p.i))
	}
	if p.current != "" {
		p.strList = append(p.strList, p.current)
		p.current = ""
	}
	return nil
}

// handleCloseSquareBracket method is used to update the parser state when a close square bracket ("]") is passed
func (p *parser) handleCloseSquareBracket(c string) error {
	if p.inSBkts {
		return Error("']' not expected inside a closure")
	}
	if p.inQuotes {
		p.current += c
		return nil
	}
	p.strList = append(p.strList, p.current)
	p.current = ""
	return nil
}

// handleDefault method is used to update the parser state when a normal character (which is not a symbol) is passed
func (p *parser) handleDefault(c string) error {
	p.current += c
	return nil
}

// runHandler method runs a handler function corresponding to the string passed and returns the error
// result from the handler
func (p *parser) runHandler(c string) error {
	switch c {
	case " ", "\n", "\r", "\t":
		return p.handleWhiteSpace(c)
	case "|":
		return p.handleClosure(c)
	case "\"":
		return p.handleQuote(c)
	case "[":
		return p.handleOpenSquareBracket(c)
	case "(":
		return p.handleOpenParen(c)
	case ")":
		return p.handleCloseParen(c)
	case "]":
		return p.handleCloseSquareBracket(c)
	}
	return p.handleDefault(c)
}
