package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
)

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	log.Printf("request body: %s", request.Body)
	log.Printf("Proxxy request: %+v\n", request)

	// If no name is provided in the HTTP request body, throw an error
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	var jsonBodyInterface interface{}
	err := json.Unmarshal([]byte(request.Body), &jsonBodyInterface)

	if err != nil {
		return events.APIGatewayProxyResponse{},
			errors.New("not a valid json")
	}

	jsonBody := jsonBodyInterface.(map[string]interface{})
	log.Printf("request body interface: %s", jsonBody)
	responseBody, _ := json.Marshal(request)
	return events.APIGatewayProxyResponse{
		Body:       string(responseBody),
		StatusCode: 200,
	}, nil
}

func main() {
	fmt.Println("hello world")
	lambda.Start(Handler)
}
