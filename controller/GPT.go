package controller

import (
	"encoding/json"
	"fmt"
	"mj/model"
	"mj/service"

	"github.com/gin-gonic/gin"
)

type GPT struct{}

func (s *GPT) Tuili(c *gin.Context) {
	ser := service.GPT
	var mod model.GPTReqMod
	var msg struct {
		Content     string `json:"content"`
		Number      string `json:"number"`
		UserContent string `json:"userContent"`
	}
	c.ShouldBind(&msg)
	mod.Model = "gpt-3.5-turbo-0615"

	var message model.GPTMessage
	var system model.GPTMessage
	var system_System model.GPTMessage
	system.Role = "user"
	if msg.Number == "1" {
		system.Content = "现在你是一名基于输入描述的提示词生成器，你会将我输入的自然语言想象为完整的画面生成提示词。请注意，你生成后的内容服务于一个绘画AI，它只能理解具象的提示词而非抽象的概念。我将提供简短的中文描述，生成器需要为我提供准确的提示词，必要时优化和重组以提供更准确的内容，也只输出翻译后的英文内容。请模仿示例的结构生成完美的提示词。示例输入：“一条龙在一个女孩身后” 示例输出：(mysterious:1.3), ultra-realistic mix fantasy,(1 giant eastern dragon:1.3) behind an asian woman holding a glowing sword,void energy diamond sword, in the style of dark azure and light azure, mixes realistic and fantastical elements, vibrant manga, uhd image, glassy translucence, vibrant illustrations, ultra realistic, long hair, straight hair, light purple hair,head jewelly, jewelly, shawls,light In eyes, red eyes, portrait, firefly, bokeh, mysterious, fantasy, cloud, abstract, colorful background, night sky, flame, very detailed, high resolution, sharp, sharp image, 4k, 8k, masterpiece, best quality, magic effect, (high contrast:1.4), dream art, diamond, skin detail, face detail, eyes detail, mysterious colorful background, dark blue themes 请注意示例中的(__:1.x)表示给提示词增加权重，数值范围在0.6-1.5，数值越大权重越高。请仔细阅读我的要求，并严格按照规则生成提示词，请生成我需要的英文内容。接下来我发送文本的需要你生成"
	} else {
		system.Content = msg.UserContent
	}
	system_System.Role = "system"
	system_System.Content = "Please provide the Chinese description you'd like me to generate the prompt words for."

	message.Role = "user"
	message.Content = msg.Content
	mod.Messages = append(mod.Messages, system)
	mod.Messages = append(mod.Messages, system_System)
	mod.Messages = append(mod.Messages, message)

	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}
	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resContent,
	})
}
func (s *GPT) Run(c *gin.Context) {
	ser := service.GPT
	var msg struct {
		Content string `json:"content"`
	}
	// token := c.GetHeader("token")
	// id := until.Decrypt(token)
	var mod model.GPTReqMod
	// var res model.GPTResponse
	c.ShouldBind(&msg)
	mod.Model = "gpt-4"
	var message model.GPTMessage
	message.Role = "user"
	message.Content = msg.Content
	mod.Messages = append(mod.Messages, message)
	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}
	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resContent,
	})

}

// func (s *GPT) Write(c *gin.Context) {
// 	ser := service.GPT
// 	var mod model.GPTReqMod
// 	var message model.GPTMessage
// 	var msg struct {
// 		Content string `json:"content"`
// 	}
// 	c.ShouldBind(&msg)
// 	mod.Model = "gpt-3.5-turbo"
// 	message.Role = "user"
// 	message.Content = `"` + msg.Content + `"`+
// }

