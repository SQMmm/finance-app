package data

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/sqmmm/finance-app/entities"
)

type Tag struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
}

func (t *Tag) Validate() error {
	err := validation.Errors{
		"title length is incorrect": validation.Validate(
			t.Title,
			validation.Required,
			validation.RuneLength(1, 120),
		),
	}
	return err.Filter()
}

func (t *Tag) GetEntity() *entities.Tag {
	return &entities.Tag{
		Title: t.Title,
	}
}

func GetTagFromEntity(tag *entities.Tag) *Tag {
	return &Tag{
		ID:     tag.ID,
		UserID: tag.User.ID,
		Title:  tag.Title,
	}
}
