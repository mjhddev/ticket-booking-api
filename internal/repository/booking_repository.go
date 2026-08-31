package repository

import (
	"context"
	"errors"

	"github.com/mjhddev/ticket-booking-api/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BookingRepository interface {
	CreateBooking(
		ctx context.Context,
		tx *gorm.DB,
		booking *model.Booking,
	) error

	CreateBookingItems(
		ctx context.Context,
		tx *gorm.DB,
		items []model.BookingItem,
	) error

	FindByID(
		ctx context.Context,
		id uint64,
	) (*model.Booking, error)

	FindByIDWithItems(
		ctx context.Context,
		id uint64,
	) (*model.Booking, error)

	FindByIDForUpdate(
		ctx context.Context,
		tx *gorm.DB,
		id uint64,
	) (*model.Booking, error)

	FindItemsByBookingID(
		ctx context.Context,
		tx *gorm.DB,
		bookingID uint64,
	) ([]model.BookingItem, error)

	UpdateStatus(
		ctx context.Context,
		tx *gorm.DB,
		bookingID uint64,
		status model.BookingStatus,
	) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) CreateBooking(
	ctx context.Context,
	tx *gorm.DB,
	booking *model.Booking,
) error {
	return tx.WithContext(ctx).
		Create(booking).
		Error
}

func (r *bookingRepository) CreateBookingItems(
	ctx context.Context,
	tx *gorm.DB,
	items []model.BookingItem,
) error {
	return tx.WithContext(ctx).
		Create(&items).
		Error
}

func (r *bookingRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*model.Booking, error) {
	var booking model.Booking

	err := r.db.WithContext(ctx).
		First(&booking, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByIDWithItems(
	ctx context.Context,
	id uint64,
) (*model.Booking, error) {
	var booking model.Booking

	err := r.db.WithContext(ctx).
		Preload("Items").
		First(&booking, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindByIDForUpdate(
	ctx context.Context,
	tx *gorm.DB,
	id uint64,
) (*model.Booking, error) {
	var booking model.Booking

	err := tx.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		First(&booking, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &booking, nil
}

func (r *bookingRepository) FindItemsByBookingID(
	ctx context.Context,
	tx *gorm.DB,
	bookingID uint64,
) ([]model.BookingItem, error) {
	var items []model.BookingItem

	err := tx.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("seat_id ASC").
		Find(&items).
		Error

	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *bookingRepository) UpdateStatus(
	ctx context.Context,
	tx *gorm.DB,
	bookingID uint64,
	status model.BookingStatus,
) error {
	return tx.WithContext(ctx).
		Model(&model.Booking{}).
		Where("id = ?", bookingID).
		Update("status", status).
		Error
}
