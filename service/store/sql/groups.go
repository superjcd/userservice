package sql

import (
	"context"

	v1 "github.com/HooYa-Bigdata/microservices/grpc_service/userservice/genproto/v1"
	"github.com/HooYa-Bigdata/microservices/grpc_service/userservice/service/store"
	"gorm.io/gorm"
)

type groups struct {
	db *gorm.DB
}

var _ store.GroupStore = (*groups)(nil)

func (u *groups) Create(ctx context.Context, rq *v1.CreateGroupRequest) error {

	group := store.Group{
		Name: rq.Groupname,
	}

	return u.db.Create(&group).Error // 我只存储了用户， 但没有处理和用户group有关的逻辑
}

func (u *groups) List(ctx context.Context, rq *v1.ListGroupRequest) (*store.GroupList, error) {
	result := &store.GroupList{}

	d := u.db.Where("name like ?", rq.Part).
		Offset(int(rq.Offset)).
		Limit(int(rq.Limit)).
		Find(&result.Items).
		Offset(-1).
		Limit(-1).
		Count(&result.TotalCount)

	return result, d.Error
}

func (u *groups) Update(ctx context.Context, rq *v1.UpdateGroupRequest) error {
	group := &store.Group{}
	if err := u.db.Where("name = ?", rq.OldName).First(group).Error; err != nil {
		return err
	}

	group.Name = rq.NewName

	return u.db.Save(&group).Error
}

func (u *groups) Delete(ctx context.Context, rq *v1.DeleteGroupRequest) error {
	return u.db.Where("name = ?", rq.Name).Delete(&store.Group{}).Error
}
