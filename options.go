package participle

import (
	"fmt"
	"reflect"

	"github.com/alecthomas/participle/v2/lexer"
)

// An Option to modify the behaviour of the Parser.
type Option func(p *Parser) error

// Lexer is an Option that sets the lexer to use with the given grammar.
func Lexer(def lexer.Definition) Option {
	return func(p *Parser) error {
		p.lex = def
		return nil
	}
}

// UseLookahead allows branch lookahead up to "n" tokens.
//
// If parsing cannot be disambiguated before "n" tokens of lookahead, parsing will fail.
//
// Note that increasing lookahead has a minor performance impact, but also
// reduces the accuracy of error reporting.
func UseLookahead(n int) Option {
	return func(p *Parser) error {
		p.useLookahead = n
		return nil
	}
}

// CaseInsensitive allows the specified token types to be matched case-insensitively.
//
// Note that the lexer itself will also have to be case-insensitive; this option
// just controls whether literals in the grammar are matched case insensitively.
func CaseInsensitive(tokens ...string) Option {
	return func(p *Parser) error {
		for _, token := range tokens {
			p.caseInsensitive[token] = true
		}
		return nil
	}
}

func UseInterface(parseFn interface{}) Option {
	errType := reflect.TypeOf((*error)(nil)).Elem()
	lexType := reflect.TypeOf((*lexer.PeekingLexer)(nil))
	return func(parser *Parser) error {
		fv := reflect.ValueOf(parseFn)
		ft := fv.Type()
		if ft.Kind() != reflect.Func {
			return fmt.Errorf("parseFn must be a function (got %s)", ft)
		}
		if ft.NumIn() != 1 || ft.In(0) != lexType {
			return fmt.Errorf("parseFn must only take one arg of type %s", lexType)
		}
		if ft.NumOut() != 2 || ft.Out(0).Kind() != reflect.Interface || ft.Out(1) != errType {
			return fmt.Errorf("parseFn must return two values, the interface to produce & an error")
		}
		parser.interfaceParsers[ft.Out(0)] = fv
		return nil
	}
}

// ParseOption modifies how an individual parse is applied.
type ParseOption func(p *parseContext)

// AllowTrailing tokens without erroring.
//
// That is, do not error if a full parse completes but additional tokens remain.
func AllowTrailing(ok bool) ParseOption {
	return func(p *parseContext) {
		p.allowTrailing = ok
	}
}
