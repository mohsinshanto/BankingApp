package controllers

import (
	"banking/database"
	"banking/dto"
	"banking/models"
	"banking/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateAccount(c *gin.Context){
	value, ok := c.Get("user-id")
	if !ok {
		c.JSON(http.StatusUnauthorized,gin.H{"error":"user id not found"})
	}
	userID:=uint(value.(float64))
	var existingAccount models.Account
	err:=database.DB.Where("user_id=?",userID).Take(&existingAccount).Error; 
	if err == nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":"user already exists"})
		return
	}
	if err != gorm.ErrRecordNotFound{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	accountNo,err := getUniqueAccNo()
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	newAccount := models.Account{
		AccountNo: accountNo,
		UserID: userID,
	}
	if err:=database.DB.Create(&newAccount).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusCreated,gin.H{
		"msg": "New Account created successfully",
		"AccNo":newAccount.AccountNo,
	})



}
// No outside access
func getUniqueAccNo()(string,error){
	const maxAttempts = 5
	for i:=0;i<maxAttempts;i++{
		accountNo, err := utils.GenerateAccountNo()
		if err != nil{
            return "",err
		}
		var account models.Account
		err = database.DB.Where("account_no=?",accountNo).Take(&account).Error
		if errors.Is(err,gorm.ErrRecordNotFound){
			return accountNo,nil
		}
		if err !=nil{
			return "",err
		}
	}
	return "",errors.New("Couldn't generate a unique account number")
}
func Deposit(c *gin.Context){
	var InputDepo dto.Deposit
	if err:=c.ShouldBindJSON(&InputDepo); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	tx:= database.DB.Begin()
	if tx.Error != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":tx.Error.Error()})
		return
	}
	var account models.Account
	err:=tx.Clauses(clause.Locking{Strength: "update"}).Where("account_no=?",InputDepo.AccountNo).Take(&account).Error
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			tx.Rollback()
			c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	account.Balance += InputDepo.Amount
	err=tx.Save(&account).Error
	if err!= nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if err:=tx.Commit().Error;err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Successfully Deposited",
		"NewBalance": account.Balance,
	})
}
func Withdraw(c *gin.Context){
	var withdrawInput dto.Withdraw
	if err:= c.ShouldBindJSON(&withdrawInput); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	tx:= database.DB.Begin()
	if tx.Error != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":tx.Error.Error()})
		return
	}
	var account models.Account
	err:=tx.Clauses(clause.Locking{Strength: "update"}).Where("account_no=?",withdrawInput.AccountNo).Take(&account).Error
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			tx.Rollback()
			c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if account.Balance < withdrawInput.Amount{
		tx.Rollback()
		c.JSON(http.StatusBadRequest,gin.H{"msg":"Insufficient Account Balance"})
		return
	}
	account.Balance -= withdrawInput.Amount
	err=tx.Save(&account).Error
	if err!= nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if err:= tx.Commit().Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Successfully Wtihdrawn",
		"NewBalance": account.Balance,
	})


}
func MoneyTransfer(c *gin.Context){
	var transferInput dto.TransferInput
	if err:= c.ShouldBindJSON(&transferInput); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	tx := database.DB.Begin()
	if tx.Error != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": tx.Error.Error(),
    })
    return
}
	var senderAccount models.Account
	if err:=tx.Clauses(clause.Locking{Strength: "update"}).
	Where("account_no=?",transferInput.SenderAccNo).
	Take(&senderAccount).Error; err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			tx.Rollback()
			c.JSON(http.StatusNotFound,gin.H{"error":"Sender doesn't exists"})
			return

		}
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	var receiverAccount models.Account
	if err:= tx.Clauses(clause.Locking{Strength: "update"}).
	Where("account_no=?",transferInput.ReceiverAccNo).
	Take(&receiverAccount).Error; err  != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
           tx.Rollback()
		   c.JSON(http.StatusNotFound,gin.H{"error":"Receiver doesn't exists"})
		   return
		}
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if transferInput.SenderAccNo == transferInput.ReceiverAccNo {
    tx.Rollback()

    c.JSON(http.StatusBadRequest, gin.H{
        "error":"sender and receiver cannot be same",
    })

    return
}
	if senderAccount.Status != "ACTIVE"|| receiverAccount.Status != "ACTIVE"{
		tx.Rollback()
		c.JSON(http.StatusBadRequest,gin.H{"error":"both account must be active state"})
		return
	}

	if transferInput.Amount > senderAccount.Balance{
		tx.Rollback()
		c.JSON(http.StatusBadRequest,gin.H{"error":"Insufficient Balance"})
		return

	}
	senderAccount.Balance -= transferInput.Amount
	if err:= tx.Save(&senderAccount).Error; err != nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	receiverAccount.Balance += transferInput.Amount
	if err:= tx.Save(&receiverAccount).Error; err != nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	if err := tx.Commit().Error; err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
    })
    return
}
	c.JSON(http.StatusOK,gin.H{"SenderNewBalance":senderAccount.Balance,"ReceiverNewBalance":receiverAccount.Balance})
}