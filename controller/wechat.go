package controller

import (
	"mj/model"
	"mj/service"

	"github.com/gin-gonic/gin"
)

type Wechat struct{}

func (s *Wechat) Get(c *gin.Context) {
	out_trade_no := c.Query("out_trade_no")
	token := c.GetHeader("token")
	var mod model.PayInfo
	id, err := service.Redis.Get(c, token).Result()
	if err != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "token err",
			Data: nil,
		})
		return
	}
	service.DB.Model(&model.PayInfo{}).Where("user_id=? AND out_trade_no=?", id, out_trade_no).Find(&mod)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: mod,
	})
}
func (s *Wechat) Pay(c *gin.Context) {
	var mod struct {
		Number      int64  `json:"number"`
		Description string `json:"description"`
	}
	c.ShouldBind(&mod)
	token := c.GetHeader("token")
	if token == "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "token err",
			Data: nil,
		})
		return
	}
	id, err := service.Redis.Get(c, token).Result()
	if err != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "token err",
			Data: nil,
		})
		return
	}
	var res struct {
		Url  string `json:"url"`
		Note string `json:"note"`
	}
	res.Url, res.Note = service.Wechat.Pay(id, mod.Description, mod.Number)
	if res.Url == "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}

	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: res,
	})
	return
}
