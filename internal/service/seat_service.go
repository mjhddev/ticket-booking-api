package service

import (
	"context"
	"strings"

	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/model"
	"github.com/mjhddev/ticket-booking-api/internal/repository"
)

type SeatService interface {
	CreateSeats(ctx context.Context, organizerID uint64, eventID uint64, req dto.CreateSeatsRequest) ([]dto.SeatResponse, error)
	GetSeatsByEvent(ctx context.Context, eventID uint64) ([]dto.SeatResponse, error)
}

type seatService struct {
	seatRepo  repository.SeatRepository
	eventRepo repository.EventRepository
}

func NewSeatService(seatRepo repository.SeatRepository, eventRepo repository.EventRepository) SeatService {
	return &seatService{
		seatRepo:  seatRepo,
		eventRepo: eventRepo,
	}
}

func (s *seatService) CreateSeats(ctx context.Context, organizerID uint64, eventID uint64, req dto.CreateSeatsRequest) ([]dto.SeatResponse, error) {

	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errs.ErrEventNotFound
	}

	if event.OrganizerID != organizerID {
		return nil, errs.ErrEventNotOwned
	}

	seats := make([]model.Seat, 0, len(req.Seats))
	seen := make(map[string]struct{}, len(req.Seats))

	for _, reqSeat := range req.Seats {
		seatNumber := strings.TrimSpace(reqSeat.SeatNumber)

		if seatNumber == "" {
			continue
		}

		seatNumber = strings.ToUpper(seatNumber)

		if _, exists := seen[seatNumber]; exists {
			return nil, errs.ErrSeatAlreadyExists
		}

		seen[seatNumber] = struct{}{}

		seats = append(seats, model.Seat{
			EventID:    eventID,
			SeatNumber: seatNumber,
			Price:      reqSeat.Price,
			Status:     model.SeatStatusAvailable,
		})
	}

	if len(seats) == 0 {
		return nil, errs.ErrSeatNotFound
	}

	if err := s.seatRepo.CreateBatch(ctx, seats); err != nil {
		return nil, err
	}

	responses := make([]dto.SeatResponse, 0, len(seats))

	for _, seat := range seats {
		responses = append(responses, dto.SeatResponse{
			ID:         seat.ID,
			EventID:    seat.EventID,
			SeatNumber: seat.SeatNumber,
			Price:      seat.Price,
			Status:     string(seat.Status),
		})
	}

	return responses, nil
}

func (s *seatService) GetSeatsByEvent(ctx context.Context, eventID uint64) ([]dto.SeatResponse, error) {

	event, err := s.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errs.ErrEventNotFound
	}

	seats, err := s.seatRepo.FindByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.SeatResponse, 0, len(seats))

	for _, seat := range seats {
		responses = append(responses, dto.SeatResponse{
			ID:         seat.ID,
			EventID:    seat.EventID,
			SeatNumber: seat.SeatNumber,
			Price:      seat.Price,
			Status:     string(seat.Status),
		})
	}

	return responses, nil
}
