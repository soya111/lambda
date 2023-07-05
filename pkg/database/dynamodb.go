package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Subscriber struct {
	MemberName string `dynamo:"member_name,hash" json:"member_name"`
	UserId     string `dynamo:"user_id,range" json:"user_id"`
}

type Dynamo struct {
	db *dynamodb.DynamoDB
}

func NewDynamo() (*Dynamo, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	db := dynamodb.New(sess)
	return &Dynamo{db}, nil
}

func (d *Dynamo) GetDestination(memberName string) ([]string, error) {
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

	result, err := d.db.Query(params)
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
