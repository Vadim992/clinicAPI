package validate

import "regexp"

func ValidatePhoneNum(phoneNum string) bool {
	rePhone := regexp.MustCompile(`^(\+7|8)(( )?(\(\d{3}\)|\d{3}))(( )?\d{3})(( |-)?\d{2}){2}$`)

	if rePhone.MatchString(phoneNum) {
		return true
	}
	
	return false
}