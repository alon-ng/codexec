package checkers

type CheckerType string

const (
	CheckerTypeIO   CheckerType = "io"
	CheckerTypeCode CheckerType = "code"
	CheckerTypeQuiz CheckerType = "quiz"
)

type CheckerResult struct {
	Type    CheckerType `json:"type"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}
