package model

import "time"

type SeatStatus string

const (
	SeatStatusAvailable SeatStatus = "available"
	SeatStatusReserved  SeatStatus = "reserved"
	SeatStatusSold      SeatStatus = "sold"
)

type Seat struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement"`
	EventID    uint64     `gorm:"not null;index"`
	SeatNumber string     `gorm:"size:20;not null"`
	Price      int64      `gorm:"not null"`
	Status     SeatStatus `gorm:"type:seat_status;not null;default:available"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}

func (Seat) TableName() string {
	return "seats"
}
