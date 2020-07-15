package utils

import (
	"encoding/json"

	"github.com/google/uuid"
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

func (c *EncryptedKeyV3) ToString() string {
	j, _ := json.Marshal(c)
	return string(j)
}

func (c *EncryptedKeyV3) Generate(crypto CryptoJSON, publicAddress string) {
	uuid := uuid.New()
	c.Address = publicAddress
	c.Crypto = crypto
	c.Id = uuid.String()
	c.Version = 3
}
