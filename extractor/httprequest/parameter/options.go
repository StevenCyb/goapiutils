package parameter

const (
	base      = 10
	bitSize32 = 32
	bitSize64 = 64
)

// Option provides options for parameter extraction.
type Option struct {
	Key          string
	Default      string
	RegexPattern string
	Required     bool
}
