package libpass

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/sign"
)

var Zeros [32]byte = [32]byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

type Keys struct {
	SignPriv [64]byte
	SignPub  [32]byte
	BoxPriv  [32]byte
	BoxPub   [32]byte
}

func (k *Keys) Generate() (err error) {
	p, s, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	k.SignPub, k.SignPriv = *p, *s

	p, t, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return
	}

	k.BoxPub, k.BoxPriv = *p, *t
	return nil
}
