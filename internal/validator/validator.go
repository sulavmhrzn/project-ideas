package validator

import "regexp"

type Validator struct {
	Errors map[string]any
}

var EmailPatter = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func New() *Validator {
	return &Validator{
		Errors: make(map[string]any),
	}
}

func (v *Validator) AddError(key, message string) {
	if _, found := v.Errors[key]; !found {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func ValidEmail(email string) bool {
	return EmailPatter.MatchString(email)
}

func Unique(value ...string) bool {
	uniqueValues := make(map[string]bool)
	for _, v := range value {
		if _, ok := uniqueValues[v]; !ok {
			uniqueValues[v] = true
		}
	}
	return len(uniqueValues) == len(value)
}
