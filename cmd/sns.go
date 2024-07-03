package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func Publish(input *sns.PublishBatchInput) (*sns.PublishBatchOutput, error) {
	mySession := session.Must(session.NewSession())
	client := sns.New(mySession)
	return client.PublishBatch(input)
}
