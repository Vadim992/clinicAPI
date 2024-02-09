package validate

func ValidateId(id int) bool {

	if id < 1 {
		return false
	}

	return true
}