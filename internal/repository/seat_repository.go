package repository

import (
	"context"
	"errors"

	"github.com/mjhddev/ticket-booking-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SeatRepository interface {
	Create(
		ctx context.Context,
		seat *model.Seat,
	) error

	CreateMany(
		ctx context.Context,
		seats []model.Seat,
	) error

	FindByID(
		ctx context.Context,
		id uint64,
	) (*model.Seat, error)

	FindByEventID(
		ctx context.Context,
		eventID uint64,
	) ([]model.Seat, error)

	FindSeatForUpdate(
		ctx context.Context,
		tx *gorm.DB,
		eventID uint64,
		seatID uint64,
	) (*model.Seat, error)

	UpdateStatus(
		ctx context.Context,
		tx *gorm.DB,
		seatID uint64,
		status model.SeatStatus,
	) error
}

type seatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{
		db: db,
	}
}

func (r *seatRepository) Create(
	ctx context.Context,
	seat *model.Seat,
) error {
	return r.db.WithContext(ctx).
		Create(seat).
		Error
}

func (r *seatRepository) CreateMany(
	ctx context.Context,
	seats []model.Seat,
) error {
	return r.db.WithContext(ctx).
		Create(&seats).
		Error
}

func (r *seatRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Seat, error) {
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

func (r *seatRepository) FindByEventID(
	ctx context.Context,
	eventID uint64,
) ([]model.Seat, error) {
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

func (r *seatRepository) FindSeatForUpdate(
	ctx context.Context,
	tx *gorm.DB,
	eventID uint64,
	seatID uint64,
) (*model.Seat, error) {
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

func (r *seatRepository) UpdateStatus(
	ctx context.Context,
	tx *gorm.DB,
	seatID uint64,
	status model.SeatStatus,
) error {
	return tx.WithContext(ctx).
		Model(&model.Seat{}).
		Where("id = ?", seatID).
		Update("status", status).
		Error
}
