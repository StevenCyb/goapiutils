package subset

import (
	"strings"

	"github.com/StevenCyb/goquery/errs"
	"github.com/StevenCyb/goquery/tokenizer"
)

// Types that are used in this parser
const (
	TYPE_SKIP           tokenizer.Type = "SKIP"
	TYPE_JOIN           tokenizer.Type = ","
	TYPE_PATH_SEPARATOR tokenizer.Type = "."
	TYPE_FIELD_NAME     tokenizer.Type = "FIELD_NAME"
)

// specialEncode is the map for encoding
// a list of special characters
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	` `: "%20",
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{}
}

// Parser provides the logic to parse
// rsql statements
type Parser struct {
	tokenizer *tokenizer.Tokenizer
	lookahead *tokenizer.Token
}

// eat return a token with expected type
func (p *Parser) eat(tokenType tokenizer.Type) (*tokenizer.Token, error) {
	token := p.lookahead

	if token == nil {
		return nil, errs.NewErrUnexpectedInputEnd(tokenType.String())
	}
	if token.Type != tokenType {
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPostion(),
			token.Type.String(),
			tokenType.String(),
		)
	}

	var err error
	p.lookahead, err = p.tokenizer.GetNextToken()
	return token, err
}

// Parse a given query
func (p *Parser) Parse(query string, fullObject interface{}) (interface{}, error) {
	if query == "" {
		return fullObject, nil
	}

	for dec, enc := range specialEncode {
		query = strings.ReplaceAll(query, enc, dec)
	}

	p.tokenizer = tokenizer.NewTokenizer(
		query,
		TYPE_SKIP, TYPE_SKIP,
		[]*tokenizer.Spec{
			tokenizer.NewSpec(`^\s+`, TYPE_SKIP),
			tokenizer.NewSpec(`^,`, TYPE_JOIN),
			tokenizer.NewSpec(`^\.`, TYPE_PATH_SEPARATOR),
			tokenizer.NewSpec(`^[^\.,=]*`, TYPE_FIELD_NAME),
		},
		nil,
	)

	var err error
	p.lookahead, err = p.tokenizer.GetNextToken()
	if err != nil {
		return nil, err
	}

	subsetObject := map[string]interface{}{}
	err = p.expression(fullObject, &subsetObject)
	return subsetObject, err
}

/**
 * <expression>
 * 	: <subset_spec>
 * 	| <subset_spec> ',' <expression>
 */
func (p *Parser) expression(object interface{}, subsetObject *map[string]interface{}) error {
	err := p.subsetSpec(object, subsetObject)
	if err != nil {
		return err
	}

	if p.tokenizer.HasMoreTokens() && p.lookahead.Type == TYPE_JOIN {
		_, err = p.eat(TYPE_JOIN)
		if err != nil {
			return err
		}

		return p.expression(object, subsetObject)
	}

	return nil
}

/**
 * <subset_spec>
 * 	: <field_name> '=' <field_name>
 *	| <field_name> "." <subset_spec>
 */
func (p *Parser) subsetSpec(object interface{}, subsetObject *map[string]interface{}) error {
	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(TYPE_FIELD_NAME.String())
	}

	fieldNameToken, err := p.eat(TYPE_FIELD_NAME)
	if err != nil {
		return err
	}

	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(TYPE_FIELD_NAME.String())
	}
	if p.lookahead.Type == 

	return nil
}
