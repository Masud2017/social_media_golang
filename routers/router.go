package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/Masud2017/social_media_golang/controllers"
)



func SetupRouter()  *gin.Engine {
	router := gin.Default();


	
	controller := new(controllers.Controller)

	router.GET("/",controller.Index)
	
	router.GET("/signup",controller.Signup) // done
	router.GET("/userlist",controller.UserList) // done
	router.GET("/me/:my_id",controller.Me) // done
	router.GET("/acceptreq/:user_id/:req_id",controller.AcceptReq) // done
	router.GET("/cancelreq/:user_id/:req_id",controller.CancelReq) // done

	router.GET("/addfriend",controller.AddFriend) // done
	router.GET("/addfather",controller.AddFather) // done
	router.GET("/addmother",controller.AddMother) // done
	router.GET("/addson",controller.AddSon) // done

	router.GET("/myrelationlist/:user_id",controller.MyRelationList) //done
	router.GET("/relationship_reqs/:user_id",controller.RelationShipRequests)
	router.GET("/my_relationship_reqs/:user_id",controller.MyRelationShipRequests)

	return router;
}