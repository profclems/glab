package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

		got := CheckPathExists(string(dir))
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
	filepath := file.Name()
	defer os.Remove(filepath)

	t.Run("exists", func(t *testing.T) {
		got := CheckFileExists(string(filepath))
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
		filepath := file.Name()
		defer os.Remove(filepath)

		err = BackupConfigFile(filepath)
		if err != nil {
			t.Errorf("Unexpected error = %s", err)
		}

		got := CheckFileExists(filepath + ".bak")
		assert.True(t, got)
	})
	t.Run("failure", func(t *testing.T) {
		err := BackupConfigFile("/Path/Not/Exist")
		assert.EqualError(t, err, "rename /Path/Not/Exist /Path/Not/Exist.bak: no such file or directory")
	})
}

func Test_CheckFileHasLine(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		file, err := ioutil.TempFile("", "")
		if err != nil {
			t.Skipf("Unexpected error creeating temporary file for testing = %s", err)
		}
		filepath := file.Name()
		defer os.Remove(filepath)

		file.WriteString("profclems/glab")

		got := CheckFileHasLine(filepath, "profclems/glab")
		assert.True(t, got)
	})
	t.Run("failed", func(t *testing.T) {
		t.Run("no-line-present", func(t *testing.T) {
			file, err := ioutil.TempFile("", "")
			if err != nil {
				t.Skipf("Unexpected error creeating temporary file for testing = %s", err)
			}
			filepath := file.Name()
			defer os.Remove(filepath)

			file.WriteString("profclems/glab")

			got := CheckFileHasLine(filepath, "maxice8/glab")
			assert.False(t, got)
		})
		t.Run("no-file-present", func(t *testing.T) {
			got := CheckFileHasLine("/Path/Not/Exist", "profclems/glab")
			assert.False(t, got)
		})
	})
}

func Test_ReadAndAppend(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Run("write", func(t *testing.T) {
			file, err := ioutil.TempFile("", "")
			if err != nil {
				t.Skipf("Unexpected error creating temporary file for testing = %s", err)
			}
			filepath := file.Name()
			defer os.Remove(filepath)

			err = ReadAndAppend(filepath, "profclems/glab")
			assert.NoError(t, err)
			got := CheckFileHasLine(filepath, "profclems/glab")
			assert.True(t, got)
		})
		t.Run("create-and-write", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				t.Skipf("Unexpected error creating temporary directory for testing = %s", err)
			}
			defer os.RemoveAll(dir)

			fpath := filepath.Join(dir, "file")

			err = ReadAndAppend(fpath, "profclems/glab")
			assert.NoError(t, err)
			got := CheckFileHasLine(fpath, "profclems/glab")
			assert.True(t, got)
		})
		t.Run("create-and-write-and-append", func(t *testing.T) {
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				t.Skipf("Unexpected error creating temporary directory for testing = %s", err)
			}
			defer os.RemoveAll(dir)

			fpath := filepath.Join(dir, "file")

			err = ReadAndAppend(fpath, "profclems/glab")
			assert.NoError(t, err)
			err = ReadAndAppend(fpath, "maxice8/glab")
			assert.NoError(t, err)
		})
		t.Run("write-and-append", func(t *testing.T) {
			file, err := ioutil.TempFile("", "")
			if err != nil {
				t.Skipf("Unexpected error creating temporary file for testing = %s", err)
			}
			filepath := file.Name()
			defer os.Remove(filepath)

			err = ReadAndAppend(filepath, "profclems/glab")
			assert.NoError(t, err)

			err = ReadAndAppend(filepath, "maxice8/glab")
			assert.NoError(t, err)
		})
	})
	t.Run("failed", func(t *testing.T) {
		t.Run("no-permissions", func(t *testing.T) {
			err := ReadAndAppend("/no-perm", "profclems/glab")
			assert.EqualError(t, err, "open /no-perm: permission denied")
		})
	})
}
