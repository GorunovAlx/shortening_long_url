package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Set environment variables.
func SetEnvs(m map[string]string) {
	for name, value := range m {
		os.Setenv(name, value)
	}
}

// Check environment variables.
func CheckEnvs(t *testing.T, m map[string]string) {
	for name, value := range m {
		actual := os.Getenv(name)
		assert.Equal(t, value, actual)
	}
}

// Unset environment variables.
func UnsetEnvs(m map[string]string) {
	for name := range m {
		os.Unsetenv(name)
	}
}

// Testing environment variables.
func TestEnvSetConfig(t *testing.T) {
	tests := []struct {
		name string
		vars map[string]string
	}{
		{
			name: "test environment variables with all vars",
			vars: map[string]string{
				"SERVER_ADDRESS":    ":34567",
				"BASE_URL":          "http://localhost:34567",
				"FILE_STORAGE_PATH": "file.txt",
			},
		},
		{
			name: "test environment variables with BASE_URL and FILE_STORAGE_PATH vars",
			vars: map[string]string{
				"BASE_URL":          "http://localhost:34567",
				"FILE_STORAGE_PATH": "file.txt",
			},
		},
		{
			name: "test environment variables with FILE_STORAGE_PATH var",
			vars: map[string]string{
				"FILE_STORAGE_PATH": "file.txt",
			},
		},
		{
			name: "test environment variables without vars",
			vars: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetEnvs(tt.vars)
			SetConfig()
			CheckEnvs(t, tt.vars)
			UnsetEnvs(tt.vars)
		})
	}
}
