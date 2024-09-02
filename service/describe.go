package service

import (
	"fmt"
	"mj/model"
	"mj/until"
)

func (s *discord) UntilDescribe(req model.Queue, channel string) (string, error) {
	for _, j := range MjInfoMod {
		for _, jj := range j.RequestInfo {
			if jj.ChannelID == channel {
				msg, err := s.Describe(req.FileName, req.UploadFileName, jj.ApplicationID, jj.GuildId, jj.ChannelID, jj.SessionID, jj.UserToken)
				//发送失败
				var mod model.Message
				if msg != "" || err != nil {
					mod.RequestMessage = req.UploadFileName
					fmt.Println(msg)
					mod.ResponseMessage = "msg1:::" + msg
					mod.UserId = req.UserId
					mod.Type = 2
					mod.Url = ""
					DB.Model(&model.Message{}).Create(&mod)
				}
				//设置这个用户已经占用了这个频道
				ChannelStateMap.Store(jj.ChannelID, req.UserId)

				botToken := ""
				for _, j := range MjInfoMod {
					for _, jj := range j.RequestInfo {
						if jj.ChannelID == channel {
							botToken = j.BotToken
						}
					}
				}
				//减少这个频道的并发数
				if value, ok := MjNumberMap.Load(botToken); ok {
					val := value.(int)
					val = val - 4
					MjNumberMap.Store(botToken, val)
					DiscribeNumber.Store(channel, 4)
				}

				//再Mysql创建一个等待消息的记录
				// mod.RequestMessage = req.Prompt
				// mod.ResponseMessage = ""
				// mod.UserId = req.UserId
				// mod.Type = 1
				// mod.Url = ""
				// DB.Model(&model.Message{}).Create(&mod)
			}
		}
	}

	// mod := DescribeRequestModel
	// att := model.Attachments{}
	// att.ID = "0"
	// att.Filename = fileName
	// att.UploadedFilename = uploadFileName
	// mod.Data.Attachments = append(mod.Data.Attachments, att)
	// for _, j := range MjInfoMod {
	// 	for _, jj := range j.RequestInfo {
	// 		if jj.ChannelID == channel {
	// 			mod.ApplicationID = jj.ApplicationID
	// 			mod.GuildID = jj.GuildId
	// 			mod.ChannelID = jj.ChannelID
	// 			mod.SessionID = jj.SessionID
	// 		}
	// 	}
	// }

	return until.NewRequest("https://discord.com/api/v9/interactions", nil, "*.*.*-*")

}
