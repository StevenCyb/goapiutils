package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/operation"
	"github.com/StevenCyb/goapiutils/parser/mongo/jsonpatch/rule"
)

const (
	prefix = "jp_"
)

// Validator interprets reference to validate JSON patch operations.
type Validator struct {
	knownTagRules map[string]rule.Rule
	rules         map[operation.Path][]rule.Rule
	wildcardRules map[operation.Path][]rule.Rule
	generalRules  []rule.Rule
}

// RegisterRule register addition rule for specific key (must have prefix `jp_`).
func (v *Validator) RegisterRule(key string, rule rule.Rule) error {
	if rule == nil {
		return ErrNilRule
	}

	if !strings.HasPrefix(key, prefix) {
		return ErrMissingPrefix
	}

	if _, exists := v.knownTagRules[key]; exists {
		return ErrDuplicateRuleTags
	}

	v.knownTagRules[key] = rule

	return nil
}

// Validate a given JSON patch operations again rules.
func (v Validator) Validate(operationSpec operation.Spec) error {
	if rules, match := v.rules[operationSpec.Path]; match {
		for _, rule := range rules {
			err := rule.Validate(operationSpec)
			if err != nil {
				return fmt.Errorf("operation no allowed: %w", err)
			}
		}
	}

	for path, rules := range v.wildcardRules {
		if path.Equal(operationSpec.Path) {
			for _, rule := range rules {
				err := rule.Validate(operationSpec)
				if err != nil {
					return fmt.Errorf("operation no allowed: %w", err)
				}
			}

			break
		}
	}

	return nil
}

// UseReference interpret given reference to model rule set.
func (v *Validator) UseReference(referenceType reflect.Type) error {
	return v.parseReference(referenceType, "", nil)
}

func (v *Validator) parseReference(objectType reflect.Type, path string, inheritedRules []rule.Rule) error {
	var err error

	if objectType == nil {
		return ErrReferenceIsNil
	}

	if path != "" && objectType.Kind() != reflect.Ptr {
		path += "."
	}

	switch objectType.Kind() { //nolint:exhaustive
	case reflect.Ptr:
		err = v.parseReference(reflect.Zero(objectType.Elem()).Type(), path, inheritedRules)
	case reflect.Array, reflect.Map, reflect.Slice:
		err = v.parseReferenceIterable(objectType, path, inheritedRules)
	case reflect.Struct:
		err = v.parseReferenceStruct(objectType, path, inheritedRules)
	}

	return err
}

func (v *Validator) parseReferenceIterable(objectType reflect.Type, path string, inheritedRules []rule.Rule) error {
	var (
		zeroValue         reflect.Value
		valueType         = objectType.Elem()
		kind              = valueType.Kind()
		err               error
		newInheritedRules []rule.Rule
	)

	switch kind { //nolint:exhaustive
	case reflect.Ptr:
		zeroValue = reflect.Zero(valueType.Elem())
		kind = zeroValue.Kind()
	case reflect.Array, reflect.Slice:
		zeroValue = reflect.MakeSlice(valueType, 0, 0)
	case reflect.Map:
		zeroValue = reflect.MakeMap(valueType)
		kind = zeroValue.Kind()
	default:
		zeroValue = reflect.Zero(valueType)
	}

	if newInheritedRules, err = v.addRule(path+"*", kind, "", zeroValue.Interface(), inheritedRules); err != nil {
		return err
	}

	if err := v.parseReference(zeroValue.Type(), path+"*", newInheritedRules); err != nil {
		return err
	}

	return nil
}

func (v *Validator) parseReferenceStruct(objectType reflect.Type, path string, inheritedRules []rule.Rule) error {
	for i := 0; i < objectType.NumField(); i++ {
		var (
			zeroValue         reflect.Value
			field             = objectType.Field(i)
			bsonName          = field.Tag.Get("bson")
			kind              = field.Type.Kind()
			err               error
			newInheritedRules []rule.Rule
		)

		if bsonName == "" || bsonName == "-" {
			continue
		}

		switch kind { //nolint:exhaustive
		case reflect.Ptr:
			zeroValue = reflect.Zero(field.Type.Elem())
			kind = zeroValue.Kind()
		case reflect.Array, reflect.Slice:
			zeroValue = reflect.MakeSlice(field.Type, 0, 0)
		case reflect.Map:
			zeroValue = reflect.MakeMap(field.Type)
			kind = zeroValue.Kind()
		default:
			zeroValue = reflect.Zero(field.Type)
		}

		if newInheritedRules, err = v.addRule(
			path+bsonName, kind, field.Tag, zeroValue.Interface(), inheritedRules,
		); err != nil {
			return err
		}

		if kind == reflect.Struct ||
			kind == reflect.Slice ||
			kind == reflect.Array ||
			kind == reflect.Ptr ||
			kind == reflect.Map {
			if err := v.parseReference(zeroValue.Type(), path+bsonName, newInheritedRules); err != nil {
				return err
			}
		}
	}

	return nil
}

