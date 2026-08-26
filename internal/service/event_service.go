package service

import (
	"context"

	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/model"
	"github.com/mjhddev/ticket-booking-api/internal/repository"
)

type EventService interface {
	Create(ctx context.Context, organizerID uint64, req dto.CreateEventRequest) (*dto.EventResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.EventResponse, error)
	GetPublished(ctx context.Context) ([]dto.EventResponse, error)
	Update(ctx context.Context, organizerID uint64, id uint64, req dto.UpdateEventRequest) (*dto.EventResponse, error)
	Delete(ctx context.Context, organizerID uint64, id uint64) error
	Publish(ctx context.Context, userID uint64, isAdmin bool, id uint64) error
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{
		eventRepo: eventRepo,
	}
}

func (s *eventService) Create(ctx context.Context, organizerID uint64, req dto.CreateEventRequest) (*dto.EventResponse, error) {

	if !req.EndTime.After(req.StartTime) {
		return nil, errs.ErrInvalidEventTime
	}

	event := &model.Event{
		OrganizerID: organizerID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      model.EventStatusDraft,
	}

	if err := s.eventRepo.Create(ctx, event); err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
}

func (s *eventService) GetByID(ctx context.Context, id uint64) (*dto.EventResponse, error) {
	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errs.ErrEventNotFound
	}

	return toEventResponse(event), nil
}

func (s *eventService) GetPublished(ctx context.Context) ([]dto.EventResponse, error) {

	events, err := s.eventRepo.FindAllPublished(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.EventResponse, 0, len(events))

	for _, event := range events {
		responses = append(
			responses,
			*toEventResponse(&event),
		)
	}

	return responses, nil
}

func (s *eventService) Update(ctx context.Context, organizerID uint64, id uint64, req dto.UpdateEventRequest) (*dto.EventResponse, error) {

	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errs.ErrEventNotFound
	}

	if event.OrganizerID != organizerID {
		return nil, errs.ErrEventNotOwned
	}

	if !req.EndTime.After(req.StartTime) {
		return nil, errs.ErrInvalidEventTime
	}

	event.Title = req.Title
	event.Description = req.Description
	event.Location = req.Location
	event.StartTime = req.StartTime
	event.EndTime = req.EndTime

	if err := s.eventRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return toEventResponse(event), nil
}

func (s *eventService) Delete(ctx context.Context, organizerID uint64, id uint64) error {

	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if event == nil {
		return errs.ErrEventNotFound
	}

	if event.OrganizerID != organizerID {
		return errs.ErrEventNotOwned
	}

	return s.eventRepo.Delete(ctx, event)
}

func (s *eventService) Publish(ctx context.Context, userID uint64, isAdmin bool, id uint64) error {

	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if event == nil {
		return errs.ErrEventNotFound
	}

	if !isAdmin && event.OrganizerID != userID {
		return errs.ErrEventNotOwned
	}

	if event.Status != model.EventStatusDraft {
		return errs.ErrEventCannotBePublished
	}

	event.Status = model.EventStatusPublished

	return s.eventRepo.Update(ctx, event)
}

func toEventResponse(event *model.Event) *dto.EventResponse {
	return &dto.EventResponse{
		ID:          event.ID,
		OrganizerID: event.OrganizerID,
		Title:       event.Title,
		Description: event.Description,
		Location:    event.Location,
		StartTime:   event.StartTime,
		EndTime:     event.EndTime,
		Status:      string(event.Status),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}
