package service

import (
	"context"
	"errors"

	"github.com/bricksocoolxd/bengi-investment-system/module/account/dto"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/model"
	"github.com/bricksocoolxd/bengi-investment-system/module/account/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrAccountNotFound      = errors.New("account not found")
	ErrAccountAlreadyExists = errors.New("account with this currency already exists")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrAccountFrozen        = errors.New("account is frozen")
	ErrSameAccount          = errors.New("cannot transfer to same account")
	ErrInvaludAmount        = errors.New("amount must be greater than 0")
)

type AccountService struct {
	repository *repository.AccountRepository
}

func NewAccountService(repository *repository.AccountRepository) *AccountService {
	return &AccountService{
		repository: repository,
	}
}

func (s *AccountService) CreateAccount(ctx context.Context, userID string, req *dto.CreateAccountRequest) (*dto.AccountResponse, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	exists, err := s.repository.ExistsByUserAndCurrency(ctx, userObjectID, req.Currency)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrAccountAlreadyExists
	}

	account := &model.Account{
		UserID:   userObjectID,
		Currency: req.Currency,
		Balance:  0,
	}
	if err := s.repository.Create(ctx, account); err != nil {
		return nil, err
	}

	return s.toAccountResponse(account), nil
}

func (s *AccountService) GetAccounts(ctx context.Context, userID string) ([]dto.AccountResponse, error) {
	accounts, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.AccountResponse
	for _, acc := range accounts {
		responses = append(responses, *s.toAccountResponse(&acc))
	}
	return responses, nil
}

func (s *AccountService) GetAccountByID(ctx context.Context, accountID, userID string) (*dto.AccountResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	// Check ownership
	if account.UserID.Hex() != userID {
		return nil, ErrAccountNotFound
	}

	return s.toAccountResponse(account), nil
}

func (s *AccountService) Deposit(ctx context.Context, accountID, userID string, req *dto.DepositRequest) (*dto.TransactionResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.UserID.Hex() != userID {
		return nil, ErrAccountNotFound
	}

	if account.Status == model.AccountStatusFrozen {
		return nil, ErrAccountFrozen
	}

	balanceBefore := account.Balance
	balanceAfter := account.Balance + req.Amount

	if err := s.repository.UpdateBalance(ctx, account.ID, balanceAfter); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		AccountID:     account.ID,
		Type:          model.TransactionTypeDeposit,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        model.TransactionStatusCompleted,
		Description:   req.Description,
	}

	if tx.Description == "" {
		tx.Description = "Deposit"
	}

	if err := s.repository.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return s.toTransactionResponse(tx), nil
}

func (s *AccountService) Withdraw(ctx context.Context, accountID, userID string, req *dto.WithdrawRequest) (*dto.TransactionResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.UserID.Hex() != userID {
		return nil, ErrAccountNotFound
	}

	if account.Status == model.AccountStatusFrozen {
		return nil, ErrAccountFrozen
	}

	if req.Amount <= 0 {
		return nil, ErrInvaludAmount
	}

	if req.Amount > account.Balance {
		return nil, ErrInsufficientBalance
	}

	balanceBefore := account.Balance
	balanceAfter := account.Balance - req.Amount

	if err := s.repository.UpdateBalance(ctx, account.ID, balanceAfter); err != nil {
		return nil, err
	}

	tx := &model.Transaction{
		AccountID:     account.ID,
		Type:          model.TransactionTypeWithdraw,
		Amount:        req.Amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Status:        model.TransactionStatusCompleted,
		Description:   req.Description,
	}

	if tx.Description == "" {
		tx.Description = "Withdraw"
	}

	if err := s.repository.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	return s.toTransactionResponse(tx), nil
}

func (s *AccountService) GetTransactions(ctx context.Context, accountID, userID string, limit, offset int) ([]dto.TransactionResponse, error) {
	account, err := s.repository.FindByID(ctx, accountID)
	if err != nil {
		return nil, ErrAccountNotFound
	}

	if account.UserID.Hex() != userID {
		return nil, ErrAccountNotFound
	}

	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	transactions, err := s.repository.GetTransactionsByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []dto.TransactionResponse
	for _, tx := range transactions {
		responses = append(responses, *s.toTransactionResponse(&tx))
	}

	return responses, nil
}

// Helper: Convert Account to AccountResponse
func (s *AccountService) toAccountResponse(acc *model.Account) *dto.AccountResponse {
	return &dto.AccountResponse{
		ID:       acc.ID.Hex(),
		UserID:   acc.UserID.Hex(),
		Currency: acc.Currency,
		Balance:  acc.Balance,
		Status:   string(acc.Status),
	}
}

// Helper: Convert Transaction to TransactionResponse
func (s *AccountService) toTransactionResponse(tx *model.Transaction) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		ID:            tx.ID.Hex(),
		AccountID:     tx.AccountID.Hex(),
		Type:          string(tx.Type),
		Amount:        tx.Amount,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		Status:        string(tx.Status),
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
