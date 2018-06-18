package libpass

import (
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	//"github.com/vmihailenco/msgpack"
	msgpack "encoding/json"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/sign"
)

type Keys struct {
	SignPriv *[64]byte
	SignPub  *[32]byte
	BoxPriv  *[32]byte
	BoxPub   *[32]byte
}

func (k *Keys) Generate() (err error) {
	if k == nil {
		return
	}

	if k.SignPriv == nil {
		k.SignPub, k.SignPriv, err = sign.GenerateKey(rand.Reader)
		if err != nil {
			return
		}
	}

	if k.SignPub == nil {
		k.SignPub = new([32]byte)
		copy((*k.SignPub)[:], (*k.SignPriv)[32:])
	}

	if k.BoxPriv == nil {
		k.BoxPub, k.BoxPriv, err = box.GenerateKey(rand.Reader)
		if err != nil {
			return
		}
	}

	if k.BoxPub == nil {
		curve25519.ScalarBaseMult(k.BoxPub, k.BoxPriv)
	}

	return nil
}

type KeyFile struct {
	Client, Phone *Keys
	file          *os.File
}

func NewKeyFile(file *os.File) (kf KeyFile, err error) {
	kf.file = file
	kf.Client = new(Keys)
	kf.Phone = new(Keys)

	err = kf.Load()

	// If key is invalid, remake
	if err != nil {
		// Create new key

		err = kf.Client.Generate()
		if err != nil {
			return kf, fmt.Errorf("Error generating key. %s", err.Error())
		}

		// Write KeyFile to file
		file.Truncate(0)
		file.Seek(0, 0)

		err = kf.Save()
		if err != nil {
			return kf, fmt.Errorf("Error writing to file. %s", err.Error())
		}
	}

	return
}

func (kf KeyFile) Load() error {
	// Read contents of file
	b, err := ioutil.ReadAll(kf.file)
	if err != nil {
		return err
	}

	if len(b) < 64 {
		return errors.New("Bad file format")
	}

	var checksum [64]byte
	copy(checksum[:], b[:64])
	if sum := sha512.Sum512(b[64:]); sum != checksum {
		return errors.New("Could not verify file integrity")
	}

	return msgpack.Unmarshal(b[64:], &kf)
}

func (kf KeyFile) Save() error {
	b, err := msgpack.Marshal(kf)
	if err != nil {
		return err
	}

	sum := sha512.Sum512(b)
	b = append(sum[:], b...)

	_, err = kf.file.Write(b)
	return err
}
