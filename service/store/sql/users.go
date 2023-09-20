package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/pkg/passwd"
	"github.com/superjcd/userservice/service/store"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

var _ store.UserStore = (*users)(nil)

func (u *users) Create(ctx context.Context, rq *v1.InviteUserRequest) error {
	isAdmin := 0

	if rq.Role == 0 {
		isAdmin = 1
	}

	user := store.User{
		Name:    rq.Invitee.Username,
		Email:   rq.Invitee.Email,
		IsAdmin: isAdmin,
	}

	return u.db.Create(&user).Error // 我只存储了用户， 但没有处理和用户group有关的逻辑
}

func (u *users) List(ctx context.Context, rq *v1.ListUserRequest) (*store.UserList, error) {
	result := &store.UserList{}

	var where_clause string
	if rq.Part == "" {
		where_clause = "1 = 1"
	} else {
		where_clause = fmt.Sprintf("name like '%%%s%%'", rq.Part)
	}

	d := u.db.Where(where_clause).
		Offset(int(rq.Offset)).
		Limit(int(rq.Limit)).
		Find(&result.Items).
		Offset(-1).
		Limit(-1).
		Count(&result.TotalCount)

	return result, d.Error
}

func (u *users) Update(ctx context.Context, rq *v1.UpdateUserRequest) error {
	user := &store.User{}
	if err := u.db.Where("email = ?", rq.User.Email).First(user).Error; err != nil {
		return err
	}

	user.Name = rq.User.Username
	isAdmin := 0
	if rq.Role == 0 {
		isAdmin = 1
	}
	user.IsAdmin = isAdmin

	return u.db.Save(user).Error
}

func (u *users) UpdatePassword(ctx context.Context, rq *v1.UpdateUserPasswordRequest) error {
	var err error
	user := store.User{}
	if err := u.db.Where("email = ?", rq.User.Email).First(&user).Error; err != nil {
		return err
	}
	if user.Password, err = passwd.Encrypt(rq.Password); err != nil {
		return err
	}
	return u.db.Save(&user).Error
}

func (u *users) Delete(ctx context.Context, rq *v1.RemoveUserRequest) error {
	return u.db.Where("email = ?", rq.Email).Delete(&store.User{}).Error
}
