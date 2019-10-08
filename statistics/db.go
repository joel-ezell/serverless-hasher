package statistics

import (
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

type HashTotals struct {
	Index int64  `json:"index"`
	Hash  string `json:"hash"`
}

func getStats(index int) (*HashTotals, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("statistics"),
		Key: map[string]*dynamodb.AttributeValue{
			"statsName": {
				S: aws.String("hashStats"),
			},
		},
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		log.Printf("Error recevied looking up stats from DB: %s", err)
		return new(HashTotals), err
	}
	if result.Item == nil {
		return new(HashTotals), nil
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	totals := new(HashTotals)
	err = dynamodbattribute.UnmarshalMap(result.Item, totals)
	if err != nil {
		log.Printf("Error received unmarshalling stats from DB: %s", err)
		return new(HashTotals), err
	}

	return totals, nil
}

// type statsItem struct {
// 	TotalCount    int `json:"totalCount"`
// 	TotalDuration int `json:"totalDuration"`
// }

func updateStats(newDuration int64) error {
	input := &dynamodb.UpdateItemInput{
		TableName:    aws.String("statistics"),
		ReturnValues: aws.String("UPDATED_NEW"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				N: aws.String(strconv.Itoa(1)),
			},
			":b": {
				N: aws.String(strconv.FormatInt(newDuration, 10)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#v": aws.String("totalCount"),
			"#w": aws.String("totalDuration"),
		},
		UpdateExpression: aws.String("[SET #v = #v + :a] [SET #w = #w + :b]"),
		Key: map[string]*dynamodb.AttributeValue{
			"counterName": {
				S: aws.String("importantCounter"),
			},
		},
	}

	_, err := db.UpdateItem(input)
	return err
}
