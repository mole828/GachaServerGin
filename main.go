package main

import (
	"GachaServerGin/implements"
	"GachaServerGin/src"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func main() {
	data, mongoErr := implements.NewMongoData()
	if mongoErr != nil {
		src.Logger.Error(mongoErr)
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
		analysis := analyst.Analysis(uid)
		user := data.GetUser(uid)
		user.Token = ""
		analysis.User = user
		context.JSON(200, analysis)
	})

	service := src.NewGachaService(data, analyst)
	var lastTime, lastLastTime time.Time
	go func() {
		lastTime = time.Now()
		for {
			for _, user := range data.GetUsers() {
				//src.Logger.Infof("user: %s begin, uid: %s", user.NickName, user.Uid)
				service.UpdateChannel <- user.Uid
				time.Sleep(time.Minute)
			}
			lastLastTime = lastTime
			lastTime = time.Now()
		}
	}()
	app.GET("/updates", func(context *gin.Context) {
		context.JSON(200, service.UpdateTimes)
	})
	app.GET("/users.invalid", func(context *gin.Context) {
		pick := lo.PickBy(service.UpdateTimes, func(key string, value time.Time) bool {
			src.Logger.WithFields(logrus.Fields{
				"lastTime":                    lastTime,
				"lastLastTime":                lastLastTime,
				"value":                       value,
				"lastLastTime.Sub(value) > 0": lastLastTime.Sub(value) > 0,
			}).Info()
			return lastLastTime.Sub(value) > 0
		})
		src.Logger.WithField("pick", pick).Info()
		context.JSON(200, lo.Keys(pick))
	})
	app.POST("/register", func(context *gin.Context) {
		token := context.Query("token")
		var msg = "ok"
		regErr := service.NewToken(token)
		if regErr != nil {
			msg = regErr.Error()
		}
		context.JSON(200, gin.H{
			"msg":   msg,
			"token": token,
		})
	})

	mongoErr = app.Run(":8000")
	if mongoErr != nil {
		src.Logger.Error(mongoErr)
	}
}
