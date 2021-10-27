package text

import "testing"

func TestJoin(t *testing.T) {
	type args struct {
		list  []string
		delim string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "spaces as delim",
			args: args{
				list:  []string{"my", "name", "is"},
				delim: " ",
			},
			want: "my name is",
		},
		{
			name: "newline as delims",
			args: args{
				list:  []string{"my", "name", "is"},
				delim: "\n",
			},
			want: "my\nname\nis",
		},
		{
			name: "empty list",
			args: args{
				list:  []string{},
				delim: "\n",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Join(tt.args.list, tt.args.delim); got != tt.want {
				t.Errorf("Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrip(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{"\033[0;30mSome text"},
			want: "Some text",
		},
		{
			args: args{"\033]8;;https://example.com\033\\Example\033]8;;\033\\"},
			want: "Example",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Strip(tt.args.s); got != tt.want {
				t.Errorf("Strip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringWidth(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "withcolors",
			args: args{"\033[0;30mSome text"},
			want: 9,
		},
		{
			name: "without colors",
			args: args{"Some text"},
			want: 9,
		},
		{
			name: "hyperlink",
			args: args{"\033]8;;https://example.com\033\\Example\033]8;;\033\\"},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringWidth(tt.args.s); got != tt.want {
				t.Errorf("StringWidth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapString(t *testing.T) {
	type args struct {
		text      string
		lineWidth int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "long",
			args: args{
				text:      "this is a really long sentence",
				lineWidth: 10,
			},
			want: "this is a\nreally\nlong\nsentence",
		},
		{
			name: "short",
			args: args{
				text:      "shortword",
				lineWidth: 10,
			},
			want: "shortword",
		},
		{
			name: "empty",
			args: args{
				text:      "",
				lineWidth: 10,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WrapString(tt.args.text, tt.args.lineWidth); got != tt.want {
				t.Errorf("WrapString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	type args struct {
		str    string
		length int
		pad    byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				str:    "hello world",
				length: 20,
				pad:    ' ',
			},
			want: "hello world         ",
		},
		{
			args: args{
				str:    "a string longer than the specified length",
				length: 20,
				pad:    ' ',
			},
			want: "a string longer than the specified length",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadRight(tt.args.str, tt.args.length, tt.args.pad); got != tt.want {
				t.Errorf("PadRight() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	type args struct {
		str    string
		length int
		pad    byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				str:    "hello world",
				length: 20,
				pad:    ' ',
			},
			want: "         hello world",
		},
		{
			args: args{
				str:    "a string longer than the specified length",
				length: 20,
				pad:    ' ',
			},
			want: "a string longer than the specified length",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PadLeft(tt.args.str, tt.args.length, tt.args.pad); got != tt.want {
				t.Errorf("PadLeft() = %v, want %v", got, tt.want)
			}
		})
	}
}
