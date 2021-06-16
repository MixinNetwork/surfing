package mobile

import (
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/btcsuite/btcutil/base58"
)

func LocalGenerateAddress(publicSpend, publicView string) (string, error) {
	publicSpendKey, err := crypto.KeyFromString(string(publicSpend))
	if err != nil {
		return "", err
	}
	publicViewKey, err := crypto.KeyFromString(string(publicView))
	if err != nil {
		return "", err
	}
	data := append([]byte(common.MainNetworkId), publicSpendKey[:]...)
	data = append(data, publicViewKey[:]...)
	checksum := crypto.NewHash(data)
	data = append(publicSpendKey[:], publicViewKey[:]...)
	data = append(data, checksum[:4]...)
	return common.MainNetworkId + base58.Encode(data), nil
}
