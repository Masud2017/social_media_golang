package controllers

import (
	"github.com/gin-gonic/gin"
	
)

type IndexController struct {

}



func (indexController *IndexController) Index(c *gin.Context) {

	c.JSON(200, gin.H{
		"data":"hello",
	})
}