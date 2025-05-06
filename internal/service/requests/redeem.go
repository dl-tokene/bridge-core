package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
	"regexp"
)

type RedeemRequest resources.RedeemRequest

func NewRedeemRequest(r *http.Request) (resources.RedeemRequest, error) {
	request := struct {
		Data resources.RedeemRequest
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request.Data, errors.Wrap(err, "failed to decode request")
	}

	if request.Data.EventIndex == nil {
		i := 0
		request.Data.EventIndex = &i
	}

	if err := RedeemRequest(request.Data).Validate(); err != nil {
		return request.Data, errors.Wrap(err, "invalid request")
	}

	return request.Data, nil
}

func (r RedeemRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TokenId, validation.Required, is.Int),
		validation.Field(&r.ChainFrom, validation.Required, is.Int),
		validation.Field(&r.TxHash, validation.Required, validation.By(isHash)),

		validation.Field(&r.Sender, validation.By(isHexAddress)),
		validation.Field(&r.RawTxData, validation.Match(regexp.MustCompile("^0x[0-9a-fA-F]+$"))),
	)
}
