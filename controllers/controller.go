package controllers

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"bytes"
	"log"
	
)


var (
	buf    bytes.Buffer
	logger = log.New(&buf, "INFO: ", log.Lshortfile)

	infof = func(info string) {
		logger.Output(2, info)
	}
)



type Controller struct {

}

func (controller *Controller) Signup(c *gin.Context) {
	
	name := c.Query("name") 
	email := c.Query("email")
	password := c.Query("password")

	fmt.Printf("Hello world your name is :%s , email : %s and password is : %s",name,email,password)
	infof("YOYO")

	c.JSON(200, gin.H{
		"name":name,
		"email":email,
		"password":password,
	})
}
