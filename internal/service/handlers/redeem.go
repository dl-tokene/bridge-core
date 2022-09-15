package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/internal/data"
	"gitlab.com/tokend/bridge/core/internal/proxy/types"
	"gitlab.com/tokend/bridge/core/internal/service/requests"
	"gitlab.com/tokend/bridge/core/resources"
	"net/http"
)

func Redeem(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewRedeemRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	token, err := TokensQ(r).FilterByID(req.TokenId).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get token")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if token == nil {
		Log(r).Debug("token not found")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/token_id": errors.New("token that you have sent does not exist"),
		})...)
		return
	}

	sourceChain, err := TokenChainsQ(r).
		FilterByTokenID(req.TokenId).
		FilterByChainID(req.ChainFrom).
		Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get token chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if sourceChain == nil {
		Log(r).WithError(err).Debug("token chain not found")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/chain_from": errors.New("token that you have sent does not connected to this chain"),
		})...)
		return
	}

	switch token.Type {
	case resources.FUNGIBLE:
		redeemFungibleToken(w, r, req, *sourceChain)
	case resources.NON_FUNGIBLE:
		redeemNonFungibleToken(w, r, req, *sourceChain)
	default:
		Log(r).WithError(err).
			WithField("token", token.ID).
			WithField("type", token.Type).
			Error("token type not supported")
		ape.RenderErr(w, problems.InternalError())
	}
}

func redeemFungibleToken(w http.ResponseWriter, r *http.Request, req resources.RedeemRequest, sourceChain data.TokenChain) {
	event, err := ProxyRepo(r).Get(sourceChain.ChainID).CheckFungibleLockEvent(req.TxHash, *req.EventIndex, sourceChain)
	if err != nil {
		renderCheckEventError(w, r, err)
		return
	}

	ape.Render(w, event)

	//destChain, err := TokenChainsQ(r).
	//	FilterByTokenID(req.TokenId).
	//	FilterByChainID(event.DestinationChain).
	//	Get()
	//if err != nil {
	//	Log(r).WithError(err).Error("failed to get token chain")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if destChain == nil {
	//	Log(r).WithError(err).Debug("token chain not found")
	//	ape.RenderErr(w, problems.BadRequest(validation.Errors{
	//		"data/chain_to": errors.New("token that you have sent does not connected to this chain"),
	//	})...)
	//	return
	//}
}

func redeemNonFungibleToken(w http.ResponseWriter, r *http.Request, req resources.RedeemRequest, sourceChain data.TokenChain) {
	event, err := ProxyRepo(r).Get(sourceChain.ChainID).CheckNonFungibleLockEvent(req.TxHash, *req.EventIndex, sourceChain)
	if err != nil {
		renderCheckEventError(w, r, err)
		return
	}

	ape.Render(w, event)
}

func renderCheckEventError(w http.ResponseWriter, r *http.Request, err error) {
	if err == types.TxNotConfirmed {
		Log(r).WithError(err).Debug("tx not confirmed")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/tx_hash": errors.New("tx not confirmed"),
		})...)
		return
	}
	if err == types.EventNotFound || err == types.WrongLockEvent || err == types.TxFailed {
		Log(r).WithError(err).Debug("invalid transaction")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data/tx_hash": errors.New("invalid transaction"),
		})...)
		return
	}
	Log(r).WithError(err).Error("failed to check fungible lock event")
	ape.RenderErr(w, problems.InternalError())
	return
}