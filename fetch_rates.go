package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"net/http"
	"os"
)

var errorLogger = log.New(os.Stderr, "Error ", log.Llongfile)

type rate struct {
	ID   string  `json:"id"`   // a UUID
	Bank string  `json:"bank"` // bank id
	Day  string  `json:"day"`  // an ISO representation of the date YYYY-MM-DD
	Buy  float32 `json:"buy"`  // the value of the dollar in HTG when the bank buys
	Sell float32 `json:"sell"` // the value of the dollar in HTG when the bank sells
}

func list(event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rates, err := getItems()
	if err != nil {
		return serverError(err)
	}

	js, err := json.Marshal(rates)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(list)
}
