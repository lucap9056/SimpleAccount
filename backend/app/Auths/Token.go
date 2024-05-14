package Auths

import (
	"encoding/base64"
	"encoding/json"
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
	"strings"
	"time"
)

type Signature struct {
	Playload   Playload  `json:"playload"`
	Secret     string    `json:"secret"`
	CreateTime time.Time `json:"-"`
}

type Playload struct {
	User AccountStruct.User `json:"user"`
	Iat  int64              `json:"iat"`
}

func (auth *Auth) GenerateToken(playload Playload, secret string) (string, error) {

	playloadBytes, _ := json.Marshal(playload)

	playloadStr := base64.StdEncoding.EncodeToString(playloadBytes)

	sign := Signature{
		Playload: playload,
		Secret:   secret,
	}

	signJson, _ := json.Marshal(sign)

	signBytes, err := auth.rsaEncode(signJson)
	if err != nil {
		return "", err
	}

	signStr := base64.StdEncoding.EncodeToString(signBytes)

	token := playloadStr + "." + signStr
	return token, nil
}

func (auth *Auth) DecodeToken(jwtStr string) (*Signature, int, error) {
	var err error

	token := strings.Split(jwtStr, ".")
	if len(token) != 2 {
		return nil, Error.AUTHORIZATION_INVALID, nil
	}

	signBase64 := token[1]
	signEncoded, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}
	signDecoded, err := auth.RsaDecode(signEncoded)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}

	var sign Signature
	err = json.Unmarshal(*signDecoded, &sign)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, nil
	}

	return &sign, Error.NULL, nil
}
