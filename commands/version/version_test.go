package version

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/profclems/glab/internal/utils"

	"github.com/stretchr/testify/assert"
)

func Test_Version(t *testing.T) {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	ios, _, _, _ := utils.IOTest()

	NewCmdVersion(ios, "v1.0.0", "2020-01-01").Execute()

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

	assert.Contains(t, out, "lab version 1.0.0 (2020-01-01)")
}
