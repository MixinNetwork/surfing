package mobile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	assert := assert.New(t)
	// address:	XINQTmRReDuPEUAVEyDyE2mBgxa1ojVRAvpYcKs5nSA7FDBBfAEeVRn8s9vAm3Cn1qzQ7JtjG62go4jSJU6yWyRUKHpamWAM
	// view key:	3f3683c539b95291253c364d766f83bf256210b30a72ff9aad22d7f417deff06
	// spend key:	8f337fb4a8b4ae6ed02b93473b879f25e574471dfd7be50248aae1b96927460c
	addr := "XINQTmRReDuPEUAVEyDyE2mBgxa1ojVRAvpYcKs5nSA7FDBBfAEeVRn8s9vAm3Cn1qzQ7JtjG62go4jSJU6yWyRUKHpamWAM"

	publicSpendKey := "b35480550384abc5933e86ff7770aa52f882ad06228acd6afa856c41d74ab60c"
	publicViewKey := "7b990ccd51697140fdc611b9f09ac8c97c6c7a283136e07a7cfa8ab162c82bc3"

	address, err := LocalGenerateAddress(publicSpendKey, publicViewKey)
	assert.Nil(err)
	assert.Equal(address, addr)
}
