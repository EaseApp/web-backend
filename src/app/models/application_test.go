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
			"/howdy",
			"/",
			map[string]interface{}{"yeah": "wassuuuup", "multiple": []int{1, 2}},
			map[string]interface{}{"howdy": map[string]interface{}{"yeah": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}}},
		},
		{
			"/hello/world/hi",
			"/hello/world/hi",
			10,
			float64(10),
		},
		{
			"/a/b/c",
			"/a/b",
			10,
			map[string]interface{}{"c": float64(10)},
		},
		{
			"/yes/no/maybe",
			"/yes/no/maybe",
			map[string]interface{}{"hello": "wassuuuup", "multiple": []int{1, 2}},
			map[string]interface{}{"hello": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}},
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
		p, err := lib.ParsePath("/hello")
		querier.SaveApplicationData(app, p, nil)
	}
}

func TestDeleteApplicationData(t *testing.T) {
	querier := getModelQuerier(t)

	// Create a user with an application.
	user, err := NewUser("user", "pass")
	require.NoError(t, err)

	_, err = querier.Save(user)
	require.NoError(t, err)

	app, err := querier.CreateApplication(user, "app1")
	require.NoError(t, err)

	// Create default data to be deleted.
	path, err := lib.ParsePath("/hello/world")
	require.NoError(t, err)
	path2, err := lib.ParsePath("/oh")
	require.NoError(t, err)

	err = querier.SaveApplicationData(app, path, map[string]interface{}{"yes": "wassuuuup", "multiple": []int{1, 2}})
	require.NoError(t, err)
	err = querier.SaveApplicationData(app, path2, 5)
	require.NoError(t, err)

	testcases := []struct {
		deletePath       string
		expectedRootData interface{}
	}{
		{
			"/hello/world/i/dont/exist",
			map[string]interface{}{"oh": float64(5), "hello": map[string]interface{}{"world": map[string]interface{}{"yes": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}}}},
		},
		{
			"/oh",
			map[string]interface{}{"hello": map[string]interface{}{"world": map[string]interface{}{"yes": "wassuuuup", "multiple": []interface{}{float64(1), float64(2)}}}},
		},
		{
			"/hello/world/multiple",
			map[string]interface{}{"hello": map[string]interface{}{"world": map[string]interface{}{"yes": "wassuuuup"}}},
		},
		{
			"/hello/world",
			map[string]interface{}{"hello": map[string]interface{}{}},
		},
		{
			"/",
			map[string]interface{}{},
		},
	}

	rootPath, err := lib.ParsePath("/")
	require.NoError(t, err)
	for _, testcase := range testcases {
		path, err = lib.ParsePath(testcase.deletePath)
		require.NoError(t, err)

		err = querier.DeleteApplicationData(app, path)
		assert.NoError(t, err)

		data, err := querier.ReadApplicationData(app, rootPath)
		assert.NoError(t, err)
		assert.Equal(t, testcase.expectedRootData, data)
	}
}
