package model

type SDPromptConfig struct {
	Type        int    `json:"type"`
	OldFileName string `json:"oldFileName"`
	ImageBytes  string `json:"imageBytes"`

	Prompt           string `json:"prompt"`
	NegativePrompt   string `json:"negative_prompt"`
	Steps            int    `json:"steps"`
	Width            int    `json:"width"`
	Height           int    `json:"height"`
	CfgScale         int    `json:"cfg_scale"`
	SamplerName      string `json:"sampler_name"`
	NIter            int    `json:"n_iter"`
	BatchSize        int    `json:"batch_size"`
	UserId           string `json:"userId"`
	OverrideSettings struct {
		Sd_model_checkpoint string `json:"sd_model_checkpoint"`
	} `json:"override_settings"`
}

type SDTxt2ImgResponse struct {
	Images     []string    `json:"images"`
	Parameters interface{} `json:"parameters"`
	Info       string      `json:"info"`
	Detail     interface{} `json:"detail"`
}
