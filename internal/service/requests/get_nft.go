package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetNftRequest struct {
	TokenId string `url:"-"`
	NftId   string `url:"-"`
	Chain   string `url:"chain"`
}

func NewGetNftRequest(r *http.Request) (GetNftRequest, error) {
	var req GetNftRequest
	if err := urlval.DecodeSilently(r.URL.Query(), &req); err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req.TokenId = chi.URLParam(r, "id")
	req.NftId = chi.URLParam(r, "nft_id")

	if err := req.Validate(); err != nil {
		return req, errors.Wrap(err, "invalid request")
	}

	return req, nil
}

func (r GetNftRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TokenId, validation.Required),
		validation.Field(&r.NftId, validation.Required),
		validation.Field(&r.Chain, validation.Required),
	)
}
