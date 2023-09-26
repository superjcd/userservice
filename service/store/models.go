package store

import (
	"time"

	v1 "github.com/superjcd/userservice/genproto/v1"
	"gorm.io/gorm"
)

// gorm 模型定义

type User struct {
	gorm.Model
	Name      string    `json:"name" gorm:"column:name"`
	Password  string    `json:"password" gorm:"column:password"`
	Email     string    `json:"email" gorm:"column:email;uniqueIndex;size:30" `
	IsAdmin   int       `json:"is_admin,omitempty" gorm:"column:is_admin"`
	RoleLevel int       `json:"role_level,omitempty" gorm:"column:role_level"`
	Creator   string    `json:"creator" gorm:"column:creator"` // 这个会是email
	LogindAt  time.Time `json:"logine_at,omitempty" gorm:"column:logine_at"`
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
			Creator:  user.Creator,
		})
	}

	return v1.ListUserResponse{
		Msg:    msg,
		Status: status,
		Users:  users,
	}
}

type Group struct {
	gorm.Model
	Name    string `json:"name" gorm:"column:name;uniqueIndex;size:30" validate:"required,min=1,max=30"`
	Creator string `json:"creator" gorm:"column:creator"`
}

type GroupList struct {
	TotalCount int64    `json:"totalCount"`
	Items      []*Group `json:"items"`
}

func (gl *GroupList) ConvertToListGroupResponse(msg string, status v1.Status) v1.ListGroupResponse {
	groups := make([]*v1.Group, 8)

	for _, group := range gl.Items {
		groups = append(groups, &v1.Group{
			Name:    group.Name,
			Creator: group.Creator,
		})
	}

	return v1.ListGroupResponse{
		Msg:    msg,
		Status: status,
		Groups: groups,
	}

}

func MigrateDatabase(db *gorm.DB) error {
	if err := db.AutoMigrate(User{}, Group{}); err != nil {
		return err
	}
	return nil
}
