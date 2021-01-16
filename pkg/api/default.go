package api

var DefaultListLimit = 30

//IsValidToken checks if a token provided is valid.
// The token string must be 20 characters in length to be recognized as
// a valid personal access token.
//
// TODO: check if token has minimum scopes required by glab
func IsValidToken(token string) bool {
	return len(token) == 20
}
