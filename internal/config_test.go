//go:build unit

package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestLoadConfig(t *testing.T) {
	t.Run("successfully loads Config", func(t *testing.T) {
		// GIVEN
		givenConfig := Config{
			InstallItems: []InstallItem{
				{
					Type:       ItemTypeMod,
					SourcePath: "SomeMod",
				},
				{
					Type:       ItemTypeMap,
					SourcePath: "SomeMap",
				},
			},
		}
		content, err := yaml.Marshal(givenConfig)
		require.NoError(t, err)
		configFilePath, err := buildConfigFilePath()
		require.NoError(t, err)
		err = writeConfigFile(configFilePath, content)
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.Remove(configFilePath)
		})

		// WHEN
		config, err := LoadConfig()
		require.NoError(t, err)

		// THEN
		assert.Equal(t, &givenConfig, config)
	})

	t.Run("error if config file does not exist", func(t *testing.T) {
		// WHEN
		config, err := LoadConfig()

		// THEN
		require.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, config)
	})

	t.Run("error if Config file contains invalid yaml", func(t *testing.T) {
		// GIVEN
		configFilePath, err := buildConfigFilePath()
		require.NoError(t, err)
		err = writeConfigFile(configFilePath, []byte("this-is-not-valid-yaml"))
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = os.Remove(configFilePath)
		})

		// WHEN
		config, err := LoadConfig()

		// THEN
		require.ErrorContains(t, err, "cannot unmarshal")
		assert.Nil(t, config)
	})
}

func buildConfigFilePath() (string, error) {
	wd, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(wd), ConfigFilename), nil
}

func writeConfigFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0666)
}
