package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// getSymbols fetch ticker symbols
func newUser(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	fmt.Printf("Postconfirmation for user: %s\n", event.UserName)
	return event, nil
}

func main() {
	lambda.Start(newUser)
}
