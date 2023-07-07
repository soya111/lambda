package dynamodb

import (
	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type DynamoDiaryRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

func NewDynamoDiaryRepository(sess *session.Session, tableName string) *DynamoDiaryRepository {
	db := dynamo.New(sess)
	table := db.Table(tableName)
	return &DynamoDiaryRepository{
		db:    db,
		table: table,
	}
}

func (r *DynamoDiaryRepository) GetDiary(memberName string, diaryId int) (*model.Diary, error) {
	diary := new(model.Diary)
	err := r.table.Get("member_name", memberName).Range("diary_id", dynamo.Equal, diaryId).One(diary)
	if err != nil {
		if err == dynamo.ErrNotFound {
			return nil, model.ErrDiaryNotFound
		}
		return nil, err
	}
	return diary, nil
}

func (r *DynamoDiaryRepository) PostDiary(diary *model.Diary) error {
	return r.table.Put(diary).Run()
}
