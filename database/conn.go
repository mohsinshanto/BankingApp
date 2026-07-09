package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var DB *gorm.DB

func ConnectDB()error{
    var err error
	dsn := "root:root123@tcp(127.0.0.1:3306)/banking_system?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn),&gorm.Config{})
	if err != nil{
		return err
	}
	fmt.Println("Db connectec successfully")
	return nil
}