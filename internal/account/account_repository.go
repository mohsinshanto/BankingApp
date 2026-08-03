package account

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
	AccountStatusUpdate(account *models.Account) error
	GetTransactionStat(accountNo string) (*TransactionStatistics, error)
	GetAccountSummary(accountNo string) (*AccountSummary, error)
	GetTransactionsByAccount(accountNo string, filter *TransactionFilter) (*TransactionQueryResult, error)
}
type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {

	return &accountRepository{
		db: db,
	}

}
func (r *accountRepository) GetTransactionsByAccount(accountNo string, filter *TransactionFilter) (*TransactionQueryResult, error) {
	query := r.db.Model(&models.Transaction{})

	query = query.Where(
		"from_account = ? OR to_account = ?",
		accountNo,
		accountNo,
	)

	if filter.TransactionType != "" {
		query = query.Where("type = ?", filter.TransactionType)
	}

	if !filter.From.IsZero() {
		query = query.Where("created_at >= ?", filter.From)
	}

	if !filter.To.IsZero() {
		query = query.Where("created_at <= ?", filter.To)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var transactions []models.Transaction

	if err := query.
		Order(filter.OrderBy).
		Limit(filter.PageSize).
		Offset(filter.Offset).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	return &TransactionQueryResult{
		Total:        total,
		Transactions: transactions,
	}, nil
}
func (r *accountRepository) GetAccountSummary(accountNo string) (*AccountSummary, error) {
	var totalDeposit, totalWithdraw, totalTransferSent, totalTransferReceived float64
	var totalTransactions int64
	if err := r.db.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("to_account=? AND type=?", accountNo, "DEPOSIT").
		Scan(&totalDeposit).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("from_account=? AND type=?", accountNo, "WITHDRAW").
		Scan(&totalWithdraw).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("from_account=? AND type=?", accountNo, "TRANSFER").
		Scan(&totalTransferSent).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("to_account=? AND type=?", accountNo, "TRANSFER").
		Scan(&totalTransferReceived).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Transaction{}).Where("from_account=? OR to_account=?", accountNo, accountNo).
		Count(&totalTransactions).Error; err != nil {
		return nil, err
	}
	return &AccountSummary{
		TotalTransactions:     totalTransactions,
		TotalDeposit:          totalDeposit,
		TotalWithdraw:         totalWithdraw,
		TotalTransferSent:     totalTransferSent,
		TotalTransferReceived: totalTransferReceived,
	}, nil
}
func (r *accountRepository) GetTransactionStat(accountNo string) (*TransactionStatistics, error) {
	var todayCount, weekCount, thisMonthCount, totalCount int64
	if err := r.db.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND DATE(created_at)= CURDATE()", accountNo, accountNo).
		Count(&todayCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND YEARWEEK(created_at, 1)= YEARWEEK(CURDATE(), 1)", accountNo, accountNo).
		Count(&weekCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND YEAR(created_at)= YEAR(CURDATE()) AND MONTH(created_at)= MONTH(CURDATE())", accountNo, accountNo).
		Count(&thisMonthCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&models.Transaction{}).Where("from_account= ? OR to_account= ?", accountNo, accountNo).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	return &TransactionStatistics{
		AccountNo:            accountNo,
		TodayTransaction:     todayCount,
		ThisWeekTransaction:  weekCount,
		ThisMonthTransaction: thisMonthCount,
		TotalTransaction:     totalCount,
	}, nil
}
func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}
func (r *accountRepository) FindByUserID(userID uint) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("user_id=?", userID).Take(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appError.ErrAccountNotFound
		}
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
func (r *accountRepository) AccountStatusUpdate(account *models.Account) error {
	return r.db.Save(account).Error
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appError.ErrAccountNotFound
		}
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
