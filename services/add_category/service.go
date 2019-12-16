package add_category

import (
	"context"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

type CategoryAdder interface {
	AddCategory(ctx context.Context, category *entities.Category) (*entities.Category, error)
}

type service struct {
	manager            logger.LoggerManager
	categoryRepository categoryRepository
	iconRepository     iconRepository
}

func NewService(m logger.LoggerManager, c categoryRepository, i iconRepository) *service {
	return &service{manager: m, categoryRepository: c, iconRepository: i}
}

func (s service) AddCategory(ctx context.Context, category *entities.Category) (*entities.Category, error) {
	log := s.manager.LogCtx(ctx)
	var err error

	//todo: check log (%+v) or (%#v)
	log.Infof("adding new category: %+v", category)
	category.ID, err = s.categoryRepository.Add(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to add category: %s", err)
	}
	log.Debugf("category was added with id=%v", category.ID)

	if category.Icon != nil {
		category.Icon, err = s.iconRepository.GetIconByID(ctx, category.Icon.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get icon: %s", err)
		}

		if category.Color != nil {
			category.Color, err = s.iconRepository.GetColorByID(ctx, category.Color.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get color: %s", err)
			}
		}
	}

	return category, nil
}
