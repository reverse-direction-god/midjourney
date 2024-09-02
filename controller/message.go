package controller

import (
	"mj/model"
	"mj/service"
	"strings"

	"github.com/gin-gonic/gin"
)

type Message struct{}

func (s *Message) Get(c *gin.Context) {
	var mod []model.Message
	var res []struct {
		Url             []string `json:"url"`
		ResponseMessage string   `json:"response_message"`
	}
	service.DB.Model(&model.Message{}).Find(&mod)
	for _, j := range mod {
		var resMod struct {
			Url             []string `json:"url"`
			ResponseMessage string   `json:"response_message"`
		}
		if j.Url != "" {

			resMod.ResponseMessage = j.ResponseMessage
			urls := strings.Split(j.Url, "|")
			for i := 0; i < len(urls)-1; i++ {
				resMod.Url = append(resMod.Url, urls[i])
			}
		}
		res = append(res, resMod)
	}
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: res,
	})
}
