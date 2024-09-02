package controller

import (
	"mj/model"
	"mj/service"

	"github.com/gin-gonic/gin"
)

type MjAccount struct{}

func (s *MjAccount) AddReq(c *gin.Context) {
	var mod model.RequestInfo
	c.ShouldBind(&mod)

	service.DB.Model(&model.RequestInfo{}).Create(&mod)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
func (s *MjAccount) EditReq(c *gin.Context) {
	var mod model.RequestInfo
	c.ShouldBind(&mod)
	service.DB.Model(&model.RequestInfo{}).Where("id=?", mod.ID).Save(&mod)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
func (s *MjAccount) DeleteReq(c *gin.Context) {
	id := c.Query("id")

	service.DB.Model(&model.RequestInfo{}).Where("id=?", id).Delete(nil)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
func (s *MjAccount) Add(c *gin.Context) {
	var mod model.MjAccount
	c.ShouldBind(&mod)
	if mod.Name == "" || mod.BotToken == "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "name or botToken is nil",
			Data: nil,
		})
		return
	}
	service.DB.Model(&model.MjAccount{}).Create(&mod)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
func (s *MjAccount) Edit(c *gin.Context) {
	var mod model.MjAccount
	c.ShouldBind(&mod)
	if mod.Name == "" || mod.BotToken == "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "name or botToken is nil",
			Data: nil,
		})
		return
	}
	service.DB.Model(&model.MjAccount{}).Where("id=?", mod.ID).Save(&mod)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
func (s *MjAccount) Delete(c *gin.Context) {

	id := c.Query("id")
	service.DB.Model(&model.MjAccount{}).Where("id=?", id).Delete(nil)
	service.DB.Model(&model.RequestInfo{}).Where("mj_id=?", id).Delete(nil)
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})

	return
}
func (s *MjAccount) Get(c *gin.Context) {
	var mod []model.MjInfo

	service.DB.Model(&model.MjAccount{}).Find(&mod)

	for i, j := range mod {
		service.DB.Model(&model.RequestInfo{}).Where("mj_id=?", j.ID).Find(&mod[i].RequestInfo)

	}

	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: mod,
	})
	return
}
