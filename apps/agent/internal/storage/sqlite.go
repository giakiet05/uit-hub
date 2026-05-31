package storage

import (
	"encoding/json"
	"time"

	"github.com/giakiet05/uit-hub/apps/agent/internal/conversation"
	"github.com/giakiet05/uit-hub/apps/agent/internal/session"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SessionRecord is the GORM model for storing a session.
type SessionRecord struct {
	ID        string    `gorm:"primaryKey"`
	StartedAt time.Time
	UpdatedAt time.Time
	Messages  string    // JSON encoded []conversation.Message
	Usage     string    // JSON encoded session.Usage
}

// TableName overrides the table name used by SessionRecord to `sessions`.
func (SessionRecord) TableName() string {
	return "sessions"
}

// DB wraps the GORM database instance.
type DB struct {
	gormDB *gorm.DB
}

// InitDB connects to the SQLite database and migrates the schema.
func InitDB(dsn string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&SessionRecord{})
	if err != nil {
		return nil, err
	}

	return &DB{gormDB: db}, nil
}

// Envelope is used to serialize and deserialize the conversation.Message interface
type messageEnvelope struct {
	Role string          `json:"role"`
	Data json.RawMessage `json:"data"`
}

func marshalMessages(messages []conversation.Message) ([]byte, error) {
	envelopes := make([]messageEnvelope, 0, len(messages))
	for _, msg := range messages {
		role := string(conversation.RoleOf(msg))
		data, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}
		envelopes = append(envelopes, messageEnvelope{Role: role, Data: data})
	}
	return json.Marshal(envelopes)
}

func unmarshalMessages(raw string) ([]conversation.Message, error) {
	if raw == "" {
		return nil, nil
	}
	var envelopes []messageEnvelope
	if err := json.Unmarshal([]byte(raw), &envelopes); err != nil {
		return nil, err
	}

	var messages []conversation.Message
	for _, env := range envelopes {
		switch conversation.Role(env.Role) {
		case conversation.RoleSystem:
			var m conversation.SystemMessage
			if err := json.Unmarshal(env.Data, &m); err == nil {
				messages = append(messages, m)
			}
		case conversation.RoleUser:
			var m conversation.UserMessage
			if err := json.Unmarshal(env.Data, &m); err == nil {
				messages = append(messages, m)
			}
		case conversation.RoleAssistant:
			var m conversation.AssistantMessage
			if err := json.Unmarshal(env.Data, &m); err == nil {
				messages = append(messages, m)
			}
		case conversation.RoleTool:
			var m conversation.ToolResultMessage
			if err := json.Unmarshal(env.Data, &m); err == nil {
				messages = append(messages, m)
			}
		}
	}
	return messages, nil
}

// UpsertSession saves or updates a session's history and usage metrics.
func (db *DB) UpsertSession(id string, startedAt time.Time, messages []conversation.Message, usage session.Usage) error {
	messagesJSON, err := marshalMessages(messages)
	if err != nil {
		return err
	}

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return err
	}

	record := SessionRecord{
		ID:        id,
		StartedAt: startedAt,
		Messages:  string(messagesJSON),
		Usage:     string(usageJSON),
	}

	// Save will update if primary key exists, or insert if it doesn't.
	return db.gormDB.Save(&record).Error
}

// GetSession retrieves a session by its ID.
func (db *DB) GetSession(id string) (*SessionRecord, error) {
	var record SessionRecord
	err := db.gormDB.First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetLastSession retrieves the most recently updated session.
func (db *DB) GetLastSession() (*SessionRecord, error) {
	var record SessionRecord
	err := db.gormDB.Order("updated_at desc").First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// ListSessions retrieves recent sessions ordered by updated_at descending.
func (db *DB) ListSessions(limit int) ([]SessionRecord, error) {
	var records []SessionRecord
	err := db.gormDB.Order("updated_at desc").Limit(limit).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// Decode converts a SessionRecord back into domain models.
func (r *SessionRecord) Decode() ([]conversation.Message, session.Usage, error) {
	messages, err := unmarshalMessages(r.Messages)
	if err != nil {
		return nil, session.Usage{}, err
	}

	var usage session.Usage
	if r.Usage != "" {
		if err := json.Unmarshal([]byte(r.Usage), &usage); err != nil {
			return nil, session.Usage{}, err
		}
	}

	return messages, usage, nil
}
