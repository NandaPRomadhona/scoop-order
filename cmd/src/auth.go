package src

import (
	"encoding/base64"
	"encoding/json"
	"scoop-order/internal/schemas"
	"strings"
)

func DecodeToken(token string) (schemas.JWTTokenBody, error) {
	//Split token part
	var jwtBody schemas.JWTTokenBody
	parts := strings.Split(token, ".")

	tokBody, _ := base64.RawURLEncoding.DecodeString(parts[1])
	err := json.Unmarshal(tokBody, &jwtBody)
	if err != nil {
		return jwtBody, err
	}

	return jwtBody, err
}
