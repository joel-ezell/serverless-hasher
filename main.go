package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joel-ezell/serverless-hasher/passwords"
	"github.com/joel-ezell/serverless-hasher/statistics"
)

var srv http.Server

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("HTTP method is this: %s\n", req.HTTPMethod)
	fmt.Printf("Path is %s\n", req.Path)

	parts := strings.Split(req.Path, "/")
	fmt.Printf("First path element: %s", parts[1])

	// In some ways, it would be cleaner to have a separate Lambda function for each resource (hash and stats).
	// That would make it so that we don't have to have this level of routing in our code; API Gateway would take
	// care of that. I will likely make that change soon.
	switch parts[1] {
	case "hash":
		return hashRouter(req)
	case "stats":
		return statsRouter(req)
	default:
		return clientError(http.StatusNotFound)
	}
}

func hashRouter(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return getHash(req)
	case "POST":
		return putHash(req)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func statsRouter(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		stats, err := statistics.GetStats()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       stats,
		}, err
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}

func getHash(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var indexStr string
	indexStr = req.PathParameters["index"]
	index, _ := strconv.Atoi(indexStr)
	hashedPwd, err := passwords.GetHashedPassword(index)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       hashedPwd,
	}, err
}

func putHash(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Body is: %s", req.Body)
	parts := strings.Split(req.Body, "password=")
	fmt.Printf("Hopefully this is the value: %s", parts[len(parts)-1])
	index, err := passwords.HashAndStore(parts[len(parts)-1])

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       strconv.Itoa(index),
	}, err
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
