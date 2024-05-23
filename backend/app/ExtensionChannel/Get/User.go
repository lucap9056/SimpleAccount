package ExtensionGet

import (
	"encoding/json"
	"fmt"
	"simple_account/app/AccountStruct"
	"simple_account/app/Error"
	"simple_account/app/ExtensionChannel/Get/QueryScript"
	"simple_account/app/Http/Message"
)

func GetUser(context *Message.Context) (string, int, error) {

	includes := QueryScript.Includes(context.Url)
	if includes == "" {
		includes = "id"
	}

	wheres, whereValues := QueryScript.Wheres(context, "User")
	limit := QueryScript.Limit(context)

	connect := context.Database.Connect()

	query := fmt.Sprintf("SELECT %s FROM User%s%s", includes, wheres, limit)

	rows, err := connect.Query(query, whereValues...)
	if err != nil {
		return err.Error(), Error.EXTENSION_SQL_ERROR, nil
	}

	var users []AccountStruct.User
	for rows.Next() {

		var user AccountStruct.User
		columns := user.MappingTable(includes)

		err := rows.Scan(columns...)
		if err != nil {
			return err.Error(), Error.EXTENSION_SQL_ERROR, nil
		}
		user.MoveTempToFinal()
		users = append(users, user)
	}

	result, err := json.Marshal(users)
	if err != nil {
		return err.Error(), Error.EXTENSION_JSON_STRINGIFY_ERROR, nil
	}

	return string(result), Error.NULL, nil
}
