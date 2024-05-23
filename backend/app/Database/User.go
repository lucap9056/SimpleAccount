package Database

import (
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
)

func (db *API) UserExist(userName string, userEmail string, filterId string) (int, error) {

	connect := db.Connect()
	include := "username,email"

	query := "SELECT " + include + " FROM User WHERE (username=? OR email=?) AND id!=?"
	rows, err := connect.Query(query, userName, userEmail, filterId)
	if err != nil {
		return Error.SYSTEM, err
	}
	if rows.Next() {
		user := AccountStruct.User{}
		columns := user.MappingTable(include)
		err := rows.Scan(columns...)
		if err != nil {
			return Error.SYSTEM, err
		}
		if user.Username == userName {
			return Error.USERNAME_EXISTED, nil
		}
		return Error.EMAIL_EXISTED, nil
	}

	return Error.NULL, nil
}

func (db *API) GetUser(userId int) (AccountStruct.User, int, error) {
	var user AccountStruct.User
	connect := db.Connect()
	include := "id,username,email,salt,hash,last_edit,register_time,deleted"
	query := "SELECT " + include + " FROM User WHERE id=?"
	rows, err := connect.Query(query, userId)
	if err != nil {
		return user, Error.SYSTEM, err
	}

	if rows.Next() {
		columns := user.MappingTable(include)
		err := rows.Scan(columns...)
		if err != nil {
			return user, Error.SYSTEM, err
		}
		user.MoveTempToFinal()
		return user, Error.NULL, nil
	}

	return user, Error.USER_NOT_EXIST, nil
}
