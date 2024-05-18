package config

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                   string
		envVariables           map[string]string
		flagArgs               []string
		expectedRunAddress     string
		expectedDatabaseURI    string
		expectedAccrualAddress string
		expectedMigrations     string
		expectedJWTTokenTTL    int
		expectedLogLevel       string
	}{
		{
			name:                   "DefaultValues",
			expectedRunAddress:     defaultRunAddress,
			expectedDatabaseURI:    defaultDatabaseURI,
			expectedAccrualAddress: defaultAccrualSystemAddress,
			expectedMigrations:     defaultMigrationsFolder,
			expectedJWTTokenTTL:    defaultJWTTokenTTLMinutes,
			expectedLogLevel:       defaultLogLevel,
		},
		{
			name: "PriorityIsEnvThenFlagThenDefault",
			envVariables: map[string]string{
				"DATABASE_URI":           "postgresql://localhost:5432/test_db",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:9090",
				"MIGRATIONS_FOLDER":      "/path/to/custom/migrations",
				"JWT_TOKEN_TTL_MINUTES":  "60",
			},
			flagArgs:               []string{"-l", "info", "-j", "30"},
			expectedRunAddress:     defaultRunAddress,
			expectedDatabaseURI:    "postgresql://localhost:5432/test_db",
			expectedAccrualAddress: "http://localhost:9090",
			expectedMigrations:     "/path/to/custom/migrations",
			expectedJWTTokenTTL:    60,
			expectedLogLevel:       "INFO",
		},
		{
			name: "UnknownFlagDoesNotCauseError",
			envVariables: map[string]string{
				"DATABASE_URI":           "postgresql://localhost:5432/test_db",
				"ACCRUAL_SYSTEM_ADDRESS": "http://localhost:9090",
				"MIGRATIONS_FOLDER":      "/path/to/custom/migrations",
				"JWT_TOKEN_TTL_MINUTES":  "60",
			},
			flagArgs:               []string{"-unknown", "something", "-j", "30"},
			expectedRunAddress:     defaultRunAddress,
			expectedDatabaseURI:    "postgresql://localhost:5432/test_db",
			expectedAccrualAddress: "http://localhost:9090",
			expectedMigrations:     "/path/to/custom/migrations",
			expectedJWTTokenTTL:    60,
			expectedLogLevel:       "INFO",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envVariables {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			os.Args = []string{"cmd"}
			os.Args = append(os.Args, tt.flagArgs...)

			err := New(context.Background())
			require.NoError(t, err)

			if Applied.RunAddress != tt.expectedRunAddress {
				t.Errorf("RunAddress: got %s, want %s", Applied.RunAddress, tt.expectedRunAddress)
			}
			if Applied.DatabaseURI != tt.expectedDatabaseURI {
				t.Errorf("DatabaseURI: got %s, want %s", Applied.DatabaseURI, tt.expectedDatabaseURI)
			}
			if Applied.AccrualSystemAddress != tt.expectedAccrualAddress {
				t.Errorf("AccrualSystemAddress: got %s, want %s", Applied.AccrualSystemAddress, tt.expectedAccrualAddress)
			}
			if Applied.MigrationsFolder != tt.expectedMigrations {
				t.Errorf("MigrationsFolder: got %s, want %s", Applied.MigrationsFolder, tt.expectedMigrations)
			}
			if Applied.JWTTokenTTLMinutes != tt.expectedJWTTokenTTL {
				t.Errorf("JWTTokenTTLMinutes: got %d, want %d", Applied.JWTTokenTTLMinutes, tt.expectedJWTTokenTTL)
			}
			if Applied.LogLevel != tt.expectedLogLevel {
				t.Errorf("LogLevel: got %s, want %s", Applied.LogLevel, tt.expectedLogLevel)
			}
		})
	}
}

func TestStartUp_GetAccrualSystemAddress(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultAccrualSystemAddress, Applied.GetAccrualSystemAddress())
}

func TestStartUp_GetJWTTokenTTLMinutes(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultJWTTokenTTLMinutes, Applied.GetJWTTokenTTLMinutes())
}

func TestStartUp_GetLogLevel(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultLogLevel, Applied.GetLogLevel())
}

func TestStartUp_GetRunAddress(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultRunAddress, Applied.GetRunAddress())
}

func TestStartUp_GetDatabaseURI(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultDatabaseURI, Applied.GetDatabaseURI())
}

func TestStartUp_GetMigrationsFolder(t *testing.T) {
	require.NoError(t, New(context.Background()))
	require.Equal(t, defaultMigrationsFolder, Applied.GetMigrationsFolder())
}
