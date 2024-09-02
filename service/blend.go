package service

import (
	"mj/model"
	"mj/until"
)

func (s *discord) Blend(mod []model.Attachments) (string, error) {
	abs := BlendRequestModel
	for _, j := range mod {
		abs.Data.Attachments = append(abs.Data.Attachments, j)
	}

	return until.NewRequest("https://discord.com/api/v9/interactions", abs, "*.*.*")

}
