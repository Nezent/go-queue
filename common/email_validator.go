package common

import "regexp"

// ValidateEmailWithRegex validates if the email has a proper structure and ends with 'ac.bd' or 'edu.bd'
func ValidateEmailWithRegex(email string) bool {
	// Regex to match email with domains ending in ac.bd or edu.bd
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.(ac\.bd|edu\.bd)$`

	// Compile the regex
	re := regexp.MustCompile(emailRegex)

	// Validate the email against the regex
	return re.MatchString(email)
}
