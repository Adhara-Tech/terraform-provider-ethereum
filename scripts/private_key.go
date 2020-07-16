package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

type EncryptedKeyV3 struct {
	Address string     `json:"address"`
	Crypto  CryptoJSON `json:"crypto"`
	Id      string     `json:"id"`
	Version int        `json:"version"`
}

type CryptoJSON struct {
	Cipher       string                 `json:"cipher"`
	CipherText   string                 `json:"ciphertext"`
	CipherParams CipherparamsJSON       `json:"cipherparams"`
	KDF          string                 `json:"kdf"`
	KDFParams    map[string]interface{} `json:"kdfparams"`
	MAC          string                 `json:"mac"`
}

type CipherparamsJSON struct {
	IV string `json:"iv"`
}

func keccak256(data ...[]byte) ([]byte, error) {
	var b []byte
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, err := d.Write(b)
		if err != nil {
			return nil, err
		}
	}
	b = d.Sum(b)
	return b, nil
}

func aes128(key, inText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func getKDFKey(cryptoJSON CryptoJSON, auth string) ([]byte, error) {
	salt, err := hex.DecodeString(cryptoJSON.KDFParams["salt"].(string))
	if err != nil {
		return nil, err
	}

	dkLen := int(cryptoJSON.KDFParams["dklen"].(float64))
	n := int(cryptoJSON.KDFParams["n"].(float64))
	r := int(cryptoJSON.KDFParams["r"].(float64))
	p := int(cryptoJSON.KDFParams["p"].(float64))

	return scrypt.Key([]byte(auth), salt, n, r, p, dkLen)

}

func decryptKey(cryptoJson CryptoJSON, auth string) ([]byte, error) {

	if cryptoJson.Cipher != "aes-128-ctr" {
		return nil, fmt.Errorf("not supported cipher: %v", cryptoJson.Cipher)
	}

	mac, err := hex.DecodeString(cryptoJson.MAC)
	if err != nil {
		return nil, err
	}

	iv, err := hex.DecodeString(cryptoJson.CipherParams.IV)
	if err != nil {
		return nil, err
	}

	cipherText, err := hex.DecodeString(cryptoJson.CipherText)
	if err != nil {
		return nil, err
	}

	derivedKey, err := getKDFKey(cryptoJson, auth)
	if err != nil {
		return nil, err
	}

	calculatedMAC, _ := keccak256(derivedKey[16:32], cipherText)
	if !bytes.Equal(calculatedMAC, mac) {
		return nil, errors.New("invalid password")
	}

	plainText, err := aes128(derivedKey[:16], cipherText, iv)
	if err != nil {
		return nil, err
	}
	return plainText, err
}

func main() {

	file := os.Args[1]
	passphrase := os.Args[2]
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	var keystore EncryptedKeyV3
	err = json.Unmarshal(data, &keystore)
	if err != nil {
		log.Printf("%s", err.Error())
	}
	key, err := decryptKey(keystore.Crypto, passphrase)
	if err != nil {
		log.Printf("%s", err.Error())
	}

	fmt.Println(hex.EncodeToString(key))

}
