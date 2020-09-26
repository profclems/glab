package ci

import (
	"bytes"
	"github.com/profclems/glab/commands/cmdutils"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNewCmdCI(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	assert.Nil(t, NewCmdCI(&cmdutils.Factory{}).Execute())

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

	assert.Contains(t, out, "Work with GitLab CI pipelines and jobs\n")

}
