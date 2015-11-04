package models

import (
	"testing"
	"time"

	"github.com/EaseApp/web-backend/src/db"
	r "github.com/dancannon/gorethink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("user", "pass")
	require.NoError(t, err)
	assert.Equal(t, "user", user.Username)
	assert.NotEqual(t, "pass", user.PasswordHash)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEmpty(t, user.APIToken)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Minute)
}

var localDBAddr = "localhost:28015"

func TestSaveAndFindUser(t *testing.T) {
	user, err := NewUser("user", "pass")
	require.NoError(t, err)
	querier := getModelQuerier(t)
	savedUser, err := querier.Save(user)
	require.NoError(t, err)
	assert.NotEmpty(t, savedUser.ID)

	foundUser := querier.Find("user")
	require.NotNil(t, foundUser)
	assertUsersEqual(t, savedUser, foundUser)
	assertUsersEqual(t, user, foundUser)
}

func TestAttemptLogin_Success(t *testing.T) {
	user, err := NewUser("user", "pass")
	require.NoError(t, err)
	querier := getModelQuerier(t)
	savedUser, err := querier.Save(user)
	require.NoError(t, err)

	loggedInUser, err := querier.AttemptLogin("user", "pass")
	require.NoError(t, err)
	assertUsersEqual(t, savedUser, loggedInUser)
}

func TestAttemptLogin_Fail(t *testing.T) {
	user, err := NewUser("user", "pass")
	require.NoError(t, err)
	querier := getModelQuerier(t)
	_, err = querier.Save(user)
	require.NoError(t, err)

	loggedInUser, err := querier.AttemptLogin("user", "badpass")
	assert.Equal(t, "Password was invalid", err.Error())
	assert.Nil(t, loggedInUser)
}

func assertUsersEqual(t *testing.T, u1, u2 *User) {
	assert.Equal(t, u1.ID, u2.ID)
	assert.Equal(t, u1.Username, u2.Username)
	assert.Equal(t, u1.PasswordHash, u2.PasswordHash)
	assert.Equal(t, u1.APIToken, u2.APIToken)
	assert.WithinDuration(t, u1.CreatedAt, u2.CreatedAt, time.Second)
}

func getModelQuerier(t *testing.T) *ModelQuerier {
	client := getDBClient(t)
	return NewModelQuerier(client.Session)
}

func getDBClient(t *testing.T) *db.Client {
	client, err := db.NewClient(localDBAddr)
	require.NoError(t, err)

	// Wait for the db to be ready.  Needed for Travis.
	r.Wait().Exec(client.Session)

	// Clear the user table for the tests.
	r.DB("test").Table("users").Delete().Run(client.Session)

	return client
}
