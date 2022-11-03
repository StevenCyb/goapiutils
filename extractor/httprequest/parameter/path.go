//nolint:ireturn
package parameter

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
)

// PathType define types for `FromPath`.
type PathType interface {
	string | int | bool
}

// FromQuery extracts a path value by key,
// returns default value if not provided by request or throws
// exception if parameter is required.
func FromPath[T PathType](req *http.Request, option Option) (T, error) {
	var ret T

	value, exists := mux.Vars(req)[option.Key]
	if !exists {
		if option.Required {
			return ret, MissingParameterError{option.Key}
		}

		value = option.Default
	}

	if option.RegexPattern != "" {
		matched, err := regexp.MatchString(option.RegexPattern, value)
		if !matched || err != nil {
			return ret, MalformedParameterError{option.Key}
		}
	}

	switch ptr := any(&ret).(type) {
	case *string:
		*ptr = value

	case *int:
		tmp, err := strconv.ParseInt(value, base, bitSize32)
		if err != nil {
			return ret, TypeMismatchError{ExpectedType: "int32"}
		}

		*ptr = int(tmp)

	case *bool:
		tmp, err := strconv.ParseBool(value)
		if err != nil {
			return ret, TypeMismatchError{ExpectedType: "bool"}
		}

		*ptr = tmp
	}

	return ret, nil
}
