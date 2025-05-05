package requests

import (
	"github.com/go-chi/chi"
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

	return req, nil
}
