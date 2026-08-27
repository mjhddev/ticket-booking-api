package repository

import (
	"context"
	"errors"

	"github.com/mjhddev/ticket-booking-api/internal/model"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	FindByID(ctx context.Context, id uint64) (*model.Event, error)
	FindAllPublished(ctx context.Context) ([]model.Event, error)
	Update(ctx context.Context, event *model.Event) error
	Delete(ctx context.Context, event *model.Event) error
	FindByIDTx(ctx context.Context, tx *gorm.DB, id uint64) (*model.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).
		Create(event).
		Error
}

func (r *eventRepository) FindByID(ctx context.Context, id uint64) (*model.Event, error) {
	var event model.Event

	err := r.db.WithContext(ctx).
		First(&event, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) FindAllPublished(ctx context.Context) ([]model.Event, error) {
	var events []model.Event

	err := r.db.WithContext(ctx).
		Where("status = ?", model.EventStatusPublished).
		Order("start_time ASC").
		Find(&events).
		Error

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).
		Save(event).
		Error
}

func (r *eventRepository) Delete(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).
		Delete(event).
		Error
}

func (r *eventRepository) FindByIDTx(ctx context.Context, tx *gorm.DB, id uint64) (*model.Event, error) {
	var event model.Event

	err := tx.WithContext(ctx).
		First(&event, id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &event, nil
}
