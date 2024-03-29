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
                log.Printf("Failed to send message to queue %s %v", queueName, err)
                return nil
        }

        return result
}

func getMessages(queueName string) ([]*sqs.Message) {
        qURL := getQueueUrl(queueName)

        result, err := messageService.ReceiveMessage(&sqs.ReceiveMessageInput{
                        QueueUrl: qURL,
                        AttributeNames: aws.StringSlice([]string{
                                "SentTimestamp",
                        }),
                        MaxNumberOfMessages: aws.Int64(10),
                        MessageAttributeNames: aws.StringSlice([]string{
                        "All",
                }),
                WaitTimeSeconds: aws.Int64(20),
        })

        if err != nil {
                log.Printf("Unable to get messages from queue %q, %v.", queueName, err)
                return nil
        } else if result == nil {
                log.Printf("Found no messages on queue (%q) during this poll", queueName)
                return nil
        } else {
                return result.Messages
        }
}

func deleteMessage(message *sqs.Message, queueName string) (*sqs.DeleteMessageOutput) {
        qURL := getQueueUrl(queueName)

        resultDelete, err := messageService.DeleteMessage(&sqs.DeleteMessageInput{
                QueueUrl:      qURL,
                ReceiptHandle: message.ReceiptHandle,
        })

        if err != nil {
                log.Printf("Failed to delete message from queue %q, %v.", queueName, err)
                return nil
        } else {
                return resultDelete
        }
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
