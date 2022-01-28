package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Id   string `json:"user_id" dynamodbav:"user_id"`
	Name string `json:"name" dynamodbav:"name"`
}

func GetDestination() ([]string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	params := &dynamodb.ScanInput{
		TableName: aws.String("User"),
	}
	result, err := db.Scan(params)
	if err != nil {
		return nil, err
	}

	users := []User{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, user := range users {
		res = append(res, user.Id)
	}

	return res, nil
}
