package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"syslog_queue", // Queue name
		true,           // Durable
		false,          // Delete when unused
		false,          // Exclusive
		false,          // No-wait
		nil,            // Arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}

	// Sample data to be indexed
	data := map[string]interface{}{
		"timestamp": time.Now(),
		"message":   "Hello, world!",
	}

	// Marshal data to JSON
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Error marshaling data to JSON:", err)
	}

	// Publish message to RabbitMQ
	err = ch.Publish(
		"",     // Exchange
		q.Name, // Routing key
		false,  // Mandatory
		false,  // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Fatal("Failed to publish a message:", err)
	}

	log.Println("Message published to RabbitMQ:", string(body))
}
