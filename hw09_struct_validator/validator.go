package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrNotStruct        = errors.New("value must be a struct or a pointer to a struct")
	ErrInvalidTag       = errors.New("invalid validate tag")
	ErrUnknownValidator = errors.New("unknown validator")
	ErrUnsupportedType  = errors.New("unsupported field type")

	ErrStringLength = errors.New("string length is invalid")
	ErrStringRegexp = errors.New("string does not match regexp")
	ErrStringIn     = errors.New("string is not in the allowed set")

	ErrIntMin = errors.New("number is less than minimum")
	ErrIntMax = errors.New("number is greater than maximum")
	ErrIntIn  = errors.New("number is not in the allowed set")
)

type ValidationError struct {
	Field string
	Err   error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %v", e.Field, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	parts := make([]string, 0, len(v))
	for _, err := range v {
		parts = append(parts, err.Error())
	}

	return strings.Join(parts, "; ")
}

func (v ValidationErrors) Unwrap() []error {
	result := make([]error, 0, len(v))
	for _, err := range v {
		result = append(result, err)
	}

	return result
}

// Is supports errors.Is for the underlying validation errors.
func (v ValidationErrors) Is(target error) bool {
	for _, err := range v {
		if errors.Is(err.Err, target) {
			return true
		}
	}

	return false
}

// Validate validates exported struct fields using the "validate" tag.
// All validation errors are accumulated and returned as ValidationErrors.
func Validate(v interface{}) error {
	if v == nil {
		return ErrNotStruct
	}

	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return ErrNotStruct
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	validationErrors := validateStruct(value)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func validateStruct(value reflect.Value) ValidationErrors {
	result := make(ValidationErrors, 0)
	valueType := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)
		if field.PkgPath != "" { // unexported field
			continue
		}

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		fieldValue := value.Field(i)
		if tag == "nested" {
			result = append(result, validateNested(field.Name, fieldValue)...)
			continue
		}

		fieldErrors := validateField(field.Name, fieldValue, tag)
		result = append(result, fieldErrors...)
	}

	return result
}

func validateNested(fieldName string, value reflect.Value) ValidationErrors {
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return ValidationErrors{{
			Field: fieldName,
			Err:   fmt.Errorf("%w: nested field must be a struct", ErrUnsupportedType),
		}}
	}

	result := make(ValidationErrors, 0)
	for _, err := range validateStruct(value) {
		err.Field = fieldName + "." + err.Field
		result = append(result, err)
	}

	return result
}

func validateField(fieldName string, value reflect.Value, tag string) ValidationErrors {
	switch value.Kind() { //nolint:exhaustive // only the required types are supported
	case reflect.String:
		if err := validateString(value.String(), tag); err != nil {
			return ValidationErrors{{Field: fieldName, Err: err}}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if err := validateInt(value.Int(), tag); err != nil {
			return ValidationErrors{{Field: fieldName, Err: err}}
		}
	case reflect.Slice:
		return validateSlice(fieldName, value, tag)
	default:
		return ValidationErrors{{
			Field: fieldName,
			Err:   fmt.Errorf("%w: %s", ErrUnsupportedType, value.Kind()),
		}}
	}

	return nil
}

func validateSlice(fieldName string, value reflect.Value, tag string) ValidationErrors {
	result := make(ValidationErrors, 0)

	for i := 0; i < value.Len(); i++ {
		elementName := fmt.Sprintf("%s[%d]", fieldName, i)
		element := value.Index(i)

		var err error
		switch element.Kind() { //nolint:exhaustive // only the required types are supported
		case reflect.String:
			err = validateString(element.String(), tag)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err = validateInt(element.Int(), tag)
		default:
			err = fmt.Errorf("%w: %s", ErrUnsupportedType, element.Kind())
		}

		if err != nil {
			result = append(result, ValidationError{Field: elementName, Err: err})
		}
	}

	return result
}

func validateString(value, tag string) error {
	for _, rule := range strings.Split(tag, "|") {
		name, argument, ok := strings.Cut(rule, ":")
		if !ok {
			return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
		}

		switch name {
		case "len":
			expected, err := strconv.Atoi(argument)
			if err != nil || expected < 0 {
				return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
			}
			if utf8.RuneCountInString(value) != expected {
				return fmt.Errorf("%w: expected %d, got %d", ErrStringLength, expected, utf8.RuneCountInString(value))
			}
		case "regexp":
			re, err := regexp.Compile(argument)
			if err != nil {
				return fmt.Errorf("%w: %q: %v", ErrInvalidTag, rule, err)
			}
			if !re.MatchString(value) {
				return fmt.Errorf("%w: %q", ErrStringRegexp, value)
			}
		case "in":
			if !contains(strings.Split(argument, ","), value) {
				return fmt.Errorf("%w: %q", ErrStringIn, value)
			}
		default:
			return fmt.Errorf("%w: %q", ErrUnknownValidator, name)
		}
	}

	return nil
}

func validateInt(value int64, tag string) error {
	for _, rule := range strings.Split(tag, "|") {
		name, argument, ok := strings.Cut(rule, ":")
		if !ok {
			return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
		}

		switch name {
		case "min":
			expected, err := strconv.ParseInt(argument, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
			}
			if value < expected {
				return fmt.Errorf("%w: expected at least %d, got %d", ErrIntMin, expected, value)
			}
		case "max":
			expected, err := strconv.ParseInt(argument, 10, 64)
			if err != nil {
				return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
			}
			if value > expected {
				return fmt.Errorf("%w: expected at most %d, got %d", ErrIntMax, expected, value)
			}
		case "in":
			values, err := parseInts(argument)
			if err != nil {
				return fmt.Errorf("%w: %q", ErrInvalidTag, rule)
			}
			if !containsInt(values, value) {
				return fmt.Errorf("%w: %d", ErrIntIn, value)
			}
		default:
			return fmt.Errorf("%w: %q", ErrUnknownValidator, name)
		}
	}

	return nil
}

func parseInts(argument string) ([]int64, error) {
	parts := strings.Split(argument, ",")
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		value, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}

	return result, nil
}

func contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}

	return false
}

func containsInt(values []int64, expected int64) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}

	return false
}
