package main

import (
	"bufio"
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

// todo: common func
func OnError(err error, msg string) {
	if err != nil {
		fmt.Errorf("%s: %s", msg, err)
	}
}


func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	OnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	OnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bloXroute",
		false,
		false,
		false,
		false,
		nil,   
	)
	OnError(err, "Failed to declare a queue")

		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("-> ")
			command, _ := reader.ReadString('\n')
			body := command
			 err := ch.Publish(
				"",
				q.Name,
				false,
				false,
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				})
			OnError(err, "Failed to publish a message")
		}
		fmt.Println("Successfully Published Message to Queue")
	}