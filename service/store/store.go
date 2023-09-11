package store

import (
	"context"

	v1 "github.com/HooYa-Bigdata/userservice/genproto/v1"
)

type Factory interface {
	Users() UserStore
	Groups() GroupStore
	Close() error
}

type UserStore interface {
	Create(ctx context.Context, _ *v1.InviteUserRequest) error
	List(ctx context.Context, _ *v1.ListUserRequest) (*UserList, error)
	Update(ctx context.Context, _ *v1.UpdateUserRequest) error
	UpdatePassword(ctx context.Context, _ *v1.UpdateUserPasswordRequest) error
	Delete(ctx context.Context, _ *v1.RemoveUserRequest) error
}

type GroupStore interface {
	Create(ctx context.Context, _ *v1.CreateGroupRequest) error
	List(ctx context.Context, _ *v1.ListGroupRequest) (*GroupList, error)
	Update(ctx context.Context, _ *v1.UpdateGroupRequest) error
	Delete(ctx context.Context, _ *v1.DeleteGroupRequest) error
}
