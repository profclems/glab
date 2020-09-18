package config

import (
	"fmt"
	"io"
	"os"
	"path"
)

func StubBackupConfig() func() {
	orig := BackupConfigFile
	BackupConfigFile = func(_ string) error {
		return nil
	}

	return func() {
		BackupConfigFile = orig
	}
}

func StubWriteConfig(wc io.Writer, wh io.Writer) func() {
	orig := WriteConfigFile
	WriteConfigFile = func(fn string, data []byte) error {
		switch path.Base(fn) {
		case "config.yml":
			_, err := wc.Write(data)
			return err
		case "aliases.yml":
			_, err := wh.Write(data)
			return err
		default:
			return fmt.Errorf("write to unstubbed file: %q", fn)
		}
	}
	return func() {
		WriteConfigFile = orig
	}
}

func StubConfig(main, aliases string) func() {
	orig := ReadConfigFile
	ReadConfigFile = func(fn string) ([]byte, error) {
		switch path.Base(fn) {
		case "config.yml":
			if main == "" {
				return []byte(nil), os.ErrNotExist
			} else {
				return []byte(main), nil
			}
		case "aliases.yml":
			if aliases == "" {
				return []byte(nil), os.ErrNotExist
			} else {
				return []byte(aliases), nil
			}
		default:
			return []byte(nil), fmt.Errorf("read from unstubbed file: %q", fn)
		}

	}
	return func() {
		ReadConfigFile = orig
	}
}
