package backend

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"testing"

	"anagrams/internal/database"
)

func mockRedis() (*miniredis.Miniredis, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return mr, err
	}
	database.RedisClient = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	database.Ctx = context.Background()
	return mr, nil
}

func TestGetAnagrams(t *testing.T) {
	mr, err := mockRedis()
	if err != nil {
		t.Errorf("Mock Redis error: %v", err)
	}
	defer mr.Close()

	words := []string{"foobar", "aabb", "baba", "boofar", "test"}
	if err = LoadNewWords(words); err != nil {
		t.Errorf("Function 'LoadNewWords' return error: %v", err)
	}

	tests := []struct {
		in   string
		want []string
	}{
		{
			in:   "abba",
			want: []string{"baba", "aabb"},
		},
		{
			in:   "abfoor",
			want: []string{"boofar", "foobar"},
		},
		{
			in:   "raboof",
			want: []string{"boofar", "foobar"},
		},
		{
			in:   "incorrect",
			want: nil,
		},
	}
	for _, tt := range tests {
		got, err := GetAnagrams(tt.in)
		if err != nil && err != ErrAnagramsNotFound {
			t.Errorf("Function 'GetAnagrams' return error: %v", err)
		}

		require.Equal(t, tt.want, got)
	}
}

func TestLoadNewWords(t *testing.T) {
	mr, err := mockRedis()
	if err != nil {
		t.Errorf("Mock Redis error: %v", err)
	}
	defer mr.Close()

	tests := []struct {
		in   []string
		want []string
	}{
		{
			in:   []string{"foobar", "aabb", "baba", "boofar", "test"},
			want: []string{"aabb", "abfoor", "estt"},
		},
		{
			in:   []string{"aabb", ""},
			want: []string{"", "aabb"},
		},
		{
			in:   nil,
			want: []string{},
		},
	}

	for _, tt := range tests {
		if err = LoadNewWords(tt.in); err != nil {
			t.Errorf("Function 'LoadNewWords' return error: %v", err)
		}
		got, err := database.RedisClient.Keys(database.Ctx, "*").Result()
		if err != nil {
			t.Errorf("RedisClient.Keys return error: %v", err)
		}

		require.Equal(t, tt.want, got)
	}
}

func TestAddNewWords(t *testing.T) {
	mr, err := mockRedis()
	if err != nil {
		t.Errorf("Mock Redis error: %v", err)
	}
	defer mr.Close()

	words := []string{"foobar", "boofar"}
	if err = LoadNewWords(words); err != nil {
		t.Errorf("Function 'LoadNewWords' return error: %v", err)
	}

	newWords := []string{"bofaro", "boofra"}
	if err = AddWords(newWords); err != nil {
		t.Errorf("Function 'AddWords' return error: %v", err)
	}
	want := []string{"boofra", "bofaro", "boofar", "foobar"}
	got, err := GetAnagrams("foobar")
	if err != nil {
		t.Errorf("Function 'GetAnagrams' return error: %v", err)
	}

	require.Equal(t, want, got)
}
