package main

import (
	"GachaServerGin/implements"
	"GachaServerGin/src"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
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
	app.GET("/users.invalid", func(context *gin.Context) {
		context.JSON(200, []string{"TODO"})
	})
	app.GET("/gachas", func(context *gin.Context) {
		page, err := strconv.Atoi(context.DefaultQuery("page", "0"))
		if err != nil {
			context.Status(400)
		}
		uid := context.Query("uid")
		gachas := data.GetGachasByPage(uid, page, 10)
		for index, gacha := range gachas.List {
			gachas.List[index].NickName = data.GetUser(gacha.Uid).NickName
		}
		context.JSON(200, gin.H{
			"data": gachas,
		})
	})

	go func() {
		analyst.Analyze("")
		for _, user := range data.GetUsers() {
			analyst.Analyze(user.Uid)
		}
	}()
	app.GET("/analysis", func(context *gin.Context) {
		uid := context.Query("uid")
		context.JSON(200, analyst.Analysis(uid))
	})

	service := src.NewGachaService(data, analyst)
	go func() {
		for {
			for _, user := range data.GetUsers() {
				src.Logger.Infof("user: %s begin", user.NickName)
				service.UpdateChannel <- user.Uid
				time.Sleep(time.Minute)
			}
		}
	}()

	err = app.Run(":8000")
	if err != nil {
		src.Logger.Error(err)
	}
}
