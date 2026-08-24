package service

import (
	"context"

	"github.com/mjhddev/ticket-booking-api/internal/dto"
	"github.com/mjhddev/ticket-booking-api/internal/errs"
	"github.com/mjhddev/ticket-booking-api/internal/model"
	"github.com/mjhddev/ticket-booking-api/internal/repository"
	"github.com/mjhddev/ticket-booking-api/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *utils.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *utils.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*dto.RegisterResponse, error) {

	// cek email
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errs.ErrEmailAlreadyExists
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Role:         model.RoleCustomer,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errs.ErrInvalidCredentials
	}

	if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
		return nil, errs.ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   utils.TokenTypeBearer,
	}, nil
}
