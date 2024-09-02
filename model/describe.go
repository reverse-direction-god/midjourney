package model

type Describe struct {
	Type          int    `json:"type"`
	ApplicationID string `json:"application_id"`
	GuildID       string `json:"guild_id"`
	ChannelID     string `json:"channel_id"`
	SessionID     string `json:"session_id"`
	Data          struct {
		Version string `json:"version"`
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    int    `json:"type"`
		Options []struct {
			Type  int    `json:"type"`
			Name  string `json:"name"`
			Value int    `json:"value"`
		} `json:"options"`
		ApplicationCommand struct {
			ID            string `json:"id"`
			Type          int    `json:"type"`
			ApplicationID string `json:"application_id"`
			Version       string `json:"version"`
			Name          string `json:"name"`
			Description   string `json:"description"`
			Options       []struct {
				Type                 int    `json:"type"`
				Name                 string `json:"name"`
				Description          string `json:"description"`
				Required             bool   `json:"required"`
				DescriptionLocalized string `json:"description_localized"`
				NameLocalized        string `json:"name_localized"`
			} `json:"options"`
			IntegrationTypes     []int  `json:"integration_types"`
			GlobalPopularityRank int    `json:"global_popularity_rank"`
			DescriptionLocalized string `json:"description_localized"`
			NameLocalized        string `json:"name_localized"`
		} `json:"application_command"`
		Attachments []Attachments `json:"attachments"`
	} `json:"data"`

	AnalyticsLocation string `json:"analytics_location"`
}
type Data struct {
	ComponentType int    `json:"component_type"`
	CustomID      string `json:"custom_id"`
}
type DescribeTwo struct {
	Type          int    `json:"type"`
	Nonce         string `json:"nonce"`
	GuildID       string `json:"guild_id"`
	ChannelID     string `json:"channel_id"`
	MessageFlags  int    `json:"message_flags"`
	MessageID     string `json:"message_id"`
	ApplicationID string `json:"application_id"`
	SessionID     string `json:"session_id"`
	Data          Data   `json:"data"`
}

type Queue struct {
	UserId int    `json:"userId"`
	Type   string `json:"type"` //接口类型
	My     bool   `json:"my"`   //是否自有账号
	Prompt struct {
		Content     string `json:"content"`
		SceneValues []struct {
			Text      string `json:"text"`
			SceneType string `json:"sceneType"`
		} `json:"sceneValues"`
		MjContent struct {
			Prefix []string `json:"prefix"`
			Suffix []string `json:"suffix"`
		} `json:"mjContent"`
		Size  string `json:"size"`
		Model string `json:"model"`
		Mode  string `json:"mode"`
		Cref  string `json:"cref"`
	} `json:"prompt"` //imagine需要
	FileName       string `json:"fileName"`       //describe需要
	UploadFileName string `json:"uploadFileName"` //describe需要
}
