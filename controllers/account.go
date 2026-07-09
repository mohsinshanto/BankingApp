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
func Deposit(c *gin.Context){
	var InputDepo dto.Deposit
	if err:=c.ShouldBindJSON(&InputDepo); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}
	var account models.Account
	err:=database.DB.Where("account_no=?",InputDepo.AccountNo).Take(&account).Error
	if err != nil{
		if errors.Is(err,gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound,gin.H{"error":err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	account.Balance += InputDepo.Amount
	err=database.DB.Save(&account).Error
	if err!= nil{
		c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"msg":"Successfully Deposited",
		"NewBalance": account.Balance,
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