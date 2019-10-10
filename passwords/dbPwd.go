package passwords

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Declare a new DynamoDB instance. Note that this is safe for concurrent
// use.
var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

type hashItem struct {
	Index int64  `json:"index"`
	Hash  string `json:"hash"`
}

func getHash(index int) (string, error) {
	// Prepare the input for the query.
	input := &dynamodb.GetItemInput{
		TableName: aws.String("hashes"),
		Key: map[string]*dynamodb.AttributeValue{
			"index": {
				N: aws.String(strconv.Itoa(index)),
			},
		},
	}

	// Retrieve the item from DynamoDB. If no matching item is found
	// return nil.
	result, err := db.GetItem(input)
	if err != nil {
		return "", err
	}
	if result.Item == nil {
		return "", nil
	}

	// The result.Item object returned has the underlying type
	// map[string]*AttributeValue. We can use the UnmarshalMap helper
	// to parse this straight into the fields of a struct. Note:
	// UnmarshalListOfMaps also exists if you are working with multiple
	// items.
	hash := new(hashItem)
	err = dynamodbattribute.UnmarshalMap(result.Item, hash)
	if err != nil {
		return "", err
	}

	return hash.Hash, nil
}

// Add a hash to DynamoDB.
func putHash(index int, hash string) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("hashes"),
		Item: map[string]*dynamodb.AttributeValue{
			"index": {
				N: aws.String(strconv.Itoa(index)),
			},
			"hash": {
				S: aws.String(hash),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}

type counterItem struct {
	CurrentValue int `json:"currentValue"`
}

const indexTableName string = "counters"
const indexValueName string = "currentValue"
const indexKeyName string = "hashIndex"

func nextIndex() int {
	input := &dynamodb.UpdateItemInput{
		TableName:    aws.String(indexTableName),
		ReturnValues: aws.String("UPDATED_NEW"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				N: aws.String(strconv.Itoa(1)),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#v": aws.String(indexValueName),
		},
		UpdateExpression: aws.String("SET #v = #v + :a"),
		Key: map[string]*dynamodb.AttributeValue{
			"counterName": {
				S: aws.String(indexKeyName),
			},
		},
	}

	result, err := db.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}

	counter := new(counterItem)
	err = dynamodbattribute.UnmarshalMap(result.Attributes, counter)
	if err != nil {
		return 0
	}

	return counter.CurrentValue
}
