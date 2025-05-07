package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

var validTokenTypes = []any{
	string(resources.FUNGIBLE),
	string(resources.NON_FUNGIBLE),
}

var validTokenIncludes = []any{
	string(resources.CHAIN),
}

type RawTokensRequest struct {
	FilterType    []string `filter:"token_type"`
	IncludeChains string   `url:"-"`
}

type TokensRequest struct {
	FilterType    []resources.TokenType `filter:"token_type"`
	IncludeChains bool                  `include:"chain"`
}

func NewTokensRequest(r *http.Request) (TokensRequest, error) {
	rawRequest := RawTokensRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &rawRequest); err != nil {
		return TokensRequest{}, err
	}
	rawRequest.IncludeChains = r.URL.Query().Get("include")

	if err := rawRequest.Validate(); err != nil {
		return TokensRequest{}, errors.Wrap(err, "invalid request")
	}

	request := TokensRequest{}
	if err := urlval.DecodeSilently(r.URL.Query(), &request); err != nil {
		return request, err
	}

	return request, nil
}

func (r RawTokensRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.FilterType, validation.Each(validation.In(validTokenTypes...))),
		validation.Field(&r.IncludeChains, validation.In(validTokenIncludes...)),
	)
}
