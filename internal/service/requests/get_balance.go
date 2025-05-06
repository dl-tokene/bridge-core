package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetBalanceRequest struct {
	TokenId string  `url:"-"`
	Address string  `url:"address"`
	Chain   string  `url:"chain"`
	Nft     *string `url:"nft"`
}

func NewGetBalanceRequest(r *http.Request) (GetBalanceRequest, error) {
	var req GetBalanceRequest
	if err := urlval.DecodeSilently(r.URL.Query(), &req); err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req.TokenId = chi.URLParam(r, "id")

	if err := req.Validate(); err != nil {
		return req, errors.Wrap(err, "invalid request")
	}

	return req, nil
}

func (r GetBalanceRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TokenId, validation.Required, is.Int),
		validation.Field(&r.Address, validation.Required, validation.By(isHexAddress)),
		validation.Field(&r.Chain, validation.Required, is.Int),
		validation.Field(&r.Nft, is.Int),
	)
}
