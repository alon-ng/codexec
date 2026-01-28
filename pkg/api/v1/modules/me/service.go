package me

import (
	"codim/pkg/ai"
	e "codim/pkg/api/v1/errors"
	"codim/pkg/api/v1/models"
	"codim/pkg/api/v1/modules/chat"
	"codim/pkg/api/v1/modules/progress"
	"codim/pkg/db"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	q           *db.Queries
	p           *pgxpool.Pool
	progressSvc *progress.Service
	chatSvc     *chat.Service
}

func NewService(q *db.Queries, p *pgxpool.Pool, aiClient *ai.Client) *Service {
	progressSvc := progress.NewService(q, p)
	chatSvc := chat.NewService(q, p, aiClient)
	return &Service{q: q, p: p, progressSvc: progressSvc, chatSvc: chatSvc}
}

func (s *Service) Me(ctx context.Context, meUUID uuid.UUID) (models.User, *e.APIError) {
	u, err := s.q.GetUser(ctx, meUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, e.NewAPIError(err, ErrMeFailed)
		}

		return models.User{}, e.NewAPIError(err, ErrMeFailed)
	}

	return models.ToUser(u), nil
}

func (s *Service) ListUserCoursesWithProgress(ctx context.Context, meUUID uuid.UUID, req ListUserCoursesWithProgressRequest) ([]models.UserCourseWithProgress, *e.APIError) {
	return s.progressSvc.ListUserCoursesWithProgress(ctx, meUUID, req)
}

func (s *Service) GetUserCourseFull(ctx context.Context, meUUID uuid.UUID, courseUUID uuid.UUID) (models.UserCourseFull, *e.APIError) {
	return s.progressSvc.GetUserCourseFull(ctx, meUUID, courseUUID)
}

func (s *Service) GetUserExercise(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID) (models.UserExercise, *e.APIError) {
	return s.progressSvc.GetUserExercise(ctx, meUUID, exerciseUUID)
}

func (s *Service) SaveUserExerciseSubmission(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, req SaveUserExerciseSubmissionRequest) *e.APIError {
	return s.progressSvc.SaveUserExerciseSubmission(ctx, meUUID, exerciseUUID, req)
}

func (s *Service) ListChatMessages(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, req ListChatMessagesRequest) ([]models.ChatMessage, *e.APIError) {
	return s.chatSvc.ListChatMessages(ctx, exerciseUUID, meUUID, req, nil)
}

func (s *Service) SendChatMessage(ctx context.Context, meUUID uuid.UUID, exerciseUUID uuid.UUID, req SendChatMessageRequest) (models.ChatMessage, *e.APIError) {
	return s.chatSvc.SendChatMessage(ctx, exerciseUUID, meUUID, req)
}
