package clientvalid

import "regexp"

func ValidPathAll(pathUrl string) bool {
	rePathAll := regexp.MustCompile(`^client$`)

	if !rePathAll.MatchString(pathUrl) {
		return false
	}
	return true
}

func ValidPathId(pathUrl string) bool {
	rePathId := regexp.MustCompile(`^client/[0-9]+$`)

	if !rePathId.MatchString(pathUrl) {
		return false
	}
	return true
}
