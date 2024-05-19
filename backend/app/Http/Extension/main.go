package Extension

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"strings"
	"time"
)

type Data struct {
	secret string `json:"secret"`
	token  string `json:"token"`
}

func GetUser(context *Message.Context) {
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &Data{}
	err = json.Unmarshal(body, data)
	if err != nil {
		http.Error(context.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	userId := 0

	auth := context.Auth

	secret := data.secret
	token := data.token

	if strings.Contains(token, ".") {
		sign, errCode, _ := auth.DecodeToken(token)
		if errCode != Error.NULL {
			http.Error(context.Writer, fmt.Sprintf("%d", errCode), http.StatusInternalServerError)
			return
		}

		if secret != sign.Secret {
			http.Error(context.Writer, fmt.Sprintf("%d", errCode), http.StatusInternalServerError)
			return
		}

		currentTime := time.Now()
		expiresTime := time.Unix(sign.Head.Exp, 0)

		if expiresTime.Before(currentTime) {
			http.Error(context.Writer, fmt.Sprintf("%d", errCode), http.StatusInternalServerError)
			return
		}

		userId = sign.User.Id

	}

	if strings.Index(token, "T ") == 0 {
		tempToken := secret + "-" + token[2:]

		sign := auth.Cache.Verify(tempToken)
		if sign != nil {
			userId = sign.User.Id
		} else {
			return
		}
	}

	w := context.Writer

	w.WriteHeader(http.StatusOK)
	userIdStr := fmt.Sprintf("%d", userId)
	w.Write([]byte(userIdStr))
}
