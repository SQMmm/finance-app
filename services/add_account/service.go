package add_account

import (
	"context"
	"fmt"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
)

type AccountAdder interface {
	AddAccount(context.Context, *entities.Account) (*entities.Account, error)
}

type service struct {
	manager           logger.LoggerManager
	accountRepository accountRepository
}

func NewService(m logger.LoggerManager, accounts accountRepository) *service {
	return &service{
		manager:           m,
		accountRepository: accounts,
	}
}

func (s *service) AddAccount(ctx context.Context, account *entities.Account) (*entities.Account, error) {
	log := s.manager.LogCtx(ctx)
	var err error

	log.Infof("adding new account: %#v", account)
	account.ID, err = s.accountRepository.Add(ctx, account)
	if err != nil {
		return nil, fmt.Errorf("failed to add account: %s", err)
	}
	log.Debugf("account was added with id=%v", account.ID)

	return account, nil
}
