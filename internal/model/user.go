package model

import "time"

type UserRole string

const (
	RoleAdmin     UserRole = "admin"
	RoleOrganizer UserRole = "organizer"
	RoleCustomer  UserRole = "customer"
)

type User struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;size:255;not null"`
	Role         UserRole  `gorm:"type:user_role;default:customer;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
