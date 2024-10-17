/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

import "gitlab.com/tokend/bridge/core/internal/amount"

type ApproveRequest struct {
	Address string `json:"address"`
	// amount of token to approve
	Amount  *amount.Amount `json:"amount,omitempty"`
	ChainId string         `json:"chain_id"`
	TokenId string         `json:"token_id"`
}
