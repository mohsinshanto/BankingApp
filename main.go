package main

import (
	"banking/database"
	"banking/models"
	"banking/routes"
	"log"

	"github.com/gin-gonic/gin"
)


func main(){
	if err:=database.ConnectDB(); err != nil{
		log.Fatal(err)
	}
	if err:=database.DB.AutoMigrate(&models.User{},&models.Account{},&models.Transaction{});err != nil{
        log.Fatal(err)
	}
	router:= gin.Default()
	routes.RouteHandler(router)
	router.Run()
	
    
}


