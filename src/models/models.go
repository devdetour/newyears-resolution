package models

import (
	"time"

	"gorm.io/gorm"
)

type TokenSource int

const (
	StravaTokenSource TokenSource = iota
)

// User struct
type User struct {
	gorm.Model
	Username   string `gorm:"uniqueIndex;not null" json:"username"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Password   string `gorm:"not null" json:"password"`
	Names      string `json:"names"`
	AuthTokens []ExternalAuthToken
}

type ExternalAuthToken struct {
	gorm.Model
	UserId       uint        `gorm:"not null"`
	Text         string      `gorm:"not null" json:"text"`
	Scope        string      `json:"scope"`
	Source       TokenSource `json:"source"` // e.g. "strava"
	Expires      time.Time   `json:"expires"`
	RefreshToken string      `json:"refresh"`
}

type Contract struct {
	gorm.Model
}
