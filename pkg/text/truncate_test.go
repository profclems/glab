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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Truncate(tt.args.s, tt.args.length); got != tt.want {
				t.Errorf("Truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}
