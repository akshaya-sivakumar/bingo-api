package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
	go h.run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId/:connectionType", func(c *gin.Context) {
		roomId := c.Param("roomId")
		connectionType := c.Param("connectionType")
		serveWs(c.Writer, c.Request, roomId,connectionType)
	})

	router.Run("0.0.0.0:8080")
}
