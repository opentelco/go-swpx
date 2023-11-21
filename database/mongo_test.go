package database

import "testing"

func Test_hidePassword(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "to short",
			args: args{
				input: "",
			},
			want: "",
		},
		{
			name: "to short",
			args: args{
				input: "123",
			},
			want: "123",
		},
		{
			name: "to short",
			args: args{
				input: "xfyz",
			},
			want: "**yz",
		},
		{
			name: "test1",
			args: args{
				input: "abcd1234",
			},
			want: "****1234",
		},
		{
			name: "test1",
			args: args{
				input: "asdfsdafsdfsdfdsdfs",
			},
			want: "***************sdfs",
		},
		{
			name: "test1",
			args: args{
				input: "äöäöåäöäöasäsdfölasödlfkjasdäfölasädfölaäsdfölaäsdflö",
			},
			want: "*********************************************************************flö",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hidePassword(tt.args.input); got != tt.want {
				t.Errorf("hidePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
