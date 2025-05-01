package models

import (
	"gitlab.com/tokend/bridge/core/internal/data"
	"gitlab.com/tokend/bridge/core/resources"
)

type PermitResponse struct {
	Data     interface{}        `json:"data"`
	Included resources.Included `json:"included"`
}

func NewPermitResponse(data interface{}, chain data.Chain) PermitResponse {
	response := PermitResponse{
		Data:     data,
		Included: resources.Included{},
	}

	response.Included.Add(newChainModel(chain))

	return response
}
