package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// Rule defines the interface for a patch operation rule.
type Rule interface {
	Tag() string
	UseValue(patch operation.Path, kind reflect.Kind, instance interface{}, value string) error
	Apply(operationSpec operation.Spec) error
}
