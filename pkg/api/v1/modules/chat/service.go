package chat

import (
	"codim/pkg/ai"
	"codim/pkg/db"
	"context"
	"errors"
	"fmt"

	e "codim/pkg/api/v1/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/openai/openai-go"
)

type Service struct {
	q        *db.Queries
	p        *pgxpool.Pool
	aiClient *ai.Client
}

func NewService(q *db.Queries, p *pgxpool.Pool, aiClient *ai.Client) *Service {
	return &Service{q: q, p: p, aiClient: aiClient}
}

func ToChatMessage(d db.ChatMessage) ChatMessage {
	return ChatMessage{
		Uuid:         d.Uuid,
		Ts:           d.Ts,
		ExerciseUuid: d.ExerciseUuid,
		UserUuid:     d.UserUuid,
		Role:         d.Role,
		Content:      d.Content,
	}
}

func (s *Service) ListChatMessages(ctx context.Context, exerciseUUID uuid.UUID, userUUID uuid.UUID, req ListChatMessagesRequest, qtx *db.Queries) ([]ChatMessage, *e.APIError) {
	if qtx == nil {
		qtx = s.q
	}

	messages, err := qtx.ListChatMessages(ctx, db.ListChatMessagesParams{
		ExerciseUuid: exerciseUUID,
		UserUuid:     userUUID,
		Limit:        req.Limit,
		Offset:       req.Offset,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []ChatMessage{}, nil
		}

		return nil, e.NewAPIError(err, ErrListChatMessagesFailed)
	}

	m := make([]ChatMessage, len(messages))
	for i, message := range messages {
		m[i] = ToChatMessage(message)
	}

	return m, nil
}

func (s *Service) SendChatMessage(
	ctx context.Context,
	exerciseUUID uuid.UUID,
	userUUID uuid.UUID,
	req SendChatMessageRequest,
) (ChatMessage, *e.APIError) {
	tx, err := s.p.Begin(ctx)
	if err != nil {
		return ChatMessage{}, e.NewAPIError(err, ErrSendChatMessageFailed)
	}
	defer tx.Rollback(ctx)

	qtx := s.q.WithTx(tx)

	_, err = qtx.CreateChatMessage(ctx, db.CreateChatMessageParams{
		ExerciseUuid: exerciseUUID,
		UserUuid:     userUUID,
		Role:         "user",
		Content:      req.Content,
	})
	if err != nil {
		return ChatMessage{}, e.NewAPIError(err, ErrSendChatMessageFailed)
	}

	messageHistory, apiErr := s.ListChatMessages(ctx, exerciseUUID, userUUID, ListChatMessagesRequest{
		Limit:  6,
		Offset: 0,
	}, qtx)

	if apiErr != nil {
		return ChatMessage{}, apiErr
	}

	openaiMessages := s.constructOpenaiMessages(req, messageHistory)

	r, err := s.aiClient.SendMessage(ctx, openaiMessages)
	if err != nil {
		return ChatMessage{}, e.NewAPIError(err, ErrSendChatMessageFailed)
	}

	message, err := qtx.CreateChatMessage(ctx, db.CreateChatMessageParams{
		ExerciseUuid:     exerciseUUID,
		UserUuid:         userUUID,
		Role:             "assistant",
		Content:          r.Content,
		PromptTokens:     int32(r.PromptTokens),
		CompletionTokens: int32(r.CompletionTokens),
	})
	if err != nil {
		return ChatMessage{}, e.NewAPIError(err, ErrSendChatMessageFailed)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ChatMessage{}, e.NewAPIError(err, ErrSendChatMessageFailed)
	}

	return ToChatMessage(message), nil
}

func (s *Service) constructOpenaiMessages(req SendChatMessageRequest, messageHistory []ChatMessage) []openai.ChatCompletionMessageParamUnion {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(messageHistory)+2)
	openaiMessages[0] = openai.SystemMessage(`You are an expert AI Programming Tutor on a code learning platform. Your goal is to help the user complete the exercise described in <ExerciseInstructions> based on their current progress in <UserCode>.

**CORE DIRECTIVES:**
1.  **Guide, Do Not Solve:** Under no circumstances should you write the full solution or fix the code for the user.
2.  **Socratic Debugging:** If there is an error, explain *why* it is happening or ask a leading question (e.g., "Check how you are updating the loop variable") rather than providing the fix.
3.  **Stay in Scope:** Focus strictly on the concepts required for the current exercise. Do not introduce advanced topics unless necessary.
4.  **Tone:** Be concise, encouraging, and technically precise.
5.  **Language:** Always respond in the same language as the user's question.
6.  **Straight to the point:** Be concise and to the point. Don't beat around the bush.

Write the response in markdown format:
- Use code blocks to format code examples.
- Use bullet points to list steps or key points.
- Use bold or italic text to highlight important information.
- Use headings to structure the response.

Reference the user's code lines explicitly when offering feedback.`)
	openaiMessages[1] = openai.UserMessage(fmt.Sprintf("<ExerciseInstructions>%s</ExerciseInstructions>\n<UserCode>%s</UserCode>", req.ExerciseInstructions, req.Code))
	for i, message := range messageHistory {
		switch message.Role {
		case "user":
			openaiMessages[i+2] = openai.UserMessage(message.Content)
		case "assistant":
			openaiMessages[i+2] = openai.AssistantMessage(message.Content)
		}
	}
	return openaiMessages
}
