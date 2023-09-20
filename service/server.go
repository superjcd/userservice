package service

import (
	"context"

	"github.com/superjcd/userservice/config"
	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/pkg/database"
	"github.com/superjcd/userservice/service/store"
	"github.com/superjcd/userservice/service/store/sql"
	"gorm.io/gorm"
)

var _DB *gorm.DB

// Server Server struct
type Server struct {
	v1.UnimplementedUserServiceServer
	datastore store.Factory
	client    v1.UserServiceClient
	conf      *config.Config
}

// NewServer New service grpc server
func NewServer(conf *config.Config, client v1.UserServiceClient) (v1.UserServiceServer, error) {
	_DB = database.MustPreParePostgresqlDb(&conf.Pg)
	factory, err := sql.NewSqlStoreFactory(_DB)
	if err != nil {
		return nil, err
	}

	server := &Server{
		client:    client,
		datastore: factory,
		conf:      conf,
	}

	return server, nil
}

func (s *Server) CreateGroup(ctx context.Context, rq *v1.CreateGroupRequest) (*v1.CreateGroupResponse, error) {
	if err := s.datastore.Groups().Create(ctx, rq); err != nil {
		return &v1.CreateGroupResponse{Msg: "创建失败", Status: v1.Status_failure, Group: nil}, err
	}
	return &v1.CreateGroupResponse{
		Msg:    "创建成功",
		Status: v1.Status_success,
		Group:  &v1.Group{Name: rq.Groupname},
	}, nil
}

func (s *Server) ListGroup(ctx context.Context, rq *v1.ListGroupRequest) (*v1.ListGroupResponse, error) {
	groups, err := s.datastore.Groups().List(ctx, rq)
	if err != nil {
		return &v1.ListGroupResponse{Msg: "获取列表失败", Status: v1.Status_failure}, err
	}

	resp := groups.ConvertToListGroupResponse("成功获取列表", v1.Status_success)

	return &resp, nil
}

func (s *Server) UpdateGroup(ctx context.Context, rq *v1.UpdateGroupRequest) (*v1.UpdateGroupResponse, error) {
	if err := s.datastore.Groups().Update(ctx, rq); err != nil {
		return &v1.UpdateGroupResponse{Msg: "失败", Status: v1.Status_success}, err
	}

	return &v1.UpdateGroupResponse{
		Msg:    "更新成功",
		Status: v1.Status_success,
	}, nil

}

func (s *Server) DeleteGroup(ctx context.Context, rq *v1.DeleteGroupRequest) (*v1.DeleteGroupResponse, error) {
	if err := s.datastore.Groups().Delete(ctx, rq); err != nil {
		return &v1.DeleteGroupResponse{Msg: "删除失败", Status: v1.Status_failure}, err
	}

	return &v1.DeleteGroupResponse{Msg: "删除成功", Status: v1.Status_success}, nil
}

func (s *Server) CreateUser(ctx context.Context, rq *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	newCtx, cancel := context.WithCancel(ctx) // cancel
	defer cancel()

	if err := s.datastore.Users().Create(newCtx, rq); err != nil {
		return &v1.CreateUserResponse{Msg: "创建用户失败", Status: v1.Status_failure}, err
	}

	// 调用一下notification client ,失败了上面cancel掉

	return &v1.CreateUserResponse{Msg: "创建用户成功", Status: v1.Status_success}, nil

}

func (s *Server) ListUser(ctx context.Context, rq *v1.ListUserRequest) (*v1.ListUserResponse, error) {
	users, err := s.datastore.Users().List(ctx, rq)

	if err != nil {
		return &v1.ListUserResponse{Msg: "获取列表失败", Status: v1.Status_failure, Users: nil}, err
	}

	resp := users.ConvertToListUserResponse("获取列表成功", v1.Status_success)
	return &resp, nil
}

func (s *Server) UpdateUser(ctx context.Context, rq *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	err := s.datastore.Users().Update(ctx, rq)
	if err != nil {
		return &v1.UpdateUserResponse{Msg: "更新失败", Status: v1.Status_failure}, err
	}

	return &v1.UpdateUserResponse{Msg: "更新成功", Status: v1.Status_success}, nil
}

func (s *Server) RemoveUser(ctx context.Context, rq *v1.RemoveUserRequest) (*v1.RemoveUserResponse, error) {
	if err := s.datastore.Users().Delete(ctx, rq); err != nil {
		return &v1.RemoveUserResponse{Msg: "删除用户失败", Status: v1.Status_failure}, err
	}
	return &v1.RemoveUserResponse{Msg: "删除用户成功", Status: v1.Status_failure}, nil
}

func (s *Server) UpdateUserPassword(ctx context.Context, rq *v1.UpdateUserPasswordRequest) (*v1.UpdateUserPasswordResponse, error) {
	if err := s.datastore.Users().UpdatePassword(ctx, rq); err != nil {
		return &v1.UpdateUserPasswordResponse{Msg: "更新密码失败", Status: v1.Status_failure}, err
	}

	return &v1.UpdateUserPasswordResponse{Msg: "更新密码成功", Status: v1.Status_success}, nil
}

func (s *Server) ResetUserPassword(ctx context.Context, rq *v1.ResetUserPasswordRequest) (*v1.ResetUserPasswordResponse, error) {
	// newCtx, cancel := context.WithCancel(ctx)
	// notification
	return &v1.ResetUserPasswordResponse{Msg: "更新用户成功", Status: v1.Status_success}, nil

}
