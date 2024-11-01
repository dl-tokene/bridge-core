package handlers

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"gitlab.com/tokend/bridge/core/internal/proxy/evm"
	"gitlab.com/tokend/bridge/core/internal/service/models"
	"gitlab.com/tokend/bridge/core/internal/service/requests"
	"net/http"
)

func Approve(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewApproveRequest(r)
	if err != nil {
		Log(r).WithError(err).Debug("failed to decode request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	tokenChain, err := TokenChainsQ(r).
		FilterByTokenID(req.TokenId).
		FilterByChainID(req.ChainId).
		Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get token chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if tokenChain == nil {
		Log(r).WithError(err).Debug("token chain not found")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data": errors.New("token that you have sent does not connected to this network or does not exist"),
		})...)
		return
	}
	if tokenChain.TokenType != evm.TokenTypeErc20 && req.Amount != nil {
		Log(r).WithError(err).Debug("unexpectedly received an amount field for non-ERC20 ")
		ape.RenderErr(w, problems.BadRequest(validation.Errors{
			"data": errors.New("amount is intended for ERC20 tokens"),
		})...)
		return
	}

	tx, err := ProxyRepo(r).Get(tokenChain.ChainID).Approve(*tokenChain, req.Address, req.Amount)
	if err != nil {
		Log(r).WithError(err).Error("failed to approve")
		ape.RenderErr(w, problems.InternalError())
		return
	}
	if tx == nil {
		// No approve is needed
		w.WriteHeader(http.StatusNoContent)
		return
	}

	chain, err := ChainsQ(r).FilterByID(tokenChain.ChainID).Get()
	if err != nil {
		Log(r).WithError(err).Error("failed to get chain")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, models.NewTxResponse(tx, *chain))
}
