package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type Subscriber struct {
	MemberName string `dynamo:"member_name,hash" json:"member_name"`
	UserId     string `dynamo:"user_id,range" json:"user_id"`
}

type SubscriberRepository interface {
	GetAllByMemberName(memberName string) ([]string, error)
	Subscribe(subscriber Subscriber) error
	Unsubscribe(memberName, userId string) error
	GetAllById(id string) ([]Subscriber, error)
}

type DynamoSubscriberRepository struct {
	db *dynamo.DB
}

func NewDynamoSubscriberRepository(sess *session.Session) SubscriberRepository {
	db := dynamo.New(sess)
	return &DynamoSubscriberRepository{db}
}

func (d *DynamoSubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	table := d.db.Table("Subscriber")

	var subscribers []Subscriber
	err := table.Get("member_name", memberName).All(&subscribers)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, len(subscribers))
	for i, sub := range subscribers {
		userIds[i] = sub.UserId
	}

	return userIds, nil
}

func (d *DynamoSubscriberRepository) Subscribe(subscriber Subscriber) error {
	table := d.db.Table("Subscriber")

	if err := table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

func (d *DynamoSubscriberRepository) Unsubscribe(memberName, userId string) error {
	table := d.db.Table("Subscriber")

	err := table.Delete("memberName", memberName).Range("userId", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}

	return nil
}

func (d *DynamoSubscriberRepository) GetAllById(id string) ([]Subscriber, error) {
	table := d.db.Table("Subscriber")

	var res []Subscriber
	err := table.Get("userId", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}

	return res, nil
}
