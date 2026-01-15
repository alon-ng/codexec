package auth

import (
	authProvider "codim/pkg/api/auth"
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q           *db.Queries
	p           *pgxpool.Pool
	authService *authProvider.Provider
}

func NewService(q *db.Queries, p *pgxpool.Pool, authService *authProvider.Provider) *Service {
	return &Service{q: q, p: p, authService: authService}
}

func (s *Service) Signup(ctx context.Context, req SignupRequest) (AuthResponse, *e.APIError) {
	hashedPassword := s.authService.HashPassword(req.Password)

	user, err := s.q.CreateUser(ctx, db.CreateUserParams{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsVerified:   false,
		IsAdmin:      false,
	})
	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email") {
			return AuthResponse{}, e.NewAPIError(err, ErrEmailAlreadyExists)
		}
		return AuthResponse{}, e.NewAPIError(err, ErrSignupFailed)
	}

	token, err := s.authService.GenerateToken(user.Uuid)
	if err != nil {
		return AuthResponse{}, e.NewAPIError(err, ErrTokenGenerationFailed)
	}

	return AuthResponse{Token: token, User: user}, nil
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (AuthResponse, *e.APIError) {
	user, err := s.q.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResponse{}, e.NewAPIError(err, ErrInvalidCredentials)
		}
		return AuthResponse{}, e.NewAPIError(err, ErrLoginFailed)
	}

	if !s.authService.VerifyPassword(req.Password, user.PasswordHash) {
		return AuthResponse{}, e.NewAPIError(errors.New("invalid password"), ErrInvalidCredentials)
	}

	token, err := s.authService.GenerateToken(user.Uuid)
	if err != nil {
		return AuthResponse{}, e.NewAPIError(err, ErrTokenGenerationFailed)
	}

	return AuthResponse{Token: token, User: user}, nil
}
