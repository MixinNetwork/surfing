package mobile

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"

	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
	"github.com/stretchr/testify/assert"
)

var nodes []string

func init() {
	nodes = []string{
		"http://mixin-node-01.b1.run:8239",
		"http://mixin-node-02.b1.run:8239",
		"http://mixin-node0.exinpool.com:8239",
		"http://mixin-node1.exinpool.com:8239",
	}
}

func TestDeriveGhostPrivateKey(t *testing.T) {
	assert := assert.New(t)
	// {"amount":"100.00000000","hash":"b8414aaad80d095e7ab9e4529870c54e228f92d81dfded037c1c3b74ab25f6b2","index":0,"keys":["f457125dba0355c40d3e83df5aed733b2e8723ac2d94675669e034b9a97c112a"],"mask":"2d7ff76cd75825c53c6e12afc9cc87d6179456c098c3704b0e99515160c397df","script":"fffe01","type":0}
	mask, err := crypto.KeyFromString("2d7ff76cd75825c53c6e12afc9cc87d6179456c098c3704b0e99515160c397df")
	assert.Nil(err)
	privateSpend, err := crypto.KeyFromString("07b2d5ae306b8fc96d0b40e54b42592d63d786cd13bdbda4fda3eb958987d70b")
	assert.Nil(err)
	privateView, err := crypto.KeyFromString("a8e0fe81425bcad149f87ab082eb1442e0756269279722d44af7fe58ec19e70c")
	assert.Nil(err)
	priv := crypto.DeriveGhostPrivateKey(&mask, &privateView, &privateSpend, uint64(0))
	assert.Equal("8618e47480415ea548a8c7b1d0291d8edacbdc787762b164a89a60131c7d780a", priv.String())
}

func TextCreateTx(t *testing.T) {
	hash := "b8414aaad80d095e7ab9e4529870c54e228f92d81dfded037c1c3b74ab25f6b2"
	index := 0

	outputKeys := "aed95d85cfe8249aae8b260b7b1c48c483b30f81df57de18eceb8111d323b6e8"
	outputMask := "48f4f0dbe9f2060571889921c1823c0a70e0d62168ee79a82c3e4d306ddb86af"

	raw := fmt.Sprintf(`{"version":2,"asset": "b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","inputs":[{"hash":"%s","index":%d}],"outputs":[{"type":0,"amount":"100","script":"fffe01","keys":["%s"], "mask": "%s"}]}`, hash, index, outputKeys, outputMask)

	tx, err := CreateTransaction(nodes[rand.Intn(len(nodes))], raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}

func TestSignTx(t *testing.T) {
	hash := "b8414aaad80d095e7ab9e4529870c54e228f92d81dfded037c1c3b74ab25f6b2"
	index := 0

	outputKeys := "aed95d85cfe8249aae8b260b7b1c48c483b30f81df57de18eceb8111d323b6e8"
	outputMask := "48f4f0dbe9f2060571889921c1823c0a70e0d62168ee79a82c3e4d306ddb86af"

	raw := fmt.Sprintf(`{"version":2,"asset": "b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","inputs":[{"hash":"%s","index":%d}],"outputs":[{"type":0,"amount":"100","script":"fffe01","keys":["%s"], "mask": "%s"}]}`, hash, index, outputKeys, outputMask)

	// XINcEguDnBD9nSPMJeFVoTc2MeV3ta1iBvcGke3mC77XjQpcBHvH1xUnCEQ1pjhvrVijcPQKJ5jVsG6sSQjazwckYr9NTQn
	spend, err := hex.DecodeString("07b2d5ae306b8fc96d0b40e54b42592d63d786cd13bdbda4fda3eb958987d70b")
	if err != nil {
		return
	}
	view, err := hex.DecodeString("a8e0fe81425bcad149f87ab082eb1442e0756269279722d44af7fe58ec19e70c")
	if err != nil {
		return
	}
	var account common.Address
	copy(account.PrivateViewKey[:], view[:])
	copy(account.PrivateSpendKey[:], spend[:])

	tx, err := SignTransactionRaw(nodes[rand.Intn(len(nodes))], account, raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
	err = SentRawTransaction(nodes[rand.Intn(len(nodes))], tx)
	fmt.Println(err)
}
