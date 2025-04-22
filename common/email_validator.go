package common

import "regexp"

func ValidateEmailWithRegex(email string) bool {

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Validate the email against the regex
	return re.MatchString(email)
}
