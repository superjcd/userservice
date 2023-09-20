package sql

import (
	"context"
	"fmt"

	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/service/store"
	"gorm.io/gorm"
)

type groups struct {
	db *gorm.DB
}

var _ store.GroupStore = (*groups)(nil)

func (g *groups) Create(ctx context.Context, rq *v1.CreateGroupRequest) error {

	group := store.Group{
		Name: rq.Groupname,
	}

	return g.db.Create(&group).Error // 我只存储了用户， 但没有处理和用户group有关的逻辑
}

func (g *groups) List(ctx context.Context, rq *v1.ListGroupRequest) (*store.GroupList, error) {
	result := &store.GroupList{}

	var where_clause string
	if rq.Part == "" {
		where_clause = "1 = 1"
	} else {
		where_clause = fmt.Sprintf("name like '%%%s%%'", rq.Part)
	}

	d := g.db.Where(where_clause).
		Offset(int(rq.Offset)).
		Limit(int(rq.Limit)).
		Find(&result.Items).
		Offset(-1).
		Limit(-1).
		Count(&result.TotalCount)

	return result, d.Error
}

func (g *groups) Update(ctx context.Context, rq *v1.UpdateGroupRequest) error {
	group := store.Group{}
	if err := g.db.Where("name = ?", rq.OldName).First(&group).Error; err != nil {
		return err
	}

	group.Name = rq.NewName

	return g.db.Save(&group).Error
}

func (g *groups) Delete(ctx context.Context, rq *v1.DeleteGroupRequest) error {
	return g.db.Where("name = ?", rq.Name).Delete(&store.Group{}).Error
}
