package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Item struct {
	Action string `json:"action,omitempty"`
	Id int `json:"id,omitempty"`
	TextMessage string `json:"textMessage,omitempty"`
}

var curMsgIndex int
var items =make(map[int]string)


func addItem(textMessage string) {
	curMsgIndex++
	items[curMsgIndex]=textMessage
	log.Println("----addItem----, textMessage:",textMessage)
}

func removeItem(num int) {
	delete(items, num)
	log.Println("----removeItem----, id:",num)
}

func getItem(num int)  {
	if val, ok := items[num]; ok {
		log.Println("----getItem----, id:",num, val)
	}else {
		log.Println("----getItem----, id:",num, "Does not exist")
	}
}

func getAllItems() {
	log.Println("----getAllItems-----")
	for id, element := range items {
		log.Println("itemId:", id, "=>", "item:", element)
	}
}
// todo: common func
func onError(err error, msg string) {
	if err != nil {
		fmt.Errorf("%s: %s", msg, err)
	}
}

func main() {
	// todo: log file + date
	f, err := os.OpenFile("logFile.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		onError(err,"cannot open a log file file")
	}
	defer f.Close()
	log.SetOutput(f)

	// todo: security - login user name & password in configuration file
	// todo: common func in one file - connection, open channel and declare rabbit

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	onError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	onError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bloXroute",
		false,
		false,
		false,
		false,
		nil,
	)
	onError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	onError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// todo: validations
			log.Printf("Received a message: %s", d.Body)
			var item Item
			if err1 := json.Unmarshal(d.Body, &item); err1 != nil {
				onError(err1, "cannot Unmarshal body to message")
			} else {
				switch item.Action {
				case "AddItem":
					addItem(item.TextMessage)
				case "RemoveItem":
					removeItem(item.Id)
				case "GetItem":
					getItem(item.Id)
				case "GetAllItems":
					getAllItems()
				default:
					getAllItems()
				}
			}
		}
	}()
	log.Println(" [*] - Waiting for messages")
	<-forever
}