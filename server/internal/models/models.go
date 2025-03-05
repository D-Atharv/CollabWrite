package models

import (
	"time"

	"github.com/google/uuid"
)

// User Model (OAuth2 Authentication)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"` 
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"` // Empty if OAuth2 user
	Provider     string    `gorm:"not null"` // "google", "github", etc.
	ProviderID   string    `gorm:"unique"`   // Google/GitHub OAuth2 ID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Document Model (Real-time Collaboration)
type Document struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"type:text;not null"`
	OwnerID   uuid.UUID `gorm:"not null;index"`
	Owner     User      `gorm:"foreignKey:OwnerID"`
	Version   int       `gorm:"default:1"` // Tracks document versions
	CreatedAt time.Time
	UpdatedAt time.Time
}

// DocumentHistory Model (Versioning System)
type DocumentHistory struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DocumentID uuid.UUID `gorm:"not null;index"`
	Document   Document  `gorm:"foreignKey:DocumentID"`
	Content    string    `gorm:"type:text;not null"`
	Version    int       `gorm:"not null"`
	EditedBy   uuid.UUID `gorm:"not null;index"`
	Editor     User      `gorm:"foreignKey:EditedBy"`
	EditedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// Access Control (For Sharing Documents)
type DocumentAccess struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	DocumentID uuid.UUID `gorm:"not null;index"`
	Document   Document  `gorm:"foreignKey:DocumentID"`
	UserID     uuid.UUID `gorm:"not null;index"`
	User       User      `gorm:"foreignKey:UserID"`
	Role       string    `gorm:"not null"` // "editor", "viewer"
	CreatedAt  time.Time
}

// API Tokens (For OAuth2 & Refresh Tokens)
type Token struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	Token     string    `gorm:"not null;unique"`
	ExpiresAt time.Time
}

