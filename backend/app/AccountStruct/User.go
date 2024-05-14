package AccountStruct

import "time"

type User struct {
	Id         int        `json:"id,omitempty"`
	Username   string     `json:"name,omitempty"`
	Email      string     `json:"email,omitempty"`
	Salt       string     `json:"-"`
	Hash       string     `json:"-"`
	CreateTime *time.Time `json:"-"`
}

func (user *User) Empty() bool {
	if user.Id != 0 {
		return false
	}
	if user.Username != "" {
		return false
	}
	if user.Email != "" {
		return false
	}
	return true
}

func (user *User) MappingTable(args []string) []interface{} {
	tableMap := map[string]interface{}{
		"user.id":    &user.Id,
		"user.name":  &user.Username,
		"user.Email": &user.Email,
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
