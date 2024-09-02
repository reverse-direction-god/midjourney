package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mj/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomStringWithTimestamp(length int) string {
	rand.Seed(time.Now().UnixNano())

	timestamp := time.Now().Unix()
	randomString := generateRandomString(length)
	finalString := fmt.Sprintf("%d%s", timestamp, randomString)

	return finalString
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func httpPost(url string, headers map[string]string, body []byte,
	timeout time.Duration) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	retBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return retBody, err
}
func reqMessage(txt string) (map[string]map[string]interface{}, map[string]string) {
	params := make(map[string]map[string]interface{})
	params["app"] = make(map[string]interface{})
	//填写平台申请的appid
	params["app"]["appid"] = "5137972340"
	//这部分的token不生效，填写下方的默认值就好
	params["app"]["token"] = "* Access *=="
	//填写平台上显示的集群名称
	params["app"]["cluster"] = "volcano_tts"
	params["user"] = make(map[string]interface{})
	//这部分如有需要，可以传递用户真实的ID，方便问题定位
	params["user"]["uid"] = "1"
	params["audio"] = make(map[string]interface{})
	//填写选中的音色代号
	params["audio"]["voice_type"] = "BV700_streaming"
	params["audio"]["encoding"] = "wav"
	params["audio"]["speed_ratio"] = 1.0
	params["audio"]["volume_ratio"] = 1.0
	params["audio"]["pitch_ratio"] = 1.0
	params["request"] = make(map[string]interface{})
	params["request"]["reqid"] = generateRandomStringWithTimestamp(20)
	params["request"]["text"] = txt
	params["request"]["text_type"] = "plain"
	params["request"]["operation"] = "query"

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	//bearerToken为saas平台对应的接入认证中的Token
	headers["Authorization"] = "Bearer;" + "*-Mhqq"
	return params, headers
}
func codeMessage(code int) model.Response {
	mod := model.Response{}
	mod.Code = 500

	mod.Data = nil

	switch code {
	case 3000:
		mod.Code = 200
		mod.Data = nil
		mod.Msg = ""

		return mod

	case 3001:

		mod.Msg = "一些参数的值非法，比如operation/workflow配置错误"

		return mod
	case 3003:

		mod.Msg = "超过在线设置的并发阈值"

		return mod
	case 3005:

		mod.Msg = "后端服务器负载高"

		return mod
	case 3006:

		mod.Msg = "请求已完成/失败之后，相同reqid再次请求"

		return mod
	case 3010:

		mod.Msg = "单次请求超过设置的文本长度阈值"

		return mod
	case 3011:

		mod.Msg = "参数有误或者文本为空、文本与语种不匹配、文本只含标点"

		return mod
	case 3030:

		mod.Msg = "单次请求超过服务最长时间限制"

		return mod
	case 3031:

		mod.Msg = "后端出现异常"

		return mod
	case 3032:

		mod.Msg = "后端网络异常"

		return mod

	}
	return mod
}

type TTS struct{}

func (s *TTS) Send(c *gin.Context) {
	txt := c.Query("txt")
	var respJSON model.TTSServResponse
	params, headers := reqMessage(txt)
	url := "https://openspeech.bytedance.com/api/v1/tts"
	timeo := 30 * time.Second
	bodyStr, _ := json.Marshal(params)
	synResp, err := httpPost(url, headers,
		[]byte(bodyStr), timeo)
	if err != nil {
		log.Println(err)
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}
	json.Unmarshal(synResp, &respJSON)
	if respJSON.Code == 3000 {
		c.JSON(200, model.Response{
			Code: 200,
			Msg:  "ok",
			Data: respJSON.Data,
		})
		return
	}
	c.JSON(200, codeMessage(respJSON.Code))
	return
}
