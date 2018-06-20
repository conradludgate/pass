package libpass

import (
	"crypto/rand"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/sign"
)

func (kf *KeyFile) Sign(data []byte) (b []byte) {
	return sign.Sign(nil, data, &kf.Client.SignPriv)
}

func (kf *KeyFile) Open(b []byte) (data []byte, ok bool) {
	return sign.Open(nil, b, &kf.Phone.SignPub)
}

func (kf *KeyFile) Seal(data []byte) (b []byte) {
	if kf.precomputed {
		kf.precomputed = true
		box.Precompute(&kf.sharedKey, &kf.Phone.BoxPub, &kf.Client.BoxPriv)
	}

	var nonce [24]byte
	n, err := rand.Read(nonce[:])
	if err != nil {
		panic(err)
	}
	if n != 24 {
		panic("crypto/rand did not read 24 bytes into nonce")
	}

	return box.SealAfterPrecomputation(nonce[:], data, &nonce, &kf.sharedKey)
}

func (kf *KeyFile) Unseal(b []byte) (data []byte, ok bool) {
	if kf.precomputed {
		kf.precomputed = true
		box.Precompute(&kf.sharedKey, &kf.Phone.BoxPub, &kf.Client.BoxPriv)
	}

	var nonce [24]byte
	copy(nonce[:], b)

	return box.OpenAfterPrecomputation(nil, b[24:], &nonce, &kf.sharedKey)
}
