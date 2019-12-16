package data

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/sqmmm/finance-app/entities"
)

type Category struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Title  string `json:"title"`
	Icon   *Icon  `json:"icon,omitempty"`
	Color  *Color `json:"color,omitempty"`
}

func (c *Category) Validate() error {
	err := validation.Errors{
		"title length is incorrect": validation.Validate(
			c.Title,
			validation.Required,
			validation.RuneLength(1, 120),
		),
	}
	return err.Filter()
}

func (c *Category) GetEntity() *entities.Category {
	cat := &entities.Category{Title: c.Title}

	if c.Icon != nil {
		cat.Icon = &entities.Icon{ID: c.Icon.ID}
	}
	if c.Color != nil {
		cat.Color = &entities.IconColor{ID: c.Color.ID}
	}
	return cat
}

func GetCategoryFromEntity(c *entities.Category) *Category {
	cat := &Category{
		ID:     c.ID,
		Title:  c.Title,
		UserID: c.User.ID,
	}
	if c.Icon != nil {
		cat.Icon = &Icon{
			ID:   c.Icon.ID,
			Name: c.Icon.Name,
			Path: c.Icon.Path,
		}
	}
	if c.Color != nil {
		cat.Color = &Color{
			ID:    c.Color.ID,
			Name:  c.Color.Name,
			Color: c.Color.Color,
		}
	}

	return cat
}
