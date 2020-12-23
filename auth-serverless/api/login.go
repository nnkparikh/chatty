package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/chatty/auth-serverless/api/idp"
)

// AuthRequestBody : request parameters for initAuth coginto API method
type AuthRequestBody struct {
	AuthFlow       string
	AuthParameters map[string]string
	ClientId       string
}

func login(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	formData, err := url.ParseQuery(request.Body)
	svcEndpoint := url.URL{
		Scheme: "https",
		Host:   idp.CognitoSvcEndpoint,
	}
	log.Printf("Form data: %s\n", formData)
	log.Printf("svcEndpoint: %s\n", svcEndpoint.String())
	log.Printf("username: %s\n", formData["username"][0])
	log.Printf("password: %s\n", formData["password"][0])
	log.Printf("client id: %s\n", os.Getenv("CLIENT_ID"))
	authReqBody := AuthRequestBody{
		AuthFlow: "USER_PASSWORD_AUTH",
		AuthParameters: map[string]string{
			"PASSWORD": formData["password"][0],
			"USERNAME": formData["username"][0],
		},
		ClientId: os.Getenv("CLIENT_ID"),
	}

	authReqBodyJSON, err := json.Marshal(authReqBody)
	log.Printf("req body: %s\n", authReqBodyJSON)
	if err != nil {
		log.Fatalf("Error occured marashing map to json: %s\n", err)
		os.Exit(1)
	}
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", svcEndpoint.String(), bytes.NewBuffer(authReqBodyJSON))
	if err != nil {
		log.Fatalf("Unable to build request: %s\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/x-amz-json-1.1")
	req.Header.Set("X-Amz-Target", "AWSCognitoIdentityProviderService.InitiateAuth")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("InitateAuth response error: %v\n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading response body: %v\n", err)
		os.Exit(1)
	}

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

func main() {
	lambda.Start(login)
}
