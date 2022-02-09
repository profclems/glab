package snippet

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/profclems/glab/commands/cmdutils"
)

func TestCmdSnippet_noARgs(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	assert.Nil(t, NewCmdSnippet(&cmdutils.Factory{}).Execute())

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	assert.Contains(t, out, "Use \"snippet [command] --help\" for more information about a command.\n")

}
