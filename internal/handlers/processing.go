package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"anagrams/internal/backend"
)

func Get(c *gin.Context) {
	word := c.Query("word")
	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"result":  nil,
			"error":   "couldn't found 'word' parameter",
		})
		return
	}

	anagrams, err := backend.GetAnagrams(word)
	if err == backend.ErrAnagramsNotFound {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"result":  nil,
			"error":   nil,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"result":  nil,
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  anagrams,
		"error":   nil,
	})
}

func Load(c *gin.Context) {
	var words []string

	if err := c.BindJSON(&words); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"result":  nil,
			"error":   fmt.Sprintf("failed to decode input json: %v", err),
		})
		return
	}
	if err := backend.LoadNewWords(words); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"result":  nil,
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  nil,
		"error":   nil,
	})
}

func Add(c *gin.Context) {
	var words []string

	if err := c.BindJSON(&words); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"result":  nil,
			"error":   fmt.Sprintf("failed to decode input json: %v", err),
		})
		return
	}

	if err := backend.AddWords(words); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"result":  nil,
			"error":   fmt.Sprintf("%v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"result":  nil,
		"error":   nil,
	})
}
