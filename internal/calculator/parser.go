package calculation

import (
	"fmt"
	"strconv"
	"unicode"
)

func Parse(expr string) (float64, error) {
	p := &Parser{expr: expr, pos: 0}
	return p.parseExpression()
}

type Parser struct {
	expr string
	pos  int
}

func (p *Parser) parseExpression() (float64, error) {
	left, err := p.parseTerm()
	if err != nil {
		return 0, err
	}
	for p.pos < len(p.expr) {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}
		op := p.expr[p.pos]
		if op != '+' && op != '-' {
			break
		}
		p.pos++
		right, err := p.parseTerm()
		if err != nil {
			return 0, err
		}
		if op == '+' {
			left += right
		} else {
			left -= right
		}
	}
	return left, nil
}

func (p *Parser) parseTerm() (float64, error) {
	left, err := p.parseFactor()
	if err != nil {
		return 0, err
	}
	for p.pos < len(p.expr) {
		p.skipWhitespace()
		if p.pos >= len(p.expr) {
			break
		}
		op := p.expr[p.pos]
		if op != '*' && op != '/' {
			break
		}
		p.pos++
		right, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		if op == '*' {
			left *= right
		} else {
			if right == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			left /= right
		}
	}
	return left, nil
}

func (p *Parser) parseFactor() (float64, error) {
	p.skipWhitespace()
	if p.pos >= len(p.expr) {
		return 0, fmt.Errorf("unexpected end of expression")
	}
	if p.expr[p.pos] == '+' {
		p.pos++
		return p.parseFactor()
	}
	if p.expr[p.pos] == '-' {
		p.pos++
		val, err := p.parseFactor()
		if err != nil {
			return 0, err
		}
		return -val, nil
	}
	if p.expr[p.pos] == '(' {
		p.pos++
		val, err := p.parseExpression()
		if err != nil {
			return 0, err
		}
		p.skipWhitespace()
		if p.pos >= len(p.expr) || p.expr[p.pos] != ')' {
			return 0, fmt.Errorf("missing closing parenthesis")
		}
		p.pos++
		return val, nil
	}
	return p.parseNumber()
}

func (p *Parser) parseNumber() (float64, error) {
	p.skipWhitespace()
	if p.pos >= len(p.expr) {
		return 0, fmt.Errorf("expected number")
	}
	if !isDigit(rune(p.expr[p.pos])) && p.expr[p.pos] != '.' {
		return 0, fmt.Errorf("expected number")
	}
	start := p.pos
	seenDot := false
	for p.pos < len(p.expr) {
		if p.expr[p.pos] == '.' {
			if seenDot {
				return 0, fmt.Errorf("invalid number format (multiple decimal points)")
			}
			seenDot = true
		} else if !isDigit(rune(p.expr[p.pos])) {
			break
		}
		p.pos++
	}
	numStr := p.expr[start:p.pos]
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number format: %s", numStr)
	}
	return num, nil
}

func (p *Parser) skipWhitespace() {
	for p.pos < len(p.expr) && unicode.IsSpace(rune(p.expr[p.pos])) {
		p.pos++
	}
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
