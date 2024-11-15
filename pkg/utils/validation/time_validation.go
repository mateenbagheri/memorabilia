package validation

import (
	"errors"
	"regexp"
)

var ErrInvalidJonTimeFormat = errors.New("invalid time format, must match (number)h(number)m(number)s pattern")

func ValidateJobTimeFormat(timeStr string) error {
	jobTimePattern := `^(?:(?:[1-9]\d*)h)?(?:(?:[1-9]\d*)m)?(?:(?:[1-9]\d*)s)?$`

	if timeStr == "" {
		return ErrInvalidJonTimeFormat
	}

	re := regexp.MustCompile(jobTimePattern)
	if !re.MatchString(timeStr) {
		return ErrInvalidJonTimeFormat
	}

	return nil
}
