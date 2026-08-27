package model

import "time"

type BookingItem struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	BookingID uint64    `gorm:"not null;index"`
	SeatID    uint64    `gorm:"not null;index"`
	Price     int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (BookingItem) TableName() string {
	return "booking_items"
}
