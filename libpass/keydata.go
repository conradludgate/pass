package libpass

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vmihailenco/msgpack"
)

type KeyFile struct {
	Client Keys
	Phone  struct {
		Sign, Box [32]byte
	}
	file        *os.File
	sharedKey   [32]byte
	precomputed bool
}

func NewKeyFile(file *os.File) (kf *KeyFile, err error) {
	kf = new(KeyFile)
	kf.file = file

	err = kf.Load()

	// If key is invalid, remake
	if err != nil {
		// Create new key
		err = kf.Client.Generate()
		if err != nil {
			return kf, fmt.Errorf("Error generating key. %s", err.Error())
		}

		// Write KeyFile to file

		err = kf.Save()
		if err != nil {
			return kf, fmt.Errorf("Error writing to file. %s", err.Error())
		}
	}

	return
}

func (kf *KeyFile) Load() error {
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

	return msgpack.Unmarshal(b[64:], kf)
}

func (kf *KeyFile) Save() error {
	b, err := msgpack.Marshal(kf)
	if err != nil {
		return err
	}

	sum := sha512.Sum512(b)
	b = append(sum[:], b...)

	kf.file.Truncate(0)
	kf.file.Seek(0, 0)
	_, err = kf.file.Write(b)
	return err
}
