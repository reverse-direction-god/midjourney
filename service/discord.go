package service

import (
	"encoding/json"
	"fmt"
	"mj/model"
	"mj/until"
)

type discord struct{}

var Discord = discord{}

func (s *discord) Describe(fileName string,
	filePath string,
	ApplicationID string,
	guild_id string,
	channel_id string,
	session_id string,
	userToken string) (string, error) {
	var att model.Attachments
	att.ID = "0"
	att.Filename = fileName
	att.UploadedFilename = filePath
	mod := DescribeRequestModel
	mod.ApplicationID = ApplicationID
	mod.GuildID = guild_id
	mod.ChannelID = channel_id
	mod.SessionID = session_id
	mod.Data.ApplicationCommand.ApplicationID = ApplicationID
	mod.Data.Attachments = append(mod.Data.Attachments, att)
	by, _ := json.Marshal(mod)
	fmt.Println(string(by))
	msg, err := until.NewRequest("https://discord.com/api/v9/interactions", mod, userToken)
	return msg, err
}
func (s *discord) Imagine(prompt string,
	ApplicationID string,
	guild_id string,
	channel_id string,
	session_id string,
	userToken string,
) (string, error) {
	mod := ImagineRequestModel
	mod.Data.Options[0].Value = prompt
	mod.ApplicationID = ApplicationID
	mod.GuildId = guild_id
	mod.ChannelID = channel_id
	mod.SessionID = session_id
	mod.Data.ApplicationCommand.ApplicationID = ApplicationID
	msg, err := until.NewRequest("https://discord.com/api/v9/interactions", mod, userToken)

	return msg, err
}
func Push(queue model.Imagine) {

}
