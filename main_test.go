package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestGettingRecipientService(t *testing.T) {
	serviceUrl, err := getMpesaRecipientService("sla")

	assert.Equal(t, err, nil)
	assert.Equal(t, serviceUrl.Host, "www.tumacredo.com")
	assert.Equal(t, serviceUrl.Scheme, "https")
	assert.Equal(t, serviceUrl.Path, "api_v1/billing/")
}

func TestSendRequest(t *testing.T) {
	response, statusCode, err := sendRequest(
		url.URL{Scheme: "http", Host: "httpbin", Path: "/post"},
		"name")

	assert.Equal(t, 200, statusCode)
	assert.Equal(t, err, nil)
	assert.Equal(t, "\"name\"", response["data"])
}

func TestHandler(t *testing.T) {

	tests := []struct {
		request events.APIGatewayProxyRequest
		expect  string
		err     error
	}{
		{
			// Test that the handler responds with the correct response
			// when a valid name is provided in the HTTP body
			request: events.APIGatewayProxyRequest{Body: "Paul"},
			expect:  "",
			err:     nil,
		},
	}

	for _, test := range tests {
		response, err := Handler(test.request)
		assert.IsType(t, test.err, err)
		assert.Equal(t, test.expect, response.Body)
	}

}
