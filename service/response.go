package service

import (
	"encoding/json"
	"fmt"
	"mj/model"
	"mj/until"
	"os"
	"os/signal"
	"regexp"

	"syscall"

	"github.com/bwmarrin/discordgo"
)

var isn = 0

type reqCb struct {
	Embeds  []*discordgo.MessageEmbed `json:"embeds,omitempty"`
	Discord *discordgo.MessageCreate  `json:"discord,omitempty"`
	Content string                    `json:"content,omitempty"`
	Type    string                    `json:"type"`
}

var isRequest = make(map[string]bool)

func message(bot string) {
	// 创建一个 Discord 会话

	discord, err := discordgo.New("Bot " + bot)
	if err != nil {
		panic("discord New err: " + err.Error())
	}

	discord.Identify.Intents = discordgo.IntentGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentGuildMessageReactions | discordgo.IntentDirectMessageReactions
	// 添加消息处理函数
	discord.AddHandler(messageCreate)

	// 连接到 Discord
	err = discord.Open()
	if err != nil {

		panic("Error opening connection to Discord:" + err.Error())
	}

	// 等待程序终止
	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt
	// 关闭 Discord 会话
	discord.Close()
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	by, err := json.Marshal(message)
	if err != nil {
		fmt.Println("解析失败")
	} else {
		fmt.Println(string(by))
	}

	re := regexp.MustCompile(`\*\*(.*?)\*\*`)
	// 使用正则表达式查找匹配项
	match := re.FindStringSubmatch(message.Content)
	//找到** **中间的文字
	var middleStr = ""
	if len(match) > 2 {
		middleStr = match[1]
		fmt.Println("中间的是:::" + middleStr)
	}
	//过滤自己的
	if message.Author.ID == session.State.User.ID {
		return
	}

	//拿到哪一个用户的信息
	userId, _ := ChannelStateMap.Load(message.ChannelID)
	//如果这个是空的那就是调用的imagine
	if message.MessageReference == nil {
		if len(message.Attachments) != 0 {

			s, err := json.Marshal(message.Application)
			if err != nil {
				fmt.Println("json 解析失败")
				var modError = new(model.Error)
				modError.Message = "mj 返回消息 json 解析失败" + err.Error()
				DB.Model(&model.Error{}).Create(&modError)
			}
			fmt.Println(string(s))

			var mod model.Message
			go func() {
				for _, j := range message.Attachments {
					if userId != nil {
						mod.UserId = userId.(int)
					} else {
						mod.UserId = 0
					}
					allImg := until.UploadFile(j.URL, mod.UserId)
					mod.Url, err = until.ImgToFour(allImg, mod.UserId) //消息url
					if err != nil {
						var modError = new(model.Error)
						modError.Message = "mj 分割图片:::" + err.Error()
						DB.Model(&model.Error{}).Create(&modError)
					}
					mod.ResponseMessage = message.Content //消息体
					mod.Type = 1
					DB.Model(&model.Message{}).Create(&mod)

					DisChanTime.Store(message.ChannelID, true)
					DisChanUser.Delete(message.ChannelID)
				}
			}()

		}

	}
}

// // describe有两次请求
// // 第一次没有用 只拿到 message.id就可以
// if message.Interaction != nil && message.Interaction.Name == "describe" && message.Content == "" && len(message.Attachments) == 0 {
// 	if val, ok := isRequest[message.ChannelID]; ok && val {
// 		return
// 	}
// 	isRequest[message.ChannelID] = true
// 	go func() {

// 		var describeTwo = model.DescribeTwo{
// 			Type:          3,
// 			GuildID:       message.GuildID,
// 			ChannelID:     message.ChannelID,
// 			MessageFlags:  0,
// 			MessageID:     message.ID,
// 			ApplicationID: "*",
// 			SessionID:     "*",
// 			Data: model.Data{
// 				ComponentType: 2,
// 				CustomID:      "MJ::Job::PicReader::all",
// 			},
// 		}
// 		for {
// 			time.Sleep(time.Second)
// 			msg, _ := until.NewRequest("https://discord.com/api/v9/interactions", describeTwo, "*")
// 			if msg == "" {
// 				break
// 			}
// 		}

// 	}()
// }
// // if message.Content == "" || !strings.Contains(message.Content, "(fast)") {
// // 	return
// // }
// //如果describe真正的消息来了
// // if message.Interaction != nil {
// // 	fmt.Println("不是 interaction 空")
// // }
// // if message.Interaction != nil && message.Interaction.Name == "describe" {
// // 	fmt.Println("describe有东西")
// // }
// // if message.Content != "" {
// // 	fmt.Println("不是 Content 空")
// // }
// // if len(message.Attachments) != 0 {
// // 	fmt.Println("不是 interaction 空")
// // }

// if message.MessageReference != nil {
// 	go func() {
// 		var mod model.Message
// 		for _, j := range message.Attachments {
// 			mod.UserId = userId.(int)
// 			mod.Url = j.URL                       //消息url
// 			mod.ResponseMessage = message.Content //消息体
// 			mod.Type = 2
// 			DB.Model(&model.Message{}).Create(&mod)
// 		}
// 		//频道的并发数+1
// 		if value, ok := MjNumberMap.Load(session.Token[4:]); ok {
// 			val := value.(int)
// 			val++
// 			MjNumberMap.Store(session.Token[4:], val)
// 		}
// 		//还需要获取的次数-1
// 		if value, ok := DiscribeNumber.Load(message.ChannelID); ok {
// 			val := value.(int)
// 			if val > 0 {
// 				val--
// 				DiscribeNumber.Store(message.ChannelID, val)
// 			}
// 			if val == 0 {
// 				ChannelStateMap.Store(message.ChannelID, 1)
// 				isRequest[message.ChannelID] = false
// 			}
// 		}
// 	}()

// }
