package db_test

import (
	"codim/pkg/db"
	"context"
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) db.User {
	rnd := getRandomInt()
	params := db.CreateUserParams{
		FirstName:    fmt.Sprintf("Test First %d", rnd),
		LastName:     fmt.Sprintf("Test Last %d", rnd),
		Email:        fmt.Sprintf("test%d@example.com", rnd),
		PasswordHash: fmt.Sprintf("hashed_password_%d", rnd),
		IsVerified:   false,
		IsAdmin:      false,
	}

	user, err := testQueries.CreateUser(context.Background(), params)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.FirstName, user.FirstName)
	require.Equal(t, params.LastName, user.LastName)
	require.Equal(t, params.Email, user.Email)
	require.Equal(t, params.PasswordHash, user.PasswordHash)
	require.Equal(t, params.IsVerified, user.IsVerified)
	require.Equal(t, params.IsAdmin, user.IsAdmin)

	require.NotZero(t, user.Uuid)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.ModifiedAt)
	require.Nil(t, user.DeletedAt)
	require.Zero(t, user.Streak)
	require.Zero(t, user.Score)

	return user
}

func assertUserEqual(t *testing.T, expectedUser db.User, gotUser db.User) {
	assert.NotNil(t, gotUser)

	require.Equal(t, expectedUser.Uuid, gotUser.Uuid)
	require.Equal(t, expectedUser.FirstName, gotUser.FirstName)
	require.Equal(t, expectedUser.LastName, gotUser.LastName)
	require.Equal(t, expectedUser.Email, gotUser.Email)
	require.Equal(t, expectedUser.PasswordHash, gotUser.PasswordHash)
	require.Equal(t, expectedUser.IsVerified, gotUser.IsVerified)
	require.Equal(t, expectedUser.Streak, gotUser.Streak)
	require.Equal(t, expectedUser.Score, gotUser.Score)
	require.Equal(t, expectedUser.IsAdmin, gotUser.IsAdmin)

	require.NotZero(t, gotUser.CreatedAt)
	require.NotZero(t, gotUser.ModifiedAt)
	require.Nil(t, gotUser.DeletedAt)
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), user.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	assertUserEqual(t, user, gotUser)
}

func TestGetUserByEmail(t *testing.T) {
	user := createRandomUser(t)

	gotUser, err := testQueries.GetUserByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	assertUserEqual(t, user, gotUser)
}

func TestUpdateUser(t *testing.T) {
	user := createRandomUser(t)

	rnd := getRandomInt()
	updateParams := db.UpdateUserParams{
		Uuid:       user.Uuid,
		FirstName:  fmt.Sprintf("Updated First %d", rnd),
		LastName:   fmt.Sprintf("Updated Last %d", rnd),
		Email:      fmt.Sprintf("updated%d@example.com", rnd),
		IsVerified: true,
		IsAdmin:    true,
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, updateParams.FirstName, updatedUser.FirstName)
	require.Equal(t, updateParams.LastName, updatedUser.LastName)
	require.Equal(t, updateParams.Email, updatedUser.Email)
	require.Equal(t, updateParams.IsVerified, updatedUser.IsVerified)
	require.Equal(t, updateParams.IsAdmin, updatedUser.IsAdmin)
	require.Equal(t, user.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, user.Streak, updatedUser.Streak)
	require.Equal(t, user.Score, updatedUser.Score)

	require.NotZero(t, updatedUser.ModifiedAt)
	require.Nil(t, updatedUser.DeletedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	user := createRandomUser(t)

	rnd := getRandomInt()
	newPasswordHash := fmt.Sprintf("new_hashed_password_%d", rnd)
	updateParams := db.UpdateUserPasswordParams{
		Uuid:         user.Uuid,
		PasswordHash: newPasswordHash,
	}

	updatedUser, err := testQueries.UpdateUserPassword(context.Background(), updateParams)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	require.Equal(t, newPasswordHash, updatedUser.PasswordHash)
	require.Equal(t, user.FirstName, updatedUser.FirstName)
	require.Equal(t, user.LastName, updatedUser.LastName)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.IsVerified, updatedUser.IsVerified)
	require.Equal(t, user.IsAdmin, updatedUser.IsAdmin)

	require.NotZero(t, updatedUser.ModifiedAt)
	require.Nil(t, updatedUser.DeletedAt)
}

func TestDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.Uuid)
	require.NoError(t, err)

	gotUser, err := testQueries.GetUser(context.Background(), user.Uuid)
	require.Error(t, err)
	require.Empty(t, gotUser)
}

func TestHardDeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.HardDeleteUser(context.Background(), user.Uuid)
	require.NoError(t, err)

	gotUser, err := testQueries.GetUser(context.Background(), user.Uuid)
	require.Error(t, err)
	require.Empty(t, gotUser)
}

func TestUndeleteUser(t *testing.T) {
	user := createRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), user.Uuid)
	require.NoError(t, err)

	gotUser, err := testQueries.GetUser(context.Background(), user.Uuid)
	require.Error(t, err)
	require.Empty(t, gotUser)

	err = testQueries.UndeleteUser(context.Background(), user.Uuid)
	require.NoError(t, err)

	gotUser, err = testQueries.GetUser(context.Background(), user.Uuid)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	assertUserEqual(t, user, gotUser)
}

func TestCreateUserConflict(t *testing.T) {
	user := createRandomUser(t)

	params := db.CreateUserParams{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		IsVerified:   user.IsVerified,
		IsAdmin:      user.IsAdmin,
	}

	_, err := testQueries.CreateUser(context.Background(), params)
	require.Error(t, err)
	require.True(t, db.IsDuplicateKeyErrorWithConstraint(err, "uq_users_email"))
}

func TestCountUsers(t *testing.T) {
	initialCount, err := testQueries.CountUsers(context.Background())
	require.NoError(t, err)

	user1 := createRandomUser(t)
	_ = createRandomUser(t)

	count, err := testQueries.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+2, count)

	err = testQueries.DeleteUser(context.Background(), user1.Uuid)
	require.NoError(t, err)

	count, err = testQueries.CountUsers(context.Background())
	require.NoError(t, err)
	require.Equal(t, initialCount+1, count)
}

func TestListUsers(t *testing.T) {
	user1 := createRandomUser(t)
	user2 := createRandomUser(t)

	params := db.ListUsersParams{
		Limit:  10,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(context.Background(), params)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(users), 2)

	var foundUser1, foundUser2 bool
	for _, user := range users {
		if user.Uuid == user1.Uuid {
			foundUser1 = true
			assertUserEqual(t, user1, user)
		}
		if user.Uuid == user2.Uuid {
			foundUser2 = true
			assertUserEqual(t, user2, user)
		}
	}
	require.True(t, foundUser1)
	require.True(t, foundUser2)
}
