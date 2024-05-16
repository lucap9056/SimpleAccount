package Database

import (
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
)

func (db *API) UserExist(userName string, userEmail string, filterId string) (int, error) {

	connect := db.Connect()
	rows, err := connect.Query("SELECT username,email FROM User WHERE (username=? OR email=?) AND id!=?", userName, userEmail, filterId)
	if err != nil {
		return Error.SYSTEM, err
	}
	for rows.Next() {
		var name string
		var email string
		err := rows.Scan(&name, &email)
		if err != nil {
			return Error.SYSTEM, err
		}
		if name == userName {
			return Error.USERNAME_EXISTED, nil
		}
		if email == userEmail {
			return Error.EMAIL_EXISTED, nil
		}
	}

	return Error.NULL, nil
}

func (db *API) GetUser(userId int) (AccountStruct.User, int, error) {
	var user AccountStruct.User
	connect := db.Connect()
	query := "SELECT id,username,email,salt,hash FROM User WHERE id=?"
	rows, err := connect.Query(query, userId)
	if err != nil {
		return user, Error.SYSTEM, err
	}

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Salt, &user.Hash)
		if err != nil {
			return user, Error.SYSTEM, err
		}
		return user, Error.NULL, nil
	}

	return user, Error.USER_NOT_EXIST, nil
}
