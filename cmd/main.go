package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/sns"
)

type CheckResponse struct {
	URL       string
	Available bool
	Price     float64
}

var SNS_TOPIC_ARN = os.Getenv("AWS_SNS_TOPIC_ARN")

type Event struct {
	OzoneItems []string `json:"ozoneItems"`
}

func HandleRequest(ctx *context.Context, event *Event) {
	oc := NewOzoneChecker()
	resChan := make(chan CheckResponse, 10)

	var wg sync.WaitGroup
	wg.Add(len(event.OzoneItems))

	go func() {
		defer close(resChan)

		for _, url := range event.OzoneItems {
			go func(url string) {
				defer wg.Done()
				resp, err := oc.Check(url)
				if err != nil {
					log.Fatal(err)
				}
				resChan <- resp
			}(url)
		}
		wg.Wait()
	}()

	publish := sns.PublishBatchInput{TopicArn: &SNS_TOPIC_ARN,
		PublishBatchRequestEntries: make([]*sns.PublishBatchRequestEntry, 0)}

	i := 0
	mailSubject := "Product available!"
	for resp := range resChan {
		if !resp.Available {
			continue
		}

		id := strconv.Itoa(i)
		i++

		message := fmt.Sprintf("Hello,\n\nThe %s is now available at price %.2flv!\n", resp.URL, resp.Price)
		entry := sns.PublishBatchRequestEntry{Subject: &mailSubject, Message: &message, Id: &id}
		publish.PublishBatchRequestEntries = append(publish.PublishBatchRequestEntries, &entry)
	}

	if len(publish.PublishBatchRequestEntries) == 0 {
		log.Println("no available items found.")
		return
	}

	_, err := Publish(&publish)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
