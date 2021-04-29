package dal

import "regexp"

type InputValidator struct {
	count            int
	validationSteps  []func(string) (string, bool)
	params           []interface{}
	alphaNumericRgxp *regexp.Regexp
}

func NewValidator() *InputValidator {
	r, _ := regexp.Compile("^[a-zA-Z0-9_-]*$")
	validator := &InputValidator{
		alphaNumericRgxp: r,
	}
	validator.validationSteps = []func(string) (string, bool){
		validator.isAlphaNumerical(),
		validator.replaceOperator(),
		validator.validateOperand(),
	}

	return validator
}

func (i *InputValidator) Validate(s string) (string, bool) {
	s, ok := i.validationSteps[i.count](s)
	i.count += 1
	i.count %= len(i.validationSteps)
	return s, ok
}

func (i *InputValidator) validateOperand() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !i.alphaNumericRgxp.MatchString(s) {
			return s, false
		}
		i.params = append(i.params, s)
		return "?", true
	}
}

func (i *InputValidator) replaceOperator() func(s string) (string, bool) {
	replaceMap := map[string]string{
		"eq": "=",
		"ne": "!=",
		"gt": ">",
		"lt": "<",
	}

	return func(s string) (string, bool) {
		replaceWith, ok := replaceMap[s]
		if !ok {
			return s, false
		}
		return replaceWith, true
	}
}

func (i *InputValidator) isAlphaNumerical() func(string) (string, bool) {
	return func(s string) (string, bool) {
		if !i.alphaNumericRgxp.MatchString(s) {
			return s, false
		}
		return s, true
	}
}
