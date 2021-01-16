package api

import "testing"

func TestIsValidToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Token Is Valid",
			args: args{token: "xxxxxxxxxxxxxxxxxxxx"},
			want: true,
		},
		{
			name: "Token Is inValid",
			args: args{token: "123"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidToken(tt.args.token); got != tt.want {
				t.Errorf("IsValidToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
