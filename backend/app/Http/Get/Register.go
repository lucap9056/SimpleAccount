package simple_account_http_get

import (
	"simple_account/app/Database"
	"simple_account/app/Error"
)

func RegisterVerify(db *Database.API, verifyKey string) (int, error) {
	userData := db.Verification.Verify(verifyKey)
	if userData == nil {
		return Error.REGISTER_VERIFY_USER_NOT_EXIST, nil
	}

	connect := db.Connect()
	query := "INSERT INTO User(username, email, salt, hash) VALUE(?, ?, ?, ?)"
	stmt, err := connect.Prepare(query)
	if err != nil {
		return Error.SYSTEM, err
	}
	_, err = stmt.Exec(userData.Username, userData.Email, userData.Salt, userData.Hash)
	if err != nil {
		return Error.SYSTEM, err
	}

	db.Verification.Del(verifyKey)

	return Error.NULL, nil
}
