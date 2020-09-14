package config

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func getOldGlobalConfigDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(usr.HomeDir, ".glab-cli", "config"), nil
}

func migrateGlobalConfigDir() error {
	// check if xdg directory exists, bail if so.
	newConfigDir, err := getXdgGlobalConfigDir()
	if err != nil {
		return err
	}
	if CheckPathExists(newConfigDir) {
		return nil
	}

	// check if old config dir exists, or there's nothing to migrate.
	oldConfigDir, err := getOldGlobalConfigDir()
	if err != nil {
		return err
	}
	if !CheckPathExists(oldConfigDir) {
		return nil
	}

	// do the migration
	log.Println("Migrating config dir to XDG_CONFIG_HOME.")

	// First make sure parent directory exists
	if !CheckPathExists(filepath.Join(newConfigDir, "..")) {
		if err := os.MkdirAll(filepath.Join(newConfigDir, ".."), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create new parent config dir: %v", err)
		}
	}

	if err := os.Rename(oldConfigDir, newConfigDir); err != nil {
		return fmt.Errorf("failed to move config dir: %v", err)
	}

	// cleanup: remove parent directory tree of oldConfigDir if empty
	return os.Remove(filepath.Join(oldConfigDir, ".."))
}

// migrateOldAliasFile renames alias file from old aliases.format to aliases.yml
func migrateOldAliasFile() error {
	oldAliasFile := filepath.Join(globalPathDir, "aliases.format")
	if CheckFileExists(oldAliasFile) {
		if err := os.Rename(oldAliasFile, aliasFile); err == nil {
			return fmt.Errorf("failed to rename aliases.format to aliases.yml: %v", err)
		}
		return os.Remove(oldAliasFile)
	}
	return nil
}
