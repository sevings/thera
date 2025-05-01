package thera

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var ErrNotFound = fmt.Errorf("record not found")

type DB struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

type Chat struct {
	gorm.Model
	UserID  int64
	AgentID string
}

func LoadDatabase(path string) (*DB, bool) {
	log := zap.L().Named("db").Sugar()
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		log.Error(err)
		return nil, false
	}

	err = db.AutoMigrate(&Chat{})
	if err != nil {
		log.Error(err)
		return nil, false
	}

	return &DB{
		db:  db,
		log: log,
	}, true
}
