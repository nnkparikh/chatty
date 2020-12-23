package idp

import (
	"fmt"
	"os"
)

// REGION : identifies AWS Region
var REGION string = os.Getenv("REGION")

// CognitoSvcEndpoint : cognito service endpoint
var CognitoSvcEndpoint string = fmt.Sprintf("cognito-idp.%s.amazonaws.com", REGION)
