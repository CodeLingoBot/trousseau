package openpgp

import (
	"bytes"
	"code.google.com/p/go.crypto/openpgp"
	"code.google.com/p/go.crypto/openpgp/armor"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

func Encrypt(encryptionKeys *openpgp.EntityList, s string) []byte {
	buf := &bytes.Buffer{}

	wa, err := armor.Encode(buf, "PGP MESSAGE", nil)
	if err != nil {
		log.Fatalf("Can't make armor: %v", err)
	}

	w, err := openpgp.Encrypt(wa, *encryptionKeys, nil, nil, nil)
	if err != nil {
		log.Fatalf("Error encrypting: %v", err)
	}
	_, err = io.Copy(w, strings.NewReader(s))
	if err != nil {
		log.Fatalf("Error encrypting: %v", err)
	}

	w.Close()
	wa.Close()

	return buf.Bytes()
}

func Decrypt(decryptionKeys *openpgp.EntityList, s, passphrase string) ([]byte, error) {
	if s == "" {
		return nil, nil
	}

	raw, err := armor.Decode(strings.NewReader(s))
	if err != nil {
		return nil, err
	}

	d, err := openpgp.ReadMessage(raw.Body, decryptionKeys,
		func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
			kp := []byte(passphrase)

			if symmetric {
				return kp, nil
			}

			for _, k := range keys {
				err := k.PrivateKey.Decrypt(kp)
				if err == nil {
					return nil, nil
				}
			}

			return nil, fmt.Errorf("Unable to decrypt trousseau data store. " +
				"Invalid passphrase supplied.")
		},
		nil)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(d.UnverifiedBody)
	return bytes, err
}
