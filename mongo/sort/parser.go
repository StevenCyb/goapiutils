package sort

import (
	"strings"

	"github.com/StevenCyb/goquery/errs"
	"github.com/StevenCyb/goquery/tokenizer"

	"go.mongodb.org/mongo-driver/bson"
)

// Types that are used in this parser
const (
	TYPE_SKIP           tokenizer.Type = "SKIP"
	TYPE_AND            tokenizer.Type = ","
	TYPE_SET            tokenizer.Type = "="
	TYPE_SORT_CONDITION tokenizer.Type = "SORT_CRITERIA"
	TYPE_FIELD_NAME     tokenizer.Type = "FIELD_NAME"
)

// specialEncode is the map for encoding
// a list of special characters
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	`=`: "%5C%3D",
}

// NewParser creates a new parser
func NewParser(policy *tokenizer.Policy) *Parser {
	return &Parser{
		policy: policy,
	}
}

// Parser provides the logic to parse
// rsql statements
type Parser struct {
	tokenizer *tokenizer.Tokenizer
	lookahead *tokenizer.Token
	policy    *tokenizer.Policy
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
func (p *Parser) Parse(query string) (bson.D, error) {
	if query == "" {
		return bson.D{}, nil
	}

	for dec, enc := range specialEncode {
		query = strings.ReplaceAll(query, enc, dec)
	}

	p.tokenizer = tokenizer.NewTokenizer(
		query,
		TYPE_SKIP, TYPE_FIELD_NAME,
		[]*tokenizer.Spec{
			tokenizer.NewSpec(`^\s+`, TYPE_SKIP),
			tokenizer.NewSpec(`^,`, TYPE_AND),
			tokenizer.NewSpec(`^(=)`, TYPE_SET),
			tokenizer.NewSpec(`^(asc|desc|1|-1)`, TYPE_SORT_CONDITION),
			tokenizer.NewSpec(`^[^=]*`, TYPE_FIELD_NAME),
		},
		p.policy,
	)

	var err error
	p.lookahead, err = p.tokenizer.GetNextToken()
	if err != nil {
		return nil, err
	}

	return p.expression()
}

/**
 * <expression>
 *   | <sort_statement>
 *   | <sort_statement> "," <sort_statement>
 */
func (p *Parser) expression() ([]bson.E, error) {
	sortStatements := []bson.E{}
	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(TYPE_FIELD_NAME.String())
	}

	sortStatement, err := p.sortStatement()
	if err != nil {
		return nil, err
	}
	sortStatements = append(sortStatements, *sortStatement)

	if p.lookahead != nil {
		_, err := p.eat(TYPE_AND)
		if err != nil {
			return nil, err
		}

		nextSortStatements, err := p.expression()
		if err != nil {
			return nil, err
		}
		sortStatements = append(sortStatements, nextSortStatements...)
	}

	return sortStatements, nil
}

/**
 * <sort_statement>
 *   : <key> "=" <sort_condition>
 */
func (p *Parser) sortStatement() (*bson.E, error) {
	keyToken, err := p.eat(TYPE_FIELD_NAME)
	if err != nil {
		return nil, err
	}

	_, err = p.eat(TYPE_SET)
	if err != nil {
		return nil, err
	}

	sortConditionToken, err := p.eat(TYPE_SORT_CONDITION)
	if err != nil {
		return nil, err
	}

	sort := 0
	if sortConditionToken.Value == "asc" || sortConditionToken.Value == "1" {
		sort = 1
	} else {
		sort = -1
	}

	return &bson.E{Key: keyToken.Value, Value: sort}, nil
}
