package jsonpatch

import (
	"regexp"
	"strings"
)

type Path string

// Valid check if given path is in valid format.
func (p Path) Valid() bool {
	regex := regexp.MustCompile(`^(([\w-]+)+\.?)*([\w-]+)+$`)

	return regex.MatchString(string(p))
}

// Compare check if path is equal to given path.
// Single fields can be set to `*` for wildcard.
func (p Path) Compare(p2 Path) bool {
	if strings.Contains(string(p), "*") {
		stringPath := []rune(p)
		offset := 0

		for index, char := range p2 {
			if index == len(p2)-1 && (index+offset) == len(stringPath)-1 {
				return true
			}

			if len(stringPath) <= index+offset {
				break
			}

			if stringPath[index+offset] == '*' {
				if char != '.' {
					offset--
				} else {
					offset++
				}

				continue
			}

			if stringPath[index+offset] != char {
				break
			}
		}

		return false
	}

	if p == p2 {
		return true
	}

	return false
}
