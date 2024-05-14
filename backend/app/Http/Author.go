package simple_account_http

import (
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"strings"
	"time"
)

func GetAuthor(auth *Auths.Auth, writer http.ResponseWriter, req *http.Request) *AccountStruct.User {

	secretCookie, err := req.Cookie("secret")
	if err != nil {
		return nil
	}
	secret := secretCookie.Value

	token := req.Header.Get("Authorization")
	//是否為完整Token
	if strings.Contains(token, ".") {
		//嘗試從Token取出user資料
		sign, _, err := auth.DecodeToken(token)
		if err != nil || sign == nil {
			invaildToken(writer)
			return nil
		}
		//驗證
		if secret != sign.Secret {
			invaildToken(writer)
			return nil
		}

		currentTime := time.Now()
		expiresTime := time.Unix(sign.Playload.Iat, 0)
		//如果過期 通知client
		if expiresTime.Before(currentTime) {
			invaildToken(writer)
			return nil
		}

		//是否需要更新
		if expiresTime.Before(currentTime.Add(auth.RenewTime)) {

			expiresTime = currentTime.Add(auth.ValidityDuration)

			sign.Playload.Iat = expiresTime.Unix()
			secret := Auths.Salt()
			token, err := auth.GenerateToken(sign.Playload, secret)
			if err != nil {
				return nil
			}

			cookie := http.Cookie{
				Name:     "secret",
				Value:    secret,
				HttpOnly: true,
				Expires:  expiresTime,
				SameSite: http.SameSiteLaxMode,
				Path:     "/",
			}

			http.SetCookie(writer, &cookie)
			writer.Header().Add("Authorization", token)
		} else {
			//生成臨時Token
			tempKey, err := auth.Cache.GenerateKey()
			if err != nil {
				return nil
			}

			sign.CreateTime = currentTime
			err = auth.Cache.Add(secret+"-"+tempKey, *sign)
			if err != nil {
				return nil
			}

			writer.Header().Add("Authorization", "Bearer "+tempKey)
		}
		user := sign.Playload.User
		return &user
	}
	//持有臨時Token
	if strings.Index(token, "Bearer ") == 0 {

		cacheToken := secret + "-" + token[7:]

		sign := auth.Cache.Verify(cacheToken)
		if sign != nil {
			return &sign.Playload.User
		}
		//通知 client 臨時Token無效
		writer.Header().Add("Authorization", "invalid_t")
	}
	return nil
}

func invaildToken(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:   "secret",
		MaxAge: -1,
		Path:   "/",
	})
	writer.Header().Add("Authorization", "invalid")
}
