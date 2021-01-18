package backend

import (
	"anagrams/internal/database"
	"anagrams/internal/logging"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const StartIndex = 0
const StopIndex = -1

var ErrAnagramsNotFound = errors.New("anagrams not found")
var Logger = logging.GetLogger("backend")

func GetAnagrams(word string) ([]string, error) {
	Logger.Info(fmt.Sprintf("Starting to find anagrams for word '%s'", word))

	sortedWord := SortString(strings.ToLower(word))
	Logger.Debug(fmt.Sprintf("Sorted word is '%s'", sortedWord))
	anagrams, err := database.RedisClient.LRange(database.Ctx, sortedWord, StartIndex, StopIndex).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get anagrams from redis: %w", err)
	}

	if len(anagrams) == 0 {
		Logger.Info("Anagrams not found for the input word")
		return nil, ErrAnagramsNotFound
	}

	Logger.Info(fmt.Sprintf("Found anagrams: %v", anagrams))
	return anagrams, nil
}

func AddWords(words []string) error {
	Logger.Info(fmt.Sprintf("Adding %d new words", len(words)))

	for _, word := range words {
		keyword := SortString(strings.ToLower(word))
		err := database.RedisClient.LPush(database.Ctx, keyword, word).Err()
		if err != nil {
			return fmt.Errorf("failed to add new words: %w", err)
		}
	}

	Logger.Info("New words added successfully")
	return nil
}

func LoadNewWords(words []string) error {
	Logger.Info("Starting to load new wordlist")

	Logger.Debug("Run FlushDB for clear old data")
	if err := database.RedisClient.FlushDB(database.Ctx).Err(); err != nil {
		return fmt.Errorf("failed to clear old data from redis: %w", err)
	}
	Logger.Debug("Old data cleared successfully")

	return AddWords(words)
}
