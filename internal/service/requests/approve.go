package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

type ApproveRequest resources.ApproveRequest

func NewApproveRequest(r *http.Request) (resources.ApproveRequest, error) {
	request := struct {
		Data resources.ApproveRequest
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to decode request")
	}

	if err := ApproveRequest(request.Data).Validate(); err != nil {
		return request.Data, errors.Wrap(err, "invalid request")
	}

	return request.Data, nil
}

func (r ApproveRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Address, validation.Required, validation.By(isHexAddress)),
		validation.Field(&r.ChainId, validation.Required, is.Int),
		validation.Field(&r.TokenId, validation.Required, is.Int),
	)
}
