package dynamodb

import (
	"notify/pkg/model"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type DiaryRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

func NewDiaryRepository(sess *session.Session, tableName string) *DiaryRepository {
	db := dynamo.New(sess)
	table := db.Table(tableName)
	return &DiaryRepository{
		db:    db,
		table: table,
	}
}

func (r *DiaryRepository) GetDiary(memberName string, diaryId int) (*model.Diary, error) {
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

func (r *DiaryRepository) PostDiary(diary *model.Diary) error {
	return r.table.Put(diary).Run()
}
