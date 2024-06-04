package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func publishToRabbitMQ(res []byte) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"flock-processor", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Body: []byte(res),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", res)
}

func callOpenAICompletions(c *gin.Context) {
	// Authenticate with Flock
	flockApiKey := c.GetHeader("FLOCK-AUTH")
	if flockApiKey == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "FLOCK-AUTH header is missing"})
		return
	}
	fmt.Println("Received FLOCK-AUTH:", flockApiKey) // Print FLOCK-AUTH value to console

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

	startTime := time.Now()
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	endTime := time.Now()
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

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res["flock_metrics"] = map[string]interface{}{
		"response_time": endTime.Sub(startTime).Microseconds(),
	}

	data, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Error serializing map to JSON: %v", err)
	}

	publishToRabbitMQ(data)

	c.Data(http.StatusOK, "application/json", body)
}

func main() {
	router := gin.Default()
	router.POST("/openai/v1/completions", callOpenAICompletions)

	router.Run("localhost:8080")
}
