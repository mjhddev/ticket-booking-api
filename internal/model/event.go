package model

import "time"

type EventStatus string

const (
	EventStatusDraft     EventStatus = "draft"
	EventStatusPublished EventStatus = "published"
	EventStatusCancelled EventStatus = "cancelled"
)

type Event struct {
	ID          uint64      `gorm:"primaryKey;autoIncrement"`
	OrganizerID uint64      `gorm:"not null;index"`
	Title       string      `gorm:"size:150;not null"`
	Description string      `gorm:"type:text"`
	Location    string      `gorm:"size:255;not null"`
	StartTime   time.Time   `gorm:"not null"`
	EndTime     time.Time   `gorm:"not null"`
	Status      EventStatus `gorm:"type:event_status;not null;default:draft"`
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime"`
}

func (Event) TableName() string {
	return "events"
}
