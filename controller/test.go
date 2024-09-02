package controller

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type Test struct{}

func (s *Test) Test(c *gin.Context) {
	var session = new(discordgo.Session)
	ss := `{
		"type": 2,
		"application_id": "936929561302675456",
		"guild_id": "1214874136522784778",
		"channel_id": "1214874136522784781",
		"session_id": "62efd6694334a74d038cda63aebe9355",
		"data": {
			"version": "1166847114203123795",
			"id": "938956540159881230",
			"name": "imagine",
			"type": 1,
			"options": [
				{
					"type": 3,
					"name": "prompt",
					"value": "a cat"
				}
			],
			"application_command": {
				"id": "938956540159881230",
				"type": 1,
				"application_id": "936929561302675456",
				"version": "1166847114203123795",
				"name": "imagine",
				"description": "Create images with Midjourney",
				"options": [
					{
						"type": 3,
						"name": "prompt",
						"description": "The prompt to imagine",
						"required": true,
						"description_localized": "The prompt to imagine",
						"name_localized": "prompt"
					}
				],
				"integration_types": [
					0
				],
				"global_popularity_rank": 1,
				"description_localized": "Create images with Midjourney",
				"name_localized": "imagine"
			},
			"attachments": []
		},
		"nonce": "1215683049480519680",
		"analytics_location": "slash_ui"
	}`
	mp := make(map[interface{}]interface{})
	json.Unmarshal([]byte(ss), &mp)
	by, _ := json.Marshal(mp)
	message, err := session.ChannelMessageSend("*", string(by))
	if err != nil {
		panic(err)
	}
	fmt.Println(message.Content)
}
