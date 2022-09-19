package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/arjunmalhotra1/bloXroute/common"
	"github.com/arjunmalhotra1/bloXroute/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var supportedCommands = []string{"add", "remove", "get", "get-all"}

func main() {
	fmt.Println("client...")
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// defer signal.Stop(sigChan)
	if err := godotenv.Load("./config.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
	awsSession := common.BuildSession()
	reader := bufio.NewReader(os.Stdin)

	for {
		inputText, _ := reader.ReadString('\n')
		// This is for because we don't want to send empty line to sqs.
		if inputText == "\n" {
			continue
		}
		input := strings.Split(inputText, "\n")
		command := strings.Split(input[0], " ")

		if len(command) == 0 ||
			len(command) > 2 {
			log.Println("Please see the readme on how to run the commands.")
			continue
		}

		if !isValid(command[0]) || (command[0] == "remove" && len(command) != 2) ||
			(command[0] == "add" && len(command) != 2) ||
			(command[0] == "get" && len(command) != 2) ||
			(command[0] == "get-all" && len(command) != 1) {
			log.Println("Please see the readme on how to run the commands.")
			continue
		}
		var mssg models.Message
		if command[0] == "get-all" {
			mssg = models.Message{
				Function: command[0],
				Value:    "",
			}
		} else {
			mssg = models.Message{
				Function: command[0],
				Value:    command[1],
			}
		}
		mssgBytes, _ := json.Marshal(mssg)

		sendSQS(awsSession, os.Getenv("SQS_URL"), string(mssgBytes))

	}

}

func isValid(command string) bool {
	for _, v := range supportedCommands {
		if command == v {
			return true
		}
	}
	return false

}
func sendSQS(session *session.Session, destination string, message string) {
	fmt.Println("sending message ", message)
	svc := sqs.New(session, nil)
	messageGId := uuid.New().String()

	sendInput := &sqs.SendMessageInput{
		MessageBody:    aws.String(message),
		QueueUrl:       aws.String(destination),
		MessageGroupId: &messageGId,
	}

	_, err := svc.SendMessage(sendInput)
	if err != nil {
		fmt.Println(err)
		return
	}

}
