package expand

import (
	"errors"
	"reflect"
	"testing"

	"github.com/profclems/glab/internal/config"

	"github.com/MakeNowJust/heredoc"
)

func TestExpandAlias(t *testing.T) {
	findShFunc := func() (string, error) {
		return "/usr/bin/sh", nil
	}

	cfg := config.NewFromString(heredoc.Doc(`
		aliases:
		  co: mr checkout
		  il: issue list --author="$1" --label="$2"
		  ia: issue list --author="$1" --assignee="$1"
	`))

	type args struct {
		config config.Config
		argv   []string
	}
	tests := []struct {
		name         string
		args         args
		wantExpanded []string
		wantIsShell  bool
		wantErr      error
	}{
		{
			name: "no arguments",
			args: args{
				config: cfg,
				argv:   []string{},
			},
			wantExpanded: []string(nil),
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "too few arguments",
			args: args{
				config: cfg,
				argv:   []string{"glab"},
			},
			wantExpanded: []string(nil),
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "no expansion",
			args: args{
				config: cfg,
				argv:   []string{"glab", "mr", "status"},
			},
			wantExpanded: []string{"mr", "status"},
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "simple expansion",
			args: args{
				config: cfg,
				argv:   []string{"glab", "co"},
			},
			wantExpanded: []string{"mr", "checkout"},
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "adding arguments after expansion",
			args: args{
				config: cfg,
				argv:   []string{"glab", "co", "123"},
			},
			wantExpanded: []string{"mr", "checkout", "123"},
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "not enough arguments for expansion",
			args: args{
				config: cfg,
				argv:   []string{"glab", "il"},
			},
			wantExpanded: []string{},
			wantIsShell:  false,
			wantErr:      errors.New(`not enough arguments for alias: issue list --author="$1" --label="$2"`),
		},
		{
			name: "not enough arguments for expansion 2",
			args: args{
				config: cfg,
				argv:   []string{"glab", "il", "vilmibm"},
			},
			wantExpanded: []string{},
			wantIsShell:  false,
			wantErr:      errors.New(`not enough arguments for alias: issue list --author="vilmibm" --label="$2"`),
		},
		{
			name: "satisfy expansion arguments",
			args: args{
				config: cfg,
				argv:   []string{"glab", "il", "vilmibm", "help wanted"},
			},
			wantExpanded: []string{"issue", "list", "--author=vilmibm", "--label=help wanted"},
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "mixed positional and non-positional arguments",
			args: args{
				config: cfg,
				argv:   []string{"glab", "il", "vilmibm", "epic", "-R", "monalisa/testing"},
			},
			wantExpanded: []string{"issue", "list", "--author=vilmibm", "--label=epic", "-R", "monalisa/testing"},
			wantIsShell:  false,
			wantErr:      nil,
		},
		{
			name: "dollar in expansion",
			args: args{
				config: cfg,
				argv:   []string{"glab", "ia", "$coolmoney$"},
			},
			wantExpanded: []string{"issue", "list", "--author=$coolmoney$", "--assignee=$coolmoney$"},
			wantIsShell:  false,
			wantErr:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExpanded, gotIsShell, err := ExpandAlias(tt.args.config, tt.args.argv, findShFunc)
			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("expected error")
				}
				if tt.wantErr.Error() != err.Error() {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(gotExpanded, tt.wantExpanded) {
				t.Errorf("ExpandAlias() gotExpanded = %v, want %v", gotExpanded, tt.wantExpanded)
			}
			if gotIsShell != tt.wantIsShell {
				t.Errorf("ExpandAlias() gotIsShell = %v, want %v", gotIsShell, tt.wantIsShell)
			}
		})
	}
}
