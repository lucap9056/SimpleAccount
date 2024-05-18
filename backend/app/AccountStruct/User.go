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
	LastEditTime *time.Time `json:"lastEditTime"`
	lastEditTime []uint8    `json:"-"`
	CreateTime   *time.Time `json:"createTime"`
	createTime   []uint8    `json:"-"`
	DeletedTime  *time.Time `json:"deletedTime"`
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

	if user.CreateTime != nil {
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
		"user.id":          &user.Id,
		"user.username":    &user.Username,
		"user.email":       &user.Email,
		"user.salt":        &user.Salt,
		"user.hash":        &user.Hash,
		"user.last_edit":   &user.lastEditTime,
		"user.create_time": &user.createTime,
		"user.deleted":     &user.deletedTime,
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

func (user *User) Convert() {
	if len(user.lastEditTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.lastEditTime))
		if err == nil {
			user.LastEditTime = &time
		}
		user.lastEditTime = []uint8{}
	}

	if len(user.createTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.createTime))
		if err == nil {
			user.CreateTime = &time
		}
		user.createTime = []uint8{}
	}

	if len(user.deletedTime) > 0 {
		time, err := time.Parse("2006-01-02 15:04:05", string(user.deletedTime))
		if err == nil {
			user.DeletedTime = &time
		}
		user.deletedTime = []uint8{}
	}
}
