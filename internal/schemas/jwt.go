package schemas

import (
	"errors"
	"time"
)

type JWTTokenHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type JWTTokenBody struct{
	Organizations   []int `json:"organizations"`
	UserID          int      `json:"user_id"`
	Exp             int      `json:"exp"`
	Roles           []int    `json:"roles"`
	Iss             string   `json:"iss"`
	Sig             string   `json:"sig"`
	ExpireTimedelta int      `json:"expire_timedelta"`
	UserName        string   `json:"user_name"`
	Email           string   `json:"email"`
	DeviceID        int      `json:"device_id"`
}

func (J JWTTokenBody) Valid() error {
	if time.Now().After(time.Unix(int64(J.Exp),0)) {
		return errors.New("token has expired")
	}
	return nil
}

