package repository

import (
	"context"
	"errors"

	"github.com/mjhddev/ticket-booking-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, tx *gorm.DB, booking *model.Booking) error
	CreateBookingItems(ctx context.Context, tx *gorm.DB, items []model.BookingItem) error
	FindSeatForUpdate(ctx context.Context, tx *gorm.DB, eventID uint64, seatID uint64) (*model.Seat, error)
	UpdateSeat(ctx context.Context, tx *gorm.DB, seat *model.Seat) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) CreateBooking(ctx context.Context, tx *gorm.DB, booking *model.Booking) error {
	return tx.WithContext(ctx).
		Create(booking).
		Error
}

func (r *bookingRepository) CreateBookingItems(ctx context.Context, tx *gorm.DB, items []model.BookingItem) error {
	return tx.WithContext(ctx).
		Create(&items).
		Error
}

func (r *bookingRepository) FindSeatForUpdate(ctx context.Context, tx *gorm.DB, eventID uint64, seatID uint64) (*model.Seat, error) {

	var seat model.Seat

	err := tx.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("event_id = ?", eventID).
		First(&seat, seatID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &seat, nil
}

func (r *bookingRepository) UpdateSeat(ctx context.Context, tx *gorm.DB, seat *model.Seat) error {
	return tx.WithContext(ctx).
		Save(seat).
		Error
}
