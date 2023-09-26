package sql

import (
	"context"

	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/pkg/passwd"
	"github.com/superjcd/userservice/service/store"
	"gorm.io/gorm"
)

type users struct {
	db *gorm.DB
}

var _ store.UserStore = (*users)(nil)

func (u *users) Create(ctx context.Context, rq *v1.CreateUserRequest) error {
	isAdmin := 0

	if rq.Role >= v1.Role_admin {
		isAdmin = 1
	}

	user := store.User{
		Name:      rq.Username,
		Email:     rq.Email,
		IsAdmin:   isAdmin,
		Creator:   rq.Creator,
		RoleLevel: int(rq.Role),
	}

	return u.db.Create(&user).Error // 我只存储了用户， 但没有处理和用户group有关的逻辑
}

func (u *users) List(ctx context.Context, rq *v1.ListUserRequest) (*store.UserList, error) {
	result := &store.UserList{}
	tx := u.db

	if rq.Username != "" {
		tx = tx.Where("name = ?", rq.Username)
	}

	if rq.Email != "" {
		tx = tx.Where("email = ?", rq.Email)
	}

	if rq.Creator != "" {
		tx = tx.Where("creator = ?", rq.Creator)
	}

	d := tx.
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

	if rq.Role >= 1 {
		isAdmin = 1
	}

	user.RoleLevel = int(rq.Role)
	user.IsAdmin = isAdmin
	if rq.Creator != "" {
		user.Creator = rq.Creator
	}

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
	return u.db.Unscoped().Where("email = ?", rq.Email).Delete(&store.User{}).Error
}
