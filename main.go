package main

import (
	"GachaServerGin/implements"
	"GachaServerGin/src"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	data, err := implements.NewMongoData()
	if err != nil {
		src.Logger.Error(err)
	}
	analyst := implements.NewMemAnalyst(data)
	//service := src.NewGachaService(data)
	app := gin.Default()
	app.GET("/users", func(context *gin.Context) {
		//service.data.GetUsers()
		context.JSON(200, data.GetUsers())
	})
	app.GET("/gachas", func(context *gin.Context) {
		page, err := strconv.Atoi(context.DefaultQuery("page", "0"))
		if err != nil {
			context.Status(400)
		}
		uid := context.Query("uid")
		context.JSON(200, data.GetGachasByPage(uid, page, 10))
	})

	app.GET("/analysis", func(context *gin.Context) {
		uid := context.Query("uid")
		context.JSON(200, analyst.Analysis(uid))
	})

	analyst.Analyze("69023059")

	err = app.Run(":8000")
	if err != nil {
		src.Logger.Error(err)
	}
}
