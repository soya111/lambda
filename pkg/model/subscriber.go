package model

// Subscriber is a struct that represents a subscriber.
type Subscriber struct {
	MemberName string `dynamo:"member_name,hash" json:"member_name"`
	UserId     string `dynamo:"user_id,range" json:"user_id"`
}

// SubscriberRepository is a interface that represents a repository of subscribers.
type SubscriberRepository interface {
	GetAllByMemberName(memberName string) ([]string, error)
	Subscribe(subscriber Subscriber) error
	Unsubscribe(memberName, userId string) error
	GetAllById(id string) ([]Subscriber, error)
}
