package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func callOpenAICompletions(c *gin.Context) {
	// Authenticate with Flock
	flockApiKey := c.GetHeader("FLOCK-AUTH")
	if flockApiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "FLOCK-AUTH header is missing"})
		return
	}
	fmt.Println("Received FlOCK-AUTH:", flockApiKey) // Print FLOCK-AUTH value to console

	// Read raw request
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Forward request to OpenAI
	openaiApiKey := c.GetHeader("Authorization")
	if openaiApiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	url := "https://api.openai.com/v1/completions"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", openaiApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var openAIResp map[string]interface{}
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, openAIResp)
}

func main() {
	router := gin.Default()
	router.POST("/openai/v1/completions", callOpenAICompletions)

	router.Run("localhost:8080")
}
