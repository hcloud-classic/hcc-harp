package rsautil

import "io/ioutil"

func WritePrivateAndPublicKeys() error {
	privKey, pubKey, err := GenerateKeyPair(4096)
	if err != nil {
		return err
	}

	privKeyData := PrivateKeyToBytes(privKey)

	pubKeyData, err := PublicKeyToBytes(pubKey)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("privkey.rsa", privKeyData, 0600)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("pubkey.rsa.pub", pubKeyData, 0644)
	if err != nil {
		return err
	}

	return nil
}
