package users

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/models"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q *db.Queries
	p *pgxpool.Pool
}

func NewService(q *db.Queries, p *pgxpool.Pool) *Service {
	return &Service{q: q, p: p}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (models.User, *e.APIError) {
	u, err := s.q.CreateUser(ctx, db.CreateUserParams{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		IsVerified: req.IsVerified,
		IsAdmin:    req.IsAdmin,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email") {
			return models.User{}, e.NewAPIError(err, ErrUserEmailAlreadyExists)
		}

		return models.User{}, e.NewAPIError(err, ErrUserCreationFailed)
	}

	return models.ToUser(u), nil
}

func (s *Service) Update(ctx context.Context, req UpdateUserRequest) (models.User, *e.APIError) {
	u, err := s.q.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      req.Uuid,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		return models.User{}, e.NewAPIError(err, ErrUserUpdateFailed)
	}

	return models.ToUser(u), nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.DeleteUser(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrUserDeleteFailed)
	}

	return nil
}

func (s *Service) Restore(ctx context.Context, id uuid.UUID) *e.APIError {
	err := s.q.UndeleteUser(ctx, id)
	if err != nil {
		return e.NewAPIError(err, ErrUserRestoreFailed)
	}

	return nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (models.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, e.NewAPIError(err, ErrUserNotFound)
		}

		return models.User{}, e.NewAPIError(err, ErrUserGetFailed)
	}

	return models.ToUser(u), nil
}

func (s *Service) List(ctx context.Context, req ListUsersRequest) ([]models.User, *e.APIError) {
	users, err := s.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrUserListFailed)
	}

	result := make([]models.User, len(users))
	for i, u := range users {
		result[i] = models.ToUser(u)
	}

	return result, nil
}
