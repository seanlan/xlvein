package xlredis

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	type args struct {
		uri      string
		username string
		password string
		prefix   string
		db       int
	}
	tests := []struct {
		name       string
		args       args
		wantClient Redis
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				uri:      "127.0.0.1:7001",
				username: "",
				password: "",
				prefix:   "",
				db:       0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClient, err := NewClient(tt.args.uri, tt.args.username, tt.args.password, tt.args.prefix, tt.args.db)
			if err != nil {
				t.Errorf("NewClient() error = %v", err)
				return
			}
			err = gotClient.Set(context.TODO(), "test", "test", 0).Err()
			if err != nil {
				t.Errorf("gotClient.Set() error = %v", err)
				return
			}
			r, err := gotClient.Get(context.TODO(), "test").Result()
			t.Log(r, err)
		})
	}
}
