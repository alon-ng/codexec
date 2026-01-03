package users

import (
	e "codim/pkg/api/v1/errors"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Service struct {
	q *db.Queries
}

func NewService(q *db.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) Create(ctx context.Context, req CreateUserRequest) (db.User, *e.APIError) {
	u, err := s.q.CreateUser(ctx, db.CreateUserParams{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		IsVerified: req.IsVerified,
		IsAdmin:    req.IsAdmin,
	})

	if err != nil {
		if db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email") {
			return db.User{}, e.NewAPIError(err, ErrUserEmailAlreadyExists)
		}

		return db.User{}, e.NewAPIError(err, ErrUserCreationFailed)
	}

	return u, nil
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (db.User, *e.APIError) {
	u, err := s.q.UpdateUser(ctx, db.UpdateUserParams{
		Uuid:      id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})

	if err != nil {
		return db.User{}, e.NewAPIError(err, ErrUserUpdateFailed)
	}

	return u, nil
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

func (s *Service) Get(ctx context.Context, id uuid.UUID) (db.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, e.NewAPIError(err, ErrUserNotFound)
		}

		return db.User{}, e.NewAPIError(err, ErrUserGetFailed)
	}

	return u, nil
}

type ListUsersRequest struct {
	Limit  int32 `json:"limit" form:"limit,default=10" example:"10"`
	Offset int32 `json:"offset" form:"offset,default=0" example:"0"`
}

func (s *Service) List(ctx context.Context, req ListUsersRequest) ([]db.User, *e.APIError) {
	users, err := s.q.ListUsers(ctx, db.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	})

	if err != nil {
		return nil, e.NewAPIError(err, ErrUserListFailed)
	}

	return users, nil
}
