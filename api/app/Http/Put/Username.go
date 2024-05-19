package simple_account_http_put

import (
	"simple_account/app/Auths"
	"simple_account/app/Error"
	"simple_account/app/Http/Message"
	"strconv"
)

func Username(context *Message.Context) (int, error) {
	if context.Author.User == nil {
		return Error.NOT_LOGGED_IN, nil
	}
	author := context.Author.User

	requestUser, errCode, err := context.ResquestBody()
	if errCode != Error.NULL {
		return errCode, err
	}

	if requestUser.Username == "" {
		return Error.USERNAME_IS_EMPTY, nil
	}

	if requestUser.Password == "" {
		return Error.PASSWORD_IS_EMPTY, nil
	}

	newUsername := requestUser.Username
	userId := strconv.FormatInt(int64(author.Id), 10)

	errCode, err = context.Database.UserExist(newUsername, "", userId)
	if err != nil {
		return errCode, err
	}

	connect := context.Database.Connect()
	query := "SELECT salt,hash FROM User WHERE id=?"
	rows, err := connect.Query(query, author.Id)
	if err != nil {
		return Error.SYSTEM, err
	}

	if rows.Next() {
		var salt, hash string
		err = rows.Scan(&salt, &hash)
		if err != nil {
			return Error.SYSTEM, err
		}

		if Auths.Hash(salt, requestUser.Password) != hash {
			return Error.PASSWORD_NOT_MATCH, nil
		}
	} else {
		return Error.USER_NOT_EXIST, nil
	}

	updateQuery := "UPDATE User SET username=?"
	_, err = connect.Exec(updateQuery, newUsername)
	if err != nil {
		return Error.SYSTEM, err
	}

	newUser, errCode, err := context.Database.GetUser(author.Id)
	if errCode != Error.NULL {
		return Error.SYSTEM, err
	}

	errCode, err = context.Author.GenerateToken(&newUser, context.Auth)
	if errCode != Error.NULL {
		return errCode, err
	}

	return Error.NULL, nil
}
