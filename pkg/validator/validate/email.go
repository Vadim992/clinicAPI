package validate

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

func ValidateEmail(email string) bool {
	reEmail := regexp.MustCompile(`(?i)(^([a-z0-9_.-]){1,64})(@([a-z0-9_.-]+(\.(ru|com))))$`)

	if reEmail.MatchString(email) {

		emailArr := strings.Split(email, "@")

		for _, val := range emailArr {

			if !validateNameAndDomain(val) {
				return false
			}

		}

		if utf8.RuneCountInString(email) > 320 {
			return false
		}

		return true

	}
	return false
}

func validateNameAndDomain(str string) bool {

	// func ValidateEmail guarantees that str is not empty
	chars := map[byte]struct{}{
		'.': {},
		'-': {},
		'_': {},
	}

	_, okStart := chars[str[0]]
	_, okEnd := chars[str[len(str)-1]]

	if okStart || okEnd {
		return false
	}
	return true
}
