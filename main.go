package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Status struct {
	Status string `json:"status"`
}

var status Status = Status{Status: "ok"}

func main() {
	//Your code here
	router := gin.Default()
	router.GET("/status", getServerStatus)
	router.Run(":8080")
}

func getServerStatus(c *gin.Context) {
	// Your code here
	c.IndentedJSON(http.StatusOK, status)
}
