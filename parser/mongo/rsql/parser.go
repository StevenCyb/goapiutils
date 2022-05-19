package rsql

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/StevenCyb/goquery/errs"
	"github.com/StevenCyb/goquery/tokenizer"

	"go.mongodb.org/mongo-driver/bson"
)

// Types that are used in this parser
const (
	TYPE_SKIP                                 tokenizer.Type = "SKIP"
	TYPE_AND_COMPOSITE                        tokenizer.Type = ";"
	TYPE_OR_COMPOSITE                         tokenizer.Type = ","
	TYPE_CONTEXT_START                        tokenizer.Type = "("
	TYPE_CONTEXT_END                          tokenizer.Type = ")"
	TYPE_VALUE_COMPARE_OPERATOR               tokenizer.Type = "VALUE_COMPARE_OPERATOR"
	TYPE_QUOTED_STRING_VALUE_COMPARE_OPERATOR tokenizer.Type = "QUOTED_STRING_VALUE_COMPARE_OPERATOR"
	TYPE_NUMERIC_VALUE_COMPARE_OPERATOR       tokenizer.Type = "NUMERIC_VALUE_COMPARE_OPERATOR"
	TYPE_ARRAY_COMPARE_OPERATOR               tokenizer.Type = "ARRAY_COMPARE_OPERATOR"
	TYPE_BOOL_LITERAL                         tokenizer.Type = "BOOL_LITERAL"
	TYPE_QUOTED_STRING_LITERAL                tokenizer.Type = "QUOTED_STRING_LITERAL"
	TYPE_FIELD_NAME                           tokenizer.Type = "FIELD_NAME"
	TYPE_NUMBER_LITERAL                       tokenizer.Type = "NUMERIC_LITERAL"
)

