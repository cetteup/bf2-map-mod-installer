//go:build unit

package internal

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstallItems(t *testing.T) {
	t.Run("successfully installs mod item", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		modDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		modFilePath, err := createTempFile(t, modDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), modDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: srcDir,
		}

		// WHEN
		err = InstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
		installedFilePath := filepath.Join(installDirPath, "mods", filepath.Base(modDirPath), filepath.Base(modFilePath))
		exists, err := fileExists(installedFilePath)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("successfully installs map item", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		mapDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		mapFilePath, err := createTempFile(t, mapDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), mapDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMap,
			SourcePath: srcDir,
			ForMod:     "bf2",
		}

		// WHEN
		err = InstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
		installedFilePath := filepath.Join(installDirPath, "mods", "bf2", "levels", filepath.Base(mapDirPath), filepath.Base(mapFilePath))
		exists, err := fileExists(installedFilePath)
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("fails if src dir does not exist", func(t *testing.T) {
		// GIVEN
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: "some-path",
		}

		// WHEN
		err = InstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.ErrorContains(t, err, "The system cannot find the file specified")
	})

	t.Run("does not fail if install dir does not exist", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		modDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		_, err = createTempFile(t, modDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), modDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: srcDir,
		}

		// WHEN
		err = InstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
	})
}

func TestUninstallItems(t *testing.T) {
	t.Run("successfully uninstalls mod item", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		modDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		modFilePath, err := createTempFile(t, installDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), modDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: srcDir,
		}

		// WHEN
		err = UninstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
		installedFilePath := filepath.Join(installDirPath, "mods", filepath.Base(modDirPath), filepath.Base(modFilePath))
		exists, err := fileExists(installedFilePath)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("successfully uninstalls map item", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		mapDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		mapFilePath, err := createTempFile(t, installDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), mapDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMap,
			SourcePath: srcDir,
			ForMod:     "bf2",
		}

		// WHEN
		err = UninstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
		installedFilePath := filepath.Join(installDirPath, "mods", "bf2", "levels", filepath.Base(mapDirPath), filepath.Base(mapFilePath))
		exists, err := fileExists(installedFilePath)
		require.NoError(t, err)
		require.False(t, exists)
	})

	t.Run("does not fail if src dir does not exist", func(t *testing.T) {
		// GIVEN
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: "some-path",
		}

		// WHEN
		err = UninstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
	})

	t.Run("does not fail if install dir does not exist", func(t *testing.T) {
		// GIVEN
		wd, err := os.Executable()
		require.NoError(t, err)
		modDirPath, err := createTempDir(t)
		require.NoError(t, err)
		installDirPath, err := createTempDir(t)
		require.NoError(t, err)
		_, err = createTempFile(t, modDirPath, "some-file")
		require.NoError(t, err)
		srcDir, err := filepath.Rel(filepath.Dir(wd), modDirPath)
		require.NoError(t, err)
		item := InstallItem{
			Type:       ItemTypeMod,
			SourcePath: srcDir,
		}

		// WHEN
		err = UninstallItems(installDirPath, []InstallItem{item})

		// THEN
		require.NoError(t, err)
	})
}

func createTempFile(t *testing.T, dir string, name string) (string, error) {
	f, err := os.CreateTemp(dir, name)
	if err != nil {
		return "", err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	t.Cleanup(func() {
		_ = os.RemoveAll(filepath.Dir(f.Name()))
	})
	return f.Name(), nil
}

func createTempDir(t *testing.T) (string, error) {
	dirPath, err := os.MkdirTemp("", "test-dir")
	if err != nil {
		return "", err
	}
	t.Cleanup(func() {
		_ = os.RemoveAll(dirPath)
	})
	return dirPath, nil
}

func fileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}
