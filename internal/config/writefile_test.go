package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WriteFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Skipf("unexpected error while creating temporay directory = %s", err)
	}
	defer os.RemoveAll(dir)

	fpath := filepath.Join(dir, "test-file")

	err = WriteFile(fpath, []byte("profclems/glab"), 0644)
	if err != nil {
		t.Errorf("unexpected error = %s", err)
	}

	result, err := ioutil.ReadFile(fpath)
	if err != nil {
		t.Errorf("failed to read file %q due to %q", fpath, err)
	}
	assert.Equal(t, "profclems/glab", string(result))

	permissions, err := os.Stat(fpath)
	if err != nil {
		t.Errorf("failed to get stats for file %q due to %q", fpath, err)
	}
	// TODO:
	assert.Equal(t, "0644", fmt.Sprintf("%04o", permissions.Mode()))
}
