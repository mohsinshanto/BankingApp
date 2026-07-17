package controllers

import (
	"banking/database"
	"banking/dto"
	"banking/models"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTransactionsByAccount(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "2")
	transactionType := c.Query("type")
	fromDate := c.Query("from")
	toDate := c.Query("to")
	sortBy := c.DefaultQuery("sort", "newest")
	var fromTime time.Time
	var toTime time.Time

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 2
	}

	if limit > 100 {
		limit = 100
	}

	accountNo := c.Param("accountNo")
	offset := (page - 1) * limit

	var account models.Account
	if err := database.DB.Where("account_no=?", accountNo).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User Account Not Found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	query := database.DB.Model(&models.Transaction{})
	query = query.Where("from_account=? OR to_account=?", accountNo, accountNo)

	if transactionType != "" {
		transactionType = strings.ToUpper(transactionType)

		switch transactionType {
		case "DEPOSIT", "WITHDRAW", "TRANSFER":
			query = query.Where("type=?", transactionType)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction type"})
			return
		}
	}
	if fromDate != "" {
		var err error
		fromTime, err = time.Parse("2006-01-02", fromDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format.Use YYYY-MM-DD"})
			return
		}
	}
	if toDate != "" {
		var err error
		toTime, err = time.Parse("2006-01-02", toDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format.Use YYYY-MM-DD"})
			return
		}
	}
	if !fromTime.IsZero() && !toTime.IsZero() {
		if fromTime.After(toTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "from date can't be after to date"})
			return
		}
	}
	if !fromTime.IsZero() {
		query = query.Where("created_at >= ?", fromTime)
	}
	if !toTime.IsZero() {
		toTime = toTime.Add(24*time.Hour - 1*time.Nanosecond)
		query = query.Where("created_at <= ?", toTime)
	}
	sortBy = strings.ToLower(sortBy)

	var orderBy string

	switch sortBy {
	case "newest":
		orderBy = "created_at DESC"

	case "oldest":
		orderBy = "created_at ASC"

	case "amount_asc":
		orderBy = "amount ASC"

	case "amount_desc":
		orderBy = "amount DESC"

	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid sort option",
		})
		return
	}
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transactions []models.Transaction
	err = query.
		Order(orderBy).
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"account_no":         account.AccountNo,
		"current_page":       page,
		"limit":              limit,
		"total_pages":        (total + int64(limit) - 1) / int64(limit),
		"total_transactions": total,
		"transactions":       transactions,
	})
}
func GetAccountSummary(c *gin.Context) {
	accountNo := c.Param("accountNo")
	var account models.Account
	if err := database.DB.Where("account_no=?", accountNo).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// present account balance
	currentBalance := account.Balance
	// total deposit
	var totalDeposit float64
	err := database.DB.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("to_account=? AND type=?", accountNo, "DEPOSIT").
		Scan(&totalDeposit).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// total withdraw
	var totalWithdraw float64
	err = database.DB.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("from_account=? AND type=?", accountNo, "WITHDRAW").
		Scan(&totalWithdraw).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var totalTransferSent float64
	err = database.DB.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("from_account=? AND type=?", accountNo, "TRANSFER").
		Scan(&totalTransferSent).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var totalTransferReceived float64
	err = database.DB.Model(&models.Transaction{}).Select("COALESCE(SUM(amount),0)").Where("to_account=? AND type=?", accountNo, "TRANSFER").
		Scan(&totalTransferReceived).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var totalTransactions int64
	err = database.DB.Model(&models.Transaction{}).Where("from_account=? OR to_account=?", accountNo, accountNo).
		Count(&totalTransactions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := dto.AccountSummaryResponse{
		AccountNo:             accountNo,
		CurrentBalance:        currentBalance,
		TotalTransactions:     totalTransactions,
		TotalDeposit:          totalDeposit,
		TotalWtihdraw:         totalWithdraw,
		TotalTransferSent:     totalTransferSent,
		TotalTransferReceived: totalTransferReceived,
	}
	c.JSON(http.StatusOK, response)
}
func GetTransactionStatistics(c *gin.Context) {
	accountNo := c.Param("accountNo")
	var account models.Account
	if err := database.DB.Where("account_no=?", accountNo).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var todayCount int64
	if err := database.DB.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND DATE(created_at)= CURDATE()", accountNo, accountNo).
		Count(&todayCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var thisWeek int64
	if err := database.DB.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND YEARWEEK(created_at, 1)= YEARWEEK(CURDATE(), 1)", accountNo, accountNo).
		Count(&thisWeek).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var thisMonth int64
	if err := database.DB.
		Model(&models.Transaction{}).
		Where("(from_account= ? OR to_account= ?) AND YEAR(created_at)= YEAR(CURDATE()) AND MONTH(created_at)= MONTH(CURDATE())", accountNo, accountNo).
		Count(&thisMonth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var totalTransactions int64
	if err := database.DB.Model(&models.Transaction{}).Where("from_account= ? OR to_account= ?", accountNo, accountNo).
		Count(&totalTransactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account_no":         account.AccountNo,
		"today_transactions": todayCount,
		"week_transactions":  thisWeek,
		"month_transactions": thisMonth,
		"total_transactions": totalTransactions,
	})

}
func AccountDetails(c *gin.Context) {
	accountNo := c.Param("accountNo")
	var account models.Account
	if err := database.DB.Preload("User", func(db *gorm.DB) *gorm.DB { return db.Select("id", "name", "email") }).
		Where("account_no=?", accountNo).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account Not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"account_no": account.AccountNo,
		"balance":    account.Balance,
		"status":     account.Status,
		"created_at": account.CreatedAt,
		"user": gin.H{
			"name":  account.User.Name,
			"email": account.User.Name,
		},
	})
}
func AccountStatusUpdate(c *gin.Context) {
	accountNo := c.Param("accountNo")
	var inputStatus dto.AccountStatusUpdate
	if err := c.ShouldBindJSON(&inputStatus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	inputStatus.Status = strings.ToUpper(inputStatus.Status)
	switch inputStatus.Status {
	case "ACTIVE", "BLOCKED", "CLOSED":
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status response",
		})
		return
	}
	var account models.Account
	if err := database.DB.Where("account_no=?", accountNo).Take(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if account.Status == inputStatus.Status {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account is already in this status",
		})
		return
	}
	account.Status = inputStatus.Status
	if err := database.DB.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Account status updated successfully",
		"account_no":     account.AccountNo,
		"account_status": account.Status,
	})
}
