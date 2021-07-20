package mobile

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

func createTxWithAccount(raw signerInput, account common.Address) (*common.SignedTransaction, error) {
	tx := common.NewTransaction(raw.Asset)
	for _, in := range raw.Inputs {
		tx.AddInput(in.Hash, in.Index)
	}

	for _, out := range raw.Outputs {
		if out.Type != common.OutputTypeScript {
			return nil, fmt.Errorf("invalid output type %d", out.Type)
		}

		if out.Accounts != nil {
			tx.AddRandomScriptOutput(out.Accounts, out.Script, out.Amount)
		}
		if out.Keys != nil {
			tx.Outputs = append(tx.Outputs, &common.Output{
				Type:   common.OutputTypeScript,
				Amount: out.Amount,
				Keys:   out.Keys,
				Script: common.NewThresholdScript(1),
				Mask:   out.Mask,
			})
		}
	}

	extra, err := hex.DecodeString(raw.Extra)
	if err != nil {
		return nil, err
	}
	tx.Extra = extra

	signed := &common.SignedTransaction{Transaction: *tx}
	for i := range signed.Inputs {
		signed, err = signInputWithAccount(signed, raw, i, account)
		if err != nil {
			return nil, err
		}
	}
	return signed, nil
}

func signInputWithAccount(signed *common.SignedTransaction, reader common.UTXOKeysReader, index int, acc common.Address) (*common.SignedTransaction, error) {
	msg := signed.AsLatestVersion().PayloadMarshal()

	if index >= len(signed.Inputs) {
		return nil, fmt.Errorf("invalid input index %d/%d", index, len(signed.Inputs))
	}
	in := signed.Inputs[index]
	utxo, err := reader.ReadUTXOKeys(in.Hash, in.Index)
	if err != nil {
		return nil, err
	}
	if utxo == nil {
		return nil, fmt.Errorf("input not found %s:%d", in.Hash.String(), in.Index)
	}
	if len(utxo.Keys) != 1 {
		return nil, fmt.Errorf("utxo keys found %d", len(utxo.Keys))
	}
	sigs := make(map[uint16]*crypto.Signature)
	priv := crypto.DeriveGhostPrivateKey(&utxo.Mask, &acc.PrivateViewKey, &acc.PrivateSpendKey, uint64(in.Index))
	sig := priv.Sign(msg)
	sigs[0] = &sig
	signed.SignaturesMap = append(signed.SignaturesMap, sigs)
	return signed, nil
}

func CreateTransactionWithAccount(node string, account common.Address, rawStr string) (string, error) {
	var raw signerInput
	err := json.Unmarshal([]byte(rawStr), &raw)
	if err != nil {
		return "", err
	}
	raw.Node = node
	tx, err := createTxWithAccount(raw, account)
	if err != nil {
		return "", err
	}
	d := &common.VersionedTransaction{SignedTransaction: *tx}
	return hex.EncodeToString(d.Marshal()), nil
}
