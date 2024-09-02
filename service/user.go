package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mj/model"
	"mj/until"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudflare/cfssl/log"
)

type userInfo struct{}

var UserInfo = userInfo{}
var (
	appid  = "*"
	secret = "*"
)

func loginWechat(code string) string {
	url := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + appid + "&secret=" + secret + "&code=" + code + "&grant_type=authorization_code"
	var mod struct {
		AccessToken string `json:"access_token"`
		OpenId      string `json:"openid"`
	}
	client := http.Client{}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return ""
	}

	by, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(by, &mod)
	fmt.Println(string(by))
	return mod.OpenId

}
func delRedisToken(id uint) {
	keys, err := Redis.Keys(Redis.Context(), "*").Result()
	if err != nil {
		log.Error(err)
		return
	}
	strId := strconv.Itoa(int(id))
	for _, key := range keys {
		// 获取键的值
		value, err := Redis.Get(Redis.Context(), key).Result()
		if err != nil {
			log.Error(err)
			return
		}

		// 如果值为 "1"，则删除键
		if value == strId {
			err := Redis.Del(Redis.Context(), key).Err()
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
}
func (s *userInfo) Login(invitation_code string, code string) (string, model.UserInfo) {

	var mod model.UserInfo

	var userCode model.UserCodes

	openid := loginWechat(code)
	if openid == "" {
		return "err", mod
	}
	tx := DB.Model(&model.UserInfo{}).Where("openid=?", openid)
	tx.Find(&mod)

	if mod.Openid == "" {
		mod.InvitationCode = strconv.Itoa(int(time.Now().Unix()))
		mod.Openid = openid
		tx.Create(&mod)
	}
	id := mod.ID
	if invitation_code != "" {
		//搜索这个用户是否使用过邀请码了
		//如果使用过 就跳过
		//如果没有使用过 就给用户加入优惠券 并且为邀请码主返现
		tx := DB.Model(&model.UserCodes{}).Where("user_id=?", id)
		tx.Find(&userCode)
		if userCode.ID == 0 {
			var user model.UserInfo

			DB.Model(&model.UserInfo{}).Where("invitation_code=?", invitation_code).Find(&user) //寻找此邀请码用户
			userCode.UserId = int(id)
			userCode.UsedUserId = int(user.ID)
			tx.Create(&userCode) //标记为使用过

		}
	}
	var token string
	idstr := strconv.Itoa(int(id))
	keys := Redis.Keys(context.Background(), "*").Val()
	// 遍历键
	for _, key := range keys {
		// 获取键的值
		value := Redis.Get(context.Background(), key).Val()
		if value == idstr {
			err := Redis.Del(context.Background(), key).Err()
			if err != nil {
				log.Error(err)
				return "", mod
			}
		}
	}
	token = until.Encryption(id)
	Redis.Set(Redis.Context(), token, id, 0)
	return token, mod
}
