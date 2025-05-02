package models

import (
	"fmt"
	"gitlab.com/tokend/bridge/core/resources"
)

type DecimalsResponse struct {
	Data DecimalsResponseData `json:"data"`
}

type DecimalsResponseData struct {
	resources.Key
	Attributes DecimalsResponseAttributes `json:"attributes"`
}

type DecimalsResponseAttributes struct {
	Decimals int `json:"decimals"`
}

func NewDecimalsResponse(tokenId string, decimals int) DecimalsResponse {
	id := fmt.Sprintf("%s", tokenId)

	response := DecimalsResponse{
		Data: DecimalsResponseData{
			Key: resources.Key{
				ID:   id,
				Type: "decimals",
			},
			Attributes: DecimalsResponseAttributes{
				Decimals: decimals,
			},
		},
	}

	return response
}
