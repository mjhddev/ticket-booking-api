package service

import (
	"context"
	"sort"
	"time"

	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/model"
	"github.com/mjhddev/ticket-booking-api/internal/repository"
	"gorm.io/gorm"
)

type BookingService interface {
	CreateBooking(ctx context.Context, userID uint64, req dto.CreateBookingRequest) (*dto.BookingResponse, error)
}

type bookingService struct {
	db          *gorm.DB
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
}

func NewBookingService(
	db *gorm.DB,
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
) BookingService {
	return &bookingService{
		db:          db,
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
	}
}

func (s *bookingService) CreateBooking(ctx context.Context, userID uint64, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {

	if len(req.SeatIDs) == 0 {
		return nil, errs.ErrInvalidBooking
	}

	seen := make(map[uint64]struct{}, len(req.SeatIDs))

	for _, seatID := range req.SeatIDs {
		if _, exists := seen[seatID]; exists {
			return nil, errs.ErrInvalidBooking
		}

		seen[seatID] = struct{}{}
	}

	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	event, err := s.eventRepo.FindByIDTx(
		ctx,
		tx,
		req.EventID,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if event == nil {
		tx.Rollback()
		return nil, errs.ErrEventNotFound
	}

	if event.Status != model.EventStatusPublished {
		tx.Rollback()
		return nil, errs.ErrInvalidBooking
	}

	seatIDs := make([]uint64, len(req.SeatIDs))
	copy(seatIDs, req.SeatIDs)

	sort.Slice(seatIDs, func(i, j int) bool {
		return seatIDs[i] < seatIDs[j]
	})

	seats := make([]*model.Seat, 0, len(seatIDs))

	for _, seatID := range seatIDs {
		seat, err := s.bookingRepo.FindSeatForUpdate(
			ctx,
			tx,
			req.EventID,
			seatID,
		)

		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if seat == nil {
			tx.Rollback()
			return nil, errs.ErrSeatNotFound
		}

		if seat.Status != model.SeatStatusAvailable {
			tx.Rollback()
			return nil, errs.ErrSeatNotAvailable
		}

		seats = append(seats, seat)
	}

	expiresAt := time.Now().Add(10 * time.Minute)

	booking := &model.Booking{
		UserID:    userID,
		EventID:   req.EventID,
		Status:    model.BookingStatusPending,
		ExpiresAt: expiresAt,
	}

	if err := s.bookingRepo.CreateBooking(
		ctx,
		tx,
		booking,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	items := make([]model.BookingItem, 0, len(seats))

	for _, seat := range seats {
		items = append(items, model.BookingItem{
			BookingID: booking.ID,
			SeatID:    seat.ID,
			Price:     seat.Price,
		})

		seat.Status = model.SeatStatusReserved

		if err := s.bookingRepo.UpdateSeat(
			ctx,
			tx,
			seat,
		); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := s.bookingRepo.CreateBookingItems(
		ctx,
		tx,
		items,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	responseItems := make([]dto.BookingItemResponse, 0, len(items))

	for i, item := range items {
		responseItems = append(responseItems, dto.BookingItemResponse{
			ID:         item.ID,
			SeatID:     item.SeatID,
			SeatNumber: seats[i].SeatNumber,
			Price:      item.Price,
		})
	}

	return &dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		EventID:   booking.EventID,
		Status:    string(booking.Status),
		ExpiresAt: booking.ExpiresAt,
		Items:     responseItems,
		CreatedAt: booking.CreatedAt,
	}, nil
}
