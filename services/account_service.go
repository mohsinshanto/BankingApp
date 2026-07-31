package services

import (
	"banking/dto"
	appError "banking/errors"
	"banking/models"
	"banking/repositories"
	"banking/utils"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var validAccountStatus = map[string]bool{
	"ACTIVE":  true,
	"BLOCKED": true,
	"CLOSED":  true,
}

type AccountService interface {
	CreateAccount(userID uint) (*models.Account, error)
	Deposit(accountNo string, amount float64) (*models.Account, error)
	Withdraw(accountNo string, amount float64) (*models.Account, error)
	MoneyTransfer(transferInput *dto.TransferInput) (*dto.TransferResponse, error)
	AccountDetails(accountNo string) (*dto.AccountDetails, error)
	AccountStatusUpdate(accountNo, status string) (*models.Account, error)
	GetTransactionStat(accountNo string) (*dto.TransactionStatistics, error)
}
type accountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{
		repo: repo,
	}
}
func (s *accountService) GetTransactionStat(accountNo string) (*dto.TransactionStatistics, error) {
	_, err := s.repo.FindByAccountNo(accountNo)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTransactionStat(accountNo)

}
func (s *accountService) CreateAccount(userID uint) (*models.Account, error) {
	existingAccount, err := s.repo.FindByUserID(userID)
	if err == nil && existingAccount != nil {
		return nil, appError.ErrAccountAlreadyExists
	}
	accountNo, err := s.getUniqueAccountNumber()
	if err != nil {
		return nil, err
	}
	account := &models.Account{
		AccountNo: accountNo,
		UserID:    userID,
	}
	if err := s.repo.Create(account); err != nil {
		return nil, err
	}
	return account, nil
}
func (s *accountService) Deposit(accountNo string, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, appError.ErrInvalidAmount
	}
	tx := s.repo.BeginTx()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer tx.Rollback()
	account, err := s.repo.FindByAccountNoForUpdate(tx, accountNo)
	if err != nil {
		return nil, err
	}
	if account.Status != "ACTIVE" {
		return nil, appError.ErrInvalidStatus
	}

	account.Balance += amount
	if err := s.repo.Update(tx, account); err != nil {
		return nil, err
	}
	if err := s.createTransaction(tx, "", accountNo, amount, "DEPOSIT", "Money deposited"); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return account, nil

}
func (s *accountService) Withdraw(accountNo string, amount float64) (*models.Account, error) {
	if amount <= 0 {
		return nil, appError.ErrInvalidAmount
	}
	tx := s.repo.BeginTx()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer tx.Rollback()
	account, err := s.repo.FindByAccountNoForUpdate(tx, accountNo)
	if err != nil {
		return nil, err
	}
	if account.Status != "ACTIVE" {
		return nil, appError.ErrInvalidStatus
	}
	if account.Balance < amount {
		return nil, appError.ErrInsufficientBalance
	}
	account.Balance -= amount
	if err := s.repo.Update(tx, account); err != nil {
		return nil, err
	}
	if err := s.createTransaction(tx, accountNo, "", amount, "WITHDRAW", "Money withdrawn"); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return account, nil

}
func (s *accountService) MoneyTransfer(transferInput *dto.TransferInput) (*dto.TransferResponse, error) {
	if transferInput.Amount <= 0 {
		return nil, appError.ErrInvalidAmount
	}

	if transferInput.SenderAccNo == transferInput.ReceiverAccNo {
		return nil, appError.ErrSameAccountTransfer
	}
	tx := s.repo.BeginTx()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer tx.Rollback()
	senderAccount, err := s.findAccountForTransfer(tx, transferInput.SenderAccNo)
	if err != nil {
		return nil, err
	}
	receiverAccount, err := s.findAccountForTransfer(tx, transferInput.ReceiverAccNo)
	if err != nil {
		return nil, err
	}

	if senderAccount.Status != "ACTIVE" {
		return nil, appError.ErrSenderAccountNotActive
	}

	if receiverAccount.Status != "ACTIVE" {
		return nil, appError.ErrReceiverAccountNotActive
	}
	if transferInput.Amount > senderAccount.Balance {
		return nil, appError.ErrInsufficientBalance
	}
	senderAccount.Balance -= transferInput.Amount
	if err := s.repo.Update(tx, senderAccount); err != nil {
		return nil, err
	}
	receiverAccount.Balance += transferInput.Amount
	if err := s.repo.Update(tx, receiverAccount); err != nil {
		return nil, err
	}
	if err := s.createTransaction(tx, senderAccount.AccountNo, receiverAccount.AccountNo, transferInput.Amount, "TRANSFER", "Money transfer"); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return &dto.TransferResponse{
		SenderAccountNo:   senderAccount.AccountNo,
		ReceiverAccountNo: receiverAccount.AccountNo,
		TransferAmount:    transferInput.Amount,
		SenderNewBalance:  senderAccount.Balance,
	}, nil
}
func (s *accountService) AccountDetails(accountNo string) (*dto.AccountDetails, error) {
	account, err := s.repo.FindAccountDetails(accountNo)
	if err != nil {
		return nil, err
	}
	return &dto.AccountDetails{
		AccountNo: account.AccountNo,
		Balance:   account.Balance,
		Status:    account.Status,
		CreatedAt: account.CreatedAt,
		UserName:  account.User.Name,
		UserEmail: account.User.Email,
	}, nil

}
func (s *accountService) AccountStatusUpdate(accountNo, status string) (*models.Account, error) {
	status = strings.ToUpper(status)
	if _, ok := validAccountStatus[status]; !ok {
		return nil, appError.ErrInvalidStatus
	}

	account, err := s.repo.FindByAccountNo(accountNo)
	if err != nil {
		return nil, err
	}
	if account.Status == status {
		return nil, appError.ErrStatusAlreadySet
	}
	account.Status = status
	if err := s.repo.AccountStatusUpdate(account); err != nil {
		return nil, err
	}
	return account, nil

}

// private method not for outside use
func (s *accountService) createTransaction(tx *gorm.DB, fromAccount string, toAccount string, amount float64, transType string, description string) error {
	transaction := models.Transaction{
		FromAccount: fromAccount,
		ToAccount:   toAccount,
		Amount:      amount,
		Type:        transType,
		Description: description,
		Status:      "SUCCESS",
	}
	// trans db call
	return s.repo.CreateTransaction(tx, &transaction)
}
func (s *accountService) getUniqueAccountNumber() (string, error) {
	for i := 0; i < 5; i++ {
		accountNO, err := utils.GenerateAccountNo()
		if err != nil {
			return "", err
		}
		_, err = s.repo.FindByAccountNo(accountNO)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return accountNO, nil
		}
		if err != nil {
			return "", err
		}
	}
	return "", errors.New("couldn't generate unique account number")
}
func (s *accountService) findAccountForTransfer(tx *gorm.DB, accountNo string) (*models.Account, error) {
	account, err := s.repo.FindByAccountNoForUpdate(tx, accountNo)
	if err != nil {
		return nil, err
	}
	return account, nil
}
