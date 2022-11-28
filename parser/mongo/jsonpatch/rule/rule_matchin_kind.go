//nolint:ireturn
package rule

import (
	"reflect"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
)

// MatchingKindRule is a default rule that is applied to all fields.
// This rules checks for type and name matches to prevent input for
// unknown fields or to violate types.
type MatchingKindRule struct {
	Instance interface{}
}

// UseValue instantiate new rule instance for field.
func (m *MatchingKindRule) NewInstance(_ string, _ reflect.Kind, instance interface{}, _ string) (Rule, error) {
	return &MatchingKindRule{
		Instance: instance,
	}, nil
}

// NewInheritInstance instantiate new rule instance based on given rule.
func (m *MatchingKindRule) NewInheritInstance(_ string, _ reflect.Kind, instance interface{}) (Rule, error) {
	return &MatchingKindRule{
		Instance: instance,
	}, nil
}

// Validate applies rule on given patch operation specification.
func (m MatchingKindRule) Validate(operationSpec operation.Spec) error {
	return m.deepCompareType("root value", m.Instance, operationSpec.Value)
}

// deepCompareType checks recursively one interface against a reference.
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
		err = m.deepCompareIterable(name, referenceType, objectType)
	case reflect.Struct:
		err = m.deepCompareStruct(referenceType, objectType)
	}

	return err
}

func (m MatchingKindRule) deepCompareIterable(name string, referenceType, objectType reflect.Type) error {
	var (
		referenceZeroValue = reflect.Zero(referenceType.Elem())
		objectZeroValue    = reflect.Zero(objectType.Elem())
	)

	if objectType.Kind() == reflect.Map && referenceType.Kind() == reflect.Map {
		if referenceType.Key().Kind() != objectType.Key().Kind() {
			return TypeMismatchError{
				name: name, actual: objectType.Key().Kind(), expected: referenceType.Key().Kind(), forKey: true,
			}
		}
	}

	return m.deepCompareType(name+" item", referenceZeroValue.Interface(), objectZeroValue.Interface())
}

func (m MatchingKindRule) deepCompareStruct(referenceType, objectType reflect.Type) error {
	var err error

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

	return err
}
