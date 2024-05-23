package AccountStruct

import (
	"strings"
	"time"
)

type User struct {
	Id           int        `json:"id,omitempty"`
	Username     string     `json:"name,omitempty"`
	Email        string     `json:"email,omitempty"`
	Salt         string     `json:"-"`
	Hash         string     `json:"-"`
	LastEditTime *time.Time `json:"lastEditTime,omitempty"`
	lastEditTime []uint8    `json:"-"`
	RegisterTime *time.Time `json:"registerTime,omitempty"`
	registerTime []uint8    `json:"-"`
	DeletedTime  *time.Time `json:"deletedTime,omitempty"`
	deletedTime  []uint8    `json:"-"`
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

	if user.LastEditTime != nil {
		return false
	}

	if user.RegisterTime != nil {
		return false
	}

	if user.DeletedTime != nil {
		return false
	}

	return true
}

func (user *User) MappingTable(args ...string) []interface{} {

	if len(args) == 1 && strings.Contains(args[0], ",") {
		args = strings.Split(args[0], ",")
	}

	for i, arg := range args {
		if !strings.Contains(arg, ".") {
			args[i] = "user." + arg
		}
	}

	tableMap := map[string]interface{}{
		"user.id":            &user.Id,
		"user.username":      &user.Username,
		"user.email":         &user.Email,
		"user.salt":          &user.Salt,
		"user.hash":          &user.Hash,
		"user.last_edit":     &user.lastEditTime,
		"user.register_time": &user.registerTime,
		"user.deleted":       &user.deletedTime,
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

func (user *User) MoveTempToFinal() {
	if len(user.lastEditTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.lastEditTime))
		if err == nil {
			user.LastEditTime = &time
		}
		user.lastEditTime = []uint8{}
	}

	if len(user.registerTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.registerTime))
		if err == nil {
			user.RegisterTime = &time
		}
		user.registerTime = []uint8{}
	}

	if len(user.deletedTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.deletedTime))
		if err == nil {
			user.DeletedTime = &time
		}
		user.deletedTime = []uint8{}
	}
}
