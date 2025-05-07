package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

type LockRequest resources.LockRequest

func NewLockRequest(r *http.Request) (resources.LockRequest, error) {
	request := struct {
		Data resources.LockRequest
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to decode request")
	}

	if err := LockRequest(request.Data).Validate(); err != nil {
		return request.Data, errors.Wrap(err, "invalid request")
	}

	return request.Data, nil
}

func (r LockRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TokenId, validation.Required),
		validation.Field(&r.ChainFrom, validation.Required),
		validation.Field(&r.ChainTo, validation.Required),
		validation.Field(&r.Sender, validation.Required, validation.By(isHexAddress)),
		validation.Field(&r.Receiver, validation.Required, validation.By(isHexAddress)),

		validation.Field(&r.NftId),
	)
}
