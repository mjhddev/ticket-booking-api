package repository

import (
	"context"
	"errors"

	"github.com/mjhddev/ticket-booking-api/internal/model"
	"gorm.io/gorm"
)

type SeatRepository interface {
	CreateBatch(ctx context.Context, seats []model.Seat) error
	FindByID(ctx context.Context, id uint64) (*model.Seat, error)
	FindByEventID(ctx context.Context, eventID uint64) ([]model.Seat, error)
	FindByEventAndNumber(ctx context.Context, eventID uint64, seatNumber string) (*model.Seat, error)
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{
		db: db,
	}
}

func (r *seatRepository) CreateBatch(ctx context.Context, seats []model.Seat) error {
	return r.db.WithContext(ctx).Create(&seats).Error
}

func (r *seatRepository) FindByID(ctx context.Context, id uint64) (*model.Seat, error) {
	var seat model.Seat

	err := r.db.WithContext(ctx).
		First(&seat, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &seat, nil
}

func (r *seatRepository) FindByEventID(ctx context.Context, eventID uint64) ([]model.Seat, error) {
	var seats []model.Seat

	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("seat_number ASC").
		Find(&seats).
		Error

	if err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *seatRepository) FindByEventAndNumber(ctx context.Context, eventID uint64, seatNumber string) (*model.Seat, error) {
	var seat model.Seat

	err := r.db.WithContext(ctx).
		Where(
			"event_id = ? AND seat_number = ?",
			eventID,
			seatNumber,
		).
		First(&seat).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &seat, nil
}
