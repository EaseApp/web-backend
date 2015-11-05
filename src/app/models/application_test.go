package models

import (
	"errors"
	"testing"

	"github.com/EaseApp/web-backend/src/lib"
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

func TestSaveAndReadApplicationData(t *testing.T) {
	querier := getModelQuerier(t)

	// Create a user with an application.
	user, err := NewUser("user", "pass")
	require.NoError(t, err)

	_, err = querier.Save(user)
	require.NoError(t, err)

	app, err := querier.CreateApplication(user, "app1")
	require.NoError(t, err)

	testcases := []struct {
		writePath        string
		readPath         string
		writeData        interface{}
		expectedReadData interface{}
	}{
		{
			"/hello/world/hi",
			"/hello/world/hi",
			10,
			float64(10),
		},
		{
			"/hello/world/hi",
			"/hello/world",
			10,
			map[string]interface{}{"hi": float64(10)},
		},
		{
			"/hello/world/hi",
			"/hello/world",
			10,
			map[string]interface{}{"hi": float64(10)},
		},
		{
			"/hello/world/hi",
			"/hello/world/hi",
			map[string]interface{}{"hello": "wassuuuup", "multiple": []int{1, 2}},
			map[string]interface{}{"hello": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}},
		},
		{
			"/hello",
			"/",
			map[string]interface{}{"yeah": "wassuuuup", "multiple": []int{1, 2}},
			map[string]interface{}{"hello": map[string]interface{}{"yeah": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}}},
		},
	}

	for _, testcase := range testcases {
		writePath, err := lib.ParsePath(testcase.writePath)
		require.NoError(t, err)

		err = querier.SaveApplicationData(app, writePath, testcase.writeData)
		assert.NoError(t, err)

		readPath, err := lib.ParsePath(testcase.readPath)
		data, err := querier.ReadApplicationData(app, readPath)
		assert.Equal(t, testcase.expectedReadData, data)
	}
}
