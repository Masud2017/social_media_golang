package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/Masud2017/social_media_golang/controllers"
)



func SetupRouter()  *gin.Engine {
	router := gin.Default();


	indexController := new(controllers.IndexController)
	controller := new(controllers.Controller)
	

	router.GET("/", indexController.Index)
	router.GET("/helloworld2",controller.Signup)

	return router;
}