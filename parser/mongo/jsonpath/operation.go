package jsonpath

type Operation string

const (
	// RemoveOperation is an operation to remove the value at the target location.
	// Requires `path`.
	RemoveOperation Operation = "remove"
	// AddOperation add a value or array to an array at the target location.
	// Requires `path` and `value`. Target of the path must be an array.
	AddOperation Operation = "add"
	// ReplaceOperation replaces the value at the target location
	// with a new value. Requires `path` and `value`.
	ReplaceOperation Operation = "replace"
	// MoveOperation removes the value at a specified location and
	// adds it to the target location. Requires `from` and `path`.
	MoveOperation Operation = "move"
	// CopyOperation copies the value from a specified location to the
	// target location. Requires `from` and `path`.
	CopyOperation Operation = "copy"
)

// OperationSpec specify an path operation.
type OperationSpec struct {
	From      Path        `json:"from"`
	Path      Path        `json:"path"`
	Value     interface{} `json:"value"`
	Operation Operation   `json:"op"` //nolint:tagliatelle
}

// Valid check if operation is valid.
func (o OperationSpec) Valid() bool {
	if o.Operation == "" {
		return false
	}

	switch o.Operation {
	case RemoveOperation:
		if !o.Path.Valid() {
			return false
		}
	case AddOperation:
		if !o.Path.Valid() {
			return false
		} else if o.Value == nil {
			return false
		}
	case ReplaceOperation:
		if !o.Path.Valid() {
			return false
		} else if o.Value == nil {
			return false
		}
	case MoveOperation:
		if !o.Path.Valid() {
			return false
		}

		if !o.From.Valid() {
			return false
		}
	case CopyOperation:
		if !o.Path.Valid() {
			return false
		}

		if !o.From.Valid() {
			return false
		}
	default:
		return false
	}

	return true
}
