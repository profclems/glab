package text

import "testing"

func TestTruncate(t *testing.T) {
	type args struct {
		s      string
		length int
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "short",
			args: args{
				s:      "shortword",
				length: 9,
			},
			want: "shortword",
		},

		{
			name: "long sentence",
			args: args{
				s:      "this is a really long sentence",
				length: 10,
			},
			want: "this is...",
		},

		{
			name: "hyperlink prefix",
			args: args{
				s:      "\033]8;;https://example.com\033\\this\033]8;;\033\\ is a really long sentence",
				length: 10,
			},
			want: "\033]8;;https://example.com\033\\this\033]8;;\033\\ is...",
		},

		{
			name: "hyperlink prefix but shorter",
			args: args{
				s:      "\033]8;;https://example.com\033\\this\033]8;;\033\\ is a really long sentence",
				length: 4,
			},
			want: "\033]8;;https://example.com\033\\t...\033]8;;\033\\",
		},

		{
			name: "hyperlink infix",
			args: args{
				s:      "this \033]8;;https://example.com\033\\is\033]8;;\033\\ a really long sentence",
				length: 10,
			},
			want: "this \033]8;;https://example.com\033\\is\033]8;;\033\\...",
		},

		{
			name: "hyperlink infix straddling ellipsis",
			args: args{
				s:      "this \033]8;;https://example.com\033\\is\033]8;;\033\\ a really long sentence",
				length: 9,
			},
			want: "this \033]8;;https://example.com\033\\i.\033]8;;\033\\..",
		},

		{
			name: "hyperlink suffix",
			args: args{
				s:      "this is a really long \033]8;;https://example.com\033\\sentence\033]8;;\033\\",
				length: 10,
			},
			want: "this is...",
		},
		{
			name: "hyperlink suffix take 2",
			args: args{
				s:      "this \033]8;;https://example.com\033\\is a really long sentence\033]8;;\033\\",
				length: 10,
			},
			want: "this \033]8;;https://example.com\033\\is...\033]8;;\033\\",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("Truncate() = %q, want %q", got, tt.want)
			}
		})
	}
}
