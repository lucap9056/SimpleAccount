package simple_account_http_post

import (
	"simple_account/app/AccountStruct"
	"simple_account/app/Auths"
	"simple_account/app/Error"
	"simple_account/app/Http/Author"
	"simple_account/app/Http/Message"
)

func Login(context *Message.Context) (int, error) {
	requestUser, errCode, err := context.ResquestBody()
	if errCode != Error.NULL {
		return errCode, err
	}

	if requestUser.Email == "" {
		return Error.EMAIL_IS_EMPTY, nil
	}

	if requestUser.Password == "" {
		return Error.PASSWORD_IS_EMPTY, nil
	}

	connect := context.Database.Connect()
	include := "id,username,email,salt,hash,last_edit,register_time,deleted"
	query := "SELECT " + include + " FROM User WHERE email=?"
	rows, err := connect.Query(query, requestUser.Email)
	if err != nil {
		return Error.SYSTEM, err
	}

	user := &AccountStruct.User{}
	if rows.Next() {
		columns := user.MappingTable(include)
		err := rows.Scan(columns...)
		if err != nil {
			return Error.SYSTEM, err
		}
		user.MoveTempToFinal()

		if Auths.Hash(user.Salt, requestUser.Password) != user.Hash {
			return Error.PASSWORD_NOT_MATCH, nil
		}

		if user.DeletedTime != nil {
			return Error.USER_DELETED, nil
		}
	} else {
		return Error.USER_NOT_EXIST, nil
	}

	context.Author = Author.New(context.Writer)
	errCode, err = context.Author.GenerateToken(user, context.Auth)
	if errCode != Error.NULL {
		return errCode, err
	}

	return Error.NULL, nil
}
