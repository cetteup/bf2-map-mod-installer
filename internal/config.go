package internal

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ItemType string

const (
	ItemTypeMod = "mod"
	ItemTypeMap = "map"

	ConfigFilename = "Config.yaml"
)

type Config struct {
	InstallItems []InstallItem `yaml:"items"`
}

func (c *Config) GetItemOfType(itemType ItemType) []InstallItem {
	mods := make([]InstallItem, 0, len(c.InstallItems))
	for _, item := range c.InstallItems {
		if item.Type == itemType {
			mods = append(mods, item)
		}
	}
	return mods
}

type InstallItem struct {
	Type       ItemType `yaml:"type"`
	SourcePath string   `yaml:"src_path"`
	ForMod     string   `yaml:"for_mod"`
}

func LoadConfig() (*Config, error) {
	wd, err := os.Executable()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(filepath.Dir(wd), ConfigFilename)
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
