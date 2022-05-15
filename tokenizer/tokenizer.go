package tokenizer

import (
	"github.com/StevenCyb/goquery/errs"
)

// tokenizer that lazily pulls a token from a stream
type Tokenizer struct {
	cursor          int
	query           string
	skipTokenType   Type
	spec            []*Spec
	policyCheckType Type
	policy          *Policy
}

// GetCursorPostion return the position of the cursor
func (t *Tokenizer) GetCursorPostion() int {
	return t.cursor
}

// HasMoreTokens checks eather we still have more tokens
func (t *Tokenizer) HasMoreTokens() bool {
	return t.cursor < len(t.query)
}

// GetNextToken obtains next token
func (t *Tokenizer) GetNextToken() (*Token, error) {
	if !t.HasMoreTokens() {
		return nil, nil
	}

	part := t.query[t.cursor:]

	for _, spec := range t.spec {
		matched := spec.expression.FindString(part)
		if matched == "" {
			continue
		}

		t.cursor += len(matched)
		if spec.tokenType == t.skipTokenType {
			return t.GetNextToken()
		}

		if spec.tokenType == t.policyCheckType && t.policy != nil && !t.policy.Allow(matched) {
			return nil, errs.NewErrPolicyViolation(matched)
		}

		return NewToken(
			spec.tokenType,
			matched,
		), nil
	}

	return nil, errs.NewErrUnexpectedToken(
		t.cursor,
		part[:1])
}

// NewTokenizer create a new tokenizer instance
// with given parameters
func NewTokenizer(query string, skipTokenType, policyCheckType Type, spec []*Spec, policy *Policy) *Tokenizer {
	return &Tokenizer{
		cursor:          0,
		query:           query,
		skipTokenType:   skipTokenType,
		spec:            spec,
		policyCheckType: policyCheckType,
		policy:          policy,
	}
}