func (v *Validator) addRule(
	path string, kind reflect.Kind, tag reflect.StructTag, instance interface{}, inheritedRules []rule.Rule,
) ([]rule.Rule, error) {
	if kind == reflect.Invalid {
		return nil, InvalidTypeError{path: path}
	}

	var (
		rulesToInherit = []rule.Rule{}
		rules          = []rule.Rule{}
	)

	for _, rule := range v.generalRules {
		newRule, err := rule.NewInstance(path, kind, instance, "")
		if err != nil {
			return nil, fmt.Errorf("could not instantiate rule on '%s': %w", path, err)
		}

		rules = append(rules, newRule)
	}

	if tag != "" {
		var err error

		rulesToInherit, err = v.composeRulesFromTags(path, kind, instance, tag, &rules)
		if err != nil {
			return nil, err
		}
	}

	err := v.composeRulesFromHeredity(path, kind, instance, &rules, inheritedRules)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "*") {
		v.wildcardRules[operation.Path(path)] = rules
	} else {
		v.rules[operation.Path(path)] = rules
	}

	return rulesToInherit, nil
}

func (v *Validator) composeRulesFromTags(
	path string, kind reflect.Kind, instance interface{},
	tag reflect.StructTag, rules *[]rule.Rule,
) ([]rule.Rule, error) {
	rulesToInherit := []rule.Rule{}
	tagsToInherit := []string{}

	if value := tag.Get("jp_inherit"); value != "" {
		split := strings.Split(value, ",")

		for _, tagToInherit := range split {
			if tag.Get(tagToInherit) == "" {
				return nil, InheritNonExistingTagError{name: tagToInherit}
			}

			tagsToInherit = append(tagsToInherit, tagToInherit)
		}
	}

	for name, rule := range v.knownTagRules {
		if name == "jp_inherit" {
			continue
		}

		value, use := tag.Lookup(name)
		if use {
			newRule, err := rule.NewInstance(path, kind, instance, value)
			if err != nil {
				return nil, fmt.Errorf("could not instantiate rule on '%s': %w", path, err)
			}

			*rules = append(*rules, newRule)

			for _, tagToInherit := range tagsToInherit {
				if tagToInherit == name {
					rulesToInherit = append(rulesToInherit, newRule)

					break
				}
			}
		}
	}

	return rulesToInherit, nil
}

func (v *Validator) composeRulesFromHeredity(
	path string, kind reflect.Kind, instance interface{}, rules *[]rule.Rule, inheritedRules []rule.Rule,
) error {
	for _, inheritedRule := range inheritedRules {
		found := false

		for _, rule := range *rules {
			if reflect.TypeOf(inheritedRule) == reflect.TypeOf(rule) {
				found = true

				break
			}
		}

		if !found {
			newInheritedRule, err := inheritedRule.NewInheritInstance(path, kind, instance)
			if err != nil {
				return fmt.Errorf("could not inherit rule on '%s': %w", path, err)
			}

			*rules = append(*rules, newInheritedRule)
		}
	}

	return nil
}

// NewValidator create a new instance of validator using given reference.
func NewValidator(referenceType reflect.Type) (*Validator, error) {
	validator := Validator{
		generalRules: []rule.Rule{
			&rule.MatchingKindRule{},
		},
		knownTagRules: map[string]rule.Rule{
			"jp_disallow":      &rule.DisallowRule{},
			"jp_min":           &rule.MinRule{},
			"jp_max":           &rule.MaxRule{},
			"jp_expression":    &rule.ExpressionRule{},
			"jp_op_allowed":    &rule.AllowedOperationsRule{},
			"jp_op_disallowed": &rule.DisallowedOperationsRule{},
		},
		rules:         map[operation.Path][]rule.Rule{},
		wildcardRules: map[operation.Path][]rule.Rule{},
	}

	err := validator.UseReference(referenceType)
	if err != nil {
		return nil, fmt.Errorf("initializing validator failed: %w", err)
	}

	return &validator, nil
}
