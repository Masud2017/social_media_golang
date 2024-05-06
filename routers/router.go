package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/Masud2017/social_media_golang/controllers"
)



func SetupRouter()  *gin.Engine {
	router := gin.Default();


	
	controller := new(controllers.Controller)

	router.GET("/",controller.Index)
	
	router.GET("/signup",controller.Signup)
	router.GET("/userlist",controller.UserList)
	router.GET("/me/:my_id",controller.Me)
	router.GET("/acceptreq",controller.AcceptReq)
	router.GET("/cancelreq",controller.CancelReq)

	router.GET("/addrelation",controller.AddRelation)

	router.GET("/myrelationlist",controller.MyRelationList)
	router.GET("/relationship_reqs/:user_id",controller.RelationShipRequests)
	router.GET("/my_relationship_reqs",controller.RelationShipRequests)

	return router;
}