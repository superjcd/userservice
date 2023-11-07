package sql

import (
	"context"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	v1 "github.com/superjcd/userservice/genproto/v1"
	"github.com/superjcd/userservice/pkg/passwd"
	"github.com/superjcd/userservice/service/store"
	"gorm.io/gorm"
)

var dbFile = "fake.db"

type FakeStoreTestSuite struct {
	suite.Suite
	Dbfile      string
	FakeFactory store.Factory
}

func (suite *FakeStoreTestSuite) SetupSuite() {
	file, err := os.Create(dbFile)
	assert.Nil(suite.T(), err)
	defer file.Close()

	suite.Dbfile = dbFile
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	assert.Nil(suite.T(), err)

	factory, err := NewSqlStoreFactory(db)
	assert.Nil(suite.T(), err)
	suite.FakeFactory = factory
}

func (suite *FakeStoreTestSuite) TearDownSuite() {
	var err error
	err = suite.FakeFactory.Close()
	assert.Nil(suite.T(), err)
	err = os.Remove(dbFile)
	assert.Nil(suite.T(), err)
}

// User
func (suite *FakeStoreTestSuite) TestInviteUser() {
	// Create User
	newUser := &v1.CreateUserRequest{
		Username: "jack",
		Email:    "jack@example.com",
		Role:     0,
		Creator:  "chaodi.jiang@fdsintl.com",
	}
	err := suite.FakeFactory.Users().Create(context.Background(), newUser)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestListUserWithUserName() {
	request := &v1.ListUserRequest{
		Username: "jack",
		Offset:   0,
		Limit:    10,
		Creator:  "chaodi.jiang@fdsintl.com",
	}
	userList, err := suite.FakeFactory.Users().List(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(userList.Items))
}

func (suite *FakeStoreTestSuite) TestUpdateUser() {
	request := &v1.UpdateUserRequest{
		Username: "jack2",
		Email:    "jack@example.com",
		Role:     v1.Role_superadmin,
	}

	err := suite.FakeFactory.Users().Update(context.Background(), request)
	assert.Nil(suite.T(), err)

	request2 := &v1.ListUserRequest{
		Username: "jack2",
		Offset:   0,
		Limit:    10,
	}
	userList, err := suite.FakeFactory.Users().List(context.Background(), request2)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(userList.Items))
	assert.Equal(suite.T(), 1, userList.Items[0].IsAdmin) // 是superaadmin
}

func (suite *FakeStoreTestSuite) TestZDeleteUser() {
	rq := &v1.RemoveUserRequest{
		Email: "jack@example.com",
	}

	err := suite.FakeFactory.Users().Delete(context.Background(), rq)
	assert.Nil(suite.T(), err)

	request2 := &v1.ListUserRequest{
		Email:  "jack@example.com",
		Offset: 0,
		Limit:  10,
	}
	userList, err := suite.FakeFactory.Users().List(context.Background(), request2)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(userList.Items))

}

func (suite *FakeStoreTestSuite) TestUpdateUserPasswordGroup() {
	// 创建用户
	rqCreate := &v1.CreateUserRequest{
		Username: "lucy",
		Email:    "lucy@example.com",
		Role:     0,
		Creator:  "chaodi.jiang@fdsintl.com",
	}
	err := suite.FakeFactory.Users().Create(context.Background(), rqCreate)
	assert.Nil(suite.T(), err)

	rqUpdatePass := &v1.UpdateUserPasswordRequest{
		Email:    "lucy@example.com",
		Password: "123456789",
	}

	err2 := suite.FakeFactory.Users().UpdatePassword(context.Background(), rqUpdatePass)
	assert.Nil(suite.T(), err2)

	rqGetUser := &v1.ListUserRequest{Email: "lucy@example.com", Limit: 1}
	userList, _ := suite.FakeFactory.Users().List(context.Background(), rqGetUser)
	assert.Equal(suite.T(), 1, len(userList.Items))

	assert.Equal(suite.T(), nil, passwd.Compare(userList.Items[0].Password, "123456789"))
}

func (suite *FakeStoreTestSuite) TestXResetPassword() {
	rqGetUser := &v1.ListUserRequest{Email: "lucy@example.com", Limit: 1}
	userList, _ := suite.FakeFactory.Users().List(context.Background(), rqGetUser)
	assert.Equal(suite.T(), 1, len(userList.Items))

	rqResetPass := &v1.ResetUserPasswordRequest{
		Email: "lucy@example.com",
	}

	err := suite.FakeFactory.Users().ResetPassword(context.Background(), rqResetPass)
	assert.Nil(suite.T(), err)

	rqGetUser2 := &v1.ListUserRequest{Email: "lucy@example.com", Limit: 1}
	userList2, _ := suite.FakeFactory.Users().List(context.Background(), rqGetUser2)
	assert.Equal(suite.T(), 1, len(userList2.Items))

	assert.Equal(suite.T(), "", userList2.Items[0].Password) // 重置后密码应该为空
}

func (suite *FakeStoreTestSuite) TestCreateGroup() {
	rq := &v1.CreateGroupRequest{
		Name: "testgroup",
		Type: "user",
	}

	err := suite.FakeFactory.Groups().Create(context.Background(), rq)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestDeleteGroup() {
	rq := &v1.DeleteGroupRequest{
		Name: "testgroup",
		Type: "user",
	}

	err := suite.FakeFactory.Groups().Delete(context.Background(), rq)
	assert.Nil(suite.T(), err)

	rq2 := &v1.DeleteGroupRequest{
		Name: "testgroup2", // not existed
		Type: "user",
	}
	err2 := suite.FakeFactory.Groups().Delete(context.Background(), rq2) // delete the same group
	assert.Nil(suite.T(), err2)
}

func TestFakeStoreSuite(t *testing.T) {
	suite.Run(t, new(FakeStoreTestSuite))
}
