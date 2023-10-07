package main

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (app *application) server() error {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/vterminal", app.VirtualTerminal)

	return router.Run(":" + strconv.Itoa(app.config.port))
}
