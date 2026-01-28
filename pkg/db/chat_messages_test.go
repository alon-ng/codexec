package db_test

import (
	"codim/pkg/db"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateChatMessage(t *testing.T) {
	exercise := createRandomExercise(t, nil)
	user := createRandomUser(t)

	params := db.CreateChatMessageParams{
		ExerciseUuid:     exercise.Uuid,
		UserUuid:         user.Uuid,
		Role:             "user",
		Content:          "Test message",
		PromptTokens:     10,
		CompletionTokens: 20,
	}

	message, err := testQueries.CreateChatMessage(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, message)

	require.Equal(t, params.ExerciseUuid, message.ExerciseUuid)
	require.Equal(t, params.UserUuid, message.UserUuid)
	require.Equal(t, params.Role, message.Role)
	require.Equal(t, params.Content, message.Content)
	require.Equal(t, params.PromptTokens, message.PromptTokens)
	require.Equal(t, params.CompletionTokens, message.CompletionTokens)

	require.NotZero(t, message.Uuid)
	require.NotZero(t, message.Ts)
}

func TestListChatMessages(t *testing.T) {
	exercise := createRandomExercise(t, nil)
	user := createRandomUser(t)

	// Create multiple messages
	params1 := db.CreateChatMessageParams{
		ExerciseUuid:     exercise.Uuid,
		UserUuid:         user.Uuid,
		Role:             "user",
		Content:          "First message",
		PromptTokens:     5,
		CompletionTokens: 10,
	}
	message1, err := testQueries.CreateChatMessage(context.Background(), params1)
	require.NoError(t, err)

	params2 := db.CreateChatMessageParams{
		ExerciseUuid:     exercise.Uuid,
		UserUuid:         user.Uuid,
		Role:             "assistant",
		Content:          "Second message",
		PromptTokens:     8,
		CompletionTokens: 15,
	}
	message2, err := testQueries.CreateChatMessage(context.Background(), params2)
	require.NoError(t, err)

	// List messages
	listParams := db.ListChatMessagesParams{
		ExerciseUuid: exercise.Uuid,
		UserUuid:     user.Uuid,
		Limit:        10,
		Offset:       0,
	}

	messages, err := testQueries.ListChatMessages(context.Background(), listParams)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(messages), 2)

	var foundMessage1, foundMessage2 bool
	for _, msg := range messages {
		if msg.Uuid == message1.Uuid {
			foundMessage1 = true
			require.Equal(t, message1.Content, msg.Content)
			require.Equal(t, message1.Role, msg.Role)
		}
		if msg.Uuid == message2.Uuid {
			foundMessage2 = true
			require.Equal(t, message2.Content, msg.Content)
			require.Equal(t, message2.Role, msg.Role)
		}
	}
	require.True(t, foundMessage1)
	require.True(t, foundMessage2)
}

func TestListChatMessagesWithPagination(t *testing.T) {
	exercise := createRandomExercise(t, nil)
	user := createRandomUser(t)

	// Create multiple messages
	for i := 0; i < 5; i++ {
		params := db.CreateChatMessageParams{
			ExerciseUuid:     exercise.Uuid,
			UserUuid:         user.Uuid,
			Role:             "user",
			Content:          "Message",
			PromptTokens:     5,
			CompletionTokens: 10,
		}
		_, err := testQueries.CreateChatMessage(context.Background(), params)
		require.NoError(t, err)
	}

	// Test pagination
	listParams := db.ListChatMessagesParams{
		ExerciseUuid: exercise.Uuid,
		UserUuid:     user.Uuid,
		Limit:        2,
		Offset:       0,
	}

	messages, err := testQueries.ListChatMessages(context.Background(), listParams)
	require.NoError(t, err)
	require.LessOrEqual(t, len(messages), 2)

	// Test offset
	listParams.Offset = 2
	messages, err = testQueries.ListChatMessages(context.Background(), listParams)
	require.NoError(t, err)
	require.LessOrEqual(t, len(messages), 2)
}
