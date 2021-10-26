package variableutils

import "regexp"

// IsValidKey checks if a key is valid if it follows the following criteria:
// must have no more than 255 characters;
// only A-Z, a-z, 0-9, and _ are allowed
func IsValidKey(key string) bool {
	// check if key falls within range of 1-255
	if len(key) > 255 || len(key) < 1 {
		return false
	}
	keyRE := regexp.MustCompile(`^[A-Za-z0-9_]+$`)
	return keyRE.MatchString(key)
}

var ValidKeyMsg = "A valid key must have no more than 255 characters; only A-Z, a-z, 0-9, and _ are allowed"
