package libpass

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/vmihailenco/msgpack"
	"golang.org/x/crypto/ed25519"
)

type KeyData struct {
	PrivateKey [32]byte
	PublicKey  [32]byte

	PhoneEdKey    [32]byte
	PhoneCurveKey [32]byte

	CheckSum [64]byte
}

func GetKeyFromFile(file *os.File) (key KeyData, err error) {
	key, err = ReadKeyData(file)

	// If key is invalid, remake
	if err != nil || !key.Verify() {
		// Create new key
		key, err = GenerateKeyData()
		if err != nil {
			return key, fmt.Errorf("Error generating key. %s", err.Error())
		}

		// Write keydata to file
		file.Truncate(0)
		file.Seek(0, 0)

		err = key.WriteKeyData(file)
		if err != nil {
			return key, fmt.Errorf("Error writing to file. %s", err.Error())
		}
	}

	return
}

func GenerateKeyData() (key KeyData, err error) {
	_, sec, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return key, err
	}

	copy(key.PrivateKey[:], sec)
	copy(key.PublicKey[:], sec[32:])

	return
}

func ReadKeyData(r io.Reader) (key KeyData, err error) {
	// Read contents of file
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return
	}

	err = msgpack.Unmarshal(b, &key)

	return
}

func (k KeyData) CalculateSum() [64]byte {
	b := make([]byte, 256)
	copy(b[:], k.PrivateKey[:])
	copy(b[32:], k.PublicKey[:])
	copy(b[64:], k.PhoneEdKey[:])
	copy(b[128:], k.PhoneCurveKey[:])

	return sha512.Sum512(b)
}

func (k KeyData) Verify() bool {
	return k.CheckSum == k.CalculateSum()
}

func (k KeyData) WriteKeyData(w io.Writer) error {
	k.CheckSum = k.CalculateSum()
	b, err := msgpack.Marshal(k)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}
