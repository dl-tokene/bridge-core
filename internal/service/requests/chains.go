package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

var validChainTypes = []any{
	string(resources.EVM),
}

var validChainIncludes = []any{
	string(resources.TOKEN),
}

type RawChainsRequest struct {
	FilterType    []string `filter:"chain_type"`
	IncludeTokens string   `url:"-"`
}

type ChainsRequest struct {
	FilterType    []resources.ChainType `filter:"chain_type"`
	IncludeTokens bool                  `include:"token"`
}

func NewChainsRequest(r *http.Request) (ChainsRequest, error) {
	rawRequest := RawChainsRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &rawRequest); err != nil {
		return ChainsRequest{}, err
	}
	rawRequest.IncludeTokens = r.URL.Query().Get("include")

	if err := rawRequest.Validate(); err != nil {
		return ChainsRequest{}, errors.Wrap(err, "invalid request")
	}

	request := ChainsRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &request); err != nil {
		return request, err
	}

	return request, nil
}

func (r RawChainsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FilterType, validation.Each(validation.In(validChainTypes...))),
		validation.Field(&r.IncludeTokens, validation.In(validChainIncludes...)),
	)
}
