package ethereum

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/Adhara-Tech/terraform-provider-ethereum/utils"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceEthereumKeystore() *schema.Resource {

	return &schema.Resource{
		Create: CreateAccount,
		Update: UpdateAccount,
		Delete: DeleteAccount,
		Read:   ReadAccount,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				ecdsaKey, passphrase, err := parseImportString(d.Id())
				if err != nil {
					return nil, err
				}
				if err := populateAccountFromImport(d, ecdsaKey, passphrase); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"scrypt_encryption": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "scrypt encryption difficulty ( lightkdf => 4MB and 100ms CPU approximately, standartkdf => 256MB and 1s CPU approximately",
				Default:     "lightkdf",
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "passphrase to encrypt the ethereum account",
				Default:     "",
			},
			"public_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"encrypted_json": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func getScryptValues(scryptEncryption string) (scriptN, scriptP int, err error) {

	switch scryptEncryption {

	case "lightkdf":
		scriptN = 4096
		scriptP = 6
		err = nil
		return
	case "standartkdf":
		scriptN = 262144
		scriptP = 1
		err = nil
		return
	default:
		return 0, 0, fmt.Errorf("invalid script encryption. Should be lightkdf or standartkdf")
	}
}

func CreateAccount(d *schema.ResourceData, meta interface{}) error {

	scryptEncryption := d.Get("scrypt_encryption").(string)
	passphrase := d.Get("passphrase").(string)

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("error generating ecdsa key: %s", err)
	}

	scriptN, scriptP, err := getScryptValues(scryptEncryption)
	if err != nil {
		return err
	}

	crytoJSON, err := utils.EncryptJSONV3(k.D.Bytes(), []byte(passphrase), scriptN, scriptP)
	if err != nil {
		return err
	}

	privateKey := hex.EncodeToString(k.D.Bytes())
	publicAddress := hex.EncodeToString(utils.PublicKeyToAddress(k.PublicKey))
	uuid := uuid.New()

	keystoreJSONV3 := utils.EncryptedKeyV3{
		Address: publicAddress,
		Crypto:  crytoJSON,
		Id:      uuid.String(),
		Version: 3,
	}

	encryptedJSON, err := json.Marshal(keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error parsing EncryptedKeyV3 json")
	}

	d.SetId(publicAddress)
	d.Set("private_key", privateKey)
	d.Set("public_address", publicAddress)
	d.Set("encrypted_json", string(encryptedJSON))

	return nil
}

func UpdateAccount(d *schema.ResourceData, meta interface{}) error {

	scryptEncryption := d.Get("scrypt_encryption").(string)
	passphrase := d.Get("passphrase").(string)

	scriptN, scriptP, err := getScryptValues(scryptEncryption)
	if err != nil {
		return err
	}

	privateKeyBytes, err := hex.DecodeString(d.Get("private_key").(string))
	if err != nil {
		return fmt.Errorf("error decoding private key")
	}

	var keystoreJSONV3 utils.EncryptedKeyV3
	err = json.Unmarshal([]byte(d.Get("encrypted_json").(string)), &keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error unmarshaling the encryptedKeyV3")
	}

	crytoJSON, err := utils.EncryptJSONV3(privateKeyBytes, []byte(passphrase), scriptN, scriptP)
	if err != nil {
		return err
	}

	keystoreJSONV3.Crypto = crytoJSON

	encryptedJSON, err := json.Marshal(keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error parsing EncryptedKeyV3 json")
	}

	d.Set("encrypted_json", string(encryptedJSON))

	return nil
}

func DeleteAccount(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}

func ReadAccount(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func parseImportString(importStr string) (*ecdsa.PrivateKey, string, error) {
	// example with passphrase -> ab4457b3ead76c5a7d1b947e7cbf805cbccdda3a5e72a9da02415b0055a5e06e_yourPassprase (passphrase can't contain _)
	// example with empty passfrase -> ab4457b3ead76c5a7d1b947e7cbf805cbccdda3a5e72a9da02415b0055a5e06e_

	importParts := strings.Split(strings.ToLower(importStr), "_")

	if len(importParts) < 2 {
		return nil, "", fmt.Errorf("to few import arguments")
	}

	k := new(ecdsa.PrivateKey)
	k.PublicKey.Curve = elliptic.P256()

	privateKey, err := hex.DecodeString(importParts[0])
	if err != nil {
		return nil, "", fmt.Errorf("error decoding private key")
	}
	if 8*len(privateKey) != k.Params().BitSize {
		return nil, "", fmt.Errorf("invalid private key length, need %d bits", k.Params().BitSize)
	}

	k.D = new(big.Int).SetBytes(privateKey)
	if k.D.Sign() <= 0 {
		return nil, "", fmt.Errorf("invalid private key")
	}
	k.PublicKey.X, k.PublicKey.Y = k.PublicKey.Curve.ScalarBaseMult(privateKey)

	passphrase := importParts[1]

	return k, passphrase, nil
}

func populateAccountFromImport(d *schema.ResourceData, k *ecdsa.PrivateKey, passphrase string) error {

	scriptN, scriptP, err := getScryptValues("lightkdf")
	if err != nil {
		return err
	}

	crytoJSON, err := utils.EncryptJSONV3(k.D.Bytes(), []byte(passphrase), scriptN, scriptP)
	if err != nil {
		return err
	}

	privateKey := hex.EncodeToString(k.D.Bytes())
	publicAddress := hex.EncodeToString(utils.PublicKeyToAddress(k.PublicKey))
	uuid := uuid.New()

	keystoreJSONV3 := utils.EncryptedKeyV3{
		Address: publicAddress,
		Crypto:  crytoJSON,
		Id:      uuid.String(),
		Version: 3,
	}

	encryptedJSON, err := json.Marshal(keystoreJSONV3)
	if err != nil {
		return fmt.Errorf("error parsing EncryptedKeyV3 json")
	}

	d.SetId(publicAddress)
	d.Set("scrypt_encryption", "lightkdf")
	d.Set("passphrase", passphrase)
	d.Set("private_key", privateKey)
	d.Set("public_address", publicAddress)
	d.Set("encrypted_json", string(encryptedJSON))

	return nil
}
