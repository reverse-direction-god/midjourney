package router

import (
	"fmt"
	"mj/controller"
	"mj/model"
	"mj/service"
	"mj/until"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 解决跨域问题
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
func check(timeStr string) bool {
	// 要检查的时间字符串
	targetTimeString := timeStr

	// 使用给定的时间格式解析字符串
	layout := "2006-01-02T15:04:05-07:00" // 时间格式
	targetTime, err := time.Parse(layout, targetTimeString)
	if err != nil {
		fmt.Println("解析时间错误:", err)
		return false
	}

	// 获取当前时间
	currentTime := time.Now()

	// 检查目标时间是否在当前时间之前（过期）
	if targetTime.Before(currentTime) {
		return false
	} else {
		return true
	}
}
func checkMyself(timeStr string) bool {
	if timeStr == "" {
		return false
	}
	s := timeStr[0:4]
	ints, _ := strconv.Atoi(s)
	if ints >= 2024 {
		return true
	} else {
		return false
	}
}

func CheckHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.UserInfo
		token := c.GetHeader("token")
		id := until.Decrypt(token)
		if id == "" {
			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "Token not found or expired  id==nil",
				Data: nil,
			})
			return
		}
		// id, err := service.Redis.Get(service.Redis.Context(), token).Result()
		// fmt.Println(id)
		// if err != nil {
		// 	var mod = new(model.Error)
		// 	mod.Message = "redis没找到" + token
		// 	service.DB.Model(&model.Error{}).Create(&mod)
		// 	c.AbortWithStatusJSON(500, model.Response{
		// 		Code: 500,
		// 		Msg:  "Token not found or expired:::" + err.Error(),
		// 		Data: nil,
		// 	})
		// 	return
		// }

		// 如果没有错误，则继续查询数据库
		service.DB.Model(&model.UserInfo{}).Where("id=?", id).Find(&user)
		if user.ID == 0 {
			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "User not found",
				Data: nil,
			})
			return
		}

		if check(user.TokenTime) == false {

			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "Token expired",
				Data: nil,
			})
			return
		}
		c.Set("user", id)
		// 如果一切正常，继续处理请求
		c.Next()
	}
}

func CheckHandlerMyself() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.UserInfo
		token := c.GetHeader("token")
		id := until.Decrypt(token)
		if id == "" {
			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "Token not found or expired  id==nil",
				Data: nil,
			})
			return
		}
		// id, err := service.Redis.Get(service.Redis.Context(), token).Result()
		// fmt.Println(id)
		// if err != nil {
		// 	var mod = new(model.Error)
		// 	mod.Message = "本地redis没找到" + token
		// 	service.DB.Model(&model.Error{}).Create(&mod)
		// 	c.AbortWithStatusJSON(500, model.Response{
		// 		Code: 500,
		// 		Msg:  "Token not found or expired:::" + err.Error(),
		// 		Data: nil,
		// 	})
		// 	return
		// }

		// 如果没有错误，则继续查询数据库
		service.DB.Model(&model.UserInfo{}).Where("id=?", id).Find(&user)
		if user.ID == 0 {
			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "User not found",
				Data: nil,
			})
			return
		}

		if !checkMyself(user.TokenTime) {
			c.AbortWithStatusJSON(500, model.Response{
				Code: 500,
				Msg:  "Token expired",
				Data: nil,
			})
			return
		}

		// 如果一切正常，继续处理请求
		c.Next()
	}
}

func Router() {
	f, err := os.Create("gin.log")
	if err != nil {
		fmt.Println("无法创建日志文件:", err)
		return
	}
	defer f.Close()
	gin.DefaultWriter = f

	r := gin.Default()
	r.Use(cors())

	discord := r.Group("/discord")
	{
		discord.Use(CheckHandler())
		mod := new(controller.Discord)

		// discord.GET("/imagine", mod.Imagine)
		discord.POST("/upFile", mod.Upfile)
		discord.POST("/queue", mod.Queue)
		discord.POST("/userImg", mod.GetUserImg)
		discord.GET("/delImg", mod.DelImg)

		// newProduct.POST("/blend", mod.Blend)
		// discord.POST("/describe", mod.Describe)
	}
	ownAccountImagine := r.Group("/ownAccountImagine")
	{
		ownAccountImagine.Use(CheckHandlerMyself())
		mod := new(controller.OwnAccount)

		ownAccountImagine.POST("/queue", mod.Queue)
	}
	user := r.Group("/user")
	{
		mod := new(controller.UserInfo)
		user.POST("/add", mod.Add)
		user.POST("/edit", mod.Edit)
		user.POST("/login", mod.Login)
		user.POST("/autoLogin", mod.AutoLogin)
	}
	gpt := r.Group("/gpt")
	{
		gpt.Use(CheckHandler())
		mod := new(controller.GPT)
		gpt.POST("/send", mod.Run)
		//故事分镜
		gpt.POST("/storyStoryboard", mod.StoryStoryboard)
		//帮写剧情
		gpt.POST("/writePlot", mod.WritePlot)
		//提取角色名
		gpt.POST("/characterExtraction", mod.CharacterExtraction)
		gpt.POST("/characterAndStoryboard", mod.CharacterAndStoryboard)
		gpt.POST("/tuili", mod.Tuili)

	}
	tts := r.Group("/tts")
	{
		mod := new(controller.TTS)
		tts.POST("/send", mod.Send)

	}
	wechat := r.Group("/wechat")
	{
		mod := new(controller.Wechat)
		wechat.POST("/pay", mod.Pay)
		wechat.POST("/get", mod.Get)

	}
	mj := r.Group("/mj")
	{
		mj.Use(CheckHandler())
		mod := new(controller.MjAccount)
		mj.POST("/add", mod.Add)
		mj.POST("/edit", mod.Edit)
		mj.POST("/delete", mod.Delete)
		mj.POST("/get", mod.Get)

		mj.POST("/addReq", mod.AddReq)
		mj.POST("/editReq", mod.EditReq)
		mj.POST("/deleteReq", mod.DeleteReq)
	}
	message := r.Group("/test")
	{
		mod := new(controller.Message)
		mo1 := new(controller.Test)
		message.GET("/get", mod.Get)
		message.GET("/test", mo1.Test)
	}
	SD := r.Group("/sd")
	{
		mod := new(controller.SD)
		SD.POST("/run", mod.Txt2Img)
		SD.GET("/get-models", mod.GetModels)
		SD.GET("/get-lora", mod.GetLora)

	}
	r.Run(":8080")
}
