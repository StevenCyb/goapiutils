package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// MatchingKindRule is a default rule that is applied to all fields.
// This rules checks for type and name matches to prevent input for
// unknown fields or to violate types.
type MatchingKindRule struct {
	instance interface{}
}

// Tag returns tag of the rule.
func (m MatchingKindRule) Tag() string {
	return "type"
}

// UseValue initializes the rule for specified field.
func (m *MatchingKindRule) UseValue(_ operation.Path, _ reflect.Kind, instance interface{}, _ string) error {
	m.instance = instance

	return nil
}

// Apply rule on given patch operation specification.
func (m MatchingKindRule) Apply(operationSpec operation.Spec) error {
	return m.deepCompareType("root value", m.instance, operationSpec.Value)
}

// deepCompareType checks recursively one interface against a reference.
//
//nolint:funlen
func (m MatchingKindRule) deepCompareType(name string, reference, object interface{}) error {
	var (
		err           error
		referenceType = reflect.TypeOf(reference)
		objectType    = reflect.TypeOf(object)
		referenceKind = referenceType.Kind()
		objectKind    = objectType.Kind()
	)

	if referenceKind != objectKind {
		return TypeMismatchError{name: name, actual: objectKind, expected: referenceKind}
	}

	switch objectType.Kind() { //nolint:exhaustive
	case reflect.Ptr:
		err = m.deepCompareType(name, reflect.Zero(referenceType.Elem()).Interface(),
			reflect.Zero(objectType.Elem()).Interface())
	case reflect.Array, reflect.Map, reflect.Slice:
		var (
			referenceZeroValue = reflect.Zero(referenceType.Elem())
			objectZeroValue    = reflect.Zero(objectType.Elem())
		)

		if objectType.Kind() == reflect.Map && referenceType.Kind() == reflect.Map {
			if referenceType.Key().Kind() != objectType.Key().Kind() {
				return TypeMismatchError{name: name, actual: objectType.Key().Kind(), expected: referenceType.Key().Kind()}
			}
		}

		err = m.deepCompareType(name+" item", referenceZeroValue.Interface(), objectZeroValue.Interface())
	case reflect.Struct:
		for i := 0; i < objectType.NumField(); i++ {
			var (
				objectField = objectType.Field(i)
				objectName  = objectField.Name
				found       = false
			)

			for i := 0; i < referenceType.NumField(); i++ {
				var (
					referenceField = referenceType.Field(i)
					referenceName  = referenceField.Tag.Get("bson")
					zeroValue      = reflect.Zero(referenceField.Type)
				)

				if referenceField.Type.Kind() == reflect.Ptr {
					zeroValue = reflect.Zero(referenceField.Type.Elem())
				}

				if objectName == referenceName {
					err = m.deepCompareType(objectName, zeroValue.Interface(), reflect.Zero(objectField.Type).Interface())

					found = true

					break
				}
			}

			if !found {
				err = UnknownFieldError{name: objectName}

				break
			}
		}
	}

	return err
}
