package internal

import (
	"fmt"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

const (
	modPath   = "mods"
	levelPath = "levels"
)

func InstallItems(bf2InstallPath string, items []InstallItem) error {
	for _, item := range items {
		srcPath, destPath, err := buildPaths(bf2InstallPath, item)
		if err != nil {
			return err
		}

		err = cp.Copy(srcPath, destPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func UninstallItems(bf2InstallPath string, items []InstallItem) error {
	for _, item := range items {
		_, destPath, err := buildPaths(bf2InstallPath, item)
		if err != nil {
			return err
		}

		err = os.RemoveAll(destPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildPaths(bf2InstallPath string, item InstallItem) (string, string, error) {
	wd, err := os.Executable()
	if err != nil {
		return "", "", err
	}

	srcPath := filepath.Join(filepath.Dir(wd), item.SourcePath)
	var destPath string
	switch item.Type {
	case ItemTypeMod:
		destPath = filepath.Join(bf2InstallPath, modPath, filepath.Base(item.SourcePath))
	case ItemTypeMap:
		destPath = filepath.Join(bf2InstallPath, modPath, item.ForMod, levelPath, filepath.Base(item.SourcePath))
	default:
		return "", "", fmt.Errorf("unsupported item type: %s", item.Type)
	}

	return srcPath, destPath, nil
}
