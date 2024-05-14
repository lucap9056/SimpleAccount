package AccountStruct

type UserThird struct {
	Id    int    `json:"id,omitempty"`
	User  *User  `json:"user,omitempty"`
	Third *Third `json:"third,omitempty"`
}

func (t *UserThird) Init() {
	t.User = &User{}
	t.Third = &Third{}
}

func (t *UserThird) Clean() {
	if t.User.Empty() {
		t.User = nil
	}
	if t.Third.Empty() {
		t.Third = nil
	}
}

func (userThird *UserThird) MappingTable(args []string) []interface{} {
	tableMap := map[string]interface{}{
		"user_third.id":    &userThird.Id,
		"user_third.user":  &userThird.User,
		"user_third.third": &userThird.Third,
	}

	var columns []interface{}
	for _, arg := range args {
		refer, exist := tableMap[arg]
		if !exist {
			var i interface{}
			refer = &i
		}
		columns = append(columns, refer)
	}
	return columns
}