// specialEncode is the map for encoding
// a list of special characters
var specialEncode = map[string]string{
	`,`: "%5C%2C",
	`;`: "%5C%3B",
	`=`: "%5C%3D",
	`"`: "%22",
	`'`: "%27",
	` `: "%20",
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
			tokenizer.NewSpec(`^\(`, TYPE_CONTEXT_START),
			tokenizer.NewSpec(`^\)`, TYPE_CONTEXT_END),
			tokenizer.NewSpec(`^;`, TYPE_AND_COMPOSITE),
			tokenizer.NewSpec(`^,`, TYPE_OR_COMPOSITE),
			tokenizer.NewSpec(`^(==|!=)`, TYPE_VALUE_COMPARE_OPERATOR),
			tokenizer.NewSpec(`^(=sw=|=ew=)`, TYPE_QUOTED_STRING_VALUE_COMPARE_OPERATOR),
			tokenizer.NewSpec(`^(=gt=|=ge=|=lt=|=le=)`, TYPE_NUMERIC_VALUE_COMPARE_OPERATOR),
			tokenizer.NewSpec(`^(=in=|=out=)`, TYPE_ARRAY_COMPARE_OPERATOR),
			tokenizer.NewSpec(`(?i)^(true|false)`, TYPE_BOOL_LITERAL),
			tokenizer.NewSpec(`^(-|\+)?\d+(\.\d+)?`, TYPE_NUMBER_LITERAL),
			tokenizer.NewSpec(`^("[^"]*"|'[^']*')`, TYPE_QUOTED_STRING_LITERAL),
			tokenizer.NewSpec(`^[^!=]*`, TYPE_FIELD_NAME),
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
 *   : <context>
 *   | <context> <composite_operator> <expression>
 *   | <comparison>
 *   | <comparison> <composite_operator> <expression>
 */
func (p *Parser) expression() ([]bson.E, error) {
	sortStatements := []bson.E{}
	if p.lookahead == nil {
		return nil, errs.NewErrUnexpectedInputEnd(TYPE_FIELD_NAME.String())
	}

	var left bson.E
	if p.lookahead.Type == TYPE_CONTEXT_START {
		tmp, err := p.context()
		if err != nil {
			return nil, err
		}
		left = tmp[0]
	} else {
		tmp, err := p.comparison()
		if err != nil {
			return nil, err
		}
		left = *tmp
	}

	if p.lookahead != nil && p.lookahead.Type != TYPE_CONTEXT_END {
		logicalOperation, err := p.compositeOperation()
		if err != nil {
			return nil, err
		}

		right, err := p.expression()
		if err != nil {
			return nil, err
		}

		if logicalOperation.Type == TYPE_AND_COMPOSITE {
			if right[0].Key == "$and" {
				newA := right[0].Value.(bson.A)
				newA = append(bson.A{bson.D{left}}, newA...)
				left = bson.E{Key: "$and", Value: newA}
			} else {
				left = bson.E{Key: "$and", Value: bson.A{bson.D{left}, bson.D{right[0]}}}
			}
		} else {
			if right[0].Key == "$or" {
				newA := right[0].Value.(bson.A)
				newA = append(bson.A{bson.D{left}}, newA...)
				left = bson.E{Key: "$or", Value: newA}
			} else {
				left = bson.E{Key: "$or", Value: bson.A{bson.D{left}, bson.D{right[0]}}}
			}
		}
		if err != nil {
			return nil, err
		}
	}

	sortStatements = append(sortStatements, left)

	return sortStatements, nil
}

/**
 * <context>
 *   : "(" <expression> ")"
 */
func (p *Parser) context() ([]bson.E, error) {
	_, err := p.eat(TYPE_CONTEXT_START)
	if err != nil {
		return nil, err
	}

	context, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.eat(TYPE_CONTEXT_END)
	if err != nil {
		return nil, err
	}

	return context, nil
}

/**
 * <composite_operator>
 *   : ";"
 *   | ","
 */
func (p *Parser) compositeOperation() (*tokenizer.Token, error) {
	if p.lookahead.Type == TYPE_AND_COMPOSITE {
		return p.eat(TYPE_AND_COMPOSITE)
	} else if p.lookahead.Type == TYPE_OR_COMPOSITE {
		return p.eat(TYPE_OR_COMPOSITE)
	}

	return nil, errs.NewErrUnexpectedTokenType(
		p.tokenizer.GetCursorPostion()-len(p.lookahead.Value),
		p.lookahead.Type.String(),
		";/,")
}

/**
 * <comparison>
 *   : TEXT <singular_operator> <literal>
 *   | TEXT <singular_string_operator> <quoted_string_literal>
 *   | TEXT <singular_numeric_operator> <numeric_literal>
 *   | TEXT <plural_operator> "(" <literal_list> ")"
 */
func (p *Parser) comparison() (*bson.E, error) {
	keyToken, err := p.eat(TYPE_FIELD_NAME)
	if err != nil {
		return nil, err
	}
	key := keyToken.Value

	if p.lookahead.Type == TYPE_VALUE_COMPARE_OPERATOR {
		operator, err := p.eat(TYPE_VALUE_COMPARE_OPERATOR)
		if err != nil {
			return nil, err
		}

		var literal interface{}
		if p.lookahead.Type == TYPE_CONTEXT_START {
			_, err = p.eat(TYPE_CONTEXT_START)
			if err != nil {
				return nil, err
			}

			literal, err = p.literalList()
			if err != nil {
				return nil, err
			}

			_, err = p.eat(TYPE_CONTEXT_END)
			if err != nil {
				return nil, err
			}
		} else {
			literal, err = p.literal()
			if err != nil {
				return nil, err
			}
		}

		switch operator.Value {
		case "==":
			return &bson.E{Key: key, Value: literal}, nil
		case "!=":
			return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$ne", Value: literal}}}, nil
		default:
			return nil, errs.NewErrUnexpectedTokenType(
				p.tokenizer.GetCursorPostion()-len(operator.Value),
				operator.Type.String(),
				TYPE_VALUE_COMPARE_OPERATOR.String())
		}
	} else if p.lookahead.Type == TYPE_QUOTED_STRING_VALUE_COMPARE_OPERATOR {
		operator, err := p.eat(TYPE_QUOTED_STRING_VALUE_COMPARE_OPERATOR)
		if err != nil {
			return nil, err
		}

		literal, err := p.stringLiteral()
		if err != nil {
			return nil, err
		}

		switch operator.Value {
		case "=sw=":
			wildcard, err := regexp.Compile("^" + fmt.Sprintf("%v", literal))
			return &bson.E{Key: key, Value: *wildcard}, err
		case "=ew=":
			wildcard, err := regexp.Compile(fmt.Sprintf("%v", literal) + "$")
			return &bson.E{Key: key, Value: *wildcard}, err
		default:
			return nil, errs.NewErrUnexpectedTokenType(
				p.tokenizer.GetCursorPostion()-len(operator.Value),
				operator.Type.String(),
				TYPE_VALUE_COMPARE_OPERATOR.String())
		}
	} else if p.lookahead.Type == TYPE_NUMERIC_VALUE_COMPARE_OPERATOR {
		operator, err := p.eat(TYPE_NUMERIC_VALUE_COMPARE_OPERATOR)
		if err != nil {
			return nil, err
		}

		literal, err := p.numericLiteral()
		if err != nil {
			return nil, err
		}

		switch operator.Value {
		case "=gt=":
			return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$gt", Value: literal}}}, nil
		case "=ge=":
			return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$gte", Value: literal}}}, nil
		case "=lt=":
			return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$lt", Value: literal}}}, nil
		case "=le=":
			return &bson.E{Key: key, Value: bson.D{bson.E{Key: "$lte", Value: literal}}}, nil
		default:
			return nil, errs.NewErrUnexpectedTokenType(
				p.tokenizer.GetCursorPostion()-len(operator.Value),
				operator.Type.String(),
				TYPE_VALUE_COMPARE_OPERATOR.String())
		}
	} else if p.lookahead.Type == TYPE_ARRAY_COMPARE_OPERATOR {
		operator, err := p.eat(TYPE_ARRAY_COMPARE_OPERATOR)
		if err != nil {
			return nil, err
		}

		_, err = p.eat(TYPE_CONTEXT_START)
		if err != nil {
			return nil, err
		}

		literalList, err := p.literalList()
		if err != nil {
			return nil, err
		}

		_, err = p.eat(TYPE_CONTEXT_END)
		if err != nil {
			return nil, err
		}

		switch operator.Value {
		case "=in=":
			return &bson.E{Key: key, Value: bson.E{Key: "$in", Value: literalList}}, nil
		case "=out=":
			return &bson.E{Key: key, Value: bson.E{Key: "$nin", Value: literalList}}, nil
		default:
			return nil, errs.NewErrUnexpectedTokenType(
				p.tokenizer.GetCursorPostion()-len(operator.Value),
				operator.Type.String(),
				TYPE_ARRAY_COMPARE_OPERATOR.String())
		}
	}

	return nil, errs.NewErrUnexpectedToken(
		p.tokenizer.GetCursorPostion()-len(p.lookahead.Value),
		p.lookahead.Value)
}

/**
 * <literal>
 * : <bool_literal>
 * | <quoted_string_literal>
 * | <numeric_literal>
 */
func (p *Parser) literal() (interface{}, error) {
	if p.lookahead.Type == TYPE_BOOL_LITERAL {
		token, err := p.eat(TYPE_BOOL_LITERAL)
		if err != nil {
			return nil, err
		}

		return strings.ToLower(token.Value) == "true", nil
	} else if p.lookahead.Type == TYPE_QUOTED_STRING_LITERAL {
		return p.stringLiteral()
	} else if p.lookahead.Type == TYPE_NUMBER_LITERAL {
		return p.numericLiteral()
	}

	return nil, errs.NewErrUnexpectedTokenType(
		p.tokenizer.GetCursorPostion()-len(p.lookahead.Value),
		p.lookahead.Type.String(),
		"LITERAL")
}

/**
 * <quoted_string_literal>
 * : "'" <TEXT> "'"
 * | """ <TEXT> """
 */
func (p *Parser) stringLiteral() (interface{}, error) {
	token, err := p.eat(TYPE_QUOTED_STRING_LITERAL)
	if err != nil {
		return nil, err
	}

	replacer := strings.NewReplacer(`"`, "", "'", "")
	return replacer.Replace(token.Value), nil
}

/**
 * <numeric_literal>
 * : <INT>
 * | <FLOAT>
 */
func (p *Parser) numericLiteral() (interface{}, error) {
	token, err := p.eat(TYPE_NUMBER_LITERAL)
	if err != nil {
		return nil, err
	}

	if strings.Contains(token.Value, ".") {
		return strconv.ParseFloat(token.Value, 64)
	}

	return strconv.ParseInt(token.Value, 10, 64)
}

/**
 * <literal_list>
 * : <quoted_string_literal> "," <literal_list>
 * | <numeric_literal> "," <literal_list>
 */
func (p *Parser) literalList() (bson.A, error) {
	items := bson.A{}

	body, err := p.literal()
	if err != nil {
		return nil, err
	}
	items = append(items, body)

	for p.lookahead.Type == TYPE_OR_COMPOSITE {
		p.eat(TYPE_OR_COMPOSITE)
		body, err = p.literal()
		if err != nil {
			return nil, err
		}
		items = append(items, body)
	}

	return items, nil
}
