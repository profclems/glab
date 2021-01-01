package cmdutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseAssignees(t *testing.T) {
	testCases := []struct {
		name        string
		input       []string
		wantAdd     []string
		wantRemove  []string
		wantReplace []string
	}{
		{
			name:        "simple replace",
			input:       []string{"foo"},
			wantAdd:     []string{},
			wantRemove:  []string{},
			wantReplace: []string{"foo"},
		},
		{
			name:        "only add",
			input:       []string{"+foo"},
			wantAdd:     []string{"foo"},
			wantRemove:  []string{},
			wantReplace: []string{},
		},
		{
			name:        "only remove",
			input:       []string{"-foo", "!bar"},
			wantAdd:     []string{},
			wantRemove:  []string{"foo", "bar"},
			wantReplace: []string{},
		},
		{
			name:        "only replace",
			input:       []string{"baz"},
			wantAdd:     []string{},
			wantRemove:  []string{},
			wantReplace: []string{"baz"},
		},
		{
			name:        "add and remove",
			input:       []string{"+qux", "-foo", "!bar"},
			wantAdd:     []string{"qux"},
			wantRemove:  []string{"foo", "bar"},
			wantReplace: []string{},
		},
		{
			name:        "add and replace",
			input:       []string{"+foo", "bar"},
			wantAdd:     []string{"foo"},
			wantRemove:  []string{},
			wantReplace: []string{"bar"},
		},
		{
			name:        "remove and replace",
			input:       []string{"-foo", "bar", "!baz"},
			wantAdd:     []string{},
			wantRemove:  []string{"foo", "baz"},
			wantReplace: []string{"bar"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			gotAdd, gotRemove, gotReplace := ParseAssignees(tC.input)
			assert.ElementsMatch(t, gotAdd, tC.wantAdd)
			assert.ElementsMatch(t, gotRemove, tC.wantRemove)
			assert.ElementsMatch(t, gotReplace, tC.wantReplace)
		})
	}

}
