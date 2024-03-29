package dynamodb

import (
	"zephyr/pkg/model"

	"github.com/guregu/dynamo"
)

// DiaryRepository is the struct that represents the repository of diary.
type DiaryRepository struct {
	db    *dynamo.DB
	table dynamo.Table
}

// NewDiaryRepository receives a DynamoDB instance and a table name, and returns a new DiaryRepository.
func NewDiaryRepository(db *dynamo.DB, tableName string) *DiaryRepository {
	table := db.Table(tableName)
	return &DiaryRepository{
		db:    db,
		table: table,
	}
}

// GetDiary returns the diary of the specified member and diary ID.
func (r *DiaryRepository) Get(memberName string, diaryId int) (*model.Diary, error) {
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

// PostDiary posts the diary.
func (r *DiaryRepository) Post(diary *model.Diary) error {
	return r.table.Put(diary).Run()
}
