package libpass

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
)

func Verify(kf *KeyFile) (err error) {
	if kf.Client.BoxPriv == Zeros || kf.Phone.SignPub == Zeros {
		return errors.New(`Please pair first by running 'pass'`)
	}

	b := make([]byte, 128)
	copy(b[64:], kf.Phone.SignPub[:])
	copy(b[96:], kf.Phone.BoxPub[:])

	data := base64.RawStdEncoding.EncodeToString(kf.Seal(kf.Sign(b)))

	// Display QR Code
	fmt.Println("On the devices page, select this device and press the verify button, then scan this QR code")
	qrterminal.GenerateHalfBlock(data, qrterminal.L, os.Stdout)

	return nil
}
