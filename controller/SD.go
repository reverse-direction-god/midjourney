package controller

import (
	"encoding/json"
	"io/ioutil"
	"mj/model"
	"mj/service"
	"mj/until"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SD struct{}

func (s *SD) Txt2Img(c *gin.Context) {

	var mod model.SDPromptConfig
	token := c.GetHeader("token")
	id := until.Decrypt(token)
	if id == "" {
		c.JSON(500, &model.Response{
			500,
			"id==nil",
			nil,
		})
		return
	}

	c.ShouldBind(&mod)
	mod.UserId = id
	// var ser service.RocketMq
	// err = ser.PushSD(mod)
	// if err != nil {
	// 	c.JSON(500, &model.Response{
	// 		500,
	// 		err.Error(),
	// 		nil,
	// 	})
	// 	return
	// }
	// _, ok := service.SDListRunMap.Load(id)
	// if ok {
	// 	c.JSON(500, &model.Response{
	// 		500,
	// 		"用户已在队列中",
	// 		nil,
	// 	})
	// 	return
	// } else {
	// 	service.SDListRunMap.Store(id, 1)
	// }
	service.SDMQ <- mod
	c.JSON(200, &model.Response{
		200,
		"ok",
		nil,
	})
	return
}

func (s *SD) GetLora(c *gin.Context) {
	resp, err := http.Get("*")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "get req " + err.Error(),
			"data": nil,
		})
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 400,
			"msg":  "read resp " + err.Error(),
			"data": nil,
		})
		return
	}
	type Lora struct {
		Name  string `json:"name"`
		Alias string `json:"alias"`
		Path  string `json:"path"`
	}
	var mod []Lora
	err = json.Unmarshal(buf, &mod)
	if err != nil {
		c.JSON(500, gin.H{
			"code": 400,
			"msg":  "Unmarshal" + err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": mod,
	})
	return

}

func (s *SD) GetModels(c *gin.Context) {
	resp, err := http.Get("*")
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "get req " + err.Error(),
			"data": nil,
		})
		return
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "read resp " + err.Error(),
			"data": nil,
		})
		return
	}

	type Model struct {
		Title     string      `json:"title"`
		ModelName string      `json:"model_name"`
		Hash      string      `json:"hash"`
		Sha       string      `json:"sha256"`
		Filename  string      `json:"filename"`
		Config    interface{} `json:"config"`
	}

	var models []Model

	err = json.Unmarshal(buf, &models)
	if err != nil {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "Unmarshal" + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"data": models,
	})
}

// func (s *SD) Repair(c *gin.Context) {
// 	token := c.GetHeader("token")
// 	id, _ := service.Redis.Get(c, token).Result()

// 	var image = make(map[string]string)
// 	c.ShouldBind(&image)
// 	type Image struct {
// 		Image  string `json:"image"`
// 		Title  string `json:"title"`
// 		UserId string `json:"userId"`
// 	}

// }
