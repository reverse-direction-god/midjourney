package model

import "gorm.io/gorm"

type GPTSqlInfo struct {
	gorm.Model
	Content string `json:"content"`
	UserId  uint   `json:"userId"`
}
type GPTMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type GPTReqMod struct {
	Model    string       `json:"model"`
	Messages []GPTMessage `json:"messages"`
}

type GPTResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Message      GPTMessage  `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}
