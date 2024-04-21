package customer

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
)

func GetTokenFromHeader(header string) (string, error) {
	slicedHeader := strings.Split(header, " ")

	switch {
	case len(slicedHeader) != 2:
		return "", errors.New("invalid header")
	case strings.ToLower(slicedHeader[0]) != "bearer":
		return "", errors.New("missing Bearer")
	case len(slicedHeader[1]) != TOKEN_LENGTH:
		return "", errors.New("invalid token")
	}

	return slicedHeader[1], nil
}

func GenerateToken() string {
	tokenBytes := make([]byte, TOKEN_LENGTH/2)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return ""
	}

	token := hex.EncodeToString(tokenBytes)
	return token
}