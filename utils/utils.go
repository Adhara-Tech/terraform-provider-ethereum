package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

const (
	keyHeaderKDF = "scrypt"
	scryptR      = 8
	scryptDKLen  = 32
)

func Keccak256(data ...[]byte) []byte {
	var b []byte
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	b = d.Sum(b)
	return b
}

func AES128(key, inText, iv []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(aesBlock, iv)
	outText := make([]byte, len(inText))
	stream.XORKeyStream(outText, inText)
	return outText, err
}

func PublicKeyToAddress(p ecdsa.PublicKey) []byte {
	k := elliptic.Marshal(elliptic.P256(), p.X, p.Y)
	return Keccak256(k[1:])[12:]
}

func EncryptJSONV3(data, auth []byte, scryptN, scryptP int) (CryptoJSON, error) {

	salt := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return CryptoJSON{}, fmt.Errorf("crypto/rand failed generating salt: %s", err.Error())
	}

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return CryptoJSON{}, fmt.Errorf("crypto/rand failed generating iv: %s", err.Error())
	}

	derivedKey, err := scrypt.Key(auth, salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return CryptoJSON{}, fmt.Errorf("derived key generation fail: %s", err.Error())
	}

	encryptKey := derivedKey[:16]
	cipherText, err := AES128(encryptKey, data, iv)
	if err != nil {
		return CryptoJSON{}, fmt.Errorf("aes encryption fail: %s", err.Error())
	}

	mac := Keccak256(derivedKey[16:32], cipherText)

	scryptParamsJSON := make(map[string]interface{}, 5)
	scryptParamsJSON["n"] = scryptN
	scryptParamsJSON["r"] = scryptR
	scryptParamsJSON["p"] = scryptP
	scryptParamsJSON["dklen"] = scryptDKLen
	scryptParamsJSON["salt"] = hex.EncodeToString(salt)
	cipherParamsJSON := CipherparamsJSON{
		IV: hex.EncodeToString(iv),
	}

	cryptoJSON := CryptoJSON{
		Cipher:       "aes-128-ctr",
		CipherText:   hex.EncodeToString(cipherText),
		CipherParams: cipherParamsJSON,
		KDF:          keyHeaderKDF,
		KDFParams:    scryptParamsJSON,
		MAC:          hex.EncodeToString(mac),
	}
	return cryptoJSON, nil
}
