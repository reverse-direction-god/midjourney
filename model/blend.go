package model

type Attachments struct {
	ID               string `json:"id"`
	Filename         string `json:"filename"`
	UploadedFilename string `json:"uploaded_filename"`
}
type Blend struct {
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
				Choices              []struct {
					Name          string `json:"name"`
					Value         string `json:"value"`
					NameLocalized string `json:"name_localized"`
				} `json:"choices,omitempty"`
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
