package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Detail struct {
	Data    Payload `json:"data"`
	BusName string  `json:"busName"`
}

type Payload struct {
	UserID     string   `json:"userID"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	FristName  string   `json:"fristName"`
	LastName   string   `json:"lastName"`
	PlantName  string   `json:"plantName"`
	LineUserId string   `json:"lineUserId"`
	UserTenan  string   `json:"userTenan"`
	UserType   string   `json:"userType"`
	Tel        string   `json:"tel"`
	IsProduct  []string `json:"isProduct"`
}

func handler(ctx context.Context, req Detail) (string, error) {
	var TABLENAME = "demo_user_line_id"
	var REGION = "ap-southeast-1"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.New(sess)

	masrhalData, err := dynamodbattribute.MarshalMap(req)
	if err != nil {
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}

	input := &dynamodb.PutItemInput{
		Item:      masrhalData,
		TableName: aws.String(TABLENAME),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("logging inserted!")

	return "Inserted", nil
}

func main() {
	lambda.Start(handler)
}
