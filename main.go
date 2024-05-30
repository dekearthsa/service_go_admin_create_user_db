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

type Request struct {
	Detail Detail `json:"detail"`
}

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
	UserID     string   `json:"UserID"`
	CreateDate int      `json:"CreateDate"`
	Email      string   `json:"Email"`
	Password   string   `json:"Password"`
	FristName  string   `json:"FristName"`
	LastName   string   `json:"LastName"`
	PlantName  string   `json:"PlantName"`
	LineUserId string   `json:"LineUserId"`
	UserTenan  string   `json:"UserTenan"`
	UserType   string   `json:"UserType"`
	Tel        string   `json:"Tel"`
	IsProduct  []string `json:"IsProduct"`
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "cannot hash password", err
	}

	return string(hashedPassword), nil
}

func handler(ctx context.Context, req Request) (string, error) {
	var TABLENAME = "demo_user_line_id"
	var REGION = "ap-southeast-1"
	// var setPayload Payload
	now := time.Now()
	ms := now.UnixNano() / int64(time.Millisecond)
	// fmt.Println("req => ", req)
	// fmt.Println("req Data => ", req.Detail.Data)
	// fmt.Println("req Data Email => ", req.Detail.Data.Email)
	hash, err := HashPassword(req.Detail.Data.Password)
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
		Email:      req.Detail.Data.Email,
		UserID:     req.Detail.Data.UserID,
		CreateDate: int(ms),
		Password:   hash,
		FristName:  req.Detail.Data.FristName,
		LastName:   req.Detail.Data.LastName,
		PlantName:  req.Detail.Data.PlantName,
		LineUserId: req.Detail.Data.LineUserId,
		UserTenan:  req.Detail.Data.UserTenan,
		UserType:   req.Detail.Data.UserType,
		Tel:        req.Detail.Data.Tel,
		IsProduct:  req.Detail.Data.IsProduct,
	}
	// fmt.Println("setPayload => ", setPayload)

	masrhalData, err := dynamodbattribute.MarshalMap(setPayload)
	if err != nil {
		fmt.Println("masrhalData err")
		log.Fatalf("Got error marshalling new movie item: %s", err)
	}
	fmt.Println("masrhalData => ", masrhalData)

	input := &dynamodb.PutItemInput{
		Item:      masrhalData,
		TableName: aws.String(TABLENAME),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("PutItem err")
		log.Fatalf("Got error calling PutItem: %s", err)
	}

	fmt.Println("logging inserted!")

	return "Inserted", nil
}

func main() {
	lambda.Start(handler)
}
