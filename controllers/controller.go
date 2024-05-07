package controllers

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"bytes"
	"log"
	"github.com/Masud2017/social_media_golang/db"
	"github.com/Masud2017/social_media_golang/models"
	// "encoding/json"
	
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

func (controller *Controller) Index(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": "Hello world this is the social media graph api",
	})
}


func (controller *Controller) Signup(c *gin.Context) {
	db := db.DB{}

	// db.InitSchema()	
	name := c.Query("name") 
	email := c.Query("email")
	password := c.Query("password")

	db.NewClient()
	db.NewClient()
	user := models.User{
		Uid:      "_:user",
		Name:     name,
		Email:    email,
		Password: password,
		// Friend: []models.Relation{{
		// 	Rel: "Father",
		// 	User: models.User{
		// 		Uid : "_:mk",
		// 		Name: "MK",

		// 	}, // Assuming UserId is the ID of the friend
		// }},
	}


	if (db.SignupUser(&user)) {
		fmt.Printf("Hello world your name is :%s , email : %s and password is : %s",name,email,password)
		infof("YOYO")

		c.JSON(200, gin.H{
			"name":name,
			"email":email,
			"password":password,
		})

	} else {
		c.JSON(200, gin.H{
			"data": "The user that you have mentioned is already exists in the db please use an unique user",
		})
	}
}


func (controller *Controller) UserList(c *gin.Context) {
	db := db.DB{}
	db.NewClient()

	userList := db.GetUserList()

	// res,_ :=json.MarshalIndent(userList, "", "\t")
	c.JSON(200, gin.H{
		"data":userList,
	})

}

func (controller *Controller) Me(c *gin.Context) {
	my_id := c.Param("my_id")
	db := db.DB{}
	db.NewClient()
	fmt.Println("My id is : "+my_id)
	
	me := db.Me(my_id)

	

	c.JSON(200, gin.H{
		"data": me,
	})
}

// func (controller *Controller) RequestRelationship(c *gin.Context) {

// }

func (controller *Controller) AddFriend(c *gin.Context) {
	my_id := c.Query("my_id")
	rel := "Friend"
	req_to := c.Query("req_to")


	db := db.DB{}
	db.NewClient()

	me := db.Me(my_id)
	req_to_user := db.Me(req_to)

	relReq := models.RelationRequest{
		Uid : "_:reqrel",
		ReqRel : rel,
		ReqTo : req_to_user,
	}

	relReq = db.RequestForRelationship(relReq,me)

	c.JSON(200, gin.H{
		"data": relReq,
	})

}

func (controller *Controller) AddFather(c *gin.Context) {
	my_id := c.Query("my_id")
	rel := "Father"
	req_to := c.Query("req_to")


	db := db.DB{}
	db.NewClient()

	me := db.Me(my_id)
	req_to_user := db.Me(req_to)

	relReq := models.RelationRequest{
		Uid : "_:req",
		ReqRel : rel,
		ReqTo : req_to_user,
	}

	relReq = db.RequestForRelationship(relReq,me)

	c.JSON(200, gin.H{
		"data": relReq,
	})

}

func (controller *Controller) AddMother(c *gin.Context) {
	my_id := c.Query("my_id")
	rel := "Mother"
	req_to := c.Query("req_to")


	db := db.DB{}
	db.NewClient()

	me := db.Me(my_id)
	req_to_user := db.Me(req_to)

	relReq := models.RelationRequest{
		Uid : "_:req",
		ReqRel : rel,
		ReqTo : req_to_user,
	}

	relReq = db.RequestForRelationship(relReq,me)

	c.JSON(200, gin.H{
		"data": relReq,
	})

}

func (controller *Controller) AddSon(c *gin.Context) {
	my_id := c.Query("my_id")
	rel := "Son"
	req_to := c.Query("req_to")


	db := db.DB{}
	db.NewClient()

	me := db.Me(my_id)
	req_to_user := db.Me(req_to)

	relReq := models.RelationRequest{
		Uid : "_:req",
		ReqRel : rel,
		ReqTo : req_to_user,
	}

	relReq = db.RequestForRelationship(relReq,me)

	c.JSON(200, gin.H{
		"data": relReq,
	})

}


func (controller *Controller) AcceptReq(c *gin.Context) {
	user_id := c.Param("user_id")
	req_id := c.Param("req_id")

 
	db := db.DB{}
	db.NewClient()

	acceptRequestStatus := db.AcceptReq(user_id,req_id)

	c.JSON(200, gin.H{
		"data": acceptRequestStatus,
	})
}

func (controller *Controller) CancelReq(c *gin.Context) {
	user_id := c.Param("user_id")
	req_id := c.Param("req_id")

 
	db := db.DB{}
	db.NewClient()

	cancelRequestStatus := db.CancelReq(user_id,req_id)

	c.JSON(200, gin.H{
		"data": cancelRequestStatus,
	})
}

func (controller *Controller) MyRelationList(c *gin.Context) {
	user_id := c.Param("user_id")

	db := db.DB{}
	db.NewClient()

	fmt.Println(user_id)

	friendList,fatherList,motherList,sonList := db.MyRelationList(user_id)

	relationList := [][]models.Relation{friendList,fatherList,motherList,sonList}

	// res,_ :=json.MarshalIndent(userList, "", "\t")
	c.JSON(200, gin.H{
		"data":relationList,
	})

}

func (controller *Controller) RelationShipRequests(c *gin.Context) {
	user_id := c.Param("user_id")

	db := db.DB{}
	db.NewClient()

	fmt.Println(user_id)

	relationShipReqList := db.RelationShipRequests(user_id)

	// res,_ :=json.MarshalIndent(userList, "", "\t")
	c.JSON(200, gin.H{
		"data":relationShipReqList,
	})
}

func (controller *Controller) MyRelationShipRequests(c *gin.Context) {
	user_id := c.Param("user_id")

	db := db.DB{}
	db.NewClient()

	fmt.Println(user_id)

	relationShipReqList := db.MyRelationShipRequests(user_id)

	// res,_ :=json.MarshalIndent(userList, "", "\t")
	c.JSON(200, gin.H{
		"data":relationShipReqList,
	})
}