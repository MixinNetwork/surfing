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

func TestOffline2Transaction(t *testing.T) {
	assert := assert.New(t)

	spend, err := crypto.KeyFromString("1FEC3DAAE66C356F539DF4F779F49A9C0C919B962A5152D7D5BD418E016A3501")
	assert.Nil(err)
	fmt.Println("spend public key", spend.Public())
	view, err := crypto.KeyFromString("EDC8B7D77A4510872AEF0C8249D851A4217E006406E63C331E04137FCAD9F903")
	assert.Nil(err)
	fmt.Println("view public key", spend.Public())

	addr := &common.Address{
		PrivateSpendKey: spend,
		PrivateViewKey:  view,
		PublicSpendKey:  spend.Public(),
		PublicViewKey:   view.Public(),
	}
	address, err := common.NewAddressFromString("XINAWFi6YShoUsRE4KWffFZqRjQUtRjL2fKJLcRnRWUzg63pT2ASveUXo9BJwcTECfkCNS1R1JFNVRT73f7XRkteedC9jWPJ")
	assert.Nil(err)
	assert.Equal(addr.String(), address.String())

	/*
		./mixin -n node-42.f1ex.io:8239 getutxo -x 9600260ff99222012bd6fe4ee226e83ec42bafcb887e3dec64ff8d917abe4ecb
		{"amount":"100.00000000","hash":"89f0785ec04815218cef99c41bffae09e65d8c366bd7c603ea6e6164c46df236","index":0,"keys":["9b8eb6677a33805a876a9b2eae212455d9cb0ef2f23789169bfa5dec99942eb6"],"mask":"654c75f4a3609a2201dfda513cb970f54500992f8c48535f80a1445312ecf946","script":"fffe01","type":0}
		{"amount":"100.00000000","hash":"9b8e052f458ec8f81bf141f6cfa5edc291e917a4a1f25635b93f668e5948178d","index":0,"keys":["0eaa6a7eda65451880386b20130e1e9ebdfe7e3fb2661e803e4e2ab8b2d30de4"],"mask":"6a2934140f9f22929d37bed742610a057c278797bbe3e0ae9b6fc7e319a7f539","script":"fffe01","type":0}
	*/
	extra := hex.EncodeToString([]byte("to address xxxx:xx"))
	//68656c6c6f
	raw := fmt.Sprintf(`{"asset":"b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8",
		"extra":"%s",
		"inputs":[{"hash":"89f0785ec04815218cef99c41bffae09e65d8c366bd7c603ea6e6164c46df236","index":0}, 
			{"hash":"9b8e052f458ec8f81bf141f6cfa5edc291e917a4a1f25635b93f668e5948178d","index":0}],
		"outputs":[{"amount":"200","keys":["326ef472c94be57692c6a3e80e8370e37ee58eecea0d3b7304ab5192a44f5892"],"mask":"b249c11c090c002bcfdae157ac5163ce44cb5e21c59ac7a8a7d91c81761ec662","script":"fffe01","type":0}]}`, extra)
	fmt.Println(raw)
	tx, err := CreateTransaction(nodes[rand.Intn(len(nodes))], raw)
	assert.Nil(err)
	fmt.Println(tx)
	txBytes, err := hex.DecodeString(tx)
	assert.Nil(err)

	mask1, err := crypto.KeyFromString("654c75f4a3609a2201dfda513cb970f54500992f8c48535f80a1445312ecf946")
	assert.Nil(err)
	index := 0
	mask2, err := crypto.KeyFromString("6a2934140f9f22929d37bed742610a057c278797bbe3e0ae9b6fc7e319a7f539")
	assert.Nil(err)
	signature1 := Sign(txBytes, &view, &spend, &mask1, uint64(index))
	fmt.Println("Signature1: ", signature1)
	signature2 := Sign(txBytes, &view, &spend, &mask2, uint64(index))
	fmt.Println("Signature2: ", signature2)
	result, err := CreateTransactionWithSignature(nodes[rand.Intn(len(nodes))], raw, signature1.String()+","+signature2.String())
	assert.Nil(err)
	fmt.Println(result)
}

