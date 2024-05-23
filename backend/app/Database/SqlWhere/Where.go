package SqlWhere

//NOT < > LIKE AND OR

const (
	SQL_OPERATOR_NOT     = 0b000001
	SQL_OPERATOR_LIKE    = 0b000010
	SQL_OPERATOR_LESS    = 0b000100
	SQL_OPERATOR_GREATER = 0b001000
	SQL_OPERATOR_AND     = 0b010000
	SQL_OPERATOR_OR      = 0b100000
)

var (
	AcceptOperators = map[string]map[string]int{
		"User": {
			"id":            0b100000,
			"username":      0b100010,
			"email":         0b100010,
			"last_edit":     0b011101,
			"register_time": 0b011101,
			"deleted":       0b011101,
		},
	}
)

type Operators struct {
	NOT     bool
	LESS    bool
	GREATER bool
	LIKE    bool
	AND     bool
	OR      bool
}

func OperatorsArray() (*Operators, []*bool) {
	operators := Operators{}

	array := []*bool{
		&operators.NOT,
		&operators.LIKE,
		&operators.LESS,
		&operators.GREATER,
		&operators.AND,
		&operators.OR,
	}

	return &operators, array
}
