package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ServiceConfig struct {
	serviceUrl  url.URL
	contentType string
}

var (
	// ErrNameNotProvided is thrown when a name is not provided
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")

	services = map[string]ServiceConfig{
		"sla": {
			contentType: "multipart/form-data",
			serviceUrl: url.URL{
				Host:   "www.tumacredo.com",
				Path:   "api_v1/billing/",
				Scheme: "https",
			},
		},
	}
)

/*
Get the mpesa service to send the transaction
*/
func getMpesaRecipientService(suffixAccount string) (ServiceConfig, error) {

	if serviceUrl, ok := services[suffixAccount]; ok {
		return serviceUrl, nil
	}
	return ServiceConfig{}, errors.New("N")
}

func sendRequest(serviceConfig ServiceConfig, body map[string]interface{}) (map[string]interface{}, int, error) {

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	var requestBody *bytes.Buffer

	if serviceConfig.contentType == "multipart/form-data" {
		requestBody = &bytes.Buffer{}
		writer := multipart.NewWriter(requestBody)
		serviceConfig.contentType = writer.FormDataContentType()
		for k, v := range body {
			err := writer.WriteField(k, v.(string))
			if err != nil {
				return nil, -1, err
			}
		}
		_ = writer.Close()
	} else {
		bytesRepresentation, err := json.Marshal(body)
		if err != nil {
			return nil, -1, err
		}

		requestBody = bytes.NewBuffer(bytesRepresentation)
	}

	// now forward the request to the mpesa service account
	resp, err := client.Post(serviceConfig.serviceUrl.String(),
		serviceConfig.contentType, requestBody)

	if err != nil {
		return nil, -1, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, resp.StatusCode, err
}

// Handler is your Lambda function handler
// It uses Amazon API Gateway request/responses provided by the aws-lambda-go/events package,
// However you could use other event sources (S3, Kinesis etc), or JSON-decoded primitive types such as 'string'.
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// stdout and stderr are sent to AWS CloudWatch Logs
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	log.Printf("request body: %s", request.Body)
	log.Printf("Proxxy request: %+v\n", request)

	// We will always return a 200 even if we get an error during processing
	gatewayResponse := events.APIGatewayProxyResponse{Body: "", StatusCode: 200}

	var jsonBodyInterface interface{}
	err := json.Unmarshal([]byte(request.Body), &jsonBodyInterface)

	if err != nil {
		log.Print("Invalid json: ", err)
		return gatewayResponse, nil
	}

	jsonBody := jsonBodyInterface.(map[string]interface{})
	log.Printf("request body interface: %s", jsonBody)

	// get BillRefNumber
	billRefNumber := jsonBody["BillRefNumber"].(string)

	// get suffix account
	billRefAsArray := strings.Split(billRefNumber, ".")
	account := billRefAsArray[len(billRefAsArray)-1]

	// now get the mpesa service to forward the request
	serviceUrl, err := getMpesaRecipientService(account)

	if err != nil {
		log.Print("err:", err)
		return gatewayResponse, nil
	}

	log.Print("sending request", "url", serviceUrl, "body", jsonBody)
	response, statusCode, err := sendRequest(serviceUrl, jsonBody)
	log.Print("response:", response, "statusCode:", statusCode, "err:", err)
	return gatewayResponse, nil
}

func main() {
	fmt.Println("starting lambda")
	lambda.Start(Handler)
}
