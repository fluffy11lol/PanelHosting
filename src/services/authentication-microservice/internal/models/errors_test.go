package models

import (
	"testing"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "invalid password error",
			err:  ErrInvalidPassword,
			want: "invalid password",
		},
		{
			name: "user exists error",
			err:  ErrUserExist,
			want: "user already exist, try other name",
		},
		{
			name: "user not exists error",
			err:  ErrUserNotExist,
			want: "no user with this name",
		},
		{
			name: "invalid expression error",
			err:  ErrInvalidExpression,
			want: "invalid format of expression",
		},
		{
			name: "empty field error",
			err:  ErrEmptyField,
			want: "empty fields not allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error message = %v, want %v", got, tt.want)
			}
		})
	}
}
