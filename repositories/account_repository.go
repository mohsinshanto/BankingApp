package repositories

import (
	appError "banking/errors"
	"banking/models"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByUserID(userID uint) (*models.Account, error)
	FindByAccountNo(accountNo string) (*models.Account, error)
	FindAccountDetails(accountNo string) (*models.Account, error)
	FindByAccountNoForUpdate(tx *gorm.DB, accountNo string) (*models.Account, error)
	Update(tx *gorm.DB, account *models.Account) error
	BeginTx() *gorm.DB
	CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error
}
type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {

	return &accountRepository{
		db: db,
	}

}
func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}
func (r *accountRepository) FindByUserID(userID uint) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("user_id=?", userID).Take(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *accountRepository) FindByAccountNoForUpdate(tx *gorm.DB, accountNo string) (*models.Account, error) {
	var account models.Account
	err := tx.Clauses(clause.Locking{Strength: "update"}).Where("account_no=?", accountNo).Take(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appError.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}
func (r *accountRepository) Update(tx *gorm.DB, account *models.Account) error {
	return tx.Save(account).Error
}
func (r *accountRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}
func (r *accountRepository) CreateTransaction(tx *gorm.DB, transaction *models.Transaction) error {
	return tx.Create(transaction).Error
}
func (r *accountRepository) FindByAccountNo(accountNo string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("account_no = ?", accountNo).Take(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}
func (r *accountRepository) FindAccountDetails(accountNo string) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "name", "email") }).
		Where("account_no=?", accountNo).Take(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appError.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}
