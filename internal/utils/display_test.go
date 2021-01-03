package utils

import (
	"testing"

	"github.com/alecthomas/assert"
)

func Test_Indent(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		indent string
		output string
	}{
		{
			name:   "4-spaces",
			input:  "Hello Glab",
			indent: "    ",
			output: "    Hello Glab",
		},
		{
			name:   "tab",
			input:  "Hello Glab",
			indent: "\t",
			output: "\tHello Glab",
		},
		{
			name:   "prefix",
			input:  "Hello Glab",
			indent: "INFO: ",
			output: "INFO: Hello Glab",
		},
		{
			name:   "nothing",
			input:  "Hello Glab",
			indent: "",
			output: "Hello Glab",
		},
		{
			name:   "empty-string",
			input:  "",
			indent: "",
			output: "",
		},
		{
			name:   "multi-line",
			input:  "Hello\nGlab",
			indent: "- ",
			output: "- Hello\n- Glab",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := Indent(tC.input, tC.indent)
			assert.Equal(t, tC.output, got)
		})
	}
}

func Test_NewListTitle(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output ListTitleOptions
	}{
		{
			name:  "simple",
			input: "simple",
			output: ListTitleOptions{
				Name:           "simple",
				ListActionType: "list",
				Page:           1,
			},
		},
		{
			name:  "whitespace/leading",
			input: "   leading",
			output: ListTitleOptions{
				Name:           "leading",
				ListActionType: "list",
				Page:           1,
			},
		},
		{
			name:  "whitespace/trailing",
			input: "trailing    ",
			output: ListTitleOptions{
				Name:           "trailing",
				ListActionType: "list",
				Page:           1,
			},
		},
		{
			name:  "whitespace/leading-and-trailing",
			input: "   leading-and-trailing     ",
			output: ListTitleOptions{
				Name:           "leading-and-trailing",
				ListActionType: "list",
				Page:           1,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := NewListTitle(tC.input)
			assert.Equal(t, tC.output.Name, got.Name)
			assert.Equal(t, tC.output.ListActionType, got.ListActionType)
			assert.Equal(t, tC.output.Page, got.Page)
		})
	}
}

func Test_pluralizeName(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		amount int
		output string
	}{
		{
			name:   "singular",
			input:  "People",
			amount: 1,
			output: "People",
		},
		{
			name:   "plural",
			input:  "Human",
			amount: 3,
			output: "Humans",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			got := pluralizeName(tC.amount, tC.input)
			assert.Equal(t, tC.output, got)
		})
	}
}

func Test_Describe(t *testing.T) {
	opts := &ListTitleOptions{
		Name:             "test",
		Page:             0,
		CurrentPageTotal: 0,
		Total:            0,
		RepoName:         "glab",
		ListActionType:   "List",
		EmptyMessage:     "nothing here",
	}

	t.Run("empty-message/present", func(t *testing.T) {
		opts := *opts

		got := opts.Describe()
		assert.Equal(t, "nothing here", got)
	})
	t.Run("empty-message/absent", func(t *testing.T) {
		opts := *opts

		opts.EmptyMessage = ""

		got := opts.Describe()
		assert.Equal(t, "No tests available on glab", got)
	})

	t.Run("currentPageTotal/single-total", func(t *testing.T) {
		opts := *opts

		opts.Total = 0
		opts.CurrentPageTotal = 1
		opts.Page = 1

		got := opts.Describe()
		assert.Equal(t, "Showing 1 test on glab (Page 1)\n", got)
	})

	t.Run("currentPageTotal/single-page", func(t *testing.T) {
		opts := *opts

		opts.Total = 200
		opts.CurrentPageTotal = 1
		opts.Page = 1

		got := opts.Describe()
		assert.Equal(t, "Showing 1 of 200 tests on glab (Page 1)\n", got)
	})
	t.Run("currentPageTotal/multi-page", func(t *testing.T) {
		opts := *opts

		opts.Total = 200
		opts.CurrentPageTotal = 1
		opts.Page = 5

		got := opts.Describe()
		assert.Equal(t, "Showing 1 of 200 tests on glab (Page 5)\n", got)
	})

	t.Run("search/match", func(t *testing.T) {
		opts := *opts

		opts.ListActionType = "search"
		opts.CurrentPageTotal = 3
		opts.Page = 1

		got := opts.Describe()
		assert.Equal(t, "Showing 3 tests in glab that match your search (Page 1)\n", got)
	})
	t.Run("search/no-match", func(t *testing.T) {
		opts := *opts

		opts.ListActionType = "search"
		opts.CurrentPageTotal = 0

		got := opts.Describe()
		assert.Equal(t, "No tests match your search in glab\n", got)
	})
}
