package main
import (
	"github.com/Masud2017/social_media_golang/routers"
)

func main() {
	
	router := routers.SetupRouter()

	router.Run(":4443")
	
} 