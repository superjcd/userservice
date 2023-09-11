package sql

import (
	"fmt"
	"sync"

	"github.com/HooYa-Bigdata/userservice/config"
	"github.com/HooYa-Bigdata/userservice/service/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type datastore struct {
	db *gorm.DB
}

func (ds *datastore) Users() store.UserStore {
	return &users{db: ds.db}
}

func (ds *datastore) Groups() store.GroupStore {
	return &groups{db: ds.db}
}

var (
	pgFactory  store.Factory
	sqlFactory store.Factory
	once       sync.Once
)

func NewPgStoreFactory(opts *config.Pg) (store.Factory, error) {
	if opts == nil && pgFactory == nil {
		return nil, fmt.Errorf("failed to get pg store fatory")
	}

	var err error
	var db *gorm.DB

	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
			opts.Host,
			opts.Username,
			opts.Password,
			opts.Database)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		store.MigrateDatabase(db)
		pgFactory = &datastore{db: db}
	})
	if pgFactory == nil || err != nil {
		return nil, err
	}
	return pgFactory, nil
}

func NewSqlStoreFactory(db *gorm.DB) (store.Factory, error) {
	if db == nil && sqlFactory == nil {
		return nil, fmt.Errorf("failed to get pg store fatory")
	}
	once.Do(func() {
		store.MigrateDatabase(db)
		sqlFactory = &datastore{db: db}
	})

	return sqlFactory, nil
}

func (ds *datastore) Close() error {
	db, _ := ds.db.DB()

	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
