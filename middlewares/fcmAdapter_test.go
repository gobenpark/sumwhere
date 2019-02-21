package middlewares

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFireBaseApp(t *testing.T) {
	app, err := NewFireBaseApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)
}

func TestFireBaseAppAdapter_SendMessage(t *testing.T) {
	app, err := NewFireBaseApp()
	assert.NoError(t, err)
	assert.NotNil(t, app)

	err = app.SendMessage("test", "test", "eIbIzANXYYE:APA91bGrXQg5ns5aQK4m979ygcqwafKI0Hxzi8fK8Z-_UHROpjrGCqgYcsljRdSkCZE6OjMJVnKFbwUqnNhfvhLAUaToKgKD4gdALgWGGtt-S8Ev7FCDFeUg1T3knfGET-dOvMsPKtLH")
	assert.NoError(t, err)
}
