package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mj/model"
	"net/http"
)

type gpt struct{}

var GPT = gpt{}

func (s *gpt) GPTPOST(mod interface{}) string {
	var response model.GPTResponse
	client := &http.Client{}
	by, _ := json.Marshal(&mod)
	req, _ := http.NewRequest("POST", "https://hk.xty.app/v1/chat/completions", bytes.NewBuffer(by))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+"*")
	res, _ := client.Do(req)
	by, _ = ioutil.ReadAll(res.Body)
	json.Unmarshal(by, &response)
	if len(response.Choices) == 0 {
		log.Println("GPT response is nil")
		return ""
	}
	return response.Choices[0].Message.Content
}
