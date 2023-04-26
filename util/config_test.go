package util_test

import (
	"os"
	"testing"

	"github.com/mohammadrabetian/ports/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Success(t *testing.T) {
	os.Setenv("PORTS_CONFIGFILE", "config-test.toml")
	defer os.Unsetenv("PORTS_CONFIGFILE")

	config, err := util.LoadConfig("../.")
	require.NoError(t, err)

	assert.Equal(t, "development", config.Environment)
	assert.NotEmpty(t, config.FilePath)
	assert.NotEmpty(t, config.HTTPServer.Address)
	assert.NotEmpty(t, config.MySQL.DBName)
	assert.NotEmpty(t, config.MySQL.User)
	assert.NotEmpty(t, config.MySQL.Password)
	assert.NotZero(t, config.MySQL.Port)
	assert.NotEmpty(t, config.MySQL.Host)
}
