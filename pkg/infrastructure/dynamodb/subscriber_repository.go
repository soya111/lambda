package dynamodb

import (
	"fmt"

	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

const subscriberTableName = "Subscriber"

// SubscriberRepository is the struct that represents the repository of subscriber.
type SubscriberRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

// NewSubscriberRepository receives a session and returns a new SubscriberRepository.
func NewSubscriberRepository(sess *session.Session) model.SubscriberRepository {
	db := dynamo.New(sess)
	table := db.Table(subscriberTableName)
	return &SubscriberRepository{db, table}
}

// GetAllByMemberName returns the list of user IDs that subscribe the specified member.
func (d *SubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	var subscribers []model.Subscriber
	err := d.table.Get("member_name", memberName).All(&subscribers)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, len(subscribers))
	for i, sub := range subscribers {
		userIds[i] = sub.UserId
	}

	return userIds, nil
}

// Subscribe inserts the subscriber.
func (d *SubscriberRepository) Subscribe(subscriber model.Subscriber) error {
	if err := d.table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

// Unsubscribe deletes the subscriber.
func (d *SubscriberRepository) Unsubscribe(memberName, userId string) error {
	err := d.table.Delete("member_name", memberName).Range("user_id", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}

	return nil
}

// GetAllById returns the list of subscribers that the specified user ID subscribes.
func (d *SubscriberRepository) GetAllById(id string) ([]model.Subscriber, error) {
	var res []model.Subscriber
	err := d.table.Get("user_id", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}

	return res, nil
}
