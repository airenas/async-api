package auth

import (
	"reflect"
	"testing"
)

func Test_loginAuth_Next(t *testing.T) {
	type fields struct {
		username string
		password string
	}
	type args struct {
		fromServer []byte
		more       bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{name: "Pass", fields: fields{username: "aaa", password: "passs"}, args: args{fromServer: []byte("Password:"), more: true},
			want: []byte("passs"), wantErr: false},
		{name: "User", fields: fields{username: "aaa", password: "passs"}, args: args{fromServer: []byte("Username:"), more: true},
			want: []byte("aaa"), wantErr: false},
		{name: "Error", fields: fields{username: "aaa", password: "passs"}, args: args{fromServer: []byte("Any:"), more: true},
			want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &loginAuth{
				username: tt.fields.username,
				password: tt.fields.password,
			}
			got, err := a.Next(tt.args.fromServer, tt.args.more)
			if (err != nil) != tt.wantErr {
				t.Errorf("loginAuth.Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loginAuth.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}
