package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrAuthHeaderDoesNotExist = errors.New("authorization header not found")
var ErrMalformedAuthHeader = errors.New("authorization header is malformed")

func GetBearerToken(header http.Header) (string, error) {
	bearerToken := header.Get("Authorization")
	if bearerToken == "" {
		return "", ErrAuthHeaderDoesNotExist
	}
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) < 2 || splitToken[0] != "Bearer" {
		return "", ErrMalformedAuthHeader
	}

	// splitToken[1] will contain the actual bearer token
	return splitToken[1], nil
}