package controller

import (
	"fmt"
	"mj/model"
	"mj/service"
	"mj/until"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Discord struct{}

var ser = service.Discord
var Queue = make(chan model.Queue, 500000)

func init() {
	go 从队列中提取数据并执行()
}

func 从队列中提取数据并执行() {
	for {
		time.Sleep(time.Second * 2)
	abs:
		channel := ""
		MjNumber := 0
		service.MjNumberMap.Range(func(key, value interface{}) bool {
			//如果这个账号有空闲
			if value.(int) > 0 {
				MjNumber = 1
				service.ChannelStateMap.Range(func(channel_key, channel_value interface{}) bool {
					//如果这个频道没人用
					if channel_value.(int) == 1 {
						channel = channel_key.(string)

						return false
					}
					return true
				})
				return false
			}
			return true // 继续遍历
		})

		//channel没有用的   请求队列空的      账号可用并发数为0
		if channel == "" || len(Queue) == 0 || MjNumber == 0 {
			goto abs
		}

		//找到了一个可以用的channel
		var req model.Queue
		//如果chan中有数据
		if len(Queue) != 0 {

			req = <-Queue
			if req.Type == "imagine" {

				//设置这个用户已经占用了这个频道
				service.ChannelStateMap.Store(channel, req.UserId)

				//减少这个频道的并发数
				botToken := ""
				for _, j := range service.MjInfoMod {
					for _, jj := range j.RequestInfo {
						if jj.ChannelID == channel {
							botToken = j.BotToken
						}
					}
				}
				if value, ok := service.MjNumberMap.Load(botToken); ok {
					val := value.(int)
					val--
					service.MjNumberMap.Store(botToken, val)

				}
				go ser.UntilImagine(req, channel)
			}
			//图生图
			if req.Type == "describe" {
				go ser.UntilDescribe(req, channel)
			}
		}

	}

}

// 排队
func (s *Discord) Queue(c *gin.Context) {
	var req model.Queue
	c.ShouldBind(&req)
	id, _ := c.Get("user")

	idInt, _ := strconv.Atoi(id.(string))
	req.UserId = idInt
	// var ss service.RocketMq
	// ss.PushMJ(req)
	// _, ok := service.MJListRunMap.Load(idInt)
	// if ok {
	// 	c.JSON(500, &model.Response{
	// 		500,
	// 		"用户已在队列中",
	// 		nil,
	// 	})
	// 	return
	// } else {
	// 	service.MJListRunMap.Store(idInt, 1)
	// }
	Queue <- req
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}

//	func (s *Discord) Describe(c *gin.Context) {
//		mod := struct {
//			FileName       string `json:"fileName"`
//			UploadFileName string `json:"uploadFileName"`
//		}{}
//		c.ShouldBind(&mod)
//		msg, err := ser.Describe(mod.FileName, mod.UploadFileName)
//		if err != nil {
//			c.JSON(500, model.Response{
//				Code: 500,
//				Msg:  "no",
//				Data: nil,
//			})
//			return
//		}
//		if msg != "" {
//			c.JSON(500, model.Response{
//				Code: 500,
//				Msg:  msg,
//				Data: nil,
//			})
//			return
//		}
//		c.JSON(200, model.Response{
//			Code: 200,
//			Msg:  "ok",
//			Data: nil,
//		})
//	}
func (s *Discord) Blend(c *gin.Context) {

	var mod []model.Attachments
	c.ShouldBind(&mod)
	msg, err := ser.Blend(mod)
	if err != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "no",
			Data: nil,
		})
		return
	}
	if msg != "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  msg,
			Data: nil,
		})
		return
	}
	// sqlMessage := model.Message{
	// 	RequestMessage:  "",
	// 	ResponseMessage: "",
	// 	Url:             "",
	// 	Type:            2,
	// }
	// service.DB.Model(&model.Message{}).Create(&sqlMessage)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})

}

// func (s *Discord) Imagine(c *gin.Context) {

// 	prompt := c.Query("prompt")
// 	if prompt == "" {
// 		c.JSON(500, model.Response{
// 			Data: 500,
// 			Msg:  "prompt nil",
// 		})
// 	}

