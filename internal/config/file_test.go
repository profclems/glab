package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_CheckPathExists(t *testing.T) {
	t.Run("exists", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "")
		if err != nil {
			t.Skipf("unexpected error creating temporary directory for testing = %s", err)
		}
		defer os.Remove(dir)

		got := CheckPathExists(dir)
		assert.True(t, got)
	})
	t.Run("doesnt-exist", func(t *testing.T) {
		got := CheckPathExists("/Path/Not/Exist")
		assert.False(t, got)
	})
}

func Test_CheckFileExists(t *testing.T) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Skipf("Unexpected error creating temporary file for testing = %s", err)
	}
	fPath := file.Name()
	defer os.Remove(fPath)

	t.Run("exists", func(t *testing.T) {
		got := CheckFileExists(fPath)
		assert.True(t, got)
	})

	t.Run("doesnt-exist", func(t *testing.T) {
		got := CheckFileExists("/Path/Not/Exist")
		assert.False(t, got)
	})
}

func Test_BackupConfigFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		file, err := ioutil.TempFile("", "")
		if err != nil {
			t.Skipf("Unexpected error creating temporary file for testing = %s", err)
		}
		fPath := file.Name()
		defer os.Remove(fPath)

		err = BackupConfigFile(fPath)
		if err != nil {
			t.Errorf("Unexpected error = %s", err)
		}

		got := CheckFileExists(fPath + ".bak")
		assert.True(t, got)
	})
	t.Run("failure", func(t *testing.T) {
		err := BackupConfigFile("/Path/Not/Exist")
		assert.EqualError(t, err, "rename /Path/Not/Exist /Path/Not/Exist.bak: no such file or directory")
	})
}
