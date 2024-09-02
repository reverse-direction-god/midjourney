package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mj/model"
	"mj/until"

	"net/http"
	"strconv"
	"time"

	"github.com/cloudflare/cfssl/log"
)

type SD struct{}

const webuiServerURL = "*"

func callAPI(payload model.SDPromptConfig) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request, err := http.NewRequest("POST", webuiServerURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return result, err
}

func (s *SD) Txt2Img(req model.SDPromptConfig, userId string) (string, error) {
	var mod model.SDTxt2ImgResponse
	var resString string
	mp, err := callAPI(req)
	if err != nil {
		var modError = new(model.Error)
		modError.Message = "sd 发送请求错误:::" + err.Error()
		DB.Model(&model.Error{}).Create(&modError)
		return "", err
	}
	by, err := json.Marshal(mp)
	if err != nil {
		log.Error(err)
		var modError = new(model.Error)
		modError.Message = "sd 发送时格式解析:::" + err.Error()
		DB.Model(&model.Error{}).Create(&modError)
		return "", errors.New("error")
	}
	json.Unmarshal(by, &mod)
	if len(mod.Images) == 0 {

		return "", errors.New("error")
	}
	for _, j := range mod.Images {
		nowTime := time.Now().Unix()
		nowStr := until.RandString(5) + userId + strconv.Itoa(int(nowTime)) + ".png"
		pngs, _ := base64.StdEncoding.DecodeString(j)
		err := until.OssEd(pngs, nowStr)
		if err != nil {
			var modError = new(model.Error)

			modError.Message = "oss 上传图片错误:::" + err.Error()
			DB.Model(&model.Error{}).Create(&modError)
		}

		resString += nowStr + "|"
	}

	return resString, nil
}
func (s *SD) Repair(image string, userId string, w int, h int) string {
	var requestInfo = make(map[string]interface{})
	jsonByte, err := ioutil.ReadFile("./config/config/repair.json")
	if err != nil {
		return ""
	}
	err = json.Unmarshal(jsonByte, &requestInfo)
	if err != nil {
		return ""
	}
	requestInfo["upscaling_resize_w"] = w
	requestInfo["upscaling_resize_h"] = h
	requestInfo["image"] = image
	reqByte, err := json.Marshal(requestInfo)
	if err != nil {
		return ""
	}
	resp, err := http.Post("*", "application/json", bytes.NewReader(reqByte))

	if err != nil {
		return ""
	}
	respBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var m map[string]interface{}
	err = json.Unmarshal(respBuf, &m)
	if err != nil {
		return ""
	}
	if m == nil || m["image"] == nil {
		return ""
	}
	resImageBase64 := m["image"].(string)
	resBin, err := base64.StdEncoding.DecodeString(resImageBase64)
	if err != nil {
		return ""
	}
	nowTime := time.Now().Unix()
	nowStr := until.RandString(5) + userId + strconv.Itoa(int(nowTime)) + ".png"
	until.OssEd(resBin, nowStr)
	return nowStr
}

func modelsOptinons(modelName string) error {
	s := `{"sd_model_checkpoint":"` + modelName + `"}`
	respByte, err := until.Post("*", []byte(s))
	if err != nil {
		var mod = new(model.Error)
		mod.Message = "sd切换模型失败:::" + err.Error()
		DB.Model(&model.Error{}).Create(&mod)
		return err
	}
	if string(respByte) == "null" {
		return nil
	} else {
		var mod = new(model.Error)
		mod.Message = "sd切换模型返回内容:::" + string(respByte)
		DB.Model(&model.Error{}).Create(&mod)
		return errors.New("sd切换模型失败")
	}
}
func init() {
	go send()
}

func send() {
	var sd SD
	for {
		select {
		case mod := <-SDMQ:
			switch mod.Type {
			case 0:
				var sqlCreate model.Message
				mod.Type = 0
				mod.OldFileName = ""
				//换模型
				fmt.Println("换模型:::" + mod.OverrideSettings.Sd_model_checkpoint)
				err := modelsOptinons(mod.OverrideSettings.Sd_model_checkpoint)
				if err != nil {
					log.Error(err)
					log.Error(err)
					sqlCreate.ResponseMessage = "err"
					sqlCreate.Url = ""
					sqlCreate.RequestMessage = mod.Prompt
					DB.Model(&model.Message{}).Create(&sqlCreate)
					continue
				}
				time.Sleep(15 * time.Second)
				mod.OverrideSettings.Sd_model_checkpoint = ""
				stringss, err := sd.Txt2Img(mod, mod.UserId)
				sqlCreate.ID = 0
				sqlCreate.UserId, _ = strconv.Atoi(mod.UserId)
				if err != nil {
					log.Error(err)
					sqlCreate.ResponseMessage = "err"
					sqlCreate.Url = ""
					sqlCreate.RequestMessage = mod.Prompt
				} else {
					sqlCreate.ResponseMessage = mod.Prompt
					sqlCreate.Url = stringss
					sqlCreate.RequestMessage = ""
				}
				SDListRunMap.Delete(mod.UserId)
				DB.Model(&model.Message{}).Create(&sqlCreate)
			case 2:
				var sqlCreate model.Message
				fileName := sd.Repair(mod.ImageBytes, mod.OldFileName, mod.Width, mod.Height)
				if fileName == "" {
					sqlCreate.ID = 0
					sqlCreate.UserId, _ = strconv.Atoi(mod.UserId)
					sqlCreate.ResponseMessage = "err"
					sqlCreate.Url = ""
					sqlCreate.RequestMessage = mod.OldFileName
				} else {
					sqlCreate.ID = 0
					sqlCreate.UserId, _ = strconv.Atoi(mod.UserId)
					sqlCreate.ResponseMessage = mod.OldFileName
					sqlCreate.Url = fileName + "|"
					sqlCreate.RequestMessage = ""
				}
				SDListRunMap.Delete(mod.UserId)
				DB.Model(&model.Message{}).Create(&sqlCreate)

			}

		}
	}
}
