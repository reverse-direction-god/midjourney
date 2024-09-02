package model

type optionImagine struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type appImagineOptions struct {
	Type                 int    `json:"type"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	Required             bool   `json:"required"`
	DescriptionLocalized string `json:"description_localized"`
	NameLocalized        string `json:"name_localized"`
}
type applicationCommandImagine struct {
	ID                   string              `json:"id"`
	Type                 int                 `json:"type"`
	ApplicationID        string              `json:"application_id"`
	Version              string              `json:"version"`
	Name                 string              `json:"name"`
	Description          string              `json:"description"`
	Options              []appImagineOptions `json:"options"`
	IntegrationTypes     []int               `json:"integration_types"`
	GlobalPopularityRank int                 `json:"global_popularity_rank"`
	DescriptionLocalized string              `json:"description_localized"`
	NameLocalized        string              `json:"name_localized"`
}

type dataImagine struct {
	Version            string                    `json:"version"`
	ID                 string                    `json:"id"`
	Name               string                    `json:"name"`
	Type               int                       `json:"type"`
	Options            []optionImagine           `json:"options"`
	ApplicationCommand applicationCommandImagine `json:"application_command"`
	Attachments        []interface{}             `json:"attachments"`
}

type Imagine struct {
	Type              int         `json:"type"`
	ApplicationID     string      `json:"application_id"`
	GuildId           string      `json:"guild_id"`
	ChannelID         string      `json:"channel_id"`
	SessionID         string      `json:"session_id"`
	Data              dataImagine `json:"data"`
	AnalyticsLocation string      `json:"analytics_location"`
}
