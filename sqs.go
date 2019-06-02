package main

import (
	"log"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/sqs"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
                SharedConfigState: session.SharedConfigEnable,
        }))
var messageService = sqs.New(sess)

func sendMessage(message string, queueName string) (*sqs.SendMessageOutput) {

	qURL := getQueueUrl(queueName)

        result, err := messageService.SendMessage(&sqs.SendMessageInput{
                MessageBody: aws.String(message),
                QueueUrl:    qURL,
        })

        if err != nil {
                log.Printf("Error", err)
                return nil
        }

        return result
}


func getQueueUrl(queueName string) *string {
	qURLOutput, err := messageService.GetQueueUrl(&sqs.GetQueueUrlInput{
                QueueName: aws.String(queueName),
        })

	var qURL = qURLOutput.QueueUrl

        if err != nil {
                if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
                        createQueueOutput, err :=  messageService.CreateQueue(&sqs.CreateQueueInput{
                                QueueName: aws.String(queueName),
                        })
                        if err != nil {
                                log.Printf("Unable to create queue %q, %v.", queueName, err)
                        } else {
                                log.Printf("Create queue %s", queueName)
                                qURL = createQueueOutput.QueueUrl
                        }
                }
                log.Printf("Unable to queue %q, %v.", queueName, err)
        }

	return qURL
}
