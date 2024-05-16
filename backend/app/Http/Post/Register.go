package simple_account_http_post

import (
	"fmt"
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Auths/Validate"
	"simple_account/app/Email/TimedKeys"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"time"
)

func Register(context *Message.Context) (int, error) {
	database := context.Database
	email := context.Email

	userData, errCode, err := context.ResquestBody()
	if errCode != Error.NULL {
		return errCode, err
	}

	errCode = Validate.Username(userData.Username)
	if errCode != Error.NULL {
		return errCode, nil
	}

	errCode = Validate.Email(userData.Email)
	if errCode != Error.NULL {
		return errCode, nil
	}

	errCode = Validate.Password(userData.Password)
	if errCode != Error.NULL {
		return errCode, nil
	}

	errCode, err = database.UserExist(userData.Username, userData.Email, "")
	if errCode != Error.NULL {
		return errCode, err
	}

	userData.Salt = Auths.Salt()
	userData.Hash = Auths.Hash(userData.Salt, userData.Password)

	currentTime := time.Now()
	user := AccountStruct.User{
		Username:   userData.Username,
		Email:      userData.Email,
		Salt:       userData.Salt,
		Hash:       userData.Hash,
		CreateTime: &currentTime,
	}

	key, err := TimedKeys.GenerateKey()
	if err != nil {
		return Error.SYSTEM, err
	}

	language := context.Author.Language
	fmt.Println(language)
	err = email.Templates.Register.SendVerify(language, user.Email, key)
	if err != nil {
		return Error.SYSTEM, err
	}
	err = email.TimedKeys.Add("register", key, user)
	if err != nil {
		return Error.SYSTEM, err
	}

	return Error.NULL, nil
}
