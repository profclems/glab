package close

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_issueClose(t *testing.T) {
	t.Parallel()
	repo := copyTestRepo(t)

	testCases := []struct {
		Name        string
		Issue       string
		ExpectedMsg []string
		wantErr     bool
	}{
		{
			Name:        "Issue Exists",
			Issue:       "1",
			ExpectedMsg: []string{"Closing Issue...", "Issue #1 closed"},
		},
		{
			Name:        "Issue Does Not Exist",
			Issue:       "0",
			ExpectedMsg: []string{"Closing Issue", "404 Not found"},
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			cmd := exec.Command(glabBinaryPath, "issue", "close", tc.Issue)
			cmd.Dir = repo

			b, err := cmd.CombinedOutput()
			if err != nil && !tc.wantErr {
				t.Log(string(b))
				t.Fatal(err)
			}
			for _, msg := range tc.ExpectedMsg {
				assert.Contains(t, string(b), msg)
			}
		})
	}
}
