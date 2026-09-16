package service

import (
	"strings"
	"unicode"

	"go_module5/app/model"
)

const minPasswordLength = 8

func ValidateRegister(req model.RegisterRequest) map[string]string {
	errs := map[string]string{}

	nim := strings.TrimSpace(req.NIM)
	switch {
	case nim == "":
		errs["nim"] = "must be filled"
	case len(nim) < 3:
		errs["nim"] = "must be a minimum of 3 characters"
	case !isValidNIM(nim):
		errs["nim"] = "can only include letters and numbers"
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		errs["name"] = "must be filled"
	}

	if req.Grade < 0 || req.Grade > 4.0 {
		errs["grade"] = "must be between 0.0 and 4.0"
	}

	if msg := checkPasswordStrength(req.Password); msg != "" {
		errs["password"] = msg
	}

	return errs
}

func ValidateLogin(req model.LoginRequest) map[string]string {
	errs := map[string]string{}

	if strings.TrimSpace(req.NIM) == "" {
		errs["nim"] = "must be filled"
	}
	if req.Password == "" {
		errs["password"] = "must be filled"
	}

	return errs
}

func checkPasswordStrength(password string) string {
	if len(password) < minPasswordLength {
		return "must be a minimum of 8 characters"
	}

	var hasLetter, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}

	if !hasLetter || !hasDigit {
		return "must contain letters and numbers"
	}

	weak := map[string]bool{
		"password1": true, "12345678": true, "qwerty123": true,
		"admin123": true, "password123": true,
	}
	if weak[strings.ToLower(password)] {
		return "password too common"
	}

	return ""
}

func isValidNIM(nim string) bool {
	for _, r := range nim {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
