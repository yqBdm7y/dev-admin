package dadmin

import "testing"

func TestUser_ValidatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		u    User
		args args
		want bool
	}{
		{name: "Validating password: 12345678", u: User{}, args: args{password: "12345678"}, want: false},
		{name: "Validating password: a1234567", u: User{}, args: args{password: "a1234567"}, want: true},
		{name: "Validating password: a12345", u: User{}, args: args{password: "a12345"}, want: false},
		{name: "Validating password: .1545.-150", u: User{}, args: args{password: ".1545.-150"}, want: false},
		{name: "Validating password: %1545!%150", u: User{}, args: args{password: "%1545!%150"}, want: true},
		{name: "Validating password: ........", u: User{}, args: args{password: "........"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ValidatePassword(tt.args.password); got != tt.want {
				t.Errorf("User.ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
