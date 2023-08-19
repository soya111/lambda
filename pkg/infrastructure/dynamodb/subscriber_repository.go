package dynamodb

import (
	"fmt"

	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

// SubscriberRepository is the struct that represents the repository of subscriber.
type SubscriberRepository struct {
	db *dynamo.DB
}

// NewSubscriberRepository receives a session and returns a new SubscriberRepository.
func NewSubscriberRepository(sess *session.Session) model.SubscriberRepository {
	db := dynamo.New(sess)
	return &SubscriberRepository{db}
}

// GetAllByMemberName returns the list of user IDs that subscribe the specified member.
func (d *SubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	table := d.db.Table("Subscriber")

	var subscribers []model.Subscriber
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

// Subscribe inserts the subscriber.
func (d *SubscriberRepository) Subscribe(subscriber model.Subscriber) error {
	table := d.db.Table("Subscriber")

	if err := table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

// Unsubscribe deletes the subscriber.
func (d *SubscriberRepository) Unsubscribe(memberName, userId string) error {
	table := d.db.Table("Subscriber")

	err := table.Delete("member_name", memberName).Range("user_id", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}

	return nil
}

// GetAllById returns the list of subscribers that the specified user ID subscribes.
func (d *SubscriberRepository) GetAllById(id string) ([]model.Subscriber, error) {
	table := d.db.Table("Subscriber")

	var res []model.Subscriber
	err := table.Get("user_id", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}

	return res, nil
}
