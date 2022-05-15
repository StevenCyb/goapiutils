package tokenizer

// PolicyType represent policy type values
type PolicyType byte

const (
	// WHITELIST_POLICY declarate a policy
	// to use whitelist logic
	WHITELIST_POLICY PolicyType = 0
	// BLACKLIST_POLICY declarate a policy
	// to use blacklist logic
	BLACKLIST_POLICY PolicyType = 1
)

// Policy handles policy checks
// based on configuration
type Policy struct {
	policyType PolicyType
	values     []string
}

// Allow check if a value is allowed
func (p *Policy) Allow(value string) bool {
	for _, policyValue := range p.values {
		if policyValue == value {
			return p.policyType == WHITELIST_POLICY
		}
	}

	return p.policyType == BLACKLIST_POLICY
}

// NewPolicy create a new policy instance with given arguments
func NewPolicy(policyType PolicyType, values ...string) *Policy {
	return &Policy{
		policyType: policyType,
		values:     values,
	}
}
