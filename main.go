package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"golang.org/x/crypto/bcrypt"
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

type PayloadDB struct {
	UserID     string   `json:"userID"`
	CreateDate int      `json:"createDate"`
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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "cannot hash password", err
	}

	return string(hashedPassword), nil
}

func handler(ctx context.Context, req Detail) (string, error) {
	var TABLENAME = "demo_user_line_id"
	var REGION = "ap-southeast-1"
	// var setPayload Payload
	now := time.Now()
	ms := now.UnixNano() / int64(time.Millisecond)
	fmt.Println("req => ", req.Data)
	hash, err := HashPassword(req.Data.Password)
	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(REGION)},
	)
	if err != nil {
		log.Fatal(err)
	}

	svc := dynamodb.New(sess)

	setPayload := PayloadDB{
		UserID:     req.Data.UserID,
		CreateDate: int(ms),
		Email:      req.Data.Email,
		Password:   hash,
		FristName:  req.Data.FristName,
		LastName:   req.Data.LastName,
		PlantName:  req.Data.PlantName,
		LineUserId: req.Data.LineUserId,
		UserTenan:  req.Data.UserTenan,
		UserType:   req.Data.UserType,
		Tel:        req.Data.Tel,
		IsProduct:  req.Data.IsProduct,
	}

	masrhalData, err := dynamodbattribute.MarshalMap(setPayload)
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
