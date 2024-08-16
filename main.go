package bot

import "github.com/gin-gonic/gin"

func main() {
	// Your code here
	router := gin.Default()
	router.GET("/status", getServerStatus)
}

func getServerStatus(c *gin.Context) {
	// Your code here

}
