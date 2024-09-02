package until

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"net/http"
	"strings"
)

// 回调接口
func request(params interface{}, url string) {
	data, err := json.Marshal(params)

	if err != nil {
		fmt.Println("json marshal error: ", err)
		return
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("http request error: ", err)
		return
	}
	defer resp.Body.Close()
}

// 发送消息
func NewRequest(url string, mod interface{}, userToken string) (string, error) {
	var client = new(http.Client)
	jsonData, err := json.Marshal(mod)
	// fmt.Println(string(jsonData))
	if err != nil {
		log.Println(err)
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {

	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", userToken)
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}

	by, _ := ioutil.ReadAll(resp.Body)
	byString := string(by)

	defer resp.Body.Close()

	return byString, nil

}

func UploadFile(path string, userId int) string {
	sss := strings.Replace(path, "format=webp", "format=png", -1)
	fmt.Println(sss)
	res, err := http.Get(sss)
	if err != nil {
		log.Println(err)
		return ""
	}
	now := time.Now().Unix()
	nowStr := strconv.Itoa(int(now))
	by, _ := ioutil.ReadAll(res.Body)
	userIdStr := strconv.Itoa(userId)
	thisPath := "./file/" + nowStr + userIdStr
	ioutil.WriteFile(thisPath+".png", by, 777)
	return thisPath
}
