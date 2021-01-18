package backend

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSortString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "foobar",
			want: "abfoor",
		},
		{
			in:   "aaabbbc",
			want: "aaabbbc",
		},
		{
			in:   "54321",
			want: "12345",
		},
		{
			in:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		got := SortString(tt.in)
		require.Equal(t, tt.want, got)
	}
}
