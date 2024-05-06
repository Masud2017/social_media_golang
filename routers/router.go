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

	router.GET("/addfriend",controller.AddFriend)
	router.GET("/addfather",controller.AddFather)
	router.GET("/addmother",controller.AddMother)
	router.GET("/addson",controller.AddSon)

	router.GET("/myrelationlist/:user_id",controller.MyRelationList)
	router.GET("/relationship_reqs/:user_id",controller.RelationShipRequests)
	router.GET("/my_relationship_reqs/:user_id",controller.MyRelationShipRequests)

	return router;
}