// 故事分镜
func (s *GPT) StoryStoryboard(c *gin.Context) {
	ser := service.GPT
	var msg struct {
		Content string `json:"content"`
	}
	// token := c.GetHeader("token")
	// id := until.Decrypt(token)
	var mod model.GPTReqMod
	// var res model.GPTResponse
	c.ShouldBind(&msg)
	var message model.GPTMessage
	mod.Model = "gpt-3.5-turbo"
	message.Role = "user"
	constString := `
	上面的 文字，我想你给我实现分镜分细一点
	用json的这种格式[
		{
			"lens_number": 1,  //情景编号
			"scene_description": "*****",//情景描述
			"detail_handling": "*****"//情景细节
		}]给我返回,说中文 除了镜头信息不要说其他的`
	message.Content = `"` + msg.Content + `"` + constString
	mod.Messages = append(mod.Messages, message)
	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}
	var resp = make([]map[string]string, 0)
	// json.Unmarshal([]byte(resContent)[7:len(resContent)-3], &resp)  //4.0的返回格式
	json.Unmarshal([]byte(resContent), &resp)
	fmt.Println(resContent)
	fmt.Println(resp)
	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resp,
	})

}

// 帮写剧情
func (s *GPT) WritePlot(c *gin.Context) {
	ser := service.GPT
	var msg struct {
		Content     string `json:"content"`
		Require     string `json:"require"`
		Description string `json:"description"`
	}
	var mod model.GPTReqMod
	var message model.GPTMessage
	c.ShouldBind(&msg)

	mod.Model = "gpt-3.5-turbo"
	message.Role = "user"
	var 原文 = ""
	if msg.Content != "" {
		原文 = `原文:"` + msg.Content + `"` + "\n"
	}
	发展剧情描述 := `情节描述:"` + msg.Description + `"` + "\n"
	要求 := `要求:"` + msg.Require + `"` + "\n"
	message.Content = 原文 + 发展剧情描述 + 要求 + "帮我写下面的剧情不要说与剧情主体无关的直接回复文章内容"
	mod.Messages = append(mod.Messages, message)
	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}

	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resContent,
	})

}

// 角色提取
func (s *GPT) CharacterExtraction(c *gin.Context) {
	ser := service.GPT
	var msg struct {
		Content string `json:"content"`
	}
	var mod model.GPTReqMod
	var message model.GPTMessage
	c.ShouldBind(&msg)
	mod.Model = "gpt-3.5-turbo"
	message.Role = "user"
	constString := `
	帮我提取文中的所有角色
	用json的这种格式[
		{
			"name": "****",  //角色名称
		}]给我返回,说中文 除了镜头信息不要说其他的`
	message.Content = `"` + msg.Content + `"` + constString
	mod.Messages = append(mod.Messages, message)
	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}
	var resp = make([]map[string]string, 0)
	// json.Unmarshal([]byte(resContent)[7:len(resContent)-3], &resp)  //4.0的返回格式
	json.Unmarshal([]byte(resContent), &resp)
	fmt.Println(resContent)
	fmt.Println(resp)
	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resp,
	})
}

func (s *GPT) CharacterAndStoryboard(c *gin.Context) {
	ser := service.GPT
	var msg struct {
		Content string `json:"content"`
		Max     string `json:"max"`
		Min     string `json:"min"`
	}
	var mod model.GPTReqMod
	var message model.GPTMessage
	c.ShouldBind(&msg)
	mod.Model = "gpt-3.5-turbo"
	message.Role = "user"
	message.Content = `原文:"` + msg.Content + `"` + "\n" + `帮我先分镜再查每个分镜段中有什么角色,每段最小字数限制` + msg.Min + "最大字数限制" + msg.Max + `.以这种格式返回[{"content":"***"//角色在原文的剧情 "roles":["***"","***"]//角色名称如果没有返回空}].`

	mod.Messages = append(mod.Messages, message)
	resContent := ser.GPTPOST(mod)
	if resContent == "" {
		c.JSON(500, model.Response{
			Msg:  "no",
			Code: 500,
			Data: nil,
		})
		return
	}
	var resp []struct {
		Content string   `json:"content"`
		Roles   []string `json:"roles"`
	}
	// json.Unmarshal([]byte(resContent)[7:len(resContent)-3], &resp)  //4.0的返回格式
	json.Unmarshal([]byte(resContent), &resp)

	c.JSON(200, model.Response{
		Msg:  "yes",
		Code: 200,
		Data: resp,
	})
}
