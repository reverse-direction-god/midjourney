package controller

import (
	"context"
	"fmt"
	"mj/model"
	"mj/service"
	"mj/until"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserInfo struct{}

func (s *UserInfo) Add(c *gin.Context) {
	var mod model.UserInfo
	c.ShouldBind(&mod)
	tx := service.DB.Model(&model.UserInfo{}).Create(&mod)
	if tx.Error != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  tx.Error.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(500, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}
func (s *UserInfo) Edit(c *gin.Context) {
	var mod model.UserInfo
	c.ShouldBind(&mod)
	token := c.GetHeader("token")

	id, err := service.Redis.Get(context.Background(), token).Uint64()
	if err != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	if id == 0 {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "token error",
			Data: nil,
		})
		return
	}
	mod.ID = uint(id)
	fmt.Println(mod.ID)
	tx := service.DB.Model(&model.UserInfo{}).Where("id=?", mod.ID).Save(&mod)
	if tx.Error != nil {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  tx.Error.Error(),
			Data: nil,
		})
		return
	}
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
}

func (s *UserInfo) Login(c *gin.Context) {
	code := c.Query("code")
	invitation_code := c.Query("invitationCode")
	var res struct {
		UserInfo model.UserInfo `json:"userInfo"`
		Token    string         `json:"token"`
	}
	ser := service.UserInfo
	fmt.Println(invitation_code)
	res.Token, res.UserInfo = ser.Login(invitation_code, code)
	if res.Token == "" {
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "no",
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

func (s *UserInfo) AutoLogin(c *gin.Context) {
	var id string
	var res struct {
		UserInfo model.UserInfo `json:"userInfo"`
		Token    string         `json:"token"`
	}
	token := c.GetHeader("token")
	if id = until.Decrypt(token); id == "" {
		c.JSON(200, model.Response{
			Code: 500,
			Msg:  "token err",
			Data: res,
		})
	}

	service.DB.Model(&model.UserInfo{}).Where("id=?", id).Find(&res.UserInfo)
	uid, _ := strconv.Atoi(id)
	res.Token = until.Encryption(uint(uid))
	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: res,
	})

}

// func (s *UserInfo) Login(c *gin.Context) {
// 	var req model.UserInfo
// 	var mod model.UserInfo
// 	c.ShouldBind(&req)
// 	service.DB.Model(&model.UserInfo{}).Where("user_name=?", req.UserName).Find(&mod)

// 	if mod.ID == 0 {
// 		c.JSON(500, model.Response{
// 			Code: 500,
// 			Msg:  "userName err",
// 			Data: nil,
// 		})
// 		return
// 	}
// 	if req.Password != mod.Password {
// 		c.JSON(500, model.Response{
// 			Code: 500,
// 			Msg:  "password err",
// 			Data: nil,
// 		})
// 		return
// 	}
// 	mod.Password = ""
// 	token := until.Encryption(mod.ID)

// 	service.Redis.Set(context.Background(), token, mod.ID, 0)
// 	abs := struct {
// 		UserInfo model.UserInfo `json:"user_info"`
// 		Token    string         `json:"token"`
// 	}{}
// 	abs.UserInfo = mod
// 	abs.Token = token
// 	c.JSON(200, model.Response{
// 		Code: 200,
// 		Msg:  "ok",
// 		Data: abs,
// 	})
// }

// func (s *UserInfo) Delete(c *gin.Context) {
// 	var mod model.UserInfo
// 	c.ShouldBind(&mod)
// 	tx := service.DB.Model(&model.UserInfo{}).Where("id=?", mod.ID).Delete(nil)
// 	if tx.Error != nil {
// 		c.JSON(500, model.Response{
// 			Code: 500,
// 			Msg:  tx.Error.Error(),
// 			Data: nil,
// 		})
// 		return
// 	}
// 	c.JSON(500, model.Response{
// 		Code: 200,
// 		Msg:  "ok",
// 		Data: nil,
// 	})
// }
