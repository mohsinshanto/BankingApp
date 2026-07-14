package controllers

import (
	"banking/database"
	"banking/models"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func GetTransactionsByAccount(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "2")
	transactionType := c.Query("type")

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

	var total int64
	err = query.Count(&total).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transactions []models.Transaction
	err = query.
		Order("created_at DESC").
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
	
