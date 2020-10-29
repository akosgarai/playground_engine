package ui

import (
	"github.com/akosgarai/playground_engine/pkg/model"
)

type UIModel struct {
	*model.BaseModel
}

func New() *UIModel {
	return &UIModel{
		model.New(),
	}
}
