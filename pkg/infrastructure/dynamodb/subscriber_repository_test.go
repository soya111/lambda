package dynamodb

import (
	"notify/pkg/model"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
)

func TestSubscriberRepository(t *testing.T) {
	t.Skip("skipping this test for now")
	DYNAMO_ENDPOINT := "http://localhost:8000"

	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(DYNAMO_ENDPOINT),
	})
	assert.NoError(t, err)

	subscriber := NewSubscriberRepository(sess)

	s := model.Subscriber{
		MemberName: "小坂菜緒",
		UserId:     "U1234567890",
	}
	err = subscriber.Subscribe(s)
	assert.NoError(t, err)
	defer func() {
		err := subscriber.Unsubscribe("小坂菜緒", "U1234567890")
		assert.NoError(t, err)
	}()

	a, err := subscriber.GetAllById("U1234567890")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(a))
	assert.Equal(t, "小坂菜緒", a[0].MemberName)
	assert.Equal(t, "U1234567890", a[0].UserId)

	b, err := subscriber.GetAllByMemberName("小坂菜緒")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(b))
	assert.Equal(t, "U1234567890", b[0])

}
