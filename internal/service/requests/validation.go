package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

func isHexAddress(value interface{}) error {
	return validation.Validate(value, validation.Match(regexp.MustCompile("^0x[0-9a-fA-F]{40}$")))
}

func isHash(value interface{}) error {
	return validation.Validate(value, validation.Match(regexp.MustCompile("^0x[0-9a-fA-F]{64}$")))
}
