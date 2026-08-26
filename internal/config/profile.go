package config

import "reflect"

type Validator interface {
	Validate(int) error
}

type RangeValidator struct {
	Minimum int
}

func (v *RangeValidator) Validate(value int) error {
	if value < v.Minimum {
		return ErrBelowRange
	}
	return nil
}

type validationError string

func (e validationError) Error() string { return string(e) }

const ErrBelowRange validationError = "reading below range"

type Profile struct {
	Labels    map[string]string
	Validator Validator
}

func LoadDefault() Profile {
	return Profile{Labels: make(map[string]string), Validator: nil}
}

func ValidatorAvailable(validator Validator) bool {
	if validator == nil {
		return false
	}
	value := reflect.ValueOf(validator)
	return value.Kind() != reflect.Pointer || !value.IsNil()
}
