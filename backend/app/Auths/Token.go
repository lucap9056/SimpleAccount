package Auths

import (
	"encoding/base64"
	"encoding/json"
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
	"strings"
	"time"
)

type Head struct {
	KeyID string `json:"key_id"`
	Exp   int64  `json:"exp"`
}
type Signature struct {
	Head       Head               `json:"head"`
	User       AccountStruct.User `json:"playload"`
	Secret     string             `json:"secret"`
	CreateTime time.Time          `json:"-"`
}

func (auth *Auth) GenerateToken(user AccountStruct.User, secret string) (string, error) {
	keyId := auth.rsaRandomKeyId()
	expiresTime := time.Now().Add(auth.ValidityDuration)
	head := Head{
		KeyID: keyId,
		Exp:   expiresTime.Unix(),
	}

	headBytes, err := json.Marshal(head)
	if err != nil {
		return "", err
	}

	headStr := base64.StdEncoding.EncodeToString(headBytes)

	playloadBytes, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	playloadStr := base64.StdEncoding.EncodeToString(playloadBytes)

	sign := Signature{
		Head: head,
		User: AccountStruct.User{
			Id:           user.Id,
			LastEditTime: user.LastEditTime,
		},
		Secret: secret,
	}

	signJson, _ := json.Marshal(sign)

	signBytes, err := auth.rsaEncode(keyId, signJson)
	if err != nil {
		return "", err
	}

	signStr := base64.StdEncoding.EncodeToString(signBytes)

	token := headStr + "." + playloadStr + "." + signStr
	return token, nil
}

func (auth *Auth) DecodeToken(jwtStr string) (*Signature, int, error) {
	var err error

	token := strings.Split(jwtStr, ".")
	if len(token) != 3 {
		return nil, Error.AUTHORIZATION_INVALID, nil
	}

	headBase64 := token[0]
	headEncoded, err := base64.StdEncoding.DecodeString(headBase64)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}
	var head Head
	err = json.Unmarshal(headEncoded, &head)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}

	currentTime := time.Now().Unix()
	if currentTime > head.Exp {
		return nil, Error.AUTHORIZATION_INVALID, nil
	}

	signBase64 := token[2]
	signEncoded, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, err
	}
	signDecoded, errCode, err := auth.RsaDecode(head.KeyID, signEncoded)
	if err != nil {
		return nil, errCode, err
	}

	var sign Signature
	err = json.Unmarshal(*signDecoded, &sign)
	if err != nil {
		return nil, Error.AUTHORIZATION_INVALID, nil
	}

	return &sign, Error.NULL, nil
}
