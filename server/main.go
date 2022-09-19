package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/arjunmalhotra1/bloXroute/common"
	"github.com/arjunmalhotra1/bloXroute/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/joho/godotenv"
)

var wg sync.WaitGroup
var orderedMap sync.Map
var file *os.File
var fileMutex sync.Mutex

func main() {
	fmt.Println("Server ...")
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// defer signal.Stop(sigChan)
	var err error
	if err = godotenv.Load("./config.env"); err != nil {
		log.Fatalf("Error loading .env file %+v", err)
	}

	file, err = os.OpenFile(
		"./output",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666)
	if err != nil {
		log.Fatalf("Couldn't open the file %v", err)
	}
	awsSession := common.BuildSession()
	svc := sqs.New(awsSession, nil)
	maxthreads := runtime.NumCPU()

	wg.Add(maxthreads)
	for i := 0; i < maxthreads; i++ {
		go func() {
			defer wg.Done()

			for {
				messages := receiveMessages(svc, os.Getenv("SQS_URL"))

				if len(messages) == 0 {
					pollTime, err := strconv.ParseInt(os.Getenv("POLLING_TIME"), 10, 64)
					if err != nil {
						fmt.Printf("Error reading polling time from env file %v", err)
					}
					//log.Println(pollTime)
					time.Sleep(time.Duration(time.Duration(pollTime) * time.Microsecond))
				}

				for _, msg := range messages {
					if msg == nil {
						continue
					}
					processMessage(msg.Body)
					deleteMessage(svc, msg.ReceiptHandle, os.Getenv("SQS_URL"))
				}
			}
		}()
	}
	wg.Wait()
}

func deleteMessage(svc *sqs.SQS, handle *string, url string) {
	delInput := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(url),
		ReceiptHandle: handle,
	}
	_, err := svc.DeleteMessage(delInput)
	if err != nil {
		log.Println("Error while deleting message from sqs")
	}
}

func processMessage(message *string) {
	var msg models.Message
	err := json.Unmarshal([]byte(*message), &msg)
	if err != nil {
		log.Println("Error while un-marshalling the message")
	}
	switch msg.Function {
	case "add":
		orderedMap.Store(msg.Value, true)
		output := fmt.Sprintf("Added value %s \n", msg.Value)
		fileMutex.Lock()
		file.Write([]byte(output)) // TODO: Handle error
		fileMutex.Unlock()
		log.Println(output)
	case "remove":
		orderedMap.Delete(msg.Value)
		output := fmt.Sprintf("Deleted value %s \n", msg.Value)
		fileMutex.Lock()
		file.Write([]byte(output)) // TODO: Handle error
		fileMutex.Unlock()
		log.Println(output)
	case "get":
		var output string
		_, ok := orderedMap.Load(msg.Value)
		if ok {
			output = fmt.Sprintf("Value %s is present \n", msg.Value)
		} else {
			output = fmt.Sprintf("Value %s is NOT present \n", msg.Value)
		}
		fileMutex.Lock()
		file.Write([]byte(output)) // TODO: Handle error
		fileMutex.Unlock()
		log.Println(output)
	case "get-all":
		var output string
		orderedMap.Range(
			func(key, value interface{}) bool {
				output = output + fmt.Sprintf("%v ", key)
				return true
			})
		output = "Get all : " + output
		fileMutex.Lock()
		file.Write([]byte(output)) // TODO: Handle error
		fileMutex.Unlock()
		log.Println(output)
	}
}

func receiveMessages(svc *sqs.SQS, queueUrl string) []*sqs.Message {
	receiveMessagesInput := &sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: aws.Int64(10),
		WaitTimeSeconds:     aws.Int64(20),
		VisibilityTimeout:   aws.Int64(20),
	}

	receiveMessagesOuput, err := svc.ReceiveMessage(receiveMessagesInput)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}

	if receiveMessagesOuput == nil {
		return nil
	}

	return receiveMessagesOuput.Messages
}
