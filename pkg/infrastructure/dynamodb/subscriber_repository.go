package dynamodb

import (
	"fmt"

	"notify/pkg/model"

	"github.com/guregu/dynamo"
)

const subscriberTableName = "Subscriber"

// SubscriberRepository is the struct that represents the repository of subscriber.
type SubscriberRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

// NewSubscriberRepository receives a DynamoDB instance and returns a new SubscriberRepository.
func NewSubscriberRepository(db *dynamo.DB) *SubscriberRepository {
	table := db.Table(subscriberTableName)
	return &SubscriberRepository{db, table}
}

// GetAllByMemberName returns the list of user IDs that subscribe the specified member.
func (r *SubscriberRepository) GetAllByMemberName(memberName string) ([]string, error) {
	var subscribers []model.Subscriber
	err := r.table.Get("member_name", memberName).All(&subscribers)
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
func (r *SubscriberRepository) Subscribe(subscriber model.Subscriber) error {
	if err := r.table.Put(subscriber).Run(); err != nil {
		return fmt.Errorf("postSubscriber: %w", err)
	}

	return nil
}

// Unsubscribe deletes the subscriber.
func (r *SubscriberRepository) Unsubscribe(memberName, userId string) error {
	err := r.table.Delete("member_name", memberName).Range("user_id", userId).Run()

	if err != nil {
		return fmt.Errorf("deleteSubscriber: %w", err)
	}

	return nil
}

// GetAllById returns the list of subscribers that the specified user ID subscribes.
func (r *SubscriberRepository) GetAllById(id string) ([]model.Subscriber, error) {
	var res []model.Subscriber
	err := r.table.Get("user_id", id).Index("user_id-index").All(&res)
	if err != nil {
		return nil, fmt.Errorf("getSubscribeList: %w", err)
	}

	return res, nil
}

// DeleteAllById deletes all the subscribers that the specified user ID subscribes.
func (r *SubscriberRepository) DeleteAllById(id string) error {
	var subscribers []model.Subscriber

	// クエリで特定のuserIdのレコードを取得
	err := r.table.Get("user_id", id).Index("user_id-index").All(&subscribers)
	if err != nil {
		return fmt.Errorf("querying by user_id: %w", err)
	}

	for _, subscriber := range subscribers {
		fmt.Printf("deleting: %s\n", subscriber.MemberName)
	}

	// 取得した各レコードを削除
	for _, subscriber := range subscribers {
		err := r.table.Delete("member_name", subscriber.MemberName).Range("user_id", id).Run()
		if err != nil {
			return fmt.Errorf("deleting by member_name and user_id: %w", err)
		}
	}

	return nil
}
