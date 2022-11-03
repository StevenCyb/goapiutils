package subset

import (
	"reflect"
	"strings"

	"github.com/StevenCyb/goapiutils/parser/errs"
	"github.com/StevenCyb/goapiutils/parser/tokenizer"
)

// Types that are used in this parser.
const (
	SkipType          tokenizer.Type = "SKIP"
	JoinType          tokenizer.Type = ","
	PathSeparatorType tokenizer.Type = "."
	AssignmentType    tokenizer.Type = "ASSIGNMENT"
	FieldNameType     tokenizer.Type = "FIELD_NAME"
)

// specialEncode is the map for encoding
// a list of special characters.
//
//nolint:gochecknoglobals
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	` `: "%20",
	`=`: "%5C%3D",
}

// NewParser creates a new parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parser provides the logic to parse rsql statements.
type Parser struct {
	tokenizer *tokenizer.Tokenizer
	lookahead *tokenizer.Token
}

// eat return a token with expected type.
func (p *Parser) eat(tokenType tokenizer.Type) (*tokenizer.Token, error) {
	token := p.lookahead

	if token == nil {
		return nil, errs.NewErrUnexpectedInputEnd(tokenType.String())
	}

	if token.Type != tokenType {
		return nil, errs.NewErrUnexpectedTokenType(
			p.tokenizer.GetCursorPosition(),
			token.Type.String(),
			tokenType.String(),
		)
	}

	var err error
	p.lookahead, err = p.tokenizer.GetNextToken()

	return token, err //nolint:wrapcheck
}

// Parse a given query.
func (p *Parser) Parse(query string, fullObject interface{}) (interface{}, error) {
	var err error

	if query == "" {
		return fullObject, nil
	}

	for dec, enc := range specialEncode {
		query = strings.ReplaceAll(query, enc, dec)
	}

	p.tokenizer = tokenizer.NewTokenizer(
		query,
		SkipType, SkipType,
		[]*tokenizer.Spec{
			tokenizer.NewSpec(`^\s+`, SkipType),
			tokenizer.NewSpec(`^,`, JoinType),
			tokenizer.NewSpec(`^\.`, PathSeparatorType),
			tokenizer.NewSpec(`^=`, AssignmentType),
			tokenizer.NewSpec(`^[^\.,=]*`, FieldNameType),
		},
		nil,
	)

	p.lookahead, err = p.tokenizer.GetNextToken()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	subsetObject := map[string]interface{}{}
	err = p.expression(reflect.ValueOf(fullObject), &subsetObject)

	return subsetObject, err
}

/*
 * <expression>
 * 	: <subset_spec>
 * 	| <subset_spec> ',' <expression>
 * .
 */
func (p *Parser) expression(object reflect.Value, subsetObject *map[string]interface{}) error {
	err := p.subsetSpec(object, subsetObject)
	if err != nil {
		return err
	}

	if p.tokenizer.HasMoreTokens() && p.lookahead.Type == JoinType {
		_, err = p.eat(JoinType)
		if err != nil {
			return err
		}

		return p.expression(object, subsetObject)
	}

	return nil
}

/*
 * <subset_spec>
 * 	: <field_name> '=' <field_name>
 *	| <field_name> "." <subset_spec>
 * .
 */
func (p *Parser) subsetSpec(object reflect.Value, subsetObject *map[string]interface{}) error {
	var newObject reflect.Value

	if p.lookahead == nil {
		return errs.NewErrUnexpectedInputEnd(FieldNameType.String())
	}

	fieldNameToken, err := p.eat(FieldNameType)
	if err != nil {
		return err
	}

	if object.IsNil() {
		return nil
	}

	if object.Kind() == reflect.Map {
		for _, key := range object.MapKeys() {
			if key.String() == fieldNameToken.Value {
				newObject = object.MapIndex(key)

				break
			}
		}
	}

	if reflect.ValueOf(newObject).IsZero() {
		return nil
	}

	if p.lookahead != nil {
		switch p.lookahead.Type {
		case PathSeparatorType:
			_, err := p.eat(PathSeparatorType)
			if err != nil {
				return err
			}

			if newObject.Kind() == reflect.Interface {
				return p.subsetSpec(newObject.Elem(), subsetObject)
			}

			return p.subsetSpec(newObject, subsetObject)
		case AssignmentType:
			_, err := p.eat(AssignmentType)
			if err != nil {
				return err
			}

			newFieldNameToken, err := p.eat(FieldNameType)
			if err != nil {
				return err
			}

			(*subsetObject)[newFieldNameToken.Value] = newObject.Interface()

			return nil
		}
	}

	return errs.NewErrUnexpectedInputEnd(AssignmentType.String())
}
