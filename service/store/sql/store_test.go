package sql

import (
	"context"
	"os"
	"testing"

	v1 "github.com/HooYa-Bigdata/userservice/genproto/v1"
	"github.com/HooYa-Bigdata/userservice/service/store"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	newUser := &v1.InviteUserRequest{
		Invitee: &v1.User{
			Username: "jcd",
			Email:    "jcd@example.com",
		},
		Groups: []*v1.Group{{Name: "testgroup"}}, // TODO: 这一块其实需要结合RDBC实现
		Role:   0,
	}
	err := suite.FakeFactory.Users().Create(context.Background(), newUser)
	assert.Nil(suite.T(), err)
}

func (suite *FakeStoreTestSuite) TestListUserWithUserName() {
	request := &v1.ListUserRequest{
		Part:   "jcd",
		Offset: 0,
		Limit:  10,
	}
	userList, err := suite.FakeFactory.Users().List(context.Background(), request)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(userList.Items))
}

func (suite *FakeStoreTestSuite) TestUpdateUser() {
	request := &v1.UpdateUserRequest{
		User: &v1.User{
			Username: "jcd2",
			Email:    "jcd@example.com",
		},
		Role: 1,
	}

	err := suite.FakeFactory.Users().Update(context.Background(), request)
	assert.Nil(suite.T(), err)

	request2 := &v1.ListUserRequest{
		Part:   "jcd2",
		Offset: 0,
		Limit:  10,
	}
	userList, err := suite.FakeFactory.Users().List(context.Background(), request2)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(userList.Items))
	assert.Equal(suite.T(), 0, userList.Items[0].IsAdmin) // 不再是admin了
}

func TestFakeStoreSuite(t *testing.T) {
	suite.Run(t, new(FakeStoreTestSuite))
}
