package me

import (
	"codim/pkg/api/v1/models"
	"codim/pkg/api/v1/modules/chat"
	"codim/pkg/api/v1/modules/progress"
)

type ListUserCoursesWithProgressRequest = progress.ListUserCoursesWithProgressRequest
type SaveUserExerciseSubmissionRequest = progress.SaveUserExerciseSubmissionRequest
type ListChatMessagesRequest = chat.ListChatMessagesRequest
type SendChatMessageRequest = chat.SendChatMessageRequest

type RunUserExerciseCodeSubmissionRequest struct {
	Submission models.UserExerciseSubmissionCode
}
