package runner

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/pawnkit/pawntest/internal/backend"
)

const tagPublicPrefix = "ptt_"

func testTags(publics []backend.Public) map[string]map[string]bool {
	out := map[string]map[string]bool{}
	for _, public := range publics {
		if !strings.HasPrefix(public.Name, tagPublicPrefix) {
			continue
		}
		parts := strings.SplitN(strings.TrimPrefix(public.Name, tagPublicPrefix), "_", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			continue
		}
		test := "test_" + parts[1]
		if out[test] == nil {
			out[test] = map[string]bool{}
		}
		out[test][parts[0]] = true
	}
	return out
}

type tagExpression struct {
	tokens []string
	index  int
}

func parseTagExpression(input string) (*tagExpression, error) {
	var tokens []string
	for index := 0; index < len(input); {
		r := rune(input[index])
		if unicode.IsSpace(r) {
			index++
			continue
		}
		if strings.ContainsRune("!&|()", r) {
			tokens = append(tokens, string(r))
			index++
			continue
		}
		start := index
		for index < len(input) {
			r = rune(input[index])
			if !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-') {
				break
			}
			index++
		}
		if start == index {
			return nil, fmt.Errorf("invalid tag expression character %q", input[index])
		}
		tokens = append(tokens, input[start:index])
	}
	if len(tokens) == 0 {
		return nil, nil
	}
	expression := &tagExpression{tokens: tokens}
	if _, err := expression.parseOr(nil); err != nil {
		return nil, err
	}
	if expression.index != len(tokens) {
		return nil, fmt.Errorf("unexpected token %q in tag expression", tokens[expression.index])
	}
	expression.index = 0
	return expression, nil
}

func (expression *tagExpression) matches(tags map[string]bool) bool {
	if expression == nil {
		return true
	}
	expression.index = 0
	value, _ := expression.parseOr(tags)
	return value
}

func (expression *tagExpression) parseOr(tags map[string]bool) (bool, error) {
	value, err := expression.parseAnd(tags)
	for err == nil && expression.peek() == "|" {
		expression.index++
		right, nextErr := expression.parseAnd(tags)
		value, err = value || right, nextErr
	}
	return value, err
}

func (expression *tagExpression) parseAnd(tags map[string]bool) (bool, error) {
	value, err := expression.parseUnary(tags)
	for err == nil && expression.peek() == "&" {
		expression.index++
		right, nextErr := expression.parseUnary(tags)
		value, err = value && right, nextErr
	}
	return value, err
}

func (expression *tagExpression) parseUnary(tags map[string]bool) (bool, error) {
	if expression.peek() == "!" {
		expression.index++
		value, err := expression.parseUnary(tags)
		return !value, err
	}
	if expression.peek() == "(" {
		expression.index++
		value, err := expression.parseOr(tags)
		if err != nil {
			return false, err
		}
		if expression.peek() != ")" {
			return false, fmt.Errorf("missing closing parenthesis in tag expression")
		}
		expression.index++
		return value, nil
	}
	token := expression.peek()
	if token == "" || strings.Contains("&|)", token) {
		return false, fmt.Errorf("expected tag in expression")
	}
	expression.index++
	return tags[token], nil
}

func (expression *tagExpression) peek() string {
	if expression.index >= len(expression.tokens) {
		return ""
	}
	return expression.tokens[expression.index]
}
