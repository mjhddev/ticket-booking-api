package model

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
	BookingStatusExpired   BookingStatus = "expired"
)

type Booking struct {
	ID        uint64        `gorm:"primaryKey;autoIncrement"`
	UserID    uint64        `gorm:"not null;index"`
	EventID   uint64        `gorm:"not null;index"`
	Status    BookingStatus `gorm:"type:booking_status;not null;default:pending"`
	ExpiresAt time.Time     `gorm:"not null"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
	UpdatedAt time.Time     `gorm:"autoUpdateTime"`

	Items []BookingItem `gorm:"foreignKey:BookingID"`
}

func (Booking) TableName() string {
	return "bookings"
}
