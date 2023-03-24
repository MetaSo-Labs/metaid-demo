package example

import (
	"encoding/hex"
	"github.com/mvc-labs/mvc-lib-go/keys/bec"
)

func getRandPublicKey() (string, error) {
	pk, err := bec.NewPrivateKey(bec.S256())
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(pk.PubKey().SerialiseCompressed()), nil
}

