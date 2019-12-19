package add_account_group

import (
	"context"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

type AccountGroupAdder interface {
	AddAccountGroup(ctx context.Context, group *entities.AccountGroup) (*entities.AccountGroup, error)
}

type service struct {
	manager     logger.LoggerManager
	accountRepo accountRepository
	groupRepo   groupRepository
}

func NewService(m logger.LoggerManager, a accountRepository, g groupRepository) *service {
	return &service{
		manager:     m,
		accountRepo: a,
		groupRepo:   g,
	}
}

func (s *service) AddAccountGroup(ctx context.Context, group *entities.AccountGroup) (*entities.AccountGroup, error) {
	var err error

	if len(group.Accounts) > 0 {
		group.Accounts, err = s.accountRepo.GetUserAccountsByIDs(ctx, group.User.ID, group.Accounts.GetIDs())
		if err != nil {
			return nil, fmt.Errorf("failed to get accounts: %s", err)
		}
	}

	group.ID, err = s.groupRepo.Add(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("failed to add group: %s", err)
	}

	return group, nil
}
