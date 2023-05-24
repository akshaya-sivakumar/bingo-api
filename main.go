package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	go h.run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/room/:roomId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:roomId/:connectionType/:connectionLimit", func(c *gin.Context) {
		roomId := c.Param("roomId")
		connectionType := c.Param("connectionType")
		connectionLimit := c.Param("connectionLimit")
		x, err := strconv.ParseInt(connectionLimit, 10, 64)
		if err != nil {
			return
			//b.Error(err)
		}
		serveWs(c.Writer, c.Request, roomId, connectionType, x)
	})

	router.Run("0.0.0.0:8080")
}
