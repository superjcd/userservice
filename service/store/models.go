package store

import (
	"time"

	v1 "github.com/superjcd/userservice/genproto/v1"
	"gorm.io/gorm"
)

// gorm 模型定义

type User struct {
	gorm.Model
	Name      string    `json:"name" gorm:"column:name" validate:"required,min=2,max=30"`
	Password  string    `json:"password" gorm:"column:password"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex;size:30" validate:"required,email,min=1,max=30"`
	IsAdmin   int       `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`
	LoginedAt time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

type UserList struct {
	TotalCount int64   `json:"totalCount"`
	Items      []*User `json:"items"`
}

func (ul *UserList) ConvertToListUserResponse(msg string, status v1.Status) v1.ListUserResponse {
	users := make([]*v1.User, 8)

	for _, user := range ul.Items {
		users = append(users, &v1.User{
			Username: user.Name,
			Email:    user.Email,
		})
	}
BLOB/TEXT column 'email' used in key specification without a key length
	return v1.ListUserResponse{
		Msg:    msg,
		Status: status,
		Users:  users,
	}
}

type Group struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name" validate:"required,min=1,max=30"`
}

type GroupList struct {
	TotalCount int64    `json:"totalCount"`
	Items      []*Group `json:"items"`
}

func (gl *GroupList) ConvertToListGroupResponse(msg string, status v1.Status) v1.ListGroupResponse {
	groups := make([]*v1.Group, 8)

	for _, group := range gl.Items {
		groups = append(groups, &v1.Group{
			Name: group.Name,
		})
	}

	return v1.ListGroupResponse{
		Msg:    msg,
		Status: status,
		Groups: groups,
	}

}

func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(Group{}); err != nil {
		return err
	}
	return nil
}
