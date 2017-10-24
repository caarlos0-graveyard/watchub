package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaults(t *testing.T) {
	cfg := Get()
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "postgres://localhost:5432/watchub?sslmode=disable", cfg.DatabaseURL)
	assert.Equal(t, "@every 1m", cfg.Schedule)
	assert.Equal(t, "super-secret-session-secret", cfg.SessionSecret)
	assert.Equal(t, "JSESSIONID", cfg.SessionName)
}
