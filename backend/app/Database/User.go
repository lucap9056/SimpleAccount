package Database

import "simple_account/app/Error"

func (db *API) UserExist(userName string, userEmail string) (int, error) {

	connect := db.Connect()
	rows, err := connect.Query("SELECT username,email FROM User WHERE username=? OR email=?", userName, userEmail)
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
			return Error.REGISTER_USERNAME_EXISTED, nil
		}
		if email == userEmail {
			return Error.REGISTER_EMAIL_EXISTED, nil
		}
	}

	return Error.NULL, nil
}
