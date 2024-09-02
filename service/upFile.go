package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mj/model"
	"mj/until"
	"net/http"

	"github.com/gin-gonic/gin"
)

type reqUploadFile struct {
	ImgData []byte `json:"imgData"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
}

func attachments(name string, size int64) (model.ResAttachments, error) {
	var acc model.MjAccount
	var requestInfo model.RequestInfo
	DB.Model(model.MjAccount{}).Find(&acc)
	DB.Model(model.RequestInfo{}).Where("mj_id=?", acc.ID).Find(&requestInfo)

	requestBody := model.ReqAttachments{
		Files: []model.ReqFile{{
			Filename: name,
			FileSize: size,
			Id:       "1",
		}},
	}

	uploadurls := fmt.Sprintf("https://discord.com/api/v9/channels/%s/attachments", "1214874136522784781")
	body, err := until.NewRequest(uploadurls, requestBody, "*")
	var data model.ResAttachments
	json.Unmarshal([]byte(body), &data)
	return data, err
}

func (s *discord) UpFile(c *gin.Context) (string, error) {
	var body reqUploadFile
	if err := c.ShouldBindJSON(&body); err != nil {
		log.Println(err)
		return "", err
	}

	data, err := attachments(body.Name, body.Size)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(data.Attachments) == 0 {
		log.Println(err)
		return "", err
	}
	payload := bytes.NewReader(body.ImgData)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", data.Attachments[0].UploadUrl, payload)

	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "image/png")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer res.Body.Close()
	by, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(by))
	return data.Attachments[0].UploadFilename, nil
}
