package pass

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"os"

	// "golang.org/x/crypto/curve25519"
	"github.com/mdp/qrterminal"
	"github.com/vmihailenco/msgpack"
)

type Pub rsa.PublicKey

func (p *Pub) EncodeMsgpack(encoder *msgpack.Encoder) error {
	return encoder.Encode(p.E, p.N.Bytes())
}

type Data struct {
	Hostname  string
	ID        []byte
	PublicKey Pub
}

func Pair(key *rsa.PrivateKey) (err error) {

	d := Data{
		"",
		make([]byte, 8),
		Pub(key.PublicKey),
	}

	if _, err := rand.Read(d.ID); err != nil {
		return err
	}

	d.Hostname, err = os.Hostname()
	if err != nil {
		return err
	}

	b, err := msgpack.Marshal(&d)
	if err != nil {
		return err
	}

	// buf := &bytes.Buffer{}

	// if err := binary.Write(buf, binary.LittleEndian, int64(key.E)); err != nil {
	// 	panic(err)
	// }

	// if _, err := buf.Write(key.N.Bytes()); err != nil {
	// 	panic(err)
	// }

	// data := base64.RawStdEncoding.EncodeToString(buf.Bytes())
	data := base64.RawStdEncoding.EncodeToString(b)
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	return nil
}
