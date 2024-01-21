package notifications

import (
	"errors"
	"regexp"
)

type Email struct {
	Content string
	To      string
	From    string
}

type Emailer interface {
	SendEmail(mail Email) error
}

func ValidateEmailFormat(email string) error {
	validEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !validEmail.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}
