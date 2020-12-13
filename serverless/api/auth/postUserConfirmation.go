package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sony/sonyflake"
)

// genSonyFlake generated unique ID
func genSonyFlake() uint64 {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() faield with %s\n", err)
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
	fmt.Printf("Generated userId (Sonyflake): %d", userID)
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
