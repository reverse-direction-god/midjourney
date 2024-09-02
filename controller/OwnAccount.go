package controller

import (
	"fmt"
	"mj/model"
	"mj/service"
	"mj/until"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/cloudflare/cfssl/log"
	"github.com/gin-gonic/gin"
)

type OwnAccount struct{}

var ChannelFindUser sync.Map
var channelFindUserMap = make(map[string]int)
var channelTimeOut sync.Map
var accountLife sync.Map
var queue = make(chan model.OwnAccountReq, 100)

func 格式化请求参数(req model.Queue) string {
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
	return prompt
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	//过滤自己的
	if message.Author.ID == session.State.User.ID {
		return
	}
	v, _ := ChannelFindUser.Load(message.ChannelID)

	var userId int
	if v != nil {
		userId = v.(int)
	} else {
		return
	}

	if message.MessageReference == nil {
		if len(message.Attachments) != 0 {

			go func() {
				var mod model.Message

				for _, j := range message.Attachments {
					mod.UserId = userId
					allImg := until.UploadFile(j.URL, mod.UserId)
					url, err := until.ImgToFour(allImg, mod.UserId) //消息url

					mod.Url = url
					mod.ResponseMessage = message.Content //消息体
					mod.Type = 1
					if err != nil {
						service.DB.Model(&model.Error{}).Create(map[string]string{
							"message": "mj分割图片:::" + err.Error(),
						})
					}
					service.DB.Model(&model.Message{}).Create(&mod)
				}
			}()
		}
	}
}
func message(bot string) error {
	// 创建一个 Discord 会话

	discord, err := discordgo.New("Bot " + bot)
	if err != nil {
		log.Error("discord New err: " + err.Error())
		return err
	}
	discord.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentGuildMessageReactions | discordgo.IntentDirectMessageReactions
	// 添加消息处理函数
	discord.AddHandler(messageCreate)

	// 连接到 Discord
	err = discord.Open()
	if err != nil {
		log.Error("Error opening connection to Discord:" + err.Error())
		return err
	}

	// 等待程序终止
	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	// 关闭 Discord 会话
	discord.Close()
	return nil
}

func (s *OwnAccount) Queue(c *gin.Context) {

	var req model.OwnAccountReq
	c.ShouldBind(&req)
	token := c.GetHeader("token")
	id := until.Decrypt(token)
	fmt.Println(1)

	if id == "" {
		log.Error("token错误")
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  "token err",
			Data: nil,
		})
		return
	}

	prompt := 格式化请求参数(req.Info)

	_, ok := accountLife.Load(req.Channel_id)

	msg, err := service.Discord.Imagine(prompt, req.Application_id, req.Guild_id, req.Channel_id, req.Session_id, req.User_token)

	if msg != "" || err != nil {
		log.Error(msg)
		c.JSON(500, model.Response{
			Code: 500,
			Msg:  err.Error(),
			Data: nil,
		})
		return
	}

	if !ok {
		userid, _ := strconv.Atoi(id)
		fmt.Println("userId 第二次是", userid)
		ChannelFindUser.Store(req.Channel_id, userid)
		accountLife.Store(req.Channel_id, id)
		channelFindUserMap[req.Channel_id] = userid
		go message(req.BotToken)

	}

	c.JSON(200, model.Response{
		Code: 200,
		Msg:  "ok",
		Data: nil,
	})
	return
}
