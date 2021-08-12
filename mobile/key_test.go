package mobile

import (
	"fmt"
	"testing"

	"github.com/MixinNetwork/mixin/crypto"
	"github.com/stretchr/testify/assert"
)

func TestOneTimeKey(t *testing.T) {
	assert := assert.New(t)
	mask, err := crypto.KeyFromString("2d7ff76cd75825c53c6e12afc9cc87d6179456c098c3704b0e99515160c397df")
	assert.Nil(err)
	privateViewKey, err := crypto.KeyFromString("09DD51F16693A2D9AACED3EC75D41016A3BD9886BED14B3162AA28CAF4D65906")
	assert.Nil(err)
	privateSpendKey, err := crypto.KeyFromString("30DFF38C5916792A45A788A6736BE80478D91D7F70C833471F0D0AA46B563608")
	assert.Nil(err)
	index := 0

	priv := crypto.DeriveGhostPrivateKey(&mask, &privateViewKey, &privateSpendKey, uint64(index))
	fmt.Println(priv)
	assert.Equal("ef86daf5e7139d9ebe20dd3644179c6be4c4a27aca49449f9ad51fc40ca35e0a", priv.String())

	index = 99
	priv = crypto.DeriveGhostPrivateKey(&mask, &privateViewKey, &privateSpendKey, uint64(index))
	fmt.Println(priv)
	assert.Equal("329d5941a201c7626ae30d770a44c173be9992f964681e6f4593b7bf7b25250f", priv.String())
}
