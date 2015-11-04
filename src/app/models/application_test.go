package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticateApplication(t *testing.T) {
	querier := getModelQuerier(t)

	// Create a user with two applications.
	user, err := NewUser("user", "pass")
	require.NoError(t, err)

	_, err = querier.Save(user)
	require.NoError(t, err)

	app1, err := querier.CreateApplication(user, "app1")
	require.NoError(t, err)

	app2, err := querier.CreateApplication(user, "app2")
	require.NoError(t, err)

	// Test authentication.
	authedApp, err := querier.AuthenticateApplication("baduser", "app1", app1.AppToken)
	assert.Equal(t, errors.New("Couldn't find user with that name"), err)
	assert.Nil(t, authedApp)

	authedApp, err = querier.AuthenticateApplication("user", "app1", app1.AppToken)
	assert.NoError(t, err)
	assert.Equal(t, app1, authedApp)

	authedApp, err = querier.AuthenticateApplication("user", "app1", "bad token")
	assert.Equal(t, errors.New("Invalid application token"), err)
	assert.Nil(t, authedApp)

	authedApp, err = querier.AuthenticateApplication("user", "app2", app2.AppToken)
	assert.NoError(t, err)
	assert.Equal(t, app2, authedApp)
}
