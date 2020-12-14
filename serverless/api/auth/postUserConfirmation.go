package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
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
	}
	return id
}

// newUserHandler is executed when a new user is confirmed into the Cognito User Pool
func newUserHandler(event events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open(dbDriver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userID := genSonyFlake()
	fmt.Printf("Generated userId (Sonyflake): %d\n", userID)
	fmt.Printf("UserName: %s\n", event.UserName)
	fmt.Printf("UserAttributes: %s\n", event.Request.UserAttributes)

	userName := event.UserName
	userEmail := event.Request.UserAttributes["email"]
	res, err := db.Exec("INSERT INTO users VALUES(?, ?, ?)", userID, userName, userEmail)

	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, _ := res.RowsAffected()
	fmt.Printf("Inserted new user: %d rows changed.", rowsAffected)

	return event, nil
}

func main() {
	lambda.Start(newUserHandler)
}
