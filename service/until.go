package service

import (
	"fmt"
	"log"
	"mj/model"
	"strings"
	"time"
)

func (s *discord) UntilImagine(req model.Queue, channel string) {
	prompt := req.Prompt.Content
	var role []string
	// prompt := ""
	for _, j := range req.Prompt.SceneValues {
		switch j.SceneType {
		case "role":
			//name|description|seed|link
			role = strings.Split(j.Text, "|")
			if len(role) >= 2 {
				prompt = prompt + role[1] + " "
			}

		case "action":
			prompt = prompt + " action:" + j.Text

		case "emotion":
			prompt = prompt + " emotion:" + j.Text

		case "scene":
			for _, j := range req.Prompt.MjContent.Prefix {
				prompt = prompt + j + ","
			}
			prompt = prompt + " Screen Description:" + j.Text

		}
	}

	// if req.Prompt.Cref != "" {
	// 	prompt = prompt + " --cref " + req.Prompt.Cref + " --cw 100"
	// }
	prompt = prompt + " " + req.Prompt.Size
	for _, j := range req.Prompt.MjContent.Suffix {
		prompt = prompt + " " + j + " "
	}
	if len(role) >= 4 {
		if role[2] != " " && role[3] != " " {
			prompt = prompt + " " + "--cref " + role[3] + " " + "--cw " + role[2] + " " + req.Prompt.Size + " " + req.Prompt.Model + " " + req.Prompt.Mode
		}
	} else {
		prompt = prompt + " " + req.Prompt.Size + " " + req.Prompt.Model + " " + req.Prompt.Mode
	}

	fmt.Println(prompt)
	var channelID string
	var botToken string
	for _, j := range MjInfoMod {
		for _, jj := range j.RequestInfo {
			if jj.ChannelID == channel {
				channelID = jj.ChannelID
				botToken = j.BotToken
				msg, err := s.Imagine(prompt, jj.ApplicationID, jj.GuildId, jj.ChannelID, jj.SessionID, jj.UserToken)
				//发送失败
				var mod model.Message
				if msg != "" || err != nil {
					mod.RequestMessage = prompt
					mod.ResponseMessage = msg
					mod.UserId = req.UserId
					mod.Type = 1
					mod.Url = ""
					log.Println("msg", msg)
					log.Println("err", err)
					var modError = new(model.Error)
					modError.Message = "mj发送请求:::" + err.Error()
					DB.Model(&model.Error{}).Create(&modError)
					DB.Model(&model.Message{}).Create(&mod)
					//账号剩余+1
					if value, ok := MjNumberMap.Load(botToken); ok {
						val := value.(int)
						val++
						MjNumberMap.Store(botToken, val)
						//释放频道
						ChannelStateMap.Store(channel, 1)
					}

					return
				}
				DisChanUser.Store(channelID, req.UserId)
				fmt.Println("已经设置"+channelID+"为", req.UserId)
				go func() {
					// //定时任务
					var mod model.Message
					timer := time.NewTimer(10 * time.Minute)
					for {
						v, _ := DisChanTime.Load(channelID)
						if ch, ok := v.(bool); ok && ch {
							fmt.Println("执行完成没有超时")
							DisChanTime.Store(channelID, false)
							//频道的并发数+1
							if value, ok := MjNumberMap.Load(botToken); ok {
								val := value.(int)
								val++
								MjNumberMap.Store(botToken, val)
								fmt.Println("账号剩余并发::", val)
								ChannelStateMap.Store(channelID, 1)
								fmt.Println(channelID, "频道没人")
							}
							return
						}

						select {
						//如果到时间了
						case <-timer.C:
							mod.RequestMessage = prompt
							mod.ResponseMessage = "err"
							mod.UserId = req.UserId
							mod.Type = 1
							mod.Url = ""
							fmt.Println("执行超时！！！！")
							DisChanTime.Store(channelID, false)

							DB.Model(&model.Message{}).Create(&mod)
							//频道的并发数+1
							if value, ok := MjNumberMap.Load(botToken); ok {
								val := value.(int)
								val++
								MjNumberMap.Store(botToken, val)
								ChannelStateMap.Store(channelID, 1)
								DisChanUser.Delete(channelID)
							}
							MJListRunMap.Delete(req.UserId)
							return
						//如果执行完成了
						default:

							MJListRunMap.Delete(req.UserId)
						}
					}
				}()
			}

		}
	}
}
