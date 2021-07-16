package mobile

import (
	"encoding/hex"
	"encoding/json"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

type signerInput struct {
	Inputs []struct {
		Hash  crypto.Hash   `json:"hash"`
		Index int           `json:"index"`
		Keys  []*crypto.Key `json:"keys"`
		Mask  crypto.Key    `json:"mask"`
	} `json:"inputs"`
	Outputs []struct {
		Type     uint8             `json:"type"`
		Mask     crypto.Key        `json:"mask"`
		Keys     []*crypto.Key     `json:"keys"`
		Amount   common.Integer    `json:"amount"`
		Script   common.Script     `json:"script"`
		Accounts []*common.Address `json:"accounts"`
	}
	Asset crypto.Hash `json:"asset"`
	Extra string      `json:"extra"`
	Node  string      `json:"-"`
}

func CreateTransaction(rawStr string) (string, error) {
	var raw signerInput
	err := json.Unmarshal([]byte(rawStr), &raw)
	if err != nil {
		return "", err
	}

	tx := common.NewTransaction(raw.Asset)
	for _, in := range raw.Inputs {
		tx.AddInput(in.Hash, in.Index)
	}

	for _, out := range raw.Outputs {
		if out.Mask.HasValue() {
			tx.Outputs = append(tx.Outputs, &common.Output{
				Type:   out.Type,
				Amount: out.Amount,
				Keys:   out.Keys,
				Script: out.Script,
				Mask:   out.Mask,
			})
		}
	}

	extra, err := hex.DecodeString(raw.Extra)
	if err != nil {
		return "", err
	}
	tx.Extra = extra

	signed := tx.AsLatestVersion()
	msg := signed.AsLatestVersion().PayloadMarshal()
	return hex.EncodeToString(msg), nil
}
