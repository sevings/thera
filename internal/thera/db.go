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
	UserID     int64  `gorm:"uniqueIndex:idx_user_agent"`
	IdentityID string `gorm:"uniqueIndex:idx_user_agent"`
	AgentID    string `gorm:"uniqueIndex:idx_user_agent"`
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

// CreateChat creates a new chat with the given user ID and agent ID
func (d *DB) CreateChat(userID int64, identityID, agentID string) (*Chat, error) {
	chat := &Chat{
		UserID:     userID,
		IdentityID: identityID,
		AgentID:    agentID,
	}

	result := d.db.Create(chat)
	if result.Error != nil {
		d.log.Errorf("Failed to create chat: %v", result.Error)
		return nil, result.Error
	}

	return chat, nil
}

// GetChatByUserID retrieves a chat by user ID
func (d *DB) GetChatByUserID(userID int64) (*Chat, error) {
	var chat Chat
	result := d.db.Where("user_id = ?", userID).First(&chat)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound
		}
		d.log.Errorf("Failed to retrieve chat: %v", result.Error)
		return nil, result.Error
	}

	return &chat, nil
}
