package dynamodb

import (
	"fmt"

	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type SubscriberRepository struct {
	db *dynamo.DB
}

func NewSubscriberRepository(sess *session.Session) model.SubscriberRepository {
	db := dynamo.New(sess)
	return &SubscriberRepository{db}
}

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

func (d *SubscriberRepository) Subscribe(subscriber model.Subscriber) error {
	table := d.db.Table("Subscriber")

	if err := table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

func (d *SubscriberRepository) Unsubscribe(memberName, userId string) error {
	table := d.db.Table("Subscriber")

	err := table.Delete("member_name", memberName).Range("user_id", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}

	return nil
}

func (d *SubscriberRepository) GetAllById(id string) ([]model.Subscriber, error) {
	table := d.db.Table("Subscriber")

	var res []model.Subscriber
	err := table.Get("user_id", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}

	return res, nil
}
