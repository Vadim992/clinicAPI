package docrotvalid

import "regexp"

func ValidPathAll(pathUrl string) bool {
	rePathAll := regexp.MustCompile(`^doctor$`)

	if !rePathAll.MatchString(pathUrl) {
		return false
	}
	return true
}

func ValidPathId(pathUrl string) bool {
	rePathId := regexp.MustCompile(`^doctor/[0-9]+$`)

	if !rePathId.MatchString(pathUrl) {
		return false
	}
	return true
}
