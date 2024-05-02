package controllers

import (
	"github.com/gin-gonic/gin"
	
)

type IndexController struct {

}


// @BasePath /api/v1

// PingExample godoc
// @Summary YOYO example
// @Schemes
// @Description do nothing
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router / [get]
func (indexController *IndexController) Index(c *gin.Context) {
	// og := new(db.OGM)
	// user := models.User{}
	// user.Name = "Masud karim"
	// user.Age = 34
	// session,err := og.CreateSession()

	// if (err) {
	// 	Println("Error occured")
	// }

	// if err := session.Save(&user, nil); err != nil {
	// 	panic(err)
	// }


	c.JSON(200, gin.H{
		"data":"hello",
	})
}