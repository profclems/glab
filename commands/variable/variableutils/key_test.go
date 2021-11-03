package variableutils

import "testing"

func Test_isValidKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "key is empty",
			args: args{
				key: "",
			},
			want: false,
		},
		{
			name: "key is valid",
			args: args{
				key: "abc123_",
			},
			want: true,
		},
		{
			name: "key is invalid",
			args: args{
				key: "abc-123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidKey(tt.args.key); got != tt.want {
				t.Errorf("isValidKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
