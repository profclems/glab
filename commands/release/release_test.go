package release

import (
	"bytes"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func Test_Release(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	cmd := NewCmdRelease(&cmdutils.Factory{})
	assert.NotNil(t, cmd.Root())
	assert.Nil(t, cmd.Execute())

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Contains(t, out, "Manage GitLab releases")
}
