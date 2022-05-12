package tokenizer

type PolicyType byte

const (
	WHITELIST_POLICY PolicyType = 0
	BLACKLIST_POLICY PolicyType = 1
)

type Policy struct {
	policyType PolicyType
	values     []string
}

func (p *Policy) Allow(value string) bool {
	for _, policyValue := range p.values {
		if policyValue == value {
			return p.policyType == WHITELIST_POLICY
		}
	}

	return p.policyType == BLACKLIST_POLICY
}

func NewPolicy(policyType PolicyType, values ...string) *Policy {
	return &Policy{
		policyType: policyType,
		values:     values,
	}
}
