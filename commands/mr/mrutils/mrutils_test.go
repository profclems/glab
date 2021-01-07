package mrutils

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/xanzy/go-gitlab"
)

func Test_DisplayMR(t *testing.T) {
	testCases := []struct {
		name   string
		mr     *gitlab.MergeRequest
		output string
	}{
		{
			name: "opened",
			mr: &gitlab.MergeRequest{
				IID:          1,
				State:        "opened",
				Title:        "This is open",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/1",
			},
			output: "!1 This is open (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/1\n",
		},
		{
			name: "merged",
			mr: &gitlab.MergeRequest{
				IID:          2,
				State:        "merged",
				Title:        "This is merged",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/2",
			},
			output: "!2 This is merged (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/2\n",
		},
		{
			name: "closed",
			mr: &gitlab.MergeRequest{
				IID:          3,
				State:        "closed",
				Title:        "This is closed",
				SourceBranch: "trunk",
				WebURL:       "https://gitlab.com/profclems/glab/-/merge_requests/3",
			},
			output: "!3 This is closed (trunk)\n https://gitlab.com/profclems/glab/-/merge_requests/3\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := DisplayMR(tC.mr)
			assert.Equal(t, tC.output, got)
		})
	}
}
