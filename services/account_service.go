package services

import (
	"banking/dto"
	appError "banking/errors"
	"banking/models"
	"banking/repositories"
	"banking/utils"
	"errors"

	"gorm.io/gorm"
)

type AccountService interface {
	CreateAccount(userID uint) (*models.Account, error)
	Deposit(accountNo string, amount float64) (*models.Account, error)
	Withdraw(accountNo string, amount float64) (*models.Account, error)
	MoneyTransfer(transferInput *dto.TransferInput) (*dto.TransferResponse, error)
}
type accountService struct {
	repo repositories.AccountRepository
}

func NewAccountService(repo repositories.AccountRepository) AccountService {
	return &accountService{
		repo: repo,
	}
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
	tx := s.repo.BeginTx()
	if tx.Error != nil {
		return nil, tx.Error
	}
	account, err := s.repo.FindByAccountNoForUpdate(tx, accountNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, appError.ErrAccountNotFound
		}
		tx.Rollback()
		return nil, err
	}
	if account.Status != "ACTIVE" {
		tx.Rollback()
		return nil, appError.ErrInvalidStatus
	}
	if amount <= 0 {
		tx.Rollback()
		return nil, appError.ErrInvalidAmount
	}
	account.Balance += amount
	if err := s.repo.Update(tx, account); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := s.createTransaction(tx, "", accountNo, amount, "DEPOSIT", "Money deposited"); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return account, nil

}
func (s *accountService) Withdraw(accountNo string, amount float64) (*models.Account, error) {
	tx := s.repo.BeginTx()
	if tx.Error != nil {
		return nil, tx.Error
	}
	account, err := s.repo.FindByAccountNoForUpdate(tx, accountNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, appError.ErrAccountNotFound
		}
		tx.Rollback()
		return nil, err
	}
	if account.Status != "ACTIVE" {
		tx.Rollback()
		return nil, appError.ErrInvalidStatus
	}
	if amount <= 0 {
		tx.Rollback()
		return nil, appError.ErrInvalidAmount
	}
	if account.Balance < amount {
		tx.Rollback()
		return nil, appError.ErrInsufficientBalance
	}
	account.Balance -= amount
	if err := s.repo.Update(tx, account); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := s.createTransaction(tx, accountNo, "", amount, "WITHDRAW", "Money withdrawn"); err != nil {
		tx.Rollback()
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
	if tx.Error != nil {
		return nil, tx.Error
	}

	senderAccount, err := s.repo.FindByAccountNoForUpdate(tx, transferInput.SenderAccNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, appError.ErrAccountNotFound
		}
		tx.Rollback()
		return nil, err
	}
	receiverAccount, err := s.repo.FindByAccountNoForUpdate(tx, transferInput.ReceiverAccNo)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, appError.ErrAccountNotFound
		}
		tx.Rollback()
		return nil, err
	}

	if senderAccount.Status != "ACTIVE" {
		tx.Rollback()
		return nil, appError.ErrSenderAccountNotActive
	}

	if receiverAccount.Status != "ACTIVE" {
		tx.Rollback()
		return nil, appError.ErrReceiverAccountNotActive
	}
	if transferInput.Amount > senderAccount.Balance {
		tx.Rollback()
		return nil, appError.ErrInsufficientBalance
	}
	senderAccount.Balance -= transferInput.Amount
	if err := s.repo.Update(tx, senderAccount); err != nil {
		tx.Rollback()
		return nil, err
	}
	receiverAccount.Balance += transferInput.Amount
	if err := s.repo.Update(tx, receiverAccount); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := s.createTransaction(tx, senderAccount.AccountNo, receiverAccount.AccountNo, transferInput.Amount, "TRANSFER", "Money transfer"); err != nil {
		tx.Rollback()
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
