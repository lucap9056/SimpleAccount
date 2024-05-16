package Author

import (
	"net/http"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Error"
	"strings"
	"time"
)

type Author struct {
	secret      string
	token       string
	expiresTime time.Time
	User        *AccountStruct.User
	writer      http.ResponseWriter
	Language    string
}

func New(writer http.ResponseWriter) *Author {
	author := Author{
		writer: writer,
	}
	return &author
}

func getSecret(request *http.Request) string {
	cookieSecret, err := request.Cookie("secret")
	if err != nil {
		return ""
	}
	return cookieSecret.Value
}

func (author *Author) SetSecret(secret string, expiresTime time.Time) {
	author.secret = secret
	author.expiresTime = expiresTime
}

func (author *Author) SetToken(token string) {
	author.token = token
}

func (author *Author) InvaildToken() {
	author.token = "invalid"
}

func (author *Author) InvaildTempToken() {
	author.token = "invalid_t"
}

func (author *Author) UpdateToken() {

	if author.token == "invalid" {
		http.SetCookie(author.writer, &http.Cookie{
			Name:   "secret",
			MaxAge: -1,
			Path:   "/",
		})
		author.secret = ""
	}

	if author.secret != "" {
		cookie := http.Cookie{
			Name:     "secret",
			Value:    author.secret,
			HttpOnly: true,
			Expires:  author.expiresTime,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
		}
		http.SetCookie(author.writer, &cookie)
	}

	if author.token != "" {
		author.writer.Header().Add("Authorization", author.token)
	}

}

func (author *Author) setUser(user *AccountStruct.User) *Author {
	author.User = user
	return author
}

func Get(auth *Auths.Auth, writer http.ResponseWriter, request *http.Request) *Author {
	author := &Author{
		secret: "",
		token:  "",
		writer: writer,
	}

	author.Language = request.Header.Get("Accept-Language")

	secret := getSecret(request)
	token := request.Header.Get("Authorization")

	if strings.Contains(token, ".") {
		sign, _, err := auth.DecodeToken(token)
		if err != nil {
			author.InvaildToken()
			return author
		}

		if secret != sign.Secret {
			author.InvaildToken()
			return author
		}

		currentTime := time.Now()
		expiresTime := time.Unix(sign.Playload.Iat, 0)

		if expiresTime.Before(currentTime) {
			author.InvaildToken()
			return author
		}

		if expiresTime.Before(currentTime.Add(auth.RenewTime)) {
			expiresTime = currentTime.Add(auth.ValidityDuration)

			sign.Playload.Iat = expiresTime.Unix()
			secret := Auths.Salt()

			token, err := auth.GenerateToken(sign.Playload, secret)
			if err != nil {
				return author
			}

			author.SetSecret(secret, expiresTime)
			author.SetToken(token)
		} else {
			tempToken, err := auth.Cache.GenerateToken()
			if err != nil {
				return author
			}

			sign.CreateTime = currentTime
			err = auth.Cache.Add(secret+"-"+tempToken, *sign)
			if err != nil {
				return author
			}

			author.SetToken("Bearer " + tempToken)
		}

		user := &sign.Playload.User
		return author.setUser(user)
	}

	if strings.Index(token, "Bearer ") == 0 {
		tempToken := secret + "-" + token[7:]

		sign := auth.Cache.Verify(tempToken)
		if sign != nil {
			user := &sign.Playload.User
			return author.setUser(user)
		} else {
			author.InvaildTempToken()
		}
	}
	return author
}

func (author *Author) GenerateToken(user *AccountStruct.User, auth *Auths.Auth) (int, error) {
	expiresTime := time.Now().Add(auth.ValidityDuration)

	playload := Auths.Playload{
		User: AccountStruct.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
		},
		Iat: expiresTime.Unix(),
	}
	secret := Auths.Salt()
	token, err := auth.GenerateToken(playload, secret)
	if err != nil {
		return Error.SYSTEM, err
	}

	author.SetSecret(secret, expiresTime)
	author.SetToken(token)
	return Error.NULL, nil
}
