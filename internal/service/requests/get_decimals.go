package requests

import (
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/distributed_lab/urlval"
	"net/http"
)

type GetDecimalsRequest struct {
	TokenId string `url:"-"`
	Chain   string `url:"chain"`
}

func NewGetDecimalsRequest(r *http.Request) (GetDecimalsRequest, error) {
	var req GetDecimalsRequest
	if err := urlval.DecodeSilently(r.URL.Query(), &req); err != nil {
		return req, errors.Wrap(err, "failed to decode request")
	}

	req.TokenId = chi.URLParam(r, "id")

	if err := req.Validate(); err != nil {
		return req, errors.Wrap(err, "invalid request")
	}

	return req, nil
}

func (r GetDecimalsRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TokenId, validation.Required),
		validation.Field(&r.Chain, validation.Required),
	)
}
