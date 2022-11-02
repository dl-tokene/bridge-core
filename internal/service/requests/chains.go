package requests

import (
	"gitlab.com/distributed_lab/urlval"
	"gitlab.com/tokend/bridge/core/internal/data"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

type ChainsRequest struct {
	FilterType    []resources.ChainType `filter:"chain_type"`
	IncludeTokens bool                  `include:"token"`
}

type ChainsResponse struct { //todo delete this struct
	Data     data.ChainsQList
	Included resources.Included
}

func NewChainsRequest(r *http.Request) (ChainsRequest, error) {
	request := ChainsRequest{}
	err := urlval.DecodeSilently(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
