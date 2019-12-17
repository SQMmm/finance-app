package add_tag

import (
	"context"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

type TagAdder interface {
	AddTag(ctx context.Context, tag *entities.Tag) (*entities.Tag, error)
}

type service struct {
	manager  logger.LoggerManager
	tagsRepo tagRepository
}

func NewService(m logger.LoggerManager, t tagRepository) *service {
	return &service{
		manager:  m,
		tagsRepo: t,
	}
}

func (s *service) AddTag(ctx context.Context, tag *entities.Tag) (*entities.Tag, error) {
	log := s.manager.LogCtx(ctx)
	var err error

	log.Infof("adding new tag: %+v", tag)
	tag.ID, err = s.tagsRepo.Add(ctx, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to add tag: %s", err)
	}
	log.Debugf("new tag was added with id=%v", tag.ID)

	return tag, nil
}
