package ExtensionPost

import (
	"encoding/json"
	"io"
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"strings"
	"time"
)

type Data struct {
	Secret string `json:"secret"`
	Token  string `json:"token"`
}

func GetUser(context *Message.Context) (string, int, error) {
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		return err.Error(), Error.EXTENSION_INVALID_REQUEST, nil
	}

	data := &Data{}
	err = json.Unmarshal(body, data)
	if err != nil {
		return err.Error(), Error.EXTENSION_INVALID_REQUEST, nil
	}

	userId := 0

	auth := context.Auth

	secret := data.Secret
	token := data.Token

	if strings.Contains(token, ".") {
		sign, errCode, _ := auth.DecodeToken(token)
		if errCode != Error.NULL {
			return "", Error.AUTHORIZATION_INVALID, nil
		}

		if secret != sign.Secret {
			return "", Error.AUTHORIZATION_INVALID, nil
		}

		currentTime := time.Now()
		expiresTime := time.Unix(sign.Head.Exp, 0)

		if expiresTime.Before(currentTime) {
			return "", Error.AUTHORIZATION_EXPIRED, nil
		}

		userId = sign.User.Id

	}

	if strings.Index(token, "T ") == 0 {
		tempToken := secret + "-" + token[2:]

		sign := auth.Cache.Verify(tempToken)
		if sign != nil {
			userId = sign.User.Id
		} else {
			return "", Error.AUTHORIZATION_INVALID, nil
		}
	}

	connect := context.Database.Connect()
	include := "id,username"
	query := "SELECT " + include + " FROM User WHERE id=?"
	rows, err := connect.Query(query, userId)
	if err != nil {
		return err.Error(), Error.EXTENSION_SQL_ERROR, nil
	}

	if !rows.Next() {
		return "", Error.EXTENSION_USER_NOT_EXIST, nil
	}

	var user AccountStruct.User
	columns := user.MappingTable(include)
	err = rows.Scan(columns...)
	if err != nil {
		return err.Error(), Error.EXTENSION_SYSTEM_ERROR, nil
	}
	user.MoveTempToFinal()

	userData, err := json.Marshal(user)
	if err != nil {
		return err.Error(), Error.EXTENSION_SYSTEM_ERROR, nil
	}

	return string(userData), Error.NULL, nil
}
