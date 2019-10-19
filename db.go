package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"os"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

func getItems() ([]rate, error) {
	// prepare input for the query
	input := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	result, err := db.Scan(input)
	if err != nil {
		return nil, err
	}
	if *result.Count == 0 {
		return []rate{}, nil
	}

	rates := []rate{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &rates)
	if err != nil {
		return nil, err
	}
	return rates, nil
}
