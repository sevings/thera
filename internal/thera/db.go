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
	UserID     int64     `gorm:"unique"`
	IdentityID string    `gorm:"unique"`
	AgentID    string    `gorm:"unique"`
	TalkLevel  TalkLevel `gorm:"default:0"`
}

type TalkLevel uint8

const (
	TalkLevelAloud    TalkLevel = 0
	TalkLevelThoughts TalkLevel = 1
	TalkLevelVerbose  TalkLevel = 2
)

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

// CreateChat creates a new chat
func (d *DB) CreateChat(chat *Chat) (*Chat, error) {
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

// UpdateChatTalkLevel updates the talk level for a chat by user ID
func (d *DB) UpdateChatTalkLevel(userID int64, newTalkLevel TalkLevel) error {
	result := d.db.Model(&Chat{}).Where("user_id = ?", userID).Update("talk_level", newTalkLevel)

	if result.Error != nil {
		d.log.Errorf("Failed to update talk level for user %d: %v", userID, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		d.log.Errorf("No chat found for user ID %d", userID)
		return ErrNotFound
	}

	return nil
}
