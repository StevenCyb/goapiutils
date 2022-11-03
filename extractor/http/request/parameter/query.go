//nolint:ireturn
package parameter

import (
	"net/http"
	"regexp"
	"strconv"
)

// QueryType define types for `FromQuery`.
type QueryType interface {
	string | int | float64 | bool
}

// FromQuery extracts a query value by key,
// returns default value if not provided by request or throws
// exception if parameter is required.
func FromQuery[T QueryType](req *http.Request, option Option) (T, error) {
	var ret T

	values, exists := req.URL.Query()[option.Key]
	if !exists {
		if option.Required {
			return ret, MissingParameterError{option.Key}
		}

		values = []string{option.Default}
	}

	if option.RegexPattern != "" {
		matched, err := regexp.MatchString(option.RegexPattern, values[0])
		if !matched || err != nil {
			return ret, MalformedParameterError{option.Key}
		}
	}

	switch ptr := any(&ret).(type) {
	case *string:
		*ptr = values[0]

	case *int:
		tmp, err := strconv.ParseInt(values[0], base, bitSize32)
		if err != nil {
			return ret, TypeMismatchError{"int32"}
		}

		*ptr = int(tmp)

	case *float64:
		tmp, err := strconv.ParseFloat(values[0], bitSize64)
		if err != nil {
			return ret, TypeMismatchError{"float64"}
		}

		*ptr = tmp

	case *bool:
		tmp, err := strconv.ParseBool(values[0])
		if err != nil {
			return ret, TypeMismatchError{"bool"}
		}

		*ptr = tmp
	}

	return ret, nil
}
