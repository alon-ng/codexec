package me

import (
	"codim/pkg/api/v1/modules/chat"
	"codim/pkg/api/v1/modules/progress"
)

type ListUserCoursesWithProgressRequest = progress.ListUserCoursesWithProgressRequest
type UserExerciseStatus = progress.UserExerciseStatus
type UserLessonStatus = progress.UserLessonStatus
type UserCourseFull = progress.UserCourseFull
type UserCourseWithProgress = progress.UserCourseWithProgress
type UserExercise = progress.UserExercise
type SaveUserExerciseSubmissionRequest = progress.SaveUserExerciseSubmissionRequest

type ListChatMessagesRequest = chat.ListChatMessagesRequest
type SendChatMessageRequest = chat.SendChatMessageRequest

type RunUserExerciseCodeSubmissionRequest struct {
	Submission progress.UserExerciseSubmissionCode
}

type RunUserExerciseCodeSubmissionResponse struct {
	Stdout   string  `json:"stdout"`
	Stderr   string  `json:"stderr"`
	ExitCode int     `json:"exit_code"`
	Time     float64 `json:"time"`
	Memory   int64   `json:"memory"`
	CPU      float64 `json:"cpu"`
}
