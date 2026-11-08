package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) SetupTest() {
	_, file, _, ok := runtime.Caller(0)
	require.True(s.T(), ok)
	projectRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
	oldDir, err := os.Getwd()
	require.NoError(s.T(), err)
	require.NoError(s.T(), os.Chdir(projectRoot))
	s.T().Cleanup(func() { _ = os.Chdir(oldDir) })
}

func (s *ConfigSuite) TestLoad_ReturnsConfiguredValues() {
	//act
	config, err := Load()
	require.NoError(s.T(), err)

	//assert
	assert.NotEmpty(s.T(), config)
	assert.Equal(s.T(), "5400", config.Database.Port)
	assert.Equal(s.T(), ":8888", config.Server.Addr)
	assert.IsType(s.T(), time.Duration(0), config.Database.InitTimeout)
}

func (s *ConfigSuite) TestDatabase_GetConnectionString() {
	//arrange
	config, err := Load()
	require.NoError(s.T(), err)

	//act
	connectionString := config.Database.GetConnectionString()

	//assert
	assert.Equal(s.T(), "postgres://postgres:telephone_pass@localhost:5400/telephone_db?sslmode=disable", connectionString)
}