// 	msg, err := ser.Imagine(prompt)

// 	if err != nil {
// 		c.JSON(500, model.Response{
// 			Code: 500,
// 			Msg:  err.Error(),
// 			Data: nil,
// 		})
// 	}
// 	if msg != "" {
// 		c.JSON(500, model.Response{
// 			Code: 500,
// 			Msg:  msg,
// 			Data: nil,
// 		})
// 	}
// 	mod := model.Message{

//			RequestMessage:  prompt,
//			ResponseMessage: "",
//			Url:             "",
//			Type:            1,
//		}
//		service.DB.Model(&model.Message{}).Create(&mod)
//		c.JSON(200, model.Response{
//			Code: 200,
//			Msg:  "ok",
//			Data: nil,
//		})
//	}
func (s *Discord) Upfile(c *gin.Context) {

	UploadFilename, err := ser.UpFile(c)
	if err != nil {
		c.JSON(500, &model.Response{
			Code: 500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(200, &model.Response{
		Code: 200,
		Msg:  "ok",
		Data: UploadFilename,
	})
	return
}

func (s *Discord) GetUserImg(c *gin.Context) {
	queue := 1
	var res []struct {
		RequestMessage  string   `json:"request_message"`
		ResponseMessage string   `json:"response_message"`
		Url             []string `json:"url"`
		Type            int      `json:"type"`
		UserId          int      `json:"user_id"`
		ID              uint     `gorm:"primarykey"`
		CreatedAt       time.Time
		UpdatedAt       time.Time
		DeletedAt       gorm.DeletedAt `gorm:"index"`
	}
	var mod []model.Message
	id, _ := c.Get("user")

	service.DB.Model(&model.Message{}).Where("user_id=?", id).Find(&mod)
	for _, j := range mod {
		var modTim struct {
			RequestMessage  string   `json:"request_message"`
			ResponseMessage string   `json:"response_message"`
			Url             []string `json:"url"`
			Type            int      `json:"type"`
			UserId          int      `json:"user_id"`
			ID              uint     `gorm:"primarykey"`
			CreatedAt       time.Time
			UpdatedAt       time.Time
			DeletedAt       gorm.DeletedAt `gorm:"index"`
		}
		modTim.ID = j.ID
		modTim.CreatedAt = j.CreatedAt
		modTim.DeletedAt = j.DeletedAt
		modTim.RequestMessage = j.RequestMessage
		modTim.ResponseMessage = j.ResponseMessage
		modTim.Type = j.Type
		modTim.UpdatedAt = j.UpdatedAt
		if j.Url != "" {
			urls := strings.Split(j.Url, "|")
			for i := 0; i < len(urls)-1; i++ {
				modTim.Url = append(modTim.Url, "*"+urls[i])
			}
		}
		modTim.UserId = j.UserId
		res = append(res, modTim)
	}
	idInt, _ := strconv.Atoi(id.(string))
	service.DisChanUser.Range(func(key, value interface{}) bool {
		if value != nil && value.(int) == idInt {
			queue = 0    // 如果找到值设置queue为0
			return false // 停止遍历
		}
		fmt.Println("正在运行的是", value.(int))
		return true // 继续遍历
	})
	if len(res) == 0 {
		// c.JSON(200, &model.Response{
		// 	Code: 200,
		// 	Msg:  "ok",
		// 	Data: []interface{}{},
		// })
		c.JSON(200, gin.H{
			"code":  200,
			"msg":   "ok",
			"data":  []interface{}{},
			"queue": queue,
		})

		return
	}

	c.JSON(200, gin.H{
		"code":  200,
		"msg":   "ok",
		"data":  res,
		"queue": queue,
	})
}
func (s *Discord) DelImg(c *gin.Context) {
	imgId := c.Query("id")
	token := c.GetHeader("token")
	id := until.Decrypt(token)
	if id == "0" || id == "" {
		c.JSON(500, &model.Response{
			Code: 500,
			Msg:  "token err",
			Data: nil,
		})
		return
	}
	service.DB.Model(&model.Message{}).Where("id=?", imgId).Delete(nil)
	c.JSON(200, &model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}
