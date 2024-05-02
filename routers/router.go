package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/Masud2017/social_media_golang/controllers"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/Masud2017/social_media_golang/docs"
	
)



func SetupRouter()  *gin.Engine {
	router := gin.Default();


	indexController := new(controllers.IndexController)


	docs.SwaggerInfo.BasePath = "/api/v1"
   v1 := router.Group("/api/v1")
   {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",indexController.Index)
      }
   }
	
//    url := ginSwagger.URL("http://localhost:4443/docs/swagger.json") // The url pointing to API definition

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/", indexController.Index)



	

	return router;
}