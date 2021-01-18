package database

import (
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockRedis() (*miniredis.Miniredis, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return mr, err
	}

	return mr, nil
}

func TestInitRedis(t *testing.T) {
	mr, err := mockRedis()
	if err != nil {
		t.Errorf("Mock Redis error: %v", err)
	}
	defer mr.Close()

	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "127.0.0.1",
			want: "dial tcp: address 127.0.0.1: missing port in address",
		},
	}
	for _, tt := range tests {
		got := InitRedis(tt.in)
		require.EqualError(t, got, tt.want)
	}
}
