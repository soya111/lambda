package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Subscriber struct {
	MemberName string `dynamodbav:"member_name" json:"member_name"`
	UserId     string `json:"user_id" dynamodbav:"user_id"`
}

func GetDestination(memberName string) ([]string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)

	params := &dynamodb.QueryInput{
		TableName:              aws.String("Subscriber"),
		KeyConditionExpression: aws.String("#member_name = :member_name"),
		ExpressionAttributeNames: map[string]*string{
			"#member_name": aws.String("member_name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":member_name": {
				S: aws.String(memberName),
			},
		},
	}

	result, err := db.Query(params)
	if err != nil {
		return nil, err
	}

	users := []Subscriber{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return nil, err
	}

	res := []string{}
	for _, user := range users {
		res = append(res, user.UserId)
	}

	return res, nil
}