func TestOffline1Transaction(t *testing.T) {
	assert := assert.New(t)

	spend, err := crypto.KeyFromString("e612190f96fc058a52f84cabeff5eb6bbb436139f5c6920460b084eb9517210b")
	assert.Nil(err)
	fmt.Println("spend public key", spend.Public())
	view, err := crypto.KeyFromString("a85961a291b311ffbc88c14cf6d73a5aa6d8cda4ba5e3e37827c0d5e3fa83e03")
	assert.Nil(err)
	fmt.Println(spend)
	fmt.Println(view)
	/*
		./mixin -n node-42.f1ex.io:8239 getutxo -x 9600260ff99222012bd6fe4ee226e83ec42bafcb887e3dec64ff8d917abe4ecb
		{"amount":"100.00000000","hash":"9600260ff99222012bd6fe4ee226e83ec42bafcb887e3dec64ff8d917abe4ecb","index":0,"keys":["8245a19533b5a9cdf7a67d88e74b863150832d4c804ca1abf183d3c3ea7c6598"],"mask":"51e821847480618b767e92bd54b07cee4cf72b9c335f404ac8356c14898cff1b","script":"fffe01","type":0}
	*/
	extra := hex.EncodeToString([]byte("hello"))
	//68656c6c6f
	raw := fmt.Sprintf(`{"asset":"b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","extra":%s,"inputs":[{"hash":"9600260ff99222012bd6fe4ee226e83ec42bafcb887e3dec64ff8d917abe4ecb","index":0}],"outputs":[{"amount":"100.0","keys":["73755cc77391e5ba98a01f53dee962cd078b47da29ab6933041311fbd96213e6"],"mask":"6442dbb1b7b1335618f195b66425bb0de218167fd04d3f34251189833d99b777","script":"fffe01","type":0}]}`, extra)
	fmt.Println(raw)
	tx, err := CreateTransaction(nodes[rand.Intn(len(nodes))], raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
	txBytes, err := hex.DecodeString(tx)
	assert.Nil(err)

	mask, err := crypto.KeyFromString("51e821847480618b767e92bd54b07cee4cf72b9c335f404ac8356c14898cff1b")
	assert.Nil(err)
	index := 0
	signature := Sign(txBytes, &view, &spend, &mask, uint64(index))
	fmt.Println("Signature: ", signature)
	a, err := crypto.KeyFromString("f5716468bf003ae492c83be30c17095649e9e7eab057425384426d085791c60b")
	assert.Nil(err)
	fmt.Println(a.Public())
	result, err := CreateTransactionWithSignature(nodes[rand.Intn(len(nodes))], raw, signature.String())
	assert.Nil(err)
	fmt.Println(result)
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

func TestCreateTransaction(t *testing.T) {
	hash := "a22669508f2674edc6c4e6c76c7b6614704fdcc6e9814f76753422156bbc0522"
	index := 0

	outputKeys := "a855dab1e7a22345b24e47d1d8deb44edd5515ed3ba1b5766438f0bc562796b3"
	outputMask := "64957a17c09ea480150c68fad73bc7c54601c7905428918346d66aa8cf8538b5"

	raw := fmt.Sprintf(`{"version":2,"asset": "b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","inputs":[{"hash":"%s","index":%d}],"outputs":[{"type":0,"amount":"100","script":"fffe01","keys":["%s"], "mask": "%s"}]}`, hash, index, outputKeys, outputMask)

	tx, err := CreateTransaction(nodes[rand.Intn(len(nodes))], raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}

func TestCreateTx(t *testing.T) {
	hash := "a22669508f2674edc6c4e6c76c7b6614704fdcc6e9814f76753422156bbc0522"
	index := 0

	outputKeys := "a855dab1e7a22345b24e47d1d8deb44edd5515ed3ba1b5766438f0bc562796b3"
	outputMask := "64957a17c09ea480150c68fad73bc7c54601c7905428918346d66aa8cf8538b5"

	raw := fmt.Sprintf(`{"version":2,"asset": "b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","inputs":[{"hash":"%s","index":%d}],"outputs":[{"type":0,"amount":"100","script":"fffe01","keys":["%s"], "mask": "%s"}]}`, hash, index, outputKeys, outputMask)

	account := getAccount()
	tx, err := CreateTransactionWithAccount(nodes[rand.Intn(len(nodes))], *account, raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
	hash, err = SentRawTransaction(nodes[rand.Intn(len(nodes))], tx)
	fmt.Println(hash)
	fmt.Println(err)
}

func TestSignTx(t *testing.T) {
	hash := "b8414aaad80d095e7ab9e4529870c54e228f92d81dfded037c1c3b74ab25f6b2"
	index := 0

	outputKeys := "aed95d85cfe8249aae8b260b7b1c48c483b30f81df57de18eceb8111d323b6e8"
	outputMask := "48f4f0dbe9f2060571889921c1823c0a70e0d62168ee79a82c3e4d306ddb86af"

	raw := fmt.Sprintf(`{"version":2,"asset": "b9f49cf777dc4d03bc54cd1367eebca319f8603ea1ce18910d09e2c540c630d8","inputs":[{"hash":"%s","index":%d}],"outputs":[{"type":0,"amount":"100","script":"fffe01","keys":["%s"], "mask": "%s"}]}`, hash, index, outputKeys, outputMask)

	account := getAccount()

	tx, err := SignTransactionRaw(nodes[rand.Intn(len(nodes))], *account, raw)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
	hash, err = SentRawTransaction(nodes[rand.Intn(len(nodes))], tx)
	fmt.Println(hash)
	fmt.Println(err)
}

func getAccount() *common.Address {
	// XINcEguDnBD9nSPMJeFVoTc2MeV3ta1iBvcGke3mC77XjQpcBHvH1xUnCEQ1pjhvrVijcPQKJ5jVsG6sSQjazwckYr9NTQn
	spend, err := crypto.KeyFromString("07b2d5ae306b8fc96d0b40e54b42592d63d786cd13bdbda4fda3eb958987d70b")
	if err != nil {
		return nil
	}
	view, err := crypto.KeyFromString("a8e0fe81425bcad149f87ab082eb1442e0756269279722d44af7fe58ec19e70c")
	if err != nil {
		return nil
	}
	return &common.Address{
		PrivateSpendKey: spend,
		PrivateViewKey:  view,
	}
}
