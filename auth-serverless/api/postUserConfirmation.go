package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/sony/sonyflake"
)

// genMachineID: Sony implementation doesn't work with zero valued Settings
func genMachineID() (uint16, error) {
	as, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}

	var ip net.IP
	for _, a := range as {
		ipnet, ok := a.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}
		ip = ipnet.IP.To4()
		break // only need to find one valid IP
	}

	if ip == nil {
		return 0, errors.New("unable to find interface address")
	}

	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// genSonyFlake generated unique ID
func genSonyFlake() uint64 {
	settings := sonyflake.Settings{MachineID: genMachineID}
	flake := sonyflake.NewSonyflake(settings)
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
		os.Exit(1)
	}
	return id
}

// newUserHandler is executed when a new user is confirmed into the Cognito User Pool
func newUserHandler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	userID := genSonyFlake()
	fmt.Printf("Generated userId (Sonyflake): %x\n", userID)
	fmt.Printf("UserName: %s\n", event.UserName)
	fmt.Printf("UserAttributes: %s\n", event.Request.UserAttributes)

	userIDStr := fmt.Sprintf("%x", userID)
	userPoolID := event.UserPoolID
	userName := event.UserName
	customAttrNameID := "dev:custom:id"
	userAttrID := cognitoidentityprovider.AttributeType{Name: &customAttrNameID, Value: &userIDStr}

	userAttrs := []*cognitoidentityprovider.AttributeType{&userAttrID}
	userAttrInput := cognitoidentityprovider.AdminUpdateUserAttributesInput{
		UserAttributes: userAttrs,
		UserPoolId:     &userPoolID,
		Username:       &userName,
	}
	// Create a Session with a custom region
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION")),
	})
	if err != nil {
		log.Fatalf("Error occurred while creating session.")
		os.Exit(1)
	}
	// Create a CognitoIdentityProvider client from a session.
	svc := cognitoidentityprovider.New(sess)
	resp, err := svc.AdminUpdateUserAttributes(&userAttrInput)
	if err != nil {
		log.Fatalf("Error occured updating user attribute: %s\n", err)
	}
	log.Printf("User attribute updated with response: %s\n", resp)
	return event, nil
}

func main() {
	lambda.Start(newUserHandler)
}
