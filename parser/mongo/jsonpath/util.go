package jsonpath

import "regexp"

type Path string

// Valid check if given path is in valid format.
func (p Path) Valid() bool {
	regex := regexp.MustCompile(`^(([\w-]+)+\.?)*([\w-]+)+$`)

	return regex.MatchString(string(p))
}
