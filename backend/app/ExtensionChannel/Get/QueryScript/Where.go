package QueryScript

import (
	"simple_account/app/Database/SqlWhere"
	"simple_account/app/Http/Message"
	"simple_account/app/Http/Url"
	"strconv"
	"strings"
)

func Wheres(context *Message.Context, tableName string) (string, []interface{}) {
	query := query(context.Url)

	if query == nil {
		return "", []interface{}{}
	}

	acceptOperators := SqlWhere.AcceptOperators[tableName]

	allValues := []interface{}{}

	allColumnsWhereArray := []string{}

	for key, rawValues := range *query {
		acceptOperator, exist := acceptOperators[key]
		if !exist {
			continue
		}

		columnWhereArray := []string{}
		values, isStrArray := rawValues.([]string)
		if !isStrArray {
			continue
		}

		for _, rawValue := range values {
			if rawValue == "" {
				continue
			}
			value, operators := getOperators(rawValue)
			columnWhere := "("

			finalOperators, operatorsArray := SqlWhere.OperatorsArray()
			length := len(strconv.FormatInt(int64(acceptOperator), 2))
			for i := 0; i < length; i++ {
				if (operators&(1<<i) != 0) && (acceptOperator&(1<<i) != 0) {
					(*operatorsArray[i]) = true
				}

			}

			if finalOperators.NOT {
				columnWhere += "NOT"
			}

			columnWhere += " " + key

			if finalOperators.LIKE {
				columnWhere += " LIKE"
				value = "%" + value + "%"
			}

			if finalOperators.LESS || finalOperators.GREATER {
				if finalOperators.LESS {
					columnWhere += "<"
				}
				if finalOperators.GREATER {
					columnWhere += ">"
				}
			} else if !finalOperators.LIKE {
				columnWhere += " ="
			}

			columnWhere += " ? )"

			columnWhereArray = append(columnWhereArray, columnWhere)
			allValues = append(allValues, value)
		}

		logic := " AND "
		if acceptOperator >= SqlWhere.SQL_OPERATOR_OR {
			logic = " OR "
		}

		columnWhereStr := "(" + strings.Join(columnWhereArray, logic) + ")"
		allColumnsWhereArray = append(allColumnsWhereArray, columnWhereStr)
	}

	if len(allValues) == 0 {
		return "", []interface{}{}
	}

	allColumnsWhereStr := " WHERE " + strings.Join(allColumnsWhereArray, " AND ")

	return allColumnsWhereStr, allValues
}

func query(url *Url.Url) *map[string]interface{} {
	query := url.GetQuery()

	rawWheres, exist := query["where"]

	if !exist {
		return nil
	}

	if wheres, ok := rawWheres.(map[string]interface{}); ok {
		return &wheres
	}

	return nil
}

func getOperators(value string) (string, int) {
	operator := 0
	enableOperators := SqlWhere.Operators{}
	for {
		action := false
		if value[:1] == "!" {
			value = value[1:]

			if !enableOperators.NOT {
				operator += SqlWhere.SQL_OPERATOR_NOT
				enableOperators.NOT = true
			}

			action = true
		}

		if value[:1] == "<" {
			value = value[1:]

			if !enableOperators.LESS {
				operator += SqlWhere.SQL_OPERATOR_LESS
				enableOperators.LESS = true
			}

			action = true
		}

		if value[:1] == ">" {
			value = value[1:]

			if !enableOperators.GREATER {
				operator += SqlWhere.SQL_OPERATOR_GREATER
				enableOperators.GREATER = true
			}

			action = true
		}

		if value[:1] == "%" {
			value = value[1:]

			if !enableOperators.LIKE {
				operator += SqlWhere.SQL_OPERATOR_LIKE
				enableOperators.LIKE = true
			}

			action = true
		}
		if !action {
			break
		}
	}

	if value[:1] == "\\" {
		value = value[1:]
	}

	return value, operator
}
