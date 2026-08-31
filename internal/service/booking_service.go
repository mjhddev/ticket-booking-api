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
	CreateBooking(
		ctx context.Context,
		userID uint64,
		req dto.CreateBookingRequest,
	) (*dto.BookingResponse, error)

	GetBookingByID(
		ctx context.Context,
		userID uint64,
		bookingID uint64,
	) (*dto.BookingResponse, error)

	CancelBooking(
		ctx context.Context,
		userID uint64,
		bookingID uint64,
	) error
}

type bookingService struct {
	db          *gorm.DB
	bookingRepo repository.BookingRepository
	eventRepo   repository.EventRepository
	seatRepo    repository.SeatRepository
}

func NewBookingService(
	db *gorm.DB,
	bookingRepo repository.BookingRepository,
	eventRepo repository.EventRepository,
	seatRepo repository.SeatRepository,
) BookingService {
	return &bookingService{
		db:          db,
		bookingRepo: bookingRepo,
		eventRepo:   eventRepo,
		seatRepo:    seatRepo,
	}
}

func (s *bookingService) CreateBooking(
	ctx context.Context,
	userID uint64,
	req dto.CreateBookingRequest,
) (*dto.BookingResponse, error) {

	if len(req.SeatIDs) == 0 {
		return nil, errs.ErrInvalidBooking
	}

	// Prevent duplicate seat IDs in one booking request.
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

	// Find event inside the same transaction.
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

	// Sort seat IDs so concurrent transactions
	// acquire locks in the same order.
	seatIDs := make([]uint64, len(req.SeatIDs))
	copy(seatIDs, req.SeatIDs)

	sort.Slice(seatIDs, func(i, j int) bool {
		return seatIDs[i] < seatIDs[j]
	})

	seats := make([]*model.Seat, 0, len(seatIDs))

	// Lock and validate every seat.
	for _, seatID := range seatIDs {
		seat, err := s.seatRepo.FindSeatForUpdate(
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

	// Create pending booking.
	booking := &model.Booking{
		UserID:    userID,
		EventID:   req.EventID,
		Status:    model.BookingStatusPending,
		ExpiresAt: time.Now().Add(10 * time.Minute),
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

	// Reserve seats and create booking items.
	for _, seat := range seats {
		items = append(items, model.BookingItem{
			BookingID: booking.ID,
			SeatID:    seat.ID,
			Price:     seat.Price,
		})

		if err := s.seatRepo.UpdateStatus(
			ctx,
			tx,
			seat.ID,
			model.SeatStatusReserved,
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

func (s *bookingService) GetBookingByID(
	ctx context.Context,
	userID uint64,
	bookingID uint64,
) (*dto.BookingResponse, error) {

	booking, err := s.bookingRepo.FindByIDWithItems(
		ctx,
		bookingID,
	)

	if err != nil {
		return nil, err
	}

	if booking == nil {
		return nil, errs.ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, errs.ErrBookingForbidden
	}

	items := make([]dto.BookingItemResponse, 0, len(booking.Items))

	for _, item := range booking.Items {
		items = append(items, dto.BookingItemResponse{
			ID:     item.ID,
			SeatID: item.SeatID,
			Price:  item.Price,
		})
	}

	return &dto.BookingResponse{
		ID:        booking.ID,
		UserID:    booking.UserID,
		EventID:   booking.EventID,
		Status:    string(booking.Status),
		ExpiresAt: booking.ExpiresAt,
		Items:     items,
		CreatedAt: booking.CreatedAt,
	}, nil
}

func (s *bookingService) CancelBooking(
	ctx context.Context,
	userID uint64,
	bookingID uint64,
) error {

	tx := s.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Lock the booking row.
	booking, err := s.bookingRepo.FindByIDForUpdate(
		ctx,
		tx,
		bookingID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	if booking == nil {
		tx.Rollback()
		return errs.ErrBookingNotFound
	}

	// Make sure the booking belongs to the current user.
	if booking.UserID != userID {
		tx.Rollback()
		return errs.ErrBookingForbidden
	}

	// Only pending bookings can be cancelled.
	if booking.Status != model.BookingStatusPending {
		tx.Rollback()
		return errs.ErrInvalidBookingStatus
	}

	items, err := s.bookingRepo.FindItemsByBookingID(
		ctx,
		tx,
		bookingID,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Release all seats.
	for _, item := range items {
		if err := s.seatRepo.UpdateStatus(
			ctx,
			tx,
			item.SeatID,
			model.SeatStatusAvailable,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Mark booking as cancelled.
	if err := s.bookingRepo.UpdateStatus(
		ctx,
		tx,
		bookingID,
		model.BookingStatusCancelled,
	); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